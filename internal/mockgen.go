package internal

// https://github.com/golang/mock

//go:generate mockgen -destination=./mock/producer_mock.go -package=mock github.com/arttet/reddit-feed-api/internal/broker Producer
