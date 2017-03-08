package models

import (
	"fmt"
	"time"

	"github.com/reechou/holmes"
)

const (
	ROBOT_GROUP_MASS_STATUS_START = iota
	ROBOT_GROUP_MASS_STATUS_END
)

type RobotGroupMass struct {
	ID               int64  `xorm:"pk autoincr" json:"id"`
	RobotWx          string `xorm:"not null default '' varchar(128) index" json:"robotWx"`
	GroupMassContent string `xorm:"not null default '' varchar(4096)" json:"content"`
	Status           int64  `xorm:"not null default 0 int" json:"status"`
	CreatedAt        int64  `xorm:"not null default 0 int index" json:"createdAt"`
	UpdatedAt        int64  `xorm:"not null default 0 int" json:"updatedAt"`
}

func CreateRobotGroupMass(info *RobotGroupMass) error {
	if info.RobotWx == "" {
		return fmt.Errorf("wx robot wx[%s] cannot be nil.", info.RobotWx)
	}

	now := time.Now().Unix()
	info.CreatedAt = now

	_, err := x.Insert(info)
	if err != nil {
		holmes.Error("create robot group mass error: %v", err)
		return err
	}
	holmes.Info("create robot chat robot[%s] group mass success.", info.RobotWx)

	return nil
}

func UpdateRobotGroupMassStatus(info *RobotGroupMass) error {
	info.UpdatedAt = time.Now().Unix()
	_, err := x.ID(info.ID).Cols("status", "updated_at").Update(info)
	if err != nil {
		holmes.Error("update group mass status error: %v", err)
	}
	return err
}

func GetRobotGroupMassListCount() (int64, error) {
	count, err := x.Count(&RobotGroupMass{})
	if err != nil {
		holmes.Error("get robot group mass list count error: %v", err)
		return 0, err
	}
	return count, nil
}

func GetRobotGroupMassList(offset, num int64) ([]RobotGroupMass, error) {
	var list []RobotGroupMass
	err := x.Desc("created_at").Limit(int(num), int(offset)).Find(&list)
	if err != nil {
		holmes.Error("get robot group mass list error: %v", err)
		return nil, err
	}
	return list, nil
}

func GetRobotGroupMassListFromRobot(robot string) ([]RobotGroupMass, error) {
	var list []RobotGroupMass
	err := x.Where("robot_wx = ?", robot).Desc("created_at").Find(&list)
	if err != nil {
		holmes.Error("get robot group mass list from robot error: %v", err)
		return nil, err
	}
	return list, nil
}
