package client

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/isgo-golgo13/p2p_tunnel_svc/p2pc/client_config"
	"github.com/isgo-golgo13/p2p_tunnel_svc/tun"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	log "github.com/sirupsen/logrus"
)

type ClientSvc interface {
	AddAndCheckSrcDstIP() error
	AddIPPacketLayer()
	AddUDPacketLayer()
	AddPacketLayerPayload(bytes []byte) error
	Dial() error
	Send(bytes []byte) (int, error)
	Sendn(bytes []byte, n int) (int, error)
	Close() error
	CliRecv(wg *sync.WaitGroup, done <-chan interface{}) (<-chan []byte, error)
	CliRecvSend(wg *sync.WaitGroup, done <-chan interface{}, inStream <-chan []byte) error
}

/** UDPClient struct handle for UDP/Tun client activity */
type UDPClient struct {
	ClientCxt     context.Context
	Protocol      string
	Server        string
	ServerUDPAddr *net.UDPAddr
	ServerCon     net.Conn
	SrcIP         net.IP
	DstIP         net.IP
	LayerIPv4     *layers.IPv4 /** Google gopacket: IP layer struct */
	LayerUDP      *layers.UDP  /** Google gopacket: UDP layer struct */
	TunIf         *tun.IfcCtx

	Logger *log.Logger
}

func NewUDPClient(protocol, server string) *UDPClient {
	client := &UDPClient{
		Protocol: protocol,
		Server:   server,
		LayerUDP: &layers.UDP{
			SrcPort: layers.UDPPort(client_config.CliConfig.ClientSrcPort),
			DstPort: layers.UDPPort(client_config.CliConfig.ServerDstPort),
		},
		TunIf:  tun.NewIfcCtx(client_config.CliConfig.ClientTUN),
		Logger: &log.Logger{Out: os.Stderr, Formatter: &log.TextFormatter{ForceColors: true, FullTimestamp: true}, Hooks: make(log.LevelHooks), Level: log.InfoLevel | log.ErrorLevel},
	}
	client.Logger.Infof("UDP client at %s receivng TUN at: %s", client_config.CliConfig.Client+":"+fmt.Sprintf("%d", client_config.CliConfig.ClientSrcPort), client.TunIf.Ifc.Name())
	return client
}

type PacketSerialConfig struct {
	SerialBuf  gopacket.SerializeBuffer
	SerialOpts gopacket.SerializeOptions
}

/** NewPacketSerialCondfig default constructs Google gopacket SerializeBuffer and SerializeOpts
  required for the gopacket SerializeLayers func **/
func NewPacketSerialConfig(doFixLen, doChksum bool) *PacketSerialConfig {
	return &PacketSerialConfig{
		SerialBuf: gopacket.NewSerializeBuffer(),
		SerialOpts: gopacket.SerializeOptions{
			FixLengths:       doFixLen,
			ComputeChecksums: doChksum,
		},
	}
}

/** Any pre-connection or connection objects go in here (not including UDP connections) */
func init() {

}

/* Predial inits any UDP source and destination address and port and address resolution */
func (c *UDPClient) AddAndCheckSrcDstIP() error {
	c.SrcIP = net.ParseIP(client_config.CliConfig.Client)
	if c.SrcIP == nil {
		errf := fmt.Sprintf("Invalid source IP: %v\n", c.SrcIP)
		c.Logger.Errorf(errf)
		return errors.New(errf)
	}
	c.DstIP = net.ParseIP(client_config.CliConfig.Server)
	if c.DstIP == nil {
		errf := fmt.Sprintf("Invalid destination IP: %v\n", c.DstIP)
		c.Logger.Errorf(errf)
		return errors.New(errf)
	}
	return nil
}

func (c *UDPClient) AddIPPacketLayer() {
	c.LayerIPv4 = &layers.IPv4{
		SrcIP:    c.SrcIP,
		DstIP:    c.DstIP,
		Protocol: layers.IPProtocolUDP,
	}
}

func (c *UDPClient) AddUDPacketLayer() {
	c.LayerUDP = &layers.UDP{
		SrcPort: layers.UDPPort(c.LayerUDP.SrcPort),
		DstPort: layers.UDPPort(c.LayerUDP.DstPort),
	}
}

/** AddPacketLayerPayload is the entry function receiving the channel tun packet to relay to Send */
func (c *UDPClient) AddPacketLayerPayload(bytes []byte) (*PacketSerialConfig, error) {
	/** originally called checksum func here */
	sc := NewPacketSerialConfig(true, false)
	err := gopacket.SerializeLayers(sc.SerialBuf, sc.SerialOpts, c.LayerUDP, gopacket.Payload(bytes))
	if err != nil {
		c.Logger.Error("Error in serialize of packet layers")
		return nil, err
	}
	return sc, nil
}

/** Dial dials a UDP server */
func (c *UDPClient) Dial() error {
	var err error
	c.ServerCon, err = net.Dial(c.Protocol, c.Server)
	if err != nil {
		c.Logger.Errorf("Dial to server error: %v", err)
		return err
	}
	return nil
}

/** Sendn (external send API) sends bytes on UDP socket connection */
func (c *UDPClient) Send(bytes []byte) (int, error) {
	return c.sendn(bytes, 0)
}

/** Sendn (external send API) sends n offset bytes on UDP socket connection */
func (c *UDPClient) Sendn(bytes []byte, n int) (int, error) {
	return c.sendn(bytes, n)
}

/** sendn (internal/package scope send API) sends n bytes on UDP socket connection */
func (c *UDPClient) sendn(bytes []byte, n int) (int, error) {
	var nw int
	var err error
	if n == 0 {
		nw, err = c.ServerCon.Write(bytes)
		if err != nil {
			c.Logger.Errorf("Could not write to server: %v", err)
			return 0, err
		}
	} else {
		nw, err = c.ServerCon.Write(bytes[:n])
		if err != nil {
			c.Logger.Errorf("Could not write to server: %v", err)
			return 0, err
		}
	}
	return nw, nil
}

/** Close client connection */
func (c *UDPClient) Close() error {
	err := c.ServerCon.Close()
	if err != nil {
		c.Logger.Errorf("Error client peer %s connection close", fmt.Sprintf("%s:%d", client_config.CliConfig.Client, client_config.CliConfig.ClientSrcPort))
		return err
	}
	return nil
}

/** CliRecv is the TUN interface receiver and relay to UDP client */
func (c *UDPClient) CliRecv(wg *sync.WaitGroup, done <-chan interface{}) (<-chan []byte, error) {
	outStream := make(chan []byte)
	go func() error {
		defer wg.Done()
		defer close(outStream)

		for {
			select {
			case <-done:
				return nil
			default:
				buf := make([]byte, client_config.CliConfig.Mtu)
				nr, err := c.TunIf.Ifc.Read(buf)
				if err != nil {
					c.Logger.Errorf("TUN interface %s fail to receive", c.TunIf.Ifc.Name())
					return err
				}
				outStream <- buf[:nr]
			}
		}
	}()
	return outStream, nil
}

/** UDP client relay TUN/UDP (ip4:17) packet to UDP Server */
func (c *UDPClient) CliRecvSend(wg *sync.WaitGroup, done <-chan interface{}, inStream <-chan []byte) error {
	go func() error {
		defer wg.Done()

		err := c.Dial()
		if err != nil {
			c.Logger.Errorf("Dial error to server %s", fmt.Sprintf("%s:%d", client_config.CliConfig.Server, client_config.CliConfig.ServerDstPort))
			return err
		}
		defer c.Close()

		for {
			select {
			case <-done:
				return nil
			default:
				err := c.AddAndCheckSrcDstIP()
				if err != nil {
					c.Logger.Errorf("Error adding source IP address [%s] and destination IP address [%s]", c.SrcIP.String(), c.DstIP.String())
					return err
				}
				c.AddUDPacketLayer()
				packetPayload, err := c.AddPacketLayerPayload(<-inStream)
				if err != nil {
					c.Logger.Error("Error in serialization of the packet payload")
					return err
				}
				nw, err := c.Send(packetPayload.SerialBuf.Bytes())
				if err != nil {
					c.Logger.Errorf("Error sending to server: %v", err)
					return err
				}
				c.Logger.Infof("Client sent %d bytes to UDP/TUN P2P server %s: ", nw, fmt.Sprintf("%s:%d", c.DstIP, c.LayerUDP.DstPort))
			}
		}
	}()
	return nil
}

/** LogAs is handle logger for main */
func (c *UDPClient) LogAs() *log.Logger {
	return c.Logger
}
