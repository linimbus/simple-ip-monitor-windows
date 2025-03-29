package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/astaxie/beego/logs"
)

type Config struct {
	ConnectivityURL string
	FilterInterface string
	OutputDirectory string
	RestfulHeader   string
	RestfulMethod   string
	RestfulURL      string
	Interval        int
	AutoStartup     bool
}

var configCache = Config{
	ConnectivityURL: "https://test.ipw.cn/",
	FilterInterface: "",
	OutputDirectory: "",
	RestfulHeader:   "",
	RestfulMethod:   "POST",
	RestfulURL:      "",
	Interval:        60,
	AutoStartup:     false,
}

var configFilePath string
var configLock sync.Mutex

func configSyncToFile() error {
	configLock.Lock()
	defer configLock.Unlock()

	value, err := json.MarshalIndent(configCache, "\t", " ")
	if err != nil {
		logs.Error("json marshal config fail, %s", err.Error())
		return err
	}
	return os.WriteFile(configFilePath, value, 0664)
}

func ConfigGet() *Config {
	return &configCache
}

func ConnectivityURLSave(url string) error {
	configCache.ConnectivityURL = url
	return configSyncToFile()
}

func FilterInterfaceSave(filter string) error {
	configCache.FilterInterface = filter
	return configSyncToFile()
}

func OutputDirectorySave(dir string) error {
	configCache.OutputDirectory = dir
	return configSyncToFile()
}

func RestfulHeaderSave(key, value string) error {
	configCache.RestfulHeader = fmt.Sprintf("%s:%s", key, value)
	return configSyncToFile()
}

func RestfulMethodSave(method string) error {
	configCache.RestfulMethod = method
	return configSyncToFile()
}

func RestfulURLSave(url string) error {
	configCache.RestfulURL = url
	return configSyncToFile()
}

func IntervalSave(value int) error {
	configCache.Interval = value
	return configSyncToFile()
}

func AutoStartupSave(value bool) error {

	if value {
		execPath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("call executable failed, %s", err.Error())
		}
		err = RegistryStartupSet(APPLICATION_NAME, execPath)
		if err != nil {
			return err
		}
	} else {
		err := RegistryStartupDel(APPLICATION_NAME)
		if err != nil {
			return err
		}
	}

	configCache.AutoStartup = value
	return configSyncToFile()
}

func ConfigInit() {
	var err error
	var value []byte

	configFilePath = filepath.Join(ConfigDirGet(), "config.json")

	defer func() {
		if err != nil {
			err = configSyncToFile()
			if err != nil {
				logs.Error("config sync to file fail, %s", err.Error())
			}
		}
	}()

	_, err = os.Stat(configFilePath)
	if err != nil {
		logs.Info("config file not exist, create a new one")
	}

	value, err = os.ReadFile(configFilePath)
	if err != nil {
		logs.Error("read config file from app data dir fail, %s", err.Error())
		return
	}

	err = json.Unmarshal(value, &configCache)
	if err != nil {
		logs.Error("json unmarshal config fail, %s", err.Error())
		return
	}
}
