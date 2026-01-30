package handlers

import (
	"golang.org/x/net/context"

	"github.com/shatteredsilicon/ssm-managed/api"
	"github.com/shatteredsilicon/ssm-managed/services/metric"
	"github.com/shatteredsilicon/ssm-managed/utils/logger"
)

type MetricServer struct {
	Metric *metric.Service
}

func (s *MetricServer) GetAdvisedMetricSamples(ctx context.Context, req *api.GetAdvisedMetricSamplesRequest) (*api.GetAdvisedMetricSamplesResponse, error) {
	var labels map[string]string
	if req.Instance != "" {
		labels = map[string]string{"instance": req.Instance}
	}

	samples, err := s.Metric.GetRulesSamples(ctx, labels)
	if err != nil {
		logger.Get(ctx).Errorf("%+v", err)
		return nil, err
	}

	var resp api.GetAdvisedMetricSamplesResponse
	for _, sample := range samples {
		resp.Samples = append(resp.Samples, &api.MetricSample{
			Metric: sample.Metric.String(),
			Value:  float64(sample.Value),
		})
	}

	return &resp, nil
}

// check interfaces
var (
	_ api.MetricServer = (*MetricServer)(nil)
)
