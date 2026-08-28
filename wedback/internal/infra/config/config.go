package config

import (
	"encoding/json"
	"errors"
	"io/fs"
	"log"
	"strings"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"github.com/parham-alvani/wedding/wedback/internal/infra/db"
	"github.com/parham-alvani/wedding/wedback/internal/infra/generator"
	"github.com/parham-alvani/wedding/wedback/internal/infra/logger"
	"github.com/parham-alvani/wedding/wedback/internal/infra/wedding"
	"github.com/tidwall/pretty"
	"go.uber.org/fx"
)

const (
	// prefix indicates environment variables prefix.
	prefix = "wedback_"
	// configFile is the optional TOML file loaded on top of the defaults.
	configFile = "config.toml"
)

// Config holds all configurations.
type Config struct {
	fx.Out

	Logger    logger.Config    `json:"logger"    koanf:"logger"`
	Database  db.Config        `json:"database"  koanf:"database"`
	Generator generator.Config `json:"generator" koanf:"generator"`
	Wedding   wedding.Config   `json:"wedding"   koanf:"wedding"`
}

func Provide() Config {
	k := koanf.New(".")

	// load default configuration from default function
	if err := k.Load(structs.Provider(Default(), "koanf"), nil); err != nil {
		log.Fatalf("error loading default: %s", err)
	}

	// load configuration from file, which is optional: running without one is
	// a supported setup, but a malformed one is a mistake worth stopping for.
	if err := k.Load(file.Provider(configFile), toml.Parser()); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			log.Printf("no %s found, using defaults and environment variables", configFile)
		} else {
			log.Fatalf("error loading %s: %s", configFile, err)
		}
	}

	// load environment variables
	if err := k.Load(
		// replace __ with . in environment variables so you can reference field a in struct b
		// as a__b.
		env.Provider(prefix, ".", func(source string) string {
			base := strings.ToLower(strings.TrimPrefix(source, prefix))

			return strings.ReplaceAll(base, "__", ".")
		}),
		nil,
	); err != nil {
		log.Printf("error loading environment variables: %s", err)
	}

	var instance Config
	if err := k.Unmarshal("", &instance); err != nil {
		log.Fatalf("error unmarshalling config: %s", err)
	}

	// The full configuration dump is useful when something is misconfigured but
	// is pure noise in front of every list/import/export run, so it follows the
	// configured log level.
	if instance.Logger.Level == "debug" {
		indent, err := json.MarshalIndent(instance, "", "\t")
		if err != nil {
			panic(err)
		}

		indent = pretty.Color(indent, nil)

		log.Printf(`
================ Loaded Configuration ================
%s
======================================================
	`, string(indent))
	}

	return instance
}
