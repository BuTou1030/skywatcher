package etcd_lb

import (
	"fmt"
	"log"
	"skywatcher/etcd_lb/etcd"
	"strings"
	"time"

	"go.etcd.io/etcd/client/v3"
)

// Register 非阻塞，在主线程中调用该方法一次，以心跳形式将当前服务地址写入服务注册中心。
// projectName：项目名称，用于树形存储的根路径
// moduleName：模块名称(当前服务名称)，同一项目下包含多个模块，用于树形存储。projectName + moduleName下存储同一模块的一个或多个地址
// addr：当前被注册的服务地址，一般为"ip:port"等
// targetAddr：目标注册中心地址，如etcd数据库、Redis数据库等
// userName：注册中心登录用户名
// passWord：注册中心登录密码
// interval：注册的心跳周期(单位：s)，每间隔该周期会重新注册一次
// lifeCycle：注册信息的生命周期，超过lifeCycle个心跳周期(即lifeCycle * interval)信息即失效
func Register(projectName, moduleName, addr, targetAddr, userName, passWord string, interval, lyfeCycle uint64) error {
	var key = fmt.Sprintf("/%s/%s/%s", projectName, moduleName, addr)

	go func() {
		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		for {
			// 获取LeaseID
			if leaseID, err := etcd.GetGrant(strings.Split(targetAddr, ","), userName, passWord, int64(interval*lyfeCycle)); err != nil {
				log.Printf("Grant failed, %+v\n", err)
			} else {
				if err := etcd.PutKeyValue(strings.Split(targetAddr, ","), userName, passWord, key, addr, clientv3.WithLease(leaseID)); err != nil {
					log.Printf("Put failed, %+v\n", err)
				}
			}

			<-ticker.C
		}
	}()

	return nil
}
