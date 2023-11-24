// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.25.1
// source: api/execution.proto

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

type ExecutionIDs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	SetId uint64 `protobuf:"varint,2,opt,name=set_id,json=setId,proto3" json:"set_id,omitempty"`
}

func (x *ExecutionIDs) Reset() {
	*x = ExecutionIDs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_execution_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExecutionIDs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExecutionIDs) ProtoMessage() {}

func (x *ExecutionIDs) ProtoReflect() protoreflect.Message {
	mi := &file_api_execution_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExecutionIDs.ProtoReflect.Descriptor instead.
func (*ExecutionIDs) Descriptor() ([]byte, []int) {
	return file_api_execution_proto_rawDescGZIP(), []int{0}
}

func (x *ExecutionIDs) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *ExecutionIDs) GetSetId() uint64 {
	if x != nil {
		return x.SetId
	}
	return 0
}

type Execution struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ids         *ExecutionIDs `protobuf:"bytes,1,opt,name=ids,proto3" json:"ids,omitempty"`
	Generations []*Generation `protobuf:"bytes,2,rep,name=generations,proto3" json:"generations,omitempty"`
	Pareto      *Population   `protobuf:"bytes,3,opt,name=pareto,proto3" json:"pareto,omitempty"`
}

func (x *Execution) Reset() {
	*x = Execution{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_execution_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Execution) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Execution) ProtoMessage() {}

func (x *Execution) ProtoReflect() protoreflect.Message {
	mi := &file_api_execution_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Execution.ProtoReflect.Descriptor instead.
func (*Execution) Descriptor() ([]byte, []int) {
	return file_api_execution_proto_rawDescGZIP(), []int{1}
}

func (x *Execution) GetIds() *ExecutionIDs {
	if x != nil {
		return x.Ids
	}
	return nil
}

func (x *Execution) GetGenerations() []*Generation {
	if x != nil {
		return x.Generations
	}
	return nil
}

func (x *Execution) GetPareto() *Population {
	if x != nil {
		return x.Pareto
	}
	return nil
}

type ExecutionSetIDs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	UserId string `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *ExecutionSetIDs) Reset() {
	*x = ExecutionSetIDs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_execution_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExecutionSetIDs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExecutionSetIDs) ProtoMessage() {}

func (x *ExecutionSetIDs) ProtoReflect() protoreflect.Message {
	mi := &file_api_execution_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExecutionSetIDs.ProtoReflect.Descriptor instead.
func (*ExecutionSetIDs) Descriptor() ([]byte, []int) {
	return file_api_execution_proto_rawDescGZIP(), []int{2}
}

func (x *ExecutionSetIDs) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *ExecutionSetIDs) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type ExecutionSet struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ids  *ExecutionSetIDs `protobuf:"bytes,1,opt,name=ids,proto3" json:"ids,omitempty"`
	Size uint32           `protobuf:"varint,2,opt,name=size,proto3" json:"size,omitempty"`
}

func (x *ExecutionSet) Reset() {
	*x = ExecutionSet{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_execution_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExecutionSet) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExecutionSet) ProtoMessage() {}

func (x *ExecutionSet) ProtoReflect() protoreflect.Message {
	mi := &file_api_execution_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExecutionSet.ProtoReflect.Descriptor instead.
func (*ExecutionSet) Descriptor() ([]byte, []int) {
	return file_api_execution_proto_rawDescGZIP(), []int{3}
}

func (x *ExecutionSet) GetIds() *ExecutionSetIDs {
	if x != nil {
		return x.Ids
	}
	return nil
}

func (x *ExecutionSet) GetSize() uint32 {
	if x != nil {
		return x.Size
	}
	return 0
}

var File_api_execution_proto protoreflect.FileDescriptor

var file_api_execution_proto_rawDesc = []byte{
	0x0a, 0x13, 0x61, 0x70, 0x69, 0x2f, 0x65, 0x78, 0x65, 0x63, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x61, 0x70, 0x69, 0x1a, 0x0e, 0x61, 0x70, 0x69, 0x2f,
	0x75, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x14, 0x61, 0x70, 0x69, 0x2f,
	0x70, 0x6f, 0x70, 0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x14, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x35, 0x0a, 0x0c, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x69, 0x6f, 0x6e,
	0x49, 0x44, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x15, 0x0a, 0x06, 0x73, 0x65, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x05, 0x73, 0x65, 0x74, 0x49, 0x64, 0x22, 0x8c, 0x01, 0x0a, 0x09, 0x45,
	0x78, 0x65, 0x63, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x23, 0x0a, 0x03, 0x69, 0x64, 0x73, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x45, 0x78, 0x65, 0x63,
	0x75, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x44, 0x73, 0x52, 0x03, 0x69, 0x64, 0x73, 0x12, 0x31, 0x0a,
	0x0b, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x0b, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x12, 0x27, 0x0a, 0x06, 0x70, 0x61, 0x72, 0x65, 0x74, 0x6f, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0f, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x50, 0x6f, 0x70, 0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x06, 0x70, 0x61, 0x72, 0x65, 0x74, 0x6f, 0x22, 0x3a, 0x0a, 0x0f, 0x45, 0x78, 0x65,
	0x63, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x74, 0x49, 0x44, 0x73, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x17, 0x0a, 0x07,
	0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75,
	0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x4a, 0x0a, 0x0c, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x69,
	0x6f, 0x6e, 0x53, 0x65, 0x74, 0x12, 0x26, 0x0a, 0x03, 0x69, 0x64, 0x73, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x14, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x69,
	0x6f, 0x6e, 0x53, 0x65, 0x74, 0x49, 0x44, 0x73, 0x52, 0x03, 0x69, 0x64, 0x73, 0x12, 0x12, 0x0a,
	0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x73, 0x69, 0x7a,
	0x65, 0x32, 0xdf, 0x01, 0x0a, 0x11, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x12, 0x32, 0x0a, 0x06, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x12, 0x0e, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x69, 0x6f,
	0x6e, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x2b, 0x0a, 0x04, 0x52,
	0x65, 0x61, 0x64, 0x12, 0x11, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74,
	0x69, 0x6f, 0x6e, 0x49, 0x44, 0x73, 0x1a, 0x0e, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x45, 0x78, 0x65,
	0x63, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x00, 0x12, 0x32, 0x0a, 0x06, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x12, 0x0e, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x69,
	0x6f, 0x6e, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x35, 0x0a, 0x06,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x11, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x45, 0x78, 0x65,
	0x63, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x44, 0x73, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x22, 0x00, 0x32, 0xa2, 0x02, 0x0a, 0x14, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x69, 0x6f,
	0x6e, 0x53, 0x65, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x12, 0x35, 0x0a, 0x06,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x11, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x45, 0x78, 0x65,
	0x63, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x22, 0x00, 0x12, 0x36, 0x0a, 0x04, 0x46, 0x69, 0x6e, 0x64, 0x12, 0x14, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x74, 0x49, 0x44,
	0x73, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x38, 0x0a, 0x06, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x14, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x45, 0x78, 0x65, 0x63,
	0x75, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x74, 0x49, 0x44, 0x73, 0x1a, 0x16, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x30, 0x0a, 0x04, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x14, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x74,
	0x49, 0x44, 0x73, 0x1a, 0x0e, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74,
	0x69, 0x6f, 0x6e, 0x22, 0x00, 0x30, 0x01, 0x12, 0x2f, 0x0a, 0x0b, 0x46, 0x65, 0x74, 0x63, 0x68,
	0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x12, 0x0c, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x55, 0x73, 0x65,
	0x72, 0x49, 0x44, 0x73, 0x1a, 0x0e, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x45, 0x78, 0x65, 0x63, 0x75,
	0x74, 0x69, 0x6f, 0x6e, 0x22, 0x00, 0x30, 0x01, 0x42, 0x09, 0x5a, 0x07, 0x70, 0x6b, 0x67, 0x2f,
	0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_execution_proto_rawDescOnce sync.Once
	file_api_execution_proto_rawDescData = file_api_execution_proto_rawDesc
)

func file_api_execution_proto_rawDescGZIP() []byte {
	file_api_execution_proto_rawDescOnce.Do(func() {
		file_api_execution_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_execution_proto_rawDescData)
	})
	return file_api_execution_proto_rawDescData
}

var file_api_execution_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_api_execution_proto_goTypes = []interface{}{
	(*ExecutionIDs)(nil),    // 0: api.ExecutionIDs
	(*Execution)(nil),       // 1: api.Execution
	(*ExecutionSetIDs)(nil), // 2: api.ExecutionSetIDs
	(*ExecutionSet)(nil),    // 3: api.ExecutionSet
	(*Generation)(nil),      // 4: api.Generation
	(*Population)(nil),      // 5: api.Population
	(*UserIDs)(nil),         // 6: api.UserIDs
	(*emptypb.Empty)(nil),   // 7: google.protobuf.Empty
}
var file_api_execution_proto_depIdxs = []int32{
	0,  // 0: api.Execution.ids:type_name -> api.ExecutionIDs
	4,  // 1: api.Execution.generations:type_name -> api.Generation
	5,  // 2: api.Execution.pareto:type_name -> api.Population
	2,  // 3: api.ExecutionSet.ids:type_name -> api.ExecutionSetIDs
	1,  // 4: api.ExecutionServices.Create:input_type -> api.Execution
	0,  // 5: api.ExecutionServices.Read:input_type -> api.ExecutionIDs
	1,  // 6: api.ExecutionServices.Update:input_type -> api.Execution
	0,  // 7: api.ExecutionServices.Delete:input_type -> api.ExecutionIDs
	3,  // 8: api.ExecutionSetServices.Create:input_type -> api.ExecutionSet
	2,  // 9: api.ExecutionSetServices.Find:input_type -> api.ExecutionSetIDs
	2,  // 10: api.ExecutionSetServices.Delete:input_type -> api.ExecutionSetIDs
	2,  // 11: api.ExecutionSetServices.List:input_type -> api.ExecutionSetIDs
	6,  // 12: api.ExecutionSetServices.FetchByUser:input_type -> api.UserIDs
	7,  // 13: api.ExecutionServices.Create:output_type -> google.protobuf.Empty
	1,  // 14: api.ExecutionServices.Read:output_type -> api.Execution
	7,  // 15: api.ExecutionServices.Update:output_type -> google.protobuf.Empty
	7,  // 16: api.ExecutionServices.Delete:output_type -> google.protobuf.Empty
	7,  // 17: api.ExecutionSetServices.Create:output_type -> google.protobuf.Empty
	7,  // 18: api.ExecutionSetServices.Find:output_type -> google.protobuf.Empty
	7,  // 19: api.ExecutionSetServices.Delete:output_type -> google.protobuf.Empty
	1,  // 20: api.ExecutionSetServices.List:output_type -> api.Execution
	1,  // 21: api.ExecutionSetServices.FetchByUser:output_type -> api.Execution
	13, // [13:22] is the sub-list for method output_type
	4,  // [4:13] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_api_execution_proto_init() }
func file_api_execution_proto_init() {
	if File_api_execution_proto != nil {
		return
	}
	file_api_user_proto_init()
	file_api_population_proto_init()
	file_api_generation_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_api_execution_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExecutionIDs); i {
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
		file_api_execution_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Execution); i {
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
		file_api_execution_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExecutionSetIDs); i {
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
		file_api_execution_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExecutionSet); i {
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
			RawDescriptor: file_api_execution_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_api_execution_proto_goTypes,
		DependencyIndexes: file_api_execution_proto_depIdxs,
		MessageInfos:      file_api_execution_proto_msgTypes,
	}.Build()
	File_api_execution_proto = out.File
	file_api_execution_proto_rawDesc = nil
	file_api_execution_proto_goTypes = nil
	file_api_execution_proto_depIdxs = nil
}
