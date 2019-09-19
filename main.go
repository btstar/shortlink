package main

import (
	"flag"
	"fmt"
	"github.com/fonzie1006/shortlink/pkg/gredis"
	"github.com/fonzie1006/shortlink/pkg/setting"
	"github.com/fonzie1006/shortlink/routers"
	logging "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const (
	point = `  _________.__                   __  .____    .__        __               ___________                  .__        
 /   _____/|  |__   ____________/  |_|    |   |__| ____ |  | __           \_   _____/___   ____ _______|__| ____  
 \_____  \ |  |  \ /  _ \_  __ \   __\    |   |  |/    \|  |/ /   ______   |    __)/  _ \ /    \\___   /  |/ __ \ 
 /        \|   Y  (  <_> )  | \/|  | |    |___|  |   |  \    <   /_____/   |     \(  <_> )   |  \/    /|  \  ___/ 
/_______  /|___|  /\____/|__|   |__| |_______ \__|___|  /__|_ \            \___  / \____/|___|  /_____ \__|\___  >
        \/      \/                           \/       \/     \/                \/             \/      \/       \/  `
)

func ServerInitAndListen() {
	newMux := http.NewServeMux()
	router := routers.RouterInit()

	newMux.Handle("/", router)

	endpoint := fmt.Sprintf("%s:%d", setting.ServerSetting.Host, setting.ServerSetting.Port)
	maxHeaderBytes := 1 << 20

	server := http.Server{
		Addr:           endpoint,
		Handler:        router,
		ReadTimeout:    setting.ServerSetting.ReadTimeout,
		WriteTimeout:   setting.ServerSetting.WriteTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	logging.Infof("The http server starts and listens on http://%s", endpoint)
	err := server.ListenAndServe()
	if err != nil {
		logging.Fatalf("There was an error starting the server, err: %v", err)
	}

}

func main() {
	configMethod := flag.String("config", "config/config-example.ini", "Configuration file required for service startup")
	flag.Parse()
	fmt.Println(point)
	time.Sleep(1 * time.Second)

	// 初始化服务配置
	setting.Setup(*configMethod)
	// 初始化redis链接
	err := gredis.Setup()
	if err != nil {
		logging.Fatalf("Redis init err : %v", err)
	}

	// 启动http服务并监听
	ServerInitAndListen()
}
