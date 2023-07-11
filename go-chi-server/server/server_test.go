package server_test

import (
	"net"
	"strconv"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/patilchinmay/go-experiments/go-chi-server/server"
	"github.com/phayes/freeport"
	"github.com/rs/zerolog"
)

var _ = Describe("Server", Serial, func() {
	var ts = &server.Server{}
	host := "127.0.0.1"
	var portInt, _ = freeport.GetFreePort()
	var port = strconv.Itoa(portInt)

	BeforeEach(func() {
		// Create server
		ts = server.New().WithLogger(zerolog.Nop()).WithHost(host).WithPort(port)
		go ts.Serve()
	})

	AfterEach(func() {
		ts.Shutdown()
		ts = nil
	})

	Context("Server", func() {

		It("Server started and is reachable", func() {
			timeout := 1 * time.Second
			addr := host + ":" + port

			conn, err := net.DialTimeout("tcp", addr, timeout)
			Expect(err).ShouldNot(HaveOccurred())
			defer conn.Close()
		})
	})
})
