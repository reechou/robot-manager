package models

import (
	"fmt"
	"time"

	"github.com/reechou/holmes"
)

type RobotGroup struct {
	ID             int64  `xorm:"pk autoincr"`
	RobotId        int64  `xorm:"not null default 0 int index" json:"robotId"`
	RobotWx        string `xorm:"varchar(128) not null default '' unique(uni_robot_username)" json:"robotWx"`
	UserName       string `xorm:"varchar(128) not null default '' unique(uni_robot_username)" json:"userName"`
	GroupNickName  string `xorm:"not null default '' varchar(256) index" json:"groupNickName"`
	GroupMemberNum int64  `xorm:"not null default 0 int" json:"groupMemberNum"`
	CreatedAt      int64  `xorm:"not null default 0 int" json:"createAt"`
	UpdatedAt      int64  `xorm:"not null default 0 int" json:"-"`
}

func CreateRobotGroup(info *RobotGroup) error {
	if info.RobotWx == "" {
		return fmt.Errorf("wx robot wx[%s] cannot be nil.", info.RobotWx)
	}

	now := time.Now().Unix()
	info.CreatedAt = now
	info.UpdatedAt = now

	_, err := x.Insert(info)
	if err != nil {
		holmes.Error("create robot group error: %v", err)
		return err
	}
	holmes.Info("create robot group[%v] success.", info)

	return nil
}

func CreateRobotGroupList(list []RobotGroup) error {
	if len(list) == 0 {
		return nil
	}
	_, err := x.Insert(&list)
	if err != nil {
		holmes.Error("create robot group list error: %v", err)
		return err
	}
	return nil
}

func GetRobotGroupListCount(robotId int64) (int64, error) {
	count, err := x.Where("robot_id = ?", robotId).Count(&RobotGroup{RobotId: robotId})
	if err != nil {
		holmes.Error("robot_id[%d] get robot group list count error: %v", robotId, err)
		return 0, err
	}
	return count, nil
}

func GetRobotGroupList(robotId, offset, num int64) ([]RobotGroup, error) {
	var list []RobotGroup
	err := x.Where("robot_id = ?", robotId).Desc("created_at").Limit(int(num), int(offset)).Find(&list)
	if err != nil {
		holmes.Error("robot_id[%d] get robot group list error: %v", robotId, err)
		return nil, err
	}
	return list, nil
}

func GetAllRobotGroupList(robotId int64) ([]RobotGroup, error) {
	var list []RobotGroup
	err := x.Where("robot_id = ?", robotId).Find(&list)
	if err != nil {
		holmes.Error("robot_id[%d] get robot all group list error: %v", robotId, err)
		return nil, err
	}
	return list, nil
}

func GetRobotGroup(info *RobotGroup) (bool, error) {
	has, err := x.Where("robot_wx = ?", info.RobotWx).And("group_nick_name = ?", info.GroupNickName).Get(info)
	if err != nil {
		return false, err
	}
	if !has {
		holmes.Debug("cannot find robot group from unionid[%s-%s]", info.RobotWx, info.GroupNickName)
		return false, nil
	}
	return true, nil
}

func GetRobotGroupFromUserName(info *RobotGroup) (bool, error) {
	has, err := x.Where("robot_wx = ?", info.RobotWx).And("user_name = ?", info.UserName).Get(info)
	if err != nil {
		return false, err
	}
	if !has {
		holmes.Debug("cannot find robot group from user_name[%s]", info.UserName)
		return false, nil
	}
	return true, nil
}

func UpdateRobotGroupName(info *RobotGroup) error {
	info.UpdatedAt = time.Now().Unix()
	_, err := x.ID(info.ID).Cols("group_nick_name", "updated_at").Update(info)
	return err
}

func UpdateRobotGroupUserName(info *RobotGroup) error {
	info.UpdatedAt = time.Now().Unix()
	_, err := x.ID(info.ID).Cols("user_name", "updated_at").Update(info)
	return err
}
