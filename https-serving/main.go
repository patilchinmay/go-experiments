package main

import (
	"context"
	"crypto/tls"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/httplog"
	"github.com/joho/godotenv"
	"github.com/patilchinmay/go-experiments/https-serving/app"
	"github.com/patilchinmay/go-experiments/https-serving/server"
)

func main() {
	// If we crash the go code, we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Create context that listens for the interrupt signal from the OS.
	// We do this at the beginning of the program as
	// We should be capable of catching signal as soon as the program starts
	// https://henvic.dev/posts/signal-notify-context/
	// https://millhouse.dev/posts/graceful-shutdowns-in-golang-with-signal-notify-context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Create logger
	logger := httplog.NewLogger("My App")

	// Load env variables.
	err := godotenv.Load()
	if err != nil {
		logger.Error().Msg(".env file is not found")
	}

	// Set app env
	env := os.Getenv("ENV")
	logger.Info().Str("ENV", env).Msg("")

	// Logger level setting
	var jsonLogs bool
	var logLevel string
	if env == "local" {
		jsonLogs = false
		logLevel = "debug"
	} else {
		jsonLogs = true
		logLevel = "info"
	}

	// Redefine Logger with proper cconfig
	logger = httplog.NewLogger("My App", httplog.Options{
		JSON:     jsonLogs,
		Concise:  true,
		LogLevel: logLevel,
		Tags: map[string]string{
			"env": env,
		},
	})

	// key and certificate for tls
	var (
		tlscert, tlskey string
	)
	flag.StringVar(&tlscert, "tlscert", "./certs/server.crt", "File containing the x509 Certificate for HTTPS. This is a public key.")
	flag.StringVar(&tlskey, "tlskey", "./certs/server.key", "File containing the x509 private key to --tlscert. This is a private key.")

	flag.Parse()

	certs, err := tls.LoadX509KeyPair(tlscert, tlskey)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to load key pair")
	}
	// Create app with routes handlers
	app := app.New().WithLogger(logger).CreateApp()

	// Create server
	server := server.New().WithLogger(logger).WithReadTimeout(5 * time.Second).WithHandlers(app).WithTLS(certs)

	// The server is started on a separate goroutine as
	// ListenAndServe is a blocking function,
	// It fields all incoming requests on separate goroutine
	// https://medium.com/honestbee-tw-engineer/gracefully-shutdown-in-go-http-server-5f5e6b83da5a
	go server.ServeTLS()

	// Block main goroutine to wait for receiving os interrupt signal
	<-ctx.Done()

	// Received interrupt signal. Proceed to graceful shutdown
	logger.Info().Msg("Received interrupt, shutting down gracefully")
	server.Shutdown()
}
