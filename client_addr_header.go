package client_addr_header

import (
	"context"
	"fmt"
	"net"
	"net/http"
)

// Config the plugin configuration.
type Config struct {
	host	string	`json:",omitempty"`
	port	string	`json:",omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

// ClientAddrHeader a ClientAddrHeader plugin.
type ClientAddrHeader struct {
	next   http.Handler
	name   string
	config *Config
}

// New created a new ClientAddrHeader plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if config.host == "" {
		return nil, fmt.Errorf("host cannot be empty")
	}

	if config.port != "" && config.host == config.port {
		return nil, fmt.Errorf("host cannot be the same as port")
	}

	return &ClientAddrHeader{
		config: config,
		next:   next,
		name:   name,
	}, nil
}

func (a *ClientAddrHeader) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	remoteAddr := req.RemoteAddr
	host, port, splitError := net.SplitHostPort(remoteAddr)
	if splitError != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	req.Header.Set(a.config.host, host)
	if a.config.port != "" {
		req.Header.Set(a.config.port, port)
	}
	a.next.ServeHTTP(rw, req)
}
