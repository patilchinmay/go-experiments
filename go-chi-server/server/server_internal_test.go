package server

import (
	"os"
	"time"

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
			NonDefaultHost     = "192.168.10.10"
			NonDefaultPort     = "9000"
			DefaultHost        = "0.0.0.0"
			DefaultPort        = "8080"
			DefaultReadTimeout = 5 * time.Second
		)

		It("host non-default, port default", func() {
			ts := New().WithLogger(zerolog.Nop()).WithHost(NonDefaultHost)
			Expect(ts.Host).To(Equal(NonDefaultHost))
			Expect(ts.Port).To(Equal(DefaultPort))
			ts = nil
		})

		It("host default, port non-default", func() {
			ts := New().WithLogger(zerolog.Nop()).WithPort(NonDefaultPort)
			Expect(ts.Port).To(Equal(NonDefaultPort))
			Expect(ts.Host).To(Equal(DefaultHost))
			ts = nil
		})

		It("host non-default, port non-default", func() {
			ts := New().WithLogger(zerolog.Nop()).WithHost(NonDefaultHost).WithPort(NonDefaultPort)
			Expect(ts.Port).To(Equal(NonDefaultPort))
			Expect(ts.Host).To(Equal(NonDefaultHost))
			ts = nil
		})

		It("host non-default, port non-default, with env var", func() {
			os.Setenv("HOST", NonDefaultHost)
			defer os.Unsetenv("HOST")

			os.Setenv("PORT", NonDefaultPort)
			defer os.Unsetenv("PORT")

			ts := New().WithLogger(zerolog.Nop())
			Expect(ts.Port).To(Equal(NonDefaultPort))
			Expect(ts.Host).To(Equal(NonDefaultHost))
			ts = nil
		})

		It("host default, port default", func() {
			ts := New().WithLogger(zerolog.Nop())
			Expect(ts.Port).To(Equal(DefaultPort))
			Expect(ts.Host).To(Equal(DefaultHost))
			ts = nil
		})

		It("Default ReadTimeout", func() {
			ts := New().WithLogger(zerolog.Nop())
			Expect(ts.server.ReadTimeout).To(Equal(DefaultReadTimeout))
			ts = nil
		})

		It("Custom ReadTimeout", func() {
			ts := New().WithLogger(zerolog.Nop()).WithReadTimeout(10 * time.Second)
			Expect(ts.server.ReadTimeout).To(Equal(10 * time.Second))
			ts = nil
		})

		It("Default ReadTimeout from env var", func() {
			os.Setenv("READ_TIMEOUT", "15s")
			defer os.Unsetenv("READ_TIMEOUT")

			ts := New().WithLogger(zerolog.Nop())
			Expect(ts.server.ReadTimeout).To(Equal(15 * time.Second))
			ts = nil
		})

		It("Custom ReadTimeout should be preferred in the presence of env var", func() {
			os.Setenv("READ_TIMEOUT", "15s")
			defer os.Unsetenv("READ_TIMEOUT")

			ts := New().WithLogger(zerolog.Nop()).WithReadTimeout(10 * time.Second)
			Expect(ts.server.ReadTimeout).To(Equal(10 * time.Second))
			ts = nil
		})

	})
})
