package api

import (
    "net/http"
    "time"

    "github.com/gorilla/mux"
    "github.com/incu6us/asterisk-dialer/platform/database"

    "github.com/incu6us/asterisk-dialer/platform/api/handler"
)

const (
    dialerPathPrefix  = "/api/v1/dialer"
    generalPathPrefix = "/api/v1"
    httpTimeout       = 120 * time.Second
)

type API struct {
    Name        string
    Method      string
    Path        string
    Pattern     string
    HandlerFunc http.HandlerFunc
}

type APIs []API

func NewHandler(db *database.DB) http.Handler {

    handler.InitHandlers(db)
    router := mux.NewRouter().StrictSlash(true)

    for _, api := range apis {
        if api.Name != "ready" {
            router.
                PathPrefix(dialerPathPrefix).
                Methods(api.Method).
                Path(api.Path+api.Pattern).
                Name(api.Name).
                Handler(api.HandlerFunc)
        } else {
            router.
                PathPrefix(generalPathPrefix).
                Methods(api.Method).
                Path(api.Path+api.Pattern).
                Name(api.Name).
                Handler(api.HandlerFunc)
        }
    }

    // static files
    router.PathPrefix("/").Handler(http.FileServer(http.Dir("./ui/build/")))

    middlewareHandler := http.TimeoutHandler(router, httpTimeout, "Server timeout!")

    return middlewareHandler
}

var apis = APIs{
    API{
        "start",
        "GET",
        "",
        "/start",
        handler.GetHandler().Start,
    },
    API{
        "stop",
        "GET",
        "",
        "/stop",
        handler.GetHandler().Stop,
    },
    API{
        "status",
        "GET",
        "",
        "/status",
        handler.GetHandler().DialerStatus,
    },
    API{
        "ready",
        "GET",
        "",
        "/ready",
        handler.GetHandler().Ready,
    },
    API{
        "registeredUsers",
        "GET",
        "",
        "/registeredUsers",
        handler.GetHandler().GetRegisteredUsers,
    },
    API{
        "upload-msisdn",
        "POST",
        "",
        "/upload-msisdn",
        handler.GetHandler().UploadMSISDNList,
    },
    API{
        "msisdnList",
        "GET",
        "/msisdn",
        "/all",
        handler.GetHandler().GetMsisdnList,
    },
    API{
        "msisdnListInProgress",
        "GET",
        "/msisdn",
        "/inProgress", // ?page=1&limit=20
        handler.GetHandler().GetMsisdnListInProgress,
    },
    API{
        "msisdnListInProgressUpdatePriority",
        "PUT",
        "/msisdn",
        "/inProgress/{id}",
        handler.GetHandler().GetMsisdnListInProgressUpdatePriority,
    },
    API{
        "deleteMsisdn",
        "DELETE",
        "/msisdn",
        "/inProgress/{id}",
        handler.GetHandler().DeleteMsisdn,
    },
    API{
        "clearAllMsisdn",
        "DELETE",
        "/msisdn",
        "/clearAll",
        handler.GetHandler().ClearAll,
    },
    API{
        "test",
        "GET",
        "",
        "/test",
        handler.GetHandler().Test,
    },
}
