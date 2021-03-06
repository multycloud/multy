// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: api/proto/resourcespb/network_interface.proto

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

type CreateNetworkInterfaceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Resource *NetworkInterfaceArgs `protobuf:"bytes,1,opt,name=resource,proto3" json:"resource,omitempty"`
}

func (x *CreateNetworkInterfaceRequest) Reset() {
	*x = CreateNetworkInterfaceRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_resourcespb_network_interface_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateNetworkInterfaceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateNetworkInterfaceRequest) ProtoMessage() {}

func (x *CreateNetworkInterfaceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_resourcespb_network_interface_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateNetworkInterfaceRequest.ProtoReflect.Descriptor instead.
func (*CreateNetworkInterfaceRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_resourcespb_network_interface_proto_rawDescGZIP(), []int{0}
}

func (x *CreateNetworkInterfaceRequest) GetResource() *NetworkInterfaceArgs {
	if x != nil {
		return x.Resource
	}
	return nil
}

type ReadNetworkInterfaceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ResourceId string `protobuf:"bytes,1,opt,name=resource_id,json=resourceId,proto3" json:"resource_id,omitempty"`
}

func (x *ReadNetworkInterfaceRequest) Reset() {
	*x = ReadNetworkInterfaceRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_resourcespb_network_interface_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadNetworkInterfaceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadNetworkInterfaceRequest) ProtoMessage() {}

func (x *ReadNetworkInterfaceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_resourcespb_network_interface_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadNetworkInterfaceRequest.ProtoReflect.Descriptor instead.
func (*ReadNetworkInterfaceRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_resourcespb_network_interface_proto_rawDescGZIP(), []int{1}
}

func (x *ReadNetworkInterfaceRequest) GetResourceId() string {
	if x != nil {
		return x.ResourceId
	}
	return ""
}

type UpdateNetworkInterfaceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ResourceId string                `protobuf:"bytes,1,opt,name=resource_id,json=resourceId,proto3" json:"resource_id,omitempty"`
	Resource   *NetworkInterfaceArgs `protobuf:"bytes,2,opt,name=resource,proto3" json:"resource,omitempty"`
}

func (x *UpdateNetworkInterfaceRequest) Reset() {
	*x = UpdateNetworkInterfaceRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_resourcespb_network_interface_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateNetworkInterfaceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateNetworkInterfaceRequest) ProtoMessage() {}

func (x *UpdateNetworkInterfaceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_resourcespb_network_interface_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateNetworkInterfaceRequest.ProtoReflect.Descriptor instead.
func (*UpdateNetworkInterfaceRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_resourcespb_network_interface_proto_rawDescGZIP(), []int{2}
}

func (x *UpdateNetworkInterfaceRequest) GetResourceId() string {
	if x != nil {
		return x.ResourceId
	}
	return ""
}

func (x *UpdateNetworkInterfaceRequest) GetResource() *NetworkInterfaceArgs {
	if x != nil {
		return x.Resource
	}
	return nil
}

type DeleteNetworkInterfaceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ResourceId string `protobuf:"bytes,1,opt,name=resource_id,json=resourceId,proto3" json:"resource_id,omitempty"`
}

func (x *DeleteNetworkInterfaceRequest) Reset() {
	*x = DeleteNetworkInterfaceRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_resourcespb_network_interface_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteNetworkInterfaceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteNetworkInterfaceRequest) ProtoMessage() {}

func (x *DeleteNetworkInterfaceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_resourcespb_network_interface_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteNetworkInterfaceRequest.ProtoReflect.Descriptor instead.
func (*DeleteNetworkInterfaceRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_resourcespb_network_interface_proto_rawDescGZIP(), []int{3}
}

func (x *DeleteNetworkInterfaceRequest) GetResourceId() string {
	if x != nil {
		return x.ResourceId
	}
	return ""
}

type NetworkInterfaceArgs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CommonParameters *commonpb.ResourceCommonArgs `protobuf:"bytes,1,opt,name=common_parameters,json=commonParameters,proto3" json:"common_parameters,omitempty"`
	Name             string                       `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	SubnetId         string                       `protobuf:"bytes,3,opt,name=subnet_id,json=subnetId,proto3" json:"subnet_id,omitempty"`
	PublicIpId       string                       `protobuf:"bytes,4,opt,name=public_ip_id,json=publicIpId,proto3" json:"public_ip_id,omitempty"`
	AvailabilityZone int32                        `protobuf:"varint,5,opt,name=availability_zone,json=availabilityZone,proto3" json:"availability_zone,omitempty"`
}

func (x *NetworkInterfaceArgs) Reset() {
	*x = NetworkInterfaceArgs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_resourcespb_network_interface_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NetworkInterfaceArgs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NetworkInterfaceArgs) ProtoMessage() {}

func (x *NetworkInterfaceArgs) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_resourcespb_network_interface_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NetworkInterfaceArgs.ProtoReflect.Descriptor instead.
func (*NetworkInterfaceArgs) Descriptor() ([]byte, []int) {
	return file_api_proto_resourcespb_network_interface_proto_rawDescGZIP(), []int{4}
}

func (x *NetworkInterfaceArgs) GetCommonParameters() *commonpb.ResourceCommonArgs {
	if x != nil {
		return x.CommonParameters
	}
	return nil
}

func (x *NetworkInterfaceArgs) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *NetworkInterfaceArgs) GetSubnetId() string {
	if x != nil {
		return x.SubnetId
	}
	return ""
}

func (x *NetworkInterfaceArgs) GetPublicIpId() string {
	if x != nil {
		return x.PublicIpId
	}
	return ""
}

func (x *NetworkInterfaceArgs) GetAvailabilityZone() int32 {
	if x != nil {
		return x.AvailabilityZone
	}
	return 0
}

type NetworkInterfaceAwsOutputs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NetworkInterfaceId string `protobuf:"bytes,1,opt,name=network_interface_id,json=networkInterfaceId,proto3" json:"network_interface_id,omitempty"`
	EipAssociationId   string `protobuf:"bytes,2,opt,name=eip_association_id,json=eipAssociationId,proto3" json:"eip_association_id,omitempty"`
}

func (x *NetworkInterfaceAwsOutputs) Reset() {
	*x = NetworkInterfaceAwsOutputs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_resourcespb_network_interface_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NetworkInterfaceAwsOutputs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NetworkInterfaceAwsOutputs) ProtoMessage() {}

func (x *NetworkInterfaceAwsOutputs) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_resourcespb_network_interface_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NetworkInterfaceAwsOutputs.ProtoReflect.Descriptor instead.
func (*NetworkInterfaceAwsOutputs) Descriptor() ([]byte, []int) {
	return file_api_proto_resourcespb_network_interface_proto_rawDescGZIP(), []int{5}
}

func (x *NetworkInterfaceAwsOutputs) GetNetworkInterfaceId() string {
	if x != nil {
		return x.NetworkInterfaceId
	}
	return ""
}

func (x *NetworkInterfaceAwsOutputs) GetEipAssociationId() string {
	if x != nil {
		return x.EipAssociationId
	}
	return ""
}

type NetworkInterfaceAzureOutputs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NetworkInterfaceId string `protobuf:"bytes,1,opt,name=network_interface_id,json=networkInterfaceId,proto3" json:"network_interface_id,omitempty"`
}

func (x *NetworkInterfaceAzureOutputs) Reset() {
	*x = NetworkInterfaceAzureOutputs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_resourcespb_network_interface_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NetworkInterfaceAzureOutputs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NetworkInterfaceAzureOutputs) ProtoMessage() {}

func (x *NetworkInterfaceAzureOutputs) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_resourcespb_network_interface_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NetworkInterfaceAzureOutputs.ProtoReflect.Descriptor instead.
func (*NetworkInterfaceAzureOutputs) Descriptor() ([]byte, []int) {
	return file_api_proto_resourcespb_network_interface_proto_rawDescGZIP(), []int{6}
}

func (x *NetworkInterfaceAzureOutputs) GetNetworkInterfaceId() string {
	if x != nil {
		return x.NetworkInterfaceId
	}
	return ""
}

type NetworkInterfaceResource struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CommonParameters *commonpb.CommonResourceParameters `protobuf:"bytes,1,opt,name=common_parameters,json=commonParameters,proto3" json:"common_parameters,omitempty"`
	Name             string                             `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	SubnetId         string                             `protobuf:"bytes,3,opt,name=subnet_id,json=subnetId,proto3" json:"subnet_id,omitempty"`
	PublicIpId       string                             `protobuf:"bytes,4,opt,name=public_ip_id,json=publicIpId,proto3" json:"public_ip_id,omitempty"`
	AvailabilityZone int32                              `protobuf:"varint,5,opt,name=availability_zone,json=availabilityZone,proto3" json:"availability_zone,omitempty"`
	AwsOutputs       *NetworkInterfaceAwsOutputs        `protobuf:"bytes,6,opt,name=aws_outputs,json=awsOutputs,proto3" json:"aws_outputs,omitempty"`
	AzureOutputs     *NetworkInterfaceAzureOutputs      `protobuf:"bytes,7,opt,name=azure_outputs,json=azureOutputs,proto3" json:"azure_outputs,omitempty"`
}

func (x *NetworkInterfaceResource) Reset() {
	*x = NetworkInterfaceResource{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_resourcespb_network_interface_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NetworkInterfaceResource) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NetworkInterfaceResource) ProtoMessage() {}

func (x *NetworkInterfaceResource) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_resourcespb_network_interface_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NetworkInterfaceResource.ProtoReflect.Descriptor instead.
func (*NetworkInterfaceResource) Descriptor() ([]byte, []int) {
	return file_api_proto_resourcespb_network_interface_proto_rawDescGZIP(), []int{7}
}

func (x *NetworkInterfaceResource) GetCommonParameters() *commonpb.CommonResourceParameters {
	if x != nil {
		return x.CommonParameters
	}
	return nil
}

func (x *NetworkInterfaceResource) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *NetworkInterfaceResource) GetSubnetId() string {
	if x != nil {
		return x.SubnetId
	}
	return ""
}

func (x *NetworkInterfaceResource) GetPublicIpId() string {
	if x != nil {
		return x.PublicIpId
	}
	return ""
}

func (x *NetworkInterfaceResource) GetAvailabilityZone() int32 {
	if x != nil {
		return x.AvailabilityZone
	}
	return 0
}

func (x *NetworkInterfaceResource) GetAwsOutputs() *NetworkInterfaceAwsOutputs {
	if x != nil {
		return x.AwsOutputs
	}
	return nil
}

func (x *NetworkInterfaceResource) GetAzureOutputs() *NetworkInterfaceAzureOutputs {
	if x != nil {
		return x.AzureOutputs
	}
	return nil
}

var File_api_proto_resourcespb_network_interface_proto protoreflect.FileDescriptor

var file_api_proto_resourcespb_network_interface_proto_rawDesc = []byte{
	0x0a, 0x2d, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x73, 0x70, 0x62, 0x2f, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x5f,
	0x69, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x13, 0x64, 0x65, 0x76, 0x2e, 0x6d, 0x75, 0x6c, 0x74, 0x79, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x1a, 0x1f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x70, 0x62, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x66, 0x0a, 0x1d, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4e,
	0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x45, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x64, 0x65, 0x76, 0x2e, 0x6d,
	0x75, 0x6c, 0x74, 0x79, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x4e,
	0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x41,
	0x72, 0x67, 0x73, 0x52, 0x08, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x22, 0x3e, 0x0a,
	0x1b, 0x52, 0x65, 0x61, 0x64, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x49, 0x6e, 0x74, 0x65,
	0x72, 0x66, 0x61, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0a, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64, 0x22, 0x87, 0x01,
	0x0a, 0x1d, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x49,
	0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x1f, 0x0a, 0x0b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64,
	0x12, 0x45, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x29, 0x2e, 0x64, 0x65, 0x76, 0x2e, 0x6d, 0x75, 0x6c, 0x74, 0x79, 0x2e, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b,
	0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x41, 0x72, 0x67, 0x73, 0x52, 0x08, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x22, 0x40, 0x0a, 0x1d, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64, 0x22, 0xe9, 0x01, 0x0a, 0x14, 0x4e, 0x65,
	0x74, 0x77, 0x6f, 0x72, 0x6b, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x41, 0x72,
	0x67, 0x73, 0x12, 0x51, 0x0a, 0x11, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5f, 0x70, 0x61, 0x72,
	0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e,
	0x64, 0x65, 0x76, 0x2e, 0x6d, 0x75, 0x6c, 0x74, 0x79, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x2e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x41,
	0x72, 0x67, 0x73, 0x52, 0x10, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x50, 0x61, 0x72, 0x61, 0x6d,
	0x65, 0x74, 0x65, 0x72, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x75, 0x62,
	0x6e, 0x65, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x75,
	0x62, 0x6e, 0x65, 0x74, 0x49, 0x64, 0x12, 0x20, 0x0a, 0x0c, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63,
	0x5f, 0x69, 0x70, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x70, 0x75,
	0x62, 0x6c, 0x69, 0x63, 0x49, 0x70, 0x49, 0x64, 0x12, 0x2b, 0x0a, 0x11, 0x61, 0x76, 0x61, 0x69,
	0x6c, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x5f, 0x7a, 0x6f, 0x6e, 0x65, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x10, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74,
	0x79, 0x5a, 0x6f, 0x6e, 0x65, 0x22, 0x7c, 0x0a, 0x1a, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b,
	0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x41, 0x77, 0x73, 0x4f, 0x75, 0x74, 0x70,
	0x75, 0x74, 0x73, 0x12, 0x30, 0x0a, 0x14, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x5f, 0x69,
	0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x12, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66,
	0x61, 0x63, 0x65, 0x49, 0x64, 0x12, 0x2c, 0x0a, 0x12, 0x65, 0x69, 0x70, 0x5f, 0x61, 0x73, 0x73,
	0x6f, 0x63, 0x69, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x10, 0x65, 0x69, 0x70, 0x41, 0x73, 0x73, 0x6f, 0x63, 0x69, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x49, 0x64, 0x22, 0x50, 0x0a, 0x1c, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x49, 0x6e,
	0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x41, 0x7a, 0x75, 0x72, 0x65, 0x4f, 0x75, 0x74, 0x70,
	0x75, 0x74, 0x73, 0x12, 0x30, 0x0a, 0x14, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x5f, 0x69,
	0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x12, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66,
	0x61, 0x63, 0x65, 0x49, 0x64, 0x22, 0x9d, 0x03, 0x0a, 0x18, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72,
	0x6b, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x12, 0x57, 0x0a, 0x11, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5f, 0x70, 0x61, 0x72,
	0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2a, 0x2e,
	0x64, 0x65, 0x76, 0x2e, 0x6d, 0x75, 0x6c, 0x74, 0x79, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x50,
	0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x52, 0x10, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x1b, 0x0a, 0x09, 0x73, 0x75, 0x62, 0x6e, 0x65, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x73, 0x75, 0x62, 0x6e, 0x65, 0x74, 0x49, 0x64, 0x12, 0x20, 0x0a, 0x0c,
	0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x69, 0x70, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x49, 0x70, 0x49, 0x64, 0x12, 0x2b,
	0x0a, 0x11, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x5f, 0x7a,
	0x6f, 0x6e, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x10, 0x61, 0x76, 0x61, 0x69, 0x6c,
	0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x5a, 0x6f, 0x6e, 0x65, 0x12, 0x50, 0x0a, 0x0b, 0x61,
	0x77, 0x73, 0x5f, 0x6f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x2f, 0x2e, 0x64, 0x65, 0x76, 0x2e, 0x6d, 0x75, 0x6c, 0x74, 0x79, 0x2e, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x49, 0x6e,
	0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x41, 0x77, 0x73, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74,
	0x73, 0x52, 0x0a, 0x61, 0x77, 0x73, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x73, 0x12, 0x56, 0x0a,
	0x0d, 0x61, 0x7a, 0x75, 0x72, 0x65, 0x5f, 0x6f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x73, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x31, 0x2e, 0x64, 0x65, 0x76, 0x2e, 0x6d, 0x75, 0x6c, 0x74, 0x79,
	0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x4e, 0x65, 0x74, 0x77, 0x6f,
	0x72, 0x6b, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x41, 0x7a, 0x75, 0x72, 0x65,
	0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x73, 0x52, 0x0c, 0x61, 0x7a, 0x75, 0x72, 0x65, 0x4f, 0x75,
	0x74, 0x70, 0x75, 0x74, 0x73, 0x42, 0x5e, 0x0a, 0x17, 0x64, 0x65, 0x76, 0x2e, 0x6d, 0x75, 0x6c,
	0x74, 0x79, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73,
	0x42, 0x0e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x50, 0x72, 0x6f, 0x74, 0x6f,
	0x50, 0x01, 0x5a, 0x31, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d,
	0x75, 0x6c, 0x74, 0x79, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2f, 0x6d, 0x75, 0x6c, 0x74, 0x79, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x73, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_proto_resourcespb_network_interface_proto_rawDescOnce sync.Once
	file_api_proto_resourcespb_network_interface_proto_rawDescData = file_api_proto_resourcespb_network_interface_proto_rawDesc
)

func file_api_proto_resourcespb_network_interface_proto_rawDescGZIP() []byte {
	file_api_proto_resourcespb_network_interface_proto_rawDescOnce.Do(func() {
		file_api_proto_resourcespb_network_interface_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_proto_resourcespb_network_interface_proto_rawDescData)
	})
	return file_api_proto_resourcespb_network_interface_proto_rawDescData
}

var file_api_proto_resourcespb_network_interface_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_api_proto_resourcespb_network_interface_proto_goTypes = []interface{}{
	(*CreateNetworkInterfaceRequest)(nil),     // 0: dev.multy.resources.CreateNetworkInterfaceRequest
	(*ReadNetworkInterfaceRequest)(nil),       // 1: dev.multy.resources.ReadNetworkInterfaceRequest
	(*UpdateNetworkInterfaceRequest)(nil),     // 2: dev.multy.resources.UpdateNetworkInterfaceRequest
	(*DeleteNetworkInterfaceRequest)(nil),     // 3: dev.multy.resources.DeleteNetworkInterfaceRequest
	(*NetworkInterfaceArgs)(nil),              // 4: dev.multy.resources.NetworkInterfaceArgs
	(*NetworkInterfaceAwsOutputs)(nil),        // 5: dev.multy.resources.NetworkInterfaceAwsOutputs
	(*NetworkInterfaceAzureOutputs)(nil),      // 6: dev.multy.resources.NetworkInterfaceAzureOutputs
	(*NetworkInterfaceResource)(nil),          // 7: dev.multy.resources.NetworkInterfaceResource
	(*commonpb.ResourceCommonArgs)(nil),       // 8: dev.multy.common.ResourceCommonArgs
	(*commonpb.CommonResourceParameters)(nil), // 9: dev.multy.common.CommonResourceParameters
}
var file_api_proto_resourcespb_network_interface_proto_depIdxs = []int32{
	4, // 0: dev.multy.resources.CreateNetworkInterfaceRequest.resource:type_name -> dev.multy.resources.NetworkInterfaceArgs
	4, // 1: dev.multy.resources.UpdateNetworkInterfaceRequest.resource:type_name -> dev.multy.resources.NetworkInterfaceArgs
	8, // 2: dev.multy.resources.NetworkInterfaceArgs.common_parameters:type_name -> dev.multy.common.ResourceCommonArgs
	9, // 3: dev.multy.resources.NetworkInterfaceResource.common_parameters:type_name -> dev.multy.common.CommonResourceParameters
	5, // 4: dev.multy.resources.NetworkInterfaceResource.aws_outputs:type_name -> dev.multy.resources.NetworkInterfaceAwsOutputs
	6, // 5: dev.multy.resources.NetworkInterfaceResource.azure_outputs:type_name -> dev.multy.resources.NetworkInterfaceAzureOutputs
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_api_proto_resourcespb_network_interface_proto_init() }
func file_api_proto_resourcespb_network_interface_proto_init() {
	if File_api_proto_resourcespb_network_interface_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_proto_resourcespb_network_interface_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateNetworkInterfaceRequest); i {
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
		file_api_proto_resourcespb_network_interface_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReadNetworkInterfaceRequest); i {
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
		file_api_proto_resourcespb_network_interface_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateNetworkInterfaceRequest); i {
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
		file_api_proto_resourcespb_network_interface_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteNetworkInterfaceRequest); i {
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
		file_api_proto_resourcespb_network_interface_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NetworkInterfaceArgs); i {
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
		file_api_proto_resourcespb_network_interface_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NetworkInterfaceAwsOutputs); i {
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
		file_api_proto_resourcespb_network_interface_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NetworkInterfaceAzureOutputs); i {
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
		file_api_proto_resourcespb_network_interface_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NetworkInterfaceResource); i {
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
			RawDescriptor: file_api_proto_resourcespb_network_interface_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_proto_resourcespb_network_interface_proto_goTypes,
		DependencyIndexes: file_api_proto_resourcespb_network_interface_proto_depIdxs,
		MessageInfos:      file_api_proto_resourcespb_network_interface_proto_msgTypes,
	}.Build()
	File_api_proto_resourcespb_network_interface_proto = out.File
	file_api_proto_resourcespb_network_interface_proto_rawDesc = nil
	file_api_proto_resourcespb_network_interface_proto_goTypes = nil
	file_api_proto_resourcespb_network_interface_proto_depIdxs = nil
}
