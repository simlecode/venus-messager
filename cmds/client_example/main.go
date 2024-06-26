package main

import (
	"context"
	"fmt"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	venustypes "github.com/filecoin-project/venus/pkg/types"
	"github.com/ipfs-force-community/venus-messager/api/client"
	"github.com/ipfs-force-community/venus-messager/config"
	"github.com/ipfs-force-community/venus-messager/types"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.ReadConfig("/Users/lijunlong/Desktop/workload/venus-messager/messager.toml")
	if err != nil {
		log.Fatal(err)
		return
	}

	header := http.Header{}
	client, closer, err := client.NewMessageRPC(context.Background(), "http://"+cfg.API.Address+"/rpc/v0", header)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer closer()
	addr, _ := address.NewFromString("t3v3vpwnsfqmp6cj65lzflxxv3vkiiwdj33l3f4bgrwkio7w27jmwtsby372m677g7x5h6eqsaieelwr6n74la")
	uid, err := client.PushMessage(context.Background(), &types.Message{
		ID: types.NewUUID(),
		UnsignedMessage: venustypes.UnsignedMessage{
			Version: 0,
			To:      addr,
			From:    addr,
			Nonce:   1,
			Value:   abi.NewTokenAmount(100),
			Method:  0,
		},
		Meta: &types.MsgMeta{
			ExpireEpoch:       1000000,
			GasOverEstimation: 1.25,
			MaxFee:            big.NewInt(10000000000000000),
			MaxFeeCap:         big.NewInt(10000000000000000),
		},
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("send message " + uid)

	msg, err := client.WaitMessage(context.Background(), uid)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("wait for message ", msg.SignedCid)
	fmt.Println("code:", msg.Receipt.ExitCode)
	fmt.Println("gas_used:", msg.Receipt.GasUsed)
	fmt.Println("return_value:", msg.Receipt.ReturnValue)
}
