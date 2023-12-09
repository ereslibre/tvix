// SPDX-FileCopyrightText: edef <edef@unfathomable.blue>
// SPDX-License-Identifier: OSL-3.0 OR MIT OR Apache-2.0

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: tvix/castore/protos/castore.proto

package castorev1

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

// A Directory can contain Directory, File or Symlink nodes.
// Each of these nodes have a name attribute, which is the basename in that
// directory and node type specific attributes.
// The name attribute:
//   - MUST not contain slashes or null bytes
//   - MUST not be '.' or '..'
//   - MUST be unique across all three lists
//
// Elements in each list need to be lexicographically ordered by the name
// attribute.
type Directory struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Directories []*DirectoryNode `protobuf:"bytes,1,rep,name=directories,proto3" json:"directories,omitempty"`
	Files       []*FileNode      `protobuf:"bytes,2,rep,name=files,proto3" json:"files,omitempty"`
	Symlinks    []*SymlinkNode   `protobuf:"bytes,3,rep,name=symlinks,proto3" json:"symlinks,omitempty"`
}

func (x *Directory) Reset() {
	*x = Directory{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tvix_castore_protos_castore_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Directory) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Directory) ProtoMessage() {}

func (x *Directory) ProtoReflect() protoreflect.Message {
	mi := &file_tvix_castore_protos_castore_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Directory.ProtoReflect.Descriptor instead.
func (*Directory) Descriptor() ([]byte, []int) {
	return file_tvix_castore_protos_castore_proto_rawDescGZIP(), []int{0}
}

func (x *Directory) GetDirectories() []*DirectoryNode {
	if x != nil {
		return x.Directories
	}
	return nil
}

func (x *Directory) GetFiles() []*FileNode {
	if x != nil {
		return x.Files
	}
	return nil
}

func (x *Directory) GetSymlinks() []*SymlinkNode {
	if x != nil {
		return x.Symlinks
	}
	return nil
}

// A DirectoryNode represents a directory in a Directory.
type DirectoryNode struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The (base)name of the directory
	Name []byte `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// The blake3 hash of a Directory message, serialized in protobuf canonical form.
	Digest []byte `protobuf:"bytes,2,opt,name=digest,proto3" json:"digest,omitempty"`
	// Number of child elements in the Directory referred to by `digest`.
	// Calculated by summing up the numbers of `directories`, `files` and
	// `symlinks`, and for each directory, its size field. Used for inode number
	// calculation.
	// This field is precisely as verifiable as any other Merkle tree edge.
	// Resolve `digest`, and you can compute it incrementally. Resolve the entire
	// tree, and you can fully compute it from scratch.
	// A credulous implementation won't reject an excessive size, but this is
	// harmless: you'll have some ordinals without nodes. Undersizing is obvious
	// and easy to reject: you won't have an ordinal for some nodes.
	Size uint64 `protobuf:"varint,3,opt,name=size,proto3" json:"size,omitempty"`
}

func (x *DirectoryNode) Reset() {
	*x = DirectoryNode{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tvix_castore_protos_castore_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DirectoryNode) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DirectoryNode) ProtoMessage() {}

func (x *DirectoryNode) ProtoReflect() protoreflect.Message {
	mi := &file_tvix_castore_protos_castore_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DirectoryNode.ProtoReflect.Descriptor instead.
func (*DirectoryNode) Descriptor() ([]byte, []int) {
	return file_tvix_castore_protos_castore_proto_rawDescGZIP(), []int{1}
}

func (x *DirectoryNode) GetName() []byte {
	if x != nil {
		return x.Name
	}
	return nil
}

func (x *DirectoryNode) GetDigest() []byte {
	if x != nil {
		return x.Digest
	}
	return nil
}

func (x *DirectoryNode) GetSize() uint64 {
	if x != nil {
		return x.Size
	}
	return 0
}

// A FileNode represents a regular or executable file in a Directory.
type FileNode struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The (base)name of the file
	Name []byte `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// The blake3 digest of the file contents
	Digest []byte `protobuf:"bytes,2,opt,name=digest,proto3" json:"digest,omitempty"`
	// The file content size
	Size uint64 `protobuf:"varint,3,opt,name=size,proto3" json:"size,omitempty"`
	// Whether the file is executable
	Executable bool `protobuf:"varint,4,opt,name=executable,proto3" json:"executable,omitempty"`
}

func (x *FileNode) Reset() {
	*x = FileNode{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tvix_castore_protos_castore_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileNode) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileNode) ProtoMessage() {}

func (x *FileNode) ProtoReflect() protoreflect.Message {
	mi := &file_tvix_castore_protos_castore_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileNode.ProtoReflect.Descriptor instead.
func (*FileNode) Descriptor() ([]byte, []int) {
	return file_tvix_castore_protos_castore_proto_rawDescGZIP(), []int{2}
}

func (x *FileNode) GetName() []byte {
	if x != nil {
		return x.Name
	}
	return nil
}

func (x *FileNode) GetDigest() []byte {
	if x != nil {
		return x.Digest
	}
	return nil
}

func (x *FileNode) GetSize() uint64 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *FileNode) GetExecutable() bool {
	if x != nil {
		return x.Executable
	}
	return false
}

// A SymlinkNode represents a symbolic link in a Directory.
type SymlinkNode struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The (base)name of the symlink
	Name []byte `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// The target of the symlink.
	Target []byte `protobuf:"bytes,2,opt,name=target,proto3" json:"target,omitempty"`
}

func (x *SymlinkNode) Reset() {
	*x = SymlinkNode{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tvix_castore_protos_castore_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SymlinkNode) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SymlinkNode) ProtoMessage() {}

func (x *SymlinkNode) ProtoReflect() protoreflect.Message {
	mi := &file_tvix_castore_protos_castore_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SymlinkNode.ProtoReflect.Descriptor instead.
func (*SymlinkNode) Descriptor() ([]byte, []int) {
	return file_tvix_castore_protos_castore_proto_rawDescGZIP(), []int{3}
}

func (x *SymlinkNode) GetName() []byte {
	if x != nil {
		return x.Name
	}
	return nil
}

func (x *SymlinkNode) GetTarget() []byte {
	if x != nil {
		return x.Target
	}
	return nil
}

// A Node is either a DirectoryNode, FileNode or SymlinkNode.
type Node struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Node:
	//
	//	*Node_Directory
	//	*Node_File
	//	*Node_Symlink
	Node isNode_Node `protobuf_oneof:"node"`
}

func (x *Node) Reset() {
	*x = Node{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tvix_castore_protos_castore_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Node) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Node) ProtoMessage() {}

func (x *Node) ProtoReflect() protoreflect.Message {
	mi := &file_tvix_castore_protos_castore_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Node.ProtoReflect.Descriptor instead.
func (*Node) Descriptor() ([]byte, []int) {
	return file_tvix_castore_protos_castore_proto_rawDescGZIP(), []int{4}
}

func (m *Node) GetNode() isNode_Node {
	if m != nil {
		return m.Node
	}
	return nil
}

func (x *Node) GetDirectory() *DirectoryNode {
	if x, ok := x.GetNode().(*Node_Directory); ok {
		return x.Directory
	}
	return nil
}

func (x *Node) GetFile() *FileNode {
	if x, ok := x.GetNode().(*Node_File); ok {
		return x.File
	}
	return nil
}

func (x *Node) GetSymlink() *SymlinkNode {
	if x, ok := x.GetNode().(*Node_Symlink); ok {
		return x.Symlink
	}
	return nil
}

type isNode_Node interface {
	isNode_Node()
}

type Node_Directory struct {
	Directory *DirectoryNode `protobuf:"bytes,1,opt,name=directory,proto3,oneof"`
}

type Node_File struct {
	File *FileNode `protobuf:"bytes,2,opt,name=file,proto3,oneof"`
}

type Node_Symlink struct {
	Symlink *SymlinkNode `protobuf:"bytes,3,opt,name=symlink,proto3,oneof"`
}

func (*Node_Directory) isNode_Node() {}

func (*Node_File) isNode_Node() {}

func (*Node_Symlink) isNode_Node() {}

var File_tvix_castore_protos_castore_proto protoreflect.FileDescriptor

var file_tvix_castore_protos_castore_proto_rawDesc = []byte{
	0x0a, 0x21, 0x74, 0x76, 0x69, 0x78, 0x2f, 0x63, 0x61, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x63, 0x61, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x0f, 0x74, 0x76, 0x69, 0x78, 0x2e, 0x63, 0x61, 0x73, 0x74, 0x6f, 0x72,
	0x65, 0x2e, 0x76, 0x31, 0x22, 0xb8, 0x01, 0x0a, 0x09, 0x44, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f,
	0x72, 0x79, 0x12, 0x40, 0x0a, 0x0b, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x69, 0x65,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x74, 0x76, 0x69, 0x78, 0x2e, 0x63,
	0x61, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x69, 0x72, 0x65, 0x63, 0x74,
	0x6f, 0x72, 0x79, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x0b, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f,
	0x72, 0x69, 0x65, 0x73, 0x12, 0x2f, 0x0a, 0x05, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x18, 0x02, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x74, 0x76, 0x69, 0x78, 0x2e, 0x63, 0x61, 0x73, 0x74, 0x6f,
	0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x05,
	0x66, 0x69, 0x6c, 0x65, 0x73, 0x12, 0x38, 0x0a, 0x08, 0x73, 0x79, 0x6d, 0x6c, 0x69, 0x6e, 0x6b,
	0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x74, 0x76, 0x69, 0x78, 0x2e, 0x63,
	0x61, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x79, 0x6d, 0x6c, 0x69, 0x6e,
	0x6b, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x08, 0x73, 0x79, 0x6d, 0x6c, 0x69, 0x6e, 0x6b, 0x73, 0x22,
	0x4f, 0x0a, 0x0d, 0x44, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x79, 0x4e, 0x6f, 0x64, 0x65,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x64, 0x69, 0x67, 0x65, 0x73, 0x74, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x64, 0x69, 0x67, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04,
	0x73, 0x69, 0x7a, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x73, 0x69, 0x7a, 0x65,
	0x22, 0x6a, 0x0a, 0x08, 0x46, 0x69, 0x6c, 0x65, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x16, 0x0a, 0x06, 0x64, 0x69, 0x67, 0x65, 0x73, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x06, 0x64, 0x69, 0x67, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x12, 0x1e, 0x0a, 0x0a,
	0x65, 0x78, 0x65, 0x63, 0x75, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x0a, 0x65, 0x78, 0x65, 0x63, 0x75, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x22, 0x39, 0x0a, 0x0b,
	0x53, 0x79, 0x6d, 0x6c, 0x69, 0x6e, 0x6b, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x16, 0x0a, 0x06, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x06, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x22, 0xb9, 0x01, 0x0a, 0x04, 0x4e, 0x6f, 0x64, 0x65,
	0x12, 0x3e, 0x0a, 0x09, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x74, 0x76, 0x69, 0x78, 0x2e, 0x63, 0x61, 0x73, 0x74, 0x6f,
	0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x79, 0x4e,
	0x6f, 0x64, 0x65, 0x48, 0x00, 0x52, 0x09, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x79,
	0x12, 0x2f, 0x0a, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19,
	0x2e, 0x74, 0x76, 0x69, 0x78, 0x2e, 0x63, 0x61, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31,
	0x2e, 0x46, 0x69, 0x6c, 0x65, 0x4e, 0x6f, 0x64, 0x65, 0x48, 0x00, 0x52, 0x04, 0x66, 0x69, 0x6c,
	0x65, 0x12, 0x38, 0x0a, 0x07, 0x73, 0x79, 0x6d, 0x6c, 0x69, 0x6e, 0x6b, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x74, 0x76, 0x69, 0x78, 0x2e, 0x63, 0x61, 0x73, 0x74, 0x6f, 0x72,
	0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x79, 0x6d, 0x6c, 0x69, 0x6e, 0x6b, 0x4e, 0x6f, 0x64, 0x65,
	0x48, 0x00, 0x52, 0x07, 0x73, 0x79, 0x6d, 0x6c, 0x69, 0x6e, 0x6b, 0x42, 0x06, 0x0a, 0x04, 0x6e,
	0x6f, 0x64, 0x65, 0x42, 0x28, 0x5a, 0x26, 0x63, 0x6f, 0x64, 0x65, 0x2e, 0x74, 0x76, 0x6c, 0x2e,
	0x66, 0x79, 0x69, 0x2f, 0x74, 0x76, 0x69, 0x78, 0x2f, 0x63, 0x61, 0x73, 0x74, 0x6f, 0x72, 0x65,
	0x2d, 0x67, 0x6f, 0x3b, 0x63, 0x61, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x76, 0x31, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_tvix_castore_protos_castore_proto_rawDescOnce sync.Once
	file_tvix_castore_protos_castore_proto_rawDescData = file_tvix_castore_protos_castore_proto_rawDesc
)

func file_tvix_castore_protos_castore_proto_rawDescGZIP() []byte {
	file_tvix_castore_protos_castore_proto_rawDescOnce.Do(func() {
		file_tvix_castore_protos_castore_proto_rawDescData = protoimpl.X.CompressGZIP(file_tvix_castore_protos_castore_proto_rawDescData)
	})
	return file_tvix_castore_protos_castore_proto_rawDescData
}

var file_tvix_castore_protos_castore_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_tvix_castore_protos_castore_proto_goTypes = []interface{}{
	(*Directory)(nil),     // 0: tvix.castore.v1.Directory
	(*DirectoryNode)(nil), // 1: tvix.castore.v1.DirectoryNode
	(*FileNode)(nil),      // 2: tvix.castore.v1.FileNode
	(*SymlinkNode)(nil),   // 3: tvix.castore.v1.SymlinkNode
	(*Node)(nil),          // 4: tvix.castore.v1.Node
}
var file_tvix_castore_protos_castore_proto_depIdxs = []int32{
	1, // 0: tvix.castore.v1.Directory.directories:type_name -> tvix.castore.v1.DirectoryNode
	2, // 1: tvix.castore.v1.Directory.files:type_name -> tvix.castore.v1.FileNode
	3, // 2: tvix.castore.v1.Directory.symlinks:type_name -> tvix.castore.v1.SymlinkNode
	1, // 3: tvix.castore.v1.Node.directory:type_name -> tvix.castore.v1.DirectoryNode
	2, // 4: tvix.castore.v1.Node.file:type_name -> tvix.castore.v1.FileNode
	3, // 5: tvix.castore.v1.Node.symlink:type_name -> tvix.castore.v1.SymlinkNode
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_tvix_castore_protos_castore_proto_init() }
func file_tvix_castore_protos_castore_proto_init() {
	if File_tvix_castore_protos_castore_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_tvix_castore_protos_castore_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Directory); i {
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
		file_tvix_castore_protos_castore_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DirectoryNode); i {
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
		file_tvix_castore_protos_castore_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileNode); i {
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
		file_tvix_castore_protos_castore_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SymlinkNode); i {
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
		file_tvix_castore_protos_castore_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Node); i {
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
	file_tvix_castore_protos_castore_proto_msgTypes[4].OneofWrappers = []interface{}{
		(*Node_Directory)(nil),
		(*Node_File)(nil),
		(*Node_Symlink)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_tvix_castore_protos_castore_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_tvix_castore_protos_castore_proto_goTypes,
		DependencyIndexes: file_tvix_castore_protos_castore_proto_depIdxs,
		MessageInfos:      file_tvix_castore_protos_castore_proto_msgTypes,
	}.Build()
	File_tvix_castore_protos_castore_proto = out.File
	file_tvix_castore_protos_castore_proto_rawDesc = nil
	file_tvix_castore_protos_castore_proto_goTypes = nil
	file_tvix_castore_protos_castore_proto_depIdxs = nil
}
