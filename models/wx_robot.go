package models

import (
	"github.com/reechou/holmes"
)

type Robot struct {
	ID            int64  `xorm:"pk autoincr" json:"id"`
	RobotWx       string `xorm:"not null default '' varchar(128) unique" json:"robotWx"`
	RobotType     int    `xorm:"not null default 0 int index"` // 0: just robot 1: robot group manager 2: robot wechat business
	IfSaveFriend  int64  `xorm:"not null default 0 int" json:"-"`
	IfSaveGroup   int64  `xorm:"not null default 0 int" json:"-"`
	Ip            string `xorm:"not null default '' varchar(64) index" json:"-"`
	OfPort        string `xorm:"not null default '' varchar(64) index" json:"-"`
	LastLoginTime int64  `xorm:"not null default 0 int" json:"lastLoginTime"`
	Argv          string `xorm:"not null default '' varchar(2048)" json:"-"`
	BaseLoginInfo string `xorm:"not null default '' varchar(2048)" json:"-"`
	WebwxCookie   string `xorm:"not null default '' varchar(2048)" json:"-"`
	CreatedAt     int64  `xorm:"not null default 0 int" json:"createAt"`
	UpdatedAt     int64  `xorm:"not null default 0 int" json:"-"`
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
