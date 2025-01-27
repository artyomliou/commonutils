package commonutils

import (
	"log/slog"
	"path"

	"github.com/ilyakaznacheev/cleanenv"
)

func GetConfigPath(appname, filename string) []string {
	return []string{
		path.Join("/var", appname, filename),
		filename,
	}
}

// Try reading config file, it's ok to fail.
func TryReadConfig(fullpath string, cfg interface{}) error {
	err := cleanenv.ReadConfig(fullpath, cfg)
	if err != nil {
		slog.Warn(err.Error())
	}

	// The file may not be properly loaded.
	// Should read env var anyway.
	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return err
	}

	return nil
}
