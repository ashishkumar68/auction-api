package actions

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"io"
	"net/http"
	"os"
)

var _ = Describe("Index Tests", func() {
	protocol := "http"
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	prefix := "/"

	Context("The service should be available.", func() {

		When("I request Index route", func() {
			It("should find service available", func() {
				resp, err := http.Get(fmt.Sprintf("%s://%s:%s%s", protocol, host, port, prefix))
				fmt.Println(fmt.Sprintf("%s://%s:%s%s", protocol, host, port, prefix))
				Expect(err).To(BeNil(), "Could not detect service available.")
				Expect(resp).To(Not(BeNil()), "Could not detect service available.")
				defer resp.Body.Close()

				body, err := io.ReadAll(resp.Body)
				Expect(err).To(BeNil(), "Could not read response from HTTP message")
				Expect(string(body)).To(Equal("{}"))
			})
		})
	})
})