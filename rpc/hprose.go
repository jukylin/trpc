package rpc

import (
	"github.com/Aqiling/hprose-golang/rpc"
	"github.com/Aqiling/hprose-golang/io"
	"reflect"
	"time"
)

func Hprose(args *RpcArgs) (interface{}, error) {

	client := rpc.NewHTTPClient(args.Url)
	client.MaxIdleConnsPerHost = 128
	client.SetMaxConcurrentRequests(128)
	client.IsBench = args.Bench

	var in []reflect.Value
	var unResutl interface{}

	inputArgs := GetArgs(args.Args)


	for _,v := range inputArgs {
		in = append(in, reflect.ValueOf(v))
	}

	var refType []reflect.Type
	settings := &rpc.InvokeSettings{
		Timeout:        time.Duration(0),
		ResultTypes: refType,
	}

	_, err := client.Invoke(args.Fn, in, settings)
	if err != nil {
		return "", err
	}

	if args.Bench {
		return string(client.Request), nil
	}else{
		returnString := client.Reponse[1:len(client.Reponse) - 1]
		io.Unserialize(returnString, &unResutl, false)

		return unResutl, nil
	}
}