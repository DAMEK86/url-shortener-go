package config

var (
	BuildVersion  = "local-BuildVersion"
	SourceVersion = "local-SourceVersion"
)

type Config struct {
	Swagger StaticContent `json:"swagger"`
}

type StaticContent struct {
	Disabled bool   `json:"disabled"`
	FilePath string `json:"filePath"` // path where the static content html files are located (absolute or relative to the service binary)
	HttpPath string `json:"httpPath"` // http route under that the content shall be hosted
}
