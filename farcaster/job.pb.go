// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v5.29.3
// source: job.proto

package farcaster

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type RevokeMessagesBySignerJobPayload struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Fid           uint32                 `protobuf:"varint,1,opt,name=fid,proto3" json:"fid,omitempty"`
	Signer        []byte                 `protobuf:"bytes,2,opt,name=signer,proto3" json:"signer,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RevokeMessagesBySignerJobPayload) Reset() {
	*x = RevokeMessagesBySignerJobPayload{}
	mi := &file_job_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RevokeMessagesBySignerJobPayload) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RevokeMessagesBySignerJobPayload) ProtoMessage() {}

func (x *RevokeMessagesBySignerJobPayload) ProtoReflect() protoreflect.Message {
	mi := &file_job_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RevokeMessagesBySignerJobPayload.ProtoReflect.Descriptor instead.
func (*RevokeMessagesBySignerJobPayload) Descriptor() ([]byte, []int) {
	return file_job_proto_rawDescGZIP(), []int{0}
}

func (x *RevokeMessagesBySignerJobPayload) GetFid() uint32 {
	if x != nil {
		return x.Fid
	}
	return 0
}

func (x *RevokeMessagesBySignerJobPayload) GetSigner() []byte {
	if x != nil {
		return x.Signer
	}
	return nil
}

type UpdateNameRegistryEventExpiryJobPayload struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Fname         []byte                 `protobuf:"bytes,1,opt,name=fname,proto3" json:"fname,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateNameRegistryEventExpiryJobPayload) Reset() {
	*x = UpdateNameRegistryEventExpiryJobPayload{}
	mi := &file_job_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateNameRegistryEventExpiryJobPayload) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateNameRegistryEventExpiryJobPayload) ProtoMessage() {}

func (x *UpdateNameRegistryEventExpiryJobPayload) ProtoReflect() protoreflect.Message {
	mi := &file_job_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateNameRegistryEventExpiryJobPayload.ProtoReflect.Descriptor instead.
func (*UpdateNameRegistryEventExpiryJobPayload) Descriptor() ([]byte, []int) {
	return file_job_proto_rawDescGZIP(), []int{1}
}

func (x *UpdateNameRegistryEventExpiryJobPayload) GetFname() []byte {
	if x != nil {
		return x.Fname
	}
	return nil
}

var File_job_proto protoreflect.FileDescriptor

var file_job_proto_rawDesc = string([]byte{
	0x0a, 0x09, 0x6a, 0x6f, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x4c, 0x0a, 0x20, 0x52,
	0x65, 0x76, 0x6f, 0x6b, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x42, 0x79, 0x53,
	0x69, 0x67, 0x6e, 0x65, 0x72, 0x4a, 0x6f, 0x62, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x12,
	0x10, 0x0a, 0x03, 0x66, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x03, 0x66, 0x69,
	0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x06, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x72, 0x22, 0x3f, 0x0a, 0x27, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x45, 0x78, 0x70, 0x69, 0x72, 0x79, 0x4a, 0x6f, 0x62, 0x50, 0x61, 0x79,
	0x6c, 0x6f, 0x61, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x66, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x05, 0x66, 0x6e, 0x61, 0x6d, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
})

var (
	file_job_proto_rawDescOnce sync.Once
	file_job_proto_rawDescData []byte
)

func file_job_proto_rawDescGZIP() []byte {
	file_job_proto_rawDescOnce.Do(func() {
		file_job_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_job_proto_rawDesc), len(file_job_proto_rawDesc)))
	})
	return file_job_proto_rawDescData
}

var file_job_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_job_proto_goTypes = []any{
	(*RevokeMessagesBySignerJobPayload)(nil),        // 0: RevokeMessagesBySignerJobPayload
	(*UpdateNameRegistryEventExpiryJobPayload)(nil), // 1: UpdateNameRegistryEventExpiryJobPayload
}
var file_job_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_job_proto_init() }
func file_job_proto_init() {
	if File_job_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_job_proto_rawDesc), len(file_job_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_job_proto_goTypes,
		DependencyIndexes: file_job_proto_depIdxs,
		MessageInfos:      file_job_proto_msgTypes,
	}.Build()
	File_job_proto = out.File
	file_job_proto_goTypes = nil
	file_job_proto_depIdxs = nil
}
