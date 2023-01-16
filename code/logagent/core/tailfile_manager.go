package core

import (
	"code/logagent/common"
	"github.com/sirupsen/logrus"
)

// 管理 tailTask 的对象

type tailTaskMagaer struct {
	// 使用map
	tailTaskMap      map[string]*tailTask
	collectEntryList []common.CollectEntry
	confChan         chan []common.CollectEntry
}

var (
	ttMgr *tailTaskMagaer
)

// InitTail 初始化tailTask 为每一个日志文件造一个单独的tailTask
// main 函数中调用
func InitTail(allConfig []common.CollectEntry) (err error) {

	// 初始化管理对象 使用单例
	ttMgr = &tailTaskMagaer{
		tailTaskMap:      make(map[string]*tailTask, 20),   // 所有的tailTask任务
		collectEntryList: allConfig,                        // 所有配置项
		confChan:         make(chan []common.CollectEntry), // 等待新配制的通道 做一个阻塞channel
	}

	// allConfig 中存了若干个日志的收集项
	// 针对每一个日志收集项目创建一个对应的tailObj对象

	for _, conf := range allConfig {
		// 创建一个日志收集任务
		tt := newTailTask(conf.Path, conf.Topic)
		// 初始化
		err = tt.Init()
		if err != nil {
			logrus.Errorf("create tailObj for path:%s failed, err:%v\n", tt.path, err)
			continue
		}
		logrus.Infof("create a tail task for path:%s success", conf.Path)

		// 将创建的这个tailTask任务登记到tailTaskManager 方便后续管理
		ttMgr.tailTaskMap[tt.path] = tt
		// 启动一个goroutine 去收集日志
		go tt.run()
	}

	// 在后台等待新配置
	go ttMgr.watch()

	return
}

func (ttm *tailTaskMagaer) watch() {
	// 初始化新配置管道
	// 派小弟等着新配置来  如果能取到值，说明有新的配置来了
	newConf := <-ttm.confChan
	logrus.Infof("get new conf from etcd, conf:%v, start manager tailTask...", newConf)

	// 新配置来了之后应该管理一下之前启动的那些tailTask
	for _, conf := range newConf {
		// 1. 原来已经存在的任务不用动
		if ttm.isExist(conf) {
			continue
		}
		// 2. 原来没有的创建一个新的任务
		tt := newTailTask(conf.Path, conf.Topic)
		err := tt.Init()
		if err != nil {
			logrus.Errorf("create tailObj for path:%s failed, err:%v\n", tt.path, err)
			continue
		}
		logrus.Infof("create a tail task for path:%s success", conf.Path)
		ttm.tailTaskMap[tt.path] = tt
		go tt.run()
	}

	// 3. 原来有的现在没有要停掉
	// 找出tailTaskMap中存在，但是newConf中不存在的那些tailTask，把他们都停掉
	for key, task := range ttm.tailTaskMap {
		var found bool
		for _, conf := range newConf {
			if key == conf.Path {
				// 如果有，返回，什么都不做
				found = true
				break
			}
		}
		if !found {
			logrus.Infof("the task collect path:%s need to stop.", task.path)
			// 新配置中没有,要停掉tailTask 停止go routine 通过context中WithCancel停止
			delete(ttm.tailTaskMap, key) // 从管理类中删除
			task.cancel()
		}
	}
}

// isExist 判断tailTaskMap中是否存在该收集项
func (ttm *tailTaskMagaer) isExist(conf common.CollectEntry) bool {
	_, ok := ttm.tailTaskMap[conf.Path]
	return ok
}

// SendNewConf 暴露出confChan
func SendNewConf(newConf []common.CollectEntry) {
	ttMgr.confChan <- newConf
}
