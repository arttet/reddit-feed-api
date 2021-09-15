package internal

//go:generate mockgen -destination=./reddit-feed-api/mock/repo_mock.go -package=mock github.com/arttet/reddit-feed-api/internal/reddit-feed-api/repo Repo
