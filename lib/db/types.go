package db

type Msg struct {
	ID          uint64 `gorm:"type:unsigned integer;primary key;autoincrement;unique"`
	Address     string `gorm:"type:varchar(255);index:address,nonce"`
	Cid         string `gorm:"type:varchar(128);uniqueIndex"`
	Nonce       uint64 `gorm:"type:unsigned integer;index:address,nonce;default:0"`
	UnSignedMsg string `gorm:"type:varchar(1024);column:unsigned_msg"`
	SignedMsg   string `gorm:"type:varchar(1024);column:signed_msg"`
	CreateTime  uint64 `gorm:"type:varchar(1024);column:create_time"`
}
