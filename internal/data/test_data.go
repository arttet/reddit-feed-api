package data

import (
	"io/ioutil"

	"github.com/arttet/reddit-feed-api/internal/model"

	"github.com/rs/zerolog/log"

	"gopkg.in/yaml.v2"
)

type TestData struct {
	Posts         []model.Post `yaml:"posts"`
	WrongPosts    []model.Post `yaml:"wrongPosts"`
	PromotedPosts []model.Post `yaml:"promotedPosts"`
}

func LoadTestData(filename string) *TestData {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to read the file %s", filename)
		return nil
	}

	var data TestData
	err = yaml.Unmarshal(file, &data)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to load test data from the file %s", filename)
		return nil
	}

	return &data
}
