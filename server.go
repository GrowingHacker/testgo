package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mymodule/xraycoreHelper"
	"net"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var xray = &xraycoreHelper.XrayService{}

func getPublicIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// 获取服务器首个非环回 IPv4
func getServerIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		// 只处理 IP 网络地址
		if ipNet, ok := addr.(*net.IPNet); ok {
			ip := ipNet.IP
			// 跳过环回、IPv6
			if !ip.IsLoopback() && ip.To4() != nil {
				return ip.String(), nil
			}
		}
	}
	return "", fmt.Errorf("no non-loopback IPv4 address found")
}

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		// 获取服务器 IP
		serverIP, err := getServerIP()
		if err != nil {
			serverIP = "unknown"
		}
		// 获取公网 IP
		publicIP, err := getPublicIP()
		if err != nil {
			publicIP = "unknown"
		}

		// 取监听端口
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}

		// 构造 HTML
		html := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="zh-CN">
<head><meta charset="UTF-8"><title>服务器信息</title></head>
<body>
  <p style="font-size:24px; color:#e91e63; font-weight:bold;">
    您好，服务器内网IP，公网IP 和 端口：<span style="color:#3f51b5">%s,%s:%s</span>
  </p>
</body>
</html>`, serverIP, publicIP, port)

		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
	})

	// 启动接口
	r.GET("/start", func(c *gin.Context) {
		// 获取 Render 分配的端口
		port := os.Getenv("PORT")
		if port == "" {
			port = "5703"
		}
		fmt.Println("Render 端口:", port)

		// 动态读取并修改配置文件中的端口
		configPath := "./config/test.json"
		modifiedConfigPath := "./config/test_runtime.json"

		err := patchXrayPort(configPath, modifiedConfigPath, port)
		if err != nil {
			fmt.Printf("配置修改失败: %+v\n", err)
			c.JSON(500, gin.H{"error": "修改配置失败"})
			return
		}

		// 启动 Xray
		err = xray.StartFromFile(modifiedConfigPath)
		if err != nil {
			fmt.Printf("启动失败: %+v\n", err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "Xray 启动成功"})
	})

	// 停止接口
	r.GET("/stop", func(c *gin.Context) {
		xray.Stop()
		c.JSON(200, gin.H{"message": "Xray 已停止"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "5703"
	}
	log.Println("Web 服务启动于 :" + port)
	r.Run(":" + port)
}

// 修改配置文件中的端口
func patchXrayPort(inputPath, outputPath, port string) error {
	// 使用 os.ReadFile 替代 ioutil.ReadFile
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	// 反序列化配置为 map
	var config map[string]interface{}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return err
	}

	// 遍历 inbounds，替换端口
	if inbounds, ok := config["inbounds"].([]interface{}); ok {
		for _, inbound := range inbounds {
			if ib, ok := inbound.(map[string]interface{}); ok {
				ib["port"] = port // 修改端口
			}
		}
	}

	// 写入修改后的配置
	newData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	// 使用 os.WriteFile 替代 ioutil.WriteFile
	return os.WriteFile(outputPath, newData, 0644)
}
