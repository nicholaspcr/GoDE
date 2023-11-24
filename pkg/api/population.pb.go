// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.25.1
// source: api/population.proto

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

type Vector struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Elements         []float64 `protobuf:"fixed64,1,rep,packed,name=elements,proto3" json:"elements,omitempty"`
	Objectives       []float64 `protobuf:"fixed64,2,rep,packed,name=objectives,proto3" json:"objectives,omitempty"`
	CrowdingDistance float64   `protobuf:"fixed64,3,opt,name=crowding_distance,json=crowdingDistance,proto3" json:"crowding_distance,omitempty"`
}

func (x *Vector) Reset() {
	*x = Vector{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_population_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Vector) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Vector) ProtoMessage() {}

func (x *Vector) ProtoReflect() protoreflect.Message {
	mi := &file_api_population_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Vector.ProtoReflect.Descriptor instead.
func (*Vector) Descriptor() ([]byte, []int) {
	return file_api_population_proto_rawDescGZIP(), []int{0}
}

func (x *Vector) GetElements() []float64 {
	if x != nil {
		return x.Elements
	}
	return nil
}

func (x *Vector) GetObjectives() []float64 {
	if x != nil {
		return x.Objectives
	}
	return nil
}

func (x *Vector) GetCrowdingDistance() float64 {
	if x != nil {
		return x.CrowdingDistance
	}
	return 0
}

type PopulationIDs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	UserId string `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *PopulationIDs) Reset() {
	*x = PopulationIDs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_population_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PopulationIDs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PopulationIDs) ProtoMessage() {}

func (x *PopulationIDs) ProtoReflect() protoreflect.Message {
	mi := &file_api_population_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PopulationIDs.ProtoReflect.Descriptor instead.
func (*PopulationIDs) Descriptor() ([]byte, []int) {
	return file_api_population_proto_rawDescGZIP(), []int{1}
}

func (x *PopulationIDs) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *PopulationIDs) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type Population struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ids     *PopulationIDs `protobuf:"bytes,1,opt,name=ids,proto3" json:"ids,omitempty"`
	Vectors []*Vector      `protobuf:"bytes,2,rep,name=vectors,proto3" json:"vectors,omitempty"`
}

func (x *Population) Reset() {
	*x = Population{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_population_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Population) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Population) ProtoMessage() {}

func (x *Population) ProtoReflect() protoreflect.Message {
	mi := &file_api_population_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Population.ProtoReflect.Descriptor instead.
func (*Population) Descriptor() ([]byte, []int) {
	return file_api_population_proto_rawDescGZIP(), []int{2}
}

func (x *Population) GetIds() *PopulationIDs {
	if x != nil {
		return x.Ids
	}
	return nil
}

func (x *Population) GetVectors() []*Vector {
	if x != nil {
		return x.Vectors
	}
	return nil
}

type PopulationParameters struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DimensionsSize int64     `protobuf:"varint,1,opt,name=dimensions_size,json=dimensionsSize,proto3" json:"dimensions_size,omitempty"`
	ObjetivesSize  int64     `protobuf:"varint,2,opt,name=objetives_size,json=objetivesSize,proto3" json:"objetives_size,omitempty"`
	Floors         []float64 `protobuf:"fixed64,3,rep,packed,name=floors,proto3" json:"floors,omitempty"`
	Ceils          []float64 `protobuf:"fixed64,4,rep,packed,name=ceils,proto3" json:"ceils,omitempty"`
}

func (x *PopulationParameters) Reset() {
	*x = PopulationParameters{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_population_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PopulationParameters) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PopulationParameters) ProtoMessage() {}

func (x *PopulationParameters) ProtoReflect() protoreflect.Message {
	mi := &file_api_population_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PopulationParameters.ProtoReflect.Descriptor instead.
func (*PopulationParameters) Descriptor() ([]byte, []int) {
	return file_api_population_proto_rawDescGZIP(), []int{3}
}

func (x *PopulationParameters) GetDimensionsSize() int64 {
	if x != nil {
		return x.DimensionsSize
	}
	return 0
}

func (x *PopulationParameters) GetObjetivesSize() int64 {
	if x != nil {
		return x.ObjetivesSize
	}
	return 0
}

func (x *PopulationParameters) GetFloors() []float64 {
	if x != nil {
		return x.Floors
	}
	return nil
}

func (x *PopulationParameters) GetCeils() []float64 {
	if x != nil {
		return x.Ceils
	}
	return nil
}

type ParetoIDs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	UserId string `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *ParetoIDs) Reset() {
	*x = ParetoIDs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_population_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ParetoIDs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ParetoIDs) ProtoMessage() {}

func (x *ParetoIDs) ProtoReflect() protoreflect.Message {
	mi := &file_api_population_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ParetoIDs.ProtoReflect.Descriptor instead.
func (*ParetoIDs) Descriptor() ([]byte, []int) {
	return file_api_population_proto_rawDescGZIP(), []int{4}
}

func (x *ParetoIDs) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *ParetoIDs) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type Pareto struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ids        *ParetoIDs  `protobuf:"bytes,1,opt,name=ids,proto3" json:"ids,omitempty"`
	Population *Population `protobuf:"bytes,2,opt,name=population,proto3" json:"population,omitempty"`
	MaxObjs    []float64   `protobuf:"fixed64,3,rep,packed,name=max_objs,json=maxObjs,proto3" json:"max_objs,omitempty"`
}

func (x *Pareto) Reset() {
	*x = Pareto{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_population_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Pareto) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Pareto) ProtoMessage() {}

func (x *Pareto) ProtoReflect() protoreflect.Message {
	mi := &file_api_population_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Pareto.ProtoReflect.Descriptor instead.
func (*Pareto) Descriptor() ([]byte, []int) {
	return file_api_population_proto_rawDescGZIP(), []int{5}
}

func (x *Pareto) GetIds() *ParetoIDs {
	if x != nil {
		return x.Ids
	}
	return nil
}

func (x *Pareto) GetPopulation() *Population {
	if x != nil {
		return x.Population
	}
	return nil
}

func (x *Pareto) GetMaxObjs() []float64 {
	if x != nil {
		return x.MaxObjs
	}
	return nil
}

var File_api_population_proto protoreflect.FileDescriptor

var file_api_population_proto_rawDesc = []byte{
	0x0a, 0x14, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x6f, 0x70, 0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x61, 0x70, 0x69, 0x1a, 0x0e, 0x61, 0x70, 0x69,
	0x2f, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70,
	0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x71, 0x0a, 0x06, 0x56, 0x65, 0x63, 0x74,
	0x6f, 0x72, 0x12, 0x1a, 0x0a, 0x08, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x01, 0x52, 0x08, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x1e,
	0x0a, 0x0a, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x69, 0x76, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x01, 0x52, 0x0a, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x69, 0x76, 0x65, 0x73, 0x12, 0x2b,
	0x0a, 0x11, 0x63, 0x72, 0x6f, 0x77, 0x64, 0x69, 0x6e, 0x67, 0x5f, 0x64, 0x69, 0x73, 0x74, 0x61,
	0x6e, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x10, 0x63, 0x72, 0x6f, 0x77, 0x64,
	0x69, 0x6e, 0x67, 0x44, 0x69, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x22, 0x38, 0x0a, 0x0d, 0x50,
	0x6f, 0x70, 0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x44, 0x73, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x17, 0x0a, 0x07,
	0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75,
	0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x59, 0x0a, 0x0a, 0x50, 0x6f, 0x70, 0x75, 0x6c, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x24, 0x0a, 0x03, 0x69, 0x64, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x12, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x50, 0x6f, 0x70, 0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x49, 0x44, 0x73, 0x52, 0x03, 0x69, 0x64, 0x73, 0x12, 0x25, 0x0a, 0x07, 0x76, 0x65, 0x63,
	0x74, 0x6f, 0x72, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x56, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x52, 0x07, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x73,
	0x22, 0x94, 0x01, 0x0a, 0x14, 0x50, 0x6f, 0x70, 0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x50,
	0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x12, 0x27, 0x0a, 0x0f, 0x64, 0x69, 0x6d,
	0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x0e, 0x64, 0x69, 0x6d, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x53, 0x69,
	0x7a, 0x65, 0x12, 0x25, 0x0a, 0x0e, 0x6f, 0x62, 0x6a, 0x65, 0x74, 0x69, 0x76, 0x65, 0x73, 0x5f,
	0x73, 0x69, 0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0d, 0x6f, 0x62, 0x6a, 0x65,
	0x74, 0x69, 0x76, 0x65, 0x73, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x6c, 0x6f,
	0x6f, 0x72, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x01, 0x52, 0x06, 0x66, 0x6c, 0x6f, 0x6f, 0x72,
	0x73, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x65, 0x69, 0x6c, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x01,
	0x52, 0x05, 0x63, 0x65, 0x69, 0x6c, 0x73, 0x22, 0x34, 0x0a, 0x09, 0x50, 0x61, 0x72, 0x65, 0x74,
	0x6f, 0x49, 0x44, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x76, 0x0a,
	0x06, 0x50, 0x61, 0x72, 0x65, 0x74, 0x6f, 0x12, 0x20, 0x0a, 0x03, 0x69, 0x64, 0x73, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x50, 0x61, 0x72, 0x65, 0x74,
	0x6f, 0x49, 0x44, 0x73, 0x52, 0x03, 0x69, 0x64, 0x73, 0x12, 0x2f, 0x0a, 0x0a, 0x70, 0x6f, 0x70,
	0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x50, 0x6f, 0x70, 0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0a,
	0x70, 0x6f, 0x70, 0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x19, 0x0a, 0x08, 0x6d, 0x61,
	0x78, 0x5f, 0x6f, 0x62, 0x6a, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x01, 0x52, 0x07, 0x6d, 0x61,
	0x78, 0x4f, 0x62, 0x6a, 0x73, 0x32, 0x96, 0x02, 0x0a, 0x12, 0x50, 0x6f, 0x70, 0x75, 0x6c, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x12, 0x33, 0x0a, 0x06,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x0f, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x50, 0x6f, 0x70,
	0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22,
	0x00, 0x12, 0x2d, 0x0a, 0x04, 0x52, 0x65, 0x61, 0x64, 0x12, 0x12, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x50, 0x6f, 0x70, 0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x44, 0x73, 0x1a, 0x0f, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x50, 0x6f, 0x70, 0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x00,
	0x12, 0x33, 0x0a, 0x06, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x0f, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x50, 0x6f, 0x70, 0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x16, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x36, 0x0a, 0x06, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12,
	0x12, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x50, 0x6f, 0x70, 0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x49, 0x44, 0x73, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x2f, 0x0a,
	0x0a, 0x4c, 0x69, 0x73, 0x74, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x12, 0x0c, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x73, 0x1a, 0x0f, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x50, 0x6f, 0x70, 0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x00, 0x30, 0x01, 0x32, 0xfa,
	0x01, 0x0a, 0x0e, 0x50, 0x61, 0x72, 0x65, 0x74, 0x6f, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x73, 0x12, 0x2f, 0x0a, 0x06, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x0b, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x50, 0x61, 0x72, 0x65, 0x74, 0x6f, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x22, 0x00, 0x12, 0x25, 0x0a, 0x04, 0x52, 0x65, 0x61, 0x64, 0x12, 0x0e, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x50, 0x61, 0x72, 0x65, 0x74, 0x6f, 0x49, 0x44, 0x73, 0x1a, 0x0b, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x50, 0x61, 0x72, 0x65, 0x74, 0x6f, 0x22, 0x00, 0x12, 0x2f, 0x0a, 0x06, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x12, 0x0b, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x50, 0x61, 0x72, 0x65, 0x74, 0x6f,
	0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x32, 0x0a, 0x06, 0x44, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x12, 0x0e, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x50, 0x61, 0x72, 0x65, 0x74,
	0x6f, 0x49, 0x44, 0x73, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x2b,
	0x0a, 0x0a, 0x4c, 0x69, 0x73, 0x74, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x12, 0x0c, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x73, 0x1a, 0x0b, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x50, 0x61, 0x72, 0x65, 0x74, 0x6f, 0x22, 0x00, 0x30, 0x01, 0x42, 0x09, 0x5a, 0x07, 0x70,
	0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_population_proto_rawDescOnce sync.Once
	file_api_population_proto_rawDescData = file_api_population_proto_rawDesc
)

func file_api_population_proto_rawDescGZIP() []byte {
	file_api_population_proto_rawDescOnce.Do(func() {
		file_api_population_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_population_proto_rawDescData)
	})
	return file_api_population_proto_rawDescData
}

var file_api_population_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_api_population_proto_goTypes = []interface{}{
	(*Vector)(nil),               // 0: api.Vector
	(*PopulationIDs)(nil),        // 1: api.PopulationIDs
	(*Population)(nil),           // 2: api.Population
	(*PopulationParameters)(nil), // 3: api.PopulationParameters
	(*ParetoIDs)(nil),            // 4: api.ParetoIDs
	(*Pareto)(nil),               // 5: api.Pareto
	(*UserIDs)(nil),              // 6: api.UserIDs
	(*emptypb.Empty)(nil),        // 7: google.protobuf.Empty
}
var file_api_population_proto_depIdxs = []int32{
	1,  // 0: api.Population.ids:type_name -> api.PopulationIDs
	0,  // 1: api.Population.vectors:type_name -> api.Vector
	4,  // 2: api.Pareto.ids:type_name -> api.ParetoIDs
	2,  // 3: api.Pareto.population:type_name -> api.Population
	2,  // 4: api.PopulationServices.Create:input_type -> api.Population
	1,  // 5: api.PopulationServices.Read:input_type -> api.PopulationIDs
	2,  // 6: api.PopulationServices.Update:input_type -> api.Population
	1,  // 7: api.PopulationServices.Delete:input_type -> api.PopulationIDs
	6,  // 8: api.PopulationServices.ListByUser:input_type -> api.UserIDs
	5,  // 9: api.ParetoServices.Create:input_type -> api.Pareto
	4,  // 10: api.ParetoServices.Read:input_type -> api.ParetoIDs
	5,  // 11: api.ParetoServices.Update:input_type -> api.Pareto
	4,  // 12: api.ParetoServices.Delete:input_type -> api.ParetoIDs
	6,  // 13: api.ParetoServices.ListByUser:input_type -> api.UserIDs
	7,  // 14: api.PopulationServices.Create:output_type -> google.protobuf.Empty
	2,  // 15: api.PopulationServices.Read:output_type -> api.Population
	7,  // 16: api.PopulationServices.Update:output_type -> google.protobuf.Empty
	7,  // 17: api.PopulationServices.Delete:output_type -> google.protobuf.Empty
	2,  // 18: api.PopulationServices.ListByUser:output_type -> api.Population
	7,  // 19: api.ParetoServices.Create:output_type -> google.protobuf.Empty
	5,  // 20: api.ParetoServices.Read:output_type -> api.Pareto
	7,  // 21: api.ParetoServices.Update:output_type -> google.protobuf.Empty
	7,  // 22: api.ParetoServices.Delete:output_type -> google.protobuf.Empty
	5,  // 23: api.ParetoServices.ListByUser:output_type -> api.Pareto
	14, // [14:24] is the sub-list for method output_type
	4,  // [4:14] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_api_population_proto_init() }
func file_api_population_proto_init() {
	if File_api_population_proto != nil {
		return
	}
	file_api_user_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_api_population_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Vector); i {
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
		file_api_population_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PopulationIDs); i {
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
		file_api_population_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Population); i {
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
		file_api_population_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PopulationParameters); i {
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
		file_api_population_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ParetoIDs); i {
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
		file_api_population_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Pareto); i {
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
			RawDescriptor: file_api_population_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_api_population_proto_goTypes,
		DependencyIndexes: file_api_population_proto_depIdxs,
		MessageInfos:      file_api_population_proto_msgTypes,
	}.Build()
	File_api_population_proto = out.File
	file_api_population_proto_rawDesc = nil
	file_api_population_proto_goTypes = nil
	file_api_population_proto_depIdxs = nil
}
