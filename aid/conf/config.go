package conf

import (
	"fmt"
	"github.com/qinchende/gofast/aid/lang"
	"github.com/qinchende/gofast/store/bind"
	"log"
	"os"
	"path"
)

// 系统目前支持两种格式的配置文件：
// 1. JSON
// 2. Yaml
var loaders = map[string]func(any, []byte) error{
	".json": LoadConfigFromJsonBytes,
	".yaml": LoadConfigFromYamlBytes,
	".yml":  LoadConfigFromYamlBytes,
}

// 必须加载配置，否则应用无法启动，直接退出
func MustLoad(path string, dst any) {
	if err := LoadConfig(path, dst); err != nil {
		log.Fatalf("error: config file %s, %s", path, err.Error())
	}
}

func LoadConfig(file string, dst any) error {
	if content, err := os.ReadFile(file); err != nil {
		return err
	} else if loader, ok := loaders[path.Ext(file)]; ok {
		return loader(dst, lang.STB(os.ExpandEnv(string(content))))
	} else {
		return fmt.Errorf("unrecoginized file type: %s", file)
	}
}

func LoadConfigFromJsonBytes(dst any, content []byte) error {
	return bind.BindJsonBytes(dst, content, bind.AsConfig)
}

func LoadConfigFromYamlBytes(dst any, content []byte) error {
	return bind.BindYamlBytes(dst, content, bind.AsConfig)
}
