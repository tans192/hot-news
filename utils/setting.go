package utils

import (
	"fmt"
	"github.com/go-ini/ini"
	"log"
	"os"
)

//db配置文件
//type mongoSetting struct {
//	Host     string `ini:"host"`
//	Port     string `ini:"port"`
//	Database string `ini:"database"`
//	Username string `ini:"username"`
//	Password string `ini:"password"`
//	Charset  string `ini:"charset"`
//}
//
//var MongoSetting = &mongoSetting{}

// redis 配置文件
type rdSetting struct {
	Host     string `ini:"redis_host"`
	User     string `ini:"redis_user"`
	Password string `ini:"redis_password"`
	Port     int    `ini:"redis_port"`
	Db       int    `ini:"redis_db"`
}

var RedisSetting = &rdSetting{}

// app 配置文件
type appSetting struct {
	Protocol  string `ini:"protocol"`
	Port      string `ini:"http_port"`
	AppMode   string `ini:"app_mode"`
	DebugMode string `ini:"debug_mode"`
}

var AppSetting = &appSetting{}

//db配置文件
type mysqlSetting struct {
	Host      string `ini:"host"`
	Port      string `ini:"port"`
	Database  string `ini:"database"`
	Username  string `ini:"username"`
	Password  string `ini:"password"`
	Charset   string `ini:"charset"`
	Collation string `ini:"collation"`
}

var MysqlSetting = &mysqlSetting{}

var cfg *ini.File

func init() {
	pwd, _ := os.Getwd()
	var err error
	cfg, err = ini.Load(pwd + "/app.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	//mapTo("mongo", MongoSetting)
	mapTo("server", AppSetting)
	mapTo("redis", RedisSetting)
	mapTo("mysql", MysqlSetting)
	fmt.Println("init app setting ok")
}

func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo "+section+" err: %v", err)
	}
}
