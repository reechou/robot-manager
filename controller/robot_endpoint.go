package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/reechou/holmes"
	"github.com/reechou/robot-manager/models"
)

func (self *Logic) RobotSaveGroups(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteJSON(w, http.StatusOK, nil)
		return
	}

	req := &RobotSaveGroupsReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("RobotSaveGroups json decode error: %v", err)
		return
	}
	rsp := &Response{Code: RESPONSE_OK}

	robot := &models.Robot{
		RobotWx: req.RobotWx,
	}
	has, err := models.GetRobot(robot)
	if err != nil {
		holmes.Error("get robot error: %v", err)
		rsp.Code = RESPONSE_ERR
		WriteJSON(w, http.StatusOK, rsp)
		return
	}
	if !has {
		holmes.Error("cannot found this robot[%s]", req.RobotWx)
		rsp.Code = RESPONSE_ERR
		WriteJSON(w, http.StatusOK, rsp)
		return
	}

	now := time.Now().Unix()
	var list []models.RobotGroup
	for _, v := range req.Groups {
		list = append(list, models.RobotGroup{
			RobotId:        robot.ID,
			RobotWx:        req.RobotWx,
			UserName:       v.UserName,
			GroupNickName:  v.NickName,
			GroupMemberNum: int64(v.GroupMemberNum),
			CreatedAt:      now,
			UpdatedAt:      now,
		})
	}
	err = models.CreateRobotGroupList(list)
	if err != nil {
		holmes.Error("Error robot save groups error: %v", err)
		rsp.Code = RESPONSE_ERR
		rsp.Msg = fmt.Sprintf("Error robot save groups error: %v", err)
	}

	WriteJSON(w, http.StatusOK, rsp)
}

func (self *Logic) RobotReceiveMsg(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteJSON(w, http.StatusOK, nil)
		return
	}
	req := &ReceiveMsgInfo{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("RobotReceiveMsg json decode error: %v", err)
		return
	}
	self.HandleReceiveMsg(req)

	rsp := &Response{Code: RESPONSE_OK}
	WriteJSON(w, http.StatusOK, rsp)
}

func (self *Logic) RobotGroupMass(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteJSON(w, http.StatusOK, nil)
		return
	}
	req := &GroupMassInfo{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("RobotGroupMass json decode error: %v", err)
		return
	}
	self.groupMass.DoGroupMass(req)

	rsp := &Response{Code: RESPONSE_OK}
	WriteJSON(w, http.StatusOK, rsp)
}

func (self *Logic) CreateRobotManager(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteJSON(w, http.StatusOK, nil)
		return
	}
	req := &RobotCreateManagerReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("CreateRobotManager json decode error: %v", err)
		return
	}

	rsp := &Response{Code: RESPONSE_OK}

	err := self.handleCreateManager(req)
	if err != nil {
		holmes.Error("handle create manager error: %v", err)
		rsp.Code = RESPONSE_ERR
	}

	WriteJSON(w, http.StatusOK, rsp)
}
