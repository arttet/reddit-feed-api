package internal

// https://github.com/golang/mock

//go:generate mockgen -destination=./mock/repo_mock.go -package=mock github.com/arttet/reddit-feed-api/internal/repo Repo
