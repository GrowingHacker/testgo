package xraycore

import (
	"github.com/xtls/xray-core/core"
)

type XrayService struct {
	instance *core.Instance
}

// 初始化并启动 Xray
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
