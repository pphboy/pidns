package main

import (
	"pi_dns/core"
	"pi_dns/server"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Println("PiDns Start")

	ds := core.NewServer(&core.DnsServer{
		BindAddress: ":53",
		RejectType: []uint16{
			255,
		},
	})

	ds.Hosts.AddHosts("pi.g", []string{
		"192.168.224.88",
	})

	ds.Hosts.AddHosts("node1.pi.g", []string{
		"192.168.224.88",
	})

	var ms server.ManageServer

	ms = &server.HostServer{}
	
	ms.NewMngServ(":50051", "tcp", ds)

	ms.RunServ()
	
	ds.Run()

}
