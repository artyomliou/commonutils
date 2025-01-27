package commonutils_test

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/artyomliou/commonutils"
	"github.com/stretchr/testify/assert"
)

func TestTryReadConfig(t *testing.T) {
	type testConfig struct {
		AppName  string `yaml:"app_name" env:"APP_NAME"`                       // test load value from various source
		LogLevel string `yaml:"log_level" env:"LOG_LEVEL" env-default:"debug"` // test default value
		DB       struct {
			Host string `yaml:"host" env:"HOST"`                                    // test load value into embed struct
			Port int    `yaml:"port" env:"PORT" env-default:"5432" env-layout:"%d"` //  test default value in embed struct
		} `yaml:"db" env-prefix:"DB_"` // test load value with prefix
	}

	type testcase struct {
		title     string
		cfg       interface{}
		setupFunc func(*testing.T) (*os.File, error)
	}
	testcases := []testcase{
		{
			title: "read yaml",
			cfg:   &testConfig{},
			setupFunc: func(t *testing.T) (*os.File, error) {
				f, err := os.CreateTemp("", "TestReadConfig*.yaml")
				if err != nil {
					return nil, err
				}

				contents := []string{
					`app_name: "testing"`,
					`db:`,
					`  host: mydb`,
				}
				_, err = f.WriteString(strings.Join(contents, "\n"))
				if err != nil {
					return nil, err
				}
				_, err = f.Seek(0, io.SeekStart)
				if err != nil {
					return nil, err
				}
				return f, nil
			},
		},
		{
			title: "read env file",
			cfg:   &testConfig{},
			setupFunc: func(t *testing.T) (*os.File, error) {
				f, err := os.CreateTemp("", "TestReadConfig*.env")
				if err != nil {
					return nil, err
				}

				contents := []string{
					`APP_NAME = "testing"`,
					`DB_HOST = "mydb"`,
				}
				_, err = f.WriteString(strings.Join(contents, "\n"))
				if err != nil {
					return nil, err
				}
				_, err = f.Seek(0, io.SeekStart)
				if err != nil {
					return nil, err
				}
				return f, nil
			},
		},
		{
			title: "read env vars",
			cfg:   &testConfig{},
			setupFunc: func(t *testing.T) (*os.File, error) {
				// create temp file to satisfy function signature
				f, err := os.CreateTemp("", "TestReadConfig*.unused")
				if err != nil {
					return nil, err
				}

				// don't use os.Setenv to avoid pollution
				t.Setenv("APP_NAME", "testing")
				t.Setenv("DB_HOST", "mydb")
				return f, nil
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.title, func(t *testing.T) {
			f, err := tc.setupFunc(t)
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(f.Name())

			assert.NoError(t, commonutils.TryReadConfig(f.Name(), tc.cfg))

			cfg := tc.cfg.(*testConfig)
			assert.Equal(t, "testing", cfg.AppName)
			assert.Equal(t, "debug", cfg.LogLevel)
			assert.Equal(t, "mydb", cfg.DB.Host)
			assert.Equal(t, 5432, cfg.DB.Port)
		})
	}
}
