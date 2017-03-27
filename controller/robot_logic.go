package controller

import (
	"errors"

	"github.com/reechou/holmes"
	"github.com/reechou/robot-manager/models"
)

var (
	MSG_GROUP_GET_NONE = errors.New("cannot found msg group.")
)

func (self *Logic) HandleReceiveMsg(msg *ReceiveMsgInfo) {
	switch msg.BaseInfo.ReceiveEvent {
	case RECEIVE_EVENT_MSG:
		self.AddChat(msg)
	case RECEIVE_EVENT_MOD_GROUP_ADD_DETAIL:
		self.RobotGroupMod(msg)
	}
}

func (self *Logic) RobotGroupMod(msg *ReceiveMsgInfo) {
	holmes.Debug("robot group mod: %v", msg)
	rbl := &models.RobotBlacklist{
		RobotWx:  msg.BaseInfo.WechatNick,
		UserName: msg.BaseInfo.FromMemberUserName,
	}
	has, err := models.GetRobotBlacklistFromUsername(rbl)
	if err != nil {
		holmes.Error("get robot blacklist from username error: %v", err)
		return
	}
	if has {
		holmes.Debug("has in this.")
		tirenReq := &RobotGroupTirenReq{
			WechatNick:     msg.BaseInfo.WechatNick,
			GroupUserName:  msg.BaseInfo.FromUserName,
			GroupNickName:  msg.BaseInfo.FromGroupName,
			MemberUserName: msg.BaseInfo.FromMemberUserName,
			MemberNickName: msg.BaseInfo.FromNickName,
		}
		_, err := self.robotExt.GroupTiren(tirenReq)
		if err != nil {
			holmes.Error("group tiren error: %v", err)
			return
		}
	}
}

func (self *Logic) AddChat(msg *ReceiveMsgInfo) {
	switch msg.BaseInfo.FromType {
	case CHAT_TYPE_GROUP:
		group, err := self.getRobotGroup(msg.BaseInfo.WechatNick, msg.BaseInfo.FromUserName, msg.BaseInfo.FromGroupName)
		if err != nil {
			holmes.Error("get msg robot group error: %v", err)
			return
		}
		group.GroupMemberNum = int64(msg.GroupMemberNum)
		models.UpdateRobotGroupGroupMemberNum(group)

		rgc := &models.RobotGroupChat{
			RobotId:        group.RobotId,
			RobotWx:        group.RobotWx,
			GroupId:        group.ID,
			GroupName:      msg.BaseInfo.FromGroupName,
			GroupUserName:  msg.BaseInfo.FromUserName,
			MemberUserName: msg.BaseInfo.FromMemberUserName,
			FromName:       msg.BaseInfo.FromNickName,
			MsgType:        msg.MsgType,
			Content:        msg.Msg,
			MediaTempUrl:   msg.MediaTempUrl,
		}
		err = models.CreateRobotGroupChat(rgc)
		if err != nil {
			holmes.Error("create robot group chat error: %v", err)
		}
		self.msgFilter.FilterMsg(msg)
	}
}

func (self *Logic) getRobotGroup(robotWx, username, nickname string) (*models.RobotGroup, error) {
	robot := &models.Robot{
		RobotWx: robotWx,
	}
	has, err := models.GetRobot(robot)
	if err != nil {
		holmes.Error("get robot error: %v", err)
		return nil, err
	}
	if !has {
		return nil, err
	}
	group := &models.RobotGroup{
		RobotWx:  robotWx,
		UserName: username,
	}
	has, err = models.GetRobotGroupFromUserName(group)
	if err != nil {
		holmes.Error("get robot group from username[%v] error: %v", group, err)
		return nil, err
	}
	if has {
		if nickname != group.GroupNickName && nickname != "" {
			group.GroupNickName = nickname
			err = models.UpdateRobotGroupName(group)
			if err != nil {
				holmes.Error("update robot group name[%v] error: %v", group, err)
			}
		}
		return group, nil
	}
	group = &models.RobotGroup{
		RobotWx:       robotWx,
		GroupNickName: nickname,
	}
	has, err = models.GetRobotGroup(group)
	if err != nil {
		holmes.Error("get robot group from [%v] error: %v", group, err)
		return nil, err
	}
	if !has {
		holmes.Error("cannot found this[%s %s %s] group.", robotWx, username, nickname)
		return nil, MSG_GROUP_GET_NONE
	}
	if group.UserName != username {
		group.UserName = username
		err = models.UpdateRobotGroupUserName(group)
		if err != nil {
			holmes.Error("update group[%v] username error: %v", group, err)
		}
	}
	return group, nil
}
