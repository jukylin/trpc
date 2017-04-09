package rpc

import (
	//"fmt"
	//"trpc/hey"
	"github.com/weixinhost/yar.go/client"
	yar "github.com/weixinhost/yar.go"
	"errors"
)

func Yar(args *RpcArgs) (interface{}, error) {

	client, err := client.NewClient(args.Url)

	if err != nil {
		//fmt.Println("error", err)
		return nil, errors.New(err.String())
	}

	//这是默认值
	client.Opt.Timeout = 1000 * 30 //30s
	//这是默认值
	client.Opt.Packager = "json"
	//这是默认值
	client.Opt.Encrypt = false
	//这是默认值
	client.Opt.EncryptPrivateKey = ""
	//这是默认值
	client.Opt.MagicNumber = yar.MagicNumber

	client.IsBenchClient = args.Bench

	var ret interface{}

	s := GetArgs(args.Args)


	callErr := client.Call(args.Fn, &ret, s...)
	if callErr != nil {
		return nil, errors.New(callErr.String())
	}

	if args.Bench {
		return string(client.PackBody), nil
	}


	return ret,nil

}
