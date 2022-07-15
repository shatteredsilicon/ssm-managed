// Code generated by protoc-gen-go. DO NOT EDIT.
// source: logs.proto

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

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Log struct {
	// Last lines of log file
	Lines                []string `protobuf:"bytes,1,rep,name=lines,proto3" json:"lines,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Log) Reset()         { *m = Log{} }
func (m *Log) String() string { return proto.CompactTextString(m) }
func (*Log) ProtoMessage()    {}
func (*Log) Descriptor() ([]byte, []int) {
	return fileDescriptor_logs_eb2aff7eafdc83fb, []int{0}
}
func (m *Log) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Log.Unmarshal(m, b)
}
func (m *Log) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Log.Marshal(b, m, deterministic)
}
func (dst *Log) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Log.Merge(dst, src)
}
func (m *Log) XXX_Size() int {
	return xxx_messageInfo_Log.Size(m)
}
func (m *Log) XXX_DiscardUnknown() {
	xxx_messageInfo_Log.DiscardUnknown(m)
}

var xxx_messageInfo_Log proto.InternalMessageInfo

func (m *Log) GetLines() []string {
	if m != nil {
		return m.Lines
	}
	return nil
}

type LogsAllRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LogsAllRequest) Reset()         { *m = LogsAllRequest{} }
func (m *LogsAllRequest) String() string { return proto.CompactTextString(m) }
func (*LogsAllRequest) ProtoMessage()    {}
func (*LogsAllRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_logs_eb2aff7eafdc83fb, []int{1}
}
func (m *LogsAllRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LogsAllRequest.Unmarshal(m, b)
}
func (m *LogsAllRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LogsAllRequest.Marshal(b, m, deterministic)
}
func (dst *LogsAllRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LogsAllRequest.Merge(dst, src)
}
func (m *LogsAllRequest) XXX_Size() int {
	return xxx_messageInfo_LogsAllRequest.Size(m)
}
func (m *LogsAllRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_LogsAllRequest.DiscardUnknown(m)
}

var xxx_messageInfo_LogsAllRequest proto.InternalMessageInfo

type LogsAllResponse struct {
	// Maps log file name to content
	Logs                 map[string]*Log `protobuf:"bytes,1,rep,name=logs,proto3" json:"logs,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *LogsAllResponse) Reset()         { *m = LogsAllResponse{} }
func (m *LogsAllResponse) String() string { return proto.CompactTextString(m) }
func (*LogsAllResponse) ProtoMessage()    {}
func (*LogsAllResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_logs_eb2aff7eafdc83fb, []int{2}
}
func (m *LogsAllResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LogsAllResponse.Unmarshal(m, b)
}
func (m *LogsAllResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LogsAllResponse.Marshal(b, m, deterministic)
}
func (dst *LogsAllResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LogsAllResponse.Merge(dst, src)
}
func (m *LogsAllResponse) XXX_Size() int {
	return xxx_messageInfo_LogsAllResponse.Size(m)
}
func (m *LogsAllResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_LogsAllResponse.DiscardUnknown(m)
}

var xxx_messageInfo_LogsAllResponse proto.InternalMessageInfo

func (m *LogsAllResponse) GetLogs() map[string]*Log {
	if m != nil {
		return m.Logs
	}
	return nil
}

func init() {
	proto.RegisterType((*Log)(nil), "api.Log")
	proto.RegisterType((*LogsAllRequest)(nil), "api.LogsAllRequest")
	proto.RegisterType((*LogsAllResponse)(nil), "api.LogsAllResponse")
	proto.RegisterMapType((map[string]*Log)(nil), "api.LogsAllResponse.LogsEntry")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// LogsClient is the client API for Logs service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type LogsClient interface {
	// All returns last lines of all log files.
	All(ctx context.Context, in *LogsAllRequest, opts ...grpc.CallOption) (*LogsAllResponse, error)
}

type logsClient struct {
	cc *grpc.ClientConn
}

func NewLogsClient(cc *grpc.ClientConn) LogsClient {
	return &logsClient{cc}
}

func (c *logsClient) All(ctx context.Context, in *LogsAllRequest, opts ...grpc.CallOption) (*LogsAllResponse, error) {
	out := new(LogsAllResponse)
	err := c.cc.Invoke(ctx, "/api.Logs/All", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LogsServer is the server API for Logs service.
type LogsServer interface {
	// All returns last lines of all log files.
	All(context.Context, *LogsAllRequest) (*LogsAllResponse, error)
}

func RegisterLogsServer(s *grpc.Server, srv LogsServer) {
	s.RegisterService(&_Logs_serviceDesc, srv)
}

func _Logs_All_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogsAllRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogsServer).All(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Logs/All",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogsServer).All(ctx, req.(*LogsAllRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Logs_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.Logs",
	HandlerType: (*LogsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "All",
			Handler:    _Logs_All_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "logs.proto",
}

func init() { proto.RegisterFile("logs.proto", fileDescriptor_logs_eb2aff7eafdc83fb) }

var fileDescriptor_logs_eb2aff7eafdc83fb = []byte{
	// 243 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xca, 0xc9, 0x4f, 0x2f,
	0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x4e, 0x2c, 0xc8, 0x94, 0x92, 0x49, 0xcf, 0xcf,
	0x4f, 0xcf, 0x49, 0xd5, 0x4f, 0x2c, 0xc8, 0xd4, 0x4f, 0xcc, 0xcb, 0xcb, 0x2f, 0x49, 0x2c, 0xc9,
	0xcc, 0xcf, 0x83, 0x2a, 0x51, 0x92, 0xe6, 0x62, 0xf6, 0xc9, 0x4f, 0x17, 0x12, 0xe1, 0x62, 0xcd,
	0xc9, 0xcc, 0x4b, 0x2d, 0x96, 0x60, 0x54, 0x60, 0xd6, 0xe0, 0x0c, 0x82, 0x70, 0x94, 0x04, 0xb8,
	0xf8, 0x7c, 0xf2, 0xd3, 0x8b, 0x1d, 0x73, 0x72, 0x82, 0x52, 0x0b, 0x4b, 0x53, 0x8b, 0x4b, 0x94,
	0x3a, 0x18, 0xb9, 0xf8, 0xe1, 0x42, 0xc5, 0x05, 0xf9, 0x79, 0xc5, 0xa9, 0x42, 0x46, 0x5c, 0x2c,
	0x20, 0x3b, 0xc1, 0x5a, 0xb9, 0x8d, 0xe4, 0xf4, 0x12, 0x0b, 0x32, 0xf5, 0xd0, 0xd4, 0x80, 0xf9,
	0xae, 0x79, 0x25, 0x45, 0x95, 0x41, 0x60, 0xb5, 0x52, 0x8e, 0x5c, 0x9c, 0x70, 0x21, 0x21, 0x01,
	0x2e, 0xe6, 0xec, 0xd4, 0x4a, 0x09, 0x46, 0x05, 0x46, 0x0d, 0xce, 0x20, 0x10, 0x53, 0x48, 0x8e,
	0x8b, 0xb5, 0x2c, 0x31, 0xa7, 0x34, 0x55, 0x82, 0x49, 0x81, 0x51, 0x83, 0xdb, 0x88, 0x03, 0x66,
	0x66, 0x10, 0x44, 0xd8, 0x8a, 0xc9, 0x82, 0xd1, 0xc8, 0x8b, 0x8b, 0x05, 0x64, 0x84, 0x90, 0x13,
	0x17, 0xb3, 0x63, 0x4e, 0x8e, 0x90, 0x30, 0xaa, 0xbd, 0x60, 0xe7, 0x4a, 0x89, 0x60, 0x73, 0x8c,
	0x92, 0x40, 0xd3, 0xe5, 0x27, 0x93, 0x99, 0xb8, 0x84, 0x38, 0xf4, 0xcb, 0x0c, 0xf4, 0x41, 0xce,
	0x49, 0x62, 0x03, 0x07, 0x86, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0x14, 0x0d, 0xbe, 0xd7, 0x3d,
	0x01, 0x00, 0x00,
}
