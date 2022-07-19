// Code generated by protoc-gen-go. DO NOT EDIT.
// source: rds.proto

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

type RDSNode struct {
	Region               string   `protobuf:"bytes,3,opt,name=region,proto3" json:"region,omitempty"`
	Name                 string   `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RDSNode) Reset()         { *m = RDSNode{} }
func (m *RDSNode) String() string { return proto.CompactTextString(m) }
func (*RDSNode) ProtoMessage()    {}
func (*RDSNode) Descriptor() ([]byte, []int) {
	return fileDescriptor_rds_ec0d5a9cc680e922, []int{0}
}
func (m *RDSNode) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RDSNode.Unmarshal(m, b)
}
func (m *RDSNode) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RDSNode.Marshal(b, m, deterministic)
}
func (dst *RDSNode) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RDSNode.Merge(dst, src)
}
func (m *RDSNode) XXX_Size() int {
	return xxx_messageInfo_RDSNode.Size(m)
}
func (m *RDSNode) XXX_DiscardUnknown() {
	xxx_messageInfo_RDSNode.DiscardUnknown(m)
}

var xxx_messageInfo_RDSNode proto.InternalMessageInfo

func (m *RDSNode) GetRegion() string {
	if m != nil {
		return m.Region
	}
	return ""
}

func (m *RDSNode) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type RDSService struct {
	Address              string   `protobuf:"bytes,4,opt,name=address,proto3" json:"address,omitempty"`
	Port                 uint32   `protobuf:"varint,5,opt,name=port,proto3" json:"port,omitempty"`
	Engine               string   `protobuf:"bytes,6,opt,name=engine,proto3" json:"engine,omitempty"`
	EngineVersion        string   `protobuf:"bytes,7,opt,name=engine_version,json=engineVersion,proto3" json:"engine_version,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RDSService) Reset()         { *m = RDSService{} }
func (m *RDSService) String() string { return proto.CompactTextString(m) }
func (*RDSService) ProtoMessage()    {}
func (*RDSService) Descriptor() ([]byte, []int) {
	return fileDescriptor_rds_ec0d5a9cc680e922, []int{1}
}
func (m *RDSService) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RDSService.Unmarshal(m, b)
}
func (m *RDSService) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RDSService.Marshal(b, m, deterministic)
}
func (dst *RDSService) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RDSService.Merge(dst, src)
}
func (m *RDSService) XXX_Size() int {
	return xxx_messageInfo_RDSService.Size(m)
}
func (m *RDSService) XXX_DiscardUnknown() {
	xxx_messageInfo_RDSService.DiscardUnknown(m)
}

var xxx_messageInfo_RDSService proto.InternalMessageInfo

func (m *RDSService) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *RDSService) GetPort() uint32 {
	if m != nil {
		return m.Port
	}
	return 0
}

func (m *RDSService) GetEngine() string {
	if m != nil {
		return m.Engine
	}
	return ""
}

func (m *RDSService) GetEngineVersion() string {
	if m != nil {
		return m.EngineVersion
	}
	return ""
}

type RDSAgent struct {
	QanDbInstanceUuid    string   `protobuf:"bytes,1,opt,name=qan_db_instance_uuid,json=qanDbInstanceUuid,proto3" json:"qan_db_instance_uuid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RDSAgent) Reset()         { *m = RDSAgent{} }
func (m *RDSAgent) String() string { return proto.CompactTextString(m) }
func (*RDSAgent) ProtoMessage()    {}
func (*RDSAgent) Descriptor() ([]byte, []int) {
	return fileDescriptor_rds_ec0d5a9cc680e922, []int{2}
}
func (m *RDSAgent) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RDSAgent.Unmarshal(m, b)
}
func (m *RDSAgent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RDSAgent.Marshal(b, m, deterministic)
}
func (dst *RDSAgent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RDSAgent.Merge(dst, src)
}
func (m *RDSAgent) XXX_Size() int {
	return xxx_messageInfo_RDSAgent.Size(m)
}
func (m *RDSAgent) XXX_DiscardUnknown() {
	xxx_messageInfo_RDSAgent.DiscardUnknown(m)
}

var xxx_messageInfo_RDSAgent proto.InternalMessageInfo

func (m *RDSAgent) GetQanDbInstanceUuid() string {
	if m != nil {
		return m.QanDbInstanceUuid
	}
	return ""
}

type RDSInstanceID struct {
	Region               string   `protobuf:"bytes,1,opt,name=region,proto3" json:"region,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RDSInstanceID) Reset()         { *m = RDSInstanceID{} }
func (m *RDSInstanceID) String() string { return proto.CompactTextString(m) }
func (*RDSInstanceID) ProtoMessage()    {}
func (*RDSInstanceID) Descriptor() ([]byte, []int) {
	return fileDescriptor_rds_ec0d5a9cc680e922, []int{3}
}
func (m *RDSInstanceID) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RDSInstanceID.Unmarshal(m, b)
}
func (m *RDSInstanceID) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RDSInstanceID.Marshal(b, m, deterministic)
}
func (dst *RDSInstanceID) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RDSInstanceID.Merge(dst, src)
}
func (m *RDSInstanceID) XXX_Size() int {
	return xxx_messageInfo_RDSInstanceID.Size(m)
}
func (m *RDSInstanceID) XXX_DiscardUnknown() {
	xxx_messageInfo_RDSInstanceID.DiscardUnknown(m)
}

var xxx_messageInfo_RDSInstanceID proto.InternalMessageInfo

func (m *RDSInstanceID) GetRegion() string {
	if m != nil {
		return m.Region
	}
	return ""
}

func (m *RDSInstanceID) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type RDSInstance struct {
	Node                 *RDSNode    `protobuf:"bytes,1,opt,name=node,proto3" json:"node,omitempty"`
	Service              *RDSService `protobuf:"bytes,2,opt,name=service,proto3" json:"service,omitempty"`
	Agent                *RDSAgent   `protobuf:"bytes,3,opt,name=agent,proto3" json:"agent,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *RDSInstance) Reset()         { *m = RDSInstance{} }
func (m *RDSInstance) String() string { return proto.CompactTextString(m) }
func (*RDSInstance) ProtoMessage()    {}
func (*RDSInstance) Descriptor() ([]byte, []int) {
	return fileDescriptor_rds_ec0d5a9cc680e922, []int{4}
}
func (m *RDSInstance) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RDSInstance.Unmarshal(m, b)
}
func (m *RDSInstance) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RDSInstance.Marshal(b, m, deterministic)
}
func (dst *RDSInstance) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RDSInstance.Merge(dst, src)
}
func (m *RDSInstance) XXX_Size() int {
	return xxx_messageInfo_RDSInstance.Size(m)
}
func (m *RDSInstance) XXX_DiscardUnknown() {
	xxx_messageInfo_RDSInstance.DiscardUnknown(m)
}

var xxx_messageInfo_RDSInstance proto.InternalMessageInfo

func (m *RDSInstance) GetNode() *RDSNode {
	if m != nil {
		return m.Node
	}
	return nil
}

func (m *RDSInstance) GetService() *RDSService {
	if m != nil {
		return m.Service
	}
	return nil
}

func (m *RDSInstance) GetAgent() *RDSAgent {
	if m != nil {
		return m.Agent
	}
	return nil
}

type RDSDiscoverRequest struct {
	AwsAccessKeyId       string   `protobuf:"bytes,1,opt,name=aws_access_key_id,json=awsAccessKeyId,proto3" json:"aws_access_key_id,omitempty"`
	AwsSecretAccessKey   string   `protobuf:"bytes,2,opt,name=aws_secret_access_key,json=awsSecretAccessKey,proto3" json:"aws_secret_access_key,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RDSDiscoverRequest) Reset()         { *m = RDSDiscoverRequest{} }
func (m *RDSDiscoverRequest) String() string { return proto.CompactTextString(m) }
func (*RDSDiscoverRequest) ProtoMessage()    {}
func (*RDSDiscoverRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_rds_ec0d5a9cc680e922, []int{5}
}
func (m *RDSDiscoverRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RDSDiscoverRequest.Unmarshal(m, b)
}
func (m *RDSDiscoverRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RDSDiscoverRequest.Marshal(b, m, deterministic)
}
func (dst *RDSDiscoverRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RDSDiscoverRequest.Merge(dst, src)
}
func (m *RDSDiscoverRequest) XXX_Size() int {
	return xxx_messageInfo_RDSDiscoverRequest.Size(m)
}
func (m *RDSDiscoverRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RDSDiscoverRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RDSDiscoverRequest proto.InternalMessageInfo

func (m *RDSDiscoverRequest) GetAwsAccessKeyId() string {
	if m != nil {
		return m.AwsAccessKeyId
	}
	return ""
}

func (m *RDSDiscoverRequest) GetAwsSecretAccessKey() string {
	if m != nil {
		return m.AwsSecretAccessKey
	}
	return ""
}

type RDSDiscoverResponse struct {
	Instances            []*RDSInstance `protobuf:"bytes,1,rep,name=instances,proto3" json:"instances,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *RDSDiscoverResponse) Reset()         { *m = RDSDiscoverResponse{} }
func (m *RDSDiscoverResponse) String() string { return proto.CompactTextString(m) }
func (*RDSDiscoverResponse) ProtoMessage()    {}
func (*RDSDiscoverResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_rds_ec0d5a9cc680e922, []int{6}
}
func (m *RDSDiscoverResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RDSDiscoverResponse.Unmarshal(m, b)
}
func (m *RDSDiscoverResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RDSDiscoverResponse.Marshal(b, m, deterministic)
}
func (dst *RDSDiscoverResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RDSDiscoverResponse.Merge(dst, src)
}
func (m *RDSDiscoverResponse) XXX_Size() int {
	return xxx_messageInfo_RDSDiscoverResponse.Size(m)
}
func (m *RDSDiscoverResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RDSDiscoverResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RDSDiscoverResponse proto.InternalMessageInfo

func (m *RDSDiscoverResponse) GetInstances() []*RDSInstance {
	if m != nil {
		return m.Instances
	}
	return nil
}

type RDSListRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RDSListRequest) Reset()         { *m = RDSListRequest{} }
func (m *RDSListRequest) String() string { return proto.CompactTextString(m) }
func (*RDSListRequest) ProtoMessage()    {}
func (*RDSListRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_rds_ec0d5a9cc680e922, []int{7}
}
func (m *RDSListRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RDSListRequest.Unmarshal(m, b)
}
func (m *RDSListRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RDSListRequest.Marshal(b, m, deterministic)
}
func (dst *RDSListRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RDSListRequest.Merge(dst, src)
}
func (m *RDSListRequest) XXX_Size() int {
	return xxx_messageInfo_RDSListRequest.Size(m)
}
func (m *RDSListRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RDSListRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RDSListRequest proto.InternalMessageInfo

type RDSListResponse struct {
	Instances            []*RDSInstance `protobuf:"bytes,1,rep,name=instances,proto3" json:"instances,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *RDSListResponse) Reset()         { *m = RDSListResponse{} }
func (m *RDSListResponse) String() string { return proto.CompactTextString(m) }
func (*RDSListResponse) ProtoMessage()    {}
func (*RDSListResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_rds_ec0d5a9cc680e922, []int{8}
}
func (m *RDSListResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RDSListResponse.Unmarshal(m, b)
}
func (m *RDSListResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RDSListResponse.Marshal(b, m, deterministic)
}
func (dst *RDSListResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RDSListResponse.Merge(dst, src)
}
func (m *RDSListResponse) XXX_Size() int {
	return xxx_messageInfo_RDSListResponse.Size(m)
}
func (m *RDSListResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RDSListResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RDSListResponse proto.InternalMessageInfo

func (m *RDSListResponse) GetInstances() []*RDSInstance {
	if m != nil {
		return m.Instances
	}
	return nil
}

type RDSAddRequest struct {
	AwsAccessKeyId       string         `protobuf:"bytes,1,opt,name=aws_access_key_id,json=awsAccessKeyId,proto3" json:"aws_access_key_id,omitempty"`
	AwsSecretAccessKey   string         `protobuf:"bytes,2,opt,name=aws_secret_access_key,json=awsSecretAccessKey,proto3" json:"aws_secret_access_key,omitempty"`
	Id                   *RDSInstanceID `protobuf:"bytes,3,opt,name=id,proto3" json:"id,omitempty"`
	Username             string         `protobuf:"bytes,4,opt,name=username,proto3" json:"username,omitempty"`
	Password             string         `protobuf:"bytes,5,opt,name=password,proto3" json:"password,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *RDSAddRequest) Reset()         { *m = RDSAddRequest{} }
func (m *RDSAddRequest) String() string { return proto.CompactTextString(m) }
func (*RDSAddRequest) ProtoMessage()    {}
func (*RDSAddRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_rds_ec0d5a9cc680e922, []int{9}
}
func (m *RDSAddRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RDSAddRequest.Unmarshal(m, b)
}
func (m *RDSAddRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RDSAddRequest.Marshal(b, m, deterministic)
}
func (dst *RDSAddRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RDSAddRequest.Merge(dst, src)
}
func (m *RDSAddRequest) XXX_Size() int {
	return xxx_messageInfo_RDSAddRequest.Size(m)
}
func (m *RDSAddRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RDSAddRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RDSAddRequest proto.InternalMessageInfo

func (m *RDSAddRequest) GetAwsAccessKeyId() string {
	if m != nil {
		return m.AwsAccessKeyId
	}
	return ""
}

func (m *RDSAddRequest) GetAwsSecretAccessKey() string {
	if m != nil {
		return m.AwsSecretAccessKey
	}
	return ""
}

func (m *RDSAddRequest) GetId() *RDSInstanceID {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *RDSAddRequest) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *RDSAddRequest) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

type RDSAddResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RDSAddResponse) Reset()         { *m = RDSAddResponse{} }
func (m *RDSAddResponse) String() string { return proto.CompactTextString(m) }
func (*RDSAddResponse) ProtoMessage()    {}
func (*RDSAddResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_rds_ec0d5a9cc680e922, []int{10}
}
func (m *RDSAddResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RDSAddResponse.Unmarshal(m, b)
}
func (m *RDSAddResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RDSAddResponse.Marshal(b, m, deterministic)
}
func (dst *RDSAddResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RDSAddResponse.Merge(dst, src)
}
func (m *RDSAddResponse) XXX_Size() int {
	return xxx_messageInfo_RDSAddResponse.Size(m)
}
func (m *RDSAddResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RDSAddResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RDSAddResponse proto.InternalMessageInfo

type RDSRemoveRequest struct {
	Id                   *RDSInstanceID `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *RDSRemoveRequest) Reset()         { *m = RDSRemoveRequest{} }
func (m *RDSRemoveRequest) String() string { return proto.CompactTextString(m) }
func (*RDSRemoveRequest) ProtoMessage()    {}
func (*RDSRemoveRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_rds_ec0d5a9cc680e922, []int{11}
}
func (m *RDSRemoveRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RDSRemoveRequest.Unmarshal(m, b)
}
func (m *RDSRemoveRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RDSRemoveRequest.Marshal(b, m, deterministic)
}
func (dst *RDSRemoveRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RDSRemoveRequest.Merge(dst, src)
}
func (m *RDSRemoveRequest) XXX_Size() int {
	return xxx_messageInfo_RDSRemoveRequest.Size(m)
}
func (m *RDSRemoveRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RDSRemoveRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RDSRemoveRequest proto.InternalMessageInfo

func (m *RDSRemoveRequest) GetId() *RDSInstanceID {
	if m != nil {
		return m.Id
	}
	return nil
}

type RDSRemoveResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RDSRemoveResponse) Reset()         { *m = RDSRemoveResponse{} }
func (m *RDSRemoveResponse) String() string { return proto.CompactTextString(m) }
func (*RDSRemoveResponse) ProtoMessage()    {}
func (*RDSRemoveResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_rds_ec0d5a9cc680e922, []int{12}
}
func (m *RDSRemoveResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RDSRemoveResponse.Unmarshal(m, b)
}
func (m *RDSRemoveResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RDSRemoveResponse.Marshal(b, m, deterministic)
}
func (dst *RDSRemoveResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RDSRemoveResponse.Merge(dst, src)
}
func (m *RDSRemoveResponse) XXX_Size() int {
	return xxx_messageInfo_RDSRemoveResponse.Size(m)
}
func (m *RDSRemoveResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RDSRemoveResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RDSRemoveResponse proto.InternalMessageInfo

type RDSDetailRequest struct {
	QanDbInstanceUuid    string   `protobuf:"bytes,1,opt,name=qan_db_instance_uuid,json=qanDbInstanceUuid,proto3" json:"qan_db_instance_uuid,omitempty"`
	ServerUsername       string   `protobuf:"bytes,2,opt,name=server_username,json=serverUsername,proto3" json:"server_username,omitempty"`
	ServerPassword       string   `protobuf:"bytes,3,opt,name=server_password,json=serverPassword,proto3" json:"server_password,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RDSDetailRequest) Reset()         { *m = RDSDetailRequest{} }
func (m *RDSDetailRequest) String() string { return proto.CompactTextString(m) }
func (*RDSDetailRequest) ProtoMessage()    {}
func (*RDSDetailRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_rds_ec0d5a9cc680e922, []int{13}
}
func (m *RDSDetailRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RDSDetailRequest.Unmarshal(m, b)
}
func (m *RDSDetailRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RDSDetailRequest.Marshal(b, m, deterministic)
}
func (dst *RDSDetailRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RDSDetailRequest.Merge(dst, src)
}
func (m *RDSDetailRequest) XXX_Size() int {
	return xxx_messageInfo_RDSDetailRequest.Size(m)
}
func (m *RDSDetailRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RDSDetailRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RDSDetailRequest proto.InternalMessageInfo

func (m *RDSDetailRequest) GetQanDbInstanceUuid() string {
	if m != nil {
		return m.QanDbInstanceUuid
	}
	return ""
}

func (m *RDSDetailRequest) GetServerUsername() string {
	if m != nil {
		return m.ServerUsername
	}
	return ""
}

func (m *RDSDetailRequest) GetServerPassword() string {
	if m != nil {
		return m.ServerPassword
	}
	return ""
}

type RDSDetailResponse struct {
	AwsAccessKeyId       string   `protobuf:"bytes,1,opt,name=aws_access_key_id,json=awsAccessKeyId,proto3" json:"aws_access_key_id,omitempty"`
	AwsSecretAccessKey   string   `protobuf:"bytes,2,opt,name=aws_secret_access_key,json=awsSecretAccessKey,proto3" json:"aws_secret_access_key,omitempty"`
	Region               string   `protobuf:"bytes,3,opt,name=region,proto3" json:"region,omitempty"`
	Instance             string   `protobuf:"bytes,4,opt,name=instance,proto3" json:"instance,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RDSDetailResponse) Reset()         { *m = RDSDetailResponse{} }
func (m *RDSDetailResponse) String() string { return proto.CompactTextString(m) }
func (*RDSDetailResponse) ProtoMessage()    {}
func (*RDSDetailResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_rds_ec0d5a9cc680e922, []int{14}
}
func (m *RDSDetailResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RDSDetailResponse.Unmarshal(m, b)
}
func (m *RDSDetailResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RDSDetailResponse.Marshal(b, m, deterministic)
}
func (dst *RDSDetailResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RDSDetailResponse.Merge(dst, src)
}
func (m *RDSDetailResponse) XXX_Size() int {
	return xxx_messageInfo_RDSDetailResponse.Size(m)
}
func (m *RDSDetailResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RDSDetailResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RDSDetailResponse proto.InternalMessageInfo

func (m *RDSDetailResponse) GetAwsAccessKeyId() string {
	if m != nil {
		return m.AwsAccessKeyId
	}
	return ""
}

func (m *RDSDetailResponse) GetAwsSecretAccessKey() string {
	if m != nil {
		return m.AwsSecretAccessKey
	}
	return ""
}

func (m *RDSDetailResponse) GetRegion() string {
	if m != nil {
		return m.Region
	}
	return ""
}

func (m *RDSDetailResponse) GetInstance() string {
	if m != nil {
		return m.Instance
	}
	return ""
}

func init() {
	proto.RegisterType((*RDSNode)(nil), "api.RDSNode")
	proto.RegisterType((*RDSService)(nil), "api.RDSService")
	proto.RegisterType((*RDSAgent)(nil), "api.RDSAgent")
	proto.RegisterType((*RDSInstanceID)(nil), "api.RDSInstanceID")
	proto.RegisterType((*RDSInstance)(nil), "api.RDSInstance")
	proto.RegisterType((*RDSDiscoverRequest)(nil), "api.RDSDiscoverRequest")
	proto.RegisterType((*RDSDiscoverResponse)(nil), "api.RDSDiscoverResponse")
	proto.RegisterType((*RDSListRequest)(nil), "api.RDSListRequest")
	proto.RegisterType((*RDSListResponse)(nil), "api.RDSListResponse")
	proto.RegisterType((*RDSAddRequest)(nil), "api.RDSAddRequest")
	proto.RegisterType((*RDSAddResponse)(nil), "api.RDSAddResponse")
	proto.RegisterType((*RDSRemoveRequest)(nil), "api.RDSRemoveRequest")
	proto.RegisterType((*RDSRemoveResponse)(nil), "api.RDSRemoveResponse")
	proto.RegisterType((*RDSDetailRequest)(nil), "api.RDSDetailRequest")
	proto.RegisterType((*RDSDetailResponse)(nil), "api.RDSDetailResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// RDSClient is the client API for RDS service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type RDSClient interface {
	Discover(ctx context.Context, in *RDSDiscoverRequest, opts ...grpc.CallOption) (*RDSDiscoverResponse, error)
	List(ctx context.Context, in *RDSListRequest, opts ...grpc.CallOption) (*RDSListResponse, error)
	Add(ctx context.Context, in *RDSAddRequest, opts ...grpc.CallOption) (*RDSAddResponse, error)
	Remove(ctx context.Context, in *RDSRemoveRequest, opts ...grpc.CallOption) (*RDSRemoveResponse, error)
	Detail(ctx context.Context, in *RDSDetailRequest, opts ...grpc.CallOption) (*RDSDetailResponse, error)
}

type rDSClient struct {
	cc *grpc.ClientConn
}

func NewRDSClient(cc *grpc.ClientConn) RDSClient {
	return &rDSClient{cc}
}

func (c *rDSClient) Discover(ctx context.Context, in *RDSDiscoverRequest, opts ...grpc.CallOption) (*RDSDiscoverResponse, error) {
	out := new(RDSDiscoverResponse)
	err := c.cc.Invoke(ctx, "/api.RDS/Discover", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rDSClient) List(ctx context.Context, in *RDSListRequest, opts ...grpc.CallOption) (*RDSListResponse, error) {
	out := new(RDSListResponse)
	err := c.cc.Invoke(ctx, "/api.RDS/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rDSClient) Add(ctx context.Context, in *RDSAddRequest, opts ...grpc.CallOption) (*RDSAddResponse, error) {
	out := new(RDSAddResponse)
	err := c.cc.Invoke(ctx, "/api.RDS/Add", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rDSClient) Remove(ctx context.Context, in *RDSRemoveRequest, opts ...grpc.CallOption) (*RDSRemoveResponse, error) {
	out := new(RDSRemoveResponse)
	err := c.cc.Invoke(ctx, "/api.RDS/Remove", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rDSClient) Detail(ctx context.Context, in *RDSDetailRequest, opts ...grpc.CallOption) (*RDSDetailResponse, error) {
	out := new(RDSDetailResponse)
	err := c.cc.Invoke(ctx, "/api.RDS/Detail", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RDSServer is the server API for RDS service.
type RDSServer interface {
	Discover(context.Context, *RDSDiscoverRequest) (*RDSDiscoverResponse, error)
	List(context.Context, *RDSListRequest) (*RDSListResponse, error)
	Add(context.Context, *RDSAddRequest) (*RDSAddResponse, error)
	Remove(context.Context, *RDSRemoveRequest) (*RDSRemoveResponse, error)
	Detail(context.Context, *RDSDetailRequest) (*RDSDetailResponse, error)
}

func RegisterRDSServer(s *grpc.Server, srv RDSServer) {
	s.RegisterService(&_RDS_serviceDesc, srv)
}

func _RDS_Discover_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RDSDiscoverRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RDSServer).Discover(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.RDS/Discover",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RDSServer).Discover(ctx, req.(*RDSDiscoverRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RDS_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RDSListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RDSServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.RDS/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RDSServer).List(ctx, req.(*RDSListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RDS_Add_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RDSAddRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RDSServer).Add(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.RDS/Add",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RDSServer).Add(ctx, req.(*RDSAddRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RDS_Remove_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RDSRemoveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RDSServer).Remove(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.RDS/Remove",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RDSServer).Remove(ctx, req.(*RDSRemoveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RDS_Detail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RDSDetailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RDSServer).Detail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.RDS/Detail",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RDSServer).Detail(ctx, req.(*RDSDetailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _RDS_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.RDS",
	HandlerType: (*RDSServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Discover",
			Handler:    _RDS_Discover_Handler,
		},
		{
			MethodName: "List",
			Handler:    _RDS_List_Handler,
		},
		{
			MethodName: "Add",
			Handler:    _RDS_Add_Handler,
		},
		{
			MethodName: "Remove",
			Handler:    _RDS_Remove_Handler,
		},
		{
			MethodName: "Detail",
			Handler:    _RDS_Detail_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rds.proto",
}

func init() { proto.RegisterFile("rds.proto", fileDescriptor_rds_ec0d5a9cc680e922) }

var fileDescriptor_rds_ec0d5a9cc680e922 = []byte{
	// 737 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x55, 0xcd, 0x6e, 0xd3, 0x4a,
	0x14, 0x96, 0x63, 0x37, 0x3f, 0x27, 0x4d, 0xe2, 0x4e, 0xda, 0x5e, 0x2b, 0xf7, 0x2e, 0xa2, 0xb9,
	0x42, 0xb4, 0x5d, 0x24, 0x10, 0x24, 0x16, 0x74, 0x95, 0xca, 0x2c, 0xd2, 0x22, 0x40, 0x63, 0x95,
	0x05, 0x1b, 0x6b, 0x9a, 0x19, 0x45, 0x23, 0x5a, 0x3b, 0xf5, 0x38, 0x89, 0xba, 0x42, 0x62, 0xc3,
	0x82, 0x25, 0xe2, 0x11, 0x78, 0x18, 0xd6, 0xbc, 0x02, 0x0f, 0x82, 0x3c, 0x9e, 0x71, 0x9c, 0xf2,
	0x23, 0x81, 0x04, 0x3b, 0x9f, 0xbf, 0x6f, 0xbe, 0xf3, 0x9d, 0x73, 0x12, 0x68, 0x24, 0x4c, 0x0e,
	0xe6, 0x49, 0x9c, 0xc6, 0xc8, 0xa6, 0x73, 0xd1, 0xfb, 0x6f, 0x16, 0xc7, 0xb3, 0x4b, 0x3e, 0xa4,
	0x73, 0x31, 0xa4, 0x51, 0x14, 0xa7, 0x34, 0x15, 0x71, 0xa4, 0x53, 0xf0, 0x18, 0x6a, 0xc4, 0x0f,
	0x9e, 0xc6, 0x8c, 0xa3, 0x7d, 0xa8, 0x26, 0x7c, 0x26, 0xe2, 0xc8, 0xb3, 0xfb, 0xd6, 0x41, 0x83,
	0x68, 0x0b, 0x21, 0x70, 0x22, 0x7a, 0xc5, 0x3d, 0x47, 0x79, 0xd5, 0xf7, 0xa9, 0x53, 0xb7, 0xdc,
	0xca, 0xa9, 0x53, 0xaf, 0xb8, 0x36, 0x7e, 0x67, 0x01, 0x10, 0x3f, 0x08, 0x78, 0xb2, 0x14, 0x53,
	0x8e, 0x3c, 0xa8, 0x51, 0xc6, 0x12, 0x2e, 0xa5, 0xae, 0x30, 0x66, 0x06, 0x34, 0x8f, 0x93, 0xd4,
	0xdb, 0xea, 0x5b, 0x07, 0x2d, 0xa2, 0xbe, 0xb3, 0x47, 0x79, 0x34, 0x13, 0x11, 0xf7, 0xaa, 0xf9,
	0xa3, 0xb9, 0x85, 0xee, 0x40, 0x3b, 0xff, 0x0a, 0x97, 0x3c, 0x91, 0x19, 0xa9, 0x9a, 0x8a, 0xb7,
	0x72, 0xef, 0x8b, 0xdc, 0x59, 0xe6, 0x71, 0xea, 0xd4, 0x6d, 0xd7, 0xc1, 0xc7, 0x50, 0x27, 0x7e,
	0x30, 0x9e, 0xf1, 0x28, 0x45, 0x43, 0xd8, 0xbd, 0xa6, 0x51, 0xc8, 0x2e, 0x42, 0x11, 0xc9, 0x94,
	0x46, 0x53, 0x1e, 0x2e, 0x16, 0x82, 0x79, 0x96, 0x82, 0xda, 0xb9, 0xa6, 0x91, 0x7f, 0x31, 0xd1,
	0x91, 0xf3, 0x85, 0x60, 0xf8, 0x18, 0x5a, 0xc4, 0x0f, 0x8c, 0x6b, 0xe2, 0x97, 0x34, 0xb1, 0xbe,
	0xab, 0x49, 0x65, 0xad, 0x09, 0x7e, 0x0d, 0xcd, 0x52, 0x31, 0xea, 0x83, 0x13, 0xc5, 0x8c, 0xab,
	0xc2, 0xe6, 0x68, 0x7b, 0x40, 0xe7, 0x62, 0xa0, 0xa5, 0x26, 0x2a, 0x82, 0x0e, 0xa1, 0x26, 0x73,
	0xd1, 0x14, 0x4e, 0x73, 0xd4, 0x31, 0x49, 0x5a, 0x4b, 0x62, 0xe2, 0xe8, 0x7f, 0xd8, 0xa2, 0x59,
	0x4b, 0x6a, 0x34, 0xcd, 0x51, 0xcb, 0x24, 0xaa, 0x3e, 0x49, 0x1e, 0xc3, 0x09, 0x20, 0xe2, 0x07,
	0xbe, 0x90, 0xd3, 0x78, 0xc9, 0x13, 0xc2, 0xaf, 0x17, 0x5c, 0xa6, 0xe8, 0x10, 0x76, 0xe8, 0x4a,
	0x86, 0x74, 0x3a, 0xe5, 0x52, 0x86, 0xaf, 0xf8, 0x4d, 0x58, 0x28, 0xd0, 0xa6, 0x2b, 0x39, 0x56,
	0xfe, 0x33, 0x7e, 0x33, 0x61, 0xe8, 0x3e, 0xec, 0x65, 0xa9, 0x92, 0x4f, 0x13, 0x9e, 0x96, 0x2a,
	0x74, 0x9b, 0x88, 0xae, 0x64, 0xa0, 0x62, 0x45, 0x11, 0x7e, 0x0c, 0xdd, 0x8d, 0x37, 0xe5, 0x3c,
	0x8e, 0x24, 0x47, 0x03, 0x68, 0x18, 0xc9, 0xa5, 0x67, 0xf5, 0xed, 0x83, 0xe6, 0xc8, 0x35, 0x9c,
	0x8d, 0x42, 0x64, 0x9d, 0x82, 0x5d, 0x68, 0x13, 0x3f, 0x78, 0x22, 0x64, 0xaa, 0x69, 0xe3, 0x31,
	0x74, 0x0a, 0xcf, 0x6f, 0x82, 0x7e, 0xb2, 0xd4, 0x38, 0xc7, 0x8c, 0xfd, 0x15, 0x2d, 0x10, 0x86,
	0x8a, 0x60, 0x7a, 0x42, 0xe8, 0x36, 0xb1, 0x89, 0x4f, 0x2a, 0x82, 0xa1, 0x1e, 0xd4, 0x17, 0x92,
	0x27, 0xa5, 0x83, 0x2a, 0xec, 0x2c, 0x36, 0xa7, 0x52, 0xae, 0xe2, 0x84, 0xa9, 0x1b, 0x69, 0x90,
	0xc2, 0xd6, 0x02, 0xa9, 0x56, 0x72, 0x35, 0xf0, 0x43, 0x70, 0x89, 0x1f, 0x10, 0x7e, 0x15, 0x2f,
	0xb9, 0xe9, 0x2f, 0x67, 0x60, 0xfd, 0x8c, 0x01, 0xee, 0xc2, 0x4e, 0xa9, 0x4e, 0x83, 0x7d, 0xb0,
	0x14, 0x9a, 0xcf, 0x53, 0x2a, 0x2e, 0x0d, 0xda, 0xaf, 0x9e, 0x0f, 0xba, 0x0b, 0x9d, 0x6c, 0x61,
	0x79, 0x12, 0x16, 0x3d, 0xe6, 0x6a, 0xb5, 0x73, 0xf7, 0xb9, 0xe9, 0x74, 0x9d, 0x58, 0x34, 0x6c,
	0x97, 0x13, 0x9f, 0x9b, 0xb6, 0x3f, 0x5a, 0x8a, 0xad, 0xe1, 0xa5, 0x17, 0xe1, 0xcf, 0x8e, 0xf1,
	0x47, 0xbf, 0x83, 0x3d, 0xa8, 0x1b, 0x1d, 0xcc, 0xe8, 0x8c, 0x3d, 0x7a, 0x6b, 0x83, 0x4d, 0xfc,
	0x00, 0xbd, 0x84, 0xba, 0xb9, 0x05, 0xf4, 0x8f, 0x19, 0xc0, 0xad, 0x8b, 0xec, 0x79, 0xdf, 0x06,
	0xf4, 0x18, 0xfe, 0x7d, 0xf3, 0xf9, 0xcb, 0xfb, 0xca, 0x1e, 0x76, 0x87, 0xcb, 0x7b, 0xc3, 0x84,
	0xc9, 0x21, 0xd3, 0x19, 0x8f, 0xac, 0x23, 0x74, 0x02, 0x4e, 0x76, 0x0e, 0xa8, 0x6b, 0xca, 0x4b,
	0xe7, 0xd2, 0xdb, 0xdd, 0x74, 0x6a, 0xbc, 0x8e, 0xc2, 0x6b, 0xa0, 0x9a, 0xc6, 0x43, 0x27, 0x60,
	0x8f, 0x19, 0x43, 0xc5, 0x6e, 0xac, 0x6f, 0xa3, 0xd7, 0xdd, 0xf0, 0x69, 0x00, 0xa4, 0x00, 0xb6,
	0xb1, 0x01, 0xc8, 0x78, 0x9c, 0x41, 0x35, 0xdf, 0x1e, 0xb4, 0x67, 0x4a, 0x36, 0xb6, 0xb0, 0xb7,
	0x7f, 0xdb, 0xbd, 0x09, 0x76, 0x54, 0x06, 0x7b, 0x06, 0xd5, 0x7c, 0xb8, 0x6b, 0xb0, 0x8d, 0x25,
	0x5c, 0x83, 0x6d, 0xee, 0x00, 0xde, 0x57, 0x60, 0x2e, 0x6a, 0x17, 0x52, 0xa9, 0xf8, 0x45, 0x55,
	0xfd, 0xaf, 0x3d, 0xf8, 0x1a, 0x00, 0x00, 0xff, 0xff, 0x44, 0xc5, 0xba, 0x30, 0x07, 0x07, 0x00,
	0x00,
}
