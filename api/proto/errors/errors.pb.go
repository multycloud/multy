// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.6.1
// source: api/proto/errors/errors.proto

package errors

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

type ResourceValidationError struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ResourceId   string `protobuf:"bytes,1,opt,name=resource_id,json=resourceId,proto3" json:"resource_id,omitempty"`
	ErrorMessage string `protobuf:"bytes,2,opt,name=error_message,json=errorMessage,proto3" json:"error_message,omitempty"`
	// this is tentative, it might not exist or might have a different name from the proto request
	FieldName string `protobuf:"bytes,3,opt,name=field_name,json=fieldName,proto3" json:"field_name,omitempty"`
}

func (x *ResourceValidationError) Reset() {
	*x = ResourceValidationError{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_errors_errors_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResourceValidationError) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResourceValidationError) ProtoMessage() {}

func (x *ResourceValidationError) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_errors_errors_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResourceValidationError.ProtoReflect.Descriptor instead.
func (*ResourceValidationError) Descriptor() ([]byte, []int) {
	return file_api_proto_errors_errors_proto_rawDescGZIP(), []int{0}
}

func (x *ResourceValidationError) GetResourceId() string {
	if x != nil {
		return x.ResourceId
	}
	return ""
}

func (x *ResourceValidationError) GetErrorMessage() string {
	if x != nil {
		return x.ErrorMessage
	}
	return ""
}

func (x *ResourceValidationError) GetFieldName() string {
	if x != nil {
		return x.FieldName
	}
	return ""
}

var File_api_proto_errors_errors_proto protoreflect.FileDescriptor

var file_api_proto_errors_errors_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x73, 0x2f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x10, 0x64, 0x65, 0x76, 0x2e, 0x6d, 0x75, 0x6c, 0x74, 0x79, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x22, 0x7e, 0x0a, 0x17, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x56, 0x61, 0x6c,
	0x69, 0x64, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x1f, 0x0a, 0x0b,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0a, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64, 0x12, 0x23, 0x0a,
	0x0d, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x4e, 0x61, 0x6d,
	0x65, 0x42, 0x52, 0x0a, 0x14, 0x64, 0x65, 0x76, 0x2e, 0x6d, 0x75, 0x6c, 0x74, 0x79, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x42, 0x0a, 0x4d, 0x75, 0x6c, 0x74, 0x79,
	0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x2c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x75, 0x6c, 0x74, 0x79, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2f, 0x6d,
	0x75, 0x6c, 0x74, 0x79, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65,
	0x72, 0x72, 0x6f, 0x72, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_proto_errors_errors_proto_rawDescOnce sync.Once
	file_api_proto_errors_errors_proto_rawDescData = file_api_proto_errors_errors_proto_rawDesc
)

func file_api_proto_errors_errors_proto_rawDescGZIP() []byte {
	file_api_proto_errors_errors_proto_rawDescOnce.Do(func() {
		file_api_proto_errors_errors_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_proto_errors_errors_proto_rawDescData)
	})
	return file_api_proto_errors_errors_proto_rawDescData
}

var file_api_proto_errors_errors_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_api_proto_errors_errors_proto_goTypes = []interface{}{
	(*ResourceValidationError)(nil), // 0: dev.multy.config.ResourceValidationError
}
var file_api_proto_errors_errors_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_api_proto_errors_errors_proto_init() }
func file_api_proto_errors_errors_proto_init() {
	if File_api_proto_errors_errors_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_proto_errors_errors_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResourceValidationError); i {
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
			RawDescriptor: file_api_proto_errors_errors_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_proto_errors_errors_proto_goTypes,
		DependencyIndexes: file_api_proto_errors_errors_proto_depIdxs,
		MessageInfos:      file_api_proto_errors_errors_proto_msgTypes,
	}.Build()
	File_api_proto_errors_errors_proto = out.File
	file_api_proto_errors_errors_proto_rawDesc = nil
	file_api_proto_errors_errors_proto_goTypes = nil
	file_api_proto_errors_errors_proto_depIdxs = nil
}
