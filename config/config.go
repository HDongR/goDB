package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var Database string = "mysql"
var Owner string = ""
var ConnectionString string
var Port string
var TempPath string
var MailServer string
var MailSender string
var ServiceUrl string

var UploadDir string

var SmsUser string
var SmsKey string
var SmsSender string

var AdminEmail string

var (
	Version string
	Build   string
	DEBUG   uint64
)

func init() {
	DEBUG = 0
	if os.Getenv("GIN_MODE") == "release" {
		DEBUG = 0
	}
	if DEBUG > 0 {
		fmt.Printf("Debug: MODE=true, flag=%+v \n", DEBUG)
	}

	if value := viper.Get("connectionString"); value != nil {
		ConnectionString = value.(string)
	}

	if value := viper.Get("mailServer"); value != nil {
		MailServer = value.(string)
	}

	if value := viper.Get("mailSender"); value != nil {
		MailSender = value.(string)
	}

	if Port == "" {
		Port = "80"
	}

	if value := viper.Get("smsSender"); value != nil {
		SmsSender = value.(string)
	}

	if value := viper.Get("adminEmail"); value != nil {
		AdminEmail = value.(string)
	}

	Database = "mysql"
	ConnectionString = "root:@tcp(127.0.0.1:3306)/crawler"
}
