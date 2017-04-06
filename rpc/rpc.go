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

	if args.Type == "yar" {
		Yar(args)
	}else{
		Hprose(args)
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

		WriteString(func(){
			buf.WriteString("[\n")
		}, i)

		switchTypeWrite(result, i)
		if i == 1 {
			WriteString(func() {
				buf.WriteString(" ],\r\n")
			}, i)
		}else {
			WriteString(func() {
				buf.WriteString("],\r\n")
			}, i)
		}
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
	default:
		fmt.Println("unknown type", resultValue)
	}

	WriteString(func() {
		buf.WriteString(" ],\r\n")
	}, i)
}
