package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/incu6us/asterisk-dialer/internal/platform/database"
	"github.com/incu6us/asterisk-dialer/internal/platform/dialer"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/rs/xlog"
	"github.com/secsy/goftp"

	"github.com/incu6us/asterisk-dialer/internal/platform/api"
	"github.com/incu6us/asterisk-dialer/internal/platform/ftpclient"
	"github.com/incu6us/asterisk-dialer/internal/platform/scheduler"
	"github.com/incu6us/asterisk-dialer/internal/utils/config"
	"github.com/incu6us/asterisk-dialer/internal/utils/dialer-helper"
	"github.com/incu6us/asterisk-dialer/internal/utils/srvlogger"
)

func main() {

	var conf = config.GetConfig()
	var err error

	log.SetFlags(0)
	xlogger := xlog.New(srvlogger.LoggerConfig(config.GetConfig()))
	log.SetOutput(xlogger)

	var host = conf.Ami.Host + ":" + strconv.Itoa(conf.Ami.Port)

	dialerInstance := dialer.CreateDialer(host, conf.Ami.Username, conf.Ami.Password)
	if err = dialerInstance.Connect(); err != nil {
		xlog.Errorf("Error: %s", err)
	}
	dialerInstance.StartDialing()

	srv := &http.Server{
		Addr: conf.General.Listen,
		//ReadTimeout:  api.HTTP_TIMEOUT,
		//WriteTimeout: api.HTTP_TIMEOUT,
		Handler: api.NewHandler(),
	}

	sc := scheduler.SchedulerProvider()
	//sc.Run(dialer_helper.Run, 60)
	sc.Run(callToNumbersFromFtp, 30)
	sc.Run(database.DeleteMSISDNOlderThenWeek, 3600)
	sc.Start()

	if err = srv.ListenAndServe(); err != nil {
		xlog.Fatal(err)
	}
}

func callToNumbersFromFtp() {

	if config.GetConfig().Ftp.Host == "" {
		return
	}

	xlog.Debug("Checking on FTP")

	ftpClient := ftpclient.FtpClient{Client: new(goftp.Client), Config: config.GetConfig()}
	if err := ftpClient.Connect(true); err != nil {
		xlog.Error(err)
	}

	defer ftpClient.Disconnect()

	nums, err := ftpClient.GetNumbers()
	if err != nil {
		xlog.Error(err)
	}

	//
	if len(nums) > 0 {
		if err := dialer_helper.AddAndRun(nums); err != nil {
			xlog.Error(err)
		}
	} else {
		xlog.Debug("No files on FTP server")
	}
}
