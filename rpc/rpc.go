package rpc

import (
	yar "github.com/weixinhost/yar.go"
	"github.com/weixinhost/yar.go/client"
	"fmt"
	"time"
	"strings"
	"os"
	"io/ioutil"
	"encoding/json"
	"reflect"
	"bytes"
	"trpc/hey"
)

var buf bytes.Buffer

const ENPTY_NUM = 6

func DebugStart(url string, fn string, format bool, bench bool,nrun int, ncon int, args []string){
	client, err := client.NewClient(url)

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

	client.IsBenchClient = bench

	var ret interface{}

	t1 := time.Now() // get current time

	s := make([]interface{}, len(args))
	for i, v := range args {
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
		}else if strings.Contains(strings.ToLower(v), "arr:"){
			replaceStr := strings.Replace(v, "arr:", "", 1)
			replaceStr = strings.Trim(replaceStr, "#")
			splitStr := strings.Split(replaceStr, "#")
			arrMap := make(map[string]string,len(splitStr))
			for _, spV := range splitStr {
				spspV := strings.Split(spV, "=")
				arrMap[spspV[0]] = spspV[1]
			}

			s[i] = arrMap
		}else {
			s[i] = v
		}
	}

	callErr := client.Call(fn, &ret,s...)

	if callErr != nil {
		fmt.Println("error", callErr)
	}

	if bench {
		hey := new(hey.Hey)
		hey.Url = url
		hey.Method = "post"
		hey.Num = nrun
		hey.Con = ncon
		hey.Body = client.PackBody
		hey.ContentType = "application/json"
		hey.RunHey()
	}else {

		elapsed := time.Since(t1)

		if format == true {
			is_map := FormatResutl(ret, 1)
			if is_map == true {
				fmt.Println("result:\r\n", buf.String())
			}else{
				fmt.Println("result:", buf.String())
			}
		} else {
			fmt.Println("result:", ret)
		}
		fmt.Println("runtime: ", elapsed)
	}
}

/**
@param interface{} result 需要格式化的数据
@param int i 层级 默认 0
 */
func FormatResutl(result interface{}, i int) bool {

	var is_map bool
	if reflect.ValueOf(result).Kind() == reflect.Map {

		is_map = true

		WriteString(func(){
			buf.WriteString("[\n")
		}, i)

		switch resultValue := result.(type){
		case map[string]interface{}:
			for k, vv := range resultValue{
				if reflect.ValueOf(vv).Kind() == reflect.Map {

					WriteString(func(){
						buf.WriteString(fmt.Sprintf("   '%s' => \r\n", k))
					}, i)
					FormatResutl(vv, i + 1)

				}else if reflect.ValueOf(vv).Kind() == reflect.String {

					WriteString(func(){
						buf.WriteString(fmt.Sprintf("   '%s' => '%s',\r\n", k, vv))
					}, i)

				} else {

					WriteString(func(){
						buf.WriteString(fmt.Sprintf("  '%s' => %s,\r\n", k, vv))
					}, i)

				}
			}
		default:
			fmt.Println("unknown type", resultValue)
		}
		if i == 1 {
			WriteString(func() {
				buf.WriteString(" ],\r\n")
			}, i)
		}else {
			WriteString(func() {
				buf.WriteString("],\r\n")
			}, i)
		}
	}else{
		is_map = false
		buf.WriteString(fmt.Sprintf("%s", result))
	}


	return is_map
}


func WriteString(fn func(), i int){
	buf.WriteString(strings.Repeat(" ", (i - 1) * ENPTY_NUM))
	fn()
}

/**
读取文件json数组
 */
func ReadJson(path string) (string,error) {
	fi,err := os.Open(path)
	if err != nil{
		return "",err
	}
	defer fi.Close()
	fd,err := ioutil.ReadAll(fi)
	if err != nil {
		return "", err
	}

	return string(fd), nil
}
