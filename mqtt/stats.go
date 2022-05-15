package mqtt

import (
	"gmqtt/config"
	"sync/atomic"
)

type Stats struct {
	MergerPending   int64 `json:"mergerPending"`
	MergerRoomTotal int64 `json:"mergerRoomTotal"`
	MergerRoomFail  int64 `json:"mergerRoomFail"`
}

var (
	GStats *Stats
)

func InitStats() (err error) {
	GStats = &Stats{}
	return
}

func MergerPending_INCR() {
	num := atomic.AddInt64(&GStats.MergerPending, 1)
	if num >= config.GConfig.MergerPendingWarnNum {
		warnProFun.DingRobot("MergerPending 预警")
	}
}

func MergerPending_DESC() {
	atomic.AddInt64(&GStats.MergerPending, -1)
}

func MergerRoomTotal_INCR(batchSize int64) {
	atomic.AddInt64(&GStats.MergerRoomTotal, batchSize)
}

func MergerRoomFail_INCR(batchSize int64) {
	atomic.AddInt64(&GStats.MergerRoomFail, batchSize)
}
