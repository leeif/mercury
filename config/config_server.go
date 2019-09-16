package config

type ServerConfig struct {
	API APIServer `kiper_config:"name:api"`
	WS  WSServer  `kiper_config:"name:ws"`
}

type APIServer struct {
	Port *Port `kiper_value:"name:port;help:server listen port"`
}

type WSServer struct {
	Port *Port `kiper_value:"name:port;help:server listen port"`
}

type Port struct {
	s string
}

func (p *Port) Set(s string) error {
	p.s = s
	return nil
}

func (p *Port) String() string {
	return p.s
}

func newServerConfig() ServerConfig {
	return ServerConfig{
		API: APIServer{
			Port: &Port{},
		},
		WS: WSServer{
			Port: &Port{},
		},
	}
}
