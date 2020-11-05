package confkit

import (
	"fmt"
	"os"
	"path/filepath"

	"testing"
)

var (
	ConfData = GetConfInfo()
	Path, _  = filepath.Abs(filepath.Dir(os.Args[0]))
)

func GetConfInfo() AllConfig {
	confDir := Path + "/confkit/tracker.conf"
	confFilePath := InitCfgFilePath(confDir)
	config := AllConfig{}
	err := GetConfig(confFilePath, &config)
	if err != nil {
		panic(err)
	}
	return config
}

type AllConfig struct {
	Tracker  Trackerconfig `toml:"tracker"`
	GeoDb    GeoDb         `toml:"geo"`
	RedisDb  RedisDB       `toml:"redis_db"`
	RedisAss RedisAss      `toml:"redis_ass"`
}
type Trackerconfig struct {
	Sid        string `toml:"id"`
	Ips        string `toml:"ips"`
	Controller string `toml:"controller"`
	Passwd     string `toml:"passwd"`
	Loglevel   int    `toml:"loglevel"`
	LogPath    string `toml:"log_path"`
}
type GeoDb struct {
	GeodbPath string `toml:"geodb_path"`
}
type RedisDB struct {
	Addr      string `toml:"addr"`
	Pwd       string `toml:"pwd"`
	MaxConn   int    `toml:"max_conn"`
	IsExpire  bool   `toml:"is_expire"`
	ExpireSec int64  `toml:"expire_sec"`
}
type RedisAss struct {
	IsRel        bool  `toml:"is_rel"`
	RelExpireSec int64 `toml:"rel_expire_sec"`
}

func TestGetConfig(t *testing.T) {
	fmt.Println(ConfData.GeoDb.GeodbPath)
}
