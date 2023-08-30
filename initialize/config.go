package initialize

import (
	"claude2/model"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var ServerConfig = model.ServerConfig {
	Claude: &model.Claude{},
}

const (
	ConfigDefaultFile = "config.yaml"
)

var (
	listenHostflag = ""
	tlsCertFlag   = ""
	tlsKeyFlag    = ""
	configFlag    = ""
	httpProxyFlag = ""
	baseUrlFlag   = ""
)

func init() {
	flag.StringVar(&listenHostflag, "listen_host", listenHostflag, "Listen host, for example 0.0.0.0:8000")
	flag.StringVar(&tlsCertFlag, "tls_cert", tlsCertFlag, "TLS cert path")
	flag.StringVar(&tlsKeyFlag, "tls_key", tlsKeyFlag, "TLS key path")
	flag.StringVar(&configFlag, "c", configFlag, "Choose config file.")
	flag.StringVar(&httpProxyFlag, "http_proxy", httpProxyFlag, "Set http_proxy, for example http://127.0.0.1:8000")
	flag.StringVar(&baseUrlFlag, "base_url", baseUrlFlag, "Set base_url, for example https://claude.ai")
}

func NewViper() {
	flag.Parse()
	if configFlag == "" {
		configFlag = ConfigDefaultFile
	}
	// check config file
	_, err := os.Stat(configFlag)
	if os.IsNotExist(err) {
		file, err := os.Create(configFlag)
		// 其他处理
		if err != nil {
			return
		}
		defer file.Close()
		encoder := yaml.NewEncoder(file)
		encoder.SetIndent(2)
		if err := encoder.Encode(&ServerConfig); err != nil {
			panic(err)
		}
		fmt.Println("File created and data written successfully.")
	}
	v := viper.New()
	v.SetConfigFile(configFlag)
	v.SetConfigType("yaml")
	// 设置默认值
	v.SetDefault("base-url", "https://claude.ai")
	v.SetDefault("listen-host", "0.0.0.0:8000")
	err = v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		// 因为viper值如果为空（删除）不会复写原来的值，数组内的值删除会出现不生效问题，SessionKeys先置为空再赋值
		ServerConfig.Claude.SessionKeys = nil
		if err = v.Unmarshal(&ServerConfig); err != nil {
			fmt.Println(err)
		}
		PrintServerConfig()
	})
	if err = v.Unmarshal(&ServerConfig); err != nil {
		fmt.Println(err)
	}
	// 设置命令参数
	if baseUrlFlag != "" {
		ServerConfig.BaseUrl = baseUrlFlag
	}
	if httpProxyFlag != "" {
		ServerConfig.HttpProxy = httpProxyFlag
	}
	if tlsCertFlag != "" {
		ServerConfig.TlsCert = tlsCertFlag
	}
	if tlsKeyFlag != "" {
		ServerConfig.TlsKey = tlsKeyFlag
	}
	if listenHostflag != "" {
		ServerConfig.ListenHost = listenHostflag
	}
	// 设置环境变量
	keysEnv := os.Getenv("CLAUDE_SESSION_KEYS")
	if keysEnv != "" {
		keys := strings.Split(keysEnv, ",")
		ServerConfig.Claude.SessionKeys = append(ServerConfig.Claude.SessionKeys, keys...)
	}
	baseUrlEnv := os.Getenv("CLAUDE_BASE_URL")
	if baseUrlEnv != "" {
		ServerConfig.BaseUrl = baseUrlEnv
	}
	httpProxyEnv := os.Getenv("CLAUDE_HTTP_PROXY")
	if httpProxyEnv != "" {
		ServerConfig.HttpProxy = httpProxyEnv
	}

	tlsCertEnv := os.Getenv("CLAUDE_TLS_CERT")
	if tlsCertEnv != "" {
		ServerConfig.TlsCert = tlsCertEnv
	}

	tlsKeyEnv := os.Getenv("CLAUDE_TLS_KEY")
	if tlsKeyEnv != "" {
		ServerConfig.TlsKey = tlsKeyEnv
	}

	listenHostEnv := os.Getenv("CLAUDE_LISTEN_HOST")
	if listenHostEnv != "" {
		ServerConfig.ListenHost = listenHostEnv
	}
}

func PrintServerConfig() {
	indent, _ := json.MarshalIndent(ServerConfig, "", "    ")
	fmt.Println(string(indent))
}
