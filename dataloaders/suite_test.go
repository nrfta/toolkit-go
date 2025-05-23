package dataloaders_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestDataloaders(t *testing.T) {
	gomega.RegisterFailHandler(Fail)
	RunSpecs(t, "Dataloaders Suite")
}
