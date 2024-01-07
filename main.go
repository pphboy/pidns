package main

import (
	"pi_dns/core"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Println("PiDns Start")

	ds := core.DnsServer{}

	ds.NewServer()

	ds.Run()

}
