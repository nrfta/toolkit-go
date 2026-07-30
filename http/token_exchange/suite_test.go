package token_exchange_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestTokenExchange(t *testing.T) {
	gomega.RegisterFailHandler(Fail)
	RunSpecs(t, "Token Exchange Suite")
}
