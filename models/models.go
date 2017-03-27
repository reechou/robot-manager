package models

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"github.com/reechou/holmes"
	"github.com/reechou/robot-manager/config"
)

var x *xorm.Engine

func InitDB(cfg *config.Config) {
	var err error
	x, err = xorm.NewEngine("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4",
		cfg.DBInfo.User,
		cfg.DBInfo.Pass,
		cfg.DBInfo.Host,
		cfg.DBInfo.DBName))
	if err != nil {
		holmes.Fatal("Fail to init new engine: %v", err)
	}
	//x.SetLogger(nil)
	x.SetMapper(core.GonicMapper{})
	x.TZLocation, _ = time.LoadLocation("Asia/Shanghai")
	// if need show raw sql in log
	if cfg.IfShowSqlLog {
		x.ShowSQL(true)
	}

	// sync tables
	if err = x.Sync2(new(RobotBlacklist)); err != nil {
		holmes.Fatal("Fail to sync database: %v", err)
	}
	if err = x.Sync2(new(RobotManagerUser)); err != nil {
		holmes.Fatal("Fail to sync database: %v", err)
	}
	if err = x.Sync2(new(RobotGroup)); err != nil {
		holmes.Fatal("Fail to sync database: %v", err)
	}
	if err = x.Sync2(new(Robot)); err != nil {
		holmes.Fatal("Fail to sync database: %v", err)
	}
	if err = x.Sync2(new(RobotGroupMass)); err != nil {
		holmes.Fatal("Fail to sync database: %v", err)
	}

	for i := 0; i < ROBOT_GROUP_CHAT_TABLE_NUM; i++ {
		if err = x.Sync2(&RobotGroupChat{RobotId: int64(i)}); err != nil {
			holmes.Fatal("Fail to sync database: %v", err)
		}
	}
}
