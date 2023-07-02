package framework

import (
	"os"
	"runtime"
	"strings"

	"github.com/TeemoKill/WanZBlog/constants"
	"github.com/TeemoKill/WanZBlog/log"
	"github.com/jinzhu/configor"
)

func (e *BlogEngine) loadConfig(autoReloadCallbackFn func(interface{})) (err error) {
	return configor.New(
		&configor.Config{
			AutoReload:         true,
			AutoReloadCallback: autoReloadCallbackFn,
		},
	).
		Load(e.Cfg, constants.WanZBlogConfigFilePath)
}

// GenerateExampleConfig creates example configs and write to config file
func (e *BlogEngine) generateExampleConfig(filePath string) (err error) {
	logger := log.CurrentModuleLogger()

	err = os.WriteFile(filePath, []byte(exampleConfig()), 0755)
	if err != nil {
		logger.
			WithError(err).
			WithField("config_filepath", filePath).
			Errorf("failed to generate example config")
		return err
	}
	logger.
		WithField("config_filepath", filePath).
		Infof("Minimum configuration has been generated. Please modify as needed and rerun. " +
			"For advanced configuration, please refer to the help documentation.")
	return err
}

func exampleConfig() string {
	c := `
log_level = "info"

listen_addr = "0.0.0.0:8080"
gin_mode = "debug"
server_shutdown_max_wait_seconds = 5
sqlite_path = "./sqlite.db"

`
	// To solve the issue of incorrect line breaks when opening a file in Notepad on Windows
	if runtime.GOOS == "windows" {
		c = strings.ReplaceAll(c, "\n", "\r\n")
	}
	return c
}
