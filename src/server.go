package main

import (
	"conf"
	"flag"
	"github.com/golang/glog"
	"github.com/golang/sync/errgroup"
	"net/http"
	"router"
	"time"
)

var g errgroup.Group
func main() {
	flag.Parse()
	defer glog.Flush()

	// 后台管理
	AdminServ := &http.Server{
		Addr:         ":" + conf.GetConfig().AdminPort,
		Handler:      router.GetAdminRouter(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	g.Go(func() error {
		return AdminServ.ListenAndServe()
	})
	// 打印错误日志
	if err := g.Wait(); err != nil {
		glog.Fatal(err)
	}
}