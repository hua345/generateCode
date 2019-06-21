package config

import (
	"codegen/pkg/util"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

type Server struct {
	DefaultHttpMethod string            `yaml:"defaultHttpMethod"`
	DefaultHttpPort   string            `yaml:"defaultHttpPort"`
	ApiBaseUrl        string            `yaml:"apiBaseUrl"`
	AuthorName        string            `yaml:"authorName"`
	Database          Database          `yaml:"database"`
	Redis             Redis             `yaml:"redis"`
	Springboot        SpringbootSetting `yaml:"springboot"`
}

func newServer() Server {
	return Server{
		DefaultHttpMethod: "com.github",
		DefaultHttpPort:   "api/v1",
		ApiBaseUrl:        "api/v1",
		AuthorName:        "learn",
	}
}

type SpringbootSetting struct {
	GroupId           string `yaml:"groupId"`
	ArtifactId        string `yaml:"artifactId"`
	SupportMaven      bool   `yaml:"supportMaven"`
	SupportGradle     bool   `yaml:"supportGradle"`
	SupportDocker     bool   `yaml:"supportDocker"`
	SupportI18n       bool   `yaml:"supportI18n"`
	SupportDataSource string `yaml:"supportDataSource"`
	SupportSwagger    bool   `yaml:"supportSwagger"`
}

func newSpringbootSetting() SpringbootSetting {
	return SpringbootSetting{
		GroupId:           "com.github",
		ArtifactId:        "learn",
		SupportMaven:      true,
		SupportGradle:     true,
		SupportDocker:     true,
		SupportI18n:       true,
		SupportDataSource: "druid",
		SupportSwagger:    true,
	}
}

type Database struct {
	Type         string `yaml:"type"`
	Host         string `yaml:"host"`
	DatabaseName string `yaml:"databaseName"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
}

type Redis struct {
	Host        string        `yaml:"host"`
	Password    string        `yaml:"password"`
	MaxIdle     int           `yaml:"maxIdle"`
	MaxActive   int           `yaml:"maxActive"`
	IdleTimeout time.Duration `yaml:"idleTimeout"`
}

type GinSetting struct {
}

var ServerConfig Server
var DefaultSettingFile = DefaultConfigFile

// Setup initialize the configuration instance
func Setup(configPath string) {
	ServerConfig = newServer()
	path := DefaultSettingFile
	if len(configPath) > 0 {
		path = configPath
	}
	configExist, _ := util.PathExists(path)
	if true != configExist {
		log.Println("配置文件" + DefaultConfigFile + "没有找到!")
		log.Println("初始化配置文件命令: codegen config init")
		os.Exit(-1)
	}
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err.Error())
	}
	// set default
	ServerConfig.Springboot = newSpringbootSetting()
	// 检查数据库类型是否支持

	err = yaml.Unmarshal(yamlFile, &ServerConfig)
	if err != nil {
		fmt.Println(err.Error())
	}
	ServerConfig.DefaultHttpMethod = strings.ToLower(ServerConfig.DefaultHttpMethod)
	if !(httpMethodGet == ServerConfig.DefaultHttpMethod ||
		httpMethodPost == ServerConfig.DefaultHttpMethod ||
		httpMethodUpdate == ServerConfig.DefaultHttpMethod ||
		httpMethodPut == ServerConfig.DefaultHttpMethod ||
		httpMethodDelete == ServerConfig.DefaultHttpMethod) {
		log.Println("支持的http请求方法:" + httpMethodGet + "\\" +
			httpMethodPost + "\\" + httpMethodPut + "\\" + httpMethodDelete)
		os.Exit(-1)
	}
	ServerConfig.Database.Type = strings.ToLower(ServerConfig.Database.Type)
	if !(DBTypeMysql == ServerConfig.Database.Type ||
		DBTypeMariadb == ServerConfig.Database.Type ||
		DBTypePostgresql == ServerConfig.Database.Type) {
		log.Println("支持的数据库类型:" + DBTypeMysql + "\\" +
			DBTypeMariadb + "\\" + DBTypePostgresql)
		os.Exit(-1)
	}
	if !(DataSourceDruid == ServerConfig.Springboot.SupportDataSource ||
		DataSourceHikariCP == ServerConfig.Springboot.SupportDataSource) {
		log.Println("支持的数据库连接池:" + DataSourceDruid + "\\" +
			DataSourceHikariCP)
		os.Exit(-1)
	}
}
