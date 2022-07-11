package etcd_lb

import (
	"fmt"
	"log"
	"time"
	"skywatcher/etcd_lb/etcd"
	"strings"

	"go.etcd.io/etcd/api/v3/mvccpb"
	"google.golang.org/grpc/resolver"
)

func NewWatcher(resolver *Resolver, isInitialized bool, cc resolver.ClientConn) *watcher  {
	return &watcher{
		resolver: resolver, 
		isInitialized: false,
		cc: cc,
	}
}

type watcher struct {
	resolver      *Resolver
	isInitialized bool
	cc resolver.ClientConn
}

func (w *watcher) Close() {}

func (w *watcher) watch() {
	var prefix = fmt.Sprintf("/%s/%s/", w.resolver.projectName, w.resolver.moduleName)
	var addrDict = make(map[string]resolver.Address)

	update := func() {
		addrList := make([]resolver.Address, 0, len(addrDict))
		for _, v := range addrDict {
			addrList = append(addrList, v)
		}
		w.cc.UpdateState(resolver.State{Addresses: addrList})
	}
	ticker := time.NewTicker(1 * time.Second)
	for _ = range ticker.C {
		// 首次初始化，提取全部数据
		if !w.isInitialized {
			if resp, err := etcd.GetWithPrefix(strings.Split(w.resolver.targetAddr, ","), w.resolver.userName, w.resolver.passWord, prefix); err != nil {
				log.Printf("GetWithPrefix resp=%+v, err=%+v\n", resp, err)
			} else {
				for _, ele := range resp {
					addrDict[string(ele)] = resolver.Address{Addr: string(ele)}
				}
				w.isInitialized = true
			}
		} else {
			if rch, err := etcd.WatchWithPrefix(strings.Split(w.resolver.targetAddr, ","), w.resolver.userName, w.resolver.passWord, prefix); err != nil {
				log.Printf("GetWithPrefix rch=%+v, err=%+v\n", rch, err)
			} else {
				for resp := range rch {
					for _, event := range resp.Events {
						switch event.Type {
						case mvccpb.PUT:
							if event.PrevKv == nil { // 初次创建键值对
								log.Printf("create key=%s, value=%s\n", event.Kv.Key, event.Kv.Value)
								addrDict[string(event.Kv.Key)] = resolver.Address{Addr: string(event.Kv.Value)}
							} else { // 更新已存在的键值对
								log.Printf("update key=%s\n", event.Kv.Key)
							}
						case mvccpb.DELETE:
							var addr = strings.TrimPrefix(string(event.Kv.Key), prefix) // 剔除前缀，得到addr
							log.Printf("delete key=%s, value=%s\n", event.Kv.Key, addr)
							delete(addrDict, string(event.PrevKv.Key))
						}
					}
				}
			}
		}
		update()
	}
}


