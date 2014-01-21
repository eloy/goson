package goson_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGoson(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Goson Suite")
}
