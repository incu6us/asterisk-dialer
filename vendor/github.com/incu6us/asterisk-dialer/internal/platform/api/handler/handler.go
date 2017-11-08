package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/incu6us/asterisk-dialer/internal/platform/database"
	"github.com/rs/xlog"

	"github.com/incu6us/asterisk-dialer/internal/platform/dialer"
	"github.com/incu6us/asterisk-dialer/internal/utils/config"
	"github.com/incu6us/asterisk-dialer/internal/utils/dialer-helper"
	"github.com/incu6us/asterisk-dialer/internal/utils/file-utils"
)

type ApiHandler struct {
	ContentType string
	//Scheduler   scheduler.SchedulerHelper
}

type response struct {
	Result interface{} `json:"Result"`
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
	w.Header().Set("Content-Type", a.ContentType)

	if encodeError := json.NewEncoder(w).Encode(response{message}); encodeError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		xlog.Error("Parse message error", encodeError)
	}
}

func (a *ApiHandler) Test(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	//xlog.Infof("MSISDN: %s", database.ProcesseMsisdn())
	a.print(w, r, vars)
}

func (a *ApiHandler) Start(w http.ResponseWriter, r *http.Request) {
	a.print(w, r, dialer.GetDialer().StartDialing())
}

func (a *ApiHandler) Stop(w http.ResponseWriter, r *http.Request) {
	a.print(w, r, dialer.GetDialer().StopDialing())
}

func (a *ApiHandler) DialerStatus(w http.ResponseWriter, r *http.Request) {
	a.print(w, r, dialer.GetDialer().GetDialerStatus())
}

func (a *ApiHandler) GetRegisteredUsers(w http.ResponseWriter, r *http.Request) {
	a.print(w, r, database.GetRegisteredUsers())
}

func (a *ApiHandler) UploadMSISDNList(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	msisdnList, err := ioutil.ReadAll(r.Body)
	if err != nil {
		xlog.Errorf("Can't load msisdn list: $s", err)
		return
	}

	numberList := file_utils.MSISDNNormalizer(string(msisdnList))

	if err = dialer_helper.AddAndRun(numberList); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		xlog.Error(err)
		a.print(w, r, err)
		return
	}

	a.print(w, r, "Numbers added")
}

// simple check which improve, that server is running
func (a *ApiHandler) Ready(w http.ResponseWriter, r *http.Request) {
	a.print(w, r, "Service is up and running")
}

func GetHandler() *ApiHandler {
	once.Do(func() {
		handler = &ApiHandler{ContentType: CONTENT_TYPE}
	})

	return handler
}
