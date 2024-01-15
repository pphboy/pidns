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

func (h *Hosts)DelHost(name string) {
	delete(h.Map,name)
}


func (h *Hosts) AddHosts(name string,ips []string ) {
	if h.Map == nil {
		h.Map = make(map[string][]string)
	}
	h.Map[name] = ips
}
