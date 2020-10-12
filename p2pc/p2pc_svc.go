package main

import (
	"sync"

	"github.com/isgo-golgo13/p2p_tunnel_svc/p2pc/client"
	"github.com/isgo-golgo13/p2p_tunnel_svc/p2pc/client_config"
)

const (
	GOROUTINES = 2
)

func main() {

	var wg sync.WaitGroup

	done := make(chan interface{})
	defer close(done)

	wg.Add(GOROUTINES)
	/** P2P client receiver goroutine from TUN */
	cli := client.NewUDPClient(client_config.CliConfig.ServerProto, client_config.CliConfig.Server)
	outStream, err := cli.CliRecv(&wg, done)
	if err != nil {
		cli.LogAs().Fatal("Error cannot reference or read to channel stream")
	}
	/** P2P client receiver/sender goroutine to UDP */
	cli.CliRecvSend(&wg, done, outStream)
	wg.Wait()
}
