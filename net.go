package main

import (

    "net"
    "time"

    "code.google.com/p/go.net/proxy"
)

type ConnHandler struct {

    ProxyType string
    ProxyAddress string
    ProxyUsername string
    ProxyPassword string
    ProxyTimeout int
}

func (handler *ConnHandler) HandleConnection(network, addr string) (conn net.Conn, err error) {

    forwardDialer := &net.Dialer {

        Timeout: time.Duration(handler.ProxyTimeout) * time.Minute,
        DualStack: true,
    }

    if (len(handler.ProxyAddress) > 1) {

        auth := &proxy.Auth {

            User: handler.ProxyUsername,
            Password: handler.ProxyPassword,
        }

        // setup the socks proxy
        dialer, err := proxy.SOCKS5("tcp", handler.ProxyAddress, auth, forwardDialer)
        if (err != nil) {

            return nil, err
        }

        return dialer.Dial(network, addr)
    }

    return forwardDialer.Dial(network, addr)
}