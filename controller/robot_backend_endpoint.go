package controller

import (
	"encoding/json"
	"net/http"

	"github.com/reechou/holmes"
	"github.com/reechou/robot-manager/models"
)

func (self *Logic) GetAllRobots(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteJSON(w, http.StatusOK, nil)
		return
	}

	rsp := &Response{Code: RESPONSE_OK}
	list, err := models.GetRobotList()
	if err != nil {
		holmes.Error("get all robot list error: %v", err)
		WriteErrorResponse(w, rsp)
		return
	}
	rsp.Data = list
	WriteJSON(w, http.StatusOK, rsp)
}

func (self *Logic) GetRobotGroups(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteJSON(w, http.StatusOK, nil)
		return
	}

	req := &GetRobotGroupsReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("GetRobotGroups json decode error: %v", err)
		return
	}

	rsp := &Response{Code: RESPONSE_OK}

	count, err := models.GetRobotGroupListCount(req.RobotId)
	if err != nil {
		holmes.Error("get robot group list count error: %v", err)
		WriteErrorResponse(w, rsp)
		return
	}
	list, err := models.GetRobotGroupList(req.RobotId, req.Offset, req.Num)
	if err != nil {
		holmes.Error("get robot group list error: %v", err)
		WriteErrorResponse(w, rsp)
		return
	}
	rsp.Data = &GetRobotGroupsRsp{
		Count: count,
		List:  list,
	}
	WriteJSON(w, http.StatusOK, rsp)
}

func (self *Logic) GetRobotGroupChatNew(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteJSON(w, http.StatusOK, nil)
		return
	}

	req := &GetRobotGroupChatNewReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("GetRobotGroupChatNew json decode error: %v", err)
		return
	}

	rsp := &Response{Code: RESPONSE_OK}

	list, err := models.GetRobotGroupNewChatList(req.RobotId, req.Timestamp)
	if err != nil {
		holmes.Error("get robot group chat list error: %v", err)
		WriteErrorResponse(w, rsp)
		return
	}
	rsp.Data = list
	WriteJSON(w, http.StatusOK, rsp)
}

func (self *Logic) GetRobotGroupMassList(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteJSON(w, http.StatusOK, nil)
		return
	}

	req := &GetRobotGroupMassReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("GetRobotGroupMassList json decode error: %v", err)
		return
	}

	rsp := &Response{Code: RESPONSE_OK}

	count, err := models.GetRobotGroupMassListCount()
	if err != nil {
		holmes.Error("get robot group mass list count error: %v", err)
		WriteErrorResponse(w, rsp)
		return
	}
	list, err := models.GetRobotGroupMassList(req.Offset, req.Num)
	if err != nil {
		holmes.Error("get robot group mass list error: %v", err)
		WriteErrorResponse(w, rsp)
		return
	}
	rsp.Data = &GetRobotGroupMassListRsp{
		Count: count,
		List:  list,
	}
	WriteJSON(w, http.StatusOK, rsp)
}

func (self *Logic) GetRobotGroupMassFromRobot(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteJSON(w, http.StatusOK, nil)
		return
	}

	req := &GetRobotGroupMassFromRobotReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("GetRobotGroupMassFromRobot json decode error: %v", err)
		return
	}

	rsp := &Response{Code: RESPONSE_OK}

	list, err := models.GetRobotGroupMassListFromRobot(req.RobotWx)
	if err != nil {
		holmes.Error("get robot group mass list from robot error: %v", err)
		WriteErrorResponse(w, rsp)
		return
	}
	rsp.Data = list
	WriteJSON(w, http.StatusOK, rsp)
}

func (self *Logic) RobotSendGroupMsg(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteJSON(w, http.StatusOK, nil)
		return
	}

	req := &SendGroupMsgReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("RobotSendGroupMsg json decode error: %v", err)
		return
	}
	rsp := &Response{Code: RESPONSE_OK}

	var sendReq SendMsgInfo
	sendReq.SendMsgs = append(sendReq.SendMsgs, SendBaseInfo{
		WechatNick: req.RobotWx,
		ChatType:   CHAT_TYPE_GROUP,
		UserName:   req.GroupUserName,
		NickName:   req.GroupNickName,
		MsgType:    req.MsgType,
		Msg:        req.Msg,
	})
	err := self.robotExt.SendMsgs(req.RobotWx, &sendReq)
	if err != nil {
		holmes.Error("group send msg[%v] error: %v", req, err)
		rsp.Data = RESPONSE_ERR
		WriteJSON(w, http.StatusOK, rsp)
		return
	}

	// 存储web发送消息
	robot := &models.Robot{
		RobotWx: req.RobotWx,
	}
	has, err := models.GetRobot(robot)
	if err != nil {
		holmes.Error("get robot error: %v", err)
		WriteJSON(w, http.StatusOK, rsp)
		return
	}
	if !has {
		WriteJSON(w, http.StatusOK, nil)
		return
	}
	group := &models.RobotGroup{
		RobotWx:  req.RobotWx,
		UserName: req.GroupUserName,
	}
	has, err = models.GetRobotGroupFromUserName(group)
	if err != nil {
		holmes.Error("get robot group from username[%v] error: %v", group, err)
		WriteJSON(w, http.StatusOK, nil)
		return
	}
	if !has {
		group = &models.RobotGroup{
			RobotWx:       req.RobotWx,
			GroupNickName: req.GroupNickName,
		}
		has, err = models.GetRobotGroup(group)
		if err != nil {
			holmes.Error("get robot group from [%v] error: %v", group, err)
			WriteJSON(w, http.StatusOK, nil)
			return
		}
		if !has {
			WriteJSON(w, http.StatusOK, nil)
			return
		}
	}

	rgc := &models.RobotGroupChat{
		RobotId:       robot.ID,
		RobotWx:       req.RobotWx,
		GroupId:       group.ID,
		GroupName:     req.GroupNickName,
		GroupUserName: req.GroupUserName,
		FromName:      req.RobotWx,
		MsgType:       req.MsgType,
		Content:       req.Msg,
		Source:        models.ROBOT_CHAT_SOURCE_FROM_WEB,
	}
	err = models.CreateRobotGroupChat(rgc)
	if err != nil {
		holmes.Error("create robot group chat error: %v", err)
	}

	WriteJSON(w, http.StatusOK, rsp)
}

func (self *Logic) GetLoginRobotList(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteJSON(w, http.StatusOK, nil)
		return
	}
	
	rsp := &Response{Code: RESPONSE_OK}
	list, err := self.robotExt.AllLoginRobots()
	if err != nil {
		holmes.Error("get all login robots error: %v", err)
		rsp.Data = RESPONSE_ERR
		WriteJSON(w, http.StatusOK, rsp)
		return
	}
	rsp.Data = list
	WriteJSON(w, http.StatusOK, rsp)
}

func (self *Logic) RobotGroupTiren(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteJSON(w, http.StatusOK, nil)
		return
	}
	
	req := &RobotGroupTirenReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("RobotGroupTiren json decode error: %v", err)
		return
	}
	
	rsp := &Response{Code: RESPONSE_OK}
	_, err := self.robotExt.GroupTiren(req)
	if err != nil {
		holmes.Error("robot group tiren error: %v", err)
		rsp.Data = RESPONSE_ERR
		WriteJSON(w, http.StatusOK, rsp)
		return
	}
	
	WriteJSON(w, http.StatusOK, rsp)
}
