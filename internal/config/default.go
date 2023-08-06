package config

var DefaultConfig = Config{
	Swagger: StaticContent{
		Disabled: false,
		FilePath: "./staticcontent/swagger-ui/",
		HttpPath: "/swagger-ui/",
	},
}
