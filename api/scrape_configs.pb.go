// Code generated by protoc-gen-go. DO NOT EDIT.
// source: scrape_configs.proto

package api

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "google.golang.org/genproto/googleapis/api/annotations"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Target health : unknown, down, or up.
type ScrapeTargetHealth_Health int32

const (
	ScrapeTargetHealth_UNKNOWN ScrapeTargetHealth_Health = 0
	ScrapeTargetHealth_DOWN    ScrapeTargetHealth_Health = 1
	ScrapeTargetHealth_UP      ScrapeTargetHealth_Health = 2
)

var ScrapeTargetHealth_Health_name = map[int32]string{
	0: "UNKNOWN",
	1: "DOWN",
	2: "UP",
}
var ScrapeTargetHealth_Health_value = map[string]int32{
	"UNKNOWN": 0,
	"DOWN":    1,
	"UP":      2,
}

func (x ScrapeTargetHealth_Health) String() string {
	return proto.EnumName(ScrapeTargetHealth_Health_name, int32(x))
}
func (ScrapeTargetHealth_Health) EnumDescriptor() ([]byte, []int) { return fileDescriptor3, []int{5, 0} }

type LabelPair struct {
	// Label name
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	// Label value
	Value string `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
}

func (m *LabelPair) Reset()                    { *m = LabelPair{} }
func (m *LabelPair) String() string            { return proto.CompactTextString(m) }
func (*LabelPair) ProtoMessage()               {}
func (*LabelPair) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{0} }

func (m *LabelPair) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *LabelPair) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type StaticConfig struct {
	// Hostnames or IPs followed by an optional port number: "1.2.3.4:9090"
	Targets []string `protobuf:"bytes,1,rep,name=targets" json:"targets,omitempty"`
	// Labels assigned to all metrics scraped from the targets
	Labels []*LabelPair `protobuf:"bytes,2,rep,name=labels" json:"labels,omitempty"`
}

func (m *StaticConfig) Reset()                    { *m = StaticConfig{} }
func (m *StaticConfig) String() string            { return proto.CompactTextString(m) }
func (*StaticConfig) ProtoMessage()               {}
func (*StaticConfig) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{1} }

func (m *StaticConfig) GetTargets() []string {
	if m != nil {
		return m.Targets
	}
	return nil
}

func (m *StaticConfig) GetLabels() []*LabelPair {
	if m != nil {
		return m.Labels
	}
	return nil
}

type BasicAuth struct {
	Username string `protobuf:"bytes,1,opt,name=username" json:"username,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password" json:"password,omitempty"`
}

func (m *BasicAuth) Reset()                    { *m = BasicAuth{} }
func (m *BasicAuth) String() string            { return proto.CompactTextString(m) }
func (*BasicAuth) ProtoMessage()               {}
func (*BasicAuth) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{2} }

func (m *BasicAuth) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *BasicAuth) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

type TLSConfig struct {
	InsecureSkipVerify bool `protobuf:"varint,5,opt,name=insecure_skip_verify,json=insecureSkipVerify" json:"insecure_skip_verify,omitempty"`
}

func (m *TLSConfig) Reset()                    { *m = TLSConfig{} }
func (m *TLSConfig) String() string            { return proto.CompactTextString(m) }
func (*TLSConfig) ProtoMessage()               {}
func (*TLSConfig) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{3} }

func (m *TLSConfig) GetInsecureSkipVerify() bool {
	if m != nil {
		return m.InsecureSkipVerify
	}
	return false
}

type ScrapeConfig struct {
	// The job name assigned to scraped metrics by default: "example-job" (required)
	JobName string `protobuf:"bytes,1,opt,name=job_name,json=jobName" json:"job_name,omitempty"`
	// How frequently to scrape targets from this job: "10s"
	ScrapeInterval string `protobuf:"bytes,2,opt,name=scrape_interval,json=scrapeInterval" json:"scrape_interval,omitempty"`
	// Per-scrape timeout when scraping this job: "5s"
	ScrapeTimeout string `protobuf:"bytes,3,opt,name=scrape_timeout,json=scrapeTimeout" json:"scrape_timeout,omitempty"`
	// The HTTP resource path on which to fetch metrics from targets: "/metrics"
	MetricsPath string `protobuf:"bytes,4,opt,name=metrics_path,json=metricsPath" json:"metrics_path,omitempty"`
	// Configures the protocol scheme used for requests: "http" or "https"
	Scheme string `protobuf:"bytes,5,opt,name=scheme" json:"scheme,omitempty"`
	// Sets the `Authorization` header on every scrape request with the configured username and password
	BasicAuth *BasicAuth `protobuf:"bytes,6,opt,name=basic_auth,json=basicAuth" json:"basic_auth,omitempty"`
	// Configures the scrape request's TLS settings
	TlsConfig *TLSConfig `protobuf:"bytes,7,opt,name=tls_config,json=tlsConfig" json:"tls_config,omitempty"`
	// List of labeled statically configured targets for this job
	StaticConfigs []*StaticConfig `protobuf:"bytes,8,rep,name=static_configs,json=staticConfigs" json:"static_configs,omitempty"`
}

func (m *ScrapeConfig) Reset()                    { *m = ScrapeConfig{} }
func (m *ScrapeConfig) String() string            { return proto.CompactTextString(m) }
func (*ScrapeConfig) ProtoMessage()               {}
func (*ScrapeConfig) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{4} }

func (m *ScrapeConfig) GetJobName() string {
	if m != nil {
		return m.JobName
	}
	return ""
}

func (m *ScrapeConfig) GetScrapeInterval() string {
	if m != nil {
		return m.ScrapeInterval
	}
	return ""
}

func (m *ScrapeConfig) GetScrapeTimeout() string {
	if m != nil {
		return m.ScrapeTimeout
	}
	return ""
}

func (m *ScrapeConfig) GetMetricsPath() string {
	if m != nil {
		return m.MetricsPath
	}
	return ""
}

func (m *ScrapeConfig) GetScheme() string {
	if m != nil {
		return m.Scheme
	}
	return ""
}

func (m *ScrapeConfig) GetBasicAuth() *BasicAuth {
	if m != nil {
		return m.BasicAuth
	}
	return nil
}

func (m *ScrapeConfig) GetTlsConfig() *TLSConfig {
	if m != nil {
		return m.TlsConfig
	}
	return nil
}

func (m *ScrapeConfig) GetStaticConfigs() []*StaticConfig {
	if m != nil {
		return m.StaticConfigs
	}
	return nil
}

// ScrapeTargetHealth represents Prometheus scrape target health: unknown, down, or up.
type ScrapeTargetHealth struct {
	// Original scrape job name
	JobName string `protobuf:"bytes,1,opt,name=job_name,json=jobName" json:"job_name,omitempty"`
	// "job" label value, may be different from job_name due to relabeling
	Job string `protobuf:"bytes,2,opt,name=job" json:"job,omitempty"`
	// Original target
	Target string `protobuf:"bytes,3,opt,name=target" json:"target,omitempty"`
	// "instance" label value, may be different from target due to relabeling
	Instance string                    `protobuf:"bytes,4,opt,name=instance" json:"instance,omitempty"`
	Health   ScrapeTargetHealth_Health `protobuf:"varint,5,opt,name=health,enum=api.ScrapeTargetHealth_Health" json:"health,omitempty"`
}

func (m *ScrapeTargetHealth) Reset()                    { *m = ScrapeTargetHealth{} }
func (m *ScrapeTargetHealth) String() string            { return proto.CompactTextString(m) }
func (*ScrapeTargetHealth) ProtoMessage()               {}
func (*ScrapeTargetHealth) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{5} }

func (m *ScrapeTargetHealth) GetJobName() string {
	if m != nil {
		return m.JobName
	}
	return ""
}

func (m *ScrapeTargetHealth) GetJob() string {
	if m != nil {
		return m.Job
	}
	return ""
}

func (m *ScrapeTargetHealth) GetTarget() string {
	if m != nil {
		return m.Target
	}
	return ""
}

func (m *ScrapeTargetHealth) GetInstance() string {
	if m != nil {
		return m.Instance
	}
	return ""
}

func (m *ScrapeTargetHealth) GetHealth() ScrapeTargetHealth_Health {
	if m != nil {
		return m.Health
	}
	return ScrapeTargetHealth_UNKNOWN
}

type ScrapeConfigsListRequest struct {
}

func (m *ScrapeConfigsListRequest) Reset()                    { *m = ScrapeConfigsListRequest{} }
func (m *ScrapeConfigsListRequest) String() string            { return proto.CompactTextString(m) }
func (*ScrapeConfigsListRequest) ProtoMessage()               {}
func (*ScrapeConfigsListRequest) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{6} }

type ScrapeConfigsListResponse struct {
	ScrapeConfigs []*ScrapeConfig `protobuf:"bytes,1,rep,name=scrape_configs,json=scrapeConfigs" json:"scrape_configs,omitempty"`
	// Scrape targets health for all managed scrape jobs
	ScrapeTargetsHealth []*ScrapeTargetHealth `protobuf:"bytes,2,rep,name=scrape_targets_health,json=scrapeTargetsHealth" json:"scrape_targets_health,omitempty"`
}

func (m *ScrapeConfigsListResponse) Reset()                    { *m = ScrapeConfigsListResponse{} }
func (m *ScrapeConfigsListResponse) String() string            { return proto.CompactTextString(m) }
func (*ScrapeConfigsListResponse) ProtoMessage()               {}
func (*ScrapeConfigsListResponse) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{7} }

func (m *ScrapeConfigsListResponse) GetScrapeConfigs() []*ScrapeConfig {
	if m != nil {
		return m.ScrapeConfigs
	}
	return nil
}

func (m *ScrapeConfigsListResponse) GetScrapeTargetsHealth() []*ScrapeTargetHealth {
	if m != nil {
		return m.ScrapeTargetsHealth
	}
	return nil
}

type ScrapeConfigsGetRequest struct {
	JobName string `protobuf:"bytes,1,opt,name=job_name,json=jobName" json:"job_name,omitempty"`
}

func (m *ScrapeConfigsGetRequest) Reset()                    { *m = ScrapeConfigsGetRequest{} }
func (m *ScrapeConfigsGetRequest) String() string            { return proto.CompactTextString(m) }
func (*ScrapeConfigsGetRequest) ProtoMessage()               {}
func (*ScrapeConfigsGetRequest) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{8} }

func (m *ScrapeConfigsGetRequest) GetJobName() string {
	if m != nil {
		return m.JobName
	}
	return ""
}

type ScrapeConfigsGetResponse struct {
	ScrapeConfig *ScrapeConfig `protobuf:"bytes,1,opt,name=scrape_config,json=scrapeConfig" json:"scrape_config,omitempty"`
	// Scrape targets health for this scrape job
	ScrapeTargetsHealth []*ScrapeTargetHealth `protobuf:"bytes,2,rep,name=scrape_targets_health,json=scrapeTargetsHealth" json:"scrape_targets_health,omitempty"`
}

func (m *ScrapeConfigsGetResponse) Reset()                    { *m = ScrapeConfigsGetResponse{} }
func (m *ScrapeConfigsGetResponse) String() string            { return proto.CompactTextString(m) }
func (*ScrapeConfigsGetResponse) ProtoMessage()               {}
func (*ScrapeConfigsGetResponse) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{9} }

func (m *ScrapeConfigsGetResponse) GetScrapeConfig() *ScrapeConfig {
	if m != nil {
		return m.ScrapeConfig
	}
	return nil
}

func (m *ScrapeConfigsGetResponse) GetScrapeTargetsHealth() []*ScrapeTargetHealth {
	if m != nil {
		return m.ScrapeTargetsHealth
	}
	return nil
}

type ScrapeConfigsCreateRequest struct {
	ScrapeConfig *ScrapeConfig `protobuf:"bytes,1,opt,name=scrape_config,json=scrapeConfig" json:"scrape_config,omitempty"`
	// Check that added targets can be scraped from PMM Server
	CheckReachability bool `protobuf:"varint,2,opt,name=check_reachability,json=checkReachability" json:"check_reachability,omitempty"`
}

func (m *ScrapeConfigsCreateRequest) Reset()                    { *m = ScrapeConfigsCreateRequest{} }
func (m *ScrapeConfigsCreateRequest) String() string            { return proto.CompactTextString(m) }
func (*ScrapeConfigsCreateRequest) ProtoMessage()               {}
func (*ScrapeConfigsCreateRequest) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{10} }

func (m *ScrapeConfigsCreateRequest) GetScrapeConfig() *ScrapeConfig {
	if m != nil {
		return m.ScrapeConfig
	}
	return nil
}

func (m *ScrapeConfigsCreateRequest) GetCheckReachability() bool {
	if m != nil {
		return m.CheckReachability
	}
	return false
}

type ScrapeConfigsCreateResponse struct {
}

func (m *ScrapeConfigsCreateResponse) Reset()                    { *m = ScrapeConfigsCreateResponse{} }
func (m *ScrapeConfigsCreateResponse) String() string            { return proto.CompactTextString(m) }
func (*ScrapeConfigsCreateResponse) ProtoMessage()               {}
func (*ScrapeConfigsCreateResponse) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{11} }

type ScrapeConfigsDeleteRequest struct {
	JobName string `protobuf:"bytes,1,opt,name=job_name,json=jobName" json:"job_name,omitempty"`
}

func (m *ScrapeConfigsDeleteRequest) Reset()                    { *m = ScrapeConfigsDeleteRequest{} }
func (m *ScrapeConfigsDeleteRequest) String() string            { return proto.CompactTextString(m) }
func (*ScrapeConfigsDeleteRequest) ProtoMessage()               {}
func (*ScrapeConfigsDeleteRequest) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{12} }

func (m *ScrapeConfigsDeleteRequest) GetJobName() string {
	if m != nil {
		return m.JobName
	}
	return ""
}

type ScrapeConfigsDeleteResponse struct {
}

func (m *ScrapeConfigsDeleteResponse) Reset()                    { *m = ScrapeConfigsDeleteResponse{} }
func (m *ScrapeConfigsDeleteResponse) String() string            { return proto.CompactTextString(m) }
func (*ScrapeConfigsDeleteResponse) ProtoMessage()               {}
func (*ScrapeConfigsDeleteResponse) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{13} }

type ScrapeConfigsAddStaticTargetsRequest struct {
	JobName string `protobuf:"bytes,1,opt,name=job_name,json=jobName" json:"job_name,omitempty"`
	// Hostnames or IPs followed by an optional port number: "1.2.3.4:9090"
	Targets []string `protobuf:"bytes,2,rep,name=targets" json:"targets,omitempty"`
	// Check that added targets can be scraped from PMM Server
	CheckReachability bool `protobuf:"varint,3,opt,name=check_reachability,json=checkReachability" json:"check_reachability,omitempty"`
}

func (m *ScrapeConfigsAddStaticTargetsRequest) Reset()         { *m = ScrapeConfigsAddStaticTargetsRequest{} }
func (m *ScrapeConfigsAddStaticTargetsRequest) String() string { return proto.CompactTextString(m) }
func (*ScrapeConfigsAddStaticTargetsRequest) ProtoMessage()    {}
func (*ScrapeConfigsAddStaticTargetsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor3, []int{14}
}

func (m *ScrapeConfigsAddStaticTargetsRequest) GetJobName() string {
	if m != nil {
		return m.JobName
	}
	return ""
}

func (m *ScrapeConfigsAddStaticTargetsRequest) GetTargets() []string {
	if m != nil {
		return m.Targets
	}
	return nil
}

func (m *ScrapeConfigsAddStaticTargetsRequest) GetCheckReachability() bool {
	if m != nil {
		return m.CheckReachability
	}
	return false
}

type ScrapeConfigsAddStaticTargetsResponse struct {
}

func (m *ScrapeConfigsAddStaticTargetsResponse) Reset()         { *m = ScrapeConfigsAddStaticTargetsResponse{} }
func (m *ScrapeConfigsAddStaticTargetsResponse) String() string { return proto.CompactTextString(m) }
func (*ScrapeConfigsAddStaticTargetsResponse) ProtoMessage()    {}
func (*ScrapeConfigsAddStaticTargetsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor3, []int{15}
}

type ScrapeConfigsRemoveStaticTargetsRequest struct {
	JobName string `protobuf:"bytes,1,opt,name=job_name,json=jobName" json:"job_name,omitempty"`
	// Hostnames or IPs followed by an optional port number: "1.2.3.4:9090"
	Targets []string `protobuf:"bytes,2,rep,name=targets" json:"targets,omitempty"`
}

func (m *ScrapeConfigsRemoveStaticTargetsRequest) Reset() {
	*m = ScrapeConfigsRemoveStaticTargetsRequest{}
}
func (m *ScrapeConfigsRemoveStaticTargetsRequest) String() string { return proto.CompactTextString(m) }
func (*ScrapeConfigsRemoveStaticTargetsRequest) ProtoMessage()    {}
func (*ScrapeConfigsRemoveStaticTargetsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor3, []int{16}
}

func (m *ScrapeConfigsRemoveStaticTargetsRequest) GetJobName() string {
	if m != nil {
		return m.JobName
	}
	return ""
}

func (m *ScrapeConfigsRemoveStaticTargetsRequest) GetTargets() []string {
	if m != nil {
		return m.Targets
	}
	return nil
}

type ScrapeConfigsRemoveStaticTargetsResponse struct {
}

func (m *ScrapeConfigsRemoveStaticTargetsResponse) Reset() {
	*m = ScrapeConfigsRemoveStaticTargetsResponse{}
}
func (m *ScrapeConfigsRemoveStaticTargetsResponse) String() string { return proto.CompactTextString(m) }
func (*ScrapeConfigsRemoveStaticTargetsResponse) ProtoMessage()    {}
func (*ScrapeConfigsRemoveStaticTargetsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor3, []int{17}
}

func init() {
	proto.RegisterType((*LabelPair)(nil), "api.LabelPair")
	proto.RegisterType((*StaticConfig)(nil), "api.StaticConfig")
	proto.RegisterType((*BasicAuth)(nil), "api.BasicAuth")
	proto.RegisterType((*TLSConfig)(nil), "api.TLSConfig")
	proto.RegisterType((*ScrapeConfig)(nil), "api.ScrapeConfig")
	proto.RegisterType((*ScrapeTargetHealth)(nil), "api.ScrapeTargetHealth")
	proto.RegisterType((*ScrapeConfigsListRequest)(nil), "api.ScrapeConfigsListRequest")
	proto.RegisterType((*ScrapeConfigsListResponse)(nil), "api.ScrapeConfigsListResponse")
	proto.RegisterType((*ScrapeConfigsGetRequest)(nil), "api.ScrapeConfigsGetRequest")
	proto.RegisterType((*ScrapeConfigsGetResponse)(nil), "api.ScrapeConfigsGetResponse")
	proto.RegisterType((*ScrapeConfigsCreateRequest)(nil), "api.ScrapeConfigsCreateRequest")
	proto.RegisterType((*ScrapeConfigsCreateResponse)(nil), "api.ScrapeConfigsCreateResponse")
	proto.RegisterType((*ScrapeConfigsDeleteRequest)(nil), "api.ScrapeConfigsDeleteRequest")
	proto.RegisterType((*ScrapeConfigsDeleteResponse)(nil), "api.ScrapeConfigsDeleteResponse")
	proto.RegisterType((*ScrapeConfigsAddStaticTargetsRequest)(nil), "api.ScrapeConfigsAddStaticTargetsRequest")
	proto.RegisterType((*ScrapeConfigsAddStaticTargetsResponse)(nil), "api.ScrapeConfigsAddStaticTargetsResponse")
	proto.RegisterType((*ScrapeConfigsRemoveStaticTargetsRequest)(nil), "api.ScrapeConfigsRemoveStaticTargetsRequest")
	proto.RegisterType((*ScrapeConfigsRemoveStaticTargetsResponse)(nil), "api.ScrapeConfigsRemoveStaticTargetsResponse")
	proto.RegisterEnum("api.ScrapeTargetHealth_Health", ScrapeTargetHealth_Health_name, ScrapeTargetHealth_Health_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for ScrapeConfigs service

type ScrapeConfigsClient interface {
	// List returns all scrape configs.
	List(ctx context.Context, in *ScrapeConfigsListRequest, opts ...grpc.CallOption) (*ScrapeConfigsListResponse, error)
	// Get returns a scrape config by job name.
	// Errors: NotFound(5) if no such scrape config is present.
	Get(ctx context.Context, in *ScrapeConfigsGetRequest, opts ...grpc.CallOption) (*ScrapeConfigsGetResponse, error)
	// Create creates a new scrape config.
	// Errors: InvalidArgument(3) if some argument is not valid,
	// AlreadyExists(6) if scrape config with that job name is already present,
	// FailedPrecondition(9) if reachability check was requested and some scrape target can't be reached.
	Create(ctx context.Context, in *ScrapeConfigsCreateRequest, opts ...grpc.CallOption) (*ScrapeConfigsCreateResponse, error)
	// Delete removes existing scrape config by job name.
	// Errors: NotFound(5) if no such scrape config is present.
	Delete(ctx context.Context, in *ScrapeConfigsDeleteRequest, opts ...grpc.CallOption) (*ScrapeConfigsDeleteResponse, error)
	// Add static targets to existing scrape config.
	// Errors: NotFound(5) if no such scrape config is present,
	// FailedPrecondition(9) if reachability check was requested and some scrape target can't be reached.
	AddStaticTargets(ctx context.Context, in *ScrapeConfigsAddStaticTargetsRequest, opts ...grpc.CallOption) (*ScrapeConfigsAddStaticTargetsResponse, error)
	// Remove static targets from existing scrape config.
	// Errors: NotFound(5) if no such scrape config is present.
	RemoveStaticTargets(ctx context.Context, in *ScrapeConfigsRemoveStaticTargetsRequest, opts ...grpc.CallOption) (*ScrapeConfigsRemoveStaticTargetsResponse, error)
}

type scrapeConfigsClient struct {
	cc *grpc.ClientConn
}

func NewScrapeConfigsClient(cc *grpc.ClientConn) ScrapeConfigsClient {
	return &scrapeConfigsClient{cc}
}

func (c *scrapeConfigsClient) List(ctx context.Context, in *ScrapeConfigsListRequest, opts ...grpc.CallOption) (*ScrapeConfigsListResponse, error) {
	out := new(ScrapeConfigsListResponse)
	err := grpc.Invoke(ctx, "/api.ScrapeConfigs/List", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *scrapeConfigsClient) Get(ctx context.Context, in *ScrapeConfigsGetRequest, opts ...grpc.CallOption) (*ScrapeConfigsGetResponse, error) {
	out := new(ScrapeConfigsGetResponse)
	err := grpc.Invoke(ctx, "/api.ScrapeConfigs/Get", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *scrapeConfigsClient) Create(ctx context.Context, in *ScrapeConfigsCreateRequest, opts ...grpc.CallOption) (*ScrapeConfigsCreateResponse, error) {
	out := new(ScrapeConfigsCreateResponse)
	err := grpc.Invoke(ctx, "/api.ScrapeConfigs/Create", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *scrapeConfigsClient) Delete(ctx context.Context, in *ScrapeConfigsDeleteRequest, opts ...grpc.CallOption) (*ScrapeConfigsDeleteResponse, error) {
	out := new(ScrapeConfigsDeleteResponse)
	err := grpc.Invoke(ctx, "/api.ScrapeConfigs/Delete", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *scrapeConfigsClient) AddStaticTargets(ctx context.Context, in *ScrapeConfigsAddStaticTargetsRequest, opts ...grpc.CallOption) (*ScrapeConfigsAddStaticTargetsResponse, error) {
	out := new(ScrapeConfigsAddStaticTargetsResponse)
	err := grpc.Invoke(ctx, "/api.ScrapeConfigs/AddStaticTargets", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *scrapeConfigsClient) RemoveStaticTargets(ctx context.Context, in *ScrapeConfigsRemoveStaticTargetsRequest, opts ...grpc.CallOption) (*ScrapeConfigsRemoveStaticTargetsResponse, error) {
	out := new(ScrapeConfigsRemoveStaticTargetsResponse)
	err := grpc.Invoke(ctx, "/api.ScrapeConfigs/RemoveStaticTargets", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for ScrapeConfigs service

type ScrapeConfigsServer interface {
	// List returns all scrape configs.
	List(context.Context, *ScrapeConfigsListRequest) (*ScrapeConfigsListResponse, error)
	// Get returns a scrape config by job name.
	// Errors: NotFound(5) if no such scrape config is present.
	Get(context.Context, *ScrapeConfigsGetRequest) (*ScrapeConfigsGetResponse, error)
	// Create creates a new scrape config.
	// Errors: InvalidArgument(3) if some argument is not valid,
	// AlreadyExists(6) if scrape config with that job name is already present,
	// FailedPrecondition(9) if reachability check was requested and some scrape target can't be reached.
	Create(context.Context, *ScrapeConfigsCreateRequest) (*ScrapeConfigsCreateResponse, error)
	// Delete removes existing scrape config by job name.
	// Errors: NotFound(5) if no such scrape config is present.
	Delete(context.Context, *ScrapeConfigsDeleteRequest) (*ScrapeConfigsDeleteResponse, error)
	// Add static targets to existing scrape config.
	// Errors: NotFound(5) if no such scrape config is present,
	// FailedPrecondition(9) if reachability check was requested and some scrape target can't be reached.
	AddStaticTargets(context.Context, *ScrapeConfigsAddStaticTargetsRequest) (*ScrapeConfigsAddStaticTargetsResponse, error)
	// Remove static targets from existing scrape config.
	// Errors: NotFound(5) if no such scrape config is present.
	RemoveStaticTargets(context.Context, *ScrapeConfigsRemoveStaticTargetsRequest) (*ScrapeConfigsRemoveStaticTargetsResponse, error)
}

func RegisterScrapeConfigsServer(s *grpc.Server, srv ScrapeConfigsServer) {
	s.RegisterService(&_ScrapeConfigs_serviceDesc, srv)
}

func _ScrapeConfigs_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ScrapeConfigsListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScrapeConfigsServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.ScrapeConfigs/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScrapeConfigsServer).List(ctx, req.(*ScrapeConfigsListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ScrapeConfigs_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ScrapeConfigsGetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScrapeConfigsServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.ScrapeConfigs/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScrapeConfigsServer).Get(ctx, req.(*ScrapeConfigsGetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ScrapeConfigs_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ScrapeConfigsCreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScrapeConfigsServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.ScrapeConfigs/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScrapeConfigsServer).Create(ctx, req.(*ScrapeConfigsCreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ScrapeConfigs_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ScrapeConfigsDeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScrapeConfigsServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.ScrapeConfigs/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScrapeConfigsServer).Delete(ctx, req.(*ScrapeConfigsDeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ScrapeConfigs_AddStaticTargets_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ScrapeConfigsAddStaticTargetsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScrapeConfigsServer).AddStaticTargets(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.ScrapeConfigs/AddStaticTargets",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScrapeConfigsServer).AddStaticTargets(ctx, req.(*ScrapeConfigsAddStaticTargetsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ScrapeConfigs_RemoveStaticTargets_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ScrapeConfigsRemoveStaticTargetsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScrapeConfigsServer).RemoveStaticTargets(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.ScrapeConfigs/RemoveStaticTargets",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScrapeConfigsServer).RemoveStaticTargets(ctx, req.(*ScrapeConfigsRemoveStaticTargetsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ScrapeConfigs_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.ScrapeConfigs",
	HandlerType: (*ScrapeConfigsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "List",
			Handler:    _ScrapeConfigs_List_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _ScrapeConfigs_Get_Handler,
		},
		{
			MethodName: "Create",
			Handler:    _ScrapeConfigs_Create_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _ScrapeConfigs_Delete_Handler,
		},
		{
			MethodName: "AddStaticTargets",
			Handler:    _ScrapeConfigs_AddStaticTargets_Handler,
		},
		{
			MethodName: "RemoveStaticTargets",
			Handler:    _ScrapeConfigs_RemoveStaticTargets_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "scrape_configs.proto",
}

func init() { proto.RegisterFile("scrape_configs.proto", fileDescriptor3) }

var fileDescriptor3 = []byte{
	// 932 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x56, 0xdd, 0x6e, 0xdc, 0x44,
	0x14, 0xc6, 0xbb, 0xa9, 0xb3, 0x3e, 0xf9, 0x21, 0x9d, 0x06, 0xe2, 0x9a, 0x6e, 0xbb, 0x8c, 0x08,
	0x59, 0x56, 0xdd, 0x6c, 0x09, 0x50, 0x10, 0x12, 0x17, 0x25, 0x95, 0x0a, 0x6a, 0x14, 0x22, 0x27,
	0x85, 0x3b, 0xac, 0xb1, 0x33, 0x8d, 0x27, 0xf1, 0xda, 0xc6, 0x33, 0xbb, 0xa8, 0x42, 0xdc, 0xc0,
	0x1d, 0x57, 0x48, 0x3c, 0x00, 0x52, 0x6f, 0x79, 0x15, 0xee, 0x90, 0x78, 0x02, 0x1e, 0x04, 0x79,
	0x66, 0xec, 0xda, 0xbb, 0xde, 0xa4, 0xa0, 0x5e, 0xc5, 0x73, 0xce, 0x99, 0x73, 0xbe, 0xef, 0xf3,
	0x77, 0x9c, 0x85, 0x4d, 0x1e, 0x64, 0x24, 0xa5, 0x5e, 0x90, 0xc4, 0x4f, 0xd9, 0x19, 0xdf, 0x4d,
	0xb3, 0x44, 0x24, 0xa8, 0x4d, 0x52, 0xe6, 0xdc, 0x3a, 0x4b, 0x92, 0xb3, 0x88, 0x8e, 0x48, 0xca,
	0x46, 0x24, 0x8e, 0x13, 0x41, 0x04, 0x4b, 0x62, 0x5d, 0x82, 0x3f, 0x02, 0xeb, 0x80, 0xf8, 0x34,
	0x3a, 0x22, 0x2c, 0x43, 0x08, 0x96, 0x62, 0x32, 0xa6, 0xb6, 0xd1, 0x33, 0xfa, 0x96, 0x2b, 0x9f,
	0xd1, 0x26, 0x5c, 0x9b, 0x92, 0x68, 0x42, 0xed, 0x96, 0x0c, 0xaa, 0x03, 0x3e, 0x82, 0xd5, 0xe3,
	0xbc, 0x51, 0xb0, 0x2f, 0x07, 0x22, 0x1b, 0x96, 0x05, 0xc9, 0xce, 0xa8, 0xe0, 0xb6, 0xd1, 0x6b,
	0xf7, 0x2d, 0xb7, 0x38, 0xa2, 0x77, 0xc1, 0x8c, 0xf2, 0x01, 0xdc, 0x6e, 0xf5, 0xda, 0xfd, 0x95,
	0xbd, 0xf5, 0x5d, 0x92, 0xb2, 0xdd, 0x72, 0xa6, 0xab, 0xb3, 0x78, 0x1f, 0xac, 0xcf, 0x09, 0x67,
	0xc1, 0x83, 0x89, 0x08, 0x91, 0x03, 0x9d, 0x09, 0xa7, 0x59, 0x05, 0x4c, 0x79, 0xce, 0x73, 0x29,
	0xe1, 0xfc, 0xfb, 0x24, 0x3b, 0xd5, 0x98, 0xca, 0x33, 0xfe, 0x0c, 0xac, 0x93, 0x83, 0x63, 0x8d,
	0xe9, 0x1e, 0x6c, 0xb2, 0x98, 0xd3, 0x60, 0x92, 0x51, 0x8f, 0x5f, 0xb0, 0xd4, 0x9b, 0xd2, 0x8c,
	0x3d, 0x7d, 0x66, 0x5f, 0xeb, 0x19, 0xfd, 0x8e, 0x8b, 0x8a, 0xdc, 0xf1, 0x05, 0x4b, 0xbf, 0x96,
	0x19, 0xfc, 0x67, 0x0b, 0x56, 0x8f, 0xa5, 0x90, 0xba, 0xc5, 0x4d, 0xe8, 0x9c, 0x27, 0xbe, 0x57,
	0xc1, 0xb1, 0x7c, 0x9e, 0xf8, 0x87, 0x39, 0x8c, 0x1d, 0x78, 0x5d, 0x6b, 0xce, 0x62, 0x41, 0xb3,
	0x29, 0x89, 0x34, 0x9a, 0x75, 0x15, 0xfe, 0x52, 0x47, 0xd1, 0x36, 0xe8, 0x88, 0x27, 0xd8, 0x98,
	0x26, 0x13, 0x61, 0xb7, 0x65, 0xdd, 0x9a, 0x8a, 0x9e, 0xa8, 0x20, 0x7a, 0x1b, 0x56, 0xc7, 0x54,
	0x64, 0x2c, 0xe0, 0x5e, 0x4a, 0x44, 0x68, 0x2f, 0xc9, 0xa2, 0x15, 0x1d, 0x3b, 0x22, 0x22, 0x44,
	0x6f, 0x82, 0xc9, 0x83, 0x90, 0x8e, 0xa9, 0xa4, 0x60, 0xb9, 0xfa, 0x84, 0x86, 0x00, 0x7e, 0x2e,
	0x9d, 0x47, 0x26, 0x22, 0xb4, 0xcd, 0x9e, 0x51, 0xca, 0x5c, 0x2a, 0xea, 0x5a, 0x7e, 0x29, 0xee,
	0x10, 0x40, 0x44, 0x5c, 0x5b, 0xc5, 0x5e, 0xae, 0x94, 0x97, 0xda, 0xb9, 0x96, 0x88, 0xb8, 0xd6,
	0xe0, 0x13, 0x58, 0xe7, 0xf2, 0x55, 0x17, 0xe6, 0xb2, 0x3b, 0xf2, 0x45, 0x5e, 0x97, 0x57, 0xaa,
	0x2e, 0x70, 0xd7, 0x78, 0xe5, 0xc4, 0xf1, 0xdf, 0x06, 0x20, 0x25, 0xe7, 0x89, 0x34, 0xc3, 0x17,
	0x94, 0x44, 0x22, 0xbc, 0x4c, 0xd4, 0x0d, 0x68, 0x9f, 0x27, 0xbe, 0x16, 0x32, 0x7f, 0xcc, 0x39,
	0x2b, 0x27, 0x69, 0xd5, 0xf4, 0x29, 0x77, 0x01, 0x8b, 0xb9, 0x20, 0x71, 0x40, 0xb5, 0x54, 0xe5,
	0x19, 0xdd, 0x07, 0x33, 0x94, 0xa3, 0xa4, 0x4e, 0xeb, 0x7b, 0xb7, 0x15, 0xd2, 0x39, 0x24, 0xbb,
	0xea, 0x8f, 0xab, 0xab, 0xf1, 0x0e, 0x98, 0x1a, 0xe2, 0x0a, 0x2c, 0x3f, 0x39, 0x7c, 0x7c, 0xf8,
	0xd5, 0x37, 0x87, 0x1b, 0xaf, 0xa1, 0x0e, 0x2c, 0x3d, 0xcc, 0x9f, 0x0c, 0x64, 0x42, 0xeb, 0xc9,
	0xd1, 0x46, 0x0b, 0x3b, 0x60, 0x57, 0x6d, 0xc2, 0x0f, 0x18, 0x17, 0x2e, 0xfd, 0x6e, 0x42, 0xb9,
	0xc0, 0xcf, 0x0d, 0xb8, 0xd9, 0x90, 0xe4, 0x69, 0x12, 0x73, 0x2a, 0xc5, 0xac, 0x6d, 0xaa, 0x5c,
	0x97, 0x52, 0xcc, 0xca, 0xbd, 0xc2, 0x1f, 0xba, 0x0b, 0x7a, 0x0c, 0x6f, 0x14, 0x36, 0x52, 0x9b,
	0xe5, 0x69, 0x8e, 0x6a, 0xad, 0xb6, 0x16, 0x70, 0x74, 0x6f, 0xf0, 0x4a, 0x8c, 0xab, 0x20, 0xfe,
	0x10, 0xb6, 0x6a, 0x18, 0x1f, 0xd1, 0x02, 0xff, 0x25, 0x6f, 0x07, 0xff, 0x6e, 0xcc, 0xf0, 0x96,
	0xd7, 0x34, 0xb3, 0xfb, 0xb0, 0x56, 0x63, 0x26, 0x2f, 0x37, 0x12, 0x5b, 0xad, 0x12, 0x7b, 0xb5,
	0xbc, 0x7e, 0x36, 0xc0, 0xa9, 0x21, 0xdc, 0xcf, 0x28, 0x11, 0xb4, 0xe0, 0xf6, 0x7f, 0x31, 0x0e,
	0x01, 0x05, 0x21, 0x0d, 0x2e, 0xbc, 0x8c, 0x92, 0x20, 0x24, 0x3e, 0x8b, 0x98, 0x78, 0x26, 0x5d,
	0xda, 0x71, 0xaf, 0xcb, 0x8c, 0x5b, 0x49, 0xe0, 0x2e, 0xbc, 0xd5, 0x08, 0x42, 0x29, 0x85, 0x3f,
	0x9e, 0xc1, 0xf8, 0x90, 0x46, 0xf4, 0x05, 0xc6, 0x4b, 0xf4, 0x9f, 0xed, 0x5b, 0x5c, 0xd4, 0x7d,
	0x7f, 0x31, 0xe0, 0x9d, 0x5a, 0xfe, 0xc1, 0xe9, 0xa9, 0x5a, 0x4f, 0x2d, 0xd2, 0xd5, 0x23, 0xaa,
	0xdf, 0xf1, 0x56, 0xfd, 0x3b, 0xde, 0xac, 0x41, 0x7b, 0x91, 0x06, 0x3b, 0xb0, 0x7d, 0x05, 0x16,
	0x8d, 0xfa, 0x5b, 0xd8, 0xa9, 0x15, 0xba, 0x74, 0x9c, 0x4c, 0xe9, 0x2b, 0xc3, 0x8d, 0x07, 0xd0,
	0xbf, 0xba, 0xbf, 0xc2, 0xb2, 0xf7, 0xab, 0x09, 0x6b, 0xb5, 0x62, 0x44, 0x60, 0x29, 0xdf, 0x5f,
	0xd4, 0x9d, 0xb3, 0x48, 0x75, 0xe9, 0x9d, 0xdb, 0x8b, 0xd2, 0x9a, 0xa4, 0xf3, 0xd3, 0x5f, 0xff,
	0xfc, 0xd6, 0xda, 0x44, 0x68, 0x34, 0xbd, 0x37, 0x52, 0xd6, 0x1a, 0xea, 0x0f, 0x00, 0x62, 0xd0,
	0x7e, 0x44, 0x05, 0xba, 0x35, 0xdf, 0xe2, 0xc5, 0x56, 0x3a, 0xdd, 0x05, 0x59, 0xdd, 0x7f, 0x5b,
	0xf6, 0xbf, 0x83, 0xba, 0xf3, 0xfd, 0x47, 0x3f, 0x14, 0x9a, 0xfd, 0x88, 0xce, 0xc1, 0x54, 0x5e,
	0x44, 0x77, 0xe6, 0xfb, 0xd5, 0x56, 0xc5, 0xe9, 0x2d, 0x2e, 0xd0, 0x33, 0xbb, 0x72, 0xe6, 0x16,
	0x6e, 0xe0, 0xf4, 0xa9, 0x31, 0x40, 0x19, 0x98, 0xca, 0x9f, 0x4d, 0xb3, 0x6a, 0x96, 0x6f, 0x9a,
	0x35, 0x63, 0x6d, 0xcd, 0x6f, 0x70, 0x05, 0xbf, 0xe7, 0x06, 0x6c, 0xcc, 0x1a, 0x0d, 0xbd, 0x37,
	0xdf, 0x7d, 0xc1, 0x62, 0x38, 0x83, 0x97, 0x29, 0x2d, 0xb6, 0x58, 0x42, 0x7a, 0x1f, 0xdf, 0xbd,
	0x14, 0xd2, 0x48, 0xfd, 0x47, 0x1c, 0x6a, 0x37, 0xe6, 0xc2, 0xfc, 0x61, 0xc0, 0x8d, 0x06, 0x13,
	0xa2, 0xbb, 0xf3, 0xc3, 0x17, 0xef, 0x82, 0x33, 0x7c, 0xc9, 0xea, 0x3a, 0xda, 0xc1, 0x7f, 0x45,
	0xeb, 0x9b, 0xf2, 0x67, 0xe2, 0x07, 0xff, 0x06, 0x00, 0x00, 0xff, 0xff, 0xcf, 0x1d, 0xed, 0x7b,
	0x61, 0x0a, 0x00, 0x00,
}
