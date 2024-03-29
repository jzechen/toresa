/**
 * @Time: 2023/10/20 16:13
 * @Author: jzechen
 * @File: server.go
 * @Software: GoLand collector
 */

package server

import (
	"context"
	"github.com/jzechen/toresa/pkg/common/apiserver"
	"github.com/jzechen/toresa/pkg/manager/config"
	"github.com/jzechen/toresa/pkg/manager/mdb"
	"github.com/jzechen/toresa/pkg/manager/server/engine"
	"github.com/jzechen/toresa/pkg/manager/server/engine/gin"
	"github.com/jzechen/toresa/pkg/manager/server/rate"
	"github.com/jzechen/toresa/pkg/manager/services/sina"
	"k8s.io/klog/v2"
	"net"
	"net/http"
	"strconv"
	"time"
)

type ManagerServer struct {
	ctx      context.Context
	cfg      *config.CollectorManager
	server   *http.Server
	listener net.Listener
	engine   engine.Interface
	mgo      mdb.Interface
}

func NewCollectorManagerServer(ctx context.Context, cfg *config.CollectorManager) (*ManagerServer, error) {
	klog.V(4).Info("initialize rate limiter")
	rate.InitRateLimiter(&cfg.Server)

	klog.V(4).Info("initialize server listener")
	addr := net.JoinHostPort(cfg.Server.Addr, strconv.Itoa(cfg.Server.Port))
	ln, _, err := apiserver.CreateListener("tcp", addr, net.ListenConfig{})
	if err != nil {
		return nil, err
	}

	klog.V(4).Info("initialize the mongoDB client")
	mgo, err := mdb.NewMongoDBImpl(ctx, &cfg.Mongo)
	if err != nil {
		return nil, err
	}

	klog.V(4).Info("register the sina handler")
	sinaHandler := sina.NewSinaHandler(cfg, mgo)

	klog.V(4).Info("generate a web engine")
	ginEngine := gin.NewGinEngine(sinaHandler)

	klog.V(4).Info("register http server services")
	handler := ginEngine.CreateHandler()
	srv := &http.Server{
		Addr:           ln.Addr().String(),
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	//klog.V(2).Info("create MySQL data if tables not exist")
	//err = CreateTablesIfNotExist()
	//if err != nil {
	//	return nil, err
	//}

	return &ManagerServer{
		ctx:      ctx,
		cfg:      cfg,
		server:   srv,
		listener: ln,
		engine:   ginEngine,
		mgo:      mgo,
	}, nil
}

func (rms *ManagerServer) Run() {
	// run api server
	shutdownCh, _, err := apiserver.RunServer(rms.server, rms.listener, rms.cfg.Server.RequestTimeout, rms.ctx.Done())
	if err != nil {
		panic(err)
	}
	<-shutdownCh

	// shutdown gracefully
	rms.close()
}

func (rms *ManagerServer) close() {
	// do some close operation before receives exit signals
	_ = rms.mgo.Close()
}

func CreateTablesIfNotExist() error {
	// create databases tables if needed
	return nil
}
