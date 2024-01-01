package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/paulocuambe/gloak/conf"

	"gopkg.in/ini.v1"
)

type AppConfig struct {
	Name              string `ini:"app_name"`
	Version           string `ini:"version"`
	Description       string `ini:"description"`
	*HttpServerConfig `ini:"http_server"`
	*DatabaseConfig   `ini:"database"`
}

type HttpServerConfig struct {
	Port     int    `ini:"port"`
	Hostname string `ini:"hostname"`
}

func (h *HttpServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", h.Hostname, h.Port)
}

type DriverName string

func (d DriverName) GetName() string {
	if d == "postgres" {
		return "pgx"
	}
	return "sqlite3"
}

type DatabaseConfig struct {
	Driver   DriverName `ini:"driver"`
	Path     string     `ini:"path"`
	DBName   string     `ini:"name"`
	Hostname string     `ini:"hostname"`
	Port     int        `ini:"port"`
	User     string     `ini:"user"`
	Password string     `ini:"password"`
	SSLMode  string     `ini:"sslmode"`
}

// returns the db dsn
// if the db configuration is not defined or if its sqlite,
// create a local db file in the execution path
func (d *DatabaseConfig) DSN() string {
	// use sqlite if its not postgres
	if d.Driver.GetName() != "pgx" {
		return filepath.Join(d.Path, fmt.Sprintf("%v.db", d.DBName))
	}

	return fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%s application_name=%s", d.Hostname, d.Port, d.DBName, d.User, d.Password, d.SSLMode, d.DBName)
}

// loads config and returns the config, errors and warnings
func LoadConfig() (*AppConfig, error, error) {
	data, err := conf.Files.ReadFile("app.ini")

	if err != nil {
		return nil, err, nil
	}

	cfg, err := ini.Load(data)
	if err != nil {
		return nil, err, nil
	}

	var conf AppConfig
	err = cfg.MapTo(&conf)
	if err != nil {
		return nil, err, nil
	}

	errs, warnings := validate(&conf)

	return &conf, errs, warnings
}

// returns errors and warnings and a bolean denoting if has errors not just warnings
func validate(cfg *AppConfig) (error, error) {
	errs := make([]error, 0, 1)
	warnings := make([]error, 0, 1)

	if cfg.Name == "" {
		errs = append(errs, errors.New("'app_name' should not be empty"))
	}

	if cfg.Version == "" {
		errs = append(errs, errors.New("'version' can't be empty"))
	}

	if cfg.Description == "" {
		warnings = append(warnings, errors.New("'description' is empty"))
	}

	if cfg.HttpServerConfig == nil {
		warnings = append(warnings, errors.New("'http_server' configuration is empty so the server will be started on the default configurations"))
	}

	if cfg.HttpServerConfig.Hostname == "" {
		cfg.HttpServerConfig.Hostname = "localhost"
		warnings = append(warnings, errors.New("'http_server.hostname' is empty so it will default to: 'localhost'"))
	}

	if cfg.HttpServerConfig.Port == 0 {
		cfg.HttpServerConfig.Port = 8080
		warnings = append(warnings, errors.New("'http_server.port' is empty so it will default to: '8080'"))
	}

	if cfg.DatabaseConfig == nil {
		errs = append(errs, errors.New("'database' configuration can't be empty"))
	}

	if cfg.DatabaseConfig.DBName == "" {
		errs = append(errs, errors.New("'database.name' can't be empty"))
	}

	if cfg.DatabaseConfig.Driver == "" {
		errs = append(errs, errors.New("'database.driver' can't be empty"))
	}

	if cfg.DatabaseConfig.Driver != "" {
		if cfg.DatabaseConfig.Driver != "postgres" && cfg.DatabaseConfig.Driver != "sqlite3" {
			errs = append(errs, errors.New("'database.driver' only allows the following values: [postgres, sqlite3]"))
		}
	}

	if cfg.DatabaseConfig.Driver == "postgres" {
		if cfg.DatabaseConfig.Hostname == "" {
			errs = append(errs, errors.New("'database.hostname' can't be empty"))
		}

		if cfg.DatabaseConfig.Port == 0 {
			warnings = append(warnings, errors.New("'database.port' can't be empty or 0"))
		}

		if cfg.DatabaseConfig.User == "" {
			errs = append(errs, errors.New("'database.user' can't be empty"))
		}

		if cfg.DatabaseConfig.Password == "" {
			errs = append(errs, errors.New("'database.password' can't be empty"))
		}

		if cfg.DatabaseConfig.SSLMode == "" {
			cfg.DatabaseConfig.SSLMode = "allow"
			warnings = append(warnings, errors.New("'database.sslmode' is empty so it will default to allow"))
		}
	}

	if cfg.DatabaseConfig.Driver == "sqlite3" {
		// if path is empty while using sqlite3 it will default to the path were the program is being executed in
		if cfg.DatabaseConfig.Path == "" {
			p, _ := os.Executable()
			p = filepath.Dir(p)
			cfg.DatabaseConfig.Path = p
			warnings = append(warnings, fmt.Errorf("'database.path' is empty so it will default to: %v", p))
		}

		if p := cfg.DatabaseConfig.Path; p != "" {
			_, err := os.Stat(p)

			if err != nil {
				if errors.Is(err, os.ErrExist) {
					errs = append(errs, fmt.Errorf("'database.path' %v provided does not exist", p))
				} else if errors.Is(err, os.ErrPermission) {
					errs = append(errs, fmt.Errorf("'database.path' %v not enough permissions to read from this location", p))
				} else {
					warnings = append(warnings, fmt.Errorf("'database.path' %v: %w", p, err))
				}
			}
			// create the db file
			if f, err := os.Create(cfg.DatabaseConfig.DSN()); err != nil {
				errs = append(errs, fmt.Errorf("'database.path' %v can't be written to", p))
			} else {
				defer f.Close()
			}
		}
	}

	var e error
	if len(errs) > 0 {
		e = ConfigErr{Errs: errs}
	}

	var w error
	if len(warnings) > 0 {
		w = ConfigWarnigs{Warnings: warnings}
	}

	return e, w
}
