package db

import (
	"database/sql"

	"golang.org/x/xerrors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	logging "github.com/ipfs/go-log/v2"

	"github.com/ipfs-force-community/venus-messager/node/config"
)

var log = logging.Logger("DB")

type DBProcess struct {
	DB *gorm.DB
}

var _ DBProcessInterface = &DBProcess{}

func NewDbProcess(cfg *config.SQLiteDBConfig) (DBProcessInterface, error) {
	var db, err = gorm.Open(sqlite.Open(cfg.Conn), &gorm.Config{})
	var sqldb *sql.DB
	if err != nil {
		return nil, xerrors.Errorf("open database(%s) failed:%w", cfg.Conn, err)
	}

	db = db.Debug()
	if sqldb, err = db.DB(); err != nil {
		return nil, xerrors.Errorf("sqlDb failed, %w", err)
	}

	sqldb.SetConnMaxIdleTime(300)
	sqldb.SetMaxIdleConns(8)
	sqldb.SetMaxOpenConns(64)

	dbProc := &DBProcess{DB: db}

	// db.Logger.LogMode(logger.Info)

	if err = db.AutoMigrate(&Msg{}); err != nil {
		return nil, xerrors.Errorf("migrate model('Msg') failed:%w", err)
	}

	log.Info("init db success! ...")

	return dbProc, err
}

func (D DBProcess) Delete() error {
	panic("implement me")
}

func (D DBProcess) Add() error {
	panic("implement me")
}

func (D DBProcess) Query() error {
	panic("implement me")
}
