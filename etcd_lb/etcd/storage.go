package etcd

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/coreos/etcd/clientv3"
)

var (
	_dialTimeout  int64
	_readTimeout  int64
	_writeTimeout int64

	_single *clientv3.Client // 连接客户端，单例模式
	_lock   *sync.Mutex      // 锁，用于单例模式
)

func init() {
	_dialTimeout = 5
	_readTimeout = 5
	_writeTimeout = 5

	_lock = new(sync.Mutex)
}

// 获取客户端连接
func getEtcdCli(endPoints []string, userName, passWord string) (*clientv3.Client, error) {
	var err error

	// 双重锁校验的单例模式
	if _single == nil {
		_lock.Lock()
		if _single == nil {
			_single, err = clientv3.New(clientv3.Config{
				Endpoints:   endPoints,
				DialTimeout: time.Duration(_dialTimeout) * time.Second,
				Username:    userName,
				Password:    passWord,
			})
		}
		_lock.Unlock()
	}

	return _single, err
}

// 根据前缀批量查询
func GetWithPrefix(endPoints []string, userName, passWord string, prefix string) ([][]byte, error) {
	var err error
	var result = make([][]byte, 0)

	var cli *clientv3.Client
	var cancel context.CancelFunc
	var ctx context.Context
	var resp *clientv3.GetResponse

	// cli
	if cli, err = getEtcdCli(endPoints, userName, passWord); err != nil {
		return result, err
	}

	// get
	ctx, cancel = context.WithTimeout(context.Background(), time.Duration(_readTimeout)*time.Second)
	resp, err = cli.Get(ctx, prefix, clientv3.WithPrefix())
	cancel()

	if err == nil {
		if resp.Count > 0 {
			for _, v := range resp.Kvs {
				result = append(result, v.Value)
			}
		} else {
			err = errors.New("no result")
		}
	}

	return result, err
}

// 获取Grant
func GetGrant(endPoints []string, userName, passWord string, ttl int64) (clientv3.LeaseID, error) {
	var err error
	var result clientv3.LeaseID

	var cli *clientv3.Client
	var cancel context.CancelFunc
	var ctx context.Context
	var resp *clientv3.LeaseGrantResponse

	// cli
	if cli, err = getEtcdCli(endPoints, userName, passWord); err != nil {
		return result, err
	}

	// grant
	ctx, cancel = context.WithTimeout(context.Background(), time.Duration(_readTimeout)*time.Second)
	resp, err = cli.Grant(ctx, ttl)
	cancel()
	if err == nil {
		result = resp.ID
	}

	return result, err
}

// 设置键值对
func PutKeyValue(endPoints []string, userName, passWord string, key, value string, opts ...clientv3.OpOption) error {
	var err error

	var cli *clientv3.Client
	var cancel context.CancelFunc
	var ctx context.Context

	// cli
	if cli, err = getEtcdCli(endPoints, userName, passWord); err != nil {
		return err
	}

	// put
	ctx, cancel = context.WithTimeout(context.Background(), time.Duration(_writeTimeout)*time.Second)
	_, err = cli.Put(ctx, key, value, opts...)
	cancel()

	return err
}

// 根据prefix监听
func WatchWithPrefix(endPoints []string, userName, passWord string, prefix string) (clientv3.WatchChan, error) {
	var err error
	var cli *clientv3.Client
	var watchChan clientv3.WatchChan

	// cli
	if cli, err = getEtcdCli(endPoints, userName, passWord); err != nil {
		return watchChan, err
	}

	// watch
	watchChan = cli.Watch(context.Background(), prefix, clientv3.WithPrefix(), clientv3.WithPrevKV())

	return watchChan, err
}
