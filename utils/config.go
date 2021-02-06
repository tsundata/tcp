package utils

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Name    string `json:"name"`
	Host    string `json:"host"`
	TCPPort int    `json:"port"`
	Version string `json:"version"`

	MaxPacketSize    uint32 `json:"max_packet_size"`
	MaxConn          int    `json:"max_conn"`
	WorkerPoolSize   uint32 `json:"worker_pool_size"`
	MaxWorkerTaskLen uint32 `json:"max_worker_task_len"`
}

var Setting *Config

func (c *Config) Reload() {
	data, err := ioutil.ReadFile("conf/server.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &Setting)
	if err != nil {
		panic(err)
	}
}

func init() {
	Setting = &Config{
		Name:             "server",
		Host:             "0.0.0.0",
		TCPPort:          5678,
		Version:          "v0.1",
		MaxPacketSize:    4096,
		MaxConn:          12000,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
	}
	Setting.Reload()
}
