// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: reddit/reddit_feed_api/v1/reddit_feed_api.proto

package reddit_feed_api

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type CreatePostsV1Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Posts []*Post `protobuf:"bytes,1,rep,name=posts,proto3" json:"posts,omitempty"`
}

func (x *CreatePostsV1Request) Reset() {
	*x = CreatePostsV1Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_reddit_reddit_feed_api_v1_reddit_feed_api_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreatePostsV1Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreatePostsV1Request) ProtoMessage() {}

func (x *CreatePostsV1Request) ProtoReflect() protoreflect.Message {
	mi := &file_reddit_reddit_feed_api_v1_reddit_feed_api_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreatePostsV1Request.ProtoReflect.Descriptor instead.
func (*CreatePostsV1Request) Descriptor() ([]byte, []int) {
	return file_reddit_reddit_feed_api_v1_reddit_feed_api_proto_rawDescGZIP(), []int{0}
}

func (x *CreatePostsV1Request) GetPosts() []*Post {
	if x != nil {
		return x.Posts
	}
	return nil
}

type CreatePostsV1Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NumberOfCreatedPosts int64 `protobuf:"varint,1,opt,name=number_of_created_posts,json=numberOfCreatedPosts,proto3" json:"number_of_created_posts,omitempty"`
}

func (x *CreatePostsV1Response) Reset() {
	*x = CreatePostsV1Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_reddit_reddit_feed_api_v1_reddit_feed_api_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreatePostsV1Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreatePostsV1Response) ProtoMessage() {}

func (x *CreatePostsV1Response) ProtoReflect() protoreflect.Message {
	mi := &file_reddit_reddit_feed_api_v1_reddit_feed_api_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreatePostsV1Response.ProtoReflect.Descriptor instead.
func (*CreatePostsV1Response) Descriptor() ([]byte, []int) {
	return file_reddit_reddit_feed_api_v1_reddit_feed_api_proto_rawDescGZIP(), []int{1}
}

func (x *CreatePostsV1Response) GetNumberOfCreatedPosts() int64 {
	if x != nil {
		return x.NumberOfCreatedPosts
	}
	return 0
}

type GenerateFeedV1Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PageId uint64 `protobuf:"varint,1,opt,name=page_id,json=pageId,proto3" json:"page_id,omitempty"`
}

func (x *GenerateFeedV1Request) Reset() {
	*x = GenerateFeedV1Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_reddit_reddit_feed_api_v1_reddit_feed_api_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GenerateFeedV1Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GenerateFeedV1Request) ProtoMessage() {}

func (x *GenerateFeedV1Request) ProtoReflect() protoreflect.Message {
	mi := &file_reddit_reddit_feed_api_v1_reddit_feed_api_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GenerateFeedV1Request.ProtoReflect.Descriptor instead.
func (*GenerateFeedV1Request) Descriptor() ([]byte, []int) {
	return file_reddit_reddit_feed_api_v1_reddit_feed_api_proto_rawDescGZIP(), []int{2}
}

func (x *GenerateFeedV1Request) GetPageId() uint64 {
	if x != nil {
		return x.PageId
	}
	return 0
}

type GenerateFeedV1Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Posts []*Post `protobuf:"bytes,1,rep,name=posts,proto3" json:"posts,omitempty"`
}

func (x *GenerateFeedV1Response) Reset() {
	*x = GenerateFeedV1Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_reddit_reddit_feed_api_v1_reddit_feed_api_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GenerateFeedV1Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GenerateFeedV1Response) ProtoMessage() {}

func (x *GenerateFeedV1Response) ProtoReflect() protoreflect.Message {
	mi := &file_reddit_reddit_feed_api_v1_reddit_feed_api_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GenerateFeedV1Response.ProtoReflect.Descriptor instead.
func (*GenerateFeedV1Response) Descriptor() ([]byte, []int) {
	return file_reddit_reddit_feed_api_v1_reddit_feed_api_proto_rawDescGZIP(), []int{3}
}

func (x *GenerateFeedV1Response) GetPosts() []*Post {
	if x != nil {
		return x.Posts
	}
	return nil
}

type Post struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Title     string `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Author    string `protobuf:"bytes,2,opt,name=author,proto3" json:"author,omitempty"`
	Subreddit string `protobuf:"bytes,3,opt,name=subreddit,proto3" json:"subreddit,omitempty"`
	// Types that are assignable to PostType:
	//	*Post_Link
	//	*Post_Content
	PostType       isPost_PostType `protobuf_oneof:"post_type"`
	Score          uint64          `protobuf:"varint,6,opt,name=score,proto3" json:"score,omitempty"`
	Promoted       bool            `protobuf:"varint,7,opt,name=promoted,proto3" json:"promoted,omitempty"`
	NotSafeForWork bool            `protobuf:"varint,8,opt,name=not_safe_for_work,json=notSafeForWork,proto3" json:"not_safe_for_work,omitempty"`
}

func (x *Post) Reset() {
	*x = Post{}
	if protoimpl.UnsafeEnabled {
		mi := &file_reddit_reddit_feed_api_v1_reddit_feed_api_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Post) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Post) ProtoMessage() {}

func (x *Post) ProtoReflect() protoreflect.Message {
	mi := &file_reddit_reddit_feed_api_v1_reddit_feed_api_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Post.ProtoReflect.Descriptor instead.
func (*Post) Descriptor() ([]byte, []int) {
	return file_reddit_reddit_feed_api_v1_reddit_feed_api_proto_rawDescGZIP(), []int{4}
}

func (x *Post) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Post) GetAuthor() string {
	if x != nil {
		return x.Author
	}
	return ""
}

func (x *Post) GetSubreddit() string {
	if x != nil {
		return x.Subreddit
	}
	return ""
}

func (m *Post) GetPostType() isPost_PostType {
	if m != nil {
		return m.PostType
	}
	return nil
}

func (x *Post) GetLink() string {
	if x, ok := x.GetPostType().(*Post_Link); ok {
		return x.Link
	}
	return ""
}

func (x *Post) GetContent() string {
	if x, ok := x.GetPostType().(*Post_Content); ok {
		return x.Content
	}
	return ""
}

func (x *Post) GetScore() uint64 {
	if x != nil {
		return x.Score
	}
	return 0
}

func (x *Post) GetPromoted() bool {
	if x != nil {
		return x.Promoted
	}
	return false
}

func (x *Post) GetNotSafeForWork() bool {
	if x != nil {
		return x.NotSafeForWork
	}
	return false
}

type isPost_PostType interface {
	isPost_PostType()
}

type Post_Link struct {
	Link string `protobuf:"bytes,4,opt,name=link,proto3,oneof"`
}

type Post_Content struct {
	Content string `protobuf:"bytes,5,opt,name=content,proto3,oneof"`
}

func (*Post_Link) isPost_PostType() {}

func (*Post_Content) isPost_PostType() {}

var File_reddit_reddit_feed_api_v1_reddit_feed_api_proto protoreflect.FileDescriptor

var file_reddit_reddit_feed_api_v1_reddit_feed_api_proto_rawDesc = []byte{
	0x0a, 0x2f, 0x72, 0x65, 0x64, 0x64, 0x69, 0x74, 0x2f, 0x72, 0x65, 0x64, 0x64, 0x69, 0x74, 0x5f,
	0x66, 0x65, 0x65, 0x64, 0x5f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x72, 0x65, 0x64, 0x64,
	0x69, 0x74, 0x5f, 0x66, 0x65, 0x65, 0x64, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x19, 0x72, 0x65, 0x64, 0x64, 0x69, 0x74, 0x2e, 0x72, 0x65, 0x64, 0x64, 0x69, 0x74,
	0x5f, 0x66, 0x65, 0x65, 0x64, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x1a, 0x1c, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x76, 0x32,
	0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69,
	0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x5a, 0x0a, 0x14, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x50, 0x6f, 0x73,
	0x74, 0x73, 0x56, 0x31, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x42, 0x0a, 0x05, 0x70,
	0x6f, 0x73, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x72, 0x65, 0x64,
	0x64, 0x69, 0x74, 0x2e, 0x72, 0x65, 0x64, 0x64, 0x69, 0x74, 0x5f, 0x66, 0x65, 0x65, 0x64, 0x5f,
	0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x6f, 0x73, 0x74, 0x42, 0x0b, 0xfa, 0x42, 0x08,
	0x92, 0x01, 0x05, 0x08, 0x01, 0x10, 0x80, 0x08, 0x52, 0x05, 0x70, 0x6f, 0x73, 0x74, 0x73, 0x22,
	0x4e, 0x0a, 0x15, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x50, 0x6f, 0x73, 0x74, 0x73, 0x56, 0x31,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x35, 0x0a, 0x17, 0x6e, 0x75, 0x6d, 0x62,
	0x65, 0x72, 0x5f, 0x6f, 0x66, 0x5f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x70, 0x6f,
	0x73, 0x74, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x14, 0x6e, 0x75, 0x6d, 0x62, 0x65,
	0x72, 0x4f, 0x66, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x50, 0x6f, 0x73, 0x74, 0x73, 0x22,
	0x39, 0x0a, 0x15, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x46, 0x65, 0x65, 0x64, 0x56,
	0x31, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x20, 0x0a, 0x07, 0x70, 0x61, 0x67, 0x65,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x32, 0x02,
	0x28, 0x01, 0x52, 0x06, 0x70, 0x61, 0x67, 0x65, 0x49, 0x64, 0x22, 0x4f, 0x0a, 0x16, 0x47, 0x65,
	0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x46, 0x65, 0x65, 0x64, 0x56, 0x31, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x35, 0x0a, 0x05, 0x70, 0x6f, 0x73, 0x74, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x72, 0x65, 0x64, 0x64, 0x69, 0x74, 0x2e, 0x72, 0x65, 0x64,
	0x64, 0x69, 0x74, 0x5f, 0x66, 0x65, 0x65, 0x64, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e,
	0x50, 0x6f, 0x73, 0x74, 0x52, 0x05, 0x70, 0x6f, 0x73, 0x74, 0x73, 0x22, 0xa1, 0x02, 0x0a, 0x04,
	0x50, 0x6f, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x36, 0x0a, 0x06, 0x61, 0x75,
	0x74, 0x68, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x1e, 0xfa, 0x42, 0x1b, 0x72,
	0x19, 0x28, 0x0b, 0x32, 0x10, 0x5e, 0x74, 0x32, 0x5f, 0x5b, 0x61, 0x2d, 0x7a, 0x30, 0x2d, 0x39,
	0x5d, 0x7b, 0x38, 0x7d, 0x24, 0x3a, 0x03, 0x74, 0x32, 0x5f, 0x52, 0x06, 0x61, 0x75, 0x74, 0x68,
	0x6f, 0x72, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x75, 0x62, 0x72, 0x65, 0x64, 0x64, 0x69, 0x74, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x75, 0x62, 0x72, 0x65, 0x64, 0x64, 0x69, 0x74,
	0x12, 0x1e, 0x0a, 0x04, 0x6c, 0x69, 0x6e, 0x6b, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08,
	0xfa, 0x42, 0x05, 0x72, 0x03, 0x90, 0x01, 0x01, 0x48, 0x00, 0x52, 0x04, 0x6c, 0x69, 0x6e, 0x6b,
	0x12, 0x1a, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x48, 0x00, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x1d, 0x0a, 0x05,
	0x73, 0x63, 0x6f, 0x72, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x04, 0x42, 0x07, 0xfa, 0x42, 0x04,
	0x32, 0x02, 0x28, 0x00, 0x52, 0x05, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70,
	0x72, 0x6f, 0x6d, 0x6f, 0x74, 0x65, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x70,
	0x72, 0x6f, 0x6d, 0x6f, 0x74, 0x65, 0x64, 0x12, 0x29, 0x0a, 0x11, 0x6e, 0x6f, 0x74, 0x5f, 0x73,
	0x61, 0x66, 0x65, 0x5f, 0x66, 0x6f, 0x72, 0x5f, 0x77, 0x6f, 0x72, 0x6b, 0x18, 0x08, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x0e, 0x6e, 0x6f, 0x74, 0x53, 0x61, 0x66, 0x65, 0x46, 0x6f, 0x72, 0x57, 0x6f,
	0x72, 0x6b, 0x42, 0x0b, 0x0a, 0x09, 0x70, 0x6f, 0x73, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x32,
	0xb2, 0x02, 0x0a, 0x14, 0x52, 0x65, 0x64, 0x64, 0x69, 0x74, 0x46, 0x65, 0x65, 0x64, 0x41, 0x50,
	0x49, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x85, 0x01, 0x0a, 0x0d, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x50, 0x6f, 0x73, 0x74, 0x73, 0x56, 0x31, 0x12, 0x2f, 0x2e, 0x72, 0x65, 0x64,
	0x64, 0x69, 0x74, 0x2e, 0x72, 0x65, 0x64, 0x64, 0x69, 0x74, 0x5f, 0x66, 0x65, 0x65, 0x64, 0x5f,
	0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x50, 0x6f, 0x73,
	0x74, 0x73, 0x56, 0x31, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x30, 0x2e, 0x72, 0x65,
	0x64, 0x64, 0x69, 0x74, 0x2e, 0x72, 0x65, 0x64, 0x64, 0x69, 0x74, 0x5f, 0x66, 0x65, 0x65, 0x64,
	0x5f, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x50, 0x6f,
	0x73, 0x74, 0x73, 0x56, 0x31, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x11, 0x82,
	0xd3, 0xe4, 0x93, 0x02, 0x0b, 0x22, 0x09, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x6f, 0x73, 0x74, 0x73,
	0x12, 0x91, 0x01, 0x0a, 0x0e, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x46, 0x65, 0x65,
	0x64, 0x56, 0x31, 0x12, 0x30, 0x2e, 0x72, 0x65, 0x64, 0x64, 0x69, 0x74, 0x2e, 0x72, 0x65, 0x64,
	0x64, 0x69, 0x74, 0x5f, 0x66, 0x65, 0x65, 0x64, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e,
	0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x46, 0x65, 0x65, 0x64, 0x56, 0x31, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x31, 0x2e, 0x72, 0x65, 0x64, 0x64, 0x69, 0x74, 0x2e, 0x72,
	0x65, 0x64, 0x64, 0x69, 0x74, 0x5f, 0x66, 0x65, 0x65, 0x64, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x76,
	0x31, 0x2e, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x46, 0x65, 0x65, 0x64, 0x56, 0x31,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1a, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x14,
	0x12, 0x12, 0x2f, 0x76, 0x31, 0x2f, 0x66, 0x65, 0x65, 0x64, 0x2f, 0x7b, 0x70, 0x61, 0x67, 0x65,
	0x5f, 0x69, 0x64, 0x7d, 0x42, 0xa6, 0x01, 0x5a, 0x45, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x72, 0x74, 0x74, 0x65, 0x74, 0x2f, 0x72, 0x65, 0x64, 0x64, 0x69,
	0x74, 0x2d, 0x66, 0x65, 0x65, 0x64, 0x2d, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x72,
	0x65, 0x64, 0x64, 0x69, 0x74, 0x2d, 0x66, 0x65, 0x65, 0x64, 0x2d, 0x61, 0x70, 0x69, 0x3b, 0x72,
	0x65, 0x64, 0x64, 0x69, 0x74, 0x5f, 0x66, 0x65, 0x65, 0x64, 0x5f, 0x61, 0x70, 0x69, 0x92, 0x41,
	0x5c, 0x12, 0x5a, 0x0a, 0x0f, 0x52, 0x65, 0x64, 0x64, 0x69, 0x74, 0x20, 0x46, 0x65, 0x65, 0x64,
	0x20, 0x41, 0x50, 0x49, 0x2a, 0x42, 0x0a, 0x03, 0x4d, 0x49, 0x54, 0x12, 0x3b, 0x68, 0x74, 0x74,
	0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x61, 0x72, 0x74, 0x74, 0x65, 0x74, 0x2f, 0x72, 0x65, 0x64, 0x64, 0x69, 0x74, 0x2d, 0x66, 0x65,
	0x65, 0x64, 0x2d, 0x61, 0x70, 0x69, 0x2f, 0x62, 0x6c, 0x6f, 0x62, 0x2f, 0x6d, 0x61, 0x69, 0x6e,
	0x2f, 0x4c, 0x49, 0x43, 0x45, 0x4e, 0x53, 0x45, 0x32, 0x03, 0x31, 0x2e, 0x30, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_reddit_reddit_feed_api_v1_reddit_feed_api_proto_rawDescOnce sync.Once
	file_reddit_reddit_feed_api_v1_reddit_feed_api_proto_rawDescData = file_reddit_reddit_feed_api_v1_reddit_feed_api_proto_rawDesc
)

func file_reddit_reddit_feed_api_v1_reddit_feed_api_proto_rawDescGZIP() []byte {
	file_reddit_reddit_feed_api_v1_reddit_feed_api_proto_rawDescOnce.Do(func() {
		file_reddit_reddit_feed_api_v1_reddit_feed_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_reddit_reddit_feed_api_v1_reddit_feed_api_proto_rawDescData)
	})
	return file_reddit_reddit_feed_api_v1_reddit_feed_api_proto_rawDescData
}

var file_reddit_reddit_feed_api_v1_reddit_feed_api_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_reddit_reddit_feed_api_v1_reddit_feed_api_proto_goTypes = []interface{}{
	(*CreatePostsV1Request)(nil),   // 0: reddit.reddit_feed_api.v1.CreatePostsV1Request
	(*CreatePostsV1Response)(nil),  // 1: reddit.reddit_feed_api.v1.CreatePostsV1Response
	(*GenerateFeedV1Request)(nil),  // 2: reddit.reddit_feed_api.v1.GenerateFeedV1Request
	(*GenerateFeedV1Response)(nil), // 3: reddit.reddit_feed_api.v1.GenerateFeedV1Response
	(*Post)(nil),                   // 4: reddit.reddit_feed_api.v1.Post
}
var file_reddit_reddit_feed_api_v1_reddit_feed_api_proto_depIdxs = []int32{
	4, // 0: reddit.reddit_feed_api.v1.CreatePostsV1Request.posts:type_name -> reddit.reddit_feed_api.v1.Post
	4, // 1: reddit.reddit_feed_api.v1.GenerateFeedV1Response.posts:type_name -> reddit.reddit_feed_api.v1.Post
	0, // 2: reddit.reddit_feed_api.v1.RedditFeedAPIService.CreatePostsV1:input_type -> reddit.reddit_feed_api.v1.CreatePostsV1Request
	2, // 3: reddit.reddit_feed_api.v1.RedditFeedAPIService.GenerateFeedV1:input_type -> reddit.reddit_feed_api.v1.GenerateFeedV1Request
	1, // 4: reddit.reddit_feed_api.v1.RedditFeedAPIService.CreatePostsV1:output_type -> reddit.reddit_feed_api.v1.CreatePostsV1Response
	3, // 5: reddit.reddit_feed_api.v1.RedditFeedAPIService.GenerateFeedV1:output_type -> reddit.reddit_feed_api.v1.GenerateFeedV1Response
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_reddit_reddit_feed_api_v1_reddit_feed_api_proto_init() }
func file_reddit_reddit_feed_api_v1_reddit_feed_api_proto_init() {
	if File_reddit_reddit_feed_api_v1_reddit_feed_api_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_reddit_reddit_feed_api_v1_reddit_feed_api_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreatePostsV1Request); i {
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
		file_reddit_reddit_feed_api_v1_reddit_feed_api_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreatePostsV1Response); i {
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
		file_reddit_reddit_feed_api_v1_reddit_feed_api_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GenerateFeedV1Request); i {
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
		file_reddit_reddit_feed_api_v1_reddit_feed_api_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GenerateFeedV1Response); i {
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
		file_reddit_reddit_feed_api_v1_reddit_feed_api_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Post); i {
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
	file_reddit_reddit_feed_api_v1_reddit_feed_api_proto_msgTypes[4].OneofWrappers = []interface{}{
		(*Post_Link)(nil),
		(*Post_Content)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_reddit_reddit_feed_api_v1_reddit_feed_api_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_reddit_reddit_feed_api_v1_reddit_feed_api_proto_goTypes,
		DependencyIndexes: file_reddit_reddit_feed_api_v1_reddit_feed_api_proto_depIdxs,
		MessageInfos:      file_reddit_reddit_feed_api_v1_reddit_feed_api_proto_msgTypes,
	}.Build()
	File_reddit_reddit_feed_api_v1_reddit_feed_api_proto = out.File
	file_reddit_reddit_feed_api_v1_reddit_feed_api_proto_rawDesc = nil
	file_reddit_reddit_feed_api_v1_reddit_feed_api_proto_goTypes = nil
	file_reddit_reddit_feed_api_v1_reddit_feed_api_proto_depIdxs = nil
}
