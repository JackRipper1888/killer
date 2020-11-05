package confkit

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
)

func InitCfgFilePath(path string) string {
	filePath := flag.String("f", path, "config file path")
	flag.Parse()
	return *filePath
}
func GetConfig(filePath string, cfg interface{}) error {
	if filePath == "" {
		panic("配置文件路径为空")
	}
	if _, err := toml.DecodeFile(filePath, cfg); err != nil {
		panic(fmt.Sprintf("配置文件解析错误，err:#{%s}", err))
		return err
	}
	return nil
}
