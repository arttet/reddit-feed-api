package repo_test

import (
	"io/ioutil"
	"testing"

	"github.com/arttet/reddit-feed-api/internal/reddit-feed-api/model"

	"gopkg.in/yaml.v2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var testData struct {
	Posts []model.Post `yaml:"posts"`
}

func TestRepo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Repo Suite")
}

var _ = BeforeSuite(func() {
	file, err := ioutil.ReadFile("../share/test_data/posts.yaml")
	Expect(err).Should(BeNil())

	err = yaml.Unmarshal(file, &testData)
	Expect(err).Should(BeNil())
})

var _ = AfterSuite(func() {
})
