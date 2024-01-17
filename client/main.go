package main

import (
	"context"
	"pi/op_cli/cli"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)


func main() {
	logrus.Println("Cli of PiDns")

	badr := ":50051"

	conn,err := grpc.Dial(*&badr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		logrus.Fatalf("did not connect %v", err)
	}

	defer conn.Close()


	c := cli.NewHostManagerClient(conn)


	ctx,cancel := context.WithTimeout(context.Background(),time.Second * 2)

	defer cancel()

	r,err := c.AddHosts(ctx, &cli.Host{
		Domain:"abc.pi.g",
		Ips: []string{
			"192.168.224.88",
		},
	})


	if err != nil {
		logrus.Fatalf("could not greet:%v", err)
	}

	logrus.Infof("Message %v",r.GetCode())
}
