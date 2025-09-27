package config

var lambdaConfig = Config{
	ServerApi: ServerHttp{
		Server: Server{
			Name: "desert-lambda-api",
			Mode: "release",
		},
		Listen:                    ":80",
		JWTKey:                    "lambda-desert",
		JWTAccessTokenExpireTime:  86400,
		JWTRefreshTokenExpireTime: 86400 * 30,
	},
	Otel: Otel{
		Trace: OtelConfig{
			ServiceName:  "desert-lambda-api",
			ExporterType: "console",
		},
		Log: OtelConfig{
			ServiceName:  "desert-lambda-api",
			ExporterType: "console",
		},
		Metric: OtelConfig{
			ServiceName:  "desert-lambda-api",
			ExporterType: "console",
		},
	},
}
