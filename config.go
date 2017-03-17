//    Copyright 2017 drillbits
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package lambique

import (
	"os/user"
	"strings"

	"github.com/BurntSushi/toml"
)

var (
	cfg         = defaultConfig()
	defaultAddr = ":2697" // \u2697
)

// Config is a config for web application.
type Config struct {
	Addr string `toml:"address"`
}

func defaultConfig() *Config {
	return &Config{
		Addr: defaultAddr,
	}
}

// LoadConfig reads the config file from path and returns the config.
func LoadConfig(path string) (*Config, error) {
	if path == "" {
		return cfg, nil
	}

	usr, _ := user.Current()
	path = strings.Replace(path, "~", usr.HomeDir, 1)

	_, err := toml.DecodeFile(path, &cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// GetConfig returns the config.
func GetConfig() *Config {
	return cfg
}
