package mysql

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/ipfs-force-community/sophon-messager/config"
	"github.com/ipfs-force-community/sophon-messager/models/repo"
)

type Repo struct {
	*gorm.DB
}

func (d Repo) ActorCfgRepo() repo.ActorCfgRepo {
	return newMysqlActorCfgRepo(d.DB)
}

func (d Repo) MessageRepo() repo.MessageRepo {
	return newMysqlMessageRepo(d.DB)
}

func (d Repo) AddressRepo() repo.AddressRepo {
	return newMysqlAddressRepo(d.DB)
}

func (d Repo) SharedParamsRepo() repo.SharedParamsRepo {
	return newMysqlSharedParamsRepo(d.DB)
}

func (d Repo) NodeRepo() repo.NodeRepo {
	return newMysqlNodeRepo(d.DB)
}

func (d Repo) AutoMigrate() error {
	return d.GetDb().AutoMigrate(mysqlActorCfg{}, mysqlMessage{}, mysqlAddress{}, mysqlSharedParams{}, mysqlNode{})
}

func (d Repo) GetDb() *gorm.DB {
	return d.DB
}

func (d Repo) DbClose() error {
	// return d.DbClose()
	// todo:
	return nil
}

func (d Repo) Transaction(cb func(txRepo repo.TxRepo) error) error {
	return d.DB.Transaction(func(tx *gorm.DB) error {
		txRepo := &TxMysqlRepo{tx}
		return cb(txRepo)
	})
}

var _ repo.TxRepo = (*TxMysqlRepo)(nil)

type TxMysqlRepo struct {
	*gorm.DB
}

func (t *TxMysqlRepo) ActorCfgRepo() repo.ActorCfgRepo {
	return newMysqlActorCfgRepo(t.DB)
}

func (t *TxMysqlRepo) SharedParamsRepo() repo.SharedParamsRepo {
	return newMysqlSharedParamsRepo(t.DB)
}

func (t *TxMysqlRepo) NodeRepo() repo.NodeRepo {
	return newMysqlNodeRepo(t.DB)
}

func (t *TxMysqlRepo) MessageRepo() repo.MessageRepo {
	return newMysqlMessageRepo(t.DB)
}

func (t *TxMysqlRepo) AddressRepo() repo.AddressRepo {
	return newMysqlAddressRepo(t.DB)
}

func OpenMysql(cfg *config.MySqlConfig) (repo.Repo, error) {
	db, err := gorm.Open(mysql.Open(cfg.ConnectionString), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info), // 日志配置
	})
	if err != nil {
		return nil, fmt.Errorf("[db connection failed] Database name: %s %w", cfg.ConnectionString, err)
	}

	db.Set("gorm:table_options", "CHARSET=utf8mb4")
	if cfg.Debug {
		db = db.Debug()
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConn)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConn)
	sqlDB.SetConnMaxLifetime(time.Minute * cfg.ConnMaxLifeTime)

	// 使用插件
	// db.Use(&TracePlugin{})
	return &Repo{
		db,
	}, nil
}
