package server_config

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/tkanos/gonfig"
)

type ServerConfiguration struct {
	Server      string
	ServerPort  int
	ServerProto string
	Mtu         int
	ServerTUN   string
}

var SrvConfig *ServerConfiguration

/** NewServerConfiguration factory function for centralized config on the server */
func NewServerConfiguration() (*ServerConfiguration, error) {
	config := &ServerConfiguration{}

	confDir, err := getConfDir()
	if err != nil {
		log.Errorf("%v", err)
		return nil, err
	}
	err = gonfig.GetConf(confDir, config)
	if err != nil {
		log.Errorf("%v", err)
		return nil, err
	}
	return config, nil
}

func getConfDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		log.Errorf("%v", err)
		return "", err
	}
	confDir := filepath.Join(dir, "/p2ps/server_config/config.dev.json")
	return confDir, nil
}

func init() {
	var err error
	SrvConfig, err = NewServerConfiguration()
	if err != nil {
		log.Errorf("Error: NewServerConfiguration %v", err)
	}
}
