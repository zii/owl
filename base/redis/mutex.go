// redis 实现分布式互斥锁
package redis

import "github.com/go-redsync/redsync"

var syncMap map[string]*redsync.Redsync

func init() {
	syncMap = make(map[string]*redsync.Redsync)
}

func getRedsync(redisName string) *redsync.Redsync {
	s, ok := syncMap[redisName]
	if !ok {
		pool := MustGetPool(redisName)
		s = redsync.New([]redsync.Pool{pool})
		syncMap[redisName] = s
	}
	return s
}

func NewMutex(redisName, key string) *redsync.Mutex {
	s := getRedsync(redisName)
	return s.NewMutex(key)
}
