package handlers

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	tmpl *template.Template
)

func InitTemplate() {
	var err error
	tmpl, err = template.ParseFiles(viper.GetString("server.template"))
	if err != nil {
		logrus.Fatalf("parse template filer: %v error: %v", viper.GetString("template.template"), err)
	}
}

type serverConfig struct {
	Name     string
	Type     string
	Server   string
	Port     string
	Cipher   string
	Password string
	UDP      bool
}

func Sub(ctx *gin.Context) {
	port := ctx.Param("port")
	if len(port) == 0 {
		msg := fmt.Sprintf("no port specified")
		logrus.Error(msg)
		ctx.String(200, msg)
		return
	}

	// if viper.GetString("configs."+port) == "" {
	// 	msg := fmt.Sprintf("no such port: %v", port)
	// 	logrus.Errorf(msg)
	// 	ctx.String(200, msg)
	// 	return
	// }

	var serverConfigs []serverConfig
	err := viper.UnmarshalKey("configs."+port, &serverConfigs)
	if err != nil {
		msg := fmt.Sprintf("Unmarshal configs.%v error: %v", port, err)
		logrus.Errorf(msg)
		ctx.String(200, msg)
		return
	}
	if len(serverConfigs) == 0 {
		msg := fmt.Sprintf("No config with: %v", port)
		logrus.Warnf(msg)
		ctx.String(200, msg)
		return
	}
	// cfg := &serverConfig{
	// 	Server:   viper.GetString("configs." + port + ".server"),
	// 	Port:     port,
	// 	Cipher:   viper.GetString("configs." + port + ".cipher"),
	// 	Password: viper.GetString("configs." + port + ".password"),
	// }
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, serverConfigs)
	if err != nil {
		msg := fmt.Sprintf("execute template error: %v", err)
		logrus.Error(msg)
		ctx.String(200, msg)
		return
	}
	// ctx.Header("Content-Type", "text/yaml")
	// ctx.Header("Connection", "keep-alive")
	s := buf.String()
	// ctx.Header("Content-Length", fmt.Sprintf("%d", len(s)))
	ctx.String(200, s)
}
