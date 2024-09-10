package controllers

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
)

var ErrNoProxies = errors.New("no proxies available")
var ErrUnreachable = errors.New("unreachable")

type ProxyController struct {
	Proxies map[int64]*Proxy
	Ppx     int64
}

var proxyTypes = map[int8]string{1: "socks4", 2: "socks5", 3: "http"}

type Proxy struct {
	ID        int64
	Type      int8
	IP        string
	Port      int
	Errors    int64
	UserAgent string
}

func NewProxyController() *ProxyController {
	return &ProxyController{Proxies: make(map[int64]*Proxy)}
}

func (pc *ProxyController) AddProxy(px ...*Proxy) {
	for _, el := range px {
		pc.Ppx++
		el.ID = pc.Ppx
		pc.Proxies[pc.Ppx] = el
	}
}

func (pc *ProxyController) GetProxy() (*Proxy, error) {
	if len(pc.Proxies) < 1 {
		return nil, ErrNoProxies
	}

	k := rand.Intn(len(pc.Proxies))
	for _, x := range pc.Proxies {
		if k == 0 {
			return x, nil
		}
		k--
	}
	return nil, ErrUnreachable
}

func (p *Proxy) ToString() string {
	return fmt.Sprintf("%s://%s:%d?timeout=5s", proxyTypes[p.Type], p.IP, p.Port)
}

func (p *Proxy) Error() {
	atomic.AddInt64(&p.Errors, 1)
}
