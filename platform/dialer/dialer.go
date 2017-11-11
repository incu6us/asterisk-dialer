package dialer

import (
    "errors"
    "fmt"
    "strings"
    "sync"
    "time"

    "github.com/ivahaev/amigo"
    "github.com/ivahaev/amigo/uuid"
    "github.com/rs/xlog"

    "github.com/incu6us/asterisk-dialer/platform/database"
    "github.com/incu6us/asterisk-dialer/utils/config"
)

type dialer struct {
    dialerStatus        bool
    amiConnectionStatus bool
    randGenDigit        int
    host                string
    user                string
    pass                string
    db                  *database.DB
}

var (
    once        sync.Once
    amiInstance *dialer
    amigoClient *amigo.Amigo
    conf        = config.GetConfig()
)

var (
    userIsCalling map[string]bool
)

func (a *dialer) peerDial(message map[string]string) {
    user := a.db.GetAutodialUser(strings.TrimPrefix(message["Peer"], "SIP/"))
    if user.Peer != "" &&
        !userIsCalling[strings.TrimPrefix(message["Peer"], "SIP/")] &&
        a.dialerStatus {
        userIsCalling[strings.TrimPrefix(message["Peer"], "SIP/")] = true

        xlog.Infof("Dial: %#v", message)
    }
}

func (a *dialer) peerHangup(message map[string]string) {
    user := a.db.GetAutodialUser(message["CallerIDNum"])
    if user.Peer != "" {
        callerIDNum := message["CallerIDNum"]
        callerIDName := message["CallerIDName"]
        cause := message["Cause"]
        causeTxt := message["Cause-txt"]
        event := message["Event"]
        channel := message["Channel"]
        uniqueid := message["Uniqueid"]

        xlog.Debugf("peer hangup: %#v", message)

        a.db.UpdateAfterHangup(callerIDNum, callerIDName, cause, causeTxt, event, channel, uniqueid)
        if userIsCalling[message["CallerIDNum"]] {
            a.db.UpdatePeerStatus(user.Peer, "", event, "")
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

func (a *dialer) peerStatusListening(message map[string]string) {
    user := a.db.GetAutodialUser(strings.TrimPrefix(message["Peer"], "SIP/"))
    if user.Peer != "" {

        xlog.Infof("Peer status changed -> %v", message)
        a.db.UpdatePeerStatus(user.Peer, message["PeerStatus"], "", "")

        switch message["PeerStatus"] {
        case "Reachable", "Registered":
            time.Sleep(5 * time.Second)
            if !a.dialerStatus {
                return
            }
            _, err := a.CallToUser(user.Peer)
            if err != nil {
                xlog.Warnf("CallToUser problem: %s", err)
            }
        }
    }
}

func (a *dialer) Connect() error {
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
        xlog.Errorf("Connection error to %s: %s", a.host, message)
        err = errors.New(message)
    })

    // TODO: Dial register must be added with adding event to 'userIsCalling'
    amigoClient.RegisterHandler("PeerStatus", a.peerStatusListening)
    amigoClient.RegisterHandler("Hangup", a.peerHangup)

    return err
}

func (a *dialer) callToAll() {
    users := a.db.GetRegisteredUsers()
    for _, user := range users {
        xlog.Debugf("Calling to %s", user.Peer)
        a.CallToUser(user.Peer)
    }
}

func (a *dialer) originate(params map[string]string) (map[string]string, error) {
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

func (a *dialer) CallToUser(userID string) (map[string]string, error) {
    // if user already calling by autodialer, then skip the new call request
    if userIsCalling[userID] {
        return nil, errors.New("already calling by dialer")
    }

    if msisdn := a.db.ProcessedMsisdn(userID); msisdn != "" {
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

        a.db.UpdatePeerStatus(userID, "", "Originate", msisdn)

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

func (a *dialer) StartDialing() string {

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

func (a *dialer) StopDialing() string {

    if a.dialerStatus {

        amigoClient.UnregisterHandler("Dial", a.peerDial)

        a.dialerStatus = false
    }

    xlog.Debug("Dialer has been stopped")
    return "stopped"
}

func (a *dialer) GetDialerStatus() bool {
    return a.dialerStatus
}

func (a *dialer) GetAmiConnectionStatus() bool {
    return amigoClient.Connected()
}

// CreateDialer - create a new connection to Asterisk (singleton)
func CreateDialer(cfg *config.TomlConfig, db *database.DB) *dialer {
    once.Do(func() {
        amiInstance = &dialer{
            host: fmt.Sprintf("%s:%d", cfg.Ami.Host, cfg.Ami.Port),
            user: cfg.Ami.Username,
            pass: cfg.Ami.Password,
            db:   db,
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
func GetDialer() *dialer {
    return amiInstance
}
