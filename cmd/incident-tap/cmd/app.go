package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	stdlog "log"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	_ "embed"

	"github.com/alecthomas/kingpin/v2"
	kitlog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/incident-io/singer-tap/client"
	"github.com/incident-io/singer-tap/config"
	"github.com/incident-io/singer-tap/tap"
	"github.com/pkg/errors"
)

var logger kitlog.Logger

var (
	app = kingpin.New("incident-tap", "Extract data from incident.io for use with Singer").Version(versionStanza())

	// Global flags
	debug       = app.Flag("debug", "Enable debug logging").Default("false").Bool()
	configFile  = app.Flag("config", "Configuration file").ExistingFile()
	catalogFile = app.Flag("catalog", "If set, allows filtering which streams would be synced").ExistingFile()
)

func Run(ctx context.Context) (err error) {
	_, err = app.Parse(os.Args[1:])
	if err != nil {
		return err
	}

	logger = kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr))
	if *debug {
		logger = level.NewFilter(logger, level.AllowDebug())
	} else {
		logger = level.NewFilter(logger, level.AllowInfo())
	}
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC, "caller", kitlog.DefaultCaller)
	stdlog.SetOutput(kitlog.NewStdlibAdapter(logger))

	// Root context to the application.
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling.
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	go func() {
		<-sigc
		cancel()
		<-sigc
		panic("received second signal, exiting immediately")
	}()

	cfg, err := loadConfigOrError(ctx, *configFile)
	if err != nil {
		return err
	}

	if cfg.Endpoint == "" {
		cfg.Endpoint = "https://api.incident.io"
	}
	cl, err := client.New(ctx, cfg.APIKey, cfg.Endpoint, Version())
	if err != nil {
		return err
	}

	// Singer requires taps to output to STDOUT. We log to STDERR so the debug log output
	// can be streamed separately.
	ol := tap.NewOutputLogger(os.Stdout)

	err = tap.Run(ctx, logger, ol, cl)
	if err != nil {
		return err
	}

	return nil
}

// Set via compiler flags
var (
	Commit    = "none"
	Date      = "unknown"
	GoVersion = runtime.Version()
)

//go:embed VERSION
var version string

func Version() string {
	return strings.TrimSpace(version)
}

func versionStanza() string {
	return fmt.Sprintf(
		"Version: %v\nGit SHA: %v\nGo Version: %v\nGo OS/Arch: %v/%v\nBuilt at: %v",
		Version(), Commit, GoVersion, runtime.GOOS, runtime.GOARCH, Date,
	)
}

func loadConfigOrError(ctx context.Context, configFile string) (cfg *config.Config, err error) {
	defer func() {
		if err == nil {
			return
		}
		if configFile == "" {
			OUT("No config file (--config) was provided, but is required.\n")
		} else {
			OUT("Failed to load config file!\n")
		}

		OUT(`We expect a config file in JSON format that looks like:
{
  "api_key": "<your-api-key>",
}
`)
	}()

	if configFile == "" {
		return nil, errors.New("No config file set! (--config)")
	}

	cfg, err = config.FileLoader(configFile).Load(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "loading config")
	}
	if err := cfg.Validate(); err != nil {
		data, _ := json.MarshalIndent(err, "", "  ")

		// Print the validation error in JSON. Needs improving.
		return nil, fmt.Errorf("validating config:\n%s", string(data))
	}

	return cfg, nil
}

// OUT prints progress output to stderr.
func OUT(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
}
