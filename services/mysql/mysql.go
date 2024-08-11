// pmm-managed
// Copyright (C) 2017 Percona LLC
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

// Package mysql contains business logic of working with Remote MySQL instances.
package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/go-sql-driver/mysql"
	servicelib "github.com/percona/kardianos-service"
	"github.com/pkg/errors"
	"github.com/prometheus/common/model"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"

	"github.com/shatteredsilicon/ssm-managed/models"
	"github.com/shatteredsilicon/ssm-managed/services"
	"github.com/shatteredsilicon/ssm-managed/services/consul"
	"github.com/shatteredsilicon/ssm-managed/services/prometheus"
	"github.com/shatteredsilicon/ssm-managed/services/qan"
	"github.com/shatteredsilicon/ssm-managed/utils/logger"
	"github.com/shatteredsilicon/ssm-managed/utils/ports"
	"github.com/shatteredsilicon/ssm/proto/config"
)

const (
	defaultMySQLPort uint32 = 3306

	// maximum time for connecting to the database and running all queries
	sqlCheckTimeout = 5 * time.Second
)

var versionRegexp = regexp.MustCompile(`([\d\.]+)-.*`)

type ServiceConfig struct {
	MySQLdExporterPath string

	Prometheus    *prometheus.Service
	Supervisor    services.Supervisor
	DB            *reform.DB
	PortsRegistry *ports.Registry
	QAN           *qan.Service
	Consul        *consul.Client
}

// Service is responsible for interactions with AWS RDS.
type Service struct {
	*ServiceConfig
	ssmServerNode *models.Node
}

// NewService creates a new service.
func NewService(config *ServiceConfig) (*Service, error) {
	var node models.Node
	err := config.DB.FindOneTo(&node, "type", models.SSMServerNodeType)
	if err != nil {
		return nil, err
	}

	for _, path := range []*string{
		&config.MySQLdExporterPath,
	} {
		if *path == "" {
			continue
		}
		p, err := exec.LookPath(*path)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		*path = p
	}

	svc := &Service{
		ServiceConfig: config,
		ssmServerNode: &node,
	}
	return svc, nil
}

type Instance struct {
	Node    models.RemoteNode
	Service models.MySQLService
}

func (svc *Service) ApplyPrometheusConfiguration(ctx context.Context, q *reform.Querier) error {
	mySQLHR := &prometheus.ScrapeConfig{
		JobName:        "remote-mysql-hr",
		ScrapeInterval: "1s",
		ScrapeTimeout:  "1s",
		MetricsPath:    "/metrics-hr",
		RelabelConfigs: []prometheus.RelabelConfig{{
			TargetLabel: "job",
			Replacement: "mysql",
		}},
		MetricRelabelConfigs: []prometheus.RelabelConfig{
			{
				SourceLabels: []model.LabelName{"__name__"},
				TargetLabel:  "__name__",
				Regex:        "(mysql_global_status_binlog_cache_use|mysql_global_status_binlog_cache_disk_use|mysql_global_status_binlog_stmt_cache_use|mysql_global_status_binlog_stmt_cache_disk_use|mysql_global_status_connections|mysql_global_status_queries|mysql_global_status_questions|mysql_global_status_created_tmp_tables|mysql_global_status_created_tmp_disk_tables|mysql_global_status_created_tmp_files|mysql_global_status_select_full_join|mysql_global_status_select_full_range_join|mysql_global_status_select_range|mysql_global_status_select_range_check|mysql_global_status_select_scan|mysql_global_status_sort_rows|mysql_global_status_sort_range|mysql_global_status_sort_merge_passes|mysql_global_status_sort_scan|mysql_global_status_slow_queries|mysql_global_status_aborted_connects|mysql_global_status_aborted_clients|mysql_global_status_table_locks_immediate|mysql_global_status_table_locks_waited|mysql_global_status_threads_created|mysql_global_status_bytes_received|mysql_global_status_bytes_sent|mysql_global_status_auroradb_commits|mysql_global_status_auroradb_commit_latency|mysql_global_status_auroradb_ddl_stmt_duration|mysql_global_status_auroradb_select_stmt_duration|mysql_global_status_auroradb_insert_stmt_duration|mysql_global_status_auroradb_update_stmt_duration|mysql_global_status_auroradb_delete_stmt_duration|mysql_global_status_qcache_hits|mysql_global_status_qcache_inserts|mysql_global_status_qcache_not_cached|mysql_global_status_qcache_lowmem_prunes|mysql_global_status_opened_files|mysql_global_status_opened_tables|mysql_global_status_table_open_cache_hits|mysql_global_status_table_open_cache_misses|mysql_global_status_opened_table_definitions|mysql_global_status_table_open_cache_overflows|mysql_global_status_auroradb_thread_deadlocks|mysql_global_status_aurora_missing_history_on_replica_incidents|mysql_global_status_aurora_reserved_mem_exceeded_incidents|mysql_global_status_key_reads|mysql_global_status_key_read_requests|mysql_global_status_key_writes|mysql_global_status_aria_pagecache_reads|mysql_global_status_aria_pagecache_writes|mysql_global_status_aria_transaction_log_syncs|mysql_global_status_aria_pagecache_read_requests|mysql_global_status_aria_pagecache_write_requests|mysql_global_status_key_write_requests|mysql_global_status_innodb_max_trx_id|mysql_global_status_innodb_log_writes|mysql_global_status_innodb_os_log_written|mysql_global_status_innodb_onlineddl_rowlog_pct_used|mysql_global_status_innodb_onlineddl_rowlog_rows|mysql_global_status_innodb_defragment_compression_failures|mysql_global_status_innodb_defragment_failures|mysql_global_status_innodb_row_lock_waits|mysql_global_status_innodb_row_lock_time|mysql_global_status_innodb_data_reads|mysql_global_status_innodb_data_writes|mysql_global_status_innodb_data_fsyncs|mysql_global_status_innodb_deadlocks|mysql_global_status_innodb_buffer_pool_read_requests|mysql_global_status_innodb_buffer_pool_write_requests|mysql_global_status_innodb_buffer_pool_reads|mysql_global_status_innodb_pages_created|mysql_global_status_innodb_pages_read|mysql_global_status_innodb_pages_written|mysql_global_status_innodb_buffer_pool_read_ahead|mysql_global_status_innodb_buffer_pool_read_ahead_rnd|mysql_global_status_innodb_buffer_pool_read_ahead_evicted|mysql_global_status_innodb_ibuf_merges|mysql_global_status_innodb_ibuf_merged_inserts|mysql_global_status_innodb_ibuf_merged_deletes|mysql_global_status_innodb_ibuf_merged_delete_marks|mysql_global_status_rocksdb_block_cache_hit|mysql_global_status_rocksdb_block_cache_miss|mysql_global_status_rocksdb_block_cache_add|mysql_global_status_rocksdb_block_cache_add_failures|mysql_global_status_rocksdb_block_cache_bytes_read|mysql_global_status_rocksdb_block_cache_bytes_write|mysql_global_status_rocksdb_block_cache_index_hit|mysql_global_status_rocksdb_block_cache_index_miss|mysql_global_status_rocksdb_block_cache_index_add|mysql_global_status_rocksdb_block_cache_index_bytes_insert|mysql_global_status_rocksdb_block_cache_index_bytes_evict|mysql_global_status_rocksdb_block_cache_filter_hit|mysql_global_status_rocksdb_block_cache_filter_miss|mysql_global_status_rocksdb_block_cache_filter_add|mysql_global_status_rocksdb_block_cache_filter_bytes_insert|mysql_global_status_rocksdb_block_cache_filter_bytes_evict|mysql_global_status_rocksdb_block_cache_data_bytes_insert|mysql_global_status_rocksdb_bloom_filter_useful|mysql_global_status_rocksdb_memtable_hit|mysql_global_status_rocksdb_memtable_miss|mysql_global_status_rocksdb_number_keys_read|mysql_global_status_rocksdb_number_keys_written|mysql_global_status_rocksdb_number_keys_updated|mysql_global_status_rocksdb_get_hit_l0|mysql_global_status_rocksdb_get_hit_l1|mysql_global_status_rocksdb_number_db_seek|mysql_global_status_rocksdb_number_db_seek_found|mysql_global_status_rocksdb_number_db_next|mysql_global_status_rocksdb_number_db_next_found|mysql_global_status_rocksdb_number_db_prev|mysql_global_status_rocksdb_number_db_prev_found|mysql_global_status_rocksdb_bytes_read|mysql_global_status_rocksdb_bytes_written|mysql_global_status_rocksdb_iter_bytes_read|mysql_global_status_rocksdb_write_self|mysql_global_status_rocksdb_write_other|mysql_global_status_rocksdb_write_timeout|mysql_global_status_rocksdb_write_wal|mysql_global_status_rocksdb_wal_synced|mysql_global_status_rocksdb_wal_bytes|mysql_global_status_rocksdb_number_reseeks_iteration|mysql_global_status_rocksdb_rows_inserted|mysql_global_status_rocksdb_rows_updated|mysql_global_status_rocksdb_rows_deleted|mysql_global_status_rocksdb_rows_read|mysql_global_status_rocksdb_rows_expired|mysql_global_status_rocksdb_system_rows_deleted|mysql_global_status_rocksdb_system_rows_inserted|mysql_global_status_rocksdb_system_rows_read|mysql_global_status_rocksdb_system_rows_updated|mysql_global_status_rocksdb_no_file_opens|mysql_global_status_rocksdb_no_file_closes|mysql_global_status_rocksdb_no_file_errors|mysql_global_status_rocksdb_stall_l0_file_count_limit_slowdowns|mysql_global_status_rocksdb_stall_l0_file_count_limit_stops|mysql_global_status_rocksdb_stall_locked_l0_file_count_limit_slowdowns|mysql_global_status_rocksdb_stall_locked_l0_file_count_limit_stops|mysql_global_status_rocksdb_stall_memtable_limit_slowdowns|mysql_global_status_rocksdb_stall_memtable_limit_stops|mysql_global_status_tokudb_txn_begin|mysql_global_status_tokudb_txn_begin_read_only|mysql_global_status_tokudb_txn_commits|mysql_global_status_tokudb_txn_aborts|mysql_global_status_tokudb_messages_injected_at_root|mysql_global_status_tokudb_broadcase_messages_injected_at_root|mysql_global_status_tokudb_messages_injected_at_root_bytes|mysql_global_status_tokudb_messages_flushed_from_h1_to_leaves_bytes|mysql_global_status_tokudb_leaf_compression_to_memory_seconds|mysql_global_status_tokudb_leaf_decompression_to_memory_seconds|mysql_global_status_tokudb_leaf_serialization_to_memory_seconds|mysql_global_status_tokudb_leaf_deserialization_to_memory_seconds|mysql_global_status_tokudb_nonleaf_compression_to_memory_seconds|mysql_global_status_tokudb_nonleaf_decompression_to_memory_seconds|mysql_global_status_tokudb_nonleaf_serialization_to_memory_seconds|mysql_global_status_tokudb_nonleaf_deserialization_to_memory_seconds|mysql_global_status_tokudb_promotion_injections_at_depth_0|mysql_global_status_tokudb_promotion_injections_at_depth_1|mysql_global_status_tokudb_promotion_injections_at_depth_2|mysql_global_status_tokudb_promotion_injections_at_depth_3|mysql_global_status_tokudb_promotion_injections_lower_than_depth_3|mysql_global_status_tokudb_promotion_stopped_nonempty_buffer|mysql_global_status_tokudb_promotion_stopped_at_height_1|mysql_global_status_tokudb_promotion_stopped_child_locked_or_not_in_memory|mysql_global_status_tokudb_promotion_stopped_child_not_fully_in_memory|mysql_global_status_tokudb_promotion_stopped_after_locking_child|mysql_global_status_tokudb_leaf_nodes_created|mysql_global_status_tokudb_nonleaf_nodes_created|mysql_global_status_tokudb_leaf_nodes_destroyed|mysql_global_status_tokudb_nonleaf_nodes_destroyed|mysql_global_status_tokudb_messages_ignored_by_leaf_due_to_msn|mysql_global_status_tokudb_cursor_skip_deleted_leaf_entry|mysql_global_status_tokudb_leaf_entry_apply_gc_bytes_out|mysql_global_status_tokudb_leaf_entry_apply_gc_bytes_in|mysql_global_status_tokudb_cachetable_evictions|mysql_global_status_tokudb_cachetable_miss|mysql_global_status_tokudb_cachetable_miss_time|mysql_global_status_tokudb_cachetable_prefetches|mysql_global_status_tokudb_leaf_node_partial_evictions|mysql_global_status_tokudb_nonleaf_node_partial_evictions|mysql_global_status_tokudb_leaf_node_partial_evictions_bytes|mysql_global_status_tokudb_nonleaf_node_partial_evictions_bytes|mysql_global_status_tokudb_leaf_node_full_evictions|mysql_global_status_tokudb_nonleaf_node_full_evictions|mysql_global_status_tokudb_leaf_node_full_evictions_bytes|mysql_global_status_tokudb_nonleaf_node_full_evictions_bytes|mysql_global_status_tokudb_checkpoint_taken|mysql_global_status_tokudb_checkpoint_duration|mysql_global_status_tokudb_checkpoint_begin_time|mysql_global_status_tokudb_checkpoint_end_time|mysql_global_status_tokudb_leaf_nodes_flushed_checkpoint_seconds|mysql_global_status_tokudb_nonleaf_nodes_flushed_to_disk_checkpoint_seconds|mysql_global_status_tokudb_leaf_nodes_flushed_not_checkpoint_seconds|mysql_global_status_tokudb_nonleaf_nodes_flushed_to_disk_not_checkpoint_seconds|mysql_global_status_tokudb_leaf_nodes_flushed_checkpoint_bytes|mysql_global_status_tokudb_nonleaf_nodes_flushed_to_disk_checkpoint_bytes|mysql_global_status_tokudb_leaf_nodes_flushed_checkpoint_uncompressed_bytes|mysql_global_status_tokudb_nonleaf_nodes_flushed_to_disk_checkpoint_uncompressed_by|mysql_global_status_tokudb_leaf_nodes_flushed_not_checkpoint_bytes|mysql_global_status_tokudb_nonleaf_nodes_flushed_to_disk_not_checkpoint_bytes|mysql_global_status_tokudb_pivots_fetched_for_query_bytes|mysql_global_status_tokudb_pivots_fetched_for_prefetch_bytes|mysql_global_status_tokudb_pivots_fetched_for_write_bytes|mysql_global_status_tokudb_pivots_fetched_for_query_seconds|mysql_global_status_tokudb_pivots_fetched_for_prefetch_seconds|mysql_global_status_tokudb_pivots_fetched_for_write_seconds|mysql_global_status_tokudb_basements_fetched_target_query_bytes|mysql_global_status_tokudb_basements_fetched_prefetch_bytes|mysql_global_status_tokudb_basements_fetched_for_write_bytes|mysql_global_status_tokudb_basements_fetched_prelocked_range_bytes|mysql_global_status_tokudb_basements_fetched_target_query_seconds|mysql_global_status_tokudb_basements_fetched_prefetch_seconds|mysql_global_status_tokudb_basements_fetched_for_write_seconds|mysql_global_status_tokudb_basements_fetched_prelocked_range_seconds|mysql_global_status_tokudb_buffers_fetched_target_query_bytes|mysql_global_status_tokudb_buffers_fetched_prefetch_bytes|mysql_global_status_tokudb_buffers_fetched_for_write_bytes|mysql_global_status_tokudb_buffers_fetched_prelocked_range_bytes|mysql_global_status_tokudb_buffers_fetched_target_query_seconds|mysql_global_status_tokudb_buffers_fetched_prefetch_seconds|mysql_global_status_tokudb_buffers_fetched_for_write_seconds|mysql_global_status_tokudb_buffers_fetched_prelocked_range_seconds|mysql_global_status_tokudb_logger_writes|mysql_global_status_tokudb_logger_writes_bytes|mysql_global_status_tokudb_filesystem_fsync_num|mysql_global_status_tokudb_logger_writes_seconds|mysql_global_status_tokudb_filesystem_fsync_time|mysql_global_status_tokudb_locktree_wait_time|mysql_global_status_tokudb_cachetable_cleaner_executions|mysql_global_status_tokudb_cachetable_pool_client_total_items_processed|mysql_global_status_tokudb_cachetable_pool_cachetable_total_items_processed|mysql_global_status_tokudb_cachetable_pool_checkpoint_total_items_processed|mysql_global_status_wsrep_flow_control_paused_ns|mysql_global_status_wsrep_local_cert_failures|mysql_global_status_wsrep_local_bf_aborts|mysql_global_status_wsrep_flow_control_sent|mysql_global_status_wsrep_received_bytes|mysql_global_status_wsrep_replicated_bytes|mysql_global_status_wsrep_received|mysql_global_status_wsrep_replicated|mysql_global_status_wsrep_flow_control_recv)",
				Replacement:  "${1}_total",
			},
		},
	}
	mySQLMR := &prometheus.ScrapeConfig{
		JobName:        "remote-mysql-mr",
		ScrapeInterval: "5s",
		ScrapeTimeout:  "1s",
		MetricsPath:    "/metrics-mr",
		RelabelConfigs: []prometheus.RelabelConfig{{
			TargetLabel: "job",
			Replacement: "mysql",
		}},
		MetricRelabelConfigs: []prometheus.RelabelConfig{
			{
				SourceLabels: []model.LabelName{"__name__"},
				TargetLabel:  "__name__",
				Regex:        "(mysql_slave_status_relay_log_space)",
				Replacement:  "${1}_total",
			},
		},
	}
	mySQLLR := &prometheus.ScrapeConfig{
		JobName:        "remote-mysql-lr",
		ScrapeInterval: "60s",
		ScrapeTimeout:  "5s",
		MetricsPath:    "/metrics-lr",
		RelabelConfigs: []prometheus.RelabelConfig{{
			TargetLabel: "job",
			Replacement: "mysql",
		}},
		MetricRelabelConfigs: []prometheus.RelabelConfig{
			{
				SourceLabels: []model.LabelName{"__name__"},
				TargetLabel:  "__name__",
				Regex:        "(mysql_info_schema_user_statistics_total_connections|mysql_global_variables_binlog_cache_size|mysql_global_variables_binlog_stmt_cache_size)",
				Replacement:  "${1}_count",
			},
		},
	}

	nodes, err := q.FindAllFrom(models.RemoteNodeTable, "type", models.RemoteNodeType, models.SSMServerNodeType)
	if err != nil {
		return errors.WithStack(err)
	}
	for _, n := range nodes {
		node := n.(*models.RemoteNode)

		mySQLServices, err := q.SelectAllFrom(models.MySQLServiceTable, "WHERE node_id = ? AND type = ?", node.ID, models.MySQLServiceType)
		if err != nil {
			return errors.WithStack(err)
		}
		if len(mySQLServices) == 0 {
			continue
		}
		service := mySQLServices[0].(*models.MySQLService)
		if service.Type != models.MySQLServiceType {
			continue
		}

		agents, err := models.AgentsForServiceID(q, service.ID)
		if err != nil {
			return err
		}
		for _, agent := range agents {
			switch agent.Type {
			case models.MySQLdExporterAgentType:
				a := models.MySQLdExporter{ID: agent.ID}
				if e := q.Reload(&a); e != nil {
					return errors.WithStack(e)
				}
				logger.Get(ctx).WithField("component", "mysql").Infof("%s %s %s %d", a.Type, node.Name, node.Region, *a.ListenPort)

				var labels []prometheus.LabelPair
				if node.Type == models.SSMServerNodeType {
					labels = []prometheus.LabelPair{
						{Name: "instance", Value: string(node.Type)}, // ssm-server node uses type as name
					}
				} else {
					labels = []prometheus.LabelPair{
						{Name: "instance", Value: node.Name},
						{Name: "region", Value: string(models.RemoteNodeRegion)},
					}
				}

				sc := prometheus.StaticConfig{
					Targets: []string{fmt.Sprintf("127.0.0.1:%d", *a.ListenPort)},
					Labels:  labels,
				}
				mySQLHR.StaticConfigs = append(mySQLHR.StaticConfigs, sc)
				mySQLMR.StaticConfigs = append(mySQLMR.StaticConfigs, sc)
				mySQLLR.StaticConfigs = append(mySQLLR.StaticConfigs, sc)
			}
		}
	}

	// sort by instance
	sorterFor := func(sc []prometheus.StaticConfig) func(int, int) bool {
		return func(i, j int) bool {
			return sc[i].Labels[0].Value < sc[j].Labels[0].Value
		}
	}
	sort.Slice(mySQLHR.StaticConfigs, sorterFor(mySQLHR.StaticConfigs))
	sort.Slice(mySQLMR.StaticConfigs, sorterFor(mySQLMR.StaticConfigs))
	sort.Slice(mySQLLR.StaticConfigs, sorterFor(mySQLLR.StaticConfigs))

	return svc.Prometheus.SetScrapeConfigs(ctx, false, mySQLHR, mySQLMR, mySQLLR)
}

func (svc *Service) List(ctx context.Context) ([]Instance, error) {
	var res []Instance
	err := svc.DB.InTransaction(func(tx *reform.TX) error {
		structs, e := tx.SelectAllFrom(models.RemoteNodeTable, "WHERE type = ? ORDER BY id", models.RemoteNodeType)
		if e != nil {
			return e
		}
		nodes := make([]models.RemoteNode, len(structs))
		for i, str := range structs {
			nodes[i] = *str.(*models.RemoteNode)
		}

		structs, e = tx.SelectAllFrom(models.MySQLServiceTable, "WHERE type = ? ORDER BY id", models.MySQLServiceType)
		if e != nil {
			return e
		}
		services := make([]models.MySQLService, len(structs))
		for i, str := range structs {
			services[i] = *str.(*models.MySQLService)
		}

		for _, node := range nodes {
			for _, service := range services {
				if node.ID == service.NodeID {
					res = append(res, Instance{
						Node:    node,
						Service: service,
					})
				}
			}
		}
		return nil
	})
	return res, err
}

func (svc *Service) addMySQLdExporter(ctx context.Context, tx *reform.TX, service *models.MySQLService, username, password string) error {
	// insert mysqld_exporter agent and association
	port, err := svc.PortsRegistry.Reserve()
	if err != nil {
		return err
	}
	agent := &models.MySQLdExporter{
		Type:         models.MySQLdExporterAgentType,
		RunsOnNodeID: svc.ssmServerNode.ID,

		ServiceUsername: &username,
		ServicePassword: &password,
		ListenPort:      &port,
	}
	if err = tx.Insert(agent); err != nil {
		return errors.WithStack(err)
	}
	if err = tx.Insert(&models.AgentService{AgentID: agent.ID, ServiceID: service.ID}); err != nil {
		return errors.WithStack(err)
	}

	// check connection and a number of tables
	var tableCount int
	dsn := agent.DSN(service)
	db, err := sql.Open("mysql", dsn)
	if err == nil {
		sqlCtx, cancel := context.WithTimeout(ctx, sqlCheckTimeout)
		err = db.QueryRowContext(sqlCtx, "SELECT COUNT(*) FROM information_schema.tables").Scan(&tableCount)
		cancel()
		db.Close()
		agent.MySQLDisableTablestats = pointer.ToBool(tableCount > 1000)
	}
	if err != nil {
		if err, ok := err.(*mysql.MySQLError); ok {
			switch err.Number {
			case 0x414: // 1044
				return status.Error(codes.PermissionDenied, err.Message)
			case 0x415: // 1045
				return status.Error(codes.Unauthenticated, err.Message)
			}
		}
		return errors.WithStack(err)
	}

	// start mysqld_exporter agent
	if svc.MySQLdExporterPath != "" {
		cfg := svc.mysqlExporterCfg(agent, dsn)
		if err = svc.Supervisor.Start(ctx, cfg); err != nil {
			return err
		}
	}

	return nil
}

func (svc *Service) mysqlExporterCfg(agent *models.MySQLdExporter, dsn string) *servicelib.Config {
	name := models.NameForSupervisor(agent.Type, *agent.ListenPort)

	arguments := []string{
		"-collect.binlog_size",
		"-collect.global_status",
		"-collect.global_variables",
		"-collect.info_schema.innodb_metrics",
		"-collect.info_schema.processlist",
		"-collect.info_schema.query_response_time",
		"-collect.info_schema.userstats",
		"-collect.perf_schema.eventswaits",
		"-collect.perf_schema.file_events",
		"-collect.slave_status",
		fmt.Sprintf("-web.listen-address=127.0.0.1:%d", *agent.ListenPort),
		"-web.auth-file=\"\"",
	}
	if agent.MySQLDisableTablestats == nil || !*agent.MySQLDisableTablestats {
		// enable tablestats and a few related collectors just like pmm-admin
		// https://github.com/shatteredsilicon/ssm-client/blob/e94b61ed0e5482a27039f0d1b0b36076731f0c29/pmm/plugin/mysql/metrics/metrics.go#L98-L105
		arguments = append(arguments, "-collect.auto_increment.columns")
		arguments = append(arguments, "-collect.info_schema.tables")
		arguments = append(arguments, "-collect.info_schema.tablestats")
		arguments = append(arguments, "-collect.perf_schema.indexiowaits")
		arguments = append(arguments, "-collect.perf_schema.tableiowaits")
		arguments = append(arguments, "-collect.perf_schema.tablelocks")
	}
	sort.Strings(arguments)

	return &servicelib.Config{
		Name:        name,
		DisplayName: name,
		Description: name,
		Executable:  svc.MySQLdExporterPath,
		Arguments:   arguments,
		Environment: []string{fmt.Sprintf("DATA_SOURCE_NAME=%s", dsn)},
	}
}

func (svc *Service) addQanAgent(
	ctx context.Context,
	tx *reform.TX,
	service *models.MySQLService,
	node *models.RemoteNode,
	username, password string,
	qanConfig *config.QAN,
) error {
	// Despite running a single qan-agent process on PMM Server, we use one database record per MySQL instance
	// to store username/password and UUID.

	// insert qan-agent agent and association
	agent := &models.QanAgent{
		Type:         models.QanAgentAgentType,
		RunsOnNodeID: svc.ssmServerNode.ID,

		ServiceUsername: &username,
		ServicePassword: &password,
		ListenPort:      pointer.ToUint16(models.QanAgentPort),
	}
	var err error
	if err = tx.Insert(agent); err != nil {
		return errors.WithStack(err)
	}
	if err = tx.Insert(&models.AgentService{AgentID: agent.ID, ServiceID: service.ID}); err != nil {
		return errors.WithStack(err)
	}

	// DSNs for mysqld_exporter and qan-agent are currently identical,
	// so we do not check connection again

	// start or reconfigure qan-agent
	if svc.QAN != nil {
		if qanConfig == nil {
			qanConfig = &config.QAN{CollectFrom: qan.PerfschemaCollectFrom}
		}

		nodeName := node.Name
		if node.Type == models.SSMServerNodeType {
			nodeName = string(node.Type) // ssm-server node uses type as name
		}
		if err = svc.QAN.AddMySQL(ctx, nodeName, service, agent, *qanConfig); err != nil {
			return err
		}

		// re-save agent with set QANDBInstanceUUID
		if err = tx.Save(agent); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (svc *Service) Add(
	ctx context.Context,
	name, address string,
	port uint32,
	username, password string,
	nodeType models.NodeType,
	region string,
	qanConfig *config.QAN,
) (int32, error) {
	address = strings.TrimSpace(address)
	username = strings.TrimSpace(username)
	name = strings.TrimSpace(name)
	if address == "" {
		return 0, status.Error(codes.InvalidArgument, "MySQL instance host is not given.")
	}
	if username == "" {
		return 0, status.Error(codes.InvalidArgument, "Username is not given.")
	}
	if port == 0 {
		port = defaultMySQLPort
	}
	if name == "" {
		name = address
	}

	// check if instance is added on client side
	added, err := svc.clientInstanceAdded(ctx, name)
	if err != nil {
		return 0, err
	}
	if added {
		return 0, status.Error(codes.AlreadyExists, fmt.Sprintf("MySQL instance with name %s is already added", name))
	}

	var id int32
	err = svc.DB.InTransaction(func(tx *reform.TX) error {
		// insert node
		node := &models.RemoteNode{
			Type:   nodeType,
			Name:   name,
			Region: region,
		}
		if err := tx.Insert(node); err != nil {
			if err, ok := err.(*mysql.MySQLError); !ok || err.Number != 0x426 {
				return errors.WithStack(err)
			}

			err = tx.SelectOneTo(node, "WHERE type = ? AND name = ? AND region = ?", nodeType, name, region)
			if err != nil {
				return errors.WithStack(err)
			}
		}
		id = node.ID

		engine, engineVersion, err := svc.EngineAndEngineVersion(ctx, address, port, username, password)
		if err != nil {
			return errors.WithStack(err)
		}

		// insert service
		service := &models.MySQLService{
			Type:   models.MySQLServiceType,
			NodeID: node.ID,

			Address:       &address,
			Port:          pointer.ToUint16(uint16(port)),
			Engine:        &engine,
			EngineVersion: &engineVersion,
		}
		if err = tx.Insert(service); err != nil {
			return errors.WithStack(err)
		}

		if err = svc.addMySQLdExporter(ctx, tx, service, username, password); err != nil {
			return err
		}
		if err = svc.addQanAgent(ctx, tx, service, node, username, password, qanConfig); err != nil {
			return err
		}

		return svc.ApplyPrometheusConfiguration(ctx, tx.Querier)
	})

	return id, err
}

func (svc *Service) clientInstanceAdded(ctx context.Context, name string) (bool, error) {
	node, err := svc.Consul.GetNode(name)
	if err != nil {
		logger.Get(ctx).Errorf("get consul services from node failed: %+v", err)
		return false, err
	}

	if node == nil {
		return false, nil
	}

	for _, service := range node.Services {
		t := models.AgentType(service.Service)
		if t == models.ClientMySQLdExporterAgentType || t == models.ClientMySQLQanAgentAgentType {
			// instance is added on client side
			return true, nil
		}
	}

	return false, nil
}

func (svc *Service) Remove(ctx context.Context, id int32) error {
	var err error
	return svc.DB.InTransaction(func(tx *reform.TX) error {
		var node models.RemoteNode
		if err = tx.SelectOneTo(
			&node,
			"WHERE ( type = ? OR type = ? ) AND id = ?",
			models.RemoteNodeType, models.SSMServerNodeType, id,
		); err != nil {
			if err == reform.ErrNoRows {
				return status.Errorf(codes.NotFound, "MySQL instance with ID %d not found.", id)
			}
			return errors.WithStack(err)
		}

		var service models.MySQLService
		if err = tx.SelectOneTo(&service, "WHERE node_id = ? and type = ?", node.ID, models.MySQLServiceType); err != nil {
			return errors.WithStack(err)
		}

		// remove associations of the service and agents
		agentsForService, err := models.AgentsForServiceID(tx.Querier, service.ID)
		if err != nil {
			return err
		}
		for _, agent := range agentsForService {
			var deleted uint
			deleted, err = tx.DeleteFrom(models.AgentServiceView, "WHERE service_id = ? AND agent_id = ?", service.ID, agent.ID)
			if err != nil {
				return errors.WithStack(err)
			}
			if deleted != 1 {
				return errors.Errorf("expected to delete 1 record, deleted %d", deleted)
			}
		}

		// remove associations of the node and agents
		agentsForNode, err := models.AgentsForNodeID(tx.Querier, node.ID)
		if err != nil {
			return err
		}
		for _, agent := range agentsForNode {
			var deleted uint
			deleted, err = tx.DeleteFrom(models.AgentNodeView, "WHERE node_id = ? AND agent_id = ?", node.ID, agent.ID)
			if err != nil {
				return errors.WithStack(err)
			}
			if deleted != 1 {
				return errors.Errorf("expected to delete 1 record, deleted %d", deleted)
			}
		}

		// stop agents
		agents := make(map[int32]models.Agent, len(agentsForService)+len(agentsForNode))
		for _, agent := range agentsForService {
			agents[agent.ID] = agent
		}
		for _, agent := range agentsForNode {
			agents[agent.ID] = agent
		}
		for _, agent := range agents {
			switch agent.Type {
			case models.MySQLdExporterAgentType:
				a := models.MySQLdExporter{ID: agent.ID}
				if err = tx.Reload(&a); err != nil {
					return errors.WithStack(err)
				}
				if svc.MySQLdExporterPath != "" {
					if err = svc.Supervisor.Stop(ctx, models.NameForSupervisor(a.Type, *a.ListenPort)); err != nil {
						return err
					}
				}

			case models.QanAgentAgentType:
				a := models.QanAgent{ID: agent.ID}
				if err = tx.Reload(&a); err != nil {
					return errors.WithStack(err)
				}
				if svc.QAN != nil {
					if node.Type == models.SSMServerNodeType {
						err = svc.QAN.RemoveMySQL(ctx, &a, true)
					} else {
						err = svc.QAN.RemoveMySQL(ctx, &a, false)
					}
					if err != nil {
						return err
					}
				}
			}
		}

		// remove agents
		for _, agent := range agents {
			if err = tx.Delete(&agent); err != nil {
				return errors.WithStack(err)
			}
		}

		if err = tx.Delete(&service); err != nil {
			return errors.WithStack(err)
		}

		if node.Type != models.SSMServerNodeType {
			if err = tx.Delete(&node); err != nil {
				return errors.WithStack(err)
			}
		}

		return svc.ApplyPrometheusConfiguration(ctx, tx.Querier)
	})
}

// Restore configuration from database.
func (svc *Service) Restore(ctx context.Context, tx *reform.TX) error {
	nodes, err := tx.FindAllFrom(models.RemoteNodeTable, "type", models.RemoteNodeType, models.SSMServerNodeType)
	if err != nil {
		return errors.WithStack(err)
	}
	for _, n := range nodes {
		node := n.(*models.RemoteNode)

		mySQLServices, err := tx.SelectAllFrom(models.MySQLServiceTable, "WHERE node_id = ? AND type = ?", node.ID, models.MySQLServiceType)
		if err != nil {
			return errors.WithStack(err)
		}
		if len(mySQLServices) == 0 {
			continue
		}
		if len(mySQLServices) != 1 {
			return errors.Errorf("expected to fetch 1 record, fetched %d. %v", len(mySQLServices), mySQLServices)
		}
		service := mySQLServices[0].(*models.MySQLService)
		if service.Type != models.MySQLServiceType {
			continue
		}

		agents, err := models.AgentsForServiceID(tx.Querier, service.ID)
		if err != nil {
			return err
		}
		for _, agent := range agents {
			switch agent.Type {
			case models.MySQLdExporterAgentType:
				a := &models.MySQLdExporter{ID: agent.ID}
				if err = tx.Reload(a); err != nil {
					return errors.WithStack(err)
				}
				if svc.MySQLdExporterPath != "" {
					name := models.NameForSupervisor(a.Type, *a.ListenPort)
					err := svc.Supervisor.Status(ctx, name)
					if err == nil {
						if err = svc.Supervisor.Stop(ctx, name); err != nil {
							return err
						}
					}

					dsn := a.DSN(service)
					cfg := svc.mysqlExporterCfg(a, dsn)
					if err = svc.Supervisor.Start(ctx, cfg); err != nil {
						return err
					}
				}

			case models.QanAgentAgentType:
				a := models.QanAgent{ID: agent.ID}
				if err = tx.Reload(&a); err != nil {
					return errors.WithStack(err)
				}
				if svc.QAN != nil {
					name := models.NameForSupervisor(a.Type, *a.ListenPort)
					err := svc.Supervisor.Status(ctx, name)
					if err == nil {
						if err = svc.Supervisor.Stop(ctx, name); err != nil {
							return err
						}
					}

					if err = svc.QAN.Restore(ctx, name, a); err != nil {
						if _, ok := err.(qan.QANCommandError); ok {
							// if it's a QAN command error, we should have already
							// restored the qan configs (although may not be perfectly),
							// one should check what happens on the qan-agent side, ssm-managed
							// should just continue on.
							logrus.WithField("component", "rds").Warnf("Got a QAN API error when restoring qan for %s: %s\n", node.Name, err.Error())
							return nil
						}

						return err
					}
				}
			}
		}
	}

	return nil
}

// EngineAndEngineVersion get mysql engine version
func (svc *Service) EngineAndEngineVersion(ctx context.Context, host string, port uint32, username string, password string) (string, string, error) {
	var version string
	var versionComment string
	agent := models.MySQLdExporter{
		ServiceUsername: pointer.ToString(username),
		ServicePassword: pointer.ToString(password),
	}
	service := &models.MySQLService{
		Address: &host,
		Port:    pointer.ToUint16(uint16(port)),
	}
	dsn := agent.DSN(service)
	db, err := sql.Open("mysql", dsn)
	if err == nil {
		sqlCtx, cancel := context.WithTimeout(ctx, sqlCheckTimeout)
		err = db.QueryRowContext(sqlCtx, "SELECT @@version, @@version_comment").Scan(&version, &versionComment)
		cancel()
		db.Close()
	}
	if err != nil {
		return "", "", errors.WithStack(err)
	}
	return NormalizeEngineAndEngineVersion(versionComment, version)
}

// NormalizeEngineAndEngineVersion normalizes the output from
// SQL 'SELECT @@version, @@version_comment'
func NormalizeEngineAndEngineVersion(engine string, engineVersion string) (string, string, error) {
	if versionRegexp.MatchString(engineVersion) {
		submatch := versionRegexp.FindStringSubmatch(engineVersion)
		engineVersion = submatch[1]
	}

	lowerEngine := strings.ToLower(engine)
	switch {
	case strings.Contains(lowerEngine, "mariadb"):
		return "MariaDB", engineVersion, nil
	case strings.Contains(lowerEngine, "percona"):
		return "Percona Server", engineVersion, nil
	default:
		return "MySQL", engineVersion, nil
	}
}
