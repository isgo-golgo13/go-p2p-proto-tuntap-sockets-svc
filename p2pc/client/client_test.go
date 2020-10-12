package client

import (
	"context"
	"net"
	"network-poc/p2p_tunnelsvc/tun"
	"reflect"
	"sync"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/google/gopacket/layers"
)

func TestUDPClient_CliRecv(t *testing.T) {
	type fields struct {
		ClientCxt     context.Context
		Protocol      string
		Server        string
		ServerUDPAddr *net.UDPAddr
		ServerCon     net.Conn
		SrcIP         net.IP
		DstIP         net.IP
		LayerIPv4     *layers.IPv4
		LayerUDP      *layers.UDP
		TunIf         *tun.IfcCxt
		Logger        *log.Logger
	}
	type args struct {
		wg   *sync.WaitGroup
		done <-chan interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    <-chan []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &UDPClient{
				ClientCxt:     tt.fields.ClientCxt,
				Protocol:      tt.fields.Protocol,
				Server:        tt.fields.Server,
				ServerUDPAddr: tt.fields.ServerUDPAddr,
				ServerCon:     tt.fields.ServerCon,
				SrcIP:         tt.fields.SrcIP,
				DstIP:         tt.fields.DstIP,
				LayerIPv4:     tt.fields.LayerIPv4,
				LayerUDP:      tt.fields.LayerUDP,
				TunIf:         tt.fields.TunIf,
				Logger:        tt.fields.Logger,
			}
			got, err := c.CliRecv(tt.args.wg, tt.args.done)
			if (err != nil) != tt.wantErr {
				t.Errorf("UDPClient.CliRecv() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UDPClient.CliRecv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUDPClient_CliRecvSend(t *testing.T) {
	type fields struct {
		ClientCxt     context.Context
		Protocol      string
		Server        string
		ServerUDPAddr *net.UDPAddr
		ServerCon     net.Conn
		SrcIP         net.IP
		DstIP         net.IP
		LayerIPv4     *layers.IPv4
		LayerUDP      *layers.UDP
		TunIf         *tun.IfcCxt
		Logger        *log.Logger
	}
	type args struct {
		wg       *sync.WaitGroup
		done     <-chan interface{}
		inStream <-chan []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &UDPClient{
				ClientCxt:     tt.fields.ClientCxt,
				Protocol:      tt.fields.Protocol,
				Server:        tt.fields.Server,
				ServerUDPAddr: tt.fields.ServerUDPAddr,
				ServerCon:     tt.fields.ServerCon,
				SrcIP:         tt.fields.SrcIP,
				DstIP:         tt.fields.DstIP,
				LayerIPv4:     tt.fields.LayerIPv4,
				LayerUDP:      tt.fields.LayerUDP,
				TunIf:         tt.fields.TunIf,
				Logger:        tt.fields.Logger,
			}
			if err := c.CliRecvSend(tt.args.wg, tt.args.done, tt.args.inStream); (err != nil) != tt.wantErr {
				t.Errorf("UDPClient.CliRecvSend() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUDPClient_Send(t *testing.T) {
	type fields struct {
		ClientCxt     context.Context
		Protocol      string
		Server        string
		ServerUDPAddr *net.UDPAddr
		ServerCon     net.Conn
		SrcIP         net.IP
		DstIP         net.IP
		LayerIPv4     *layers.IPv4
		LayerUDP      *layers.UDP
		TunIf         *tun.IfcCxt
		Logger        *log.Logger
	}
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &UDPClient{
				ClientCxt:     tt.fields.ClientCxt,
				Protocol:      tt.fields.Protocol,
				Server:        tt.fields.Server,
				ServerUDPAddr: tt.fields.ServerUDPAddr,
				ServerCon:     tt.fields.ServerCon,
				SrcIP:         tt.fields.SrcIP,
				DstIP:         tt.fields.DstIP,
				LayerIPv4:     tt.fields.LayerIPv4,
				LayerUDP:      tt.fields.LayerUDP,
				TunIf:         tt.fields.TunIf,
				Logger:        tt.fields.Logger,
			}
			got, err := c.Send(tt.args.bytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("UDPClient.Send() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UDPClient.Send() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUDPClient_Sendn(t *testing.T) {
	type fields struct {
		ClientCxt     context.Context
		Protocol      string
		Server        string
		ServerUDPAddr *net.UDPAddr
		ServerCon     net.Conn
		SrcIP         net.IP
		DstIP         net.IP
		LayerIPv4     *layers.IPv4
		LayerUDP      *layers.UDP
		TunIf         *tun.IfcCxt
		Logger        *log.Logger
	}
	type args struct {
		bytes []byte
		n     int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &UDPClient{
				ClientCxt:     tt.fields.ClientCxt,
				Protocol:      tt.fields.Protocol,
				Server:        tt.fields.Server,
				ServerUDPAddr: tt.fields.ServerUDPAddr,
				ServerCon:     tt.fields.ServerCon,
				SrcIP:         tt.fields.SrcIP,
				DstIP:         tt.fields.DstIP,
				LayerIPv4:     tt.fields.LayerIPv4,
				LayerUDP:      tt.fields.LayerUDP,
				TunIf:         tt.fields.TunIf,
				Logger:        tt.fields.Logger,
			}
			got, err := c.Sendn(tt.args.bytes, tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("UDPClient.Sendn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UDPClient.Sendn() = %v, want %v", got, tt.want)
			}
		})
	}
}
