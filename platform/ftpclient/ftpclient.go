package ftpclient

import (
	"bufio"
	"bytes"

	"github.com/rs/xlog"
	"github.com/secsy/goftp"

	"github.com/incu6us/asterisk-dialer/utils/config"
	"github.com/incu6us/asterisk-dialer/utils/file-utils"
)

type FtpClient struct {
	Client *goftp.Client
	Config *config.TomlConfig
}

var connected bool

func (f *FtpClient) Connect(activeMode bool) error {

	if connected {
		return nil
	}

	ftpConfig := goftp.Config{
		User:            f.Config.Ftp.User,
		Password:        f.Config.Ftp.Password,
		ActiveTransfers: activeMode,
	}

	client, err := goftp.DialConfig(ftpConfig, f.Config.Ftp.Host)
	if err != nil {
		return err
	}

	connected = true
	f.Client = client

	return nil
}

func (f *FtpClient) GetNumbers() ([]string, error) {

	entries, _ := f.Client.ReadDir(".")

	var numberList []string

	for _, entry := range entries {
		xlog.Infof("Found a new file on FTP: %s", entry.Name())
		name := entry.Name()

		var b bytes.Buffer
		writer := bufio.NewWriter(&b)

		err := f.Client.Retrieve(name, writer)
		if err != nil {
			return nil, err
		}

		numberList = file_utils.MSISDNNormalizer(b.String())

		f.Client.Delete(name)

	}

	return numberList, nil
}

func (f *FtpClient) Disconnect() {
	connected = false

	if err := f.Client.Close(); err != nil {
		xlog.Error(err)
	}
}
