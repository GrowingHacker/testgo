package xraycoreHelper

import (
	// JSON 配置註冊器
	_ "github.com/xtls/xray-core/main/json"

	//出入站管理配置组
	_ "github.com/xtls/xray-core/app/proxyman"
	_ "github.com/xtls/xray-core/app/proxyman/inbound"
	_ "github.com/xtls/xray-core/app/proxyman/outbound"

	// 協議（按需添加）
	_ "github.com/xtls/xray-core/proxy/freedom"
	_ "github.com/xtls/xray-core/proxy/vmess"
	_ "github.com/xtls/xray-core/proxy/vmess/inbound"
	_ "github.com/xtls/xray-core/proxy/vmess/outbound"

	// 傳輸層（按你配置中出現的來添加）
	_ "github.com/xtls/xray-core/transport/internet/tcp"
	_ "github.com/xtls/xray-core/transport/internet/tls"
	_ "github.com/xtls/xray-core/transport/internet/websocket"

	//系统包
	"bytes"
	"log"
	"os"

	//core 核心包
	"github.com/xtls/xray-core/core"
)

type XrayService struct {
	instance *core.Instance
}

// 从文件加载配置并启动 Xray
func (x *XrayService) StartFromFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	//cfg, err := core.LoadConfig("json", f)
	//cfg, err := core.LoadConfig("auto", cmdarg.Arg{path})
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	cfg, err := core.LoadConfig("json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	instance, err := core.New(cfg)
	if err != nil {
		return err
	}
	x.instance = instance
	return x.instance.Start()

}

// 使用配置启动 Xray
func (x *XrayService) Start(config string) error {
	cfg, err := core.LoadConfig("json", []byte(config))
	if err != nil {
		return err
	}
	instance, err := core.New(cfg)
	if err != nil {
		return err
	}
	x.instance = instance
	return x.instance.Start()
}

// 关闭 Xray 实例
func (x *XrayService) Stop() {
	if x.instance != nil {
		x.instance.Close()
	}
}
