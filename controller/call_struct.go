package controller

import (
	"github.com/reechou/robot-manager/models"
)

const (
	RESPONSE_OK = iota
	RESPONSE_ERR
)

type Response struct {
	Code int64       `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

type GetRobotGroupsRsp struct {
	Count int64               `json:"count"`
	List  []models.RobotGroup `json:"list"`
}

type GetRobotGroupMassListRsp struct {
	Count int64                   `json:"count"`
	List  []models.RobotGroupMass `json:"list"`
}

type RobotSaveGroupsReq struct {
	RobotWx string    `json:"robotWx"`
	Groups  []WxGroup `json:"groups"`
}

type RobotCreateManagerReq struct {
	RobotWx  string `json:"robotWx"`
	Nickname string `json:"nickname"`
}

type GetRobotGroupsReq struct {
	RobotId int64 `json:"robotId"`
	Offset  int64 `json:"offset"`
	Num     int64 `json:"num"`
}

type GetRobotGroupChatNewReq struct {
	RobotId   int64 `json:"robotId"`
	Timestamp int64 `json:"timestamp"`
}

type GetRobotGroupMassReq struct {
	Offset int64 `json:"offset"`
	Num    int64 `json:"num"`
}

type GetRobotGroupMassFromRobotReq struct {
	RobotWx string `json:"robotWx"`
}

type SendGroupMsgReq struct {
	RobotWx       string `json:"robotWx"`
	GroupUserName string `json:"groupUserName"`
	GroupNickName string `json:"groupNickName"`
	MsgType       string `json:"msgType"`
	Msg           string `json:"msg"`
}
