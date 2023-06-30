package framework

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/TeemoKill/WanZBlog/config"
	"github.com/TeemoKill/WanZBlog/constants"
	"github.com/TeemoKill/WanZBlog/log"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type BlogEngine struct {
	Cfg *config.Config

	Logger *logrus.Logger
	Server *http.Server
}

func New() *BlogEngine {
	return &BlogEngine{Cfg: &config.Config{}}
}

func (e *BlogEngine) Init() (err error) {
	logger := log.CurrentModuleLogger()

	fileInfo, err := os.Stat(constants.WanZBlogConfigFilePath)
	switch {
	case os.IsNotExist(err):
		logger.WithField("config_filepath", constants.WanZBlogConfigFilePath).
			Warnf("No configuration file detected, generating one. You can ignore this warning if it is the first time running.")
		err = e.generateExampleConfig()
		if err != nil {
			logger.WithError(err).
				Errorf("failed to generate example config, exitting")
			return err
		}
		return constants.ExampleConfigGeneratedErr

	case err != nil:
		logger.WithError(err).
			WithField("config_filepath", constants.WanZBlogConfigFilePath).
			Errorf("Failed to check the configuration file, exiting")
		return err

	default:
		if fileInfo.IsDir() {
			logger.WithField("config_filepath", constants.WanZBlogConfigFilePath).
				Errorf("Detected configuration file path, but the path is a directory instead of a file")
			return constants.ConfigPathNotFileErr
		} else {
			logger.WithField("config_filepath", constants.WanZBlogConfigFilePath).
				Infof("Detected configuration file, starting with existing configuration file")
		}
	}

	confReloadCallback := func(config interface{}) {
		logger.Infof("%v changed", config)
	}
	err = e.loadConfig(confReloadCallback)
	if err != nil {
		logger.WithError(err).
			Errorf("Failed to load configuration file! Please check if the configuration file format is correct")
		return err
	}

	// init logger
	logLevel, err := logrus.ParseLevel(e.Cfg.LogLevel)
	if err != nil {
		logrus.SetLevel(logrus.DebugLevel)
		logger.Warn("Unable to recognize logLevel, will use Debug level")
	} else {
		logrus.SetLevel(logLevel)
		logger.Infof("Set log level to %s", logLevel.String())
	}
	e.Logger = log.StandardLogger()

	// init router
	e.InitHTTPServer()

	return nil
}

func (e *BlogEngine) InitHTTPServer() {
	gin.DisableConsoleColor()
	gin.SetMode(e.Cfg.GinMode)

	router := gin.New()
	router.Use(sentrygin.New(sentrygin.Options{}))

	pprof.Register(router)

	// register router
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World!")
	})

	e.Server = &http.Server{
		Addr:    e.Cfg.ListenAddr,
		Handler: router,
	}
}

func (e *BlogEngine) StartService() {

	go func() {
		if err := e.Server.ListenAndServe(); err != nil {
			e.Logger.WithField("listen_error", err).
				Infof("router error while listening")
		}
	}()

}

func (e *BlogEngine) Stop() error {
	e.Logger.Info("wanzblog framework stopping ...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(e.Cfg.ServerShutdownMaxWaitSeconds)*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		e.Logger.Error("server shutdown timeout")
	default:
		e.Logger.Info("server exiting")
	}
	e.Logger.Info("server exited")

	return e.Server.Shutdown(ctx)
}
