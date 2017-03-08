package models

import (
	"fmt"
	"strconv"
	"time"

	"github.com/reechou/holmes"
)

const (
	ROBOT_GROUP_CHAT_TABLE_NUM = 10
)

const (
	ROBOT_CHAT_SOURCE_FROM_USER  = "来自用户"
	ROBOT_CHAT_SOURCE_FROM_WEB   = "来自web cms"
	ROBOT_CHAT_SOURCE_FROM_PHONE = "来自手机"
)

type RobotGroupChat struct {
	ID             int64  `xorm:"pk autoincr" json:"id"`
	RobotId        int64  `xorm:"not null default 0 int index" json:"robotId"`
	RobotWx        string `xorm:"not null default '' varchar(128) index" json:"robotWx"`
	GroupId        int64  `xorm:"not null default 0 int index" json:"groupId"`
	GroupName      string `xorm:"not null default '' varchar(256)" json:"groupName"`
	GroupUserName  string `xorm:"not null default '' varchar(128)" json:"groupUserName"`
	MemberUserName string `xorm:"not null default '' varchar(128)" json:"memberUserName"`
	FromName       string `xorm:"not null default '' varchar(256)" json:"fromName"`
	MsgType        string `xorm:"not null default '' varchar(16)" json:"msgType"`
	Content        string `xorm:"not null default '' varchar(768)" json:"content"`
	MediaTempUrl   string `xorm:"not null default '' varchar(256)" json:"mediaTempUrl,omitempty"`
	Source         string `xorm:"not null default '' varchar(64)" json:"source"`
	CreatedAt      int64  `xorm:"not null default 0 int index" json:"createdAt"`
}

func (self *RobotGroupChat) TableName() string {
	return "robot_group_chat_" + strconv.Itoa(int(self.RobotId)%ROBOT_GROUP_CHAT_TABLE_NUM)
}

func CreateRobotGroupChat(info *RobotGroupChat) error {
	if info.RobotWx == "" {
		return fmt.Errorf("wx robot wx[%s] cannot be nil.", info.RobotWx)
	}

	now := time.Now().Unix()
	info.CreatedAt = now

	_, err := x.Insert(info)
	if err != nil {
		holmes.Error("create robot group chat error: %v", err)
		return err
	}
	holmes.Info("create robot chat robot[%s] group[%s] fromUser[%s] content[%s] success.", info.RobotWx, info.GroupName, info.FromName, info.Content)

	return nil
}

func GetRobotGroupChatListCount(robotId, groupId int64) (int64, error) {
	count, err := x.Where("group_id = ?", groupId).Count(&RobotGroupChat{RobotId: robotId})
	if err != nil {
		holmes.Error("groupId[%d] get robot group chat list count error: %v", groupId, err)
		return 0, err
	}
	return count, nil
}

func GetRobotGroupChatList(robotId, groupId, offset, num int64) ([]RobotGroupChat, error) {
	var list []RobotGroupChat
	err := x.Table(&RobotGroupChat{RobotId: robotId}).Where("robot_id = ?", robotId).And("group_id = ?", groupId).Desc("created_at").Limit(int(num), int(offset)).Find(&list)
	if err != nil {
		holmes.Error("group_id[%d] get robot group chat list error: %v", groupId, err)
		return nil, err
	}
	return list, nil
}

func GetRobotGroupNewChatList(robotId, timestamp int64) ([]RobotGroupChat, error) {
	holmes.Debug("robotid[%d] t[%d] get new group chat list", robotId, timestamp)
	var list []RobotGroupChat
	err := x.Table(&RobotGroupChat{RobotId: robotId}).Where("robot_id = ?", robotId).And("created_at > ?", timestamp).Desc("created_at").Limit(50).Find(&list)
	if err != nil {
		holmes.Error("robot_id[%d] get robot group new chat list error: %v", robotId, err)
		return nil, err
	}
	return list, nil
}
