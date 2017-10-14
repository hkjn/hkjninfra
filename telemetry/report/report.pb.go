// Code generated by protoc-gen-go. DO NOT EDIT.
// source: report.proto

/*
Package report is a generated protocol buffer package.

It is generated from these files:
	report.proto

It has these top-level messages:
	ReportRequest
	ReportResponse
	InfoRequest
	DiskInfo
	ClientInfo
	InfoResponse
*/
package report

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "google.golang.org/genproto/googleapis/api/annotations"
import google_protobuf1 "github.com/golang/protobuf/ptypes/timestamp"

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

// ReportRequest describes the request to report in from a client.
type ReportRequest struct {
	Ts   *google_protobuf1.Timestamp `protobuf:"bytes,2,opt,name=ts" json:"ts,omitempty"`
	Info *ClientInfo                 `protobuf:"bytes,3,opt,name=info" json:"info,omitempty"`
}

func (m *ReportRequest) Reset()                    { *m = ReportRequest{} }
func (m *ReportRequest) String() string            { return proto.CompactTextString(m) }
func (*ReportRequest) ProtoMessage()               {}
func (*ReportRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *ReportRequest) GetTs() *google_protobuf1.Timestamp {
	if m != nil {
		return m.Ts
	}
	return nil
}

func (m *ReportRequest) GetInfo() *ClientInfo {
	if m != nil {
		return m.Info
	}
	return nil
}

// ReportResponse describes the response from the server when a client reports in.
type ReportResponse struct {
	Message string `protobuf:"bytes,1,opt,name=message" json:"message,omitempty"`
}

func (m *ReportResponse) Reset()                    { *m = ReportResponse{} }
func (m *ReportResponse) String() string            { return proto.CompactTextString(m) }
func (*ReportResponse) ProtoMessage()               {}
func (*ReportResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *ReportResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

// InfoRequest describes a request to look up info on known clients.
type InfoRequest struct {
}

func (m *InfoRequest) Reset()                    { *m = InfoRequest{} }
func (m *InfoRequest) String() string            { return proto.CompactTextString(m) }
func (*InfoRequest) ProtoMessage()               {}
func (*InfoRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

// DiskInfo describes info on one disk partition.
type DiskInfo struct {
	Source      string `protobuf:"bytes,1,opt,name=source" json:"source,omitempty"`
	Size        string `protobuf:"bytes,2,opt,name=size" json:"size,omitempty"`
	PercentUsed string `protobuf:"bytes,3,opt,name=percent_used,json=percentUsed" json:"percent_used,omitempty"`
	Target      string `protobuf:"bytes,4,opt,name=target" json:"target,omitempty"`
}

func (m *DiskInfo) Reset()                    { *m = DiskInfo{} }
func (m *DiskInfo) String() string            { return proto.CompactTextString(m) }
func (*DiskInfo) ProtoMessage()               {}
func (*DiskInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *DiskInfo) GetSource() string {
	if m != nil {
		return m.Source
	}
	return ""
}

func (m *DiskInfo) GetSize() string {
	if m != nil {
		return m.Size
	}
	return ""
}

func (m *DiskInfo) GetPercentUsed() string {
	if m != nil {
		return m.PercentUsed
	}
	return ""
}

func (m *DiskInfo) GetTarget() string {
	if m != nil {
		return m.Target
	}
	return ""
}

// ClientInfo describes info for one client.
type ClientInfo struct {
	Id              string      `protobuf:"bytes,13,opt,name=id" json:"id,omitempty"`
	AllowedSshKeys  string      `protobuf:"bytes,1,opt,name=allowed_ssh_keys,json=allowedSshKeys" json:"allowed_ssh_keys,omitempty"`
	CpuArch         string      `protobuf:"bytes,2,opt,name=cpu_arch,json=cpuArch" json:"cpu_arch,omitempty"`
	Disks           []*DiskInfo `protobuf:"bytes,3,rep,name=disks" json:"disks,omitempty"`
	Hostname        string      `protobuf:"bytes,4,opt,name=hostname" json:"hostname,omitempty"`
	KernelName      string      `protobuf:"bytes,5,opt,name=kernel_name,json=kernelName" json:"kernel_name,omitempty"`
	KernelVersion   string      `protobuf:"bytes,6,opt,name=kernel_version,json=kernelVersion" json:"kernel_version,omitempty"`
	CpuArchitecture string      `protobuf:"bytes,7,opt,name=cpu_architecture,json=cpuArchitecture" json:"cpu_architecture,omitempty"`
	Platform        string      `protobuf:"bytes,8,opt,name=platform" json:"platform,omitempty"`
	MemoryTotalMb   string      `protobuf:"bytes,9,opt,name=memory_total_mb,json=memoryTotalMb" json:"memory_total_mb,omitempty"`
	MemoryAvailMb   string      `protobuf:"bytes,10,opt,name=memory_avail_mb,json=memoryAvailMb" json:"memory_avail_mb,omitempty"`
	Tags            []string    `protobuf:"bytes,11,rep,name=tags" json:"tags,omitempty"`
	Zone            string      `protobuf:"bytes,12,opt,name=zone" json:"zone,omitempty"`
}

func (m *ClientInfo) Reset()                    { *m = ClientInfo{} }
func (m *ClientInfo) String() string            { return proto.CompactTextString(m) }
func (*ClientInfo) ProtoMessage()               {}
func (*ClientInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *ClientInfo) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *ClientInfo) GetAllowedSshKeys() string {
	if m != nil {
		return m.AllowedSshKeys
	}
	return ""
}

func (m *ClientInfo) GetCpuArch() string {
	if m != nil {
		return m.CpuArch
	}
	return ""
}

func (m *ClientInfo) GetDisks() []*DiskInfo {
	if m != nil {
		return m.Disks
	}
	return nil
}

func (m *ClientInfo) GetHostname() string {
	if m != nil {
		return m.Hostname
	}
	return ""
}

func (m *ClientInfo) GetKernelName() string {
	if m != nil {
		return m.KernelName
	}
	return ""
}

func (m *ClientInfo) GetKernelVersion() string {
	if m != nil {
		return m.KernelVersion
	}
	return ""
}

func (m *ClientInfo) GetCpuArchitecture() string {
	if m != nil {
		return m.CpuArchitecture
	}
	return ""
}

func (m *ClientInfo) GetPlatform() string {
	if m != nil {
		return m.Platform
	}
	return ""
}

func (m *ClientInfo) GetMemoryTotalMb() string {
	if m != nil {
		return m.MemoryTotalMb
	}
	return ""
}

func (m *ClientInfo) GetMemoryAvailMb() string {
	if m != nil {
		return m.MemoryAvailMb
	}
	return ""
}

func (m *ClientInfo) GetTags() []string {
	if m != nil {
		return m.Tags
	}
	return nil
}

func (m *ClientInfo) GetZone() string {
	if m != nil {
		return m.Zone
	}
	return ""
}

// InfoResponse describes a response for info on known clients.
type InfoResponse struct {
	// The info field describes each known client and their info.
	Info map[string]*ClientInfo `protobuf:"bytes,1,rep,name=info" json:"info,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *InfoResponse) Reset()                    { *m = InfoResponse{} }
func (m *InfoResponse) String() string            { return proto.CompactTextString(m) }
func (*InfoResponse) ProtoMessage()               {}
func (*InfoResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *InfoResponse) GetInfo() map[string]*ClientInfo {
	if m != nil {
		return m.Info
	}
	return nil
}

func init() {
	proto.RegisterType((*ReportRequest)(nil), "report.ReportRequest")
	proto.RegisterType((*ReportResponse)(nil), "report.ReportResponse")
	proto.RegisterType((*InfoRequest)(nil), "report.InfoRequest")
	proto.RegisterType((*DiskInfo)(nil), "report.DiskInfo")
	proto.RegisterType((*ClientInfo)(nil), "report.ClientInfo")
	proto.RegisterType((*InfoResponse)(nil), "report.InfoResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Report service

type ReportClient interface {
	// Send report to server.
	Send(ctx context.Context, in *ReportRequest, opts ...grpc.CallOption) (*ReportResponse, error)
	// Query for info on known clients.
	Info(ctx context.Context, in *InfoRequest, opts ...grpc.CallOption) (*InfoResponse, error)
}

type reportClient struct {
	cc *grpc.ClientConn
}

func NewReportClient(cc *grpc.ClientConn) ReportClient {
	return &reportClient{cc}
}

func (c *reportClient) Send(ctx context.Context, in *ReportRequest, opts ...grpc.CallOption) (*ReportResponse, error) {
	out := new(ReportResponse)
	err := grpc.Invoke(ctx, "/report.Report/Send", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *reportClient) Info(ctx context.Context, in *InfoRequest, opts ...grpc.CallOption) (*InfoResponse, error) {
	out := new(InfoResponse)
	err := grpc.Invoke(ctx, "/report.Report/Info", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Report service

type ReportServer interface {
	// Send report to server.
	Send(context.Context, *ReportRequest) (*ReportResponse, error)
	// Query for info on known clients.
	Info(context.Context, *InfoRequest) (*InfoResponse, error)
}

func RegisterReportServer(s *grpc.Server, srv ReportServer) {
	s.RegisterService(&_Report_serviceDesc, srv)
}

func _Report_Send_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReportRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReportServer).Send(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/report.Report/Send",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReportServer).Send(ctx, req.(*ReportRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Report_Info_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReportServer).Info(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/report.Report/Info",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReportServer).Info(ctx, req.(*InfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Report_serviceDesc = grpc.ServiceDesc{
	ServiceName: "report.Report",
	HandlerType: (*ReportServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Send",
			Handler:    _Report_Send_Handler,
		},
		{
			MethodName: "Info",
			Handler:    _Report_Info_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "report.proto",
}

func init() { proto.RegisterFile("report.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 645 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x53, 0xcf, 0x6e, 0xd3, 0x4e,
	0x10, 0xfe, 0x39, 0x49, 0xd3, 0x64, 0xf2, 0xa7, 0xd1, 0xfe, 0x4a, 0xe5, 0x46, 0x88, 0x06, 0x4b,
	0x54, 0xa1, 0x07, 0x47, 0xa4, 0x17, 0xd4, 0x5b, 0x0b, 0x1c, 0x50, 0x55, 0x84, 0xdc, 0xd2, 0x6b,
	0xb4, 0x71, 0x26, 0xc9, 0x12, 0x7b, 0xd7, 0xdd, 0x5d, 0x07, 0xa5, 0x47, 0x5e, 0x80, 0x03, 0xe2,
	0x2d, 0x78, 0x1b, 0x5e, 0x81, 0x07, 0x41, 0xde, 0x5d, 0x97, 0x16, 0xf5, 0x62, 0xed, 0xf7, 0xcd,
	0x37, 0xe3, 0x6f, 0x76, 0x66, 0xa1, 0x2d, 0x31, 0x13, 0x52, 0x87, 0x99, 0x14, 0x5a, 0x90, 0xba,
	0x45, 0xfd, 0xa7, 0x0b, 0x21, 0x16, 0x09, 0x8e, 0x68, 0xc6, 0x46, 0x94, 0x73, 0xa1, 0xa9, 0x66,
	0x82, 0x2b, 0xab, 0xea, 0x1f, 0xb8, 0xa8, 0x41, 0xd3, 0x7c, 0x3e, 0xd2, 0x2c, 0x45, 0xa5, 0x69,
	0x9a, 0x59, 0x41, 0x10, 0x43, 0x27, 0x32, 0x85, 0x22, 0xbc, 0xc9, 0x51, 0x69, 0x72, 0x04, 0x15,
	0xad, 0xfc, 0xca, 0xc0, 0x1b, 0xb6, 0xc6, 0xfd, 0xd0, 0xa6, 0x87, 0x65, 0x7a, 0x78, 0x55, 0xa6,
	0x47, 0x15, 0xad, 0xc8, 0x21, 0xd4, 0x18, 0x9f, 0x0b, 0xbf, 0x6a, 0xd4, 0x24, 0x74, 0x06, 0xdf,
	0x24, 0x0c, 0xb9, 0x7e, 0xcf, 0xe7, 0x22, 0x32, 0xf1, 0xe0, 0x08, 0xba, 0xe5, 0x4f, 0x54, 0x26,
	0xb8, 0x42, 0xe2, 0xc3, 0x76, 0x8a, 0x4a, 0xd1, 0x05, 0xfa, 0xde, 0xc0, 0x1b, 0x36, 0xa3, 0x12,
	0x06, 0x1d, 0x68, 0x99, 0x4c, 0x6b, 0x27, 0xb8, 0x81, 0xc6, 0x5b, 0xa6, 0x56, 0x05, 0x45, 0xf6,
	0xa0, 0xae, 0x44, 0x2e, 0xe3, 0x32, 0xc7, 0x21, 0x42, 0xa0, 0xa6, 0xd8, 0x2d, 0x1a, 0xd3, 0xcd,
	0xc8, 0x9c, 0xc9, 0x73, 0x68, 0x67, 0x28, 0x63, 0xe4, 0x7a, 0x92, 0x2b, 0x9c, 0x19, 0x8b, 0xcd,
	0xa8, 0xe5, 0xb8, 0x4f, 0x0a, 0x67, 0x45, 0x39, 0x4d, 0xe5, 0x02, 0xb5, 0x5f, 0xb3, 0xe5, 0x2c,
	0x0a, 0x7e, 0x56, 0x01, 0xfe, 0xb6, 0x40, 0xba, 0x50, 0x61, 0x33, 0xbf, 0x63, 0x24, 0x15, 0x36,
	0x23, 0x43, 0xe8, 0xd1, 0x24, 0x11, 0x5f, 0x70, 0x36, 0x51, 0x6a, 0x39, 0x59, 0xe1, 0x46, 0x39,
	0x3f, 0x5d, 0xc7, 0x5f, 0xaa, 0xe5, 0x39, 0x6e, 0x14, 0xd9, 0x87, 0x46, 0x9c, 0xe5, 0x13, 0x2a,
	0xe3, 0xa5, 0xf3, 0xb6, 0x1d, 0x67, 0xf9, 0xa9, 0x8c, 0x97, 0xe4, 0x10, 0xb6, 0x66, 0x4c, 0xad,
	0x94, 0x5f, 0x1d, 0x54, 0x87, 0xad, 0x71, 0xaf, 0xbc, 0xba, 0xb2, 0xd7, 0xc8, 0x86, 0x49, 0x1f,
	0x1a, 0x4b, 0xa1, 0x34, 0xa7, 0x29, 0x3a, 0x97, 0x77, 0x98, 0x1c, 0x40, 0x6b, 0x85, 0x92, 0x63,
	0x32, 0x31, 0xe1, 0x2d, 0x13, 0x06, 0x4b, 0x7d, 0x28, 0x04, 0x2f, 0xa0, 0xeb, 0x04, 0x6b, 0x94,
	0x8a, 0x09, 0xee, 0xd7, 0x8d, 0xa6, 0x63, 0xd9, 0x6b, 0x4b, 0x92, 0x97, 0xd0, 0x2b, 0x6d, 0x32,
	0x8d, 0xb1, 0xce, 0x25, 0xfa, 0xdb, 0x46, 0xb8, 0xe3, 0xec, 0x96, 0x74, 0x61, 0x27, 0x4b, 0xa8,
	0x9e, 0x0b, 0x99, 0xfa, 0x0d, 0x6b, 0xa7, 0xc4, 0xe4, 0x10, 0x76, 0x52, 0x4c, 0x85, 0xdc, 0x4c,
	0xb4, 0xd0, 0x34, 0x99, 0xa4, 0x53, 0xbf, 0x69, 0x7f, 0x67, 0xe9, 0xab, 0x82, 0xbd, 0x98, 0xde,
	0xd3, 0xd1, 0x35, 0x65, 0x46, 0x07, 0xf7, 0x75, 0xa7, 0x05, 0x7b, 0x31, 0x2d, 0xa6, 0xaa, 0xe9,
	0x42, 0xf9, 0xad, 0x41, 0xb5, 0x98, 0x6a, 0x71, 0x2e, 0xb8, 0x5b, 0xc1, 0xd1, 0x6f, 0xdb, 0x49,
	0x17, 0xe7, 0xe0, 0x9b, 0x07, 0x6d, 0xbb, 0x31, 0x6e, 0xb7, 0xc6, 0x6e, 0x2b, 0x3d, 0x73, 0xb5,
	0xcf, 0xca, 0xab, 0xbd, 0xaf, 0x31, 0xe0, 0x1d, 0xd7, 0x72, 0x63, 0x37, 0xb4, 0x7f, 0x0e, 0xcd,
	0x3b, 0x8a, 0xf4, 0xa0, 0xba, 0xc2, 0x8d, 0x1b, 0x6a, 0x71, 0x24, 0x43, 0xd8, 0x5a, 0xd3, 0x24,
	0x47, 0xf7, 0x2e, 0x1e, 0xdb, 0x74, 0x2b, 0x38, 0xa9, 0xbc, 0xf6, 0xc6, 0x3f, 0x3c, 0xa8, 0xdb,
	0x7d, 0x27, 0xd7, 0x50, 0xbb, 0x44, 0x3e, 0x23, 0x4f, 0xca, 0x8c, 0x07, 0x8f, 0xad, 0xbf, 0xf7,
	0x2f, 0x6d, 0xed, 0x05, 0x07, 0x5f, 0x7f, 0xfd, 0xfe, 0x5e, 0xd9, 0x0f, 0x76, 0x47, 0xeb, 0x57,
	0x23, 0x8d, 0x09, 0xa6, 0xa8, 0xe5, 0x66, 0x64, 0xc5, 0x27, 0xde, 0x11, 0x39, 0x86, 0x9a, 0x59,
	0xce, 0xff, 0x1f, 0x76, 0x67, 0xab, 0xee, 0x3e, 0xd6, 0x72, 0xf0, 0xdf, 0x59, 0x08, 0x41, 0x8a,
	0xe1, 0x72, 0xf5, 0x99, 0x9b, 0x0f, 0xe3, 0x73, 0x49, 0xc3, 0xbb, 0xea, 0x2e, 0xe9, 0xcc, 0x59,
	0xff, 0xe8, 0x4d, 0xeb, 0xe6, 0xd9, 0x1f, 0xff, 0x09, 0x00, 0x00, 0xff, 0xff, 0xa7, 0x53, 0x5a,
	0xc5, 0x79, 0x04, 0x00, 0x00,
}
