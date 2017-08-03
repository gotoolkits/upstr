# upStr
`针对Docker服务注册与发现后Api网关动态加载(Nginx Upstream)处理工具。
使用Consul KV与Service功能作为服务存储后端，动态加载Upstream实例，并
加载配置。`

## 安装

 >  依赖go环境，请确保go环境安装正常
    
     1) go get github.com/gotoolkits/upstr
     2) cd $GOPATH/github.com/gotoolkits/
        go install upstr.go

## API接口

 >  接口： /    
 >  方法： Get  
 >  说明： 返回系统info状态信息  

     {
     "WorkPath": "/usr/local/bin",
     "ConfigPath": "/usr/local/orange/conf/nginx.conf",
     "Consul": "192.168.X.X:8500",
     "KvPath": "paas/ngx/upstream_name?raw",
     "UpstremNum": 0,
     "Updated": 0,
     "Error": 0,
     "LastUpdate": ""
      }

 >  接口： /list   
 >  方法： Get     
 >  说明： 返回当前Nginx服务已配置的upstream列表

     {
       "app": [
         "127.0.0.1:8001"
       ],
       "bkoffice": [
         "127.0.0.1:8001"
       ],
        "default_upstream": [
         "127.0.0.1:8001"
       ],
       "grafana": [
         "127.0.0.1:8001"
       ]
     }


 >  接口： /reload   
 >  方法： Get      
 >  说明： 拉取KV配置更新upstream，Reload配置

    {
       "status": "successful",
       "errCount": 0,
       "updateConut": 4,
       "UpdateTime": "2017-08-03 09:41:56"
    }

 >  接口： /status   
 >  方法： Get    
 >  说明： 获取upstr服务自身状态，可用于监控

    ok


## 配置文件

>  默认查找路径：
>         1.   当前目录config.json (与执行文件同目录) `优先`
>         2.   /etc/upstr/config.json

    {
       "setting": {
          "port":"18083"
       },
       "consul":{
          "host":"192.168.Ｘ.Ｘ:8500"
       },
      "orange":{
         "work_path":"/usr/local/bin",
         "config_path":"/usr/local/orange/conf/nginx.conf"
       }
     }

## Todo
 >   1)  API接口访问安全功能
 >   2)  启动命令参数

 
