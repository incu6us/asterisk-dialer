package dialer_helper

import (
	"errors"

	"github.com/rs/xlog"

	"github.com/incu6us/asterisk-dialer/internal/platform/database"
	"github.com/incu6us/asterisk-dialer/internal/platform/dialer"
)

// AddAndRun - add numbers & run dialer
func AddAndRun(nums []string) error {
	var err error

	if len(nums) > 0 {
		xlog.Info("New files founded on FTP")

		if err = database.AddNewNumbers(nums); err != nil {
			xlog.Error(err)
		}

		// When AMI is connected, then start dialer
		switch dialer.GetDialer().GetAmiConnectionStatus() {
		case true:
			xlog.Debug("AMI is connected")

			if !dialer.GetDialer().GetDialerStatus() {
				xlog.Info("Starting dialer...")
				dialer.GetDialer().StartDialing()
			}

		case false:
			xlog.Error("AMI is not connected! Can't start dialer")
		}

	} else {
		err = errors.New("Number list is empty")
	}

	return err
}

func Call() {
	// When AMI is connected, then start dialer
	switch dialer.GetDialer().GetAmiConnectionStatus() {
	case true:
		xlog.Debug("AMI is connected")

		if !dialer.GetDialer().GetDialerStatus() {
			xlog.Info("Starting dialer...")
			dialer.GetDialer().StartDialing()
		}

	case false:
		xlog.Error("AMI is not connected! Can't start dialer")
	}
}
