package internal

// https://github.com/golang/mock

//go:generate mockgen -destination=./mock/span_mock.go -package=mock github.com/opentracing/opentracing-go Span
