package config

import (
	"github.com/BurntSushi/toml"
	log "github.com/laohanlinux/utils/gokitlog"
)

type Config struct {
	Riot RiotConfig `toml:"riot"`
	Srv  Server     `toml:"server"`
	Log  LogConfig  `toml:"log"`
}

type Server struct {
	RPCAddr string `toml:"rpc-addr"`
	Token   string `toml:"token"`
}

type RiotConfig struct {
	RpcAddr  []string `toml:"rpc-addrs"`
	PoolSize int      `toml:"pool-size"`
}

type LogConfig struct {
	Dir   string `toml:"dir"`
	Name  string `toml:"name"`
	Level string `toml:"level"`
}

var Conf Config

func InitConfig(configFile string) {
	_, err := toml.DecodeFile(configFile, &Conf)
	if err != nil {
		panic(err)
	}
	opt := log.LogOption{
		SegmentationThreshold: 60,
		LogDir:                Conf.Log.Dir,
		LogName:               Conf.Log.Name,
		LogLevel:              Conf.Log.Level,
	}

	log.SetGlobalLog(opt)
}
