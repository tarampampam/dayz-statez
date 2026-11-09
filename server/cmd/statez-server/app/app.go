package app

import (
	"context"
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
	envListen     evnVar = "STATEZ_SERVER_LISTEN"
	envModDataDir evnVar = "STATEZ_SERVER_MOD_DATA_DIR"
	envAPIToken   evnVar = "STATEZ_SERVER_API_TOKEN"
)

type App struct {
	name string
}

func NewApp(name string) *App {
	return &App{name: name}
}

func (a *App) Run(ctx context.Context, args []string) error {
	var (
		flags = flag.NewFlagSet(a.name, flag.ExitOnError)
		log   = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))

		hostPort   string
		modDataDir string
		apiToken   string
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
		"",
		"Path to the directory where the StateZ mod stores its data [$"+envModDataDir.String()+"]",
	)
	flags.StringVar(
		&apiToken,
		"api-token",
		"",
		"Bearer token required to access the data endpoints (empty value disables authentication) "+
			"[$"+envAPIToken.String()+"]",
	)

	// parse the command-line arguments
	if err := flags.Parse(args); err != nil {
		return err
	}

	// override flag values with environment variables if they exist
	envListen.OverrideIfExists(&hostPort)
	envModDataDir.OverrideIfExists(&modDataDir)
	envAPIToken.OverrideIfExists(&apiToken)

	// validate the provided host and port
	host, port, hostPortErr := net.SplitHostPort(hostPort)
	if hostPortErr != nil {
		return fmt.Errorf("invalid listen address: %v", hostPortErr)
	}

	// validate the provided mod data directory
	if stat, statErr := os.Stat(modDataDir); statErr != nil || !stat.IsDir() {
		return fmt.Errorf("invalid mod data directory: %v", statErr)
	}

	// open the mod data directory for reading
	modData, fsErr := os.OpenRoot(modDataDir)
	if fsErr != nil {
		return fmt.Errorf("opening mod data directory: %v", fsErr)
	}
	defer func() { _ = modData.Close() }()

	log.Log(ctx, slog.LevelInfo, "Data directory opened successfully", slog.String("path", modDataDir))

	// open a TCP listener on the specified host and port
	ln, lnErr := (&net.ListenConfig{}).Listen(ctx, "tcp", net.JoinHostPort(host, port))
	if lnErr != nil {
		return fmt.Errorf("listening http port: %v", lnErr)
	}
	defer func() { _ = ln.Close() }()

	log.Log(ctx, slog.LevelInfo, "TCP port opened successfully", slog.String("address", ln.Addr().String()))

	defer func(now time.Time) {
		log.Log(ctx, slog.LevelInfo, "HTTP Server has stopped", slog.String("uptime", time.Since(now).String()))
	}(time.Now())

	log.Log(ctx, slog.LevelInfo, "HTTP Server has started")

	return httpserver.New(
		httpserver.NewHandler(
			log,
			mod.OpenDataDir(modData.FS()),
			apiToken,
		),
	).Serve(ctx, ln)
}
