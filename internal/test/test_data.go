package test

import (
	"embed"

	"github.com/arttet/reddit-feed-api/internal/model"

	"go.uber.org/zap"

	"gopkg.in/yaml.v3"
)

var (
	//go:embed data
	resources embed.FS
)

type TestData struct {
	Posts         model.Posts `yaml:"posts"`
	WrongPosts    model.Posts `yaml:"wrongPosts"`
	PromotedPosts model.Posts `yaml:"promotedPosts"`
}

func LoadTestData(filename string) *TestData {
	file, err := resources.ReadFile(filename)
	if err != nil {
		zap.L().Fatal("failed to read the file",
			zap.String("file", filename),
			zap.Error(err),
		)
		return nil
	}

	var data TestData
	err = yaml.Unmarshal(file, &data)
	if err != nil {
		zap.L().Fatal("failed to load test data",
			zap.String("file", filename),
			zap.Error(err),
		)
		return nil
	}

	return &data
}
