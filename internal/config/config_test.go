package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestValidate(t *testing.T) {
	cases := []struct {
		name     string
		input    *AppConfig
		errors   []string
		warnings []string
	}{
		{
			name: "empty details",
			input: &AppConfig{
				Name:           "",
				Version:        "",
				Description:    "",
				DatabaseConfig: &DatabaseConfig{},
			},
			errors: []string{
				"'app_name' can't be empty",
				"'version' can't be empty",
				"'database.name' can't be empty",
				"'database.driver' can't be empty",
			},
			warnings: []string{
				"'description' is empty",
				"'http_server' configuration is empty so the server will be started on the default configurations",
				"'http_server.hostname' is empty so it will default to: 'localhost'",
				"'http_server.port' is empty so it will default to: '8080'",
			},
		},
		{
			name: "empty postgres config",
			input: &AppConfig{
				Name:             "Gloak",
				Version:          "1.0.1",
				Description:      "An approchable IdP",
				HttpServerConfig: &HttpServerConfig{Port: 8080, Hostname: "localhost"},
				DatabaseConfig:   &DatabaseConfig{DBName: "gloakdb", Driver: "postgres"},
			},
			errors: []string{
				"'database.hostname' can't be empty",
				"'database.user' can't be empty",
				"'database.password' can't be empty",
			},
			warnings: []string{
				"'database.sslmode' is empty so it will default to: 'allow'",
				"'database.port' is empty or 0 so it will default to: 5432",
			},
		},
		{
			name: "complete postgres config",
			input: &AppConfig{
				Name:             "Gloak",
				Version:          "1.0.1",
				Description:      "An approchable IdP",
				HttpServerConfig: &HttpServerConfig{Port: 8080, Hostname: "localhost"},
				DatabaseConfig: &DatabaseConfig{
					Driver:   "postgres",
					Hostname: "localhost",
					Port:     5432,
					Password: "supersecret",
					User:     "gloak",
					DBName:   "gloakdb",
					SSLMode:  "prefer",
				},
			},
			errors:   []string{},
			warnings: []string{},
		},
		{
			name: "empty sqlite3 config",
			input: &AppConfig{
				Name:             "Gloak",
				Version:          "1.0.1",
				Description:      "An approchable IdP",
				HttpServerConfig: &HttpServerConfig{Port: 8080, Hostname: "localhost"},
				DatabaseConfig:   &DatabaseConfig{Driver: "sqlite3"},
			},
			errors: []string{"'database.name' can't be empty"},
			warnings: []string{
				fmt.Sprintf("'database.path' is empty so it will default to: '%v'", getPath()),
			},
		},
		{
			name: "sqlite3 path do not exist",
			input: &AppConfig{
				Name:             "Gloak",
				Version:          "1.0.1",
				Description:      "An approchable IdP",
				HttpServerConfig: &HttpServerConfig{Port: 8080, Hostname: "localhost"},
				DatabaseConfig: &DatabaseConfig{
					DBName: "gloakdb",
					Driver: "sqlite3",
					Path:   "/path/do/not/exist-for/sure",
				},
			},
			errors: []string{
				fmt.Sprintf("'database.path' %v provided does not exists", "/path/do/not/exist-for/sure"),
			},
			warnings: []string{},
		},
		{
			name: "sqlite3 insufficient permissions",
			input: &AppConfig{
				Name:             "Gloak",
				Version:          "1.0.1",
				Description:      "An approchable IdP",
				HttpServerConfig: &HttpServerConfig{Port: 8080, Hostname: "localhost"},
				DatabaseConfig: &DatabaseConfig{
					DBName: "gloakdb",
					Driver: "sqlite3",
					Path:   "/",
				},
			},
			errors: []string{
				fmt.Sprintf("'database.path' %v not enough permissions to read and write in this directory", "/"),
			},
			warnings: []string{},
		},
		{
			name: "unsupported database driver",
			input: &AppConfig{
				Name:             "Gloak",
				Version:          "1.0.1",
				Description:      "An approchable IdP",
				HttpServerConfig: &HttpServerConfig{Port: 8080, Hostname: "localhost"},
				DatabaseConfig: &DatabaseConfig{
					DBName: "gloakdb",
					Driver: "post",
					Path:   "/",
				},
			},
			errors: []string{
				fmt.Sprintf("'database.driver' %v not supported, only [postgres, sqlite3] are supported", "post"),
			},
			warnings: []string{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			errs, ws := validate(tc.input)
			me := make(map[string]string, 0)
			mw := make(map[string]string, 0)

			for _, err := range errs {
				me[err.Error()] = "exists"
			}

			for _, w := range ws {
				mw[w.Error()] = "exists"
			}

			for _, err := range tc.errors {
				if _, ok := me[err]; !ok {
					t.Errorf("expected %#v to be in the errors list", err)
					t.Log(me)
				}
			}

			for _, w := range tc.warnings {
				if _, ok := mw[w]; !ok {
					t.Errorf("expected %#v to be in the warnings list", w)
					t.Log(mw)
				}
			}
		})
	}
}

func getPath() string {
	p, _ := os.Executable()
	return filepath.Dir(p)
}
