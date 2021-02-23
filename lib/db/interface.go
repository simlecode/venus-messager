package db

import "github.com/ipfs-force-community/venus-messager/node/config"

type DBProcessInterface interface {
	NewDbProcess(config *config.SQLiteDBConfig) (DBProcessInterface, error)
	Delete() error
	Add() error
	Query() error
}
