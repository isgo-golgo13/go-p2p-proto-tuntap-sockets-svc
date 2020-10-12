package main

import (
	"sync"

	"github.com/isgo-golgo13/p2p_tunnel_svc/p2ps/server"
)

const (
	GOROUTINES = 2
)

func main() {

	done := make(chan interface{})
	defer close(done)

	var wg sync.WaitGroup

	wg.Add(GOROUTINES)
	srv := server.NewUDPServer()
	outStream := srv.SrvRecv(&wg, done)
	srv.SrvRecvSend(&wg, done, outStream)
	wg.Wait()

}
