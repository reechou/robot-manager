package models

import (
	"github.com/reechou/holmes"
)

type Robot struct {
	ID            int64  `xorm:"pk autoincr"`
	RobotWx       string `xorm:"not null default '' varchar(128) unique"`
	RobotType     int    `xorm:"not null default 0 int index"` // 0: just robot 1: robot group manager 2: robot wechat business
	IfSaveFriend  int64  `xorm:"not null default 0 int"`
	IfSaveGroup   int64  `xorm:"not null default 0 int"`
	Ip            string `xorm:"not null default '' varchar(64) index"`
	OfPort        string `xorm:"not null default '' varchar(64) index"`
	LastLoginTime int64  `xorm:"not null default 0 int"`
	Argv          string `xorm:"not null default '' varchar(2048)"`
	BaseLoginInfo string `xorm:"not null default '' varchar(2048)"`
	WebwxCookie   string `xorm:"not null default '' varchar(2048)"`
	CreatedAt     int64  `xorm:"not null default 0 int"`
	UpdatedAt     int64  `xorm:"not null default 0 int"`
}

func GetRobot(info *Robot) (bool, error) {
	has, err := x.Where("robot_wx = ?", info.RobotWx).Get(info)
	if err != nil {
		return false, err
	}
	if !has {
		holmes.Debug("cannot find robot from robot_wx[%s]", info.RobotWx)
		return false, nil
	}
	return true, nil
}

func GetRobotList() ([]Robot, error) {
	var list []Robot
	err := x.Find(&list)
	if err != nil {
		holmes.Error("get robot list error: %v", err)
		return nil, err
	}
	return list, nil
}
