package controller

import (
	"fmt"

	"github.com/reechou/holmes"
	"github.com/reechou/robot-manager/models"
)

const (
	ROBOT_MANAGER_PREFIX = "REAL_MANAGER_WOWOWO_"
)

func (self *Logic) handleCreateManager(req *RobotCreateManagerReq) error {
	findFriendReq := &RobotFindFriendReq{
		WechatNick: req.RobotWx,
		NickName:   req.Nickname,
	}
	uf, err := self.robotExt.FindFriend(findFriendReq)
	if err != nil {
		holmes.Error("robot find friend error: %v", err)
		return err
	}
	if uf == nil || (uf.UserName == "" && uf.NickName == "") {
		holmes.Error("cannot found this man[%s] from robot[%s]", req.Nickname, req.RobotWx)
		return fmt.Errorf("cannot found this man[%s] from robot[%s]", req.Nickname, req.RobotWx)
	}
	remark := fmt.Sprintf("%s%s", ROBOT_MANAGER_PREFIX, uf.NickName)
	remarkFriendReq := &RobotRemarkFriendReq{
		WechatNick: req.RobotWx,
		UserName:   uf.UserName,
		NickName:   uf.NickName,
		Remark:     remark,
	}
	err = self.robotExt.RemarkFriend(remarkFriendReq)
	if err != nil {
		holmes.Error("robot remark friend error: %v", err)
		return err
	}
	rmu := &models.RobotManagerUser{
		Nickname: remark,
	}
	err = models.CreateRobotManager(rmu)
	if err != nil {
		holmes.Error("create robot manager error: %v", err)
		return err
	}
	return nil
}
