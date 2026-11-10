package app

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"
	"time"

	"server/internal/httpserver"
	"server/internal/mod"
)

type evnVar string

// String returns the string representation of the environment variable.
func (e evnVar) String() string { return string(e) }

// OverrideIfExists overrides the value of the provided destination string pointer with the value of the
// environment variable if it exists.
func (e evnVar) OverrideIfExists(dest *string) {
	if v, ok := os.LookupEnv(string(e)); ok {
		*dest = v
	}
}

const (
	envListen      evnVar = "STATEZ_SERVER_LISTEN"
	envModDataDir  evnVar = "STATEZ_SERVER_MOD_DATA_DIR"
	envAPIToken    evnVar = "STATEZ_SERVER_API_TOKEN" //nolint:gosec // this is not a secret, it's just an env name
	envTLSCertFile evnVar = "STATEZ_SERVER_TLS_CERT"
	envTLSKeyFile  evnVar = "STATEZ_SERVER_TLS_KEY"
)

// App is the main CLI application struct.
type App struct {
	name string
}

// NewApp creates a new instance of the App with the given name.
func NewApp(name string) *App {
	return &App{name: name}
}

// Run executes the application with the provided context and command-line arguments.
func (a *App) Run(ctx context.Context, args []string) error { //nolint:funlen
	var (
		flags = flag.NewFlagSet(a.name, flag.ExitOnError)
		log   = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))

		hostPort    string
		modDataDir  string
		apiToken    string
		tlsCertFile string
		tlsKeyFile  string
	)

	// attach flags to the flag set
	flags.StringVar(
		&hostPort,
		"listen",
		"0.0.0.0:8080",
		"Host and port to run the server on [$"+envListen.String()+"]",
	)
	flags.StringVar(
		&modDataDir,
		"mod-data-dir",
		"/nonexistent",
		"REQUIRED. Path to the directory where the StateZ mod stores its data [$"+envModDataDir.String()+"]",
	)
	flags.StringVar(
		&apiToken,
		"api-token",
		"",
		"Bearer token required to access the data endpoints (empty value disables authentication) "+
			"[$"+envAPIToken.String()+"]",
	)
	flags.StringVar(
		&tlsCertFile,
		"tls-cert",
		"",
		"Path to the TLS certificate file for HTTPS (enables HTTPS if provided) [$"+envTLSCertFile.String()+"]",
	)
	flags.StringVar(
		&tlsKeyFile,
		"tls-key",
		"",
		"Path to the TLS private key file [$"+envTLSKeyFile.String()+"]",
	)

	// parse the command-line arguments
	if err := flags.Parse(args); err != nil {
		return err
	}

	// override flag values with environment variables if they exist
	envListen.OverrideIfExists(&hostPort)
	envModDataDir.OverrideIfExists(&modDataDir)
	envAPIToken.OverrideIfExists(&apiToken)
	envTLSCertFile.OverrideIfExists(&tlsCertFile)
	envTLSKeyFile.OverrideIfExists(&tlsKeyFile)

	// validate the provided host and port
	host, port, hostPortErr := net.SplitHostPort(hostPort)
	if hostPortErr != nil {
		return fmt.Errorf("invalid listen address: %w", hostPortErr)
	}

	// validate the provided mod data directory
	if stat, statErr := os.Stat(modDataDir); statErr != nil || !stat.IsDir() {
		return fmt.Errorf("invalid mod data directory: %w", statErr)
	}

	var tlsKeyPair *tls.Certificate // nil if TLS is not enabled

	// validate and load the TLS certificate and key files if provided
	if tlsCertFile != "" && tlsKeyFile != "" {
		if stat, err := os.Stat(tlsCertFile); err != nil || !stat.Mode().IsRegular() {
			return fmt.Errorf("invalid TLS certificate file: %w", err)
		}

		if stat, err := os.Stat(tlsKeyFile); err != nil || !stat.Mode().IsRegular() {
			return fmt.Errorf("invalid TLS key file: %w", err)
		}

		tlsCert, certErr := tls.LoadX509KeyPair(tlsCertFile, tlsKeyFile)
		if certErr != nil {
			return fmt.Errorf("load TLS cert and key: %w", certErr)
		}

		tlsKeyPair = &tlsCert

		log.Log(ctx, slog.LevelInfo, "TLS configuration loaded successfully",
			slog.String("cert_file", tlsCertFile),
			slog.String("key_file", tlsKeyFile),
		)
	}

	// open the mod data directory for reading
	modData, fsErr := os.OpenRoot(modDataDir)
	if fsErr != nil {
		return fmt.Errorf("opening mod data directory: %w", fsErr)
	}
	defer func() { _ = modData.Close() }()

	log.Log(ctx, slog.LevelInfo, "Data directory opened successfully", slog.String("path", modDataDir))

	// open a TCP listener on the specified host and port
	ln, lnErr := (&net.ListenConfig{}).Listen(ctx, "tcp", net.JoinHostPort(host, port))
	if lnErr != nil {
		return fmt.Errorf("listening http port: %w", lnErr)
	}
	defer func() { _ = ln.Close() }()

	log.Log(ctx, slog.LevelInfo, "TCP port opened successfully", slog.String("address", ln.Addr().String()))

	defer func(now time.Time) {
		log.Log(ctx, slog.LevelInfo, "HTTP Server has stopped", slog.String("uptime", time.Since(now).String()))
	}(time.Now())

	if apiToken == "" {
		log.Log(ctx, slog.LevelWarn, "API token is empty, authentication is DISABLED")
	}

	opts := []httpserver.Option{
		httpserver.WithErrorLog(slog.NewLogLogger(log.Handler(), slog.LevelError)),
	}

	// if TLS is enabled, configure the HTTP server to use the provided TLS certificate and key
	if tlsKeyPair != nil {
		opts = append(opts, httpserver.WithTLSConfig(&tls.Config{Certificates: []tls.Certificate{*tlsKeyPair}}))
	}

	log.Log(ctx, slog.LevelInfo, "HTTP Server has started")

	return httpserver.New(
		httpserver.NewHandler(
			log,
			mod.OpenDataDir(modData.FS()),
			apiToken,
		),
		opts...,
	).Serve(ctx, ln)
}
