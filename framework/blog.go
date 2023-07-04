package framework

import (
	"context"
	"gorm.io/gorm"
	"net/http"
	"os"
	"time"

	"github.com/TeemoKill/WanZBlog/apirouter"
	"github.com/TeemoKill/WanZBlog/config"
	"github.com/TeemoKill/WanZBlog/constants"
	"github.com/TeemoKill/WanZBlog/log"
	"github.com/sirupsen/logrus"
)

type BlogEngine struct {
	Cfg *config.Config

	Logger *logrus.Logger
	Server *http.Server
	Db     *gorm.DB
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
		err = e.generateExampleConfig(constants.WanZBlogConfigFilePath)
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

	// init database
	if err := e.connectDB(); err != nil {
		logger.WithError(err).
			Errorf("connect database error")
		return err
	}

	// init router
	e.InitHTTPServer()

	return nil
}

func (e *BlogEngine) InitHTTPServer() {

	apiRouter := apirouter.New(
		e.Cfg, e.Db,
	)

	e.Server = &http.Server{
		Addr:    e.Cfg.ListenAddr,
		Handler: apiRouter,
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

	if err := e.Server.Shutdown(ctx); err != nil {
		e.Logger.WithError(err).
			Errorf("Shutdown http Server error")
		return err
	}

	dbConn, err := e.Db.DB()
	if err != nil {
		e.Logger.WithError(err).
			Errorf("get gorm db connection for close error")
	}
	if err := dbConn.Close(); err != nil {
		e.Logger.WithError(err).
			Errorf("Close Sqlite connection error")
		return err
	}

	return nil
}
