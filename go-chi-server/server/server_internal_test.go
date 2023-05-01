package server

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog"
)

// the package is server instead of server_test
// so that we can test internal methods of server
// this file will be run automatically by ginkgo command
var _ = Describe("Server Internal Tests", func() {

	Context("Server Parameters", func() {

		It("host non-default, port default", func() {
			ts := New().WithLogger(zerolog.Nop()).WithHost("127.0.0.1")
			Expect(ts.host).To(Equal("127.0.0.1"))
			Expect(ts.port).To(Equal("8080"))
			ts = nil
		})

		It("host default, port non-default", func() {
			ts := New().WithLogger(zerolog.Nop()).WithPort("9000")
			Expect(ts.port).To(Equal("9000"))
			Expect(ts.host).To(Equal("0.0.0.0"))
			ts = nil
		})

		It("host non-default, port non-default", func() {
			ts := New().WithLogger(zerolog.Nop()).WithHost("127.0.0.1").WithPort("9000")
			Expect(ts.port).To(Equal("9000"))
			Expect(ts.host).To(Equal("127.0.0.1"))
			ts = nil
		})

		It("host default, port default", func() {
			ts := New().WithLogger(zerolog.Nop())
			Expect(ts.port).To(Equal("8080"))
			Expect(ts.host).To(Equal("0.0.0.0"))
			ts = nil
		})
	})
})
