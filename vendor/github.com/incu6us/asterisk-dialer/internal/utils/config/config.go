package config

import (
	"flag"
	"io/ioutil"
	"os"
	"time"

	"github.com/naoina/toml"
	"github.com/rs/xlog"
)

// TomlConfig - config structure
type TomlConfig struct {
	General struct {
		Listen      string
		CallTimeout time.Duration `toml:"call-timeout"`
		LogLevel    xlog.Level    `toml:"log-level"`
	}
	Ami struct {
		Host     string
		Port     int
		Username string
		Password string
	} `toml:"ami"`
	DB struct {
		Host     string
		Database string
		Username string
		Password string
		Debug    bool
	} `toml:"db"`
	Asterisk struct {
		Context         string `toml:"call-context"`
		PlaybackContext string `toml:"playback-context"`
	}
	Mq struct {
		User     string
		Password string
		Host     string
		Port     string
	}
	Ftp struct {
		Host     string
		User     string
		Password string
	}
}

var config *TomlConfig

//var logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile)

// GetConfig - get current TOML-config
func GetConfig() *TomlConfig {
	if config == nil {
		config = new(TomlConfig).readConfig()
	}

	return config
}

func (t *TomlConfig) readConfig() *TomlConfig {
	confFlag := t.readFlags()
	f, err := os.Open(*confFlag)
	if err != nil {
		xlog.Errorf("Can't open file: %v", err)
	}

	defer f.Close()

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		xlog.Errorf("IO read error: %v", err)
	}

	if err := toml.Unmarshal(buf, t); err != nil {
		xlog.Errorf("TOML error: %v", err)
	}

	return t
}

func (t *TomlConfig) readFlags() (confFlag *string) {
	confFlag = flag.String("conf", "app.conf", "app.conf")
	flag.Parse()
	return
}
