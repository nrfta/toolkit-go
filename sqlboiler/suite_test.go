package sqlboiler_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestSQLBoiler(t *testing.T) {
	gomega.RegisterFailHandler(Fail)
	RunSpecs(t, "SQLBoiler Suite")
}
