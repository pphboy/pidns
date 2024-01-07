package core

import "net"

type Hosts struct {
	Map map[string][]string
}

func (h *Hosts) FindHosts(name string) (nips []net.IP) {

	ipv4s, ok := h.Map[name]
	if ok {
		for _, ip := range ipv4s {
			nips = append(nips, net.ParseIP(ip))
		}
	}

	return
}

func MoveHosts(name string) {

}

func AddHosts() {

}
