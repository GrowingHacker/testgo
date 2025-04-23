package main

import (
	"encoding/json"
	"fmt"
	"log"
	"mymodule/xraycoreHelper"
	"os"

	"github.com/gin-gonic/gin"
)

var xray = &xraycoreHelper.XrayService{}

func main() {
	r := gin.Default()

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
