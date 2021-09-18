package api_test

import (
	"testing"

	"github.com/arttet/reddit-feed-api/internal/data"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var testData *data.TestData

func TestAPI(t *testing.T) {
	RegisterFailHandler(Fail)

	testData = data.LoadTestData("../data/data/posts.yaml")
	Expect(testData).ShouldNot(BeNil())

	RunSpecs(t, "API Suite")
}
