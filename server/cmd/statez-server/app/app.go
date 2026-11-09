package app

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"

	"server/internal/httpserver"
	"server/internal/mod"
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

		hostPort   string
		modDataDir string
		apiToken   string
	)

	// attach flags to the flag set
	flags.StringVar(&hostPort, "listen", "0.0.0.0:8080", "Host and port to run the server on")
	flags.StringVar(&modDataDir, "mod-data-dir", "", "Path to the directory where the StateZ mod stores its data")
	flags.StringVar(&apiToken, "api-token", "",
		"Bearer token required to access the data endpoints (empty value disables authentication)")

	// parse the command-line arguments
	if err := flags.Parse(args); err != nil {
		return err
	}

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

	// open a TCP listener on the specified host and port
	ln, lnErr := (&net.ListenConfig{}).Listen(ctx, "tcp", net.JoinHostPort(host, port))
	if lnErr != nil {
		return fmt.Errorf("listening http port: %v", lnErr)
	}
	defer func() { _ = ln.Close() }()

	log := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	dd := mod.OpenDataDir(modData.FS())

	return httpserver.New(httpserver.NewHandler(log, dd, apiToken)).Serve(ctx, ln)
}
