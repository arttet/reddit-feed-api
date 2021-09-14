// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: reddit/reddit_feed_api/v1/reddit_feed_api.proto

package reddit_feed_api

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
)

// Validate checks the field values on CreatePostsV1Request with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *CreatePostsV1Request) Validate() error {
	if m == nil {
		return nil
	}

	if l := len(m.GetPosts()); l < 1 || l > 1024 {
		return CreatePostsV1RequestValidationError{
			field:  "Posts",
			reason: "value must contain between 1 and 1024 items, inclusive",
		}
	}

	for idx, item := range m.GetPosts() {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return CreatePostsV1RequestValidationError{
					field:  fmt.Sprintf("Posts[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	return nil
}

// CreatePostsV1RequestValidationError is the validation error returned by
// CreatePostsV1Request.Validate if the designated constraints aren't met.
type CreatePostsV1RequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreatePostsV1RequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreatePostsV1RequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreatePostsV1RequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreatePostsV1RequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreatePostsV1RequestValidationError) ErrorName() string {
	return "CreatePostsV1RequestValidationError"
}

// Error satisfies the builtin error interface
func (e CreatePostsV1RequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreatePostsV1Request.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreatePostsV1RequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreatePostsV1RequestValidationError{}

// Validate checks the field values on CreatePostsV1Response with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *CreatePostsV1Response) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for NumberOfCreatedPosts

	return nil
}

// CreatePostsV1ResponseValidationError is the validation error returned by
// CreatePostsV1Response.Validate if the designated constraints aren't met.
type CreatePostsV1ResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreatePostsV1ResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreatePostsV1ResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreatePostsV1ResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreatePostsV1ResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreatePostsV1ResponseValidationError) ErrorName() string {
	return "CreatePostsV1ResponseValidationError"
}

// Error satisfies the builtin error interface
func (e CreatePostsV1ResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreatePostsV1Response.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreatePostsV1ResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreatePostsV1ResponseValidationError{}

// Validate checks the field values on GenerateFeedV1Request with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *GenerateFeedV1Request) Validate() error {
	if m == nil {
		return nil
	}

	if val := m.GetLimit(); val < 3 || val > 27 {
		return GenerateFeedV1RequestValidationError{
			field:  "Limit",
			reason: "value must be inside range [3, 27]",
		}
	}

	// no validation rules for Offset

	return nil
}

// GenerateFeedV1RequestValidationError is the validation error returned by
// GenerateFeedV1Request.Validate if the designated constraints aren't met.
type GenerateFeedV1RequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GenerateFeedV1RequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GenerateFeedV1RequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GenerateFeedV1RequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GenerateFeedV1RequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GenerateFeedV1RequestValidationError) ErrorName() string {
	return "GenerateFeedV1RequestValidationError"
}

// Error satisfies the builtin error interface
func (e GenerateFeedV1RequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGenerateFeedV1Request.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GenerateFeedV1RequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GenerateFeedV1RequestValidationError{}

// Validate checks the field values on GenerateFeedV1Response with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *GenerateFeedV1Response) Validate() error {
	if m == nil {
		return nil
	}

	if v, ok := interface{}(m.GetFeed()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return GenerateFeedV1ResponseValidationError{
				field:  "Feed",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	return nil
}

// GenerateFeedV1ResponseValidationError is the validation error returned by
// GenerateFeedV1Response.Validate if the designated constraints aren't met.
type GenerateFeedV1ResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GenerateFeedV1ResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GenerateFeedV1ResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GenerateFeedV1ResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GenerateFeedV1ResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GenerateFeedV1ResponseValidationError) ErrorName() string {
	return "GenerateFeedV1ResponseValidationError"
}

// Error satisfies the builtin error interface
func (e GenerateFeedV1ResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGenerateFeedV1Response.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GenerateFeedV1ResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GenerateFeedV1ResponseValidationError{}

// Validate checks the field values on Post with the rules defined in the proto
// definition for this message. If any rules are violated, an error is returned.
func (m *Post) Validate() error {
	if m == nil {
		return nil
	}

	return nil
}

// PostValidationError is the validation error returned by Post.Validate if the
// designated constraints aren't met.
type PostValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PostValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PostValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PostValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PostValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PostValidationError) ErrorName() string { return "PostValidationError" }

// Error satisfies the builtin error interface
func (e PostValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPost.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PostValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PostValidationError{}

// Validate checks the field values on Feed with the rules defined in the proto
// definition for this message. If any rules are violated, an error is returned.
func (m *Feed) Validate() error {
	if m == nil {
		return nil
	}

	return nil
}

// FeedValidationError is the validation error returned by Feed.Validate if the
// designated constraints aren't met.
type FeedValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e FeedValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e FeedValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e FeedValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e FeedValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e FeedValidationError) ErrorName() string { return "FeedValidationError" }

// Error satisfies the builtin error interface
func (e FeedValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sFeed.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = FeedValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = FeedValidationError{}
