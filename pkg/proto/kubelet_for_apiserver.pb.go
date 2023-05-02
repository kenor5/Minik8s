// grpc server is Kubelet and grpc client is Api Server

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.21.9
// source: proto/kubelet_for_apiserver.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_proto_kubelet_for_apiserver_proto protoreflect.FileDescriptor

var file_proto_kubelet_for_apiserver_proto_rawDesc = []byte{
	0x0a, 0x21, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6b, 0x75, 0x62, 0x65, 0x6c, 0x65, 0x74, 0x5f,
	0x66, 0x6f, 0x72, 0x5f, 0x61, 0x70, 0x69, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x15, 0x6b, 0x75, 0x62, 0x65, 0x6c, 0x65, 0x74, 0x5f, 0x66, 0x6f, 0x72,
	0x5f, 0x61, 0x70, 0x69, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x1a, 0x11, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x73, 0x68, 0x61, 0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0x8c, 0x01,
	0x0a, 0x17, 0x6b, 0x75, 0x62, 0x65, 0x6c, 0x65, 0x74, 0x41, 0x70, 0x69, 0x53, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x35, 0x0a, 0x08, 0x53, 0x61, 0x79,
	0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x12, 0x13, 0x2e, 0x73, 0x68, 0x61, 0x72, 0x65, 0x2e, 0x48, 0x65,
	0x6c, 0x6c, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x73, 0x68, 0x61,
	0x72, 0x65, 0x2e, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x3a, 0x0a, 0x09, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x50, 0x6f, 0x64, 0x12, 0x16, 0x2e,
	0x73, 0x68, 0x61, 0x72, 0x65, 0x2e, 0x41, 0x70, 0x70, 0x6c, 0x79, 0x50, 0x6f, 0x64, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x73, 0x68, 0x61, 0x72, 0x65, 0x2e, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x13, 0x5a, 0x11,
	0x6d, 0x69, 0x6e, 0x69, 0x6b, 0x38, 0x73, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_proto_kubelet_for_apiserver_proto_goTypes = []interface{}{
	(*HelloRequest)(nil),    // 0: share.HelloRequest
	(*ApplyPodRequest)(nil), // 1: share.ApplyPodRequest
	(*HelloResponse)(nil),   // 2: share.HelloResponse
	(*StatusResponse)(nil),  // 3: share.StatusResponse
}
var file_proto_kubelet_for_apiserver_proto_depIdxs = []int32{
	0, // 0: kubelet_for_apiserver.kubeletApiServerService.SayHello:input_type -> share.HelloRequest
	1, // 1: kubelet_for_apiserver.kubeletApiServerService.CreatePod:input_type -> share.ApplyPodRequest
	2, // 2: kubelet_for_apiserver.kubeletApiServerService.SayHello:output_type -> share.HelloResponse
	3, // 3: kubelet_for_apiserver.kubeletApiServerService.CreatePod:output_type -> share.StatusResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_kubelet_for_apiserver_proto_init() }
func file_proto_kubelet_for_apiserver_proto_init() {
	if File_proto_kubelet_for_apiserver_proto != nil {
		return
	}
	file_proto_share_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_kubelet_for_apiserver_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_kubelet_for_apiserver_proto_goTypes,
		DependencyIndexes: file_proto_kubelet_for_apiserver_proto_depIdxs,
	}.Build()
	File_proto_kubelet_for_apiserver_proto = out.File
	file_proto_kubelet_for_apiserver_proto_rawDesc = nil
	file_proto_kubelet_for_apiserver_proto_goTypes = nil
	file_proto_kubelet_for_apiserver_proto_depIdxs = nil
}
