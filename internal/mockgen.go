package internal

//go:generate mockgen -destination=./mock/producer_mock.go -package=mock github.com/arttet/reddit-feed-api/internal/app/reddit-feed-api/broker Producer

//go:generate mockgen -destination=./mock/repo_mock.go -package=mock github.com/arttet/reddit-feed-api/internal/app/reddit-feed-api/service/repository Repository
