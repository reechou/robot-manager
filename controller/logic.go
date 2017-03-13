package controller

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/reechou/holmes"
	"github.com/reechou/robot-manager/config"
	"github.com/reechou/robot-manager/models"
)

type Logic struct {
	sync.Mutex
	robotExt  *RobotExt
	msgFilter *RobotMsgFilter
	groupMass *RobotGroupMass

	cfg *config.Config
}

func NewLogic(cfg *config.Config) *Logic {
	l := &Logic{
		cfg:      cfg,
		robotExt: NewRobotExt(cfg),
	}

	l.msgFilter = NewRobotMsgFilter(cfg, l.robotExt)
	l.groupMass = NewRobotGroupMass(cfg, l.robotExt)

	models.InitDB(cfg)
	l.init()

	return l
}

func (self *Logic) init() {
	http.HandleFunc("/robot/save_groups", self.RobotSaveGroups)
	http.HandleFunc("/robot/receive_msg", self.RobotReceiveMsg)
	http.HandleFunc("/robot/group_mass", self.RobotGroupMass)
	http.HandleFunc("/robot/create_manager", self.CreateRobotManager)

	http.HandleFunc("/manager/get_all_robots", self.GetAllRobots)
	http.HandleFunc("/manager/get_groups", self.GetRobotGroups)
	http.HandleFunc("/manager/get_new_group_chat", self.GetRobotGroupChatNew)
	http.HandleFunc("/manager/get_group_mass", self.GetRobotGroupMassList)
	http.HandleFunc("/manager/get_robot_group_mass", self.GetRobotGroupMassFromRobot)
	http.HandleFunc("/manager/send_group_msg", self.RobotSendGroupMsg)
	http.HandleFunc("/manager/get_all_login_robots", self.GetLoginRobotList)
	http.HandleFunc("/manager/robot_group_tiren", self.RobotGroupTiren)
}

func (self *Logic) Run() {
	defer holmes.Start(holmes.LogFilePath("./log"),
		holmes.EveryDay,
		holmes.AlsoStdout,
		holmes.DebugLevel).Stop()

	if self.cfg.Debug {
		EnableDebug()
	}

	holmes.Info("server starting on[%s]..", self.cfg.Host)
	holmes.Infoln(http.ListenAndServe(self.cfg.Host, nil))
}

func WriteErrorResponse(w http.ResponseWriter, rsp *Response) {
	rsp.Code = RESPONSE_ERR
	WriteJSON(w, http.StatusOK, rsp)
}

func WriteJSON(w http.ResponseWriter, code int, v interface{}) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "x-requested-with,content-type")
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

func WriteBytes(w http.ResponseWriter, code int, v []byte) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "x-requested-with,content-type")
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(code)
	w.Write(v)
}

func EnableDebug() {

}
