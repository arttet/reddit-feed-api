package api_test

import (
	"testing"

	"github.com/arttet/reddit-feed-api/internal/data"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var testDataPost *data.TestData

func TestAPI(t *testing.T) {
	RegisterFailHandler(Fail)

	testDataPost = data.LoadTestData("../data/data/posts.yaml")
	Expect(testDataPost).ShouldNot(BeNil())

	RunSpecs(t, "API Suite")
}
