// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.6.1
// source: api/proto/resourcespb/kubernetes_node_pool.proto

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

type CreateKubernetesNodePoolRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Resource *KubernetesNodePoolArgs `protobuf:"bytes,1,opt,name=resource,proto3" json:"resource,omitempty"`
}

func (x *CreateKubernetesNodePoolRequest) Reset() {
	*x = CreateKubernetesNodePoolRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_resourcespb_kubernetes_node_pool_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateKubernetesNodePoolRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateKubernetesNodePoolRequest) ProtoMessage() {}

func (x *CreateKubernetesNodePoolRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_resourcespb_kubernetes_node_pool_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateKubernetesNodePoolRequest.ProtoReflect.Descriptor instead.
func (*CreateKubernetesNodePoolRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_resourcespb_kubernetes_node_pool_proto_rawDescGZIP(), []int{0}
}

func (x *CreateKubernetesNodePoolRequest) GetResource() *KubernetesNodePoolArgs {
	if x != nil {
		return x.Resource
	}
	return nil
}

type ReadKubernetesNodePoolRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ResourceId string `protobuf:"bytes,1,opt,name=resource_id,json=resourceId,proto3" json:"resource_id,omitempty"`
}

func (x *ReadKubernetesNodePoolRequest) Reset() {
	*x = ReadKubernetesNodePoolRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_resourcespb_kubernetes_node_pool_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadKubernetesNodePoolRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadKubernetesNodePoolRequest) ProtoMessage() {}

func (x *ReadKubernetesNodePoolRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_resourcespb_kubernetes_node_pool_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadKubernetesNodePoolRequest.ProtoReflect.Descriptor instead.
func (*ReadKubernetesNodePoolRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_resourcespb_kubernetes_node_pool_proto_rawDescGZIP(), []int{1}
}

func (x *ReadKubernetesNodePoolRequest) GetResourceId() string {
	if x != nil {
		return x.ResourceId
	}
	return ""
}

type UpdateKubernetesNodePoolRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ResourceId string                  `protobuf:"bytes,1,opt,name=resource_id,json=resourceId,proto3" json:"resource_id,omitempty"`
	Resource   *KubernetesNodePoolArgs `protobuf:"bytes,2,opt,name=resource,proto3" json:"resource,omitempty"`
}

func (x *UpdateKubernetesNodePoolRequest) Reset() {
	*x = UpdateKubernetesNodePoolRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_resourcespb_kubernetes_node_pool_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateKubernetesNodePoolRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateKubernetesNodePoolRequest) ProtoMessage() {}

func (x *UpdateKubernetesNodePoolRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_resourcespb_kubernetes_node_pool_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateKubernetesNodePoolRequest.ProtoReflect.Descriptor instead.
func (*UpdateKubernetesNodePoolRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_resourcespb_kubernetes_node_pool_proto_rawDescGZIP(), []int{2}
}

func (x *UpdateKubernetesNodePoolRequest) GetResourceId() string {
	if x != nil {
		return x.ResourceId
	}
	return ""
}

func (x *UpdateKubernetesNodePoolRequest) GetResource() *KubernetesNodePoolArgs {
	if x != nil {
		return x.Resource
	}
	return nil
}

type DeleteKubernetesNodePoolRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ResourceId string `protobuf:"bytes,1,opt,name=resource_id,json=resourceId,proto3" json:"resource_id,omitempty"`
}

func (x *DeleteKubernetesNodePoolRequest) Reset() {
	*x = DeleteKubernetesNodePoolRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_resourcespb_kubernetes_node_pool_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteKubernetesNodePoolRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteKubernetesNodePoolRequest) ProtoMessage() {}

func (x *DeleteKubernetesNodePoolRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_resourcespb_kubernetes_node_pool_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteKubernetesNodePoolRequest.ProtoReflect.Descriptor instead.
func (*DeleteKubernetesNodePoolRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_resourcespb_kubernetes_node_pool_proto_rawDescGZIP(), []int{3}
}

func (x *DeleteKubernetesNodePoolRequest) GetResourceId() string {
	if x != nil {
		return x.ResourceId
	}
	return ""
}

type KubernetesNodePoolArgs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CommonParameters  *commonpb.ChildResourceCommonArgs `protobuf:"bytes,1,opt,name=common_parameters,json=commonParameters,proto3" json:"common_parameters,omitempty"`
	Name              string                            `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	SubnetIds         []string                          `protobuf:"bytes,3,rep,name=subnet_ids,json=subnetIds,proto3" json:"subnet_ids,omitempty"`
	ClusterId         string                            `protobuf:"bytes,4,opt,name=cluster_id,json=clusterId,proto3" json:"cluster_id,omitempty"`
	IsDefaultPool     bool                              `protobuf:"varint,5,opt,name=is_default_pool,json=isDefaultPool,proto3" json:"is_default_pool,omitempty"`
	StartingNodeCount int32                             `protobuf:"varint,6,opt,name=starting_node_count,json=startingNodeCount,proto3" json:"starting_node_count,omitempty"`
	MinNodeCount      int32                             `protobuf:"varint,7,opt,name=min_node_count,json=minNodeCount,proto3" json:"min_node_count,omitempty"`
	MaxNodeCount      int32                             `protobuf:"varint,8,opt,name=max_node_count,json=maxNodeCount,proto3" json:"max_node_count,omitempty"`
	VmSize            commonpb.VmSize_Enum              `protobuf:"varint,9,opt,name=vm_size,json=vmSize,proto3,enum=dev.multy.common.VmSize_Enum" json:"vm_size,omitempty"`
	DiskSizeGb        int64                             `protobuf:"varint,10,opt,name=disk_size_gb,json=diskSizeGb,proto3" json:"disk_size_gb,omitempty"`
	Labels            map[string]string                 `protobuf:"bytes,11,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *KubernetesNodePoolArgs) Reset() {
	*x = KubernetesNodePoolArgs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_resourcespb_kubernetes_node_pool_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KubernetesNodePoolArgs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KubernetesNodePoolArgs) ProtoMessage() {}

func (x *KubernetesNodePoolArgs) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_resourcespb_kubernetes_node_pool_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KubernetesNodePoolArgs.ProtoReflect.Descriptor instead.
func (*KubernetesNodePoolArgs) Descriptor() ([]byte, []int) {
	return file_api_proto_resourcespb_kubernetes_node_pool_proto_rawDescGZIP(), []int{4}
}

func (x *KubernetesNodePoolArgs) GetCommonParameters() *commonpb.ChildResourceCommonArgs {
	if x != nil {
		return x.CommonParameters
	}
	return nil
}

func (x *KubernetesNodePoolArgs) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *KubernetesNodePoolArgs) GetSubnetIds() []string {
	if x != nil {
		return x.SubnetIds
	}
	return nil
}

func (x *KubernetesNodePoolArgs) GetClusterId() string {
	if x != nil {
		return x.ClusterId
	}
	return ""
}

func (x *KubernetesNodePoolArgs) GetIsDefaultPool() bool {
	if x != nil {
		return x.IsDefaultPool
	}
	return false
}

func (x *KubernetesNodePoolArgs) GetStartingNodeCount() int32 {
	if x != nil {
		return x.StartingNodeCount
	}
	return 0
}

func (x *KubernetesNodePoolArgs) GetMinNodeCount() int32 {
	if x != nil {
		return x.MinNodeCount
	}
	return 0
}

func (x *KubernetesNodePoolArgs) GetMaxNodeCount() int32 {
	if x != nil {
		return x.MaxNodeCount
	}
	return 0
}

func (x *KubernetesNodePoolArgs) GetVmSize() commonpb.VmSize_Enum {
	if x != nil {
		return x.VmSize
	}
	return commonpb.VmSize_UNKNOWN_VM_SIZE
}

func (x *KubernetesNodePoolArgs) GetDiskSizeGb() int64 {
	if x != nil {
		return x.DiskSizeGb
	}
	return 0
}

func (x *KubernetesNodePoolArgs) GetLabels() map[string]string {
	if x != nil {
		return x.Labels
	}
	return nil
}

type KubernetesNodePoolResource struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CommonParameters  *commonpb.CommonChildResourceParameters `protobuf:"bytes,1,opt,name=common_parameters,json=commonParameters,proto3" json:"common_parameters,omitempty"`
	Name              string                                  `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	SubnetIds         []string                                `protobuf:"bytes,3,rep,name=subnet_ids,json=subnetIds,proto3" json:"subnet_ids,omitempty"`
	ClusterId         string                                  `protobuf:"bytes,4,opt,name=cluster_id,json=clusterId,proto3" json:"cluster_id,omitempty"`
	IsDefaultPool     bool                                    `protobuf:"varint,5,opt,name=is_default_pool,json=isDefaultPool,proto3" json:"is_default_pool,omitempty"`
	StartingNodeCount int32                                   `protobuf:"varint,6,opt,name=starting_node_count,json=startingNodeCount,proto3" json:"starting_node_count,omitempty"`
	MinNodeCount      int32                                   `protobuf:"varint,7,opt,name=min_node_count,json=minNodeCount,proto3" json:"min_node_count,omitempty"`
	MaxNodeCount      int32                                   `protobuf:"varint,8,opt,name=max_node_count,json=maxNodeCount,proto3" json:"max_node_count,omitempty"`
	VmSize            commonpb.VmSize_Enum                    `protobuf:"varint,9,opt,name=vm_size,json=vmSize,proto3,enum=dev.multy.common.VmSize_Enum" json:"vm_size,omitempty"`
	DiskSizeGb        int64                                   `protobuf:"varint,10,opt,name=disk_size_gb,json=diskSizeGb,proto3" json:"disk_size_gb,omitempty"`
	Labels            map[string]string                       `protobuf:"bytes,11,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *KubernetesNodePoolResource) Reset() {
	*x = KubernetesNodePoolResource{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_resourcespb_kubernetes_node_pool_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KubernetesNodePoolResource) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KubernetesNodePoolResource) ProtoMessage() {}

func (x *KubernetesNodePoolResource) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_resourcespb_kubernetes_node_pool_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KubernetesNodePoolResource.ProtoReflect.Descriptor instead.
func (*KubernetesNodePoolResource) Descriptor() ([]byte, []int) {
	return file_api_proto_resourcespb_kubernetes_node_pool_proto_rawDescGZIP(), []int{5}
}

func (x *KubernetesNodePoolResource) GetCommonParameters() *commonpb.CommonChildResourceParameters {
	if x != nil {
		return x.CommonParameters
	}
	return nil
}

func (x *KubernetesNodePoolResource) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *KubernetesNodePoolResource) GetSubnetIds() []string {
	if x != nil {
		return x.SubnetIds
	}
	return nil
}

func (x *KubernetesNodePoolResource) GetClusterId() string {
	if x != nil {
		return x.ClusterId
	}
	return ""
}

func (x *KubernetesNodePoolResource) GetIsDefaultPool() bool {
	if x != nil {
		return x.IsDefaultPool
	}
	return false
}

func (x *KubernetesNodePoolResource) GetStartingNodeCount() int32 {
	if x != nil {
		return x.StartingNodeCount
	}
	return 0
}

func (x *KubernetesNodePoolResource) GetMinNodeCount() int32 {
	if x != nil {
		return x.MinNodeCount
	}
	return 0
}

func (x *KubernetesNodePoolResource) GetMaxNodeCount() int32 {
	if x != nil {
		return x.MaxNodeCount
	}
	return 0
}

func (x *KubernetesNodePoolResource) GetVmSize() commonpb.VmSize_Enum {
	if x != nil {
		return x.VmSize
	}
	return commonpb.VmSize_UNKNOWN_VM_SIZE
}

func (x *KubernetesNodePoolResource) GetDiskSizeGb() int64 {
	if x != nil {
		return x.DiskSizeGb
	}
	return 0
}

func (x *KubernetesNodePoolResource) GetLabels() map[string]string {
	if x != nil {
		return x.Labels
	}
	return nil
}

var File_api_proto_resourcespb_kubernetes_node_pool_proto protoreflect.FileDescriptor

var file_api_proto_resourcespb_kubernetes_node_pool_proto_rawDesc = []byte{
	0x0a, 0x30, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x73, 0x70, 0x62, 0x2f, 0x6b, 0x75, 0x62, 0x65, 0x72, 0x6e, 0x65, 0x74,
	0x65, 0x73, 0x5f, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x70, 0x6f, 0x6f, 0x6c, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x13, 0x64, 0x65, 0x76, 0x2e, 0x6d, 0x75, 0x6c, 0x74, 0x79, 0x2e, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x1a, 0x1f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x70, 0x62, 0x2f, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x6a, 0x0a, 0x1f, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x4b, 0x75, 0x62, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x65, 0x73, 0x4e, 0x6f, 0x64, 0x65,
	0x50, 0x6f, 0x6f, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x47, 0x0a, 0x08, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2b, 0x2e,
	0x64, 0x65, 0x76, 0x2e, 0x6d, 0x75, 0x6c, 0x74, 0x79, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x73, 0x2e, 0x4b, 0x75, 0x62, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x65, 0x73, 0x4e, 0x6f,
	0x64, 0x65, 0x50, 0x6f, 0x6f, 0x6c, 0x41, 0x72, 0x67, 0x73, 0x52, 0x08, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x22, 0x40, 0x0a, 0x1d, 0x52, 0x65, 0x61, 0x64, 0x4b, 0x75, 0x62, 0x65,
	0x72, 0x6e, 0x65, 0x74, 0x65, 0x73, 0x4e, 0x6f, 0x64, 0x65, 0x50, 0x6f, 0x6f, 0x6c, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x49, 0x64, 0x22, 0x8b, 0x01, 0x0a, 0x1f, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x4b, 0x75, 0x62, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x65, 0x73, 0x4e, 0x6f, 0x64, 0x65, 0x50,
	0x6f, 0x6f, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0a, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64, 0x12, 0x47, 0x0a, 0x08, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2b, 0x2e,
	0x64, 0x65, 0x76, 0x2e, 0x6d, 0x75, 0x6c, 0x74, 0x79, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x73, 0x2e, 0x4b, 0x75, 0x62, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x65, 0x73, 0x4e, 0x6f,
	0x64, 0x65, 0x50, 0x6f, 0x6f, 0x6c, 0x41, 0x72, 0x67, 0x73, 0x52, 0x08, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x22, 0x42, 0x0a, 0x1f, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4b, 0x75,
	0x62, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x65, 0x73, 0x4e, 0x6f, 0x64, 0x65, 0x50, 0x6f, 0x6f, 0x6c,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64, 0x22, 0xcc, 0x04, 0x0a, 0x16, 0x4b, 0x75, 0x62,
	0x65, 0x72, 0x6e, 0x65, 0x74, 0x65, 0x73, 0x4e, 0x6f, 0x64, 0x65, 0x50, 0x6f, 0x6f, 0x6c, 0x41,
	0x72, 0x67, 0x73, 0x12, 0x56, 0x0a, 0x11, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5f, 0x70, 0x61,
	0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x29,
	0x2e, 0x64, 0x65, 0x76, 0x2e, 0x6d, 0x75, 0x6c, 0x74, 0x79, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2e, 0x43, 0x68, 0x69, 0x6c, 0x64, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x43,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x41, 0x72, 0x67, 0x73, 0x52, 0x10, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x1d, 0x0a, 0x0a, 0x73, 0x75, 0x62, 0x6e, 0x65, 0x74, 0x5f, 0x69, 0x64, 0x73, 0x18, 0x03, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x09, 0x73, 0x75, 0x62, 0x6e, 0x65, 0x74, 0x49, 0x64, 0x73, 0x12, 0x1d,
	0x0a, 0x0a, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x49, 0x64, 0x12, 0x26, 0x0a,
	0x0f, 0x69, 0x73, 0x5f, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x5f, 0x70, 0x6f, 0x6f, 0x6c,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0d, 0x69, 0x73, 0x44, 0x65, 0x66, 0x61, 0x75, 0x6c,
	0x74, 0x50, 0x6f, 0x6f, 0x6c, 0x12, 0x2e, 0x0a, 0x13, 0x73, 0x74, 0x61, 0x72, 0x74, 0x69, 0x6e,
	0x67, 0x5f, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x11, 0x73, 0x74, 0x61, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x4e, 0x6f, 0x64, 0x65,
	0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x24, 0x0a, 0x0e, 0x6d, 0x69, 0x6e, 0x5f, 0x6e, 0x6f, 0x64,
	0x65, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0c, 0x6d,
	0x69, 0x6e, 0x4e, 0x6f, 0x64, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x24, 0x0a, 0x0e, 0x6d,
	0x61, 0x78, 0x5f, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x08, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x0c, 0x6d, 0x61, 0x78, 0x4e, 0x6f, 0x64, 0x65, 0x43, 0x6f, 0x75, 0x6e,
	0x74, 0x12, 0x36, 0x0a, 0x07, 0x76, 0x6d, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x09, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x1d, 0x2e, 0x64, 0x65, 0x76, 0x2e, 0x6d, 0x75, 0x6c, 0x74, 0x79, 0x2e, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x56, 0x6d, 0x53, 0x69, 0x7a, 0x65, 0x2e, 0x45, 0x6e, 0x75,
	0x6d, 0x52, 0x06, 0x76, 0x6d, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x20, 0x0a, 0x0c, 0x64, 0x69, 0x73,
	0x6b, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x5f, 0x67, 0x62, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x0a, 0x64, 0x69, 0x73, 0x6b, 0x53, 0x69, 0x7a, 0x65, 0x47, 0x62, 0x12, 0x4f, 0x0a, 0x06, 0x6c,
	0x61, 0x62, 0x65, 0x6c, 0x73, 0x18, 0x0b, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x37, 0x2e, 0x64, 0x65,
	0x76, 0x2e, 0x6d, 0x75, 0x6c, 0x74, 0x79, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x73, 0x2e, 0x4b, 0x75, 0x62, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x65, 0x73, 0x4e, 0x6f, 0x64, 0x65,
	0x50, 0x6f, 0x6f, 0x6c, 0x41, 0x72, 0x67, 0x73, 0x2e, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x1a, 0x39, 0x0a, 0x0b,
	0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0xda, 0x04, 0x0a, 0x1a, 0x4b, 0x75, 0x62, 0x65,
	0x72, 0x6e, 0x65, 0x74, 0x65, 0x73, 0x4e, 0x6f, 0x64, 0x65, 0x50, 0x6f, 0x6f, 0x6c, 0x52, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x5c, 0x0a, 0x11, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x5f, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x2f, 0x2e, 0x64, 0x65, 0x76, 0x2e, 0x6d, 0x75, 0x6c, 0x74, 0x79, 0x2e, 0x63, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x43, 0x68, 0x69, 0x6c, 0x64,
	0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65,
	0x72, 0x73, 0x52, 0x10, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65,
	0x74, 0x65, 0x72, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x75, 0x62, 0x6e,
	0x65, 0x74, 0x5f, 0x69, 0x64, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x09, 0x73, 0x75,
	0x62, 0x6e, 0x65, 0x74, 0x49, 0x64, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x6c, 0x75, 0x73, 0x74,
	0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x6c, 0x75,
	0x73, 0x74, 0x65, 0x72, 0x49, 0x64, 0x12, 0x26, 0x0a, 0x0f, 0x69, 0x73, 0x5f, 0x64, 0x65, 0x66,
	0x61, 0x75, 0x6c, 0x74, 0x5f, 0x70, 0x6f, 0x6f, 0x6c, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x0d, 0x69, 0x73, 0x44, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x50, 0x6f, 0x6f, 0x6c, 0x12, 0x2e,
	0x0a, 0x13, 0x73, 0x74, 0x61, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x5f, 0x6e, 0x6f, 0x64, 0x65, 0x5f,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x11, 0x73, 0x74, 0x61,
	0x72, 0x74, 0x69, 0x6e, 0x67, 0x4e, 0x6f, 0x64, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x24,
	0x0a, 0x0e, 0x6d, 0x69, 0x6e, 0x5f, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0c, 0x6d, 0x69, 0x6e, 0x4e, 0x6f, 0x64, 0x65, 0x43,
	0x6f, 0x75, 0x6e, 0x74, 0x12, 0x24, 0x0a, 0x0e, 0x6d, 0x61, 0x78, 0x5f, 0x6e, 0x6f, 0x64, 0x65,
	0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0c, 0x6d, 0x61,
	0x78, 0x4e, 0x6f, 0x64, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x36, 0x0a, 0x07, 0x76, 0x6d,
	0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1d, 0x2e, 0x64, 0x65,
	0x76, 0x2e, 0x6d, 0x75, 0x6c, 0x74, 0x79, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x56,
	0x6d, 0x53, 0x69, 0x7a, 0x65, 0x2e, 0x45, 0x6e, 0x75, 0x6d, 0x52, 0x06, 0x76, 0x6d, 0x53, 0x69,
	0x7a, 0x65, 0x12, 0x20, 0x0a, 0x0c, 0x64, 0x69, 0x73, 0x6b, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x5f,
	0x67, 0x62, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x64, 0x69, 0x73, 0x6b, 0x53, 0x69,
	0x7a, 0x65, 0x47, 0x62, 0x12, 0x53, 0x0a, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x18, 0x0b,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x3b, 0x2e, 0x64, 0x65, 0x76, 0x2e, 0x6d, 0x75, 0x6c, 0x74, 0x79,
	0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x4b, 0x75, 0x62, 0x65, 0x72,
	0x6e, 0x65, 0x74, 0x65, 0x73, 0x4e, 0x6f, 0x64, 0x65, 0x50, 0x6f, 0x6f, 0x6c, 0x52, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x52, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x1a, 0x39, 0x0a, 0x0b, 0x4c, 0x61, 0x62,
	0x65, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x3a, 0x02, 0x38, 0x01, 0x42, 0x5e, 0x0a, 0x17, 0x64, 0x65, 0x76, 0x2e, 0x6d, 0x75, 0x6c, 0x74,
	0x79, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x42,
	0x0e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50,
	0x01, 0x5a, 0x31, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x75,
	0x6c, 0x74, 0x79, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2f, 0x6d, 0x75, 0x6c, 0x74, 0x79, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x73, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_proto_resourcespb_kubernetes_node_pool_proto_rawDescOnce sync.Once
	file_api_proto_resourcespb_kubernetes_node_pool_proto_rawDescData = file_api_proto_resourcespb_kubernetes_node_pool_proto_rawDesc
)

func file_api_proto_resourcespb_kubernetes_node_pool_proto_rawDescGZIP() []byte {
	file_api_proto_resourcespb_kubernetes_node_pool_proto_rawDescOnce.Do(func() {
		file_api_proto_resourcespb_kubernetes_node_pool_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_proto_resourcespb_kubernetes_node_pool_proto_rawDescData)
	})
	return file_api_proto_resourcespb_kubernetes_node_pool_proto_rawDescData
}

var file_api_proto_resourcespb_kubernetes_node_pool_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_api_proto_resourcespb_kubernetes_node_pool_proto_goTypes = []interface{}{
	(*CreateKubernetesNodePoolRequest)(nil),        // 0: dev.multy.resources.CreateKubernetesNodePoolRequest
	(*ReadKubernetesNodePoolRequest)(nil),          // 1: dev.multy.resources.ReadKubernetesNodePoolRequest
	(*UpdateKubernetesNodePoolRequest)(nil),        // 2: dev.multy.resources.UpdateKubernetesNodePoolRequest
	(*DeleteKubernetesNodePoolRequest)(nil),        // 3: dev.multy.resources.DeleteKubernetesNodePoolRequest
	(*KubernetesNodePoolArgs)(nil),                 // 4: dev.multy.resources.KubernetesNodePoolArgs
	(*KubernetesNodePoolResource)(nil),             // 5: dev.multy.resources.KubernetesNodePoolResource
	nil,                                            // 6: dev.multy.resources.KubernetesNodePoolArgs.LabelsEntry
	nil,                                            // 7: dev.multy.resources.KubernetesNodePoolResource.LabelsEntry
	(*commonpb.ChildResourceCommonArgs)(nil),       // 8: dev.multy.common.ChildResourceCommonArgs
	(commonpb.VmSize_Enum)(0),                      // 9: dev.multy.common.VmSize.enumValues
	(*commonpb.CommonChildResourceParameters)(nil), // 10: dev.multy.common.CommonChildResourceParameters
}
var file_api_proto_resourcespb_kubernetes_node_pool_proto_depIdxs = []int32{
	4,  // 0: dev.multy.resources.CreateKubernetesNodePoolRequest.resource:type_name -> dev.multy.resources.KubernetesNodePoolArgs
	4,  // 1: dev.multy.resources.UpdateKubernetesNodePoolRequest.resource:type_name -> dev.multy.resources.KubernetesNodePoolArgs
	8,  // 2: dev.multy.resources.KubernetesNodePoolArgs.common_parameters:type_name -> dev.multy.common.ChildResourceCommonArgs
	9,  // 3: dev.multy.resources.KubernetesNodePoolArgs.vm_size:type_name -> dev.multy.common.VmSize.enumValues
	6,  // 4: dev.multy.resources.KubernetesNodePoolArgs.labels:type_name -> dev.multy.resources.KubernetesNodePoolArgs.LabelsEntry
	10, // 5: dev.multy.resources.KubernetesNodePoolResource.common_parameters:type_name -> dev.multy.common.CommonChildResourceParameters
	9,  // 6: dev.multy.resources.KubernetesNodePoolResource.vm_size:type_name -> dev.multy.common.VmSize.enumValues
	7,  // 7: dev.multy.resources.KubernetesNodePoolResource.labels:type_name -> dev.multy.resources.KubernetesNodePoolResource.LabelsEntry
	8,  // [8:8] is the sub-list for method output_type
	8,  // [8:8] is the sub-list for method input_type
	8,  // [8:8] is the sub-list for extension type_name
	8,  // [8:8] is the sub-list for extension extendee
	0,  // [0:8] is the sub-list for field type_name
}

func init() { file_api_proto_resourcespb_kubernetes_node_pool_proto_init() }
func file_api_proto_resourcespb_kubernetes_node_pool_proto_init() {
	if File_api_proto_resourcespb_kubernetes_node_pool_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_proto_resourcespb_kubernetes_node_pool_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateKubernetesNodePoolRequest); i {
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
		file_api_proto_resourcespb_kubernetes_node_pool_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReadKubernetesNodePoolRequest); i {
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
		file_api_proto_resourcespb_kubernetes_node_pool_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateKubernetesNodePoolRequest); i {
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
		file_api_proto_resourcespb_kubernetes_node_pool_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteKubernetesNodePoolRequest); i {
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
		file_api_proto_resourcespb_kubernetes_node_pool_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KubernetesNodePoolArgs); i {
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
		file_api_proto_resourcespb_kubernetes_node_pool_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KubernetesNodePoolResource); i {
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
			RawDescriptor: file_api_proto_resourcespb_kubernetes_node_pool_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_proto_resourcespb_kubernetes_node_pool_proto_goTypes,
		DependencyIndexes: file_api_proto_resourcespb_kubernetes_node_pool_proto_depIdxs,
		MessageInfos:      file_api_proto_resourcespb_kubernetes_node_pool_proto_msgTypes,
	}.Build()
	File_api_proto_resourcespb_kubernetes_node_pool_proto = out.File
	file_api_proto_resourcespb_kubernetes_node_pool_proto_rawDesc = nil
	file_api_proto_resourcespb_kubernetes_node_pool_proto_goTypes = nil
	file_api_proto_resourcespb_kubernetes_node_pool_proto_depIdxs = nil
}
