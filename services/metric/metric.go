package metric

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	sAPI "github.com/shatteredsilicon/ssm-managed/api"
	"github.com/shatteredsilicon/ssm-managed/services/consul"
	"github.com/shatteredsilicon/ssm-managed/services/prometheus"
	"github.com/sirupsen/logrus"
)

const (
	metricServiceInterval = time.Minute

	// prometheus APIs seems to have a maxmium data points
	// limit of 11000, we set the batch size to 10000 just in case
	batchSizeInterval = 10000 * time.Second

	relabelKeyTpl = "%s/relabel"
)

// Service metric service
type Service struct {
	consul        *consul.Client
	logger        *logrus.Entry
	prometheusSvc *prometheus.Service
	prometheusAPI api.Client
}

// NewService returns a new metric service
func NewService(consulClient *consul.Client, prometheusSvc *prometheus.Service, prometheusAPI api.Client, logger *logrus.Entry) *Service {
	return &Service{
		consul:        consulClient,
		prometheusSvc: prometheusSvc,
		prometheusAPI: prometheusAPI,
		logger:        logger,
	}
}

// Run runs the metric service
func (s *Service) Run(ctx context.Context) {
	ticker := time.NewTicker(metricServiceInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.relabel(ctx)
		case <-ctx.Done():
			return
		}
	}
}

type relabelConsulValue struct {
	OldName   string `json:"old_name"`
	Start     int64  `json:"start"` // milliseconds
	End       int64  `json:"end"`   // milliseconds
	Relabeled bool   `json:"relabeled"`
}

func (s *Service) relabel(ctx context.Context) {
	l := s.logger.WithField("function", "relabel")

	nodes, err := s.consul.GetNodes()
	if err != nil {
		l.Errorf("Failed to retrieve nodes from consul: %s\n", err)
		return
	}

	v1api := v1.NewAPI(s.prometheusAPI)

FOR_EACH_NODE:
	for i := range nodes {
		node, err := s.consul.GetNode(nodes[i].Node)
		if err != nil {
			l.Errorf("Failed to retrieve catalog node from consul: %s\n", err)
			continue FOR_EACH_NODE
		}

		var relabelValues []relabelConsulValue
		relabelKey := fmt.Sprintf(relabelKeyTpl, node.Node.Node)
		relabelBytes, err := s.consul.GetExactKV(relabelKey)
		if err != nil {
			l.Errorf("Failed to get relabel KV from consul for key %s: %s", relabelKey, err)
			continue FOR_EACH_NODE
		}

		if relabelBytes == nil {
			// the KV is not set, meaning
			// no relabel process needed for this node
			continue FOR_EACH_NODE
		}
		if err = json.Unmarshal(relabelBytes, &relabelValues); err != nil {
			l.Errorf("Failed to decode relabel KV from consul for key %s: %s", relabelKey, err)
			continue FOR_EACH_NODE
		}

		for rvi, relabelValue := range relabelValues {
			if relabelValue.Relabeled || relabelValue.OldName == node.Node.Node {
				continue
			}

			oldName, startTime, endTime := relabelValue.OldName, time.Unix(0, relabelValue.Start*int64(time.Millisecond)), time.Unix(0, relabelValue.End*int64(time.Millisecond))

			series, _, err := v1api.Series(ctx, []string{fmt.Sprintf(`{instance="%s"}`, oldName)}, startTime, endTime)
			if err != nil {
				l.Errorf("Failed to get series list from propmetheus for instance %s in range [%s, %s]: %s\n", oldName, startTime, endTime, err)
				continue FOR_EACH_NODE
			}

			seriesChan := make(chan struct{}, 100)
			wg := sync.WaitGroup{}
			for _, se := range series {
				serieCopy := se
				seriesChan <- struct{}{}
				wg.Add(1)
				go func(serie model.LabelSet) {
					defer func() {
						<-seriesChan
						wg.Done()
					}()

					labels := make([]string, 0)
					reqLabels := make([]*sAPI.Label, 0)
					for k, v := range serie {
						labels = append(labels, fmt.Sprintf(`%s="%s"`, k, v))
						if k != "instance" {
							reqLabels = append(reqLabels, &sAPI.Label{
								Name:  string(k),
								Value: string(v),
							})
						}
					}
					reqLabels = append(reqLabels, &sAPI.Label{
						Name:  "instance",
						Value: node.Node.Node,
					})

					batchStartTime := startTime
					matchStr := fmt.Sprintf("{%s}", strings.Join(labels, ","))
					for batchStartTime.Before(endTime) {
						batchEndTime := batchStartTime.Add(batchSizeInterval)
						if batchEndTime.After(endTime) {
							batchEndTime = endTime
						}

						// fetch un-relabeled prometheus data in a certain
						// time range for processing
						value, _, err := v1api.QueryRange(ctx, matchStr, v1.Range{Start: batchStartTime, End: batchEndTime, Step: time.Second})
						if err != nil {
							l.Errorf("Failed to get query_range from propmetheus for match[] => %s in range [%s, %s]: %s\n", matchStr, batchStartTime, batchEndTime, err)
							return
						}

						writeRequest := sAPI.WriteRequest{
							Timeseries: make([]*sAPI.TimeSeries, 0),
						}

						switch value.Type() {
						case model.ValScalar:
							scalar := value.(*model.Scalar)
							writeRequest.Timeseries = append(writeRequest.Timeseries, &sAPI.TimeSeries{
								Labels:  reqLabels,
								Samples: []*sAPI.Sample{{Value: float64(scalar.Value), Timestamp: int64(scalar.Timestamp)}},
							})
						case model.ValMatrix:
							matrix := value.(model.Matrix)
							timeserie := &sAPI.TimeSeries{
								Labels:  reqLabels,
								Samples: make([]*sAPI.Sample, 0),
							}
							for matrixI := range matrix {
								for valueI := range matrix[matrixI].Values {
									timeserie.Samples = append(timeserie.Samples, &sAPI.Sample{
										Value:     float64(matrix[matrixI].Values[valueI].Value),
										Timestamp: int64(matrix[matrixI].Values[valueI].Timestamp),
									})
								}
							}
							writeRequest.Timeseries = append(writeRequest.Timeseries, timeserie)
						case model.ValVector:
							vector := value.(model.Vector)
							timeserie := &sAPI.TimeSeries{
								Labels:  reqLabels,
								Samples: make([]*sAPI.Sample, 0),
							}
							for vectorI := range vector {
								timeserie.Samples = append(timeserie.Samples, &sAPI.Sample{Value: float64(vector[vectorI].Value), Timestamp: int64(vector[vectorI].Timestamp)})
							}
							writeRequest.Timeseries = append(writeRequest.Timeseries, timeserie)
						}

						// call prometheus remote write API to write
						// relabeled data
						if len(writeRequest.Timeseries) > 0 {
							err = s.prometheusSvc.WriteMetrics(ctx, writeRequest)
							if err != nil {
								l.Errorf("Failed to write metrics to propmetheus for match[] => %s in range [%s, %s]: %s\n", matchStr, batchStartTime, batchEndTime, err)
								return
							}
						}

						batchStartTime = batchEndTime
					}
				}(serieCopy)
			}
			wg.Wait()
			close(seriesChan)

			relabelValues[rvi].Relabeled = true
			if relabelBytes, err = json.Marshal(relabelValues); err != nil {
				s.logger.Warnf("failed to unmarshal relabel consul KV (%s): %s\n", relabelKey, err.Error())
			} else if err = s.consul.PutExactKV(relabelKey, relabelBytes); err != nil {
				s.logger.Warnf("failed to put relabel consul KV (%s): %s\n", relabelKey, err.Error())
			}
		}
	}
}
