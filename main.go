package main

import (
	"log"
	"net/http"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/rs/cors"
	"github.com/rs/xlog"
	"github.com/secsy/goftp"

	"github.com/incu6us/asterisk-dialer/platform/api"
	"github.com/incu6us/asterisk-dialer/platform/database"
	"github.com/incu6us/asterisk-dialer/platform/dialer"
	"github.com/incu6us/asterisk-dialer/platform/ftpclient"
	"github.com/incu6us/asterisk-dialer/platform/scheduler"
	"github.com/incu6us/asterisk-dialer/utils/config"
	"github.com/incu6us/asterisk-dialer/utils/srvlogger"
)

func main() {

	var conf = config.GetConfig()
	var err error

	log.SetFlags(0)
	xlogger := xlog.New(srvlogger.LoggerConfig(config.GetConfig()))
	log.SetOutput(xlogger)

	if err = new(database.DB).Connect(conf); err != nil {
		xlog.Errorf("DB connection problems: %s", err)
	}
	dialerInstance := dialer.CreateDialer(conf, database.GetDB())
	if err = dialerInstance.Connect(); err != nil {
		xlog.Errorf("Error: %s", err)
	}
	dialerInstance.StartDialing()

	srv := &http.Server{
		Addr: conf.General.Listen,
		//ReadTimeout:  api.HTTP_TIMEOUT,
		//WriteTimeout: api.HTTP_TIMEOUT,
		Handler: cors.Default().Handler(api.NewHandler(database.GetDB())),
	}

	sc := scheduler.SchedulerProvider()
	sc.Run(getMsisdnsFromFtp, 30)
	sc.Run(database.GetDB().DeleteMSISDNOlderThenWeek, 3600)
	sc.Start()

	if err = srv.ListenAndServe(); err != nil {
		xlog.Fatal(err)
	}
}

func getMsisdnsFromFtp() {

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
		if err := database.GetDB().AddNewNumbers(nums); err != nil {
			xlog.Error(err)
		}
	} else {
		xlog.Debug("No files on FTP server")
	}
}
