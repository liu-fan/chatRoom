package mqtt

import (
	"context"
	"encoding/json"
	"gmqtt/common"
	"gmqtt/config"
	"log"
	"time"
)

var GMergeWorker *Merger

type Merger struct {
	RoutineWorkers []*MergeWorker // 多协程处理
}

//MergeWorker 合并发送
type MergeWorker struct {
	contextChan chan *MergeMsgStruct
	timeoutChan chan *PushBatch
	msgBatch    map[string]*PushBatch //key=topic
}

// PushBatch 根据topic批量发送
type PushBatch struct {
	items       []*common.MsgStruct
	topic       string
	commitTimer *time.Timer
}

// MergeMsgStruct 携带context的push msg
type MergeMsgStruct struct {
	ctx     context.Context
	mqttMsg *common.MsgStruct
}

//InitMerger 初始化
func InitMerger() {
	routineWorker := &Merger{
		RoutineWorkers: make([]*MergeWorker, config.GConfig.MergerWorkerCount),
	}
	for idx := 0; idx < config.GConfig.MergerWorkerCount; idx++ {
		routineWorker.RoutineWorkers[idx] = initMergeWorker()
	}
	GMergeWorker = routineWorker
	return
}
func initMergeWorker() (worker *MergeWorker) {
	worker = &MergeWorker{
		contextChan: make(chan *MergeMsgStruct, config.GConfig.MergerChannelSize),
		timeoutChan: make(chan *PushBatch, config.GConfig.MergerChannelSize),
		msgBatch:    make(map[string]*PushBatch),
	}
	go worker.mergeWorkerMain()
	return worker
}

//HandPushMsg 收集消息分发到不同协程的chan中
func (merger *Merger) HandPushMsg(ctx context.Context, msg *common.MsgStruct) (err error) {
	//计算topic具体分发到哪个协程
	var (
		workerIdx uint32 = 0
		ch        byte
		topic     string
	)
	topic = msg.Topic
	for _, ch = range []byte(topic) {
		workerIdx = (workerIdx + uint32(ch)*33) % uint32(config.GConfig.MergerWorkerCount)
	}
	log.Printf("trace %v:workerId %v", ctx.Value(common.TRACEID), workerIdx)
	return merger.RoutineWorkers[workerIdx].PushMsg(ctx, msg)
}

//PushMsg 把收集到的mqtt msg 发送到chan
func (worker *MergeWorker) PushMsg(ctx context.Context, mqttMsg *common.MsgStruct) (err error) {
	var mergeMsg *MergeMsgStruct
	mergeMsg = &MergeMsgStruct{
		ctx:     ctx,
		mqttMsg: mqttMsg,
	}
	select {
	case worker.contextChan <- mergeMsg:
		MergerPending_INCR()
	default:
		err = common.ERR_MERGE_CHANNEL_FULL
	}
	return
}

// mergeWorkerMain 处理chan中的信息，合并+定时发送
func (worker *MergeWorker) mergeWorkerMain() {
	var (
		msg          *MergeMsgStruct
		batch        *PushBatch
		timeoutBatch *PushBatch
		isCreated    bool
		existed      bool
	)
	for {
		select {
		case msg = <-worker.contextChan:
			MergerPending_DESC()
			log.Printf("trace %v:mergeWorkerMain", msg.ctx.Value(common.TRACEID))
			isCreated = false
			if batch, existed = worker.msgBatch[msg.mqttMsg.Topic]; !existed {
				batch = &PushBatch{topic: msg.mqttMsg.Topic}
				worker.msgBatch[msg.mqttMsg.Topic] = batch
				isCreated = true
			}
			// 合并消息
			batch.items = append(batch.items, msg.mqttMsg)

			// 新建批次, 启动超时自动提交
			if isCreated {
				log.Printf("trace %v:isCreated", msg.ctx.Value(common.TRACEID))
				batch.commitTimer = time.AfterFunc(time.Duration(config.GConfig.MaxMergerDelay)*time.Millisecond, worker.autoCommit(batch))
			}

			// 批次未满, 继续等待下次提交
			if len(batch.items) < config.GConfig.MaxMergerBatchSize {
				continue
			}
		case timeoutBatch = <-worker.timeoutChan:
			// 定时器触发时, 批次已被提交
			if batch, existed = worker.msgBatch[timeoutBatch.topic]; !existed {
				continue
			}

			// 定时器触发时, 前一个批次已提交, 下一个批次已建立
			if batch != timeoutBatch {
				continue
			}

		}
		err := worker.commitBatch(batch)
		MergerRoomTotal_INCR(int64(len(batch.items)))
		if err != nil {
			MergerRoomFail_INCR(int64(len(batch.items)))
		}
	}
}

//把合并消息发送到pushMqtt并删除此消息
func (worker *MergeWorker) commitBatch(batch *PushBatch) (err error) {
	var (
		buf []byte
	)
	if buf, err = json.Marshal(batch.items); err != nil {
		return
	}

	// 打包发送
	delete(worker.msgBatch, batch.topic)
	GConnMgr.PushMqtt(batch.topic, string(buf))
	return
}

//autoCommit 合并的消息定时自动提交到队列
func (worker *MergeWorker) autoCommit(batch *PushBatch) func() {
	return func() {
		worker.timeoutChan <- batch
	}
}
