package server

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog"
)

// the package is server instead of server_test
// so that we can test internal methods of server
// this file will be run automatically by ginkgo command
var _ = Describe("Server Internal Tests", func() {

	Context("Server Parameters", func() {
		const (
			NonDefaultHost = "192.168.10.10"
			NonDefaultPort = "9000"
		)

		It("host non-default, port default", func() {
			ts := New().WithLogger(zerolog.Nop()).WithHost(NonDefaultHost)
			Expect(ts.host).To(Equal(NonDefaultHost))
			Expect(ts.port).To(Equal(DefaultPort))
			ts = nil
		})

		It("host default, port non-default", func() {
			ts := New().WithLogger(zerolog.Nop()).WithPort(NonDefaultPort)
			Expect(ts.port).To(Equal(NonDefaultPort))
			Expect(ts.host).To(Equal(DefaultHost))
			ts = nil
		})

		It("host non-default, port non-default", func() {
			ts := New().WithLogger(zerolog.Nop()).WithHost(NonDefaultHost).WithPort(NonDefaultPort)
			Expect(ts.port).To(Equal(NonDefaultPort))
			Expect(ts.host).To(Equal(NonDefaultHost))
			ts = nil
		})

		It("host non-default, port non-default, with env var", func() {
			os.Setenv("HOST", NonDefaultHost)
			defer os.Unsetenv("HOST")

			os.Setenv("PORT", NonDefaultPort)
			defer os.Unsetenv("PORT")

			ts := New().WithLogger(zerolog.Nop())
			Expect(ts.port).To(Equal(NonDefaultPort))
			Expect(ts.host).To(Equal(NonDefaultHost))
			ts = nil
		})

		It("host default, port default", func() {
			ts := New().WithLogger(zerolog.Nop())
			Expect(ts.port).To(Equal(DefaultPort))
			Expect(ts.host).To(Equal(DefaultHost))
			ts = nil
		})
	})
})
