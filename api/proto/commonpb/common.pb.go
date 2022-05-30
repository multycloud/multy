// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: api/proto/commonpb/common.proto

package commonpb

import (
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

type Location int32

const (
	Location_UNKNOWN_LOCATION Location = 0
	// N. Virginia
	Location_US_EAST_1 Location = 1
	// Ohio / Virginia
	Location_US_EAST_2 Location = 2
	// N. California / Washington
	Location_US_WEST_1 Location = 3
	// Oregon / Phoenix
	Location_US_WEST_2 Location = 4
	// Ireland
	Location_EU_WEST_1 Location = 5
	// England
	Location_EU_WEST_2 Location = 6
	// France
	Location_EU_WEST_3 Location = 7
	// Sweden
	Location_EU_NORTH_1 Location = 8
)

// Enum value maps for Location.
var (
	Location_name = map[int32]string{
		0: "UNKNOWN_LOCATION",
		1: "US_EAST_1",
		2: "US_EAST_2",
		3: "US_WEST_1",
		4: "US_WEST_2",
		5: "EU_WEST_1",
		6: "EU_WEST_2",
		7: "EU_WEST_3",
		8: "EU_NORTH_1",
	}
	Location_value = map[string]int32{
		"UNKNOWN_LOCATION": 0,
		"US_EAST_1":        1,
		"US_EAST_2":        2,
		"US_WEST_1":        3,
		"US_WEST_2":        4,
		"EU_WEST_1":        5,
		"EU_WEST_2":        6,
		"EU_WEST_3":        7,
		"EU_NORTH_1":       8,
	}
)

func (x Location) Enum() *Location {
	p := new(Location)
	*p = x
	return p
}

func (x Location) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Location) Descriptor() protoreflect.EnumDescriptor {
	return file_api_proto_commonpb_common_proto_enumTypes[0].Descriptor()
}

func (Location) Type() protoreflect.EnumType {
	return &file_api_proto_commonpb_common_proto_enumTypes[0]
}

func (x Location) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Location.Descriptor instead.
func (Location) EnumDescriptor() ([]byte, []int) {
	return file_api_proto_commonpb_common_proto_rawDescGZIP(), []int{0}
}

type CloudProvider int32

const (
	CloudProvider_UNKNOWN_PROVIDER CloudProvider = 0
	CloudProvider_AWS              CloudProvider = 1
	CloudProvider_AZURE            CloudProvider = 2
)

// Enum value maps for CloudProvider.
var (
	CloudProvider_name = map[int32]string{
		0: "UNKNOWN_PROVIDER",
		1: "AWS",
		2: "AZURE",
	}
	CloudProvider_value = map[string]int32{
		"UNKNOWN_PROVIDER": 0,
		"AWS":              1,
		"AZURE":            2,
	}
)

func (x CloudProvider) Enum() *CloudProvider {
	p := new(CloudProvider)
	*p = x
	return p
}

func (x CloudProvider) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CloudProvider) Descriptor() protoreflect.EnumDescriptor {
	return file_api_proto_commonpb_common_proto_enumTypes[1].Descriptor()
}

func (CloudProvider) Type() protoreflect.EnumType {
	return &file_api_proto_commonpb_common_proto_enumTypes[1]
}

func (x CloudProvider) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CloudProvider.Descriptor instead.
func (CloudProvider) EnumDescriptor() ([]byte, []int) {
	return file_api_proto_commonpb_common_proto_rawDescGZIP(), []int{1}
}

type OperatingSystem_Enum int32

const (
	OperatingSystem_UNKNOWN_OS OperatingSystem_Enum = 0
	OperatingSystem_LINUX      OperatingSystem_Enum = 1
)

// Enum value maps for OperatingSystem_Enum.
var (
	OperatingSystem_Enum_name = map[int32]string{
		0: "UNKNOWN_OS",
		1: "LINUX",
	}
	OperatingSystem_Enum_value = map[string]int32{
		"UNKNOWN_OS": 0,
		"LINUX":      1,
	}
)

func (x OperatingSystem_Enum) Enum() *OperatingSystem_Enum {
	p := new(OperatingSystem_Enum)
	*p = x
	return p
}

func (x OperatingSystem_Enum) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (OperatingSystem_Enum) Descriptor() protoreflect.EnumDescriptor {
	return file_api_proto_commonpb_common_proto_enumTypes[2].Descriptor()
}

func (OperatingSystem_Enum) Type() protoreflect.EnumType {
	return &file_api_proto_commonpb_common_proto_enumTypes[2]
}

func (x OperatingSystem_Enum) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use OperatingSystem_Enum.Descriptor instead.
func (OperatingSystem_Enum) EnumDescriptor() ([]byte, []int) {
	return file_api_proto_commonpb_common_proto_rawDescGZIP(), []int{0, 0}
}

type DatabaseSize_Enum int32

const (
	DatabaseSize_UNKNOWN_VM_SIZE DatabaseSize_Enum = 0
	DatabaseSize_MICRO           DatabaseSize_Enum = 2
	DatabaseSize_SMALL           DatabaseSize_Enum = 4
	DatabaseSize_MEDIUM          DatabaseSize_Enum = 3
)

// Enum value maps for DatabaseSize_Enum.
var (
	DatabaseSize_Enum_name = map[int32]string{
		0: "UNKNOWN_VM_SIZE",
		2: "MICRO",
		4: "SMALL",
		3: "MEDIUM",
	}
	DatabaseSize_Enum_value = map[string]int32{
		"UNKNOWN_VM_SIZE": 0,
		"MICRO":           2,
		"SMALL":           4,
		"MEDIUM":          3,
	}
)

func (x DatabaseSize_Enum) Enum() *DatabaseSize_Enum {
	p := new(DatabaseSize_Enum)
	*p = x
	return p
}

func (x DatabaseSize_Enum) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (DatabaseSize_Enum) Descriptor() protoreflect.EnumDescriptor {
	return file_api_proto_commonpb_common_proto_enumTypes[3].Descriptor()
}

func (DatabaseSize_Enum) Type() protoreflect.EnumType {
	return &file_api_proto_commonpb_common_proto_enumTypes[3]
}

func (x DatabaseSize_Enum) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use DatabaseSize_Enum.Descriptor instead.
func (DatabaseSize_Enum) EnumDescriptor() ([]byte, []int) {
	return file_api_proto_commonpb_common_proto_rawDescGZIP(), []int{1, 0}
}

type VmSize_Enum int32

const (
	VmSize_UNKNOWN_VM_SIZE  VmSize_Enum = 0
	VmSize_MICRO            VmSize_Enum = 1
	VmSize_MEDIUM           VmSize_Enum = 2
	VmSize_LARGE            VmSize_Enum = 3
	VmSize_GENERAL_NANO     VmSize_Enum = 4
	VmSize_GENERAL_LARGE    VmSize_Enum = 5
	VmSize_GENERAL_XLARGE   VmSize_Enum = 6
	VmSize_GENERAL_2XLARGE  VmSize_Enum = 7
	VmSize_COMPUTE_LARGE    VmSize_Enum = 8
	VmSize_COMPUTE_XLARGE   VmSize_Enum = 9
	VmSize_COMPUTE_2XLARGE  VmSize_Enum = 10
	VmSize_COMPUTE_4XLARGE  VmSize_Enum = 11
	VmSize_COMPUTE_8XLARGE  VmSize_Enum = 12
	VmSize_COMPUTE_12XLARGE VmSize_Enum = 13
	VmSize_COMPUTE_16XLARGE VmSize_Enum = 14
)

// Enum value maps for VmSize_Enum.
var (
	VmSize_Enum_name = map[int32]string{
		0:  "UNKNOWN_VM_SIZE",
		1:  "MICRO",
		2:  "MEDIUM",
		3:  "LARGE",
		4:  "GENERAL_NANO",
		5:  "GENERAL_LARGE",
		6:  "GENERAL_XLARGE",
		7:  "GENERAL_2XLARGE",
		8:  "COMPUTE_LARGE",
		9:  "COMPUTE_XLARGE",
		10: "COMPUTE_2XLARGE",
		11: "COMPUTE_4XLARGE",
		12: "COMPUTE_8XLARGE",
		13: "COMPUTE_12XLARGE",
		14: "COMPUTE_16XLARGE",
	}
	VmSize_Enum_value = map[string]int32{
		"UNKNOWN_VM_SIZE":  0,
		"MICRO":            1,
		"MEDIUM":           2,
		"LARGE":            3,
		"GENERAL_NANO":     4,
		"GENERAL_LARGE":    5,
		"GENERAL_XLARGE":   6,
		"GENERAL_2XLARGE":  7,
		"COMPUTE_LARGE":    8,
		"COMPUTE_XLARGE":   9,
		"COMPUTE_2XLARGE":  10,
		"COMPUTE_4XLARGE":  11,
		"COMPUTE_8XLARGE":  12,
		"COMPUTE_12XLARGE": 13,
		"COMPUTE_16XLARGE": 14,
	}
)

func (x VmSize_Enum) Enum() *VmSize_Enum {
	p := new(VmSize_Enum)
	*p = x
	return p
}

func (x VmSize_Enum) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (VmSize_Enum) Descriptor() protoreflect.EnumDescriptor {
	return file_api_proto_commonpb_common_proto_enumTypes[4].Descriptor()
}

func (VmSize_Enum) Type() protoreflect.EnumType {
	return &file_api_proto_commonpb_common_proto_enumTypes[4]
}

func (x VmSize_Enum) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use VmSize_Enum.Descriptor instead.
func (VmSize_Enum) EnumDescriptor() ([]byte, []int) {
	return file_api_proto_commonpb_common_proto_rawDescGZIP(), []int{2, 0}
}

type OperatingSystem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *OperatingSystem) Reset() {
	*x = OperatingSystem{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_commonpb_common_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OperatingSystem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OperatingSystem) ProtoMessage() {}

func (x *OperatingSystem) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_commonpb_common_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OperatingSystem.ProtoReflect.Descriptor instead.
func (*OperatingSystem) Descriptor() ([]byte, []int) {
	return file_api_proto_commonpb_common_proto_rawDescGZIP(), []int{0}
}

type DatabaseSize struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DatabaseSize) Reset() {
	*x = DatabaseSize{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_commonpb_common_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DatabaseSize) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DatabaseSize) ProtoMessage() {}

func (x *DatabaseSize) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_commonpb_common_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DatabaseSize.ProtoReflect.Descriptor instead.
func (*DatabaseSize) Descriptor() ([]byte, []int) {
	return file_api_proto_commonpb_common_proto_rawDescGZIP(), []int{1}
}

type VmSize struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *VmSize) Reset() {
	*x = VmSize{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_commonpb_common_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VmSize) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VmSize) ProtoMessage() {}

func (x *VmSize) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_commonpb_common_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VmSize.ProtoReflect.Descriptor instead.
func (*VmSize) Descriptor() ([]byte, []int) {
	return file_api_proto_commonpb_common_proto_rawDescGZIP(), []int{2}
}

// Common messages for READ requests
type CommonResourceParameters struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ResourceId      string        `protobuf:"bytes,1,opt,name=resource_id,json=resourceId,proto3" json:"resource_id,omitempty"`
	ResourceGroupId string        `protobuf:"bytes,2,opt,name=resource_group_id,json=resourceGroupId,proto3" json:"resource_group_id,omitempty"`
	Location        Location      `protobuf:"varint,3,opt,name=location,proto3,enum=dev.multy.common.Location" json:"location,omitempty"`
	CloudProvider   CloudProvider `protobuf:"varint,4,opt,name=cloud_provider,json=cloudProvider,proto3,enum=dev.multy.common.CloudProvider" json:"cloud_provider,omitempty"`
	NeedsUpdate     bool          `protobuf:"varint,5,opt,name=needs_update,json=needsUpdate,proto3" json:"needs_update,omitempty"`
}

func (x *CommonResourceParameters) Reset() {
	*x = CommonResourceParameters{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_commonpb_common_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CommonResourceParameters) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommonResourceParameters) ProtoMessage() {}

func (x *CommonResourceParameters) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_commonpb_common_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommonResourceParameters.ProtoReflect.Descriptor instead.
func (*CommonResourceParameters) Descriptor() ([]byte, []int) {
	return file_api_proto_commonpb_common_proto_rawDescGZIP(), []int{3}
}

func (x *CommonResourceParameters) GetResourceId() string {
	if x != nil {
		return x.ResourceId
	}
	return ""
}

func (x *CommonResourceParameters) GetResourceGroupId() string {
	if x != nil {
		return x.ResourceGroupId
	}
	return ""
}

func (x *CommonResourceParameters) GetLocation() Location {
	if x != nil {
		return x.Location
	}
	return Location_UNKNOWN_LOCATION
}

func (x *CommonResourceParameters) GetCloudProvider() CloudProvider {
	if x != nil {
		return x.CloudProvider
	}
	return CloudProvider_UNKNOWN_PROVIDER
}

func (x *CommonResourceParameters) GetNeedsUpdate() bool {
	if x != nil {
		return x.NeedsUpdate
	}
	return false
}

type CommonChildResourceParameters struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ResourceId  string `protobuf:"bytes,1,opt,name=resource_id,json=resourceId,proto3" json:"resource_id,omitempty"`
	NeedsUpdate bool   `protobuf:"varint,2,opt,name=needs_update,json=needsUpdate,proto3" json:"needs_update,omitempty"`
}

func (x *CommonChildResourceParameters) Reset() {
	*x = CommonChildResourceParameters{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_commonpb_common_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CommonChildResourceParameters) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommonChildResourceParameters) ProtoMessage() {}

func (x *CommonChildResourceParameters) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_commonpb_common_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommonChildResourceParameters.ProtoReflect.Descriptor instead.
func (*CommonChildResourceParameters) Descriptor() ([]byte, []int) {
	return file_api_proto_commonpb_common_proto_rawDescGZIP(), []int{4}
}

func (x *CommonChildResourceParameters) GetResourceId() string {
	if x != nil {
		return x.ResourceId
	}
	return ""
}

func (x *CommonChildResourceParameters) GetNeedsUpdate() bool {
	if x != nil {
		return x.NeedsUpdate
	}
	return false
}

// Common messages for CREATE and UPDATE requests
type ResourceCommonArgs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ResourceGroupId string        `protobuf:"bytes,1,opt,name=resource_group_id,json=resourceGroupId,proto3" json:"resource_group_id,omitempty"`
	Location        Location      `protobuf:"varint,2,opt,name=location,proto3,enum=dev.multy.common.Location" json:"location,omitempty"`
	CloudProvider   CloudProvider `protobuf:"varint,3,opt,name=cloud_provider,json=cloudProvider,proto3,enum=dev.multy.common.CloudProvider" json:"cloud_provider,omitempty"`
}

func (x *ResourceCommonArgs) Reset() {
	*x = ResourceCommonArgs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_commonpb_common_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResourceCommonArgs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResourceCommonArgs) ProtoMessage() {}

func (x *ResourceCommonArgs) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_commonpb_common_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResourceCommonArgs.ProtoReflect.Descriptor instead.
func (*ResourceCommonArgs) Descriptor() ([]byte, []int) {
	return file_api_proto_commonpb_common_proto_rawDescGZIP(), []int{5}
}

func (x *ResourceCommonArgs) GetResourceGroupId() string {
	if x != nil {
		return x.ResourceGroupId
	}
	return ""
}

func (x *ResourceCommonArgs) GetLocation() Location {
	if x != nil {
		return x.Location
	}
	return Location_UNKNOWN_LOCATION
}

func (x *ResourceCommonArgs) GetCloudProvider() CloudProvider {
	if x != nil {
		return x.CloudProvider
	}
	return CloudProvider_UNKNOWN_PROVIDER
}

type ChildResourceCommonArgs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ChildResourceCommonArgs) Reset() {
	*x = ChildResourceCommonArgs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_commonpb_common_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChildResourceCommonArgs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChildResourceCommonArgs) ProtoMessage() {}

func (x *ChildResourceCommonArgs) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_commonpb_common_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChildResourceCommonArgs.ProtoReflect.Descriptor instead.
func (*ChildResourceCommonArgs) Descriptor() ([]byte, []int) {
	return file_api_proto_commonpb_common_proto_rawDescGZIP(), []int{6}
}

// Common messages for DELETE requests
type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_commonpb_common_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_commonpb_common_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_api_proto_commonpb_common_proto_rawDescGZIP(), []int{7}
}

type ListResourcesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Resources []*ListResourcesResponse_ResourceMetadata `protobuf:"bytes,1,rep,name=resources,proto3" json:"resources,omitempty"`
}

func (x *ListResourcesResponse) Reset() {
	*x = ListResourcesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_commonpb_common_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListResourcesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListResourcesResponse) ProtoMessage() {}

func (x *ListResourcesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_commonpb_common_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListResourcesResponse.ProtoReflect.Descriptor instead.
func (*ListResourcesResponse) Descriptor() ([]byte, []int) {
	return file_api_proto_commonpb_common_proto_rawDescGZIP(), []int{8}
}

func (x *ListResourcesResponse) GetResources() []*ListResourcesResponse_ResourceMetadata {
	if x != nil {
		return x.Resources
	}
	return nil
}

type ListResourcesResponse_ResourceMetadata struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ResourceId   string `protobuf:"bytes,1,opt,name=resource_id,json=resourceId,proto3" json:"resource_id,omitempty"`
	ResourceType string `protobuf:"bytes,2,opt,name=resource_type,json=resourceType,proto3" json:"resource_type,omitempty"`
}

func (x *ListResourcesResponse_ResourceMetadata) Reset() {
	*x = ListResourcesResponse_ResourceMetadata{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_commonpb_common_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListResourcesResponse_ResourceMetadata) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListResourcesResponse_ResourceMetadata) ProtoMessage() {}

func (x *ListResourcesResponse_ResourceMetadata) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_commonpb_common_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListResourcesResponse_ResourceMetadata.ProtoReflect.Descriptor instead.
func (*ListResourcesResponse_ResourceMetadata) Descriptor() ([]byte, []int) {
	return file_api_proto_commonpb_common_proto_rawDescGZIP(), []int{8, 0}
}

func (x *ListResourcesResponse_ResourceMetadata) GetResourceId() string {
	if x != nil {
		return x.ResourceId
	}
	return ""
}

func (x *ListResourcesResponse_ResourceMetadata) GetResourceType() string {
	if x != nil {
		return x.ResourceType
	}
	return ""
}

var File_api_proto_commonpb_common_proto protoreflect.FileDescriptor

var file_api_proto_commonpb_common_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x70, 0x62, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x10, 0x64, 0x65, 0x76, 0x2e, 0x6d, 0x75, 0x6c, 0x74, 0x79, 0x2e, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x22, 0x34, 0x0a, 0x0f, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6e, 0x67,
	0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x22, 0x21, 0x0a, 0x04, 0x45, 0x6e, 0x75, 0x6d, 0x12, 0x0e,
	0x0a, 0x0a, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x5f, 0x4f, 0x53, 0x10, 0x00, 0x12, 0x09,
	0x0a, 0x05, 0x4c, 0x49, 0x4e, 0x55, 0x58, 0x10, 0x01, 0x22, 0x4d, 0x0a, 0x0c, 0x44, 0x61, 0x74,
	0x61, 0x62, 0x61, 0x73, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x22, 0x3d, 0x0a, 0x04, 0x45, 0x6e, 0x75,
	0x6d, 0x12, 0x13, 0x0a, 0x0f, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x5f, 0x56, 0x4d, 0x5f,
	0x53, 0x49, 0x5a, 0x45, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x4d, 0x49, 0x43, 0x52, 0x4f, 0x10,
	0x02, 0x12, 0x09, 0x0a, 0x05, 0x53, 0x4d, 0x41, 0x4c, 0x4c, 0x10, 0x04, 0x12, 0x0a, 0x0a, 0x06,
	0x4d, 0x45, 0x44, 0x49, 0x55, 0x4d, 0x10, 0x03, 0x22, 0xa8, 0x02, 0x0a, 0x06, 0x56, 0x6d, 0x53,
	0x69, 0x7a, 0x65, 0x22, 0x9d, 0x02, 0x0a, 0x04, 0x45, 0x6e, 0x75, 0x6d, 0x12, 0x13, 0x0a, 0x0f,
	0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x5f, 0x56, 0x4d, 0x5f, 0x53, 0x49, 0x5a, 0x45, 0x10,
	0x00, 0x12, 0x09, 0x0a, 0x05, 0x4d, 0x49, 0x43, 0x52, 0x4f, 0x10, 0x01, 0x12, 0x0a, 0x0a, 0x06,
	0x4d, 0x45, 0x44, 0x49, 0x55, 0x4d, 0x10, 0x02, 0x12, 0x09, 0x0a, 0x05, 0x4c, 0x41, 0x52, 0x47,
	0x45, 0x10, 0x03, 0x12, 0x10, 0x0a, 0x0c, 0x47, 0x45, 0x4e, 0x45, 0x52, 0x41, 0x4c, 0x5f, 0x4e,
	0x41, 0x4e, 0x4f, 0x10, 0x04, 0x12, 0x11, 0x0a, 0x0d, 0x47, 0x45, 0x4e, 0x45, 0x52, 0x41, 0x4c,
	0x5f, 0x4c, 0x41, 0x52, 0x47, 0x45, 0x10, 0x05, 0x12, 0x12, 0x0a, 0x0e, 0x47, 0x45, 0x4e, 0x45,
	0x52, 0x41, 0x4c, 0x5f, 0x58, 0x4c, 0x41, 0x52, 0x47, 0x45, 0x10, 0x06, 0x12, 0x13, 0x0a, 0x0f,
	0x47, 0x45, 0x4e, 0x45, 0x52, 0x41, 0x4c, 0x5f, 0x32, 0x58, 0x4c, 0x41, 0x52, 0x47, 0x45, 0x10,
	0x07, 0x12, 0x11, 0x0a, 0x0d, 0x43, 0x4f, 0x4d, 0x50, 0x55, 0x54, 0x45, 0x5f, 0x4c, 0x41, 0x52,
	0x47, 0x45, 0x10, 0x08, 0x12, 0x12, 0x0a, 0x0e, 0x43, 0x4f, 0x4d, 0x50, 0x55, 0x54, 0x45, 0x5f,
	0x58, 0x4c, 0x41, 0x52, 0x47, 0x45, 0x10, 0x09, 0x12, 0x13, 0x0a, 0x0f, 0x43, 0x4f, 0x4d, 0x50,
	0x55, 0x54, 0x45, 0x5f, 0x32, 0x58, 0x4c, 0x41, 0x52, 0x47, 0x45, 0x10, 0x0a, 0x12, 0x13, 0x0a,
	0x0f, 0x43, 0x4f, 0x4d, 0x50, 0x55, 0x54, 0x45, 0x5f, 0x34, 0x58, 0x4c, 0x41, 0x52, 0x47, 0x45,
	0x10, 0x0b, 0x12, 0x13, 0x0a, 0x0f, 0x43, 0x4f, 0x4d, 0x50, 0x55, 0x54, 0x45, 0x5f, 0x38, 0x58,
	0x4c, 0x41, 0x52, 0x47, 0x45, 0x10, 0x0c, 0x12, 0x14, 0x0a, 0x10, 0x43, 0x4f, 0x4d, 0x50, 0x55,
	0x54, 0x45, 0x5f, 0x31, 0x32, 0x58, 0x4c, 0x41, 0x52, 0x47, 0x45, 0x10, 0x0d, 0x12, 0x14, 0x0a,
	0x10, 0x43, 0x4f, 0x4d, 0x50, 0x55, 0x54, 0x45, 0x5f, 0x31, 0x36, 0x58, 0x4c, 0x41, 0x52, 0x47,
	0x45, 0x10, 0x0e, 0x22, 0x8a, 0x02, 0x0a, 0x18, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x52, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73,
	0x12, 0x1f, 0x0a, 0x0b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49,
	0x64, 0x12, 0x2a, 0x0a, 0x11, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x67, 0x72,
	0x6f, 0x75, 0x70, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x64, 0x12, 0x36, 0x0a,
	0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x1a, 0x2e, 0x64, 0x65, 0x76, 0x2e, 0x6d, 0x75, 0x6c, 0x74, 0x79, 0x2e, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2e, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x08, 0x6c, 0x6f, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x46, 0x0a, 0x0e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x5f, 0x70,
	0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1f, 0x2e,
	0x64, 0x65, 0x76, 0x2e, 0x6d, 0x75, 0x6c, 0x74, 0x79, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x2e, 0x43, 0x6c, 0x6f, 0x75, 0x64, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x52, 0x0d,
	0x63, 0x6c, 0x6f, 0x75, 0x64, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x12, 0x21, 0x0a,
	0x0c, 0x6e, 0x65, 0x65, 0x64, 0x73, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x0b, 0x6e, 0x65, 0x65, 0x64, 0x73, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x22, 0x63, 0x0a, 0x1d, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x43, 0x68, 0x69, 0x6c, 0x64, 0x52,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72,
	0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x49, 0x64, 0x12, 0x21, 0x0a, 0x0c, 0x6e, 0x65, 0x65, 0x64, 0x73, 0x5f, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x6e, 0x65, 0x65, 0x64, 0x73, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x22, 0xc0, 0x01, 0x0a, 0x12, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x41, 0x72, 0x67, 0x73, 0x12, 0x2a, 0x0a, 0x11,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x64, 0x12, 0x36, 0x0a, 0x08, 0x6c, 0x6f, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1a, 0x2e, 0x64, 0x65, 0x76,
	0x2e, 0x6d, 0x75, 0x6c, 0x74, 0x79, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x4c, 0x6f,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x46, 0x0a, 0x0e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x5f, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64,
	0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1f, 0x2e, 0x64, 0x65, 0x76, 0x2e, 0x6d,
	0x75, 0x6c, 0x74, 0x79, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x43, 0x6c, 0x6f, 0x75,
	0x64, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x52, 0x0d, 0x63, 0x6c, 0x6f, 0x75, 0x64,
	0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x22, 0x19, 0x0a, 0x17, 0x43, 0x68, 0x69, 0x6c,
	0x64, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x41,
	0x72, 0x67, 0x73, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0xc9, 0x01, 0x0a,
	0x15, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x56, 0x0a, 0x09, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x38, 0x2e, 0x64, 0x65, 0x76, 0x2e,
	0x6d, 0x75, 0x6c, 0x74, 0x79, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x4c, 0x69, 0x73,
	0x74, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x2e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x52, 0x09, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x1a, 0x58,
	0x0a, 0x10, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0x12, 0x1f, 0x0a, 0x0b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x49, 0x64, 0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65, 0x2a, 0x99, 0x01, 0x0a, 0x08, 0x4c, 0x6f, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x10, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e,
	0x5f, 0x4c, 0x4f, 0x43, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x10, 0x00, 0x12, 0x0d, 0x0a, 0x09, 0x55,
	0x53, 0x5f, 0x45, 0x41, 0x53, 0x54, 0x5f, 0x31, 0x10, 0x01, 0x12, 0x0d, 0x0a, 0x09, 0x55, 0x53,
	0x5f, 0x45, 0x41, 0x53, 0x54, 0x5f, 0x32, 0x10, 0x02, 0x12, 0x0d, 0x0a, 0x09, 0x55, 0x53, 0x5f,
	0x57, 0x45, 0x53, 0x54, 0x5f, 0x31, 0x10, 0x03, 0x12, 0x0d, 0x0a, 0x09, 0x55, 0x53, 0x5f, 0x57,
	0x45, 0x53, 0x54, 0x5f, 0x32, 0x10, 0x04, 0x12, 0x0d, 0x0a, 0x09, 0x45, 0x55, 0x5f, 0x57, 0x45,
	0x53, 0x54, 0x5f, 0x31, 0x10, 0x05, 0x12, 0x0d, 0x0a, 0x09, 0x45, 0x55, 0x5f, 0x57, 0x45, 0x53,
	0x54, 0x5f, 0x32, 0x10, 0x06, 0x12, 0x0d, 0x0a, 0x09, 0x45, 0x55, 0x5f, 0x57, 0x45, 0x53, 0x54,
	0x5f, 0x33, 0x10, 0x07, 0x12, 0x0e, 0x0a, 0x0a, 0x45, 0x55, 0x5f, 0x4e, 0x4f, 0x52, 0x54, 0x48,
	0x5f, 0x31, 0x10, 0x08, 0x2a, 0x39, 0x0a, 0x0d, 0x43, 0x6c, 0x6f, 0x75, 0x64, 0x50, 0x72, 0x6f,
	0x76, 0x69, 0x64, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x10, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e,
	0x5f, 0x50, 0x52, 0x4f, 0x56, 0x49, 0x44, 0x45, 0x52, 0x10, 0x00, 0x12, 0x07, 0x0a, 0x03, 0x41,
	0x57, 0x53, 0x10, 0x01, 0x12, 0x09, 0x0a, 0x05, 0x41, 0x5a, 0x55, 0x52, 0x45, 0x10, 0x02, 0x42,
	0x54, 0x0a, 0x14, 0x64, 0x65, 0x76, 0x2e, 0x6d, 0x75, 0x6c, 0x74, 0x79, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x42, 0x0a, 0x4d, 0x75, 0x6c, 0x74, 0x79, 0x50, 0x72,
	0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x6d, 0x75, 0x6c, 0x74, 0x79, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2f, 0x6d, 0x75, 0x6c,
	0x74, 0x79, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_proto_commonpb_common_proto_rawDescOnce sync.Once
	file_api_proto_commonpb_common_proto_rawDescData = file_api_proto_commonpb_common_proto_rawDesc
)

func file_api_proto_commonpb_common_proto_rawDescGZIP() []byte {
	file_api_proto_commonpb_common_proto_rawDescOnce.Do(func() {
		file_api_proto_commonpb_common_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_proto_commonpb_common_proto_rawDescData)
	})
	return file_api_proto_commonpb_common_proto_rawDescData
}

var file_api_proto_commonpb_common_proto_enumTypes = make([]protoimpl.EnumInfo, 5)
var file_api_proto_commonpb_common_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_api_proto_commonpb_common_proto_goTypes = []interface{}{
	(Location)(0),                                  // 0: dev.multy.common.Location
	(CloudProvider)(0),                             // 1: dev.multy.common.CloudProvider
	(OperatingSystem_Enum)(0),                      // 2: dev.multy.common.OperatingSystem.Enum
	(DatabaseSize_Enum)(0),                         // 3: dev.multy.common.DatabaseSize.Enum
	(VmSize_Enum)(0),                               // 4: dev.multy.common.VmSize.Enum
	(*OperatingSystem)(nil),                        // 5: dev.multy.common.OperatingSystem
	(*DatabaseSize)(nil),                           // 6: dev.multy.common.DatabaseSize
	(*VmSize)(nil),                                 // 7: dev.multy.common.VmSize
	(*CommonResourceParameters)(nil),               // 8: dev.multy.common.CommonResourceParameters
	(*CommonChildResourceParameters)(nil),          // 9: dev.multy.common.CommonChildResourceParameters
	(*ResourceCommonArgs)(nil),                     // 10: dev.multy.common.ResourceCommonArgs
	(*ChildResourceCommonArgs)(nil),                // 11: dev.multy.common.ChildResourceCommonArgs
	(*Empty)(nil),                                  // 12: dev.multy.common.Empty
	(*ListResourcesResponse)(nil),                  // 13: dev.multy.common.ListResourcesResponse
	(*ListResourcesResponse_ResourceMetadata)(nil), // 14: dev.multy.common.ListResourcesResponse.ResourceMetadata
}
var file_api_proto_commonpb_common_proto_depIdxs = []int32{
	0,  // 0: dev.multy.common.CommonResourceParameters.location:type_name -> dev.multy.common.Location
	1,  // 1: dev.multy.common.CommonResourceParameters.cloud_provider:type_name -> dev.multy.common.CloudProvider
	0,  // 2: dev.multy.common.ResourceCommonArgs.location:type_name -> dev.multy.common.Location
	1,  // 3: dev.multy.common.ResourceCommonArgs.cloud_provider:type_name -> dev.multy.common.CloudProvider
	14, // 4: dev.multy.common.ListResourcesResponse.resources:type_name -> dev.multy.common.ListResourcesResponse.ResourceMetadata
	5,  // [5:5] is the sub-list for method output_type
	5,  // [5:5] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_api_proto_commonpb_common_proto_init() }
func file_api_proto_commonpb_common_proto_init() {
	if File_api_proto_commonpb_common_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_proto_commonpb_common_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OperatingSystem); i {
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
		file_api_proto_commonpb_common_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DatabaseSize); i {
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
		file_api_proto_commonpb_common_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VmSize); i {
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
		file_api_proto_commonpb_common_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CommonResourceParameters); i {
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
		file_api_proto_commonpb_common_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CommonChildResourceParameters); i {
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
		file_api_proto_commonpb_common_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResourceCommonArgs); i {
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
		file_api_proto_commonpb_common_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChildResourceCommonArgs); i {
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
		file_api_proto_commonpb_common_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
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
		file_api_proto_commonpb_common_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListResourcesResponse); i {
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
		file_api_proto_commonpb_common_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListResourcesResponse_ResourceMetadata); i {
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
			RawDescriptor: file_api_proto_commonpb_common_proto_rawDesc,
			NumEnums:      5,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_proto_commonpb_common_proto_goTypes,
		DependencyIndexes: file_api_proto_commonpb_common_proto_depIdxs,
		EnumInfos:         file_api_proto_commonpb_common_proto_enumTypes,
		MessageInfos:      file_api_proto_commonpb_common_proto_msgTypes,
	}.Build()
	File_api_proto_commonpb_common_proto = out.File
	file_api_proto_commonpb_common_proto_rawDesc = nil
	file_api_proto_commonpb_common_proto_goTypes = nil
	file_api_proto_commonpb_common_proto_depIdxs = nil
}
