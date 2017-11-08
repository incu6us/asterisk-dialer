package srvlogger

import (
	"github.com/rs/xlog"

	"github.com/incu6us/asterisk-dialer/utils/config"
)

func LoggerConfig(conf *config.TomlConfig) xlog.Config {

	xlog.Debugf("LogLevel=%s", conf.General.LogLevel)
	return xlog.Config{
		// Log info level and higher
		Level: conf.General.LogLevel,
		// Output everything on console
		Output: xlog.NewOutputChannel(xlog.NewConsoleOutput()),
	}
}
