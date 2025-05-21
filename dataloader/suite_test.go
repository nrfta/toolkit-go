package dataloader_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestDataloader(t *testing.T) {
	gomega.RegisterFailHandler(Fail)
	RunSpecs(t, "Dataloader Suite")
}
