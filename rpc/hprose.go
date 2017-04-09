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

	in = append(in, reflect.ValueOf("wo"))
	var tt []reflect.Type
	//i := make(map[string]interface{})
	//i := 13
	//tt = append(tt, reflect.TypeOf((*interface{})(nil)).Elem())
	//fmt.Println(tt)

	settings := &rpc.InvokeSettings{
		Timeout:        time.Duration(0),
		ResultTypes: tt,
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