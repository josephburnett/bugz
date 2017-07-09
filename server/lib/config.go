package colony

import (
	"encoding/json"
	"flag"
)

var Config = struct {
	WorldFile *string
	Ip        *string
	Port      *string
}{
	flag.String("world_file", "", "File for persistent world state."),
	flag.String("ip", "0.0.0.0", "HTTP server ip."),
	flag.String("port", "8080", "HTTP server port."),
}

func ConfigJson() ([]byte, error) {
	return json.Marshal(Config)
}
