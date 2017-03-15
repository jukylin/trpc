# trpc

####　RPC 调试工具，用于调试远程RPC接口，暂只支持yar和HTTP协议

## 安装

* 下载 godep
```
$ go get github.com/tools/godep
```
* 下载 trpc
```
$ git@github.com:Aqiling/trpc.git
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
$ derpc -u http://www.test.com -f test -a 1 -a 4 -a array:localfile.json
```

## 返回：
    result: 5
    runtime:  98.39678ms
    
## 注意：

&&&&& 参数按照函数参数传递，如果为数组，上例第三个"$c"，
需要把数组json化后放入"localfile.json"，再执行命令。
