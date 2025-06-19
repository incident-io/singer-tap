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
	app = kingpin.New("tap-incident", "Extract data from incident.io for use with Singer").Version(versionStanza())

	// Global flags
	debug         = app.Flag("debug", "Enable debug logging").Default("false").Bool()
	configFile    = app.Flag("config", "Configuration file").ExistingFile()
	catalogFile   = app.Flag("catalog", "If set, allows filtering which streams would be synced").ExistingFile()
	discoveryMode = app.Flag("discover", "If set, only outputs the catalog and exits").Default("false").Bool()
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

	if *discoveryMode {
		err = tap.Discover(ctx, logger, ol)
		if err != nil {
			return err
		}
	} else {
		// If we're syncing - check if we were given a catalog
		var (
			catalog *tap.Catalog
			err     error
		)

		if *catalogFile != "" {
			catalog, err = loadCatalogOrError(ctx, *catalogFile)
			if err != nil {
				return err
			}
		}

		err = tap.Sync(ctx, logger, ol, cl, catalog)
		if err != nil {
			return err
		}
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

func loadCatalogOrError(ctx context.Context, catalogFile string) (catalog *tap.Catalog, err error) {
	defer func() {
		if err == nil {
			return
		}
		OUT("Failed to load catalog file!\n")
	}()

	catalog, err = config.LoadAndParse(catalogFile, tap.Catalog{})
	if err != nil {
		return nil, errors.Wrap(err, "loading catalog")
	}

	return catalog, nil
}

func loadConfigOrError(ctx context.Context, configFile string) (cfg *config.Config, err error) {
	defer func() {
		if err == nil {
			return
		}
		if configFile == "" {
			OUT("No config file (--config) was provided.\n")
		} else {
			OUT("Failed to load config file!\n")
		}

		OUT(`You can provide configuration via:

1. Environment variables:
   export INCIDENT_API_KEY="<your-api-key>"
   export INCIDENT_ENDPOINT="https://api.incident.io" (optional)

2. Config file in JSON format:
   {
     "api_key": "<your-api-key>",
     "endpoint": "<api-endpoint>" (optional)
   }
`)
	}()

	// Try to load from environment variables first
	cfg = &config.Config{
		APIKey:   os.Getenv("INCIDENT_API_KEY"),
		Endpoint: os.Getenv("INCIDENT_ENDPOINT"),
	}

	// If a config file is provided, load it and merge with env vars
	if configFile != "" {
		fileCfg, err := config.LoadAndParse(configFile, config.Config{})
		if err != nil {
			return nil, errors.Wrap(err, "loading config")
		}
		
		// File config takes precedence over env vars
		if fileCfg.APIKey != "" {
			cfg.APIKey = fileCfg.APIKey
		}
		if fileCfg.Endpoint != "" {
			cfg.Endpoint = fileCfg.Endpoint
		}
	}

	// Validate the final config
	if err := cfg.Validate(); err != nil {
		if configFile == "" && cfg.APIKey == "" {
			return nil, errors.New("No API key provided. Set INCIDENT_API_KEY environment variable or use --config flag")
		}
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
