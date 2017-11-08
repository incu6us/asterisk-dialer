package database

import (
    "fmt"
    "log"
    "strings"
    "time"

    "github.com/jinzhu/gorm"
    "github.com/rs/xlog"

    "github.com/incu6us/asterisk-dialer/utils/config"
)

type DialerUser struct {
    Peer       string    `gorm:"type:varchar(10);unique;not null" json:"peer"`
    PeerStatus string    `gorm:"type:varchar(10);DEFAULT:'Unregister'" json:"peerStatus"`
    Time       time.Time `sql:"type:TIMESTAMP;DEFAULT:CURRENT_TIMESTAMP" json:"time"`
    Action     string    `gorm:"type:varchar(10);DEFAULT:''" json:"action"`
    Exten      string    `gorm:"type:varchar(10);DEFAULT:''" json:"exten"`
}

type DialerMsisdnList struct {
    ID           int                  `gorm:"primary_key" sql:"index" json:"id"`
    Msisdn       string               `sql:"type:varchar(20);not null;index" json:"msisdn"`
    Status       string               `sql:"type:varchar(10);DEFAULT:'';index" json:"status"`
    Time         time.Time            `sql:"type:TIMESTAMP;DEFAULT:CURRENT_TIMESTAMP;index" json:"time"`
    ActionID     string               `sql:"type:varchar(50);index" json:"actionId"`
    CauseTxt     string               `sql:"type:varchar(50)" json:"causeTxt"`
    Cause        string               `sql:"type:varchar(5);DEFAULT:''" json:"cause"`
    Event        string               `sql:"type:varchar(50)" json:"event"`
    Channel      string               `sql:"type:varchar(50);index" json:"channel"`
    CallerIDNum  string               `sql:"type:varchar(20);index" json:"callerIdNum"`
    CallerIDName string               `sql:"type:varchar(20);index" json:"callerIdName"`
    Uniqueid     string               `sql:"type:varchar(20);index" json:"uniqueId"`
    TimeCalled   time.Time            `sql:"type:TIMESTAMP;index" json:"timeCalled"`
    Priority     DialerMsisdnPriority `json:"offer" gorm:"ForeignKey:MsisdnID;AssociationForeignKey:ID;not null"`
}

type DialerMsisdnPriority struct {
    ID       int  `json:"id" gorm:"primary_key" sql:"index"`
    MsisdnID int  `sql:"index" gorm:"unique_index;not null"`
    Priority uint `sql:"DEFAULT:10"`
}

type DialerUsers []DialerUser

var dbInstance *gorm.DB

var userToId = make(map[string]int)

func Connect(tomlConfig *config.TomlConfig) (*gorm.DB, error) {
    if dbInstance == nil {
        connString := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
            tomlConfig.DB.Username, tomlConfig.DB.Password, tomlConfig.DB.Host, tomlConfig.DB.Database)

        database, err := gorm.Open("mysql", connString)
        if err != nil {
            return nil, err
        }

        database.LogMode(tomlConfig.DB.Debug)
        dbInstance = database

        if err = dbInstance.AutoMigrate(&DialerUser{}, &DialerMsisdnList{}, &DialerMsisdnPriority{}).Error; err != nil {
            xlog.Errorf("err on Automigrate: %v", err)
        }
    }

    return dbInstance, nil
}

func GetAutodialUsers() DialerUsers {
    query, err := Connect(config.GetConfig())
    if err != nil {
        xlog.Fatal(err)
    }

    users := DialerUsers{}
    //query.Table("dialer_users").Find(&users)
    query.Model(&DialerUser{}).Find(&users)

    return users
}

func GetRegisteredUsers() DialerUsers {
    query, err := Connect(config.GetConfig())
    if err != nil {
        xlog.Fatal(err)
    }

    users := DialerUsers{}
    //query.Table("dialer_users").Find(&users)
    query.Model(&DialerUser{}).Where("peer_status = ? or peer_status = ?", "Reachable", "Registered").Find(&users)

    return users
}

func GetAutodialUser(user string) DialerUser {
    query, err := Connect(config.GetConfig())
    if err != nil {
        log.Fatal(err)
    }

    u := DialerUser{}
    //query.Table("dialer_users").Where("peer = ?", user).Find(&u)
    query.Model(&DialerUser{}).Where("peer = ?", user).Find(&u)

    return u
}

func DeleteMSISDNOlderThenWeek() {
    query, err := Connect(config.GetConfig())
    if err != nil {
        log.Fatalf("error to delete racords older then 7 days: %s", err)
    }

    query.Delete(&DialerMsisdnList{}, "time < DATE(NOW()) - INTERVAL 7 DAY")
}
func UpdatePeerStatus(user, status, action, exten string) {
    query, err := Connect(config.GetConfig())
    if err != nil {
        xlog.Fatal(err)
    }

    userFiledsUpdate := make(map[string]interface{})
    userFiledsUpdate["peer_status"] = status
    if action != "" {
        userFiledsUpdate["action"] = action
    }

    if exten != "" {
        userFiledsUpdate["exten"] = exten
    }

    userFiledsUpdate["time"] = time.Now().UTC()

    //query.Table("dialer_users").Where("peer = ?", user).Updates(userFiledsUpdate) //("peer_status", status)
    query.Model(&DialerUser{}).Where("peer = ?", user).Updates(userFiledsUpdate) //("peer_status", status)

}

func ProcesseMsisdn(callerIdNum string) string {
    query, err := Connect(config.GetConfig())
    if err != nil {
        xlog.Fatal(err)
    }

    tx := query.Begin()

    msisdn := &DialerMsisdnList{}

    if err := tx.Raw(
        "SELECT * FROM `dialer_msisdn_lists` WHERE status = ? or status = ? LIMIT 1 FOR UPDATE",
        "", "recall").
        Scan(msisdn).
        Error; err != nil {
        tx.Rollback()
        xlog.Infof("Getting MSISDN error: %v", err)
        return ""
    }

    if err = tx.Model(&DialerMsisdnList{}).Where("id = ?", msisdn.ID).Updates(map[interface{}]interface{}{
        "status":        "progress",
        "caller_id_num": callerIdNum,
        "time_called":   time.Now().UTC(),
    }).Error; err != nil {
        tx.Rollback()
        xlog.Errorf("Updating MSISDN error: %v", err)
        return ""
    }

    userToId[callerIdNum] = msisdn.ID

    tx.Commit()

    return msisdn.Msisdn
}

func UpdateAfterHangup(callerIDNum, callerIDName, cause, causeTxt, event, channel, uniqueid string) {
    userToIdCopy := userToId[callerIDNum]
    query, err := Connect(config.GetConfig())
    if err != nil {
        xlog.Fatal(err)
    }

    tx := query.Begin()

    xlog.Debugf(
        "UpdateAfterHangup: callerIDNum: %s, callerIDName: %s, cause: %s, causeTxt: %s, event: %s, channel: %s, uniqueid: %s",
        callerIDNum,
        callerIDName,
        cause,
        causeTxt,
        event,
        channel,
        uniqueid,
    )

    if cause == "16" {
        if err = tx.Model(&DialerMsisdnList{}).Where(
            "id = ? and caller_id_num = ? and msisdn = ? and (status = ? or status = ?)",
            userToIdCopy,
            callerIDNum,
            strings.TrimPrefix(callerIDName, "autodial_"),
            "progress",
            "recall",
        ).Updates(map[interface{}]interface{}{
            "status":      "done",
            "time_called": time.Now().UTC(),
            "cause_txt":   causeTxt,
            "event":       event,
            "channel":     channel,
            "uniqueid":    uniqueid,
        }).Error; err != nil {
            xlog.Errorf("Updating MSISDN error: %v", err)
            tx.Rollback()
        }
    } else {
        if err = tx.Model(&DialerMsisdnList{}).Where(
            "id = ? and caller_id_num = ? and msisdn = ? and (status = ? or status = ?)",
            userToIdCopy,
            callerIDNum,
            strings.TrimPrefix(callerIDName, "autodial_"),
            "progress",
            "recall",
        ).Updates(map[interface{}]interface{}{
            "status":      "recall",
            "time_called": time.Now().UTC(),
            "cause_txt":   causeTxt,
            "event":       event,
            "channel":     channel,
            "uniqueid":    uniqueid,
        }).Error; err != nil {
            xlog.Errorf("Updating MSISDN error: %v", err)
            tx.Rollback()
        }
    }

    tx.Commit()
}

func AddNewNumbers(numbers []string) error {

    query, err := Connect(config.GetConfig())
    if err != nil {
        xlog.Fatal(err)
        return err
    }

    for _, number := range numbers {
        if err := query.Create(&DialerMsisdnList{Msisdn: number}).Error; err != nil {
            return err
        }
    }

    return nil
}
