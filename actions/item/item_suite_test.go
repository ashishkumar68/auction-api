package item

import (
	"github.com/ashishkumar68/auction-api/config"
	"github.com/ashishkumar68/auction-api/database"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"testing"
)

func TestItem(t *testing.T) {

	RegisterFailHandler(Fail)
	RunSpecs(t, "Item Suite")
}

var _ = BeforeSuite(func() {
	config.LoadDBConfig()
	database.InitialiseDatabase()
})
