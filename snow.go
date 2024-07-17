package jxsnow

import (
	"fmt"
	"sync"
	"time"
)

type Generator struct {
	machineID, sequence, timestamp int64
	lock                           sync.Mutex
}

const (
	machineBits   = 10
	sequenceBits  = 12
	timestampBits = 41

	maxMachineID   = 1<<machineBits - 1
	maxSequenceNum = 1<<sequenceBits - 1
	maxTimestamp   = 1<<timestampBits - 1

	timestampShiftBits = machineBits + sequenceBits
	machineIDShiftBits = sequenceBits
)

// NewGenerator 创建一个生成 ID 对象。每个节点的 machineID 必须不同。
func NewGenerator(machineID int64) (*Generator, error) {
	if machineID > maxMachineID {
		return nil, fmt.Errorf("%s", "machine code is too large")
	}
	return &Generator{
		machineID: machineID,
		timestamp: time.Now().UnixMilli(),
	}, nil
}

// Generate 雪花算法生成 ID。41 位毫秒时间戳，10 位工作机器 ID，12 位序列号。
func (g *Generator) Generate() (int64, error) {
	g.lock.Lock()
	defer g.lock.Unlock()

	t := time.Now().UnixMilli()
	if t < g.timestamp {
		t = g.nextTime(g.timestamp)
	}

	if t == g.timestamp {
		g.sequence = (g.sequence + 1) & maxSequenceNum
		if g.sequence == 0 {
			g.timestamp = g.nextTime(t + 1)
		}
	} else {
		g.sequence = 0
		g.timestamp = t
	}

	if g.timestamp > maxTimestamp {
		return 0, fmt.Errorf("%s", "time has exceeded its limit")
	}

	return g.timestamp<<timestampShiftBits | g.machineID<<machineIDShiftBits | g.sequence, nil
}

func (g *Generator) nextTime(timestamp int64) int64 {
	/*for {
		t := time.Now().UnixMilli()
		if t < timestamp {
			time.Sleep(time.Millisecond)
			continue
		}
		return t
	}*/
	for {
		t := time.Now().UnixMilli()
		if t >= timestamp {
			return t
		}

		// 使用指数退避策略避免频繁调用
		time.Sleep(time.Duration(1<<(timestamp-t)) * time.Microsecond)
	}
}
