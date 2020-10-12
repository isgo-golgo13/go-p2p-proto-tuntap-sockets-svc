package server

import (
	"context"
	"fmt"
	"net"
	"network-poc/p2p_tunnelsvc/p2ps/server_config"
	"network-poc/p2p_tunnelsvc/tun"
	"os"
	"sync"

	log "github.com/Sirupsen/logrus"
)

type ServerSvc interface {
	Close() error
	Recv() (int, []byte, *net.UDPAddr, error)
	SrvRecv(wg *sync.WaitGroup, done <-chan interface{}) <-chan []byte
	SrvRecvSend(wg *sync.WaitGroup, done <-chan interface{}, inStream <-chan []byte)
}

type UDPServer struct {
	ServerCtx context.Context
	ServerCon *net.UDPConn

	TunIf  *tun.IfcCtx
	Logger *log.Logger
}

func NewUDPServer() *UDPServer {
	server := &UDPServer{
		TunIf:  tun.NewIfcCtx(server_config.SrvConfig.ServerTUN),
		Logger: &log.Logger{Out: os.Stderr, Formatter: &log.TextFormatter{ForceColors: true, FullTimestamp: true}, Hooks: make(log.LevelHooks), Level: log.InfoLevel | log.ErrorLevel},
	}
	server.Logger.Infof("P2P UDP/TUN server at TUN: %s", server.TunIf.Ifc.Name())
	return server
}

func init() {

}

func (s *UDPServer) NewUDPAddr(protocol, addr string) (*net.UDPAddr, error) {
	udpAddr, err := net.ResolveUDPAddr(protocol, addr)
	if err != nil {
		s.Logger.Errorf("Error ResolveUDPAddr: %v", err)
		return nil, err
	}
	return udpAddr, nil
}

/** Recv **/
func (s *UDPServer) Recv(bytes []byte) (int, *net.UDPAddr, error) {
	nr, clientCon, err := s.ServerCon.ReadFromUDP(bytes)
	if err != nil {
		s.Logger.Errorf("Server receive error:  %v", err)
		return 0, nil, err
	}
	return nr, clientCon, nil
}

/** SrvRecv */
func (s *UDPServer) SrvRecv(wg *sync.WaitGroup, done <-chan interface{}) <-chan []byte {
	outStream := make(chan []byte)
	go func() {
		defer wg.Done()
		defer close(outStream)

		var err error
		var udpAddr *net.UDPAddr
		udpAddr, err = s.NewUDPAddr(server_config.SrvConfig.ServerProto, fmt.Sprintf("%s:%d", server_config.SrvConfig.Server, server_config.SrvConfig.ServerPort))
		if err != nil {
			s.Logger.Errorf("Error NewUDPAddr: %v", err)
			return
		}

		s.ServerCon, err = net.ListenUDP(server_config.SrvConfig.ServerProto, udpAddr)
		if err != nil {
			s.Logger.Errorf("Error ListenUDP: %v\n", err)
			return
		}
		defer s.Close()

		for {
			select {
			case <-done:
				return
			default:
				buf := make([]byte, server_config.SrvConfig.Mtu)

				nr, clientCon, err := s.Recv(buf)
				if err != nil {
					s.Logger.Errorf("Error ReadFromUDP:  %v", err)
					continue
				}
				s.Logger.Infof("\nReceived from %d bytes from client: %v %x \n", nr, clientCon, buf[:nr])
				outStream <- buf[:nr]
			}
		}
	}()
	return outStream
}

/** SrvRecvSend */
func (s *UDPServer) SrvRecvSend(wg *sync.WaitGroup, done <-chan interface{}, inStream <-chan []byte) {
	go func() {
		defer wg.Done()

		for {
			select {
			case <-done:
				return
			default:
				buf := <-inStream
				nw, err := s.TunIf.Ifc.Write(buf)
				if err != nil {
					log.Fatal(err)
				}
				s.Logger.Infof("UDP/TUN Server sent %d bytes to TUN: %s", nw, s.TunIf.Ifc.Name())
			}
		}
	}()
	return
}

/** Close client connection */
func (s *UDPServer) Close() error {
	err := s.ServerCon.Close()
	if err != nil {
		s.Logger.Errorf("Error Close : %v", err)
		return err
	}
	return nil
}

/** LogAs is handle logger for main */
func (s *UDPServer) LogAs() *log.Logger {
	return s.Logger
}
