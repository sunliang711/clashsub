package handlers

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	tmpl     *template.Template
	inbounds Inbounds
)

type Shadowsocks struct {
	User     string
	Tag      string
	Server   string
	Port     string
	Cipher   string
	Password string
	UDP      bool
	Sub      bool
}

type Vmess struct {
	User    string
	Tag     string
	Server  string
	Port    string
	UUID    string
	AlterID string
	Network string
	Sub     bool
}

type Socks5 struct {
	User     string
	Tag      string
	Server   string
	Port     string
	UDP      bool
	Auth     string
	Username string
	Password string
	Sub      bool
}

type Http struct {
	User     string
	Tag      string
	Server   string
	Port     string
	Username string
	Password string
	Sub      bool
}

type Inbounds struct {
	Vmess       []Vmess
	Shadowsocks []Shadowsocks
	Http        []Http
	Socks5      []Socks5
}

func Init() {
	var err error
	tmpl, err = template.ParseFiles(viper.GetString("server.template"))
	if err != nil {
		logrus.Fatalf("parse template filer: %v error: %v", viper.GetString("template.template"), err)
	}
	err = viper.UnmarshalKey("inbounds", &inbounds)
	if err != nil {
		msg := fmt.Sprintf("Unmarshal inbounds error: %v", err)
		logrus.Fatalf(msg)
	}
}

func Sub(ctx *gin.Context) {
	user := ctx.Param("user")
	if len(user) == 0 {
		msg := fmt.Sprintf("no user specified")
		logrus.Error(msg)
		ctx.String(200, msg)
		return
	}

	// Port
	var (
		http        []Http
		vmess       []Vmess
		shadowsocks []Shadowsocks
		socks5      []Socks5
	)
	for i := range inbounds.Http {
		fields := strings.Split(inbounds.Http[i].Tag, ":")
		if len(fields) < 3 {
			logrus.Fatalf("Http tag format incorrect")
		}
		inbounds.Http[i].Port = fields[1]

		if inbounds.Http[i].User == user && inbounds.Http[i].Sub {
			http = append(http, inbounds.Http[i])
		}
	}
	for i := range inbounds.Vmess {
		fields := strings.Split(inbounds.Vmess[i].Tag, ":")
		if len(fields) < 3 {
			logrus.Fatalf("Vmess tag format incorrect")
		}
		inbounds.Vmess[i].Port = fields[1]

		if inbounds.Vmess[i].User == user && inbounds.Vmess[i].Sub {
			vmess = append(vmess, inbounds.Vmess[i])
		}
	}
	for i := range inbounds.Shadowsocks {
		fields := strings.Split(inbounds.Shadowsocks[i].Tag, ":")
		if len(fields) < 3 {
			logrus.Fatalf("Shadowsocks tag format incorrect")
		}
		inbounds.Shadowsocks[i].Port = fields[1]

		if inbounds.Shadowsocks[i].User == user && inbounds.Shadowsocks[i].Sub {
			shadowsocks = append(shadowsocks, inbounds.Shadowsocks[i])
		}
	}
	for i := range inbounds.Socks5 {
		fields := strings.Split(inbounds.Socks5[i].Tag, ":")
		if len(fields) < 3 {
			logrus.Fatalf("Socks5 tag format incorrect")
		}
		inbounds.Socks5[i].Port = fields[1]

		if inbounds.Socks5[i].User == user && inbounds.Socks5[i].Sub {
			socks5 = append(socks5, inbounds.Socks5[i])
		}
	}

	if len(vmess) == 0 && len(shadowsocks) == 0 && len(socks5) == 0 && len(http) == 0 {
		msg := fmt.Sprintf("No such sub with user: %v", user)
		logrus.Errorf(msg)
		ctx.String(200, msg)
		return
	}

	subInbounds := Inbounds{
		Vmess:       vmess,
		Shadowsocks: shadowsocks,
		Socks5:      socks5,
		Http:        http,
	}

	logrus.Debugf("inbounds: %+v", subInbounds)

	var buf bytes.Buffer
	err := tmpl.Execute(&buf, &subInbounds)
	if err != nil {
		msg := fmt.Sprintf("execute template error: %v", err)
		logrus.Error(msg)
		ctx.String(200, msg)
		return
	}
	s := buf.String()
	ctx.String(200, s)
}
