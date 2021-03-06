package handler

import (
    "encoding/json"
    "errors"
    "io/ioutil"
    "net/http"
    "strconv"
    "sync"

    "github.com/gorilla/mux"
    "github.com/rs/xlog"

    "github.com/incu6us/asterisk-dialer/platform/database"
    "github.com/incu6us/asterisk-dialer/platform/dialer"
    "github.com/incu6us/asterisk-dialer/utils/config"
    "github.com/incu6us/asterisk-dialer/utils/file-utils"
)

type ApiHandler struct {
    DefaultContentType string
    db                 *database.DB
}

type response struct {
    Result interface{} `json:"result"`
}

type msisdnPaging struct {
    Total  int                 `json:"total"`
    Result *[]normalizedMsisdn `json:"result"`
}

type normalizedMsisdn struct {
    ID           int                     `json:"id"`
    Msisdn       string                  `json:"msisdn"`
    Status       string                  `sql:"type:varchar(10);DEFAULT:'';index" json:"status"`
    Time         string                  `json:"time"`
    ActionID     string                  `sql:"type:varchar(50);index" json:"actionId"`
    CauseTxt     string                  `sql:"type:varchar(50)" json:"causeTxt"`
    Cause        string                  `sql:"type:varchar(5);DEFAULT:''" json:"cause"`
    Event        string                  `sql:"type:varchar(50)" json:"event"`
    Channel      string                  `sql:"type:varchar(50);index" json:"channel"`
    CallerIDNum  string                  `sql:"type:varchar(20);index" json:"callerIdNum"`
    CallerIDName string                  `sql:"type:varchar(20);index" json:"callerIdName"`
    Uniqueid     string                  `sql:"type:varchar(20);index" json:"uniqueId"`
    TimeCalled   string                  `json:"timeCalled"`
    Priority     database.MsisdnPriority `json:"priority"`
}

const (
    CONTENT_TYPE = "application/json"
    TimeFormat   = "2006-01-02 15:04:05"
)

var (
    once    sync.Once
    handler ApiHandler
    conf    = config.GetConfig()
)

func (a *ApiHandler) print(w http.ResponseWriter, r *http.Request, message interface{}) {
    if err := json.NewEncoder(w).Encode(response{message}); err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        xlog.Errorf("parse message error: %s", err)
    }
}

func (a *ApiHandler) convertMsisdnObject(tmpList *[]database.MsisdnList) *[]normalizedMsisdn {
    list := make([]normalizedMsisdn, 0, len(*tmpList))
    for _, i := range *tmpList {
        var timeCalled string
        if i.TimeCalled.Format(TimeFormat) == "0001-01-01 00:00:00" {
            timeCalled = "0000-00-00 00:00:00"
        } else {
            timeCalled = i.TimeCalled.Format(TimeFormat)
        }
        item := &normalizedMsisdn{
            ID:           i.ID,
            Msisdn:       i.Msisdn,
            Status:       i.Status,
            Time:         i.Time.Format(TimeFormat),
            ActionID:     i.ActionID,
            CauseTxt:     i.CauseTxt,
            Cause:        i.Cause,
            Event:        i.Event,
            Channel:      i.Channel,
            CallerIDNum:  i.CallerIDNum,
            CallerIDName: i.CallerIDName,
            Uniqueid:     i.Uniqueid,
            TimeCalled:   timeCalled,
            Priority:     i.Priority,
        }
        list = append(list, *item)
    }
    return &list
}

func (a *ApiHandler) Test(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)

    //xlog.Infof("MSISDN: %s", database.ProcesseMsisdn())
    w.Header().Set("Content-Type", a.DefaultContentType)
    w.WriteHeader(http.StatusOK)
    a.print(w, r, vars)
}

func (a *ApiHandler) Start(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", a.DefaultContentType)
    w.WriteHeader(http.StatusOK)
    a.print(w, r, dialer.GetDialer().StartDialing())
}

func (a *ApiHandler) Stop(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", a.DefaultContentType)
    w.WriteHeader(http.StatusOK)
    a.print(w, r, dialer.GetDialer().StopDialing())
}

func (a *ApiHandler) DialerStatus(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", a.DefaultContentType)
    w.WriteHeader(http.StatusOK)
    a.print(w, r, dialer.GetDialer().GetDialerStatus())
}

func (a *ApiHandler) GetRegisteredUsers(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", a.DefaultContentType)
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

    w.Header().Set("Content-Type", a.DefaultContentType)
    w.WriteHeader(http.StatusCreated)
    a.print(w, r, "Numbers added")
}

func (a *ApiHandler) GetMsisdnList(w http.ResponseWriter, r *http.Request) {
    tmpList, err := a.db.GetMsisdnListWithPriority()
    if err != nil {
        xlog.Errorf("Can't load priority msisdn list: $s", err)
        w.WriteHeader(http.StatusInternalServerError)
        a.print(w, r, err)
        return
    }
    w.Header().Set("Content-Type", a.DefaultContentType)
    w.WriteHeader(http.StatusOK)
    a.print(w, r, a.convertMsisdnObject(tmpList))
}

// defaults: page=20; if limit=0 - show all records
func (a *ApiHandler) GetMsisdnListInProgress(w http.ResponseWriter, r *http.Request) {
    var page, limit, count int
    var sortOrder, sortBy string
    var list *[]database.MsisdnList
    var err error
    vars := r.URL.Query() // []string
    if varPage, ok := vars["page"]; ok {
        if len(varPage) > 0 {
            if page, err = strconv.Atoi(varPage[0]); err != nil {
                xlog.Errorf("can't parse 'page' param from url: %s", err)
                w.WriteHeader(http.StatusInternalServerError)
                a.print(w, r, err)
                return
            }
        }
    }
    if varLimit, ok := vars["limit"]; ok {
        if len(varLimit) > 0 {
            if limit, err = strconv.Atoi(varLimit[0]); err != nil {
                xlog.Errorf("can't parse 'limit' param from url: %s", err)
                w.WriteHeader(http.StatusInternalServerError)
                a.print(w, r, err)
                return
            }
        }
    }

    // return: sortBy, sortOrder
    sortFn := func(values map[string][]string) (string, string) {
        if sortOrderVar, ok := vars["sortOrder"]; ok {
            if len(sortOrderVar) > 0 {
                sortOrder = sortOrderVar[0]
            }
        }
        if sortByVar, ok := vars["sortBy"]; ok {
            if len(sortByVar) > 0 {
                sortBy = sortByVar[0]
            }
        }
        if sortBy == "" {
            sortBy = "id"
        }
        if sortOrder == "" {
            sortOrder = "asc"
        }
        return sortBy, sortOrder
    }

    sortBy, sortOrder = sortFn(vars)
    if limit == 0 {
        count, list, err = a.db.GetMsisdnListInProgress(sortBy, sortOrder)
    } else {
        count, list, err = a.db.GetMsisdnListInProgressWithPagination(limit, page, sortBy, sortOrder)
    }
    if err != nil {
        xlog.Errorf("get msisdn list with priority error: %s", err)
        w.WriteHeader(http.StatusInternalServerError)
        a.print(w, r, err)
        return
    }

    w.Header().Set("Content-Type", a.DefaultContentType)
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(
        msisdnPaging{
            Total:  count,
            Result: a.convertMsisdnObject(list),
        },
    ); err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        xlog.Errorf("parse message error: %s", err)
        return
    }
}

type MsisdnUpdatePriority struct {
    Priority int `json:"priority"`
}

func (a *ApiHandler) GetMsisdnListInProgressUpdatePriority(w http.ResponseWriter, r *http.Request) {
    idVar, ok := mux.Vars(r)["id"]
    if !ok {
        w.WriteHeader(http.StatusBadRequest)
        a.print(w, r, response{errors.New("id is not set")})
        return
    }
    xlog.Infof("ID: %#v", idVar)
    id, err := strconv.Atoi(idVar)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        xlog.Errorf("can't read id var: %s", err)
        a.print(w, r, response{err})
        return
    }
    prior := new(MsisdnUpdatePriority)
    data, err := ioutil.ReadAll(r.Body)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        xlog.Errorf("can't read data from body: %s", err)
        a.print(w, r, response{err})
        return
    }
    xlog.Infof("BODY: %#v", string(data))
    err = json.Unmarshal(data, prior)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        xlog.Errorf("can't serialize data from bytes: %s", err)
        a.print(w, r, response{err})
        return
    }

    if err = a.db.UpdatePriority(id, prior.Priority); err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        xlog.Errorf("can't update priority: %s", err)
        a.print(w, r, response{err})
        return
    }
    w.WriteHeader(http.StatusNoContent)
}

func (a *ApiHandler) DeleteMsisdn(w http.ResponseWriter, r *http.Request) {
    idVar, ok := mux.Vars(r)["id"]
    if !ok {
        w.WriteHeader(http.StatusBadRequest)
        a.print(w, r, response{errors.New("id is not set")})
        return
    }
    id, err := strconv.Atoi(idVar)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        xlog.Errorf("can't read id var: %s", err)
        a.print(w, r, response{err})
        return
    }
    if err := a.db.DeleteMsisdn(id); err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        xlog.Errorf("can't delete item: %s", err)
        a.print(w, r, response{err})
        return
    }
    w.WriteHeader(http.StatusNoContent)
}

func (a *ApiHandler) ClearAll(w http.ResponseWriter, r *http.Request){
    if err := a.db.ClearAll(); err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        xlog.Errorf("can't delete items: %s", err)
        a.print(w, r, response{err})
    }
    w.WriteHeader(http.StatusNoContent)
}

// simple check which improve, that server is running
func (a *ApiHandler) Ready(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", a.DefaultContentType)
    w.WriteHeader(http.StatusOK)
    a.print(w, r, "Service is up and running")
}

func InitHandlers(db *database.DB) {
    once.Do(func() {
        xlog.Debugf("DB -> %#v", db)
        handler = ApiHandler{DefaultContentType: CONTENT_TYPE, db: db}
    })
}

func GetHandler() *ApiHandler {
    return &handler
}
