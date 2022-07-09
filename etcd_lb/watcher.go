package etcd_lb

import (
	"fmt"
	"log"
	"skywatcher/etcd_lb/etcd"
	"strings"

	"github.com/coreos/etcd/mvcc/mvccpb"
	"google.golang.org/grpc/naming"
)

type watcher struct {
	resolver      *resolver
	isInitialized bool
}

func (w *watcher) Close() {}

func (w *watcher) Next() ([]*naming.Update, error) {
	var prefix = fmt.Sprintf("/%s/%s/", w.resolver.projectName, w.resolver.moduleName)

	// 首次初始化，提取全部数据
	if !w.isInitialized {
		if resp, err := etcd.GetWithPrefix(strings.Split(w.resolver.targetAddr, ","), w.resolver.userName, w.resolver.passWord, prefix); err != nil {
			log.Printf("GetWithPrefix resp=%+v, err=%+v\n", resp, err)
			return nil, err
		} else {
			var result = make([]*naming.Update, 0)
			for _, ele := range resp {
				result = append(result, &naming.Update{Op: naming.Add, Addr: string(ele)})
			}
			w.isInitialized = true
			return result, nil
		}
	} else {
		if rch, err := etcd.WatchWithPrefix(strings.Split(w.resolver.targetAddr, ","), w.resolver.userName, w.resolver.passWord, prefix); err != nil {
			return nil, err
		} else {
			for resp := range rch {
				var result = make([]*naming.Update, 0)
				for _, event := range resp.Events {
					switch event.Type {
					case mvccpb.PUT:
						if event.PrevKv == nil { // 初次创建键值对
							log.Printf("create key=%s, value=%s\n", event.Kv.Key, event.Kv.Value)
							result = append(result, &naming.Update{Op: naming.Add, Addr: string(event.Kv.Value)})
						} else { // 更新已存在的键值对
							log.Printf("update key=%s\n", event.Kv.Key)
						}
					case mvccpb.DELETE:
						var addr = strings.TrimPrefix(string(event.Kv.Key), prefix) // 剔除前缀，得到addr
						log.Printf("delete key=%s, value=%s\n", event.Kv.Key, addr)
						result = append(result, &naming.Update{Op: naming.Delete, Addr: addr})
					}
				}
				return result, err
			}
			return nil, err
		}
	}
}
