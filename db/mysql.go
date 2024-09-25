package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"hotNews/http/models"
	"hotNews/utils"
	"log"
	"os"
)

var DbEngin *gorm.DB

func init() {
	var err error
	var config = utils.MysqlSetting
	dsName := config.Username + ":" + config.Password + "@(" + config.Host + ":" + config.Port + ")/" +
		config.Database + "?charset=" + config.Charset +
		"&parseTime=True&loc=Local"
	DbEngin, err = gorm.Open("mysql", dsName)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	DbEngin.AutoMigrate(&model.Application{}, &model.Article{})
	DbEngin.LogMode(true)
	fmt.Println("init database ok")
}

func DbClose() {
	DbEngin.Close()
}
