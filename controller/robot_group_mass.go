package controller

import (
	"encoding/json"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/reechou/holmes"
	"github.com/reechou/robot-manager/config"
	"github.com/reechou/robot-manager/models"
)

const (
	GROUP_MASS_TYPE_ALL           = 1
	GROUP_MASS_TYPE_SELECT_GROUPS = 2
	GROUP_MASS_WORKER             = 1024
)

var (
	RANDOM_MSG_ADD = []string{
		".", "..", "↭",
		"★", "✔", "↧",
		"↩", "⇤", "⇜",
		"↞", "↜", "┄",
		"-", "--", "^", "^_^",
		"!", "!!", "↮",
		"！", "•", "“",
		"[机智]", "[机智][机智]",
		"♥", "♥♥", "♥♥♥",
		"─", "↕↕", "↕",
		"☈", "✓", "☑",
		"⊰", "⊱", "†",
		"↓", "ˉ", "﹀",
		"﹏", "˜", "ˆ",
		"﹡", "≑", "≐",
		"≍", "≎", "≏",
		"≖", "≗", "≡",
	}
)

type RobotMsgInfo struct {
	MsgType string `json:"msgType"`
	Msg     string `json:"msg"`
}

type GroupMassInfo struct {
	RobotWx         string       `json:"robotWx"`
	Msg             RobotMsgInfo `json:"msg"`
	GroupMassType   int          `json:"groupMassType"`
	GroupNamePrefix string       `json:"groupNamePrefix,omitempty"`
	GroupList       []string     `json:"groupList,omitempty"`
}

type RobotGroupMass struct {
	cfg      *config.Config
	robotExt *RobotExt
	wg       sync.WaitGroup

	gmChan chan *GroupMassInfo
	stop   chan struct{}
	done   chan struct{}
}

func NewRobotGroupMass(cfg *config.Config, robotExt *RobotExt) *RobotGroupMass {
	rgm := &RobotGroupMass{
		cfg:      cfg,
		robotExt: robotExt,
		gmChan:   make(chan *GroupMassInfo, 1024),
		stop:     make(chan struct{}),
		done:     make(chan struct{}),
	}

	rgm.initWorkers()

	return rgm
}

func (self *RobotGroupMass) Stop() {
	close(self.stop)
	self.wg.Wait()
}

func (self *RobotGroupMass) DoGroupMass(gm *GroupMassInfo) {
	select {
	case self.gmChan <- gm:
	case <-self.stop:
		return
	}
}

func (self *RobotGroupMass) initWorkers() {
	holmes.Info("group mass init workers: %d", GROUP_MASS_WORKER)
	for i := 0; i < GROUP_MASS_WORKER; i++ {
		self.wg.Add(1)
		go self.runWorker(self.stop)
	}
}

func (self *RobotGroupMass) runWorker(stop chan struct{}) {
	for {
		select {
		case gm := <-self.gmChan:
			self.handleGroupMass(gm)
		case <-stop:
			self.wg.Done()
			return
		}
	}
}

func (self *RobotGroupMass) handleGroupMass(gm *GroupMassInfo) {
	holmes.Debug("handle group mass[%v] start.", gm)
	massRecord := self.recordGroupMass(gm)
	switch gm.GroupMassType {
	case GROUP_MASS_TYPE_ALL:
		robot := &models.Robot{
			RobotWx: gm.RobotWx,
		}
		has, err := models.GetRobot(robot)
		if err != nil {
			holmes.Error("get robot error: %v", err)
			return
		}
		if !has {
			holmes.Error("cannot found robot[%s]", gm.RobotWx)
			return
		}
		groupList, err := models.GetAllRobotGroupList(robot.ID)
		if err != nil {
			holmes.Error("get all robot group list error: %v", err)
			return
		}
		for _, v := range groupList {
			if strings.HasPrefix(v.GroupNickName, gm.GroupNamePrefix) {
				gsmi := &GroupSendMsgInfo{
					RobotWx:       robot.RobotWx,
					GroupUserName: v.UserName,
					GroupNickName: v.GroupNickName,
					Msg:           &gm.Msg,
				}
				self.sendMsgs(gsmi)
				time.Sleep(2 * time.Second)
			}
		}
	case GROUP_MASS_TYPE_SELECT_GROUPS:
		for _, v := range gm.GroupList {
			gsmi := &GroupSendMsgInfo{
				RobotWx:       gm.RobotWx,
				GroupNickName: v,
				Msg:           &gm.Msg,
			}
			self.sendMsgs(gsmi)
			time.Sleep(time.Duration(self.cfg.GroupMassInterval) * time.Second)
		}
	}
	holmes.Debug("handle group mass[%v] success.", gm)
	if massRecord.ID != 0 {
		massRecord.Status = models.ROBOT_GROUP_MASS_STATUS_END
		models.UpdateRobotGroupMassStatus(massRecord)
	}
}

func (self *RobotGroupMass) recordGroupMass(gm *GroupMassInfo) *models.RobotGroupMass {
	groupMass, _ := json.Marshal(gm)
	record := &models.RobotGroupMass{
		RobotWx:          gm.RobotWx,
		GroupMassContent: string(groupMass),
		Status:           models.ROBOT_GROUP_MASS_STATUS_START,
	}
	models.CreateRobotGroupMass(record)
	return record
}

type GroupSendMsgInfo struct {
	RobotWx       string
	GroupUserName string
	GroupNickName string
	Msg           *RobotMsgInfo
}

func (self *RobotGroupMass) sendMsgs(msg *GroupSendMsgInfo) {
	var sendReq SendMsgInfo
	offset := rand.Intn(len(RANDOM_MSG_ADD))
	msgStr := msg.Msg.Msg + RANDOM_MSG_ADD[offset]
	sendReq.SendMsgs = append(sendReq.SendMsgs, SendBaseInfo{
		WechatNick: msg.RobotWx,
		ChatType:   CHAT_TYPE_GROUP,
		UserName:   msg.GroupUserName,
		NickName:   msg.GroupNickName,
		MsgType:    msg.Msg.MsgType,
		Msg:        msgStr,
	})
	self.robotExt.SendMsgs(msg.RobotWx, &sendReq)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
