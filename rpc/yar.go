package rpc

import (
	"strings"
	//"fmt"
	//"trpc/hey"
	"github.com/weixinhost/yar.go/client"
	yar "github.com/weixinhost/yar.go"
	"encoding/json"
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

	s := make([]interface{}, len(args.Args))
	for i, v := range args.Args {
		//解析数组
		if strings.Contains(strings.ToLower(v), "arrfile:") {
			tmp := strings.Split(v, ":")
			jsonData, err := ReadJson(tmp[1])
			if err != nil {
				//fmt.Println(err.Error())
				return nil, err
			}

			b := []byte(jsonData)
			var m interface{}
			if err := json.Unmarshal(b, &m); err != nil {
				//fmt.Println(err)
				return nil, err
			}
			s[i] = m
		} else if strings.Contains(strings.ToLower(v), "arr:") {
			replaceStr := strings.Replace(v, "arr:", "", 1)
			replaceStr = strings.Trim(replaceStr, "#")
			splitStr := strings.Split(replaceStr, "#")
			arrMap := make(map[string]string, len(splitStr))
			for _, spV := range splitStr {
				spspV := strings.Split(spV, "=")
				arrMap[spspV[0]] = spspV[1]
			}

			s[i] = arrMap
		} else {
			s[i] = v
		}
	}

	callErr := client.Call(args.Fn, &ret, s...)
	if callErr != nil {
		//fmt.Println("error", callErr)
		return nil, errors.New(callErr.String())
	}

	if args.Bench {
		return string(client.PackBody), nil
		//hey := new(hey.Hey)
		//hey.Url = args.Url
		//hey.Method = "post"
		//hey.Num = args.Nrun
		//hey.Con = args.Ncon
		//hey.Body = client.PackBody
		//hey.ContentType = "application/json"
		//hey.RunHey()
	}


	return ret,nil

}
