package database

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/rs/xlog"

	"github.com/incu6us/asterisk-dialer/internal/utils/config"
)

type DialerUser struct {
	Peer       string    `gorm:"type:varchar(10);unique;not null"`
	PeerStatus string    `gorm:"type:varchar(10);DEFAULT:'Unregister'"`
	Time       time.Time `sql:"type:TIMESTAMP;DEFAULT:CURRENT_TIMESTAMP"`
	Action     string    `gorm:"type:varchar(10);DEFAULT:''"`
	Exten      string    `gorm:"type:varchar(10);DEFAULT:''"`
}

type DialerMsisdnList struct {
	Id           int64     `sql:"primary_key;AUTO_INCREMENT;index"`
	Msisdn       string    `sql:"type:varchar(20);not null;index"`
	Status       string    `sql:"type:varchar(10);DEFAULT:'';index"`
	Time         time.Time `sql:"type:TIMESTAMP;DEFAULT:CURRENT_TIMESTAMP;index"`
	ActionID     string    `sql:"type:varchar(50);index"`
	CauseTxt     string    `sql:"type:varchar(50)"`
	Cause        string    `sql:"type:varchar(5);DEFAULT:''"`
	Event        string    `sql:"type:varchar(50)"`
	Channel      string    `sql:"type:varchar(50);index"`
	CallerIDNum  string    `sql:"type:varchar(20);index"`
	CallerIDName string    `sql:"type:varchar(20);index"`
	Uniqueid     string    `sql:"type:varchar(20);index"`
	TimeCalled   time.Time `sql:"type:TIMESTAMP;index"`
}

type DialerUsers []DialerUser

var dbInstance *gorm.DB

var userToId = make(map[string]int64)

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

		if err = dbInstance.AutoMigrate(&DialerUser{}, &DialerMsisdnList{}).Error; err != nil {
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

	if err = tx.Model(&DialerMsisdnList{}).Where("id = ?", msisdn.Id).Updates(map[interface{}]interface{}{
		"status":        "progress",
		"caller_id_num": callerIdNum,
		"time_called":   time.Now().UTC(),
	}).Error; err != nil {
		tx.Rollback()
		xlog.Errorf("Updating MSISDN error: %v", err)
		return ""
	}

	userToId[callerIdNum] = msisdn.Id

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
