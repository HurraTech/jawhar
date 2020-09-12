// Code generated by protoc-gen-go. DO NOT EDIT.
// source: hurra_agent.proto

package proto

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type MountDriveResponse struct {
	Errors               []string `protobuf:"bytes,1,rep,name=errors,proto3" json:"errors,omitempty"`
	IsSuccessful         bool     `protobuf:"varint,2,opt,name=is_successful,json=isSuccessful,proto3" json:"is_successful,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MountDriveResponse) Reset()         { *m = MountDriveResponse{} }
func (m *MountDriveResponse) String() string { return proto.CompactTextString(m) }
func (*MountDriveResponse) ProtoMessage()    {}
func (*MountDriveResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_d1a676bd7648f983, []int{0}
}

func (m *MountDriveResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MountDriveResponse.Unmarshal(m, b)
}
func (m *MountDriveResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MountDriveResponse.Marshal(b, m, deterministic)
}
func (m *MountDriveResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MountDriveResponse.Merge(m, src)
}
func (m *MountDriveResponse) XXX_Size() int {
	return xxx_messageInfo_MountDriveResponse.Size(m)
}
func (m *MountDriveResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MountDriveResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MountDriveResponse proto.InternalMessageInfo

func (m *MountDriveResponse) GetErrors() []string {
	if m != nil {
		return m.Errors
	}
	return nil
}

func (m *MountDriveResponse) GetIsSuccessful() bool {
	if m != nil {
		return m.IsSuccessful
	}
	return false
}

type MountDriveRequest struct {
	DeviceFile           string   `protobuf:"bytes,1,opt,name=device_file,json=deviceFile,proto3" json:"device_file,omitempty"`
	MountPoint           string   `protobuf:"bytes,2,opt,name=mount_point,json=mountPoint,proto3" json:"mount_point,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MountDriveRequest) Reset()         { *m = MountDriveRequest{} }
func (m *MountDriveRequest) String() string { return proto.CompactTextString(m) }
func (*MountDriveRequest) ProtoMessage()    {}
func (*MountDriveRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_d1a676bd7648f983, []int{1}
}

func (m *MountDriveRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MountDriveRequest.Unmarshal(m, b)
}
func (m *MountDriveRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MountDriveRequest.Marshal(b, m, deterministic)
}
func (m *MountDriveRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MountDriveRequest.Merge(m, src)
}
func (m *MountDriveRequest) XXX_Size() int {
	return xxx_messageInfo_MountDriveRequest.Size(m)
}
func (m *MountDriveRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_MountDriveRequest.DiscardUnknown(m)
}

var xxx_messageInfo_MountDriveRequest proto.InternalMessageInfo

func (m *MountDriveRequest) GetDeviceFile() string {
	if m != nil {
		return m.DeviceFile
	}
	return ""
}

func (m *MountDriveRequest) GetMountPoint() string {
	if m != nil {
		return m.MountPoint
	}
	return ""
}

type UnmountDriveRequest struct {
	DeviceFile           string   `protobuf:"bytes,1,opt,name=device_file,json=deviceFile,proto3" json:"device_file,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UnmountDriveRequest) Reset()         { *m = UnmountDriveRequest{} }
func (m *UnmountDriveRequest) String() string { return proto.CompactTextString(m) }
func (*UnmountDriveRequest) ProtoMessage()    {}
func (*UnmountDriveRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_d1a676bd7648f983, []int{2}
}

func (m *UnmountDriveRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UnmountDriveRequest.Unmarshal(m, b)
}
func (m *UnmountDriveRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UnmountDriveRequest.Marshal(b, m, deterministic)
}
func (m *UnmountDriveRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UnmountDriveRequest.Merge(m, src)
}
func (m *UnmountDriveRequest) XXX_Size() int {
	return xxx_messageInfo_UnmountDriveRequest.Size(m)
}
func (m *UnmountDriveRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UnmountDriveRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UnmountDriveRequest proto.InternalMessageInfo

func (m *UnmountDriveRequest) GetDeviceFile() string {
	if m != nil {
		return m.DeviceFile
	}
	return ""
}

type UnmountDriveResponse struct {
	Error                string   `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
	IsSuccessful         bool     `protobuf:"varint,2,opt,name=is_successful,json=isSuccessful,proto3" json:"is_successful,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UnmountDriveResponse) Reset()         { *m = UnmountDriveResponse{} }
func (m *UnmountDriveResponse) String() string { return proto.CompactTextString(m) }
func (*UnmountDriveResponse) ProtoMessage()    {}
func (*UnmountDriveResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_d1a676bd7648f983, []int{3}
}

func (m *UnmountDriveResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UnmountDriveResponse.Unmarshal(m, b)
}
func (m *UnmountDriveResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UnmountDriveResponse.Marshal(b, m, deterministic)
}
func (m *UnmountDriveResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UnmountDriveResponse.Merge(m, src)
}
func (m *UnmountDriveResponse) XXX_Size() int {
	return xxx_messageInfo_UnmountDriveResponse.Size(m)
}
func (m *UnmountDriveResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_UnmountDriveResponse.DiscardUnknown(m)
}

var xxx_messageInfo_UnmountDriveResponse proto.InternalMessageInfo

func (m *UnmountDriveResponse) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

func (m *UnmountDriveResponse) GetIsSuccessful() bool {
	if m != nil {
		return m.IsSuccessful
	}
	return false
}

type GetDrivesRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetDrivesRequest) Reset()         { *m = GetDrivesRequest{} }
func (m *GetDrivesRequest) String() string { return proto.CompactTextString(m) }
func (*GetDrivesRequest) ProtoMessage()    {}
func (*GetDrivesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_d1a676bd7648f983, []int{4}
}

func (m *GetDrivesRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetDrivesRequest.Unmarshal(m, b)
}
func (m *GetDrivesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetDrivesRequest.Marshal(b, m, deterministic)
}
func (m *GetDrivesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetDrivesRequest.Merge(m, src)
}
func (m *GetDrivesRequest) XXX_Size() int {
	return xxx_messageInfo_GetDrivesRequest.Size(m)
}
func (m *GetDrivesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetDrivesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetDrivesRequest proto.InternalMessageInfo

type GetDrivesResponse struct {
	Drives               []*Drive `protobuf:"bytes,1,rep,name=drives,proto3" json:"drives,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetDrivesResponse) Reset()         { *m = GetDrivesResponse{} }
func (m *GetDrivesResponse) String() string { return proto.CompactTextString(m) }
func (*GetDrivesResponse) ProtoMessage()    {}
func (*GetDrivesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_d1a676bd7648f983, []int{5}
}

func (m *GetDrivesResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetDrivesResponse.Unmarshal(m, b)
}
func (m *GetDrivesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetDrivesResponse.Marshal(b, m, deterministic)
}
func (m *GetDrivesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetDrivesResponse.Merge(m, src)
}
func (m *GetDrivesResponse) XXX_Size() int {
	return xxx_messageInfo_GetDrivesResponse.Size(m)
}
func (m *GetDrivesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetDrivesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetDrivesResponse proto.InternalMessageInfo

func (m *GetDrivesResponse) GetDrives() []*Drive {
	if m != nil {
		return m.Drives
	}
	return nil
}

type Drive struct {
	Name                 string       `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	DeviceFile           string       `protobuf:"bytes,2,opt,name=device_file,json=deviceFile,proto3" json:"device_file,omitempty"`
	SizeBytes            uint64       `protobuf:"varint,3,opt,name=size_bytes,json=sizeBytes,proto3" json:"size_bytes,omitempty"`
	IsRemovable          bool         `protobuf:"varint,4,opt,name=is_removable,json=isRemovable,proto3" json:"is_removable,omitempty"`
	Type                 string       `protobuf:"bytes,5,opt,name=type,proto3" json:"type,omitempty"`
	SerialNumber         string       `protobuf:"bytes,6,opt,name=serial_number,json=serialNumber,proto3" json:"serial_number,omitempty"`
	StorageController    string       `protobuf:"bytes,7,opt,name=storage_controller,json=storageController,proto3" json:"storage_controller,omitempty"`
	Partitions           []*Partition `protobuf:"bytes,8,rep,name=partitions,proto3" json:"partitions,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Drive) Reset()         { *m = Drive{} }
func (m *Drive) String() string { return proto.CompactTextString(m) }
func (*Drive) ProtoMessage()    {}
func (*Drive) Descriptor() ([]byte, []int) {
	return fileDescriptor_d1a676bd7648f983, []int{6}
}

func (m *Drive) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Drive.Unmarshal(m, b)
}
func (m *Drive) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Drive.Marshal(b, m, deterministic)
}
func (m *Drive) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Drive.Merge(m, src)
}
func (m *Drive) XXX_Size() int {
	return xxx_messageInfo_Drive.Size(m)
}
func (m *Drive) XXX_DiscardUnknown() {
	xxx_messageInfo_Drive.DiscardUnknown(m)
}

var xxx_messageInfo_Drive proto.InternalMessageInfo

func (m *Drive) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Drive) GetDeviceFile() string {
	if m != nil {
		return m.DeviceFile
	}
	return ""
}

func (m *Drive) GetSizeBytes() uint64 {
	if m != nil {
		return m.SizeBytes
	}
	return 0
}

func (m *Drive) GetIsRemovable() bool {
	if m != nil {
		return m.IsRemovable
	}
	return false
}

func (m *Drive) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Drive) GetSerialNumber() string {
	if m != nil {
		return m.SerialNumber
	}
	return ""
}

func (m *Drive) GetStorageController() string {
	if m != nil {
		return m.StorageController
	}
	return ""
}

func (m *Drive) GetPartitions() []*Partition {
	if m != nil {
		return m.Partitions
	}
	return nil
}

type Partition struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	DeviceFile           string   `protobuf:"bytes,2,opt,name=device_file,json=deviceFile,proto3" json:"device_file,omitempty"`
	SizeBytes            uint64   `protobuf:"varint,3,opt,name=size_bytes,json=sizeBytes,proto3" json:"size_bytes,omitempty"`
	AvailableBytes       uint64   `protobuf:"varint,4,opt,name=available_bytes,json=availableBytes,proto3" json:"available_bytes,omitempty"`
	Filesystem           string   `protobuf:"bytes,5,opt,name=filesystem,proto3" json:"filesystem,omitempty"`
	MountPoint           string   `protobuf:"bytes,6,opt,name=mount_point,json=mountPoint,proto3" json:"mount_point,omitempty"`
	Label                string   `protobuf:"bytes,7,opt,name=label,proto3" json:"label,omitempty"`
	IsReadOnly           bool     `protobuf:"varint,8,opt,name=is_read_only,json=isReadOnly,proto3" json:"is_read_only,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Partition) Reset()         { *m = Partition{} }
func (m *Partition) String() string { return proto.CompactTextString(m) }
func (*Partition) ProtoMessage()    {}
func (*Partition) Descriptor() ([]byte, []int) {
	return fileDescriptor_d1a676bd7648f983, []int{7}
}

func (m *Partition) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Partition.Unmarshal(m, b)
}
func (m *Partition) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Partition.Marshal(b, m, deterministic)
}
func (m *Partition) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Partition.Merge(m, src)
}
func (m *Partition) XXX_Size() int {
	return xxx_messageInfo_Partition.Size(m)
}
func (m *Partition) XXX_DiscardUnknown() {
	xxx_messageInfo_Partition.DiscardUnknown(m)
}

var xxx_messageInfo_Partition proto.InternalMessageInfo

func (m *Partition) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Partition) GetDeviceFile() string {
	if m != nil {
		return m.DeviceFile
	}
	return ""
}

func (m *Partition) GetSizeBytes() uint64 {
	if m != nil {
		return m.SizeBytes
	}
	return 0
}

func (m *Partition) GetAvailableBytes() uint64 {
	if m != nil {
		return m.AvailableBytes
	}
	return 0
}

func (m *Partition) GetFilesystem() string {
	if m != nil {
		return m.Filesystem
	}
	return ""
}

func (m *Partition) GetMountPoint() string {
	if m != nil {
		return m.MountPoint
	}
	return ""
}

func (m *Partition) GetLabel() string {
	if m != nil {
		return m.Label
	}
	return ""
}

func (m *Partition) GetIsReadOnly() bool {
	if m != nil {
		return m.IsReadOnly
	}
	return false
}

type Result struct {
	ExitCode             int32    `protobuf:"varint,1,opt,name=exitCode,proto3" json:"exitCode,omitempty"`
	Message              string   `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Result) Reset()         { *m = Result{} }
func (m *Result) String() string { return proto.CompactTextString(m) }
func (*Result) ProtoMessage()    {}
func (*Result) Descriptor() ([]byte, []int) {
	return fileDescriptor_d1a676bd7648f983, []int{8}
}

func (m *Result) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Result.Unmarshal(m, b)
}
func (m *Result) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Result.Marshal(b, m, deterministic)
}
func (m *Result) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Result.Merge(m, src)
}
func (m *Result) XXX_Size() int {
	return xxx_messageInfo_Result.Size(m)
}
func (m *Result) XXX_DiscardUnknown() {
	xxx_messageInfo_Result.DiscardUnknown(m)
}

var xxx_messageInfo_Result proto.InternalMessageInfo

func (m *Result) GetExitCode() int32 {
	if m != nil {
		return m.ExitCode
	}
	return 0
}

func (m *Result) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type Command struct {
	Command              string   `protobuf:"bytes,1,opt,name=command,proto3" json:"command,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Command) Reset()         { *m = Command{} }
func (m *Command) String() string { return proto.CompactTextString(m) }
func (*Command) ProtoMessage()    {}
func (*Command) Descriptor() ([]byte, []int) {
	return fileDescriptor_d1a676bd7648f983, []int{9}
}

func (m *Command) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Command.Unmarshal(m, b)
}
func (m *Command) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Command.Marshal(b, m, deterministic)
}
func (m *Command) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Command.Merge(m, src)
}
func (m *Command) XXX_Size() int {
	return xxx_messageInfo_Command.Size(m)
}
func (m *Command) XXX_DiscardUnknown() {
	xxx_messageInfo_Command.DiscardUnknown(m)
}

var xxx_messageInfo_Command proto.InternalMessageInfo

func (m *Command) GetCommand() string {
	if m != nil {
		return m.Command
	}
	return ""
}

func init() {
	proto.RegisterType((*MountDriveResponse)(nil), "proto.MountDriveResponse")
	proto.RegisterType((*MountDriveRequest)(nil), "proto.MountDriveRequest")
	proto.RegisterType((*UnmountDriveRequest)(nil), "proto.UnmountDriveRequest")
	proto.RegisterType((*UnmountDriveResponse)(nil), "proto.UnmountDriveResponse")
	proto.RegisterType((*GetDrivesRequest)(nil), "proto.GetDrivesRequest")
	proto.RegisterType((*GetDrivesResponse)(nil), "proto.GetDrivesResponse")
	proto.RegisterType((*Drive)(nil), "proto.Drive")
	proto.RegisterType((*Partition)(nil), "proto.Partition")
	proto.RegisterType((*Result)(nil), "proto.Result")
	proto.RegisterType((*Command)(nil), "proto.Command")
}

func init() {
	proto.RegisterFile("hurra_agent.proto", fileDescriptor_d1a676bd7648f983)
}

var fileDescriptor_d1a676bd7648f983 = []byte{
	// 632 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xb4, 0x54, 0xcf, 0x6f, 0xd3, 0x30,
	0x14, 0x5e, 0xbb, 0x36, 0x5b, 0x5f, 0xbb, 0x1f, 0x35, 0x13, 0x84, 0x4c, 0xc0, 0xc8, 0x90, 0xd8,
	0x85, 0x08, 0x0d, 0x09, 0xc4, 0x05, 0xc1, 0x0a, 0x0c, 0x0e, 0x40, 0x09, 0xda, 0x85, 0x4b, 0xe4,
	0x26, 0xde, 0xb0, 0x94, 0xc4, 0xc1, 0x76, 0xaa, 0x95, 0x1b, 0xff, 0x00, 0x67, 0xfe, 0x5c, 0xfc,
	0xab, 0xa1, 0x74, 0x3d, 0xc0, 0x81, 0x53, 0xfc, 0xbe, 0xf7, 0xf9, 0xb3, 0xdf, 0xe7, 0xf7, 0x02,
	0xc3, 0x2f, 0x35, 0xe7, 0x38, 0xc1, 0x17, 0xa4, 0x94, 0x51, 0xc5, 0x99, 0x64, 0xa8, 0x6b, 0x3e,
	0xe1, 0x47, 0x40, 0xef, 0x58, 0x5d, 0xca, 0x97, 0x9c, 0x4e, 0x49, 0x4c, 0x44, 0xc5, 0x4a, 0x41,
	0xd0, 0x75, 0xf0, 0x08, 0xe7, 0x8c, 0x0b, 0xbf, 0x75, 0xb0, 0x7e, 0xd4, 0x8b, 0x5d, 0x84, 0x0e,
	0x61, 0x8b, 0x8a, 0x44, 0xd4, 0x69, 0x4a, 0x84, 0x38, 0xaf, 0x73, 0xbf, 0x7d, 0xd0, 0x3a, 0xda,
	0x8c, 0x07, 0x54, 0x7c, 0x6a, 0xb0, 0xf0, 0x0c, 0x86, 0x8b, 0x92, 0x5f, 0x6b, 0x22, 0x24, 0xba,
	0x03, 0xfd, 0x8c, 0x4c, 0x69, 0x4a, 0x92, 0x73, 0x9a, 0x13, 0x25, 0xdb, 0x52, 0xb2, 0x60, 0xa1,
	0xd7, 0x0a, 0xd1, 0x84, 0x42, 0xef, 0x4a, 0x2a, 0x46, 0x4b, 0x69, 0x84, 0x15, 0xc1, 0x40, 0x63,
	0x8d, 0x84, 0x8f, 0xe1, 0xda, 0x59, 0x59, 0xfc, 0xb3, 0xb0, 0xaa, 0x70, 0xef, 0xcf, 0x7d, 0xae,
	0xc6, 0x3d, 0xe8, 0x9a, 0xaa, 0xdc, 0x16, 0x1b, 0xfc, 0x5d, 0x85, 0x08, 0x76, 0x4f, 0x89, 0x95,
	0x13, 0xee, 0x1e, 0xe1, 0x53, 0x18, 0x2e, 0x60, 0xee, 0x8c, 0x7b, 0xe0, 0x65, 0x06, 0x31, 0x3e,
	0xf6, 0x8f, 0x07, 0xd6, 0xfc, 0xc8, 0xde, 0xc4, 0xe5, 0xc2, 0x9f, 0x6d, 0xe8, 0x1a, 0x04, 0x21,
	0xe8, 0x94, 0xb8, 0x98, 0x57, 0x61, 0xd6, 0xcb, 0x05, 0xb6, 0xaf, 0x38, 0x77, 0x0b, 0x40, 0xd0,
	0x6f, 0x24, 0x99, 0xcc, 0xa4, 0x3a, 0x68, 0x5d, 0xe5, 0x3b, 0x71, 0x4f, 0x23, 0x27, 0x1a, 0x40,
	0x77, 0x41, 0x5d, 0x3e, 0xe1, 0xa4, 0x60, 0x53, 0x3c, 0x51, 0x02, 0x1d, 0x53, 0x50, 0x9f, 0xaa,
	0x5b, 0x3a, 0x48, 0x1f, 0x2b, 0x67, 0x15, 0xf1, 0xbb, 0xf6, 0x58, 0xbd, 0xd6, 0x46, 0x08, 0xc2,
	0x29, 0xce, 0x93, 0xb2, 0x2e, 0x26, 0x84, 0xfb, 0x9e, 0x49, 0x0e, 0x2c, 0xf8, 0xde, 0x60, 0xe8,
	0x01, 0x20, 0x21, 0x19, 0x57, 0x7d, 0x95, 0xa4, 0xac, 0x94, 0x9c, 0xe5, 0xb9, 0x62, 0x6e, 0x18,
	0xe6, 0xd0, 0x65, 0x46, 0x4d, 0x02, 0x3d, 0x04, 0xa8, 0x30, 0x97, 0x54, 0x52, 0xe5, 0x8e, 0xbf,
	0x69, 0x2c, 0xd9, 0x75, 0x96, 0x8c, 0xe7, 0x89, 0x78, 0x81, 0x13, 0x7e, 0x6f, 0x43, 0xaf, 0xc9,
	0xfc, 0x17, 0x7b, 0xee, 0xc3, 0x0e, 0x9e, 0x62, 0x9a, 0x6b, 0x23, 0x1c, 0xa7, 0x63, 0x38, 0xdb,
	0x0d, 0x6c, 0x89, 0xb7, 0x01, 0xf4, 0x09, 0x62, 0x26, 0x24, 0x29, 0x9c, 0x55, 0x0b, 0xc8, 0x72,
	0x03, 0x7b, 0xcb, 0x0d, 0xac, 0x1b, 0x4e, 0xc9, 0x91, 0xdc, 0xf9, 0x63, 0x03, 0x74, 0xe0, 0x9e,
	0x07, 0x67, 0x09, 0x2b, 0xf3, 0x99, 0x72, 0x45, 0x3f, 0x0f, 0xe8, 0xe7, 0xc1, 0xd9, 0x07, 0x85,
	0x84, 0xcf, 0xc0, 0x53, 0x0d, 0x55, 0xe7, 0x12, 0x05, 0xb0, 0x49, 0x2e, 0xa9, 0x1c, 0xb1, 0xcc,
	0x7a, 0xd0, 0x8d, 0x9b, 0x18, 0xf9, 0xb0, 0x51, 0xa8, 0xf6, 0x54, 0x86, 0x3b, 0x0f, 0xe6, 0x61,
	0x78, 0x08, 0x1b, 0x23, 0x56, 0x14, 0xb8, 0xcc, 0x34, 0x29, 0xb5, 0x4b, 0xe7, 0xe1, 0x3c, 0x3c,
	0xfe, 0xd1, 0x06, 0x78, 0xa3, 0x7f, 0x12, 0x2f, 0xf4, 0x3f, 0x02, 0x8d, 0x00, 0x7e, 0xcf, 0x30,
	0xf2, 0xdd, 0x1b, 0x5d, 0x19, 0xeb, 0xe0, 0xe6, 0x8a, 0x8c, 0xed, 0xfd, 0x70, 0x0d, 0xbd, 0x85,
	0xc1, 0xe2, 0xe4, 0xa1, 0xc0, 0x91, 0x57, 0x8c, 0x71, 0xb0, 0xbf, 0x32, 0xd7, 0x48, 0x3d, 0x87,
	0x5e, 0x33, 0x5d, 0xe8, 0x86, 0xe3, 0x2e, 0xcf, 0x60, 0xe0, 0x5f, 0x4d, 0x34, 0x0a, 0x11, 0xf4,
	0x5f, 0x5d, 0x92, 0x74, 0xee, 0xc4, 0xb6, 0xa3, 0xba, 0x38, 0xd8, 0x72, 0xb1, 0x75, 0x3a, 0x5c,
	0x3b, 0x79, 0x02, 0xfb, 0x94, 0x45, 0x17, 0xbc, 0x4a, 0x23, 0x72, 0x89, 0x8b, 0x4a, 0xbd, 0x73,
	0xc4, 0x59, 0x2d, 0xc9, 0x45, 0x4d, 0x33, 0x72, 0xb2, 0x13, 0xeb, 0xf5, 0xa9, 0x5e, 0x8f, 0xf5,
	0xc6, 0x71, 0xeb, 0xb3, 0xfd, 0xa3, 0x4e, 0x3c, 0xf3, 0x79, 0xf4, 0x2b, 0x00, 0x00, 0xff, 0xff,
	0x50, 0x6a, 0x68, 0x2a, 0x74, 0x05, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// HurraAgentClient is the client API for HurraAgent service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type HurraAgentClient interface {
	MountDrive(ctx context.Context, in *MountDriveRequest, opts ...grpc.CallOption) (*MountDriveResponse, error)
	UnmountDrive(ctx context.Context, in *UnmountDriveRequest, opts ...grpc.CallOption) (*UnmountDriveResponse, error)
	GetDrives(ctx context.Context, in *GetDrivesRequest, opts ...grpc.CallOption) (*GetDrivesResponse, error)
	ExecCommand(ctx context.Context, in *Command, opts ...grpc.CallOption) (*Result, error)
}

type hurraAgentClient struct {
	cc grpc.ClientConnInterface
}

func NewHurraAgentClient(cc grpc.ClientConnInterface) HurraAgentClient {
	return &hurraAgentClient{cc}
}

func (c *hurraAgentClient) MountDrive(ctx context.Context, in *MountDriveRequest, opts ...grpc.CallOption) (*MountDriveResponse, error) {
	out := new(MountDriveResponse)
	err := c.cc.Invoke(ctx, "/proto.HurraAgent/MountDrive", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hurraAgentClient) UnmountDrive(ctx context.Context, in *UnmountDriveRequest, opts ...grpc.CallOption) (*UnmountDriveResponse, error) {
	out := new(UnmountDriveResponse)
	err := c.cc.Invoke(ctx, "/proto.HurraAgent/UnmountDrive", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hurraAgentClient) GetDrives(ctx context.Context, in *GetDrivesRequest, opts ...grpc.CallOption) (*GetDrivesResponse, error) {
	out := new(GetDrivesResponse)
	err := c.cc.Invoke(ctx, "/proto.HurraAgent/GetDrives", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hurraAgentClient) ExecCommand(ctx context.Context, in *Command, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := c.cc.Invoke(ctx, "/proto.HurraAgent/ExecCommand", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HurraAgentServer is the server API for HurraAgent service.
type HurraAgentServer interface {
	MountDrive(context.Context, *MountDriveRequest) (*MountDriveResponse, error)
	UnmountDrive(context.Context, *UnmountDriveRequest) (*UnmountDriveResponse, error)
	GetDrives(context.Context, *GetDrivesRequest) (*GetDrivesResponse, error)
	ExecCommand(context.Context, *Command) (*Result, error)
}

// UnimplementedHurraAgentServer can be embedded to have forward compatible implementations.
type UnimplementedHurraAgentServer struct {
}

func (*UnimplementedHurraAgentServer) MountDrive(ctx context.Context, req *MountDriveRequest) (*MountDriveResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MountDrive not implemented")
}
func (*UnimplementedHurraAgentServer) UnmountDrive(ctx context.Context, req *UnmountDriveRequest) (*UnmountDriveResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnmountDrive not implemented")
}
func (*UnimplementedHurraAgentServer) GetDrives(ctx context.Context, req *GetDrivesRequest) (*GetDrivesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDrives not implemented")
}
func (*UnimplementedHurraAgentServer) ExecCommand(ctx context.Context, req *Command) (*Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExecCommand not implemented")
}

func RegisterHurraAgentServer(s *grpc.Server, srv HurraAgentServer) {
	s.RegisterService(&_HurraAgent_serviceDesc, srv)
}

func _HurraAgent_MountDrive_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MountDriveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HurraAgentServer).MountDrive(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.HurraAgent/MountDrive",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HurraAgentServer).MountDrive(ctx, req.(*MountDriveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HurraAgent_UnmountDrive_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnmountDriveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HurraAgentServer).UnmountDrive(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.HurraAgent/UnmountDrive",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HurraAgentServer).UnmountDrive(ctx, req.(*UnmountDriveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HurraAgent_GetDrives_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDrivesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HurraAgentServer).GetDrives(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.HurraAgent/GetDrives",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HurraAgentServer).GetDrives(ctx, req.(*GetDrivesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HurraAgent_ExecCommand_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Command)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HurraAgentServer).ExecCommand(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.HurraAgent/ExecCommand",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HurraAgentServer).ExecCommand(ctx, req.(*Command))
	}
	return interceptor(ctx, in, info, handler)
}

var _HurraAgent_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.HurraAgent",
	HandlerType: (*HurraAgentServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "MountDrive",
			Handler:    _HurraAgent_MountDrive_Handler,
		},
		{
			MethodName: "UnmountDrive",
			Handler:    _HurraAgent_UnmountDrive_Handler,
		},
		{
			MethodName: "GetDrives",
			Handler:    _HurraAgent_GetDrives_Handler,
		},
		{
			MethodName: "ExecCommand",
			Handler:    _HurraAgent_ExecCommand_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "hurra_agent.proto",
}
