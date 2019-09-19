# ShortLink

----

短链接服务,项目使用GOlang编写,使用原生的```net/http```作为Web服务器,由于默认的```net/http```对于路由支持的不是很好，所以这里选择了使用
[mux](github.com/gorilla/mux)作为路由，配合[alice](github.com/justinas/alice)完成中间件工作.整个项目主要的接口有通过短链接获取长链
接，通过长链接获取短链接,最后一个就是Redirect 默认302,相对短链接进行访问基数.

----

### 使用方法

> 长链接转短链接

```bash
➜  ~ curl -i  -X POST http://127.0.0.1:8080/shortlink\?url\=https://www.google.com\&timeout\=-1
HTTP/1.1 200 OK
Content-Type: application/json
Date: Tue, 17 Sep 2019 14:54:04 GMT
Content-Length: 142

{"Code":200,"Data":{"longlink":"https://www.google.com","md5":"8ffdefbdec956b595d257f0aaeefd623","shortlink":"/ZZZZ","timeout":-1},"Msg":"ok"}
```

> 短链接转长链接
```
➜  ~ curl -i http://127.0.0.1:8080/shortlink/ZZZZ
HTTP/1.1 200 OK
Content-Type: application/json
Date: Tue, 17 Sep 2019 14:55:13 GMT
Content-Length: 142

{"Code":200,"Data":{"longlink":"https://www.google.com","md5":"8ffdefbdec956b595d257f0aaeefd623","shortlink":"/ZZZZ","timeout":-1},"Msg":"ok"}
```


如果不想暴露短链接转长链接的API可以直接删除即可

> 短链接跳转并且计数

```
➜  ~ curl -i http://127.0.0.1:8080/ZZZZ
HTTP/1.1 302 Found
Content-Type: text/html; charset=utf-8
Location: https://www.google.com
Date: Tue, 17 Sep 2019 14:56:45 GMT
Content-Length: 45

<a href="https://www.google.com">Found</a>.
```

### 配置服务

默认情况下使用```.ini```的配置文件配置内容很简单

```bash
[redis]
# redis的地址加端口
Host = 127.0.0.1:6379
# 链接redis的密码
Password = 123456
MaxIdle = 40
MacActive = 130
IdleTimeout = 10

[server]
# 服务器启动的端口
Port = 8080
ReadTimeout = 20
WriteTimeout = 20
```


### 运行效果

```bash
go run main.go
```

```bash
  _________.__                   __  .____    .__        __               ___________                  .__
 /   _____/|  |__   ____________/  |_|    |   |__| ____ |  | __           \_   _____/___   ____ _______|__| ____
 \_____  \ |  |  \ /  _ \_  __ \   __\    |   |  |/    \|  |/ /   ______   |    __)/  _ \ /    \\___   /  |/ __ \
 /        \|   Y  (  <_> )  | \/|  | |    |___|  |   |  \    <   /_____/   |     \(  <_> )   |  \/    /|  \  ___/
/_______  /|___|  /\____/|__|   |__| |_______ \__|___|  /__|_ \            \___  / \____/|___|  /_____ \__|\___  >
        \/      \/                           \/       \/     \/                \/             \/      \/       \/
INFO[0001] Start initializing the service configuration....
INFO[0001] Server configuration initialized successfully....
INFO[0001] The http server starts and listens on http://0.0.0.0:8081
```
