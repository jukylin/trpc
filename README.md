# trpc

### RPC 调试工具，用于调试远程RPC接口，暂只支持yar和HTTP协议

## 安装

* 下载 godep
```
$ go get github.com/tools/godep
```
* 下载 trpc
```
$ git clone git@gitlab.etcchebao.cn:go_service/trpc.git
$ cd trpc
$ godep go build .
$ mv trpc /usr/local/bin/trpc
```

## 例子：
```php
   # URL : http://www.test.com
   public function test($a, $b, $c = []){
       file_put_content("./log.log", json_encode($c));
       return $a + $b;
   }
```

## 执行：
```
$ trpc -u http://www.test.com -f test -a 1 -a 4 -a arrfile:localfile.json
```

## 返回：
    result: 5
    runtime:  98.39678ms
    
## 注意：

参数按照函数参数顺序传递，如果为数组提供2种传递方式：
* 1：-a arr:name=trpc#age=20，通过"#"把key=>val连接起来，
组成["name"=>"trpc","age" => 20]。
* 2：-a arrfile:./localfile.json，对于传递复杂的数组，
需要把数组json化后，放入"./localfile.json"，再执行命令。

## 压力测试 [-b]
测试使用[hey](https://github.com/rakyll/hey)，一款用go重写的ab压力测试工具。


## 格式化返回结果 [-m]
接口返回的结果是数组，map，使用go直接输出，对开发者不友好。

加上[-m]，可对结果美化：
```
 [
   'code' =>
      [
         'age' => '10',
         'name' => 'test',
      ],
   'name' => '1',
 ],
```````````````