# Claude2

[Claude2](https://claude.ai) 聊天功能接口转 OpenAI API 标准接口

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

## 编译

[Go](https://go.dev/dl/) 1.20 及以上版本。

```
git clone https://github.com/gngpp/claude2.git && cd claude2
go mod tidy
go build -ldflags "-s -w -extldflags -static -extldflags -static" main.go
```

### 其他

使用 `-c` 指定配置文件 `config-dev.yaml`

使用 `-http_proxy` 设置 `http_proxy` 例如 `http://127.0.0.1:8000`

```shell
go run main.go -c config-dev.yaml -http_proxy http://127.0.0.1:8000
```

## 配置

配置文件如果不存在,程序会自动创建 `config.yaml`。

如果启动后填写的配置信息有误,直接修改配置文件并保存即可,程序会自动重新加载。

| 配置项            | 说明                                                                                                                                                                                     | 示例值                   |
|----------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-----------------------|  
| base-url       | Claude服务地址,可选                                                                                                                                                                          | https://claude.ai     |
| claude         | Claude 相关配置                                                                                                                                                                            |                       |
| - session-keys | 当前对话session唯一标识数组,必填<br/>支持在 `Header Authorization` 中设置 `Bearer sessionKey`<br/>参考 [Authentication](https://platform.openai.com/docs/api-reference/authentication)<br/>Header优先级大于配置文件 | [sk-1, sk-2]          | 
| http-proxy     | 代理配置,可选<br/>(包含但不限于)注意在Docker中的连通性<br/>可能需要更换`http://127.0.0.1:8000`为宿主机IP<br/>如`http://192.168.1.2:8000`                                                                              | http://127.0.0.1:8000 |

原项目：[claude-to-chatgpt](https://github.com/oldweipro/claude-to-chatgpt)
