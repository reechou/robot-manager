package controller

import (
	"strings"
	"time"

	"github.com/reechou/holmes"
	"github.com/reechou/robot-manager/config"
	"github.com/reechou/robot-manager/models"
)

type RobotMsgFilter struct {
	cfg      *config.Config
	robotExt *RobotExt

	msgChan chan *ReceiveMsgInfo
	stop    chan struct{}
	done    chan struct{}
}

func NewRobotMsgFilter(cfg *config.Config, robotExt *RobotExt) *RobotMsgFilter {
	rmf := &RobotMsgFilter{
		cfg:      cfg,
		robotExt: robotExt,
		msgChan:  make(chan *ReceiveMsgInfo, 1024),
		stop:     make(chan struct{}),
		done:     make(chan struct{}),
	}
	go rmf.run()

	return rmf
}

func (self *RobotMsgFilter) Stop() {
	close(self.stop)
	<-self.done
}

func (self *RobotMsgFilter) FilterMsg(msg *ReceiveMsgInfo) {
	select {
	case self.msgChan <- msg:
	case <-time.After(2 * time.Second):
		holmes.Error("filter msg maybe channal is full.")
	}
}

func (self *RobotMsgFilter) run() {
	for {
		select {
		case msg := <-self.msgChan:
			self.checkMsg(msg)
		case <-self.stop:
			close(self.done)
			return
		}
	}
}

func (self *RobotMsgFilter) checkMsg(msg *ReceiveMsgInfo) {
	if self.checkTiren(msg) {
		self.handleTiren(msg)
	}
}

func (self *RobotMsgFilter) handleTiren(msg *ReceiveMsgInfo) {
	manager := &models.RobotManagerUser{
		Nickname: msg.BaseInfo.FromNickName,
	}
	has, err := models.GetRobotManager(manager)
	if err != nil {
		holmes.Error("get robot manager error: %v", err)
		return
	}
	if !has {
		holmes.Error("nickname[%s] is not manager", msg.BaseInfo.FromNickName)
		return
	}
	idx := strings.Index(msg.Msg, "@")
	if idx == -1 {
		holmes.Debug("[%s] not found @, stop tiren", msg.Msg)
		return
	}
	if idx > len(msg.Msg)-2 {
		holmes.Error("idx[%d] > len(msg.Msg)[%d] - 1", idx, len(msg.Msg))
		return
	}
	memberNickName := msg.Msg[idx+1 : len(msg.Msg)-1]
	//memberNickName = strings.Replace(memberNickName, " ", "", -1)
	holmes.Debug("nickname[%s] group[%s] tiren[%s] memberNickName: %d", msg.BaseInfo.FromNickName, msg.BaseInfo.FromGroupName, memberNickName, len(memberNickName))
	tirenReq := &RobotGroupTirenReq{
		WechatNick:     msg.BaseInfo.WechatNick,
		GroupUserName:  msg.BaseInfo.FromUserName,
		GroupNickName:  msg.BaseInfo.FromGroupName,
		MemberNickName: memberNickName,
	}
	gui, err := self.robotExt.GroupTiren(tirenReq)
	if err != nil {
		holmes.Error("group tiren error: %v", err)
		return
	}
	rbl := &models.RobotBlacklist{
		RobotWx:  msg.BaseInfo.WechatNick,
		UserName: gui.UserName,
		NickName: gui.NickName,
	}
	err = models.CreateRobotBlacklist(rbl)
	if err != nil {
		holmes.Error("create robot blacklist error: %v", err)
	}
}

func (self *RobotMsgFilter) checkTiren(msg *ReceiveMsgInfo) bool {
	for _, v := range ROBOT_TIREN {
		if strings.HasPrefix(msg.Msg, v) {
			return true
		}
	}
	return false
}
