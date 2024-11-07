package conf

import (
	"fmt"
	"github.com/qinchende/gofast/core/lang"
	"github.com/qinchende/gofast/store/bind"
	"log"
	"os"
	"path"
)

// 系统目前支持下面几种格式的配置文件：
// 1. Json
// 2. Yaml
// 3. Toml
// 4. Ini
var loaders = map[string]func(any, []byte) error{
	".json": LoadFromJson,
	".toml": LoadFromToml,
	".yaml": LoadFromYaml,
	".yml":  LoadFromYaml,
	".int":  LoadFromIni,
}

// 必须加载配置，否则应用无法启动，直接退出
func MustLoad(dst any, file string) {
	if err := LoadFile(dst, file); err != nil {
		log.Fatalf("error: config file %s, %s", file, err.Error())
	}
}

func LoadFile(dst any, file string) error {
	if content, err := os.ReadFile(file); err != nil {
		return err
	} else if loader, ok := loaders[path.Ext(file)]; ok {
		return loader(dst, lang.S2B(os.ExpandEnv(string(content))))
	} else {
		return fmt.Errorf("unsupport file type: %s", file)
	}
}

func LoadFromJson(dst any, content []byte) error {
	return bind.BindJsonBytes(dst, content, bind.AsConfig)
}

func LoadFromYaml(dst any, content []byte) error {
	return bind.BindYamlBytes(dst, content, bind.AsConfig)
}

func LoadFromToml(dst any, content []byte) error {
	return nil
}

func LoadFromIni(dst any, content []byte) error {
	return nil
}
