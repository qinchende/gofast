package conf

import (
	"fmt"
	"github.com/qinchende/gofast/aid/logx"
	"github.com/qinchende/gofast/core/lang"
	"github.com/qinchende/gofast/store/bind"
	"os"
	"path"
)

// 系统目前支持两种格式的配置文件：
// 1. JSON
// 2. Yaml
var loaders = map[string]func(any, []byte) error{
	".json": LoadFromJson,
	".yaml": LoadFromYaml,
	".yml":  LoadFromYaml,
}

// 必须加载配置，否则应用无法启动，直接退出
func MustLoad(path string, dst any) {
	if err := LoadFile(path, dst); err != nil {
		logx.ErrorFatalF("error: config file %s, %s", path, err.Error())
	}
}

func LoadFile(file string, dst any) error {
	if content, err := os.ReadFile(file); err != nil {
		return err
	} else if loader, ok := loaders[path.Ext(file)]; ok {
		return loader(dst, lang.STB(os.ExpandEnv(string(content))))
	} else {
		return fmt.Errorf("unrecoginized file type: %s", file)
	}
}

func LoadFromJson(dst any, content []byte) error {
	return bind.BindJsonBytes(dst, content, bind.AsConfig)
}

func LoadFromYaml(dst any, content []byte) error {
	return bind.BindYamlBytes(dst, content, bind.AsConfig)
}
