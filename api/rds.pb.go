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
	return fileDescriptor_rds_f63eb022365793e0, []int{0}
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
	return fileDescriptor_rds_f63eb022365793e0, []int{1}
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
	return fileDescriptor_rds_f63eb022365793e0, []int{2}
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
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *RDSInstance) Reset()         { *m = RDSInstance{} }
func (m *RDSInstance) String() string { return proto.CompactTextString(m) }
func (*RDSInstance) ProtoMessage()    {}
func (*RDSInstance) Descriptor() ([]byte, []int) {
	return fileDescriptor_rds_f63eb022365793e0, []int{3}
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
	return fileDescriptor_rds_f63eb022365793e0, []int{4}
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
	return fileDescriptor_rds_f63eb022365793e0, []int{5}
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
	return fileDescriptor_rds_f63eb022365793e0, []int{6}
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
	return fileDescriptor_rds_f63eb022365793e0, []int{7}
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
	return fileDescriptor_rds_f63eb022365793e0, []int{8}
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
	return fileDescriptor_rds_f63eb022365793e0, []int{9}
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
	return fileDescriptor_rds_f63eb022365793e0, []int{10}
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
	return fileDescriptor_rds_f63eb022365793e0, []int{11}
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
	return fileDescriptor_rds_f63eb022365793e0, []int{12}
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
	return fileDescriptor_rds_f63eb022365793e0, []int{13}
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

func init() { proto.RegisterFile("rds.proto", fileDescriptor_rds_f63eb022365793e0) }

var fileDescriptor_rds_f63eb022365793e0 = []byte{
	// 704 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x55, 0x3d, 0x6f, 0xd3, 0x40,
	0x18, 0x96, 0x63, 0x93, 0x8f, 0x37, 0x6d, 0xe2, 0x5e, 0xda, 0x12, 0x05, 0x86, 0xc8, 0x12, 0xa2,
	0xed, 0x90, 0x40, 0x90, 0x18, 0x60, 0x4a, 0x65, 0x86, 0xb4, 0x08, 0xd0, 0x59, 0x65, 0xe8, 0x62,
	0x5d, 0x73, 0xa7, 0xe8, 0x44, 0xeb, 0x4b, 0xef, 0x9c, 0x44, 0x5d, 0x59, 0x18, 0x18, 0x11, 0x3f,
	0x81, 0x1f, 0xc3, 0xcc, 0x5f, 0xe0, 0x87, 0x20, 0x9f, 0xef, 0x1c, 0xa7, 0x7c, 0x0c, 0x48, 0xb0,
	0xf9, 0xfd, 0x7a, 0xee, 0x79, 0xde, 0x8f, 0x04, 0x1a, 0x92, 0xaa, 0xc1, 0x5c, 0x8a, 0x54, 0x20,
	0x97, 0xcc, 0x79, 0xef, 0xfe, 0x4c, 0x88, 0xd9, 0x25, 0x1b, 0x92, 0x39, 0x1f, 0x92, 0x24, 0x11,
	0x29, 0x49, 0xb9, 0x48, 0x4c, 0x4a, 0x30, 0x86, 0x1a, 0x0e, 0xa3, 0x57, 0x82, 0x32, 0xb4, 0x0f,
	0x55, 0xc9, 0x66, 0x5c, 0x24, 0x5d, 0xb7, 0xef, 0x1c, 0x34, 0xb0, 0xb1, 0x10, 0x02, 0x2f, 0x21,
	0x57, 0xac, 0xeb, 0x69, 0xaf, 0xfe, 0x3e, 0xf1, 0xea, 0x8e, 0x5f, 0x39, 0xf1, 0xea, 0x15, 0xdf,
	0x0d, 0x3e, 0x3a, 0x00, 0x38, 0x8c, 0x22, 0x26, 0x97, 0x7c, 0xca, 0x50, 0x17, 0x6a, 0x84, 0x52,
	0xc9, 0x94, 0x32, 0x15, 0xd6, 0xcc, 0x80, 0xe6, 0x42, 0xa6, 0xdd, 0x3b, 0x7d, 0xe7, 0x60, 0x1b,
	0xeb, 0xef, 0xec, 0x51, 0x96, 0xcc, 0x78, 0xc2, 0xba, 0xd5, 0xfc, 0xd1, 0xdc, 0x42, 0x0f, 0xa0,
	0x95, 0x7f, 0xc5, 0x4b, 0x26, 0x55, 0x46, 0xaa, 0xa6, 0xe3, 0xdb, 0xb9, 0xf7, 0x6d, 0xee, 0x2c,
	0xf3, 0x38, 0xf1, 0xea, 0xae, 0xef, 0x05, 0xcf, 0x61, 0x1b, 0x87, 0xd1, 0x24, 0x51, 0x29, 0x49,
	0xa6, 0x6c, 0x12, 0x96, 0x64, 0x39, 0xbf, 0x94, 0x55, 0x59, 0xcb, 0x0a, 0xce, 0xa1, 0x59, 0x2a,
	0x46, 0x7d, 0xf0, 0x12, 0x41, 0x99, 0x2e, 0x6c, 0x8e, 0xb6, 0x06, 0x64, 0xce, 0x07, 0xa6, 0x5b,
	0x58, 0x47, 0xd0, 0x21, 0xd4, 0x54, 0xae, 0x5b, 0xe3, 0x34, 0x47, 0x6d, 0x9b, 0x64, 0xda, 0x81,
	0x6d, 0x3c, 0x90, 0x80, 0x70, 0x18, 0x85, 0x5c, 0x4d, 0xc5, 0x92, 0x49, 0xcc, 0xae, 0x17, 0x4c,
	0xa5, 0xe8, 0x10, 0x76, 0xc8, 0x4a, 0xc5, 0x64, 0x3a, 0x65, 0x4a, 0xc5, 0xef, 0xd8, 0x4d, 0xcc,
	0xa9, 0x21, 0xda, 0x22, 0x2b, 0x35, 0xd6, 0xfe, 0x53, 0x76, 0x33, 0xa1, 0xe8, 0x31, 0xec, 0x65,
	0xa9, 0x8a, 0x4d, 0x25, 0x4b, 0x4b, 0x15, 0x46, 0x01, 0x22, 0x2b, 0x15, 0xe9, 0x58, 0x51, 0x14,
	0xbc, 0x80, 0xce, 0xc6, 0x9b, 0x6a, 0x2e, 0x12, 0xc5, 0xd0, 0x00, 0x1a, 0xdc, 0x68, 0x54, 0x5d,
	0xa7, 0xef, 0x1e, 0x34, 0x47, 0xbe, 0xe5, 0x6d, 0xc5, 0xe3, 0x75, 0x4a, 0xe0, 0x43, 0x0b, 0x87,
	0xd1, 0x4b, 0xae, 0x52, 0x43, 0x3b, 0x18, 0x43, 0xbb, 0xf0, 0xfc, 0x25, 0xe8, 0x57, 0x47, 0x4f,
	0x6a, 0x4c, 0xe9, 0x7f, 0xe9, 0x05, 0x0a, 0xa0, 0xc2, 0xa9, 0x5e, 0xed, 0xe6, 0x08, 0xdd, 0x26,
	0x36, 0x09, 0x71, 0x85, 0x53, 0xd4, 0x83, 0xfa, 0x42, 0x31, 0x59, 0x5a, 0xf7, 0xc2, 0xce, 0x62,
	0x73, 0xa2, 0xd4, 0x4a, 0x48, 0xaa, 0x37, 0xb8, 0x81, 0x0b, 0xdb, 0x34, 0x48, 0x4b, 0xc9, 0xbb,
	0x11, 0x3c, 0x05, 0x1f, 0x87, 0x11, 0x66, 0x57, 0x62, 0xc9, 0xac, 0xbe, 0x9c, 0x81, 0xf3, 0x27,
	0x06, 0x41, 0x07, 0x76, 0x4a, 0x75, 0x06, 0xec, 0xb3, 0xa3, 0xd1, 0x42, 0x96, 0x12, 0x7e, 0x69,
	0xd1, 0x86, 0xb0, 0x7b, 0x4d, 0x92, 0x98, 0x5e, 0xc4, 0xb6, 0xa7, 0xf1, 0x62, 0x51, 0x34, 0x6c,
	0xe7, 0x9a, 0x24, 0xe1, 0x85, 0x05, 0x3f, 0x5b, 0x70, 0x8a, 0x1e, 0x42, 0x3b, 0xdb, 0x45, 0x26,
	0xe3, 0x42, 0x63, 0xde, 0xad, 0x56, 0xee, 0x3e, 0xb3, 0x4a, 0xd7, 0x89, 0x85, 0x60, 0xb7, 0x9c,
	0xf8, 0xc6, 0xca, 0xfe, 0xe2, 0x68, 0xb6, 0x96, 0x97, 0x59, 0x84, 0x7f, 0x3b, 0xc6, 0xdf, 0xfd,
	0x4a, 0xf5, 0xa0, 0x6e, 0xfb, 0x60, 0x47, 0x67, 0xed, 0xd1, 0x07, 0x17, 0x5c, 0x1c, 0x46, 0xe8,
	0x1c, 0xea, 0xf6, 0x16, 0xd0, 0x5d, 0x3b, 0x80, 0x5b, 0x17, 0xd9, 0xeb, 0xfe, 0x1c, 0x30, 0x63,
	0xb8, 0xf7, 0xfe, 0xdb, 0xf7, 0x4f, 0x95, 0xbd, 0xc0, 0x1f, 0x2e, 0x1f, 0x0d, 0x25, 0x55, 0x43,
	0x6a, 0x32, 0x9e, 0x39, 0x47, 0xe8, 0x18, 0xbc, 0xec, 0x1c, 0x50, 0xc7, 0x96, 0x97, 0xce, 0xa5,
	0xb7, 0xbb, 0xe9, 0x34, 0x78, 0x6d, 0x8d, 0xd7, 0x40, 0x35, 0x83, 0x87, 0x8e, 0xc1, 0x1d, 0x53,
	0x8a, 0x8a, 0xdd, 0x58, 0xdf, 0x46, 0xaf, 0xb3, 0xe1, 0x33, 0x00, 0x48, 0x03, 0x6c, 0x05, 0x16,
	0x20, 0xe3, 0x71, 0x0a, 0xd5, 0x7c, 0x7b, 0xd0, 0x9e, 0x2d, 0xd9, 0xd8, 0xc2, 0xde, 0xfe, 0x6d,
	0xf7, 0x26, 0xd8, 0x51, 0x19, 0xec, 0x35, 0x54, 0xf3, 0xe1, 0xae, 0xc1, 0x36, 0x96, 0x70, 0x0d,
	0xb6, 0xb9, 0x03, 0xc1, 0xbe, 0x06, 0xf3, 0x51, 0xab, 0x68, 0x95, 0x8e, 0x5f, 0x54, 0xf5, 0xbf,
	0xce, 0x93, 0x1f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x45, 0x83, 0xe2, 0x46, 0xa5, 0x06, 0x00, 0x00,
}
