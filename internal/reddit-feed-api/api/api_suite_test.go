package api_test

import (
	"io/ioutil"
	"testing"

	"github.com/arttet/reddit-feed-api/internal/reddit-feed-api/model"

	"gopkg.in/yaml.v2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var testData struct {
	Posts      []model.Post `yaml:"posts"`
	WrongPosts []model.Post `yaml:"wrongPosts"`
}

func TestAPI(t *testing.T) {
	RegisterFailHandler(Fail)

	file, err := ioutil.ReadFile("../share/test_data/posts.yaml")
	Expect(err).Should(BeNil())

	err = yaml.Unmarshal(file, &testData)
	Expect(err).Should(BeNil())

	RunSpecs(t, "API Suite")
}
