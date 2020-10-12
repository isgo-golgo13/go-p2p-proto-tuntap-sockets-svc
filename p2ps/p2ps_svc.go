package main

import (
	"network-poc/p2p_tunnelsvc/p2ps/server"
	"sync"
)

const ( GOROUTINES = 2 )


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
