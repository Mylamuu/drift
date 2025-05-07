package main

import (
	"flag"
	"strings"
)

type Config struct {
	Port             int
	BindAddress      string
	LogLevel         string
	StoragePath      string
	MaxFileSizeMB    int64
	AllowedFileTypes []string
}

type stringSliceFlag []string

func (s *stringSliceFlag) String() string {
	return strings.Join(*s, ",")
}

func (s *stringSliceFlag) Set(value string) error {
	if value == "" {
		*s = []string{}
		return nil
	}
	*s = strings.Split(value, ",")
	return nil
}

func loadConfig(args []string) (Config, error) {
	fs := flag.NewFlagSet("drift", flag.ContinueOnError)

	port := fs.Int("port", 8080, "Port to listen on")
	bindAddress := fs.String("bind", "localhost", "Address to bind to")
	logLevel := fs.String("level", "info", "Logging level: [debug, info, warn, error]")
	storagePath := fs.String("path", "/tmp/drift", "Path to save uploads to")
	maxFileSizeMB := fs.Int64("size", 1024, "Maximum file size in MB")

	var allowedFileTypes stringSliceFlag
	fs.Var(&allowedFileTypes, "allowed", "Comma seperated list of allowed MIME types")

	if err := fs.Parse(args); err != nil {
		return Config{}, err
	}

	if allowedFileTypes == nil {
		allowedFileTypes = []string{}
	}

	return Config{
		Port:             *port,
		BindAddress:      *bindAddress,
		LogLevel:         *logLevel,
		StoragePath:      *storagePath,
		MaxFileSizeMB:    *maxFileSizeMB,
		AllowedFileTypes: allowedFileTypes,
	}, nil
}
