package models

import (
	"fmt"
	"time"

	"github.com/reechou/holmes"
)

type RobotBlacklist struct {
	ID        int64  `xorm:"pk autoincr"`
	RobotWx   string `xorm:"not null default '' varchar(128) index"`
	UserName  string `xorm:"not null default '' varchar(256) index"`
	NickName  string `xorm:"not null default '' varchar(256)"`
	CreatedAt int64  `xorm:"not null default 0 int"`
	UpdatedAt int64  `xorm:"not null default 0 int"`
}

func CreateRobotBlacklist(info *RobotBlacklist) error {
	if info.RobotWx == "" {
		return fmt.Errorf("wx robot wx[%s] cannot be nil.", info.RobotWx)
	}

	now := time.Now().Unix()
	info.CreatedAt = now
	info.UpdatedAt = now

	_, err := x.Insert(info)
	if err != nil {
		holmes.Error("create robot blacklist error: %v", err)
		return err
	}
	holmes.Info("create robot blacklist[%v] success.", info)

	return nil
}

func GetRobotBlacklistFromUsername(info *RobotBlacklist) (bool, error) {
	has, err := x.Where("robot_wx = ?", info.RobotWx).And("user_name = ?", info.UserName).Get(info)
	if err != nil {
		return false, err
	}
	if !has {
		holmes.Debug("cannot find robot blacklist from robot[%s] username[%s]", info.RobotWx, info.UserName)
		return false, nil
	}
	return true, nil
}

type RobotManagerUser struct {
	ID        int64  `xorm:"pk autoincr"`
	Nickname  string `xorm:"not null default '' varchar(128) index" json:"nickname"`
	CreatedAt int64  `xorm:"not null default 0 int"`
}

func CreateRobotManager(info *RobotManagerUser) error {
	if info.Nickname == "" {
		return fmt.Errorf("robot manager nickname[%s] cannot be nil.", info.Nickname)
	}

	now := time.Now().Unix()
	info.CreatedAt = now

	_, err := x.Insert(info)
	if err != nil {
		holmes.Error("create robot manager error: %v", err)
		return err
	}
	holmes.Info("create robot manager[%v] success.", info)

	return nil
}

func GetRobotManager(info *RobotManagerUser) (bool, error) {
	has, err := x.Where("nickname = ?", info.Nickname).Get(info)
	if err != nil {
		return false, err
	}
	if !has {
		holmes.Debug("cannot find robot manager from unionid[%s]", info.Nickname)
		return false, nil
	}
	return true, nil
}
