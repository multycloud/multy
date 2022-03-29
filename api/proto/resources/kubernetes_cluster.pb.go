// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.6.1
// source: api/proto/resources/kubernetes_cluster.proto

package resources

import (
	common "github.com/multycloud/multy/api/proto/common"
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

type CreateKubernetesClusterRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Resource *KubernetesClusterArgs `protobuf:"bytes,1,opt,name=resource,proto3" json:"resource,omitempty"`
}

func (x *CreateKubernetesClusterRequest) Reset() {
	*x = CreateKubernetesClusterRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_resources_kubernetes_cluster_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateKubernetesClusterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateKubernetesClusterRequest) ProtoMessage() {}

func (x *CreateKubernetesClusterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_resources_kubernetes_cluster_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateKubernetesClusterRequest.ProtoReflect.Descriptor instead.
func (*CreateKubernetesClusterRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_resources_kubernetes_cluster_proto_rawDescGZIP(), []int{0}
}

func (x *CreateKubernetesClusterRequest) GetResource() *KubernetesClusterArgs {
	if x != nil {
		return x.Resource
	}
	return nil
}

type ReadKubernetesClusterRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ResourceId string `protobuf:"bytes,1,opt,name=resource_id,json=resourceId,proto3" json:"resource_id,omitempty"`
}

func (x *ReadKubernetesClusterRequest) Reset() {
	*x = ReadKubernetesClusterRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_resources_kubernetes_cluster_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadKubernetesClusterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadKubernetesClusterRequest) ProtoMessage() {}

func (x *ReadKubernetesClusterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_resources_kubernetes_cluster_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadKubernetesClusterRequest.ProtoReflect.Descriptor instead.
func (*ReadKubernetesClusterRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_resources_kubernetes_cluster_proto_rawDescGZIP(), []int{1}
}

func (x *ReadKubernetesClusterRequest) GetResourceId() string {
	if x != nil {
		return x.ResourceId
	}
	return ""
}

type UpdateKubernetesClusterRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ResourceId string                 `protobuf:"bytes,1,opt,name=resource_id,json=resourceId,proto3" json:"resource_id,omitempty"`
	Resource   *KubernetesClusterArgs `protobuf:"bytes,2,opt,name=resource,proto3" json:"resource,omitempty"`
}

func (x *UpdateKubernetesClusterRequest) Reset() {
	*x = UpdateKubernetesClusterRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_resources_kubernetes_cluster_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateKubernetesClusterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateKubernetesClusterRequest) ProtoMessage() {}

func (x *UpdateKubernetesClusterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_resources_kubernetes_cluster_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateKubernetesClusterRequest.ProtoReflect.Descriptor instead.
func (*UpdateKubernetesClusterRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_resources_kubernetes_cluster_proto_rawDescGZIP(), []int{2}
}

func (x *UpdateKubernetesClusterRequest) GetResourceId() string {
	if x != nil {
		return x.ResourceId
	}
	return ""
}

func (x *UpdateKubernetesClusterRequest) GetResource() *KubernetesClusterArgs {
	if x != nil {
		return x.Resource
	}
	return nil
}

type DeleteKubernetesClusterRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ResourceId string `protobuf:"bytes,1,opt,name=resource_id,json=resourceId,proto3" json:"resource_id,omitempty"`
}

func (x *DeleteKubernetesClusterRequest) Reset() {
	*x = DeleteKubernetesClusterRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_resources_kubernetes_cluster_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteKubernetesClusterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteKubernetesClusterRequest) ProtoMessage() {}

func (x *DeleteKubernetesClusterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_resources_kubernetes_cluster_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteKubernetesClusterRequest.ProtoReflect.Descriptor instead.
func (*DeleteKubernetesClusterRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_resources_kubernetes_cluster_proto_rawDescGZIP(), []int{3}
}

func (x *DeleteKubernetesClusterRequest) GetResourceId() string {
	if x != nil {
		return x.ResourceId
	}
	return ""
}

type KubernetesClusterArgs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CommonParameters *common.ResourceCommonArgs `protobuf:"bytes,1,opt,name=common_parameters,json=commonParameters,proto3" json:"common_parameters,omitempty"`
	Name             string                     `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	SubnetIds        []string                   `protobuf:"bytes,3,rep,name=subnet_ids,json=subnetIds,proto3" json:"subnet_ids,omitempty"`
}

func (x *KubernetesClusterArgs) Reset() {
	*x = KubernetesClusterArgs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_resources_kubernetes_cluster_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KubernetesClusterArgs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KubernetesClusterArgs) ProtoMessage() {}

func (x *KubernetesClusterArgs) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_resources_kubernetes_cluster_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KubernetesClusterArgs.ProtoReflect.Descriptor instead.
func (*KubernetesClusterArgs) Descriptor() ([]byte, []int) {
	return file_api_proto_resources_kubernetes_cluster_proto_rawDescGZIP(), []int{4}
}

func (x *KubernetesClusterArgs) GetCommonParameters() *common.ResourceCommonArgs {
	if x != nil {
		return x.CommonParameters
	}
	return nil
}

func (x *KubernetesClusterArgs) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *KubernetesClusterArgs) GetSubnetIds() []string {
	if x != nil {
		return x.SubnetIds
	}
	return nil
}

type KubernetesClusterResource struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CommonParameters *common.CommonResourceParameters `protobuf:"bytes,1,opt,name=common_parameters,json=commonParameters,proto3" json:"common_parameters,omitempty"`
	Name             string                           `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	SubnetIds        []string                         `protobuf:"bytes,3,rep,name=subnet_ids,json=subnetIds,proto3" json:"subnet_ids,omitempty"`
	// outputs
	Endpoint string `protobuf:"bytes,4,opt,name=endpoint,proto3" json:"endpoint,omitempty"`
}

func (x *KubernetesClusterResource) Reset() {
	*x = KubernetesClusterResource{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_resources_kubernetes_cluster_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KubernetesClusterResource) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KubernetesClusterResource) ProtoMessage() {}

func (x *KubernetesClusterResource) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_resources_kubernetes_cluster_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KubernetesClusterResource.ProtoReflect.Descriptor instead.
func (*KubernetesClusterResource) Descriptor() ([]byte, []int) {
	return file_api_proto_resources_kubernetes_cluster_proto_rawDescGZIP(), []int{5}
}

func (x *KubernetesClusterResource) GetCommonParameters() *common.CommonResourceParameters {
	if x != nil {
		return x.CommonParameters
	}
	return nil
}

func (x *KubernetesClusterResource) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *KubernetesClusterResource) GetSubnetIds() []string {
	if x != nil {
		return x.SubnetIds
	}
	return nil
}

func (x *KubernetesClusterResource) GetEndpoint() string {
	if x != nil {
		return x.Endpoint
	}
	return ""
}

var File_api_proto_resources_kubernetes_cluster_proto protoreflect.FileDescriptor

var file_api_proto_resources_kubernetes_cluster_proto_rawDesc = []byte{
	0x0a, 0x2c, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x6b, 0x75, 0x62, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x65, 0x73,
	0x5f, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x13,
	0x64, 0x65, 0x76, 0x2e, 0x6d, 0x75, 0x6c, 0x74, 0x79, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x73, 0x1a, 0x1d, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x68, 0x0a, 0x1e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4b, 0x75, 0x62, 0x65,
	0x72, 0x6e, 0x65, 0x74, 0x65, 0x73, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x46, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2a, 0x2e, 0x64, 0x65, 0x76, 0x2e, 0x6d, 0x75, 0x6c,
	0x74, 0x79, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x4b, 0x75, 0x62,
	0x65, 0x72, 0x6e, 0x65, 0x74, 0x65, 0x73, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x41, 0x72,
	0x67, 0x73, 0x52, 0x08, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x22, 0x3f, 0x0a, 0x1c,
	0x52, 0x65, 0x61, 0x64, 0x4b, 0x75, 0x62, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x65, 0x73, 0x43, 0x6c,
	0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0a, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64, 0x22, 0x89, 0x01,
	0x0a, 0x1e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4b, 0x75, 0x62, 0x65, 0x72, 0x6e, 0x65, 0x74,
	0x65, 0x73, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x1f, 0x0a, 0x0b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49,
	0x64, 0x12, 0x46, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x2a, 0x2e, 0x64, 0x65, 0x76, 0x2e, 0x6d, 0x75, 0x6c, 0x74, 0x79, 0x2e,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x4b, 0x75, 0x62, 0x65, 0x72, 0x6e,
	0x65, 0x74, 0x65, 0x73, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x41, 0x72, 0x67, 0x73, 0x52,
	0x08, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x22, 0x41, 0x0a, 0x1e, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x4b, 0x75, 0x62, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x65, 0x73, 0x43, 0x6c, 0x75,
	0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0a, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64, 0x22, 0x9d, 0x01, 0x0a,
	0x15, 0x4b, 0x75, 0x62, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x65, 0x73, 0x43, 0x6c, 0x75, 0x73, 0x74,
	0x65, 0x72, 0x41, 0x72, 0x67, 0x73, 0x12, 0x51, 0x0a, 0x11, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x5f, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x24, 0x2e, 0x64, 0x65, 0x76, 0x2e, 0x6d, 0x75, 0x6c, 0x74, 0x79, 0x2e, 0x63, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x43, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x41, 0x72, 0x67, 0x73, 0x52, 0x10, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x50,
	0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1d, 0x0a,
	0x0a, 0x73, 0x75, 0x62, 0x6e, 0x65, 0x74, 0x5f, 0x69, 0x64, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x09, 0x73, 0x75, 0x62, 0x6e, 0x65, 0x74, 0x49, 0x64, 0x73, 0x22, 0xc3, 0x01, 0x0a,
	0x19, 0x4b, 0x75, 0x62, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x65, 0x73, 0x43, 0x6c, 0x75, 0x73, 0x74,
	0x65, 0x72, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x57, 0x0a, 0x11, 0x63, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x5f, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2a, 0x2e, 0x64, 0x65, 0x76, 0x2e, 0x6d, 0x75, 0x6c, 0x74,
	0x79, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x52,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72,
	0x73, 0x52, 0x10, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74,
	0x65, 0x72, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x75, 0x62, 0x6e, 0x65,
	0x74, 0x5f, 0x69, 0x64, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x09, 0x73, 0x75, 0x62,
	0x6e, 0x65, 0x74, 0x49, 0x64, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69,
	0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69,
	0x6e, 0x74, 0x42, 0x5c, 0x0a, 0x17, 0x64, 0x65, 0x76, 0x2e, 0x6d, 0x75, 0x6c, 0x74, 0x79, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x42, 0x0e, 0x52,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a,
	0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x75, 0x6c, 0x74,
	0x79, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2f, 0x6d, 0x75, 0x6c, 0x74, 0x79, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_proto_resources_kubernetes_cluster_proto_rawDescOnce sync.Once
	file_api_proto_resources_kubernetes_cluster_proto_rawDescData = file_api_proto_resources_kubernetes_cluster_proto_rawDesc
)

func file_api_proto_resources_kubernetes_cluster_proto_rawDescGZIP() []byte {
	file_api_proto_resources_kubernetes_cluster_proto_rawDescOnce.Do(func() {
		file_api_proto_resources_kubernetes_cluster_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_proto_resources_kubernetes_cluster_proto_rawDescData)
	})
	return file_api_proto_resources_kubernetes_cluster_proto_rawDescData
}

var file_api_proto_resources_kubernetes_cluster_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_api_proto_resources_kubernetes_cluster_proto_goTypes = []interface{}{
	(*CreateKubernetesClusterRequest)(nil),  // 0: dev.multy.resources.CreateKubernetesClusterRequest
	(*ReadKubernetesClusterRequest)(nil),    // 1: dev.multy.resources.ReadKubernetesClusterRequest
	(*UpdateKubernetesClusterRequest)(nil),  // 2: dev.multy.resources.UpdateKubernetesClusterRequest
	(*DeleteKubernetesClusterRequest)(nil),  // 3: dev.multy.resources.DeleteKubernetesClusterRequest
	(*KubernetesClusterArgs)(nil),           // 4: dev.multy.resources.KubernetesClusterArgs
	(*KubernetesClusterResource)(nil),       // 5: dev.multy.resources.KubernetesClusterResource
	(*common.ResourceCommonArgs)(nil),       // 6: dev.multy.common.ResourceCommonArgs
	(*common.CommonResourceParameters)(nil), // 7: dev.multy.common.CommonResourceParameters
}
var file_api_proto_resources_kubernetes_cluster_proto_depIdxs = []int32{
	4, // 0: dev.multy.resources.CreateKubernetesClusterRequest.resource:type_name -> dev.multy.resources.KubernetesClusterArgs
	4, // 1: dev.multy.resources.UpdateKubernetesClusterRequest.resource:type_name -> dev.multy.resources.KubernetesClusterArgs
	6, // 2: dev.multy.resources.KubernetesClusterArgs.common_parameters:type_name -> dev.multy.common.ResourceCommonArgs
	7, // 3: dev.multy.resources.KubernetesClusterResource.common_parameters:type_name -> dev.multy.common.CommonResourceParameters
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_api_proto_resources_kubernetes_cluster_proto_init() }
func file_api_proto_resources_kubernetes_cluster_proto_init() {
	if File_api_proto_resources_kubernetes_cluster_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_proto_resources_kubernetes_cluster_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateKubernetesClusterRequest); i {
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
		file_api_proto_resources_kubernetes_cluster_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReadKubernetesClusterRequest); i {
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
		file_api_proto_resources_kubernetes_cluster_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateKubernetesClusterRequest); i {
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
		file_api_proto_resources_kubernetes_cluster_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteKubernetesClusterRequest); i {
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
		file_api_proto_resources_kubernetes_cluster_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KubernetesClusterArgs); i {
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
		file_api_proto_resources_kubernetes_cluster_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KubernetesClusterResource); i {
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
			RawDescriptor: file_api_proto_resources_kubernetes_cluster_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_proto_resources_kubernetes_cluster_proto_goTypes,
		DependencyIndexes: file_api_proto_resources_kubernetes_cluster_proto_depIdxs,
		MessageInfos:      file_api_proto_resources_kubernetes_cluster_proto_msgTypes,
	}.Build()
	File_api_proto_resources_kubernetes_cluster_proto = out.File
	file_api_proto_resources_kubernetes_cluster_proto_rawDesc = nil
	file_api_proto_resources_kubernetes_cluster_proto_goTypes = nil
	file_api_proto_resources_kubernetes_cluster_proto_depIdxs = nil
}
