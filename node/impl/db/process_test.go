package db

import (
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/filecoin-project/go-address"

	"github.com/ipfs-force-community/venus-messager/chain/types"
	"github.com/ipfs-force-community/venus-messager/lib/message"
	"github.com/ipfs-force-community/venus-messager/node/config"
)

func TestDBProcess(t *testing.T) {
	Convey("test sqlite process", t, func() {
		filePath, err := createFile("./")
		So(err, ShouldBeNil)

		t.Log("file path: ", filePath)
		defer func() {
			So(os.Remove(filePath), ShouldBeNil)
		}()

		cfg := config.DBConfig{
			Conn:      filePath,
			Type:      "sqlite",
			DebugMode: false,
		}

		db, err := NewDbProcess(&cfg)
		So(err, ShouldBeNil)

		to, err := address.NewFromString("f01000")
		So(err, ShouldBeNil)
		from, err := address.NewFromString("f01001")
		So(err, ShouldBeNil)

		Convey("test add message", func() {

			msg := &types.Message{
				Version: 1,
				To:      to,
				From:    from,
				Nonce:   1,
			}
			meta := &message.MsgMeta{}

			err = db.AddMessage(msg, meta)
			So(err, ShouldBeNil)

			msgs, err := db.QueryMessage(from, 0, 0)
			So(err, ShouldBeNil)
			So(len(msgs), ShouldEqual, 1)
		})

		Convey("test delete message", func() {
			err := db.DelMessage(from, 0, 0)
			So(err, ShouldBeNil)

			err = db.DelMessage(from, 1, 1)
			So(err, ShouldBeNil)
		})

		Convey("test set and query nonce", func() {
			err := db.SetNonce(from, 1)
			So(err, ShouldBeNil)

			nonce, err := db.QueryNonce(from)
			So(err, ShouldBeNil)
			So(nonce, ShouldEqual, 1)

			nonce, err = db.QueryNonce(to)
			So(err, ShouldBeNil)
			So(nonce, ShouldEqual, 0)
		})
	})
}

func createFile(p string) (string, error) {
	filePath := path.Join(p, fmt.Sprintf("%d", time.Now().Second()))
	_, err := os.Create(filePath)

	return filePath, err
}
