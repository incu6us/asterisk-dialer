package handler

import (
    "encoding/json"
    "io/ioutil"
    "net/http"
    "sync"

    "github.com/gorilla/mux"
    "github.com/rs/xlog"

    "github.com/incu6us/asterisk-dialer/platform/database"
    "github.com/incu6us/asterisk-dialer/platform/dialer"
    "github.com/incu6us/asterisk-dialer/utils/config"
    "github.com/incu6us/asterisk-dialer/utils/file-utils"
)

type ApiHandler struct {
    ContentType string
    db          *database.DB
}

type response struct {
    Result interface{} `json:"result"`
}

const (
    CONTENT_TYPE = "application/json"
)

var (
    once    sync.Once
    handler *ApiHandler
    conf    = config.GetConfig()
)

func (a *ApiHandler) print(w http.ResponseWriter, r *http.Request, message interface{}) {
    if err := json.NewEncoder(w).Encode(response{message}); err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        xlog.Errorf("parse message error: %s", err)
    }
}

func (a *ApiHandler) Test(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)

    //xlog.Infof("MSISDN: %s", database.ProcesseMsisdn())
    w.Header().Set("Content-Type", a.ContentType)
    w.WriteHeader(http.StatusOK)
    a.print(w, r, vars)
}

func (a *ApiHandler) Start(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", a.ContentType)
    w.WriteHeader(http.StatusOK)
    a.print(w, r, dialer.GetDialer().StartDialing())
}

func (a *ApiHandler) Stop(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", a.ContentType)
    w.WriteHeader(http.StatusOK)
    a.print(w, r, dialer.GetDialer().StopDialing())
}

func (a *ApiHandler) DialerStatus(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", a.ContentType)
    w.WriteHeader(http.StatusOK)
    a.print(w, r, dialer.GetDialer().GetDialerStatus())
}

func (a *ApiHandler) GetRegisteredUsers(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", a.ContentType)
    w.WriteHeader(http.StatusOK)
    a.print(w, r, a.db.GetRegisteredUsers())
}

func (a *ApiHandler) UploadMSISDNList(w http.ResponseWriter, r *http.Request) {
    msisdnList, err := ioutil.ReadAll(r.Body)
    if err != nil {
        xlog.Errorf("Can't load msisdn list: $s", err)
        w.WriteHeader(http.StatusInternalServerError)
        a.print(w, r, err)
        return
    }

    numberList := file_utils.MSISDNNormalizer(string(msisdnList))

    if err = a.db.AddNewNumbers(numberList); err != nil {
        xlog.Error(err)
        w.WriteHeader(http.StatusInternalServerError)
        a.print(w, r, err)
        return
    }

    w.Header().Set("Content-Type", a.ContentType)
    w.WriteHeader(http.StatusCreated)
    a.print(w, r, "Numbers added")
}

func (a *ApiHandler) GetMsisdnListWithPriority(w http.ResponseWriter, r *http.Request) {
    list, err := a.db.GetMsisdnListWithPriority()
    if err != nil {
        xlog.Errorf("get msisdn list with priority error: %s", err)
        w.WriteHeader(http.StatusInternalServerError)
        a.print(w, r, err)
        return
    }

    w.Header().Set("Content-Type", a.ContentType)
    w.WriteHeader(http.StatusOK)
    a.print(w, r, list)
}

// simple check which improve, that server is running
func (a *ApiHandler) Ready(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", a.ContentType)
    w.WriteHeader(http.StatusOK)
    a.print(w, r, "Service is up and running")
}

func GetHandler() *ApiHandler {
    once.Do(func() {
        handler = &ApiHandler{ContentType: CONTENT_TYPE, db: database.GetDB()}
    })

    return handler
}
