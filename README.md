# Traefik Plugin: Client Address Header

This is a fork from https://github.com/bonsai-oss/custom-source-header

Client Address Header is a middleware plugin for Traefik. It injects the remote address and port of the client into selectable headers.

## Why does this exist?
When operating a kubernetes environment inside a private network combined with a second layer of Traefik, the remote address of the client is not available to the application due to header overwriting.

## Configuration

To configure this plugin you should add its configuration to the Traefik dynamic configuration as explained [here](https://docs.traefik.io/getting-started/configuration-overview/#the-dynamic-configuration).
The following snippet shows how to configure this plugin with the File provider in TOML and YAML:

Static:

```toml
[experimental.plugins.clientaddrheader]
  modulename = "github.com/huaxzeng/client-addr-header"
  version = "v0.0.1"
```

Dynamic:

```toml
[http.middlewares]
  [http.middlewares.injectclientaddrheaders.plugin.clientaddrheader]
    host = "X-Client-IP"
    port = "X-Client-Port"
```

```yaml
http:
  middlewares:
   injectclientaddrheaders:
      plugin:
        clientaddrheader:
          host: "X-Client-IP"
          port: "X-Client-Port"
```