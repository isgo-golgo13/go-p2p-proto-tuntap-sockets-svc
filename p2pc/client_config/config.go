package client_config

import (
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/tkanos/gonfig"
)

type ClientConfiguration struct {
	Client        string
	Server        string
	ClientSrcPort int
	ServerDstPort int
	ServerProto   string
	Mtu           int
	ClientTUN     string
}

var CliConfig *ClientConfiguration

/** NewClientConfiguration factory function for centralized config on the client */
func NewClientConfiguration() (*ClientConfiguration, error) {
	config := &ClientConfiguration{}

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
	confDir := filepath.Join(dir, "/p2pc/client_config/config.dev.json")
	return confDir, nil
}

func init() {
	var err error
	CliConfig, err = NewClientConfiguration()
	if err != nil {
		log.Errorf("Error: NewClientConfiguration %v", err)
	}
}
