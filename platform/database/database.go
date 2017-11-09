package database

import (
    "fmt"
    "strings"
    "sync"
    "time"

    "github.com/jinzhu/gorm"
    "github.com/rs/xlog"

    "github.com/incu6us/asterisk-dialer/utils/config"
)

const (
    defaultMsisdnRowsCount = 20
)

type DialerUser struct {
    Peer       string    `gorm:"type:varchar(10);unique;not null" json:"peer"`
    PeerStatus string    `gorm:"type:varchar(10);DEFAULT:'Unregister'" json:"peerStatus"`
    Time       time.Time `sql:"type:TIMESTAMP;DEFAULT:CURRENT_TIMESTAMP" json:"time"`
    Action     string    `gorm:"type:varchar(10);DEFAULT:''" json:"action"`
    Exten      string    `gorm:"type:varchar(10);DEFAULT:''" json:"exten"`
}

type MsisdnList struct {
    ID           int             `gorm:"primary_key" sql:"index" json:"id"`
    Msisdn       string          `sql:"type:varchar(20);not null;index" json:"msisdn"`
    Status       string          `sql:"type:varchar(10);DEFAULT:'';index" json:"status"`
    Time         time.Time       `sql:"type:TIMESTAMP;DEFAULT:CURRENT_TIMESTAMP;index" json:"time"`
    ActionID     string          `sql:"type:varchar(50);index" json:"actionId"`
    CauseTxt     string          `sql:"type:varchar(50)" json:"causeTxt"`
    Cause        string          `sql:"type:varchar(5);DEFAULT:''" json:"cause"`
    Event        string          `sql:"type:varchar(50)" json:"event"`
    Channel      string          `sql:"type:varchar(50);index" json:"channel"`
    CallerIDNum  string          `sql:"type:varchar(20);index" json:"callerIdNum"`
    CallerIDName string          `sql:"type:varchar(20);index" json:"callerIdName"`
    Uniqueid     string          `sql:"type:varchar(20);index" json:"uniqueId"`
    TimeCalled   time.Time       `sql:"type:TIMESTAMP;index" json:"timeCalled"`
    Priority     *MsisdnPriority `json:"priority" gorm:"ForeignKey:MsisdnID;AssociationForeignKey:ID;not null"`
}

type MsisdnPriority struct {
    ID       int  `json:"id" gorm:"primary_key" sql:"index"`
    MsisdnID int  `json:"msisdnId" sql:"index" gorm:"unique_index;not null"`
    Priority uint `json:"priority" sql:"DEFAULT:10"`
}

type DialerUsers []DialerUser

type DB struct {
    *gorm.DB
}

var db *DB
var once sync.Once

var userToId = make(map[string]int)

func (d *DB) Connect(tomlConfig *config.TomlConfig) (error) {
    var err error
    once.Do(func() {
        connString := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
            tomlConfig.DB.Username, tomlConfig.DB.Password, tomlConfig.DB.Host, tomlConfig.DB.Database)

        db = new(DB)
        db.DB, err = gorm.Open("mysql", connString)
        if err != nil {
            xlog.Errorf("error open database connection: %s", err)
        }

        db.LogMode(tomlConfig.DB.Debug)

        if err = db.AutoMigrate(&DialerUser{}, &MsisdnList{}, &MsisdnPriority{}).Error; err != nil {
            xlog.Errorf("err on Automigrate: %v", err)
        }
        if err = db.Model(&MsisdnPriority{}).AddForeignKey("msisdn_id", "msisdn_lists(id)", "CASCADE", "CASCADE").Error; err != nil {
            xlog.Errorf("error to create an ForeignKey: %s", err)
        }
    })
    return err
}

func (d *DB) GetAutodialUsers() DialerUsers {
    users := DialerUsers{}
    d.Model(&DialerUser{}).Find(&users)

    return users
}

func (d *DB) GetRegisteredUsers() DialerUsers {
    users := DialerUsers{}
    d.Model(&DialerUser{}).Where("peer_status = ? or peer_status = ?", "Reachable", "Registered").Find(&users)

    return users
}

func (d *DB) GetAutodialUser(user string) DialerUser {
    u := DialerUser{}
    //query.Table("dialer_users").Where("peer = ?", user).Find(&u)
    d.Model(&DialerUser{}).Where("peer = ?", user).Find(&u)

    return u
}

func (d *DB) DeleteMSISDNOlderThenWeek() {
    d.Delete(&MsisdnList{}, "time < DATE(NOW()) - INTERVAL 7 DAY")
}
func (d *DB) UpdatePeerStatus(user, status, action, exten string) {
    userFiledsUpdate := make(map[string]interface{})
    userFiledsUpdate["peer_status"] = status
    if action != "" {
        userFiledsUpdate["action"] = action
    }

    if exten != "" {
        userFiledsUpdate["exten"] = exten
    }

    userFiledsUpdate["time"] = time.Now().UTC()
    d.Model(&DialerUser{}).Where("peer = ?", user).Updates(userFiledsUpdate) //("peer_status", status)
}

func (d *DB) ProcessedMsisdn(callerIdNum string) string {
    tx := d.Begin()

    msisdn := &MsisdnList{}

    if err := tx.Raw(
        "SELECT * FROM `msisdn_lists` WHERE status = ? or status = ? LIMIT 1 FOR UPDATE",
        "", "recall").
        Scan(msisdn).
        Error; err != nil {
        tx.Rollback()
        xlog.Infof("Getting MSISDN error: %v", err)
        return ""
    }

    if err := tx.Model(&MsisdnList{}).Where("id = ?", msisdn.ID).Updates(map[interface{}]interface{}{
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

func (d *DB) UpdateAfterHangup(callerIDNum, callerIDName, cause, causeTxt, event, channel, uniqueid string) {
    var err error
    userToIdCopy := userToId[callerIDNum]

    tx := d.Begin()

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
        if err = tx.Model(&MsisdnList{}).Where(
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
        if err = tx.Model(&MsisdnList{}).Where(
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

func (d *DB) GetMsisdnListWithPriority() (*[]MsisdnList, error) {
    list := new([]MsisdnList)
    d.Preload("Priority").Find(list)

    return list, nil
}

func (d *DB) GetMsisdnListInProgress() (*[]MsisdnList, error) {
    list := new([]MsisdnList)
    err := d.getMsisdnInProgressDB(list).Error
    return list, err
}

func (d *DB) GetMsisdnCount() int {
    var count int
    d.getMsisdnInProgressDB(&[]MsisdnList{}).Count(&count)
    return count
}

func (d *DB) GetMsisdnListInProgressWithPagination(rows, page int) (*[]MsisdnList, error) {
    list := new([]MsisdnList)
    err := d.getMsisdnInProgressWithPaginationDB(list, rows, page).Error
    return list, err
}

func (d *DB) getPreloadPriorityDB() *gorm.DB {
    return d.Preload("Priority")
}

func (d *DB) getMsisdnInProgressDB(list *[]MsisdnList) *gorm.DB {
    return d.getPreloadPriorityDB().
        Find(list, "status = ? or status = ? or status = ?", "progress", "", "recall")
}

func (d *DB) getMsisdnInProgressWithPaginationDB(list *[]MsisdnList, row, page int) *gorm.DB {
    if row == 0 {
        row = defaultMsisdnRowsCount
    }
    if page != 0 {
        page = page - 1
    }
    return d.getPreloadPriorityDB().Limit(row).Offset(page).Find(list, "status = ? or status = ? or status = ?", "progress", "", "recall")
}

func (d *DB) AddNewNumbers(numbers []string) error {
    for _, number := range numbers {
        if err := d.Create(&MsisdnList{Msisdn: number, Priority: &MsisdnPriority{}}).Error; err != nil {
            return err
        }
    }
    return nil
}

func GetDB() *DB {
    return db
}
