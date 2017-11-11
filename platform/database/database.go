package database

import (
    "database/sql"
    "errors"
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
    Priority     MsisdnPriority `json:"priority" gorm:"ForeignKey:MsisdnID;AssociationForeignKey:ID;not null"`
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

    if err := tx.Raw("SELECT m.* FROM `msisdn_lists` m, `msisdn_priorities` p "+
        "WHERE m.id = p.msisdn_id and (m.status = ? or m.status = ?) ORDER BY p.priority, p.msisdn_id "+
        "LIMIT 1 FOR UPDATE",
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
    return list, d.getPreloadPriorityDB("", "").Find(list).Error
}

func (d *DB) GetMsisdnListInProgress(sortBy, sortOrder string) (int, *[]MsisdnList, error) {
    var count int
    d.Find(&[]MsisdnList{}).Count(&count)
    list, err := d.getMsisdnInProgressDB(sortBy, sortOrder)
    return count, list, err
}

func (d *DB) GetMsisdnListInProgressWithPagination(rows, page int, sortBy, sortOrder string) (int, *[]MsisdnList, error) {
    var count int
    d.Find(&[]MsisdnList{}).Count(&count)
    list, err := d.getMsisdnInProgressWithPaginationDB(rows, page, sortBy, sortOrder)
    return count, list, err
}

func (d *DB) getPreloadPriorityDB(sortByField, order string) *gorm.DB {
    return d.Preload("Priority")
}

func (d *DB) getMsisdnInProgressDB(sortByField, order string) (*[]MsisdnList, error) {
    if sortByField == "id" {
        sortByField = "l.id"
    }
    rows, err := d.Raw(fmt.Sprintf(GetQuery(MsisdnInProgress), sortByField, order)).Rows()
    if err != nil {
        xlog.Errorf("get in progress msisdns error: %s", err)
        return nil, err
    }
    return d.scanMsisdn(rows)
}

func (d *DB) getMsisdnInProgressWithPaginationDB(row, page int, sortByField, order string) (*[]MsisdnList, error) {
    if row == 0 {
        row = defaultMsisdnRowsCount
    }
    if page != 0 {
        page = (page - 1) * row
    }
    if sortByField == "id" {
        sortByField = "l.id"
    }
    rows, err := d.Raw(fmt.Sprintf(GetQuery(MsisdnInProgressWithPagination), sortByField, order, row, page)).Rows()
    if err != nil {
        xlog.Errorf("get in progress msisdns error: %s", err)
        return nil, err
    }
    return d.scanMsisdn(rows)
}

func (d *DB) scanMsisdn(rows *sql.Rows) (*[]MsisdnList, error){
    defer rows.Close()
    m := MsisdnList{Priority: MsisdnPriority{}}
    list := new([]MsisdnList)
    for rows.Next() {
        if err := rows.Scan(&m.ID, &m.Msisdn, &m.Status, &m.Time,
            &m.ActionID, &m.CauseTxt, &m.Cause, &m.Event,
            &m.Channel, &m.CallerIDNum, &m.CallerIDName,
            &m.Uniqueid, &m.TimeCalled,
            &m.Priority.ID, &m.Priority.MsisdnID, &m.Priority.Priority); err != nil {
            xlog.Errorf("scan error: %s", err)
            return nil, err
        }
        *list = append(*list, m)
    }
    return list, nil
}

func (d *DB) UpdatePriority(id, priority int) error {
    if priority > 10 && priority < 1 && id < 0 {
        return errors.New("Priority or ID error")
    }

    return d.Table("msisdn_priorities").Where("msisdn_id = ?", id).
        Updates(map[string]interface{}{"priority": priority}).Error
}

func (d *DB) DeleteMsisdn(id int) error {
    if id < 0 {
        return errors.New("ID can't be less then 0")
    }

    return d.Delete(&MsisdnList{}, "id = ?", id).Error
}

func (d *DB) AddNewNumbers(numbers []string) error {
    for _, number := range numbers {
        if err := d.Create(&MsisdnList{Msisdn: number, Priority: MsisdnPriority{}}).Error; err != nil {
            return err
        }
    }
    return nil
}

func GetDB() *DB {
    return db
}
