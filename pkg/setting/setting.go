package setting

import (
	"encoding/json"
	"fmt"
	"github.com/go-ini/ini"
	logging "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
	"time"
)

type Config struct {
	Server *Server `yaml:"server" json:"server"`
	Redis  *Redis  `yaml:"redis" json:"redis"`
}

type Server struct {
	Host         string        `yaml:"host" json:"host"`
	Port         int           `yaml:"port" json:"port"`
	ReadTimeout  time.Duration `yaml:"readTimeout" json:"readTimeout"`
	WriteTimeout time.Duration `yaml:"writeTimeout" json:"writeTimeout"`
}

type Redis struct {
	DB          int           `yaml:"db" json:"db"`
	Host        string        `yaml:"host" json:"host"`
	Password    string        `yaml:"password" json:"password"`
	MaxIdle     int           `yaml:"maxIdle" json:"maxIdle"`
	MaxActive   int           `yaml:"maxActive" json:"maxActive"`
	IdleTimeout time.Duration `yaml:"idleTimeout" json:"idleTimeout"`
}

// 服务器配置
var ServerSetting = &Server{}

// Redis配置
var RedisSetting = &Redis{
	DB:          0,
	Host:        "127.0.0.1:6379",
	Password:    "",
	MaxIdle:     30,
	MaxActive:   100,
	IdleTimeout: 10 * time.Second,
}

var cfg *ini.File

func Setup(configMethod string) {
	// 整合结构体
	var configSetting = &Config{
		Server: ServerSetting,
		Redis:  RedisSetting,
	}

	// 如果需要配置日志模块其他操作,比如文件存储是否标准输出,详细配置请参考:https://github.com/sirupsen/logrus.git
	loggingInit()
	logging.Infof("Ready to read the configuration file: %s", configMethod)
	logging.Infof("Start initializing the service configuration....")

	var configPath = configMethod

	if strings.HasSuffix(configPath, ".yaml") {
		yamlInitConfig(configPath, configSetting)
	} else if strings.HasSuffix(configPath, ".json") {
		jsonInitConfig(configPath, configSetting)
	} else if strings.HasSuffix(configPath, ".ini") {
		iniInitConfig(configPath, configSetting)
	} else {
		logging.Fatalf("Unknown file name suffix")
	}

	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second
	RedisSetting.IdleTimeout = RedisSetting.IdleTimeout * time.Second

	fmt.Printf("%+v\n", ServerSetting)
	logging.Infof("Server configuration initialized successfully....")
}

func loggingInit() {

}

// json方式的配置文件
func jsonInitConfig(confPath string, configSetting *Config) {
	var err error
	jsonFile, err := ioutil.ReadFile(confPath)
	if err != nil {
		logging.Fatalf("Read Json config %s err : %v", confPath, err)
	}

	err = json.Unmarshal(jsonFile, configSetting)
	if err != nil {
		logging.Fatalf("Json unmarshal err : %v", err)
	}
}

// yaml方式的配置文件
func yamlInitConfig(confPath string, configSetting *Config) {
	var err error
	yamlFile, err := ioutil.ReadFile(confPath)
	if err != nil {
		logging.Fatalf("Read Yaml config %s err : %v", confPath, err)
	}

	err = yaml.Unmarshal(yamlFile, configSetting)
	if err != nil {
		logging.Printf("Yaml unmarshal err : %v", err)
	}

}

// ini方式的配置文件
func iniInitConfig(confPath string, configSetting *Config) {
	var err error

	cfg, err = ini.Load(confPath)
	if err != nil {
		logging.Fatalf("setting.iniInitConfig, fail to parse %s", confPath)
	}

	mapTo("redis", configSetting.Redis)
	mapTo("server", configSetting.Server)

}

func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		logging.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}
