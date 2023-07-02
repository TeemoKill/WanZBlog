package framework

import (
	"database/sql"
	"github.com/TeemoKill/WanZBlog/constants"
	"github.com/TeemoKill/WanZBlog/datamodel"
	"github.com/TeemoKill/WanZBlog/log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func (e *BlogEngine) connectDB() error {
	logger := log.CurrentModuleLogger()

	dbFileInfo, err := os.Stat(e.Cfg.SqlitePath)
	switch {
	case os.IsNotExist(err):
		logger.WithField("sqlite_path", e.Cfg.SqlitePath).
			Warnf("Sqlite db file not exist")
		err = e.initDB(e.Cfg.SqlitePath)
		if err != nil {
			logger.WithError(err).
				Errorf("failed to initialize database, exitting")
			return err
		}
		return err
	case err != nil:
		logger.WithError(err).
			WithField("sqlite_path", e.Cfg.SqlitePath).
			Errorf("Failed to open sqlite db file, exiting")
		return err
	case dbFileInfo.IsDir():
		logger.WithField("sqlite_path", e.Cfg.SqlitePath).
			Errorf("Detected sqlite db file path, but the path is a directory instead of a file")
		return constants.SqliteDBPathNotFileErr
	default:
		logger.WithField("sqlite_path", e.Cfg.SqlitePath).
			Infof("Detected sqlite db file")
	}

	sqliteConn, err := sql.Open("sqlite3", e.Cfg.SqlitePath)
	if err != nil {
		logger.WithError(err).
			Errorf("Failed to open sqlite!")
		return err
	}
	e.Db = sqliteConn

	return nil
}

func (e *BlogEngine) initDB(dbPath string) error {
	logger := log.CurrentModuleLogger()

	db, err := gorm.Open(
		sqlite.Open(dbPath),
	)
	dbConn, _ := db.DB()

	if err != nil {
		logger.WithError(err).
			Errorf("initialize db error")
		return err
	}

	datamodel.InitDataModel(db)

	if err := dbConn.Close(); err != nil {
		logger.WithError(err).
			Errorf("close db connection error")
		return err
	}

	return nil
}
