package server

import (
	"context"

	"net"
	"pi_dns/core"

	"github.com/golang/protobuf/ptypes"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func InitMngServ() {
	logrus.Infoln("Start GRPC Server of PiDns")
}

type ManageServer interface {
	RunServ()
	// new host manager Server
	NewMngServ(lisn string, network string,ds *core.DnsServer)
}

type HostServer struct {
	// listen server address
	lisndAddr string
	// tcp or udp
	network string
	// dnsServer instance
	dnsServer *core.DnsServer

	UnimplementedHostManagerServer
}

func (h *HostServer) NewMngServ(lisn string, network string,ds *core.DnsServer) {
	h.lisndAddr = lisn
	h.network = network
	h.dnsServer = ds
}

func (h *HostServer) RunServ() {
	lis, err := net.Listen(h.network, h.lisndAddr)

	if err != nil {
		logrus.Fatalf("failed run host server of manager, err:%v", err)
	}

	s := grpc.NewServer()

	RegisterHostManagerServer(s, h)

	logrus.Infof("host manager server running at %v", h.lisndAddr)

	if err := s.Serve(lis); err != nil {
		logrus.Fatalf("failed to serve: %v", h.lisndAddr)
	}
}

func (h *HostServer) AddHosts(c context.Context, host *Host) (*Result, error) {

	h.dnsServer.Hosts.AddHosts(host.Domain, host.Ips)

	return &Result{
		Code: 1,
	}, nil
}

func (h *HostServer) DelHosts(c context.Context, host *Host) (*Result, error) {

	h.dnsServer.Hosts.DelHost(host.Domain)
	return &Result{
		Code: 1,
	}, nil
}

func (h *HostServer) GetHosts(ctx context.Context, host *Host) (*Result, error) {

	ips := h.dnsServer.Hosts.GetHosts(host.Domain)

	copy(host.Ips, ips)

	data, err := ptypes.MarshalAny(host)

	if err != nil {
		logrus.Warnf("host dont exists, domain %v", host.Domain)
		return &Result{
			Code: 0,
			Data: nil,
		}, nil
	}

	return &Result{
		Code: 1,
		Data: data,
	}, nil
}
