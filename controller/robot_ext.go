package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/reechou/holmes"
	"github.com/reechou/robot-manager/config"
	"github.com/reechou/robot-manager/models"
)

const (
	ROBOT_ALL_ROBOTS_URI    = "/allrobots"
	ROBOT_SEND_MSGS_URI     = "/sendmsgs"
	ROBOT_FIND_FRIEND_URI   = "/findfriend"
	ROBOT_REMARK_FRIEND_URI = "/remarkfriend"
	ROBOT_GROUP_TIREN_URI   = "/grouptiren"
)

type RobotExt struct {
	client *http.Client
	cfg    *config.Config
}

func NewRobotExt(cfg *config.Config) *RobotExt {
	return &RobotExt{
		client: &http.Client{},
		cfg:    cfg,
	}
}

func (self *RobotExt) AllLoginRobots() (interface{}, error) {
	url := "http://" + self.cfg.RobotHost.Host + ROBOT_ALL_ROBOTS_URI
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		holmes.Error("http new request error: %v", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := self.client.Do(req)
	if err != nil {
		holmes.Error("http do request error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	rspBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		holmes.Error("ioutil ReadAll error: %v", err)
		return nil, err
	}
	var response WxResponse
	err = json.Unmarshal(rspBody, &response)
	if err != nil {
		holmes.Error("json decode error: %v [%s]", err, string(rspBody))
		return nil, err
	}
	if response.Code != 0 {
		holmes.Error("get all login robots result code error: %d %s", response.Code, response.Msg)
		return nil, fmt.Errorf("get all login robots result error.")
	}

	return response.Data, nil
}

func (self *RobotExt) GroupTiren(request *RobotGroupTirenReq) (*GroupUserInfo, error) {
	// get robot
	robot := &models.Robot{
		RobotWx: request.WechatNick,
	}
	has, err := models.GetRobot(robot)
	if err != nil {
		holmes.Error("get robot error: %v", err)
		return nil, err
	}
	if !has {
		return nil, fmt.Errorf("cannot found this robot[%s]", robot)
	}

	reqBytes, err := json.Marshal(request)
	if err != nil {
		holmes.Error("json encode error: %v", err)
		return nil, err
	}

	url := "http://" + robot.Ip + robot.OfPort + ROBOT_GROUP_TIREN_URI
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBytes))
	if err != nil {
		holmes.Error("http new request error: %v", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := self.client.Do(req)
	if err != nil {
		holmes.Error("http do request error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	rspBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		holmes.Error("ioutil ReadAll error: %v", err)
		return nil, err
	}
	var response WxGroupTirenResponse
	err = json.Unmarshal(rspBody, &response)
	if err != nil {
		holmes.Error("json decode error: %v [%s]", err, string(rspBody))
		return nil, err
	}
	if response.Code != 0 {
		holmes.Error("group tiren result code error: %d %s", response.Code, response.Msg)
		return nil, fmt.Errorf("group tiren result error.")
	}

	return &response.Data, nil
}

func (self *RobotExt) RemarkFriend(request *RobotRemarkFriendReq) error {
	// get robot
	robot := &models.Robot{
		RobotWx: request.WechatNick,
	}
	has, err := models.GetRobot(robot)
	if err != nil {
		holmes.Error("get robot error: %v", err)
		return err
	}
	if !has {
		return fmt.Errorf("cannot found this robot[%s]", robot)
	}

	reqBytes, err := json.Marshal(request)
	if err != nil {
		holmes.Error("json encode error: %v", err)
		return err
	}

	url := "http://" + robot.Ip + robot.OfPort + ROBOT_REMARK_FRIEND_URI
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBytes))
	if err != nil {
		holmes.Error("http new request error: %v", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := self.client.Do(req)
	if err != nil {
		holmes.Error("http do request error: %v", err)
		return err
	}
	defer resp.Body.Close()
	rspBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		holmes.Error("ioutil ReadAll error: %v", err)
		return err
	}
	var response WxResponse
	err = json.Unmarshal(rspBody, &response)
	if err != nil {
		holmes.Error("json decode error: %v [%s]", err, string(rspBody))
		return err
	}
	if response.Code != 0 {
		holmes.Error("remark friend result code error: %d %s", response.Code, response.Msg)
		return fmt.Errorf("remark friend result error.")
	}

	return nil
}

func (self *RobotExt) FindFriend(request *RobotFindFriendReq) (*UserFriend, error) {
	// get robot
	robot := &models.Robot{
		RobotWx: request.WechatNick,
	}
	has, err := models.GetRobot(robot)
	if err != nil {
		holmes.Error("get robot error: %v", err)
		return nil, err
	}
	if !has {
		return nil, fmt.Errorf("cannot found this robot[%s]", robot)
	}

	reqBytes, err := json.Marshal(request)
	if err != nil {
		holmes.Error("json encode error: %v", err)
		return nil, err
	}

	url := "http://" + robot.Ip + robot.OfPort + ROBOT_FIND_FRIEND_URI
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBytes))
	if err != nil {
		holmes.Error("http new request error: %v", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := self.client.Do(req)
	if err != nil {
		holmes.Error("http do request error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	rspBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		holmes.Error("ioutil ReadAll error: %v", err)
		return nil, err
	}
	var response WxFindFriendResponse
	err = json.Unmarshal(rspBody, &response)
	if err != nil {
		holmes.Error("json decode error: %v [%s]", err, string(rspBody))
		return nil, err
	}
	if response.Code != 0 {
		holmes.Error("find friend result code error: %d %s", response.Code, response.Msg)
		return nil, fmt.Errorf("find friend result error.")
	}

	return &response.Data, nil
}

func (self *RobotExt) SendMsgs(robotWx string, msg *SendMsgInfo) error {
	holmes.Debug("msg: %v", msg)
	// get robot
	robot := &models.Robot{
		RobotWx: robotWx,
	}
	has, err := models.GetRobot(robot)
	if err != nil {
		holmes.Error("get robot error: %v", err)
		return err
	}
	if !has {
		return fmt.Errorf("cannot found this robot[%s]", robot)
	}

	reqBytes, err := json.Marshal(msg)
	if err != nil {
		holmes.Error("json encode error: %v", err)
		return err
	}

	url := "http://" + robot.Ip + robot.OfPort + ROBOT_SEND_MSGS_URI
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBytes))
	if err != nil {
		holmes.Error("http new request error: %v", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := self.client.Do(req)
	if err != nil {
		holmes.Error("http do request error: %v", err)
		return err
	}
	defer resp.Body.Close()
	rspBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		holmes.Error("ioutil ReadAll error: %v", err)
		return err
	}
	var response SendMsgResponse
	err = json.Unmarshal(rspBody, &response)
	if err != nil {
		holmes.Error("json decode error: %v [%s]", err, string(rspBody))
		return err
	}
	if response.Code != 0 {
		holmes.Error("send msg[%v] result code error: %d %s", msg, response.Code, response.Msg)
		return fmt.Errorf("send msg result error.")
	}

	return nil
}
