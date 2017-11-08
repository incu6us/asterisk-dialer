package dialer

import (
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/ivahaev/amigo"
	"github.com/ivahaev/amigo/uuid"
	"github.com/rs/xlog"

	"github.com/incu6us/asterisk-dialer/internal/platform/database"
	"github.com/incu6us/asterisk-dialer/internal/utils/config"
)

type Dialer struct {
	dialerStatus        bool
	amiConnectionStatus bool
	randGenDigit        int
	host                string
	user                string
	pass                string
}

var (
	once        sync.Once
	amiInstance *Dialer
	amigoClient *amigo.Amigo
	conf        = config.GetConfig()
)

var (
	userIsCalling map[string]bool
)

func (a *Dialer) peerDial(message map[string]string) {
	user := database.GetAutodialUser(strings.TrimPrefix(message["Peer"], "SIP/"))
	if user.Peer != "" &&
		!userIsCalling[strings.TrimPrefix(message["Peer"], "SIP/")] &&
		a.dialerStatus {
		userIsCalling[strings.TrimPrefix(message["Peer"], "SIP/")] = true

		xlog.Infof("Dial: %#v", message)
	}
}

func (a *Dialer) peerHangup(message map[string]string) {
	user := database.GetAutodialUser(message["CallerIDNum"])
	if user.Peer != "" {

		callerIDNum := message["CallerIDNum"]
		callerIDName := message["CallerIDName"]
		cause := message["Cause"]
		causeTxt := message["Cause-txt"]
		event := message["Event"]
		channel := message["Channel"]
		uniqueid := message["Uniqueid"]

		xlog.Debugf("peer hangup: %#v", message)

		database.UpdateAfterHangup(callerIDNum, callerIDName, cause, causeTxt, event, channel, uniqueid)
		if userIsCalling[message["CallerIDNum"]] {
			database.UpdatePeerStatus(user.Peer, "", event, "")
			xlog.Infof("Hangup: %#v", message)
		}

		userIsCalling[message["CallerIDNum"]] = false
		xlog.Debugf("Is user calling: %t", userIsCalling[message["CallerIDNum"]])

		if a.dialerStatus {
			timeout := config.GetConfig().General.CallTimeout * time.Second
			xlog.Info("Sleeping " + timeout.String())

			time.Sleep(timeout)
			_, err := a.CallToUser(callerIDNum)
			if err != nil {
				xlog.Warnf("CallToUser problem: %s", err)
			}
		}

	}
}

func (a *Dialer) peerStatusListening(message map[string]string) {

	user := database.GetAutodialUser(strings.TrimPrefix(message["Peer"], "SIP/"))
	if user.Peer != "" && a.dialerStatus {

		xlog.Infof("Peer status changed -> %v", message)
		database.UpdatePeerStatus(user.Peer, message["PeerStatus"], "", "")

		switch message["PeerStatus"] {
		case "Reachable", "Registered":
			time.Sleep(5 * time.Second)
			_, err := a.CallToUser(user.Peer)
			if err != nil {
				xlog.Warnf("CallToUser problem: %s", err)
			}
		}
	}
}

func (a *Dialer) Connect() error {

	userIsCalling = make(map[string]bool)

	var err error

	host := strings.Split(a.host, ":")[0]
	port := strings.Split(a.host, ":")[1]

	settings := &amigo.Settings{Username: a.user, Password: a.pass, Host: host, Port: port}
	amigoClient = amigo.New(settings)

	amigoClient.Connect()

	// Listen for connection events
	amigoClient.On("connect", func(message string) {
		xlog.Infof("Connected: %s", message)
	})
	amigoClient.On("error", func(message string) {
		xlog.Errorf("Connection error: %s", message)
		err = errors.New(message)
	})

	// registering and un-registering sip users
	//amigoClient.RegisterHandler("PeerStatus", func(message map[string]string) {
	//	if message["Peer"] == "" {
	//		return
	//	}
	//	user := database.GetAutodialUser(strings.TrimPrefix(message["Peer"], "SIP/"))
	//	database.UpdatePeerStatus(user.Peer, message["PeerStatus"], "", "")
	//	xlog.Infof("user event: %#v", message)
	//})

	amigoClient.RegisterHandler("PeerStatus", a.peerStatusListening)
	amigoClient.RegisterHandler("Hangup", a.peerHangup)

	return err
}

// TODO: add userIsCalling to skip parallel calls
func (a *Dialer) callToAll() {
	users := database.GetRegisteredUsers()
	for _, user := range users {
		xlog.Debugf("Calling to %s", user.Peer)
		a.CallToUser(user.Peer)
	}
}

func (a *Dialer) originate(params map[string]string) (map[string]string, error) {
	params["ActionID"] = uuid.NewV4()
	params["Variable"] = "ActionID=" + params["ActionID"]
	params["Action"] = "Originate"

	xlog.Debugf("Originate: %#v", params)

	resp, err := amigoClient.Action(params)
	if err != nil {
		xlog.Error(err)
		return nil, err
	}

	return resp, nil
}

func (a *Dialer) CallToUser(userID string) (map[string]string, error) {
	// if user already calling by autodialer, then skip the new call request
	if userIsCalling[userID] {
		return nil, errors.New("already calling by dialer")
	}

	if msisdn := database.ProcesseMsisdn(userID); msisdn != "" {
		userIsCalling[userID] = true
		xlog.Debugf("Calling to %s with exten: %s", userID, msisdn)

		var params = make(map[string]string)
		params["Channel"] = "SIP/" + userID
		params["CallerID"] = "autodial_" + msisdn
		params["MaxRetries"] = "0"
		params["RetryTime"] = "0"
		params["WaitTime"] = "300"
		params["Context"] = conf.Asterisk.Context
		params["Exten"] = msisdn
		params["Priority"] = "1"
		params["Async"] = "true"

		database.UpdatePeerStatus(userID, "", "Originate", msisdn)

		amiResponse, err := a.originate(params)
		if err != nil {
			userIsCalling[userID] = false
			xlog.Errorf("AMI Action error! Error: %v, AMI Response Status: %s", err, amiResponse)
			return nil, err
		}

		return amiResponse, err
	}

	return nil, nil
}

func (a *Dialer) StartDialing() string {

	if !a.dialerStatus {
		xlog.Debug("Dialer has been started...")

		userIsCalling = make(map[string]bool)

		amigoClient.RegisterHandler("Dial", a.peerDial)

		a.dialerStatus = true

		return "started"
	}

	xlog.Warn("Dialer is already started!")
	return "Can't run. Already started"
}

func (a *Dialer) StopDialing() string {

	if a.dialerStatus {

		amigoClient.UnregisterHandler("Dial", a.peerDial)

		a.dialerStatus = false
	}

	xlog.Debug("Dialer has been stopped")
	return "stopped"
}

func (a *Dialer) GetDialerStatus() bool {
	return a.dialerStatus
}

func (a *Dialer) GetAmiConnectionStatus() bool {
	return amigoClient.Connected()
}

// CreateDialer - create a new connection to Asterisk (singleton)
func CreateDialer(host, user, pass string) *Dialer {
	once.Do(func() {
		amiInstance = &Dialer{
			host: host,
			user: user,
			pass: pass,
		}

		// call to all available users while started
		go func() {
			for {
				if amiInstance.dialerStatus {
					amiInstance.callToAll()
				}
				time.Sleep(conf.General.CallTimeout * time.Second)
			}
		}()

	})

	return amiInstance
}

// GetDialer - get an existed connection
func GetDialer() *Dialer {
	return amiInstance
}
