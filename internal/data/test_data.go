package data

import (
	"io/ioutil"

	"github.com/arttet/reddit-feed-api/internal/model"

	"go.uber.org/zap"

	"gopkg.in/yaml.v3"
)

type TestData struct {
	Posts         []model.Post `yaml:"posts"`
	WrongPosts    []model.Post `yaml:"wrongPosts"`
	PromotedPosts []model.Post `yaml:"promotedPosts"`
}

func LoadTestData(filename string) *TestData {
	file, err := ioutil.ReadFile(filename)
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
