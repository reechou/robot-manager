package controller

type RobotInfo struct {
	RobotWxNick string `json:"robot"`
	RunTime     int64  `json:"runTime"`
}

type RobotFindFriendReq struct {
	WechatNick string `json:"wechatNick"`
	UserName   string `json:"username"`
	NickName   string `json:"nickname"`
}

type RobotRemarkFriendReq struct {
	WechatNick string `json:"wechatNick"`
	UserName   string `json:"username"`
	NickName   string `json:"nickname"`
	Remark     string `json:"remark"`
}

type RobotGroupTirenReq struct {
	WechatNick     string `json:"wechatNick"`
	GroupUserName  string `json:"groupUserName"`
	GroupNickName  string `json:"groupNickName"`
	MemberUserName string `json:"memberUserName"`
	MemberNickName string `json:"memberNickName"`
}

type WxFindFriendResponse struct {
	Code int        `json:"code"`
	Msg  string     `json:"msg"`
	Data UserFriend `json:"data"`
}

type WxGroupTirenResponse struct {
	Code int           `json:"code"`
	Msg  string        `json:"msg"`
	Data GroupUserInfo `json:"data"`
}

type WxResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
