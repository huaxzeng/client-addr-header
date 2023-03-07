package traefik_plugin_client_addr_header

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func dummyHandler(w http.ResponseWriter, r *http.Request) {
	headers := r.Header
	json.NewEncoder(w).Encode(headers)
	return
}

func TestClientAddrHeader_ServeHTTP(t *testing.T) {
	for _, tt := range []struct {
		name         string
		pluginConfig *Config
	}{
		{
			name: "both host and port headers",
			pluginConfig: &Config{
				host: "X-Remote-Ip",
				port: "X-Remote-Port",
			},
		},
		{
			name: "only host header",
			pluginConfig: &Config{
				host: "X-Remote-Ip",
			},
		},
		{
			name: "host and port header with non standard casing",
			pluginConfig: &Config{
				host: "x-client-host",
                port: "X-CLIENT-PORT",
			},
		},

	} {
		t.Run(tt.name, func(t *testing.T) {
			pluginHandler, pluginHandlerCreateError := New(context.Background(), http.HandlerFunc(dummyHandler), tt.pluginConfig, tt.name)
			if pluginHandlerCreateError != nil {
				t.Fatal(pluginHandlerCreateError)
			}

			svr := httptest.NewServer(
				pluginHandler,
			)
			defer svr.Close()

			req, err := http.NewRequest("GET", svr.URL, nil)
			if err != nil {
				t.Fatal(err)
			}
			rsp, _ := (&http.Client{}).Do(req)
			defer rsp.Body.Close()

			responseHeaderData := make(map[string][]string)
			json.NewDecoder(rsp.Body).Decode(&responseHeaderData)

            t.Logf("response header: %s", responseHeaderData)

            hostHeader := http.CanonicalHeaderKey(tt.pluginConfig.host)
			if _, ok := responseHeaderData[hostHeader]; !ok {
				t.Errorf("expected header %s to be set", hostHeader)
			}

			if tt.pluginConfig.port != "" {
                portHeader := http.CanonicalHeaderKey(tt.pluginConfig.port)
				if _, ok := responseHeaderData[portHeader]; !ok {
					t.Errorf("expected header %s to be set", portHeader)
				}
			}
		})
	}
}

func TestCreateConfig(t *testing.T) {
	config := CreateConfig()

	if fmt.Sprintf("%T", config) != "*client_addr_header.Config" {
		t.Errorf("expected config to be of type *client_addr_header.Config")
	}
}
