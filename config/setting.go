package config

import (
	"os"
	"path/filepath"
	"strconv"
)

var CF Configer //配置

const (
	CONFIG					string = "./config.ini"
	LOG_DIR					string = "./logs"
)

// 配置文件涉及的默认配置。
const (
	serverport				int    = 9000
)

func NewConfiger() Configer {
	os.MkdirAll(filepath.Clean(LOG_DIR), 0777)

	iniconf, err := NewConfig("ini", CONFIG)
	if err != nil {
		file, err := os.Create(CONFIG)
		file.Close()
		iniconf, err = NewConfig("ini", CONFIG)
		if err != nil {
			panic(err)
		}
		defaultConfig(iniconf)
		iniconf.SaveConfigFile(CONFIG)
	} else {
		trySet(iniconf)
	}

	return iniconf
}

func defaultConfig(iniconf Configer) {
	iniconf.Set("serverport", strconv.Itoa(serverport))
}

func trySet(iniconf Configer) {

	if v, e := iniconf.Int64("serverport"); v <= 0 || e != nil {
		iniconf.Set("serverport", strconv.Itoa(serverport))
	}
	
	iniconf.SaveConfigFile(CONFIG)
}

func init() {
	CF = NewConfiger()
}
