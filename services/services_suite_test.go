package services

import (
	"github.com/ashishkumar68/auction-api/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"testing"
)

func TestIndex(t *testing.T) {
	config.InitialiseConfig()

	RegisterFailHandler(Fail)
	RunSpecs(t, "Services Suite")
}
