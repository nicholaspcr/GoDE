// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: api/v1/generation.proto

package api

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type GenerationIDs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GenerationIDs) Reset() {
	*x = GenerationIDs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_generation_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GenerationIDs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GenerationIDs) ProtoMessage() {}

func (x *GenerationIDs) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_generation_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GenerationIDs.ProtoReflect.Descriptor instead.
func (*GenerationIDs) Descriptor() ([]byte, []int) {
	return file_api_v1_generation_proto_rawDescGZIP(), []int{0}
}

func (x *GenerationIDs) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type Generation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          uint64        `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	ExecutionId uint64        `protobuf:"varint,2,opt,name=execution_id,json=executionId,proto3" json:"execution_id,omitempty"`
	Populations []*Population `protobuf:"bytes,3,rep,name=populations,proto3" json:"populations,omitempty"`
}

func (x *Generation) Reset() {
	*x = Generation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_generation_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Generation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Generation) ProtoMessage() {}

func (x *Generation) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_generation_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Generation.ProtoReflect.Descriptor instead.
func (*Generation) Descriptor() ([]byte, []int) {
	return file_api_v1_generation_proto_rawDescGZIP(), []int{1}
}

func (x *Generation) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Generation) GetExecutionId() uint64 {
	if x != nil {
		return x.ExecutionId
	}
	return 0
}

func (x *Generation) GetPopulations() []*Population {
	if x != nil {
		return x.Populations
	}
	return nil
}

type GenerationServiceCreateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Generation *Generation `protobuf:"bytes,1,opt,name=generation,proto3" json:"generation,omitempty"`
}

func (x *GenerationServiceCreateRequest) Reset() {
	*x = GenerationServiceCreateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_generation_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GenerationServiceCreateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GenerationServiceCreateRequest) ProtoMessage() {}

func (x *GenerationServiceCreateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_generation_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GenerationServiceCreateRequest.ProtoReflect.Descriptor instead.
func (*GenerationServiceCreateRequest) Descriptor() ([]byte, []int) {
	return file_api_v1_generation_proto_rawDescGZIP(), []int{2}
}

func (x *GenerationServiceCreateRequest) GetGeneration() *Generation {
	if x != nil {
		return x.Generation
	}
	return nil
}

type GenerationServiceGetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GenerationIds *GenerationIDs `protobuf:"bytes,1,opt,name=generation_ids,json=generationIds,proto3" json:"generation_ids,omitempty"`
}

func (x *GenerationServiceGetRequest) Reset() {
	*x = GenerationServiceGetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_generation_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GenerationServiceGetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GenerationServiceGetRequest) ProtoMessage() {}

func (x *GenerationServiceGetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_generation_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GenerationServiceGetRequest.ProtoReflect.Descriptor instead.
func (*GenerationServiceGetRequest) Descriptor() ([]byte, []int) {
	return file_api_v1_generation_proto_rawDescGZIP(), []int{3}
}

func (x *GenerationServiceGetRequest) GetGenerationIds() *GenerationIDs {
	if x != nil {
		return x.GenerationIds
	}
	return nil
}

type GenerationServiceGetResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Generation *Generation `protobuf:"bytes,1,opt,name=generation,proto3" json:"generation,omitempty"`
}

func (x *GenerationServiceGetResponse) Reset() {
	*x = GenerationServiceGetResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_generation_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GenerationServiceGetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GenerationServiceGetResponse) ProtoMessage() {}

func (x *GenerationServiceGetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_generation_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GenerationServiceGetResponse.ProtoReflect.Descriptor instead.
func (*GenerationServiceGetResponse) Descriptor() ([]byte, []int) {
	return file_api_v1_generation_proto_rawDescGZIP(), []int{4}
}

func (x *GenerationServiceGetResponse) GetGeneration() *Generation {
	if x != nil {
		return x.Generation
	}
	return nil
}

type GenerationServiceUpdateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Generation *Generation `protobuf:"bytes,1,opt,name=generation,proto3" json:"generation,omitempty"`
}

func (x *GenerationServiceUpdateRequest) Reset() {
	*x = GenerationServiceUpdateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_generation_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GenerationServiceUpdateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GenerationServiceUpdateRequest) ProtoMessage() {}

func (x *GenerationServiceUpdateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_generation_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GenerationServiceUpdateRequest.ProtoReflect.Descriptor instead.
func (*GenerationServiceUpdateRequest) Descriptor() ([]byte, []int) {
	return file_api_v1_generation_proto_rawDescGZIP(), []int{5}
}

func (x *GenerationServiceUpdateRequest) GetGeneration() *Generation {
	if x != nil {
		return x.Generation
	}
	return nil
}

type GenerationServiceDeleteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GenerationIds *GenerationIDs `protobuf:"bytes,1,opt,name=generation_ids,json=generationIds,proto3" json:"generation_ids,omitempty"`
}

func (x *GenerationServiceDeleteRequest) Reset() {
	*x = GenerationServiceDeleteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_generation_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GenerationServiceDeleteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GenerationServiceDeleteRequest) ProtoMessage() {}

func (x *GenerationServiceDeleteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_generation_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GenerationServiceDeleteRequest.ProtoReflect.Descriptor instead.
func (*GenerationServiceDeleteRequest) Descriptor() ([]byte, []int) {
	return file_api_v1_generation_proto_rawDescGZIP(), []int{6}
}

func (x *GenerationServiceDeleteRequest) GetGenerationIds() *GenerationIDs {
	if x != nil {
		return x.GenerationIds
	}
	return nil
}

var File_api_v1_generation_proto protoreflect.FileDescriptor

var file_api_v1_generation_proto_rawDesc = []byte{
	0x0a, 0x17, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x61, 0x70, 0x69, 0x2e, 0x76,
	0x31, 0x1a, 0x17, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x6f, 0x70, 0x75, 0x6c, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74,
	0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x1f, 0x0a, 0x0d, 0x47, 0x65, 0x6e, 0x65, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x44, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x22, 0x75, 0x0a, 0x0a, 0x47, 0x65, 0x6e, 0x65,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x21, 0x0a, 0x0c, 0x65, 0x78, 0x65, 0x63, 0x75, 0x74,
	0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x65, 0x78,
	0x65, 0x63, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x34, 0x0a, 0x0b, 0x70, 0x6f, 0x70,
	0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x6f, 0x70, 0x75, 0x6c, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x0b, 0x70, 0x6f, 0x70, 0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22,
	0x54, 0x0a, 0x1e, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x32, 0x0a, 0x0a, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x47,
	0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0a, 0x67, 0x65, 0x6e, 0x65, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x5b, 0x0a, 0x1b, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x3c, 0x0a, 0x0e, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x49, 0x44, 0x73, 0x52, 0x0d, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49,
	0x64, 0x73, 0x22, 0x52, 0x0a, 0x1c, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x32, 0x0a, 0x0a, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e,
	0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0a, 0x67, 0x65, 0x6e, 0x65,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x54, 0x0a, 0x1e, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x32, 0x0a, 0x0a, 0x67, 0x65, 0x6e, 0x65,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x0a, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x5e, 0x0a, 0x1e,
	0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x3c,
	0x0a, 0x0e, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x73,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e,
	0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x44, 0x73, 0x52, 0x0d, 0x67,
	0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x73, 0x32, 0xcb, 0x02, 0x0a,
	0x11, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x4a, 0x0a, 0x06, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x26, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x52,
	0x0a, 0x03, 0x47, 0x65, 0x74, 0x12, 0x23, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x47,
	0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x12, 0x4a, 0x0a, 0x06, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x26, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x4a,
	0x0a, 0x06, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x26, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76,
	0x31, 0x2e, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x42, 0x09, 0x5a, 0x07, 0x70, 0x6b,
	0x67, 0x2f, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_v1_generation_proto_rawDescOnce sync.Once
	file_api_v1_generation_proto_rawDescData = file_api_v1_generation_proto_rawDesc
)

func file_api_v1_generation_proto_rawDescGZIP() []byte {
	file_api_v1_generation_proto_rawDescOnce.Do(func() {
		file_api_v1_generation_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_v1_generation_proto_rawDescData)
	})
	return file_api_v1_generation_proto_rawDescData
}

var file_api_v1_generation_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_api_v1_generation_proto_goTypes = []interface{}{
	(*GenerationIDs)(nil),                  // 0: api.v1.GenerationIDs
	(*Generation)(nil),                     // 1: api.v1.Generation
	(*GenerationServiceCreateRequest)(nil), // 2: api.v1.GenerationServiceCreateRequest
	(*GenerationServiceGetRequest)(nil),    // 3: api.v1.GenerationServiceGetRequest
	(*GenerationServiceGetResponse)(nil),   // 4: api.v1.GenerationServiceGetResponse
	(*GenerationServiceUpdateRequest)(nil), // 5: api.v1.GenerationServiceUpdateRequest
	(*GenerationServiceDeleteRequest)(nil), // 6: api.v1.GenerationServiceDeleteRequest
	(*Population)(nil),                     // 7: api.v1.Population
	(*emptypb.Empty)(nil),                  // 8: google.protobuf.Empty
}
var file_api_v1_generation_proto_depIdxs = []int32{
	7,  // 0: api.v1.Generation.populations:type_name -> api.v1.Population
	1,  // 1: api.v1.GenerationServiceCreateRequest.generation:type_name -> api.v1.Generation
	0,  // 2: api.v1.GenerationServiceGetRequest.generation_ids:type_name -> api.v1.GenerationIDs
	1,  // 3: api.v1.GenerationServiceGetResponse.generation:type_name -> api.v1.Generation
	1,  // 4: api.v1.GenerationServiceUpdateRequest.generation:type_name -> api.v1.Generation
	0,  // 5: api.v1.GenerationServiceDeleteRequest.generation_ids:type_name -> api.v1.GenerationIDs
	2,  // 6: api.v1.GenerationService.Create:input_type -> api.v1.GenerationServiceCreateRequest
	3,  // 7: api.v1.GenerationService.Get:input_type -> api.v1.GenerationServiceGetRequest
	5,  // 8: api.v1.GenerationService.Update:input_type -> api.v1.GenerationServiceUpdateRequest
	6,  // 9: api.v1.GenerationService.Delete:input_type -> api.v1.GenerationServiceDeleteRequest
	8,  // 10: api.v1.GenerationService.Create:output_type -> google.protobuf.Empty
	4,  // 11: api.v1.GenerationService.Get:output_type -> api.v1.GenerationServiceGetResponse
	8,  // 12: api.v1.GenerationService.Update:output_type -> google.protobuf.Empty
	8,  // 13: api.v1.GenerationService.Delete:output_type -> google.protobuf.Empty
	10, // [10:14] is the sub-list for method output_type
	6,  // [6:10] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_api_v1_generation_proto_init() }
func file_api_v1_generation_proto_init() {
	if File_api_v1_generation_proto != nil {
		return
	}
	file_api_v1_population_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_api_v1_generation_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GenerationIDs); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_v1_generation_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Generation); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_v1_generation_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GenerationServiceCreateRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_v1_generation_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GenerationServiceGetRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_v1_generation_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GenerationServiceGetResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_v1_generation_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GenerationServiceUpdateRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_v1_generation_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GenerationServiceDeleteRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_v1_generation_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_v1_generation_proto_goTypes,
		DependencyIndexes: file_api_v1_generation_proto_depIdxs,
		MessageInfos:      file_api_v1_generation_proto_msgTypes,
	}.Build()
	File_api_v1_generation_proto = out.File
	file_api_v1_generation_proto_rawDesc = nil
	file_api_v1_generation_proto_goTypes = nil
	file_api_v1_generation_proto_depIdxs = nil
}
