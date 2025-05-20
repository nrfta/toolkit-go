package comparators_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestComparators(t *testing.T) {
	gomega.RegisterFailHandler(Fail)
	RunSpecs(t, "Comparators Suite")
}
