package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/yqchilde/wxbot/engine/pkg/log"
	"github.com/yqchilde/wxbot/engine/pkg/net"
	"github.com/yqchilde/wxbot/engine/robot"
	"github.com/yqchilde/wxbot/framework/dean"
	"github.com/yqchilde/wxbot/framework/super"
	"github.com/yqchilde/wxbot/framework/vlw"

	"github.com/gorilla/websocket"

	// 导入插件, 变更插件请查看README
	_ "github.com/yqchilde/wxbot/engine/plugins"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type MessageResp struct {
	Code      int    `json:"Code"`
	Result    string `json:"Result"`
	ReturnStr string `json:"ReturnStr"`
	ReturnInt string `json:"ReturnInt"`
}

func main() {
	var conn *websocket.Conn
	v := viper.New()
	v.SetConfigFile("config.yaml")
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("[main] 读取配置文件失败: %s", err.Error())
	}
	c := robot.NewConfig()
	if err := v.Unmarshal(c); err != nil {
		log.Fatalf("[main] 解析配置文件失败: %s", err.Error())
	}

	f := robot.IFramework(nil)
	switch c.Framework.Name {
	case "Dean":
		f = robot.IFramework(dean.New(c.BotWxId, c.Framework.ApiUrl, c.Framework.ApiToken))
		if ipPort, err := net.CheckoutIpPort(c.Framework.ApiUrl); err == nil {
			if ping := net.PingConn(ipPort, time.Second*10); !ping {
				c.SetConnHookStatus(false)
				log.Warn("[main] 无法连接Dean框架，网络无法Ping通，请检查网络")
			}
		}
	case "VLW", "vlw":
		f = robot.IFramework(vlw.New(c.BotWxId, c.Framework.ApiUrl, c.Framework.ApiToken))
		if ipPort, err := net.CheckoutIpPort(c.Framework.ApiUrl); err == nil {
			if ping := net.PingConn(ipPort, time.Second*10); !ping {
				c.SetConnHookStatus(false)
				log.Warn("[main] 无法连接到VLW框架，网络无法Ping通，请检查网络")
			}
		}
	case "Super", "SUPER", "super":
		url := strings.Replace(c.Framework.ApiUrl, "http", "ws", 1)
		var err error
		conn, _, err = websocket.DefaultDialer.Dial(url, nil)

		if err != nil {
			log.Fatal("无法连接到 WebSocket 服务:", err)
		}

		defer func(conn *websocket.Conn) {
			// 关闭底层 TCP 网络连接
			err := conn.Close()
			if err != nil {
				log.Fatal("close tcp conn error: %v\n", err)
			}
		}(conn)
		f = robot.IFramework(super.New(c.BotWxId, c.BotNickname, c.Framework.ApiUrl, c.Framework.ApiToken, conn, c.ServerPort))
		if ipPort, err := net.CheckoutIpPort(c.Framework.ApiUrl); err == nil {
			if ping := net.PingConn(ipPort, time.Second*10); !ping {
				c.SetConnHookStatus(false)
				log.Warn("[main] 无法连接到SUPER框架，网络无法Ping通，请检查网络")
			}
		}
	default:
		log.Fatalf("[main] 请在配置文件中指定机器人框架后再启动")
	}
	robot.Run(c, f)

}
