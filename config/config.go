package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type ServerConfig struct {
	Addr string `toml:"addr"`
	Port string `toml:"port"`
}

type RpcConfig struct {
	Addr string `toml:"addr"`
	Port string `toml:"port"`
}

type RaftConfig struct {
	Addr string `toml:"addr"`
	Port string `toml:"port"`
}

type LogConfig struct {
	LogDir  string `toml:"log_dir"`
	LogName string `toml:"log_name"`
}

type Configure struct {
	SC    ServerConfig `toml:"server"`
	RpcC  RpcConfig    `toml:"rpc"`
	RaftC RaftConfig   `toml:"raft"`
	LogC  LogConfig    `toml:"log"`
}

var c *Configure

func NewConfig(data string) (*Configure, error) {
	c = new(Configure)
	_, err := toml.Decode(data, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func GetConfigure() *Configure {
	return c
}

func (c *Configure) Info() {
	fmt.Printf("raft: %v\n", c.RaftC)
	fmt.Printf("rpc: %v\n", c.RpcC)
	fmt.Printf("server:%v\n", c.SC)
}
