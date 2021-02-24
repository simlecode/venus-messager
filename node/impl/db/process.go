package db

import (
	"database/sql"
	"encoding/hex"
	"time"

	"golang.org/x/xerrors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/filecoin-project/go-address"
	logging "github.com/ipfs/go-log/v2"

	"github.com/ipfs-force-community/venus-messager/chain/types"
	"github.com/ipfs-force-community/venus-messager/lib/message"
	"github.com/ipfs-force-community/venus-messager/node/config"
)

var log = logging.Logger("DB")

type DBProcess struct {
	DB *gorm.DB
}

var _ DBProcessInterface = &DBProcess{}

func NewDbProcess(cfg *config.DBConfig) (DBProcessInterface, error) {
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

	if err = db.AutoMigrate(&NonceInfo{}); err != nil {
		return nil, xerrors.Errorf("migrate model('NonceInfo') failed:%w", err)
	}

	log.Info("init db success! ...")

	return dbProc, err
}

func (dp DBProcess) QueryMessage(addr address.Address, from, to uint64) ([]Msg, error) {
	var err error
	var ms []Msg
	if from == to {
		err := dp.DB.Raw("SELECT cid,nonce,state,unsigned_msg,signed_msg,create_time FROM `msgs` WHERE address=? AND nonce=?;",
			addr.String(), from).Scan(&ms).Error
		if err != nil {
			log.Errorf("query msg err: %v", err)
			return nil, err
		}
	} else {
		err := dp.DB.Raw("SELECT cid,nonce,state,unsigned_msg,signed_msg,create_time FROM `msgs` WHERE address=? AND nonce>=? AND nonce<=?;",
			addr.String(), from, to).Scan(&ms).Error
		if err != nil {
			log.Errorf("query msg err: %v", err)
			return nil, err
		}
	}

	return ms, err
}

func (dp DBProcess) AddMessage(msg *types.Message, msgMeta *message.MsgMeta) error {
	log.Infof("from: %v, msg: %v", msg.From.String(), msg)

	msgByte, err := msg.Serialize()
	if err != nil {
		return err
	}

	metaByte, err := msgMeta.Serialize()
	if err != nil {
		return err
	}

	log.Infof("msg len: %d, msg meta len: %d", len(hex.EncodeToString(msgByte)), len(hex.EncodeToString(metaByte)))

	now := time.Now().Unix()
	// store message, subsequent unified allocation nonce
	err = dp.DB.Exec("INSERT INTO `msgs` (address,cid,nonce,state,msg_meta,unsigned_msg,create_time) VALUES (?,?,?,?,?,?,?);",
		msg.From.String(), msg.Cid().String(), ZeroNonce, UnSinged, hex.EncodeToString(metaByte), hex.EncodeToString(msgByte), now).Error
	if err != nil {
		return err
	}

	return nil
}

func (dp DBProcess) UpdateSignedMessage(id uint64, signedMsg *types.SignedMessage) error {
	log.Infof("id: %d, from: %v, msg: %v", id, signedMsg.Message.From.String(), signedMsg)

	msgByte, err := signedMsg.Serialize()
	if err != nil {
		return err
	}
	cid := signedMsg.Cid()

	log.Infof("len: %v", len(hex.EncodeToString(msgByte)))

	err = dp.DB.Exec("UPDATE `msgs` SET `unsigned_msg`=? `cid`=? WHERE `id`=?;",
		msgByte, cid, id).Error
	if err != nil {
		return err
	}

	return nil
}

func (dp DBProcess) DelMessage(addr address.Address, from, to uint64) error {
	var err error
	if from == to {
		err = dp.DB.Exec("DELETE FROM `msgs` WHERE address=? AND nonce=?;",
			addr.String(), from).Error
	} else {
		err = dp.DB.Exec("DELETE FROM `msgs` WHERE address=? AND nonce>=? AND nonce<=?;",
			addr.String(), from, to).Error
	}

	return err
}

func (dp DBProcess) DelMessageByTime(addr address.Address, start, end uint64) error {
	if end < start {
		return xerrors.Errorf("The start time(%d) is less than the end time(%d)", start, end)
	}

	return dp.DB.Exec("DELETE FROM `msgs` WHERE address=? AND create_time>=? AND create_time<=?;",
		addr.String(), start, end).Error
}

/////// nonce ///////
func (dp DBProcess) SetNonce(addr address.Address, nonce uint64) error {
	nonceInfo := NonceInfo{}
	err := dp.DB.Table("nonce_infos").Find(&nonceInfo).Where("address=?", addr.String()).Error
	if err != nil {
		return err
	}

	if nonceInfo == EmptyNonceInfo {
		return dp.DB.Exec("INSERT INTO `nonce_infos` (address, nonce) VALUES (?,?);",
			addr.String(), nonce).Error
	}

	return dp.DB.Exec("UPDATE `nonce_infos` SET `nonce`=? WHERE `address`=?;",
		nonce, addr.String()).Error
}

func (dp DBProcess) QueryNonce(addr address.Address) (uint64, error) {
	var nonceInfo NonceInfo
	err := dp.DB.Raw("SELECT nonce FROM `nonce_infos` WHERE `address`=?;", addr.String()).Scan(&nonceInfo).Error
	if err != nil {
		return 0, err
	}

	return nonceInfo.Nonce, err
}
