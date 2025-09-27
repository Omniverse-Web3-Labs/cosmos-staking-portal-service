package config

import (
	"app/pkg/server"
	"app/pkg/utils"
	"app/pkg/utils/otelutil"
	"flag"
	"log"
	"os"
	"path"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ServerApi ServerHttp `yaml:"serverApi"`
	Otel      Otel       `yaml:"otel"`
	Databases struct {
		Postgres struct {
			Default  Database `yaml:"default"`
			Subquery Database `yaml:"subquery"`
		} `yaml:"postgres"`
	} `yaml:"databases"`
	Endpoint Endpoint `yaml:"endpoint"`
	Redis    Redis    `yaml:"redis"`
}

type Redis struct {
	Endpoint string `yaml:"endpoint"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	TLS      bool   `yaml:"tls"`
}

type Endpoint struct {
	Rest     string `yaml:"rest"`
	GRpc     string `yaml:"grpc"`
	CometBFT string `yaml:"cometbft"`
	Graph    string `yaml:"graph"`
}

type Server struct {
	Name string `yaml:"name"`
	Mode string `yaml:"mode"`
}

type ServerHttp struct {
	Server                    `yaml:",inline"`
	Listen                    string `yaml:"listen"`
	JWTKey                    string `yaml:"jwtKey"`
	JWTAccessTokenExpireTime  int    `yaml:"jwtAccessTokenExpireTime"`
	JWTRefreshTokenExpireTime int    `yaml:"jwtRefreshTokenExpireTime"`
}

type Database struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DBName   string `yaml:"dbname"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Otel struct {
	TlsCertPath string     `yaml:"tlsCertPath"`
	Trace       OtelConfig `yaml:"trace"`
	Log         OtelConfig `yaml:"log"`
	Metric      OtelConfig `yaml:"metric"`
}
type OtelConfig struct {
	ServiceName      string `yaml:"serviceName"`
	ExporterType     string `yaml:"exporterType"`
	ExporterEndpoint string `yaml:"exporterEndpoint"`
}

var config Config = Config{
	ServerApi: ServerHttp{
		Listen:                    ":8080",
		JWTAccessTokenExpireTime:  300,
		JWTRefreshTokenExpireTime: 31536000,
	},
}

func Parse(filePath string) {
	if server.InLambda() {
		config = lambdaConfig
	} else {
		var configFile string
		flag.StringVar(&configFile, "c", "", "Config file")
		flag.Parse()
		filepathes := []string{
			"configs/config.yaml",
			path.Join(utils.FileUtil.ExecuteDir(), "configs/config.yaml"),
		}
		if configFile != "" {
			filepathes = []string{
				configFile,
			}
		}
		if filePath != "" {
			filepathes = []string{
				filePath,
			}
		}
		LoadConfig(&config, filepathes)
	}
}

func LoadConfig(config interface{}, filepathes []string) {
	for i, filepath := range filepathes {
		log.Println("read config file", filepath)
		data, err := os.ReadFile(filepath)
		if err != nil {
			if i == len(filepathes)-1 {
				log.Fatal(err)
			} else {
				log.Println(err.Error())
			}
			continue
		}
		log.Println("read config file success", filepath)
		err = yaml.Unmarshal(data, config)
		if err != nil {
			log.Fatal(err)
		}
		break
	}
}

func GetConfig() *Config {
	return &config
}

func GetOtelOptions() otelutil.Options {
	otelOptions := otelutil.Options{
		TraceExporterType:      GetConfig().Otel.Trace.ExporterType,
		TraceExporterEndpoint:  GetConfig().Otel.Trace.ExporterEndpoint,
		TraceExportTimeout:     time.Minute * 1,
		MetricExporterType:     GetConfig().Otel.Metric.ExporterType,
		MetricExporterEndpoint: GetConfig().Otel.Metric.ExporterEndpoint,
		MetricExportInterval:   time.Minute * 5,
		LoggerExporterType:     GetConfig().Otel.Log.ExporterType,
		LoggerExporterEndpoint: GetConfig().Otel.Log.ExporterEndpoint,
		LoggerExportTimeout:    time.Minute * 5,
		TlsCertPath:            GetConfig().Otel.TlsCertPath,
	}
	return otelOptions
}
