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
)

func DebugStart(url string, fn string, args []string){
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

	var ret interface{}

	t1 := time.Now() // get current time

	s := make([]interface{}, len(args))
	for i, v := range args {
		//解析数组
		if strings.Contains(strings.ToLower(v), "array:") {
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
		} else {
			s[i] = v
		}
	}

	callErr := client.Call(fn, &ret, s...)

	if callErr != nil {
		fmt.Println("error", callErr)
	}
	elapsed := time.Since(t1)


	//FormatResutl(ret)
	fmt.Println("result:", ret)
	fmt.Println("runtime: ", elapsed)

}


func FormatResutl(result interface{}) interface{} {
	if reflect.ValueOf(result).Kind() == reflect.Map {

		switch resultValue := result.(type){
		case map[string]string:
			fmt.Println("it's a map, and key \"key\" is", resultValue)
		case map[string]interface{}:
			for k, vv := range resultValue{
				fmt.Println(k, "==>", vv)
			}
			fmt.Println("it's a map, and key \"key\" is", resultValue)
		default:
			fmt.Println("unknown type", resultValue)
		}

	}else{
		//return result
	}

	return 123
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
