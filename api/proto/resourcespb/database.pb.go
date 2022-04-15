// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: api/proto/resourcespb/database.proto

package resourcespb

import (
	commonpb "github.com/multycloud/multy/api/proto/commonpb"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type DatabaseEngine int32

const (
	DatabaseEngine_UNKNOWN_ENGINE DatabaseEngine = 0
	DatabaseEngine_MYSQL          DatabaseEngine = 1
)

// Enum value maps for DatabaseEngine.
var (
	DatabaseEngine_name = map[int32]string{
		0: "UNKNOWN_ENGINE",
		1: "MYSQL",
	}
	DatabaseEngine_value = map[string]int32{
		"UNKNOWN_ENGINE": 0,
		"MYSQL":          1,
	}
)

func (x DatabaseEngine) Enum() *DatabaseEngine {
	p := new(DatabaseEngine)
	*p = x
	return p
}

func (x DatabaseEngine) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (DatabaseEngine) Descriptor() protoreflect.EnumDescriptor {
	return file_api_proto_resourcespb_database_proto_enumTypes[0].Descriptor()
}

func (DatabaseEngine) Type() protoreflect.EnumType {
	return &file_api_proto_resourcespb_database_proto_enumTypes[0]
}

func (x DatabaseEngine) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use DatabaseEngine.Descriptor instead.
func (DatabaseEngine) EnumDescriptor() ([]byte, []int) {
	return file_api_proto_resourcespb_database_proto_rawDescGZIP(), []int{0}
}

type CreateDatabaseRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Resource *DatabaseArgs `protobuf:"bytes,1,opt,name=resource,proto3" json:"resource,omitempty"`
}

func (x *CreateDatabaseRequest) Reset() {
	*x = CreateDatabaseRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_resourcespb_database_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateDatabaseRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateDatabaseRequest) ProtoMessage() {}

func (x *CreateDatabaseRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_resourcespb_database_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateDatabaseRequest.ProtoReflect.Descriptor instead.
func (*CreateDatabaseRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_resourcespb_database_proto_rawDescGZIP(), []int{0}
}

func (x *CreateDatabaseRequest) GetResource() *DatabaseArgs {
	if x != nil {
		return x.Resource
	}
	return nil
}

type ReadDatabaseRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ResourceId string `protobuf:"bytes,1,opt,name=resource_id,json=resourceId,proto3" json:"resource_id,omitempty"`
}

func (x *ReadDatabaseRequest) Reset() {
	*x = ReadDatabaseRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_resourcespb_database_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadDatabaseRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadDatabaseRequest) ProtoMessage() {}

func (x *ReadDatabaseRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_resourcespb_database_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadDatabaseRequest.ProtoReflect.Descriptor instead.
func (*ReadDatabaseRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_resourcespb_database_proto_rawDescGZIP(), []int{1}
}

func (x *ReadDatabaseRequest) GetResourceId() string {
	if x != nil {
		return x.ResourceId
	}
	return ""
}

type UpdateDatabaseRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ResourceId string        `protobuf:"bytes,1,opt,name=resource_id,json=resourceId,proto3" json:"resource_id,omitempty"`
	Resource   *DatabaseArgs `protobuf:"bytes,2,opt,name=resource,proto3" json:"resource,omitempty"`
}

func (x *UpdateDatabaseRequest) Reset() {
	*x = UpdateDatabaseRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_resourcespb_database_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateDatabaseRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateDatabaseRequest) ProtoMessage() {}

func (x *UpdateDatabaseRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_resourcespb_database_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateDatabaseRequest.ProtoReflect.Descriptor instead.
func (*UpdateDatabaseRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_resourcespb_database_proto_rawDescGZIP(), []int{2}
}

func (x *UpdateDatabaseRequest) GetResourceId() string {
	if x != nil {
		return x.ResourceId
	}
	return ""
}

func (x *UpdateDatabaseRequest) GetResource() *DatabaseArgs {
	if x != nil {
		return x.Resource
	}
	return nil
}

type DeleteDatabaseRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ResourceId string `protobuf:"bytes,1,opt,name=resource_id,json=resourceId,proto3" json:"resource_id,omitempty"`
}

func (x *DeleteDatabaseRequest) Reset() {
	*x = DeleteDatabaseRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_resourcespb_database_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteDatabaseRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteDatabaseRequest) ProtoMessage() {}

func (x *DeleteDatabaseRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_resourcespb_database_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteDatabaseRequest.ProtoReflect.Descriptor instead.
func (*DeleteDatabaseRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_resourcespb_database_proto_rawDescGZIP(), []int{3}
}

func (x *DeleteDatabaseRequest) GetResourceId() string {
	if x != nil {
		return x.ResourceId
	}
	return ""
}

type DatabaseArgs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CommonParameters *commonpb.ResourceCommonArgs `protobuf:"bytes,1,opt,name=common_parameters,json=commonParameters,proto3" json:"common_parameters,omitempty"`
	Name             string                       `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Engine           DatabaseEngine               `protobuf:"varint,3,opt,name=engine,proto3,enum=dev.multy.resources.DatabaseEngine" json:"engine,omitempty"`
	EngineVersion    string                       `protobuf:"bytes,4,opt,name=engine_version,json=engineVersion,proto3" json:"engine_version,omitempty"`
	StorageGb        int64                        `protobuf:"varint,5,opt,name=storage_gb,json=storageGb,proto3" json:"storage_gb,omitempty"`
	Size             commonpb.DatabaseSize_Enum   `protobuf:"varint,6,opt,name=size,proto3,enum=dev.multy.common.DatabaseSize_Enum" json:"size,omitempty"`
	Username         string                       `protobuf:"bytes,7,opt,name=username,proto3" json:"username,omitempty"`
	Password         string                       `protobuf:"bytes,8,opt,name=password,proto3" json:"password,omitempty"`
	SubnetIds        []string                     `protobuf:"bytes,9,rep,name=subnet_ids,json=subnetIds,proto3" json:"subnet_ids,omitempty"`
}

func (x *DatabaseArgs) Reset() {
	*x = DatabaseArgs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_resourcespb_database_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DatabaseArgs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DatabaseArgs) ProtoMessage() {}

func (x *DatabaseArgs) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_resourcespb_database_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DatabaseArgs.ProtoReflect.Descriptor instead.
func (*DatabaseArgs) Descriptor() ([]byte, []int) {
	return file_api_proto_resourcespb_database_proto_rawDescGZIP(), []int{4}
}

func (x *DatabaseArgs) GetCommonParameters() *commonpb.ResourceCommonArgs {
	if x != nil {
		return x.CommonParameters
	}
	return nil
}

func (x *DatabaseArgs) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *DatabaseArgs) GetEngine() DatabaseEngine {
	if x != nil {
		return x.Engine
	}
	return DatabaseEngine_UNKNOWN_ENGINE
}

func (x *DatabaseArgs) GetEngineVersion() string {
	if x != nil {
		return x.EngineVersion
	}
	return ""
}

func (x *DatabaseArgs) GetStorageGb() int64 {
	if x != nil {
		return x.StorageGb
	}
	return 0
}

func (x *DatabaseArgs) GetSize() commonpb.DatabaseSize_Enum {
	if x != nil {
		return x.Size
	}
	return commonpb.DatabaseSize_Enum(0)
}

func (x *DatabaseArgs) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *DatabaseArgs) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *DatabaseArgs) GetSubnetIds() []string {
	if x != nil {
		return x.SubnetIds
	}
	return nil
}

type DatabaseResource struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CommonParameters *commonpb.CommonResourceParameters `protobuf:"bytes,1,opt,name=common_parameters,json=commonParameters,proto3" json:"common_parameters,omitempty"`
	Name             string                             `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Engine           DatabaseEngine                     `protobuf:"varint,3,opt,name=engine,proto3,enum=dev.multy.resources.DatabaseEngine" json:"engine,omitempty"`
	EngineVersion    string                             `protobuf:"bytes,4,opt,name=engine_version,json=engineVersion,proto3" json:"engine_version,omitempty"`
	StorageGb        int64                              `protobuf:"varint,5,opt,name=storage_gb,json=storageGb,proto3" json:"storage_gb,omitempty"`
	Size             commonpb.DatabaseSize_Enum         `protobuf:"varint,6,opt,name=size,proto3,enum=dev.multy.common.DatabaseSize_Enum" json:"size,omitempty"`
	Username         string                             `protobuf:"bytes,7,opt,name=username,proto3" json:"username,omitempty"`
	Password         string                             `protobuf:"bytes,8,opt,name=password,proto3" json:"password,omitempty"`
	SubnetIds        []string                           `protobuf:"bytes,9,rep,name=subnet_ids,json=subnetIds,proto3" json:"subnet_ids,omitempty"`
	//outputs
	Host string `protobuf:"bytes,10,opt,name=host,proto3" json:"host,omitempty"`
}

func (x *DatabaseResource) Reset() {
	*x = DatabaseResource{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_resourcespb_database_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DatabaseResource) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DatabaseResource) ProtoMessage() {}

func (x *DatabaseResource) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_resourcespb_database_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DatabaseResource.ProtoReflect.Descriptor instead.
func (*DatabaseResource) Descriptor() ([]byte, []int) {
	return file_api_proto_resourcespb_database_proto_rawDescGZIP(), []int{5}
}

func (x *DatabaseResource) GetCommonParameters() *commonpb.CommonResourceParameters {
	if x != nil {
		return x.CommonParameters
	}
	return nil
}

func (x *DatabaseResource) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *DatabaseResource) GetEngine() DatabaseEngine {
	if x != nil {
		return x.Engine
	}
	return DatabaseEngine_UNKNOWN_ENGINE
}

func (x *DatabaseResource) GetEngineVersion() string {
	if x != nil {
		return x.EngineVersion
	}
	return ""
}

func (x *DatabaseResource) GetStorageGb() int64 {
	if x != nil {
		return x.StorageGb
	}
	return 0
}

func (x *DatabaseResource) GetSize() commonpb.DatabaseSize_Enum {
	if x != nil {
		return x.Size
	}
	return commonpb.DatabaseSize_Enum(0)
}

func (x *DatabaseResource) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *DatabaseResource) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *DatabaseResource) GetSubnetIds() []string {
	if x != nil {
		return x.SubnetIds
	}
	return nil
}

func (x *DatabaseResource) GetHost() string {
	if x != nil {
		return x.Host
	}
	return ""
}

var File_api_proto_resourcespb_database_proto protoreflect.FileDescriptor

var file_api_proto_resourcespb_database_proto_rawDesc = []byte{
	0x0a, 0x24, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x73, 0x70, 0x62, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x13, 0x64, 0x65, 0x76, 0x2e, 0x6d, 0x75, 0x6c, 0x74,
	0x79, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x1a, 0x1f, 0x61, 0x70, 0x69,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x70, 0x62, 0x2f,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x56, 0x0a, 0x15,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x44, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x3d, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x64, 0x65, 0x76, 0x2e, 0x6d, 0x75,
	0x6c, 0x74, 0x79, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x44, 0x61,
	0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x41, 0x72, 0x67, 0x73, 0x52, 0x08, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x22, 0x36, 0x0a, 0x13, 0x52, 0x65, 0x61, 0x64, 0x44, 0x61, 0x74, 0x61,
	0x62, 0x61, 0x73, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0a, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64, 0x22, 0x77, 0x0a, 0x15,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x44, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x49, 0x64, 0x12, 0x3d, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x64, 0x65, 0x76, 0x2e, 0x6d,
	0x75, 0x6c, 0x74, 0x79, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x44,
	0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x41, 0x72, 0x67, 0x73, 0x52, 0x08, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x22, 0x38, 0x0a, 0x15, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x44,
	0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f,
	0x0a, 0x0b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0a, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64, 0x22,
	0x88, 0x03, 0x0a, 0x0c, 0x44, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x41, 0x72, 0x67, 0x73,
	0x12, 0x51, 0x0a, 0x11, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5f, 0x70, 0x61, 0x72, 0x61, 0x6d,
	0x65, 0x74, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x64, 0x65,
	0x76, 0x2e, 0x6d, 0x75, 0x6c, 0x74, 0x79, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x52,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x41, 0x72, 0x67,
	0x73, 0x52, 0x10, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74,
	0x65, 0x72, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x3b, 0x0a, 0x06, 0x65, 0x6e, 0x67, 0x69, 0x6e,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x23, 0x2e, 0x64, 0x65, 0x76, 0x2e, 0x6d, 0x75,
	0x6c, 0x74, 0x79, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x44, 0x61,
	0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x45, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x52, 0x06, 0x65, 0x6e,
	0x67, 0x69, 0x6e, 0x65, 0x12, 0x25, 0x0a, 0x0e, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x5f, 0x76,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x65, 0x6e,
	0x67, 0x69, 0x6e, 0x65, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1d, 0x0a, 0x0a, 0x73,
	0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x5f, 0x67, 0x62, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x09, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x47, 0x62, 0x12, 0x37, 0x0a, 0x04, 0x73, 0x69,
	0x7a, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x23, 0x2e, 0x64, 0x65, 0x76, 0x2e, 0x6d,
	0x75, 0x6c, 0x74, 0x79, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x44, 0x61, 0x74, 0x61,
	0x62, 0x61, 0x73, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x2e, 0x45, 0x6e, 0x75, 0x6d, 0x52, 0x04, 0x73,
	0x69, 0x7a, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x08, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x73,
	0x75, 0x62, 0x6e, 0x65, 0x74, 0x5f, 0x69, 0x64, 0x73, 0x18, 0x09, 0x20, 0x03, 0x28, 0x09, 0x52,
	0x09, 0x73, 0x75, 0x62, 0x6e, 0x65, 0x74, 0x49, 0x64, 0x73, 0x22, 0xa6, 0x03, 0x0a, 0x10, 0x44,
	0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12,
	0x57, 0x0a, 0x11, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5f, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x65,
	0x74, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2a, 0x2e, 0x64, 0x65, 0x76,
	0x2e, 0x6d, 0x75, 0x6c, 0x74, 0x79, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x43, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x50, 0x61, 0x72, 0x61,
	0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x52, 0x10, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x50, 0x61,
	0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x3b, 0x0a, 0x06,
	0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x23, 0x2e, 0x64,
	0x65, 0x76, 0x2e, 0x6d, 0x75, 0x6c, 0x74, 0x79, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x73, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x45, 0x6e, 0x67, 0x69, 0x6e,
	0x65, 0x52, 0x06, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x12, 0x25, 0x0a, 0x0e, 0x65, 0x6e, 0x67,
	0x69, 0x6e, 0x65, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0d, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x5f, 0x67, 0x62, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x47, 0x62, 0x12,
	0x37, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x23, 0x2e,
	0x64, 0x65, 0x76, 0x2e, 0x6d, 0x75, 0x6c, 0x74, 0x79, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x2e, 0x44, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x2e, 0x45, 0x6e,
	0x75, 0x6d, 0x52, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64,
	0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64,
	0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x75, 0x62, 0x6e, 0x65, 0x74, 0x5f, 0x69, 0x64, 0x73, 0x18, 0x09,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x09, 0x73, 0x75, 0x62, 0x6e, 0x65, 0x74, 0x49, 0x64, 0x73, 0x12,
	0x12, 0x0a, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68,
	0x6f, 0x73, 0x74, 0x2a, 0x2f, 0x0a, 0x0e, 0x44, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x45,
	0x6e, 0x67, 0x69, 0x6e, 0x65, 0x12, 0x12, 0x0a, 0x0e, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e,
	0x5f, 0x45, 0x4e, 0x47, 0x49, 0x4e, 0x45, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x4d, 0x59, 0x53,
	0x51, 0x4c, 0x10, 0x01, 0x42, 0x5e, 0x0a, 0x17, 0x64, 0x65, 0x76, 0x2e, 0x6d, 0x75, 0x6c, 0x74,
	0x79, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x42,
	0x0e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50,
	0x01, 0x5a, 0x31, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x75,
	0x6c, 0x74, 0x79, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2f, 0x6d, 0x75, 0x6c, 0x74, 0x79, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x73, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_proto_resourcespb_database_proto_rawDescOnce sync.Once
	file_api_proto_resourcespb_database_proto_rawDescData = file_api_proto_resourcespb_database_proto_rawDesc
)

func file_api_proto_resourcespb_database_proto_rawDescGZIP() []byte {
	file_api_proto_resourcespb_database_proto_rawDescOnce.Do(func() {
		file_api_proto_resourcespb_database_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_proto_resourcespb_database_proto_rawDescData)
	})
	return file_api_proto_resourcespb_database_proto_rawDescData
}

var file_api_proto_resourcespb_database_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_api_proto_resourcespb_database_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_api_proto_resourcespb_database_proto_goTypes = []interface{}{
	(DatabaseEngine)(0),                       // 0: dev.multy.resources.DatabaseEngine
	(*CreateDatabaseRequest)(nil),             // 1: dev.multy.resources.CreateDatabaseRequest
	(*ReadDatabaseRequest)(nil),               // 2: dev.multy.resources.ReadDatabaseRequest
	(*UpdateDatabaseRequest)(nil),             // 3: dev.multy.resources.UpdateDatabaseRequest
	(*DeleteDatabaseRequest)(nil),             // 4: dev.multy.resources.DeleteDatabaseRequest
	(*DatabaseArgs)(nil),                      // 5: dev.multy.resources.DatabaseArgs
	(*DatabaseResource)(nil),                  // 6: dev.multy.resources.DatabaseResource
	(*commonpb.ResourceCommonArgs)(nil),       // 7: dev.multy.common.ResourceCommonArgs
	(commonpb.DatabaseSize_Enum)(0),           // 8: dev.multy.common.DatabaseSize.Enum
	(*commonpb.CommonResourceParameters)(nil), // 9: dev.multy.common.CommonResourceParameters
}
var file_api_proto_resourcespb_database_proto_depIdxs = []int32{
	5, // 0: dev.multy.resources.CreateDatabaseRequest.resource:type_name -> dev.multy.resources.DatabaseArgs
	5, // 1: dev.multy.resources.UpdateDatabaseRequest.resource:type_name -> dev.multy.resources.DatabaseArgs
	7, // 2: dev.multy.resources.DatabaseArgs.common_parameters:type_name -> dev.multy.common.ResourceCommonArgs
	0, // 3: dev.multy.resources.DatabaseArgs.engine:type_name -> dev.multy.resources.DatabaseEngine
	8, // 4: dev.multy.resources.DatabaseArgs.size:type_name -> dev.multy.common.DatabaseSize.Enum
	9, // 5: dev.multy.resources.DatabaseResource.common_parameters:type_name -> dev.multy.common.CommonResourceParameters
	0, // 6: dev.multy.resources.DatabaseResource.engine:type_name -> dev.multy.resources.DatabaseEngine
	8, // 7: dev.multy.resources.DatabaseResource.size:type_name -> dev.multy.common.DatabaseSize.Enum
	8, // [8:8] is the sub-list for method output_type
	8, // [8:8] is the sub-list for method input_type
	8, // [8:8] is the sub-list for extension type_name
	8, // [8:8] is the sub-list for extension extendee
	0, // [0:8] is the sub-list for field type_name
}

func init() { file_api_proto_resourcespb_database_proto_init() }
func file_api_proto_resourcespb_database_proto_init() {
	if File_api_proto_resourcespb_database_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_proto_resourcespb_database_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateDatabaseRequest); i {
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
		file_api_proto_resourcespb_database_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReadDatabaseRequest); i {
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
		file_api_proto_resourcespb_database_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateDatabaseRequest); i {
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
		file_api_proto_resourcespb_database_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteDatabaseRequest); i {
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
		file_api_proto_resourcespb_database_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DatabaseArgs); i {
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
		file_api_proto_resourcespb_database_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DatabaseResource); i {
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
			RawDescriptor: file_api_proto_resourcespb_database_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_proto_resourcespb_database_proto_goTypes,
		DependencyIndexes: file_api_proto_resourcespb_database_proto_depIdxs,
		EnumInfos:         file_api_proto_resourcespb_database_proto_enumTypes,
		MessageInfos:      file_api_proto_resourcespb_database_proto_msgTypes,
	}.Build()
	File_api_proto_resourcespb_database_proto = out.File
	file_api_proto_resourcespb_database_proto_rawDesc = nil
	file_api_proto_resourcespb_database_proto_goTypes = nil
	file_api_proto_resourcespb_database_proto_depIdxs = nil
}
