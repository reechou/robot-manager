package controller

import (
	"sync"
)

type RobotUserBehavior struct {
	sync.Mutex
}
