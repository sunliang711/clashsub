package models

import (
	"context"
	"database/sql"

	// register mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"

	"github.com/sunliang711/goutils/mongodb"
	umysql "github.com/sunliang711/goutils/mysql"

	//"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/spf13/viper"
)

var (
	mysqlConn *sql.DB
	mongoConn *mongo.Client
)

//InitMysql open mysql with dsn
func InitMysql(dsn string) {
	logrus.Infof("Try to connect to mysql: '%v'", dsn)
	var err error
	mysqlConn, err = umysql.New(dsn, viper.GetInt("mysql.maxIdleConns"), viper.GetInt("mysql.maxOpenConns"))
	if err != nil {
		panic(err)
	}
	logrus.Infof("Connected to mysql: '%v'", dsn)
}

//CloseMysql close mysql connection
func CloseMysql() {
	logrus.Infoln("Close mysql.")
	mysqlConn.Close()
}

// InitMongo opens a mongodb connection
func InitMongo(url string) {
	logrus.Infof("Try to connect to mongodb: '%v'", url)
	var err error
	mongoConn, err = mongodb.New(url, 5)
	if err != nil {
		panic(err)
	}
	logrus.Infof("Connected to mongodb: '%v'", url)
}

func CloseMongo() {
	logrus.Infof("Close mongodb")
	mongoConn.Disconnect(context.Background())
}
