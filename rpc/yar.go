package rpc

import (
	"strings"
	"fmt"
	"trpc/hey"
	"time"
	"github.com/weixinhost/yar.go/client"
	yar "github.com/weixinhost/yar.go"
	"encoding/json"
)

func Yar(args *RpcArgs)  {

	client, err := client.NewClient(args.Url)

	if err != nil {
		fmt.Println("error", err)
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

	t1 := time.Now() // get current time

	s := make([]interface{}, len(args.Args))
	for i, v := range args.Args {
		//解析数组
		if strings.Contains(strings.ToLower(v), "arrfile:") {
			tmp := strings.Split(v, ":")
			jsonData, err := ReadJson(tmp[1])
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			b := []byte(jsonData)
			var m interface{}
			if err := json.Unmarshal(b, &m); err != nil {
				fmt.Println(err)
				return
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
		fmt.Println("error", callErr)
		return
	}

	if args.Bench {
		hey := new(hey.Hey)
		hey.Url = args.Url
		hey.Method = "post"
		hey.Num = args.Nrun
		hey.Con = args.Ncon
		hey.Body = client.PackBody
		hey.ContentType = "application/json"
		hey.RunHey()
	} else {

		elapsed := time.Since(t1)
		if args.Format == true {
			is_map := FormatResutl(ret, 1)
			if is_map == true {
				fmt.Println("result:\r\n", buf.String())
			} else {
				fmt.Println("result:", buf.String())
			}
		} else {
			fmt.Println("result:", ret)
		}
		fmt.Println("runtime: ", elapsed)
	}
}
