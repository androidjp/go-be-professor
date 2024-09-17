package autorate

import (
	"fmt"
	"sync/atomic"
	"time"
)

type AutoRateKey string

type AutoRateMgr struct {
	autoRateMap map[AutoRateKey]*AutoRateData // 根据key对应的config进行划分，自适应配置
}

func NewAutoRateMgr(cfgMap map[AutoRateKey]AutoRateConfig) *AutoRateMgr {
	mgr := &AutoRateMgr{
		autoRateMap: make(map[AutoRateKey]*AutoRateData),
	}
	for k, v := range cfgMap {
		// 初始化autoRateMap
		mgr.autoRateMap[k] = NewAutoRateData(&v)
	}
	return mgr
}

func NewAutoRateData(cfg *AutoRateConfig) *AutoRateData {
	return &AutoRateData{
		cfg:                cfg,
		fromTime:           time.Now(),
		funcTimeRaiseSumMS: 0,
		funcExecCount:      0,
		curNumber:          cfg.DefaultNumber,
		curRate:            1,
		stopRateReflesh:    make(chan struct{}),
	}
}

func (mgr *AutoRateMgr) Start() {
	// for循环让每一个key都开启计算动作
	for k, v := range mgr.autoRateMap {
		fmt.Printf("%s start... \n", k)
		v.Start()
	}
}

func (mgr *AutoRateMgr) Stop() {
	for k, v := range mgr.autoRateMap {
		fmt.Printf("%s stop... \n", k)
		v.Stop()
	}
}

func (mgr *AutoRateMgr) Statistic(key AutoRateKey) func() {
	startT := time.Now()
	return func() {
		dur := time.Since(startT)
		mgr.TimeRaiseSum(key, dur)
	}
}

func (mgr *AutoRateMgr) TimeRaiseSum(key AutoRateKey, dur time.Duration) {
	// 补充毫秒数
	atomic.AddUint64(&mgr.autoRateMap[key].funcTimeRaiseSumMS, uint64(dur.Milliseconds()))
	// 补充调用次数
	atomic.AddUint64(&mgr.autoRateMap[key].funcExecCount, 1)
}

func (mgr *AutoRateMgr) GetCurRate(key AutoRateKey) float64 {
	if d, ok := mgr.autoRateMap[key]; ok {
		return d.GetCurRate()
	}
	return 1 // 返回默认值为1
}

func (mgr *AutoRateMgr) GetCurNumber(key AutoRateKey) float64 {
	if d, ok := mgr.autoRateMap[key]; ok {
		return d.GetCurNumber()
	}
	return 1 // 返回默认值为1
}

//----------------------------------------------------------
//----------------------------------------------------------

type AutoRateData struct {
	cfg                *AutoRateConfig // 根据key进行划分，自适应配置
	key                AutoRateKey     // key值
	fromTime           time.Time       // 时间窗口左区间
	funcTimeRaiseSumMS uint64          // 时间窗口内累计耗时
	funcExecCount      uint64          // 时间窗口内函数执行次数
	curNumber          float64         // 当前数字
	curRate            float64         // 当前调节比例
	stopRateReflesh    chan struct{}
}

func (d *AutoRateData) Start() {
	// 启动定时计算 curRate的任务
	go func() {
		// 每秒钟尝试更新curRate
		t := time.NewTicker(time.Second)
		for {
			select {
			case <-t.C:
				fmt.Printf("%s curRate: %f \n", d.key, d.curRate)
				if subTime := time.Since(d.fromTime); subTime.Seconds() < float64(d.cfg.RefreshIntervalSec) {
					continue
				}

				// 执行函数
				// GetAutoRateFunc(d.cfg.AutoRateFuncKey)(d)
				TimeRaiseOverflow(d)

				// 重置
				d.fromTime = time.Now()
				d.funcTimeRaiseSumMS = 0
				d.funcExecCount = 0
			case <-d.stopRateReflesh:
				fmt.Printf("结束自动更新autoRate\n")
				return
			}
		}
	}()
}

func (d *AutoRateData) Stop() {
	close(d.stopRateReflesh)
}

func (d *AutoRateData) GetCurRate() float64 {
	return d.curRate
}

func (d *AutoRateData) GetCurNumber() float64 {
	return d.curNumber
}

func TimeRaiseOverflow(d *AutoRateData) {
	if d.funcExecCount == 0 {
		return
	}
	avgFuncTimeRaiseMS := d.funcTimeRaiseSumMS / d.funcExecCount

	if avgFuncTimeRaiseMS < d.cfg.AvgFuncTimeRaiseMinMS {
		fmt.Printf("avgFuncTimeRaiseMS: %d < AvgFuncTimeRaiseMinMS: %d \n", avgFuncTimeRaiseMS, d.cfg.AvgFuncTimeRaiseMinMS)
		// 小于200ms，说明无压力，可以上调。
		d.curRate = 1 + d.cfg.RateStep // 110%
	} else if avgFuncTimeRaiseMS > d.cfg.AvgFuncTimeRaiseMaxMS {
		fmt.Printf("avgFuncTimeRaiseMS: %d > AvgFuncTimeRaiseMaxMS: %d \n", avgFuncTimeRaiseMS, d.cfg.AvgFuncTimeRaiseMaxMS)
		// 大于1000ms，说明有压力，可以下调。
		d.curRate = 1 - d.cfg.RateStep // 90%
	} else {
		d.curRate = 1 // 100%
	}
	fmt.Printf("重新计算，%s curRate: %f \n", d.key, d.curRate)

	// curNumber
	newNum := d.curNumber * d.curRate
	if newNum > d.cfg.MaxNumber {
		newNum = d.cfg.MaxNumber
	}
	if newNum < d.cfg.MinNumber {
		newNum = d.cfg.MinNumber
	}
	d.curNumber = newNum
}
