package config

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	DBPath          string
	MasterPassword  string
	ServerAddr      string
	SyncIntervalSec int
}

func NewConfigFromArgs(args []string) (*Config, error) {
	fs := flag.NewFlagSet("client", flag.ContinueOnError)

	cfg := &Config{}
	fs.StringVar(&cfg.ServerAddr, "addr", "localhost:50501", "server address (host:port)")
	fs.IntVar(&cfg.SyncIntervalSec, "interval", 60, "synchronization interval in seconds")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "usage: %s [options] <db_path> <master_password>\n", os.Args[0])
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return nil, fmt.Errorf("failed to parse flags: %w", err)
	}

	remainingArgs := fs.Args()
	if len(remainingArgs) != 2 {
		fs.Usage()
		return nil, fmt.Errorf("invalid arguments")
	}

	cfg.DBPath = remainingArgs[0]
	cfg.MasterPassword = remainingArgs[1]

	return cfg, nil
}

func NewConfig() (*Config, error) {
	return NewConfigFromArgs(os.Args[1:])
}
