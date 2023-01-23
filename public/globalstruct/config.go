package globalstruct

type SystemConfigStruct struct {
	Project    map[string]interface{}
	Author     map[string]interface{}
	Paging     int    `json:"Paging"`
	Theme      string `json:"Theme"`
	Service    Services
	TimeFormat string            `json:"TimeFormat"`
	ServerAddr string            `json:"ServerAddr"`
	Renderer   map[string]string `json:"Renderer,omitempty"`
}

type RedisConfigStruct struct {
	Host     string `json:"Host"`
	Port     int    `json:"Port"`
	Password string `json:"Password"`
	DB       int    `json:"DB"`
}

type ControlConfigStruct struct {
	Password string `json:"Password"`
}

func DefaultSystemConfig() SystemConfigStruct {
	var config SystemConfigStruct
	config.ServerAddr = "0.0.0.0:9000"
	config.Service.Control = false
	return config
}

type Services struct {
	Control bool `json:"Control"`
	Comment bool `json:"Comment"`
}
