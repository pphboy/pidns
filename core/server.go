package core

import (
	"context"
	"math/rand"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/miekg/dns"
	"github.com/sirupsen/logrus"
)

type DnsServer struct {

	// 绑定的地址和端口
	BindAddress string `toml:"bindAddress"`

	rejectType []uint16

	Hosts Hosts

	ctx context.Context
}

func (ds *DnsServer) NewServer() {

	ds.BindAddress = ":53"
	ds.rejectType = []uint16{255}
	ds.Hosts.Map = map[string][]string{
		"pi.g": []string{
			"192.168.224.88",
		},
		"node1.pi.g": []string{
			"192.168.224.88",
		},
	}
	go func() {
		i := 0

		for {
			select {
			case <-time.After(10 * time.Second):
				i++
				node := "node" + strconv.Itoa(i) + ".pi.g"
				ds.Hosts.Map[node] = []string{"192.168.224.88"}
				logrus.Printf("%v", ds.Hosts.Map)
			}
		}
	}()
	ds.ctx = context.Background()
}

func (ds *DnsServer) Run() {

	mux := dns.NewServeMux()

	// 这一步就是将 自定义的 dns服务器，设置到 dns 包中去 作为处理器

	mux.Handle(".", ds)

	wg := new(sync.WaitGroup)

	wg.Add(2)

	for _, p := range [2]string{"tcp", "udp"} {
		go func(p string) {
			srv := &dns.Server{Addr: ds.BindAddress, Net: p, Handler: mux}

			go func() {
				<-ds.ctx.Done()
				logrus.Warnf("close server %s", p)
				srv.ShutdownContext(ds.ctx)
			}()

			err := srv.ListenAndServe()
			if err != nil {
				logrus.Fatalf("failed to listening %s ,%v", p, err)
				os.Exit(1)
			}

			wg.Done()
		}(p)
	}

	wg.Wait()
}

func (ds *DnsServer) Stop() {
	ds.ctx.Done()
}

func (ds *DnsServer) ServeDNS(w dns.ResponseWriter, q *dns.Msg) {
	ip, _, _ := net.SplitHostPort(w.RemoteAddr().String())

	ok := ds.validateRejectType(q)

	if ok {
		// 处理
		as := ds.Exchange(q)
		if as != nil {

			err := w.WriteMsg(as)

			if err != nil {
				logrus.Warnf("failed to write message, %v", err)
			}

			logrus.Infof("succeed wirte msg, %+v", as)
			return
		}
	} else {
		logrus.Warnf("Reject %s: %s", ip, q.Question[0].String())

		dns.HandleFailed(w, q)
	}

}

func (ds *DnsServer) Exchange(q *dns.Msg) *dns.Msg {
	qn := q.Question[0].Name
	name := qn[:len(qn)-1]
	ipv4s := ds.Hosts.FindHosts(name)

	var s dns.Msg

	if q.Question[0].Qtype == dns.TypeA && len(ipv4s) > 0 {
		// var rrl []dns.RR
		for _, ip := range ipv4s {
			rr, _ := dns.NewRR(qn + " IN A " + ip.String())
			// rrl = append(rrl, a)
			s.Answer = append(s.Answer, rr)
		}

		// 设置时间戳
		rand.Seed(time.Now().UnixNano())
		for i := range s.Answer {
			j := rand.Intn(i + 1)
			s.Answer[i], s.Answer[j] = s.Answer[j], s.Answer[i]
		}

		// 设置返回 数据包的 信息
		s.SetReply(q)

		s.RecursionAvailable = true

		SetMinimumTTL(&s, 0)

	} else {

		rand.Seed(time.Now().UnixNano())
		for i := range s.Answer {
			j := rand.Intn(i + 1)
			s.Answer[i], s.Answer[j] = s.Answer[j], s.Answer[i]
		}

		// 设置返回 数据包的 信息
		s.SetReply(q)

		s.RecursionAvailable = true

		SetMinimumTTL(&s, 0)

	}

	// 如果不返回，则说明 dns服务器挂了
	// 所以无论如何都需要返回信息
	return &s
}

// 检查是否存在 拒绝访问的类型
func (ds *DnsServer) validateRejectType(q *dns.Msg) bool {

	for _, qt := range ds.rejectType {
		if q.Question[0].Qtype == qt {
			return false
		}
	}

	return true
}

func SetMinimumTTL(msg *dns.Msg, minimumTTL uint32) {
	if minimumTTL == 0 {
		return
	}
	for _, a := range msg.Answer {
		if a.Header().Ttl < minimumTTL {
			a.Header().Ttl = minimumTTL
		}
	}
}
