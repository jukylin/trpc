package rpc

import (
	"fmt"
	"strings"
	"os"
	"io/ioutil"
	"reflect"
	"bytes"
	"time"
	"trpc/hey"
	//"github.com/weixinhost/yar.go/client"
	"encoding/json"
	"strconv"
)

var buf bytes.Buffer

const ENPTY_NUM = 6

type RpcArgs struct {
	Type string
	Url string
	Fn string
	Format bool
	Bench bool
	Nrun int
	Ncon int
	Args []string
}

func DebugStart(args *RpcArgs){
	t1 := time.Now() // get current time

	var rpcString interface{}
	var err error
	if args.Type == "yar" {
		rpcString, err = Yar(args)
	}else if args.Type == "hprose"{
		rpcString,err = Hprose(args)
	}else{
		fmt.Println("不存在",args.Type,"rpc服务")
		return
	}

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if args.Bench {

		hey := new(hey.Hey)
		hey.Url = args.Url
		hey.Method = "post"
		hey.Num = args.Nrun
		hey.Con = args.Ncon
		hey.Body = rpcString.(string)
		if args.Type == "yar" {
			hey.ContentType = "application/json"
		}else if args.Type == "hprose"{
			hey.ContentType = "application/hprose"
		}else{
			hey.ContentType = "application/json"
		}
		hey.RunHey()
	}else {
		elapsed := time.Since(t1)
		if args.Format == true {
			is_map := FormatResutl(rpcString, 1)
			if is_map == true {
				fmt.Println("result:\r\n", buf.String())
			} else {
				fmt.Println("result:", buf.String())
			}
		} else {
			fmt.Println("result:", rpcString)
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
	reflectValue := reflect.ValueOf(result)
	if reflectValue.Kind() == reflect.Map {
		is_map = true
		switchTypeWrite(result, i)
	}else if reflectValue.Kind() == reflect.Slice{
		len := reflectValue.Len()
		for ii := 0; ii < len; ii++ {
			tmp := reflect.ValueOf(result).Index(ii).Interface()
			switchTypeWrite(tmp, i)
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


func switchTypeWrite(switchValue interface{}, i int){

	WriteString(func(){
		buf.WriteString("[\n")
	}, i)

	switch resultValue := switchValue.(type){
	case map[string]interface{}:
		for k, vv := range resultValue {
			if reflect.ValueOf(vv).Kind() == reflect.Map {
				WriteString(func(){
					buf.WriteString(fmt.Sprintf("   '%s' => \r\n", k))
				}, i)
				FormatResutl(vv, i + 1)

			}else if reflect.ValueOf(vv).Kind() == reflect.String {

				WriteString(func(){
					buf.WriteString(fmt.Sprintf("   '%s' => '%s',\r\n", k, vv))
				}, i)

			}else if reflect.ValueOf(vv).Kind() == reflect.Slice {
				WriteString(func(){
					buf.WriteString(fmt.Sprintf("   '%s' => \r\n", k))
				}, i)
				FormatResutl(vv, i + 1)

			} else {

				WriteString(func(){
					buf.WriteString(fmt.Sprintf("  '%s' => %s,\r\n", k, vv))
				}, i)

			}
		}
	case map[interface{}]interface{}:
		for k, vv := range resultValue {

			if reflect.ValueOf(vv).Kind() == reflect.Map {
				WriteString(func(){
					buf.WriteString(fmt.Sprintf("   '%s' => \r\n", k))
				}, i)
				FormatResutl(vv, i + 1)

			}else if reflect.ValueOf(vv).Kind() == reflect.String {

				WriteString(func(){
					buf.WriteString(fmt.Sprintf("   '%s' => '%s',\r\n", k, vv))
				}, i)

			}else if reflect.ValueOf(vv).Kind() == reflect.Int {

				WriteString(func(){
					buf.WriteString(fmt.Sprintf("   '%s' => '%d',\r\n", k, vv))
				}, i)

			}else if reflect.ValueOf(vv).Kind() == reflect.Slice {
				WriteString(func(){
					buf.WriteString(fmt.Sprintf("   '%s' => \r\n", k))
				}, i)
				FormatResutl(vv, i + 1)

			} else {

				WriteString(func(){
					buf.WriteString(fmt.Sprintf("  '%s' => %s,\r\n", k, vv))
				}, i)

			}
		}
	case string:
		buf.WriteString(fmt.Sprintf("'%s'", resultValue))
	case int:
		buf.WriteString(fmt.Sprintf("'%d'", resultValue))
	default:
		fmt.Printf("Unknow Type:%+T\n", resultValue)
	}

	WriteString(func() {
		buf.WriteString(" ],\r\n")
	}, i)
}

/**
*解析传入的参数
*@param inputArgs []string  外部传入的参数
 */
func GetArgs(inputArgs []string) []interface{} {
	s := make([]interface{}, len(inputArgs))
	for i, v := range inputArgs {
		//解析数组
		if strings.Contains(strings.ToLower(v), "arrfile:") {
			tmp := strings.Split(v, ":")
			jsonData, err := ReadJson(tmp[1])
			if err != nil {
				panic(err.Error())
			}

			b := []byte(jsonData)
			var m interface{}
			if err := json.Unmarshal(b, &m); err != nil {
				panic(err.Error())
			}
			s[i] = m
		} else if strings.Contains(strings.ToLower(v), "arr:") {
			replaceStr := strings.Replace(v, "arr:", "", 1)
			replaceStr = strings.Trim(replaceStr, "#")
			splitStr := strings.Split(replaceStr, "#")
			arrMap := make(map[string]string, len(splitStr))
			si := 0
			for _, spV := range splitStr {
				if strings.Contains(spV, "=") {
					spspV := strings.Split(spV, "=")
					arrMap[spspV[0]] = spspV[1]
				}else{
					arrMap[strconv.Itoa(si)] = spV
					si++
				}
			}
			s[i] = arrMap
		} else if strings.HasPrefix(strings.ToLower(v), "i:") {
			tmp := strings.Split(v, ":")
			intString, err := strconv.Atoi(tmp[1])
			if err != nil {
				panic(err.Error())
			}else {
				s[i] = intString
			}
		} else {
			s[i] = v
		}
	}

	return s
}
