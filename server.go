package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"mymodule/xraycore"
  "os"
)

var xray = &xraycore.XrayService{}

func main() {
	r := gin.Default()

	// 启动接口
	r.POST("/start", func(c *gin.Context) {
		configJson := `{
  "inbounds": [
    {
      "port": 443,
      "listen": "127.0.0.1",
      "protocol": "vmess",
      "settings": {
        "clients": [
          {
            "id": "5b7a1f37-02e6-4eab-8f52-7d8be39bece0",
            "alterId": 64,
            "security": "auto"
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "headers": {
            "Host": "tunnel.honphiewon.eu.org"
          }
        },
        "tlsSettings": {
          "certificates": [
            {
              "certificateFile": "/usr/local/xray/cert.pem",
              "keyFile": "/usr/local/xray/private.key"
            }
          ]
        }
      }
    }
  ],
  "outbounds": [
    {
      "protocol": "freedom",
      "settings": {}
    }
  ],
  "dns": {
    "servers": [
      "1.1.1.1",
      "1.0.0.1",
      "8.8.8.8",
      "8.8.4.4"
    ]
  },
  "log": {
    "loglevel": "info",
    "access": "/var/log/xray/access.log",
    "error": "/var/log/xray/error.log"
  }
}
`
		err := xray.Start(configJson)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Xray 启动成功"})
	})

	// 停止接口
	r.POST("/stop", func(c *gin.Context) {
		xray.Stop()
		c.JSON(200, gin.H{"message": "Xray 已停止"})
	})


  port := os.Getenv("PORT")
  log.Println("端口分配信息：" + port)
  if port == "" {
    port = "5703"
  }
  log.Println("Web 服务启动于 :" + port)
  r.Run(":" + port)

	// r.Run(":80")
}
