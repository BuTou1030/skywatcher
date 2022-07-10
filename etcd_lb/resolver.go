package etcd_lb

import (
	"fmt"

	"google.golang.org/grpc/naming"
)

type Resolver struct {
	projectName string
	moduleName  string
	targetAddr  string
	userName    string
	passWord    string
}

// projectName：项目名称，用于树形存储的根路径
// moduleName：模块名称(当前服务名称)，同一项目下包含多个模块，用于树形存储。projectName + moduleName下存储同一模块的一个或多个地址
// targetAddr：目标注册中心地址，如etcd数据库、Redis数据库等
// userName：注册中心登录用户名
// passWord：注册中心登录密码
func NewResolver(projectName, moduleName, userName, passWord string) *Resolver {
	return &Resolver{
		projectName: projectName,
		moduleName:  moduleName,
		userName:    userName,
		passWord:    passWord,
	}
}

func (re *Resolver) Resolve(target string) (naming.Watcher, error) {
	if re.projectName == "" || re.moduleName == "" {
		return nil, fmt.Errorf("projectName or moduleName was nil")
	}
	re.targetAddr = target

	return &watcher{resolver: re, isInitialized: false}, nil
}
