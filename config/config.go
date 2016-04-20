package config

import (
	"fmt"
	"os"

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

func (rpc RpcConfig) AddrString() string {
	return fmt.Sprintf("%s:%s", rpc.Addr, rpc.Port)
}

type LeaderRpcConfig struct {
	Addr string
	Port string
}

func (lrc LeaderRpcConfig) AddrString() string {
	return fmt.Sprintf("%s:%s", lrc.Addr, lrc.Port)
}

type RaftConfig struct {
	Addr             string   `toml:"addr"`
	Port             string   `toml:"port"`
	Peers            []string `toml:"peers"`
	PeerStorage      string   `toml:"peer_storage"`
	SnapshotStorage  string   `toml:"snapshot_storage"`
	StoreBackendPath string   `toml:"storage_backend_path"`
	RaftLogPath      string   `toml:"raft_log_path"`
	ApplyLogPath     string   `toml:"apply_log_path"`
	EnableSingleNode bool     `toml:"enable_single_node"`
}

func (rc RaftConfig) AddrString() string {
	return fmt.Sprintf("%s:%s", rc.Addr, rc.Port)
}

type LogConfig struct {
	LogDir  string `toml:"log_dir"`
	LogName string `toml:"log_name"`
}

type Configure struct {
	SC         ServerConfig `toml:"server"`
	LeaderRpcC LeaderRpcConfig
	RpcC       RpcConfig  `toml:"rpc"`
	RaftC      RaftConfig `toml:"raft"`
	LogC       LogConfig  `toml:"log"`
}

var c *Configure

func NewConfig(data string) (*Configure, error) {
	c = new(Configure)
	_, err := toml.Decode(data, c)
	if err != nil {
		return nil, err
	}
	c.RpcC.Addr = c.RaftC.Addr
	// Init Application Dir
	dirs := []string{
		c.LogC.LogDir,
		c.RaftC.PeerStorage,
		c.RaftC.SnapshotStorage,
		c.RaftC.StoreBackendPath,
		c.RaftC.ApplyLogPath,
		c.RaftC.RaftLogPath,
	}

	for _, dir := range dirs {
		if _, err := os.Stat(dir); err != nil && os.IsNotExist(err) {
			if err := os.MkdirAll(dir, os.ModePerm); err != nil {
				return nil, err
			}
		}
	}
	return c, nil
}

func GetConfigure() *Configure {
	return c
}

func (c *Configure) Info() {
	fmt.Printf("raft: %v\n", c.RaftC)
	fmt.Printf("rpc: %v\n", c.RpcC)
	fmt.Printf("leader rpc: %v\n", c.LeaderRpcC)
	fmt.Printf("server:%v\n", c.SC)
}
