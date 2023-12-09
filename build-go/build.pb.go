// SPDX-License-Identifier: MIT
// Copyright © 2022 The Tvix Authors

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: tvix/build/protos/build.proto

package buildv1

import (
	castore_go "code.tvl.fyi/tvix/castore-go"
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

// A BuildRequest describes the request of something to be run on the builder.
// It is distinct from an actual [Build] that has already happened, or might be
// currently ongoing.
//
// A BuildRequest can be seen as a more normalized version of a Derivation
// (parsed from A-Term), "writing out" some of the Nix-internal details about
// how e.g. environment variables in the build are set.
//
// Nix has some impurities when building a Derivation, for example the --cores option
// ends up as an environment variable in the build, that's not part of the ATerm.
//
// As of now, we serialize this into the BuildRequest, so builders can stay dumb.
// This might change in the future.
//
// There's also a big difference when it comes to how inputs are modelled:
//   - Nix only uses store path (strings) to describe the inputs.
//     As store paths can be input-addressed, a certain store path can contain
//     different contents (as not all store paths are binary reproducible).
//     This requires that for every input-addressed input, the builder has access
//     to either the input's deriver (and needs to build it) or else a trusted
//     source for the built input.
//     to upload input-addressed paths, requiring the trusted users concept.
//   - tvix-build records a list of tvix.castore.v1.Node as inputs.
//     These map from the store path base name to their contents, relieving the
//     builder from having to "trust" any input-addressed paths, contrary to Nix.
//
// While this approach gives a better hermeticity, it has one downside:
// A BuildRequest can only be sent once the contents of all its inputs are known.
//
// As of now, we're okay to accept this, but it prevents uploading an
// entirely-non-IFD subgraph of BuildRequests eagerly.
//
// FUTUREWORK: We might be introducing another way to refer to inputs, to
// support "send all BuildRequest for a nixpkgs eval to a remote builder and put
// the laptop to sleep" usecases later.
type BuildRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The command (and its args) executed as the build script.
	// In the case of a Nix derivation, this is usually
	// ["/path/to/some-bash/bin/bash", "-e", "/path/to/some/builder.sh"].
	CommandArgs []string `protobuf:"bytes,1,rep,name=command_args,json=commandArgs,proto3" json:"command_args,omitempty"`
	// The list of outputs the build is expected to produce.
	// These are basenames inside /nix/store.
	// If the path is not produced, the build is considered to have failed.
	// Outputs are sorted.
	Outputs []string `protobuf:"bytes,2,rep,name=outputs,proto3" json:"outputs,omitempty"`
	// The list of environment variables and their values that should be set
	// inside the build environment.
	// This includes both environment vars set inside the derivation, as well as
	// more "ephemeral" ones like NIX_BUILD_CORES, controlled by the `--cores`
	// CLI option of `nix-build`.
	// For now, we consume this as an option when turning a Derivation into a BuildRequest,
	// similar to how Nix has a `--cores` option.
	// We don't want to bleed these very nix-specific sandbox impl details into
	// (dumber) builders if we don't have to.
	// Environment variables are sorted by their keys.
	EnvironmentVars []*BuildRequest_EnvVar `protobuf:"bytes,3,rep,name=environment_vars,json=environmentVars,proto3" json:"environment_vars,omitempty"`
	// The list of all root nodes that should be visible in /nix/store at the
	// time of the build.
	// As root nodes are content-addressed, no additional signatures are needed
	// to substitute / make these available in the build environment.
	// Inputs are sorted by their names.
	Inputs []*castore_go.Node `protobuf:"bytes,4,rep,name=inputs,proto3" json:"inputs,omitempty"`
	// A set of constraints that need to be satisfied on a build host before a
	// Build can be started.
	Constraints *BuildRequest_BuildConstraints `protobuf:"bytes,5,opt,name=constraints,proto3" json:"constraints,omitempty"`
}

func (x *BuildRequest) Reset() {
	*x = BuildRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tvix_build_protos_build_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BuildRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BuildRequest) ProtoMessage() {}

func (x *BuildRequest) ProtoReflect() protoreflect.Message {
	mi := &file_tvix_build_protos_build_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BuildRequest.ProtoReflect.Descriptor instead.
func (*BuildRequest) Descriptor() ([]byte, []int) {
	return file_tvix_build_protos_build_proto_rawDescGZIP(), []int{0}
}

func (x *BuildRequest) GetCommandArgs() []string {
	if x != nil {
		return x.CommandArgs
	}
	return nil
}

func (x *BuildRequest) GetOutputs() []string {
	if x != nil {
		return x.Outputs
	}
	return nil
}

func (x *BuildRequest) GetEnvironmentVars() []*BuildRequest_EnvVar {
	if x != nil {
		return x.EnvironmentVars
	}
	return nil
}

func (x *BuildRequest) GetInputs() []*castore_go.Node {
	if x != nil {
		return x.Inputs
	}
	return nil
}

func (x *BuildRequest) GetConstraints() *BuildRequest_BuildConstraints {
	if x != nil {
		return x.Constraints
	}
	return nil
}

// A Build is (one possible) outcome of executing a [BuildRequest].
type Build struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The orginal build request producing the build.
	BuildRequest *BuildRequest `protobuf:"bytes,1,opt,name=build_request,json=buildRequest,proto3" json:"build_request,omitempty"` // <- TODO: define hashing scheme for BuildRequest, refer to it by hash?
	// The outputs that were produced after successfully building.
	// They are sorted by their names.
	Outputs []*castore_go.Node `protobuf:"bytes,2,rep,name=outputs,proto3" json:"outputs,omitempty"`
}

func (x *Build) Reset() {
	*x = Build{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tvix_build_protos_build_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Build) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Build) ProtoMessage() {}

func (x *Build) ProtoReflect() protoreflect.Message {
	mi := &file_tvix_build_protos_build_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Build.ProtoReflect.Descriptor instead.
func (*Build) Descriptor() ([]byte, []int) {
	return file_tvix_build_protos_build_proto_rawDescGZIP(), []int{1}
}

func (x *Build) GetBuildRequest() *BuildRequest {
	if x != nil {
		return x.BuildRequest
	}
	return nil
}

func (x *Build) GetOutputs() []*castore_go.Node {
	if x != nil {
		return x.Outputs
	}
	return nil
}

type BuildRequest_EnvVar struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value []byte `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *BuildRequest_EnvVar) Reset() {
	*x = BuildRequest_EnvVar{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tvix_build_protos_build_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BuildRequest_EnvVar) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BuildRequest_EnvVar) ProtoMessage() {}

func (x *BuildRequest_EnvVar) ProtoReflect() protoreflect.Message {
	mi := &file_tvix_build_protos_build_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BuildRequest_EnvVar.ProtoReflect.Descriptor instead.
func (*BuildRequest_EnvVar) Descriptor() ([]byte, []int) {
	return file_tvix_build_protos_build_proto_rawDescGZIP(), []int{0, 0}
}

func (x *BuildRequest_EnvVar) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *BuildRequest_EnvVar) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

// BuildConstraints represents certain conditions that must be fulfilled
// inside the build environment to be able to build this.
// Constraints can be things like required architecture and minimum amount of memory.
// The required input paths are *not* represented in here, because it
// wouldn't be hermetic enough - see the comment around inputs too.
type BuildRequest_BuildConstraints struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The system that's needed to execute the build.
	System string `protobuf:"bytes,1,opt,name=system,proto3" json:"system,omitempty"`
	// The amount of memory required to be available for the build, in bytes.
	MinMemory uint64 `protobuf:"varint,2,opt,name=min_memory,json=minMemory,proto3" json:"min_memory,omitempty"`
	// A list of (absolute) paths that need to be available in the build
	// environment.
	// TBD, This is probably things like /dev/kvm, but no nix store paths.
	AvailableRoPaths []string `protobuf:"bytes,3,rep,name=available_ro_paths,json=availableRoPaths,proto3" json:"available_ro_paths,omitempty"`
}

func (x *BuildRequest_BuildConstraints) Reset() {
	*x = BuildRequest_BuildConstraints{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tvix_build_protos_build_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BuildRequest_BuildConstraints) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BuildRequest_BuildConstraints) ProtoMessage() {}

func (x *BuildRequest_BuildConstraints) ProtoReflect() protoreflect.Message {
	mi := &file_tvix_build_protos_build_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BuildRequest_BuildConstraints.ProtoReflect.Descriptor instead.
func (*BuildRequest_BuildConstraints) Descriptor() ([]byte, []int) {
	return file_tvix_build_protos_build_proto_rawDescGZIP(), []int{0, 1}
}

func (x *BuildRequest_BuildConstraints) GetSystem() string {
	if x != nil {
		return x.System
	}
	return ""
}

func (x *BuildRequest_BuildConstraints) GetMinMemory() uint64 {
	if x != nil {
		return x.MinMemory
	}
	return 0
}

func (x *BuildRequest_BuildConstraints) GetAvailableRoPaths() []string {
	if x != nil {
		return x.AvailableRoPaths
	}
	return nil
}

var File_tvix_build_protos_build_proto protoreflect.FileDescriptor

var file_tvix_build_protos_build_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x74, 0x76, 0x69, 0x78, 0x2f, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x73, 0x2f, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x0d, 0x74, 0x76, 0x69, 0x78, 0x2e, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x2e, 0x76, 0x31, 0x1a, 0x21,
	0x74, 0x76, 0x69, 0x78, 0x2f, 0x63, 0x61, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x73, 0x2f, 0x63, 0x61, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0xc4, 0x03, 0x0a, 0x0c, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x5f, 0x61, 0x72,
	0x67, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0b, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e,
	0x64, 0x41, 0x72, 0x67, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x73,
	0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x6f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x73, 0x12,
	0x4d, 0x0a, 0x10, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x76,
	0x61, 0x72, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x74, 0x76, 0x69, 0x78,
	0x2e, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x45, 0x6e, 0x76, 0x56, 0x61, 0x72, 0x52, 0x0f, 0x65,
	0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x56, 0x61, 0x72, 0x73, 0x12, 0x2d,
	0x0a, 0x06, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15,
	0x2e, 0x74, 0x76, 0x69, 0x78, 0x2e, 0x63, 0x61, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31,
	0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x06, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x73, 0x12, 0x4e, 0x0a,
	0x0b, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x72, 0x61, 0x69, 0x6e, 0x74, 0x73, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x74, 0x76, 0x69, 0x78, 0x2e, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x2e,
	0x76, 0x31, 0x2e, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e,
	0x42, 0x75, 0x69, 0x6c, 0x64, 0x43, 0x6f, 0x6e, 0x73, 0x74, 0x72, 0x61, 0x69, 0x6e, 0x74, 0x73,
	0x52, 0x0b, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x72, 0x61, 0x69, 0x6e, 0x74, 0x73, 0x1a, 0x30, 0x0a,
	0x06, 0x45, 0x6e, 0x76, 0x56, 0x61, 0x72, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x1a,
	0x77, 0x0a, 0x10, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x43, 0x6f, 0x6e, 0x73, 0x74, 0x72, 0x61, 0x69,
	0x6e, 0x74, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x12, 0x1d, 0x0a, 0x0a, 0x6d,
	0x69, 0x6e, 0x5f, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x09, 0x6d, 0x69, 0x6e, 0x4d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x12, 0x2c, 0x0a, 0x12, 0x61, 0x76,
	0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x5f, 0x72, 0x6f, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x73,
	0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x10, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c,
	0x65, 0x52, 0x6f, 0x50, 0x61, 0x74, 0x68, 0x73, 0x22, 0x7a, 0x0a, 0x05, 0x42, 0x75, 0x69, 0x6c,
	0x64, 0x12, 0x40, 0x0a, 0x0d, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x74, 0x76, 0x69, 0x78, 0x2e,
	0x62, 0x75, 0x69, 0x6c, 0x64, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x0c, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x2f, 0x0a, 0x07, 0x6f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x73, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x74, 0x76, 0x69, 0x78, 0x2e, 0x63, 0x61, 0x73, 0x74,
	0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x07, 0x6f, 0x75, 0x74,
	0x70, 0x75, 0x74, 0x73, 0x42, 0x24, 0x5a, 0x22, 0x63, 0x6f, 0x64, 0x65, 0x2e, 0x74, 0x76, 0x6c,
	0x2e, 0x66, 0x79, 0x69, 0x2f, 0x74, 0x76, 0x69, 0x78, 0x2f, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x2d,
	0x67, 0x6f, 0x3b, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_tvix_build_protos_build_proto_rawDescOnce sync.Once
	file_tvix_build_protos_build_proto_rawDescData = file_tvix_build_protos_build_proto_rawDesc
)

func file_tvix_build_protos_build_proto_rawDescGZIP() []byte {
	file_tvix_build_protos_build_proto_rawDescOnce.Do(func() {
		file_tvix_build_protos_build_proto_rawDescData = protoimpl.X.CompressGZIP(file_tvix_build_protos_build_proto_rawDescData)
	})
	return file_tvix_build_protos_build_proto_rawDescData
}

var file_tvix_build_protos_build_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_tvix_build_protos_build_proto_goTypes = []interface{}{
	(*BuildRequest)(nil),                  // 0: tvix.build.v1.BuildRequest
	(*Build)(nil),                         // 1: tvix.build.v1.Build
	(*BuildRequest_EnvVar)(nil),           // 2: tvix.build.v1.BuildRequest.EnvVar
	(*BuildRequest_BuildConstraints)(nil), // 3: tvix.build.v1.BuildRequest.BuildConstraints
	(*castore_go.Node)(nil),               // 4: tvix.castore.v1.Node
}
var file_tvix_build_protos_build_proto_depIdxs = []int32{
	2, // 0: tvix.build.v1.BuildRequest.environment_vars:type_name -> tvix.build.v1.BuildRequest.EnvVar
	4, // 1: tvix.build.v1.BuildRequest.inputs:type_name -> tvix.castore.v1.Node
	3, // 2: tvix.build.v1.BuildRequest.constraints:type_name -> tvix.build.v1.BuildRequest.BuildConstraints
	0, // 3: tvix.build.v1.Build.build_request:type_name -> tvix.build.v1.BuildRequest
	4, // 4: tvix.build.v1.Build.outputs:type_name -> tvix.castore.v1.Node
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_tvix_build_protos_build_proto_init() }
func file_tvix_build_protos_build_proto_init() {
	if File_tvix_build_protos_build_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_tvix_build_protos_build_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BuildRequest); i {
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
		file_tvix_build_protos_build_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Build); i {
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
		file_tvix_build_protos_build_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BuildRequest_EnvVar); i {
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
		file_tvix_build_protos_build_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BuildRequest_BuildConstraints); i {
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
			RawDescriptor: file_tvix_build_protos_build_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_tvix_build_protos_build_proto_goTypes,
		DependencyIndexes: file_tvix_build_protos_build_proto_depIdxs,
		MessageInfos:      file_tvix_build_protos_build_proto_msgTypes,
	}.Build()
	File_tvix_build_protos_build_proto = out.File
	file_tvix_build_protos_build_proto_rawDesc = nil
	file_tvix_build_protos_build_proto_goTypes = nil
	file_tvix_build_protos_build_proto_depIdxs = nil
}