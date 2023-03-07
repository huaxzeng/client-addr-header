package traefik_plugin_client_addr_header

import (
	"context"
	"fmt"
	"net"
	"net/http"
)

// Config the plugin configuration.
type Config struct {
	Host	string	`json:"host,omitempty" toml:"host,omitempty" yaml:"host,omitempty"`
	Port	string	`json:"port,omitempty" toml:"port,omitempty" yaml:"port,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Host: "",
		Port: "",
	}
}

// ClientAddrHeader a ClientAddrHeader plugin.
type ClientAddrHeader struct {
	next   http.Handler
	name   string
	config *Config
}

// New created a new ClientAddrHeader plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if config.Host == "" {
		return nil, fmt.Errorf("host cannot be empty")
	}

	if config.Port != "" && config.Host == config.Port {
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

	req.Header.Set(a.config.Host, host)
	if a.config.Port != "" {
		req.Header.Set(a.config.Port, port)
	}
	a.next.ServeHTTP(rw, req)
}
