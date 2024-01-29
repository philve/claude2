# Claude2

[Claude2](https://claude.ai) Convert to OpenAI API standard interface，Origin: [claude-to-chatgpt](https://github.com/oldweipro/claude-to-chatgpt)

```shell
curl https://claude2-0bbi.onrender.com/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-xxxxxxxxxxxxxxxx" \
  -d '{
    "model": "gpt-3.5-turbo",
    "messages": [{"role": "system", "content": "You are a helpful assistant."}, {"role": "user", "content": "Hello!"}],
    "stream": true
  }'
```

## Compile

[Go](https://go.dev/dl/) 1.20 及以上版本。

```
git clone https://github.com/gngpp/claude2.git && cd claude2
go mod tidy
go build -ldflags "-s -w -extldflags -static" main.go
```

### other

使用 `-c` 指定配置文件 `config-dev.yaml`

使用 `-http_proxy` 设置 `http_proxy` 例如 `http://127.0.0.1:8000`

```shell
go run main.go -c config-dev.yaml -http_proxy http://127.0.0.1:8000
```

## arrangement

If the configuration file does not exist, the program will automatically create 'config.yaml'.

If the configuration information filled in after startup is incorrect, you can directly modify the configuration file and save it, and the program will automatically reload.

| Configuration items            | Description                                                                                                                                                                                   | Example values                  |
|----------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-----------------------|  
| claude         | Claude related configuration                                                                                                                                                                           |                       |
| session-keys | Array of unique identifiers of the current conversation session,
Optional<br/>Support setting `Bearer sessionKey` in `Header Authorization`<br/>Reference [Authentication](https://platform.openai.com/docs/api-reference/authentication)<br/>Header priority level greater than profile | [sk-1, sk-2]          | 
| http-proxy     | Proxy configuration, optional<br/>(Including but not limited to) Pay attention to the connectivity in Docker<br/>You may need to replace `http://127.0.0.1:8000` with the host IP<br/>For example` http://192.168.1.2:8000`                            | http://127.0.0.1:8000 |
| tls-cert     | TLS certificate path                                                                        | /etc/tls.pem |
| tls-key     | TLS certificate key path                                                                      | /etc/tls.key |
| listen-host     | Listening host                                                                     | 0.0.0.0:8000 |
