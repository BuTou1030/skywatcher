package etcd_lb

import (
	"fmt"
	"google.golang.org/grpc/resolver"
)

const schema = "etcdv3_resolver"

type Resolver struct {
	projectName string
	moduleName  string
	userName    string
	passWord    string
	targetAddr  string
}

// projectName：项目名称，用于树形存储的根路径
// moduleName：模块名称(当前服务名称)，同一项目下包含多个模块，用于树形存储。projectName + moduleName下存储同一模块的一个或多个地址
// targetAddr：目标注册中心地址，如etcd数据库、Redis数据库等
// userName：注册中心登录用户名
// passWord：注册中心登录密码
func NewResolver(projectName, moduleName, userName, passWord, targetAddr string) *Resolver {
	return &Resolver{
		projectName: projectName,
		moduleName:  moduleName,
		userName:    userName,
		passWord:    passWord,
		targetAddr:  targetAddr,
	}
}

// Scheme return etcdv3 schema
func (re *Resolver) Scheme() string {
	return schema
}

func (re *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	var err error
	if re.projectName == "" || re.moduleName == "" {
		err = fmt.Errorf("projectName or moduleName was nil")
		return nil, err
	}

	var w = NewWatcher(re, false, cc)
	go w.watch()
	
	return re, nil
}

// ResolveNow
func (re *Resolver) ResolveNow(rn resolver.ResolveNowOptions) {
}

// Close
func (re *Resolver) Close() {
}

