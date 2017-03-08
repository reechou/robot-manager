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
