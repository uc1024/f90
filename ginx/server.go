package ginx

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"github.com/uc1024/f90/core/proc"
	"github.com/uc1024/f90/core/slogx"
	"github.com/uc1024/f90/ginx/trace"
)

type Server struct {
	engine *gin.Engine
	http   *http.Server
	conf   *HttpConfig
}

func NewServer(sets ...setConfigOpt) *Server {
	srv := &Server{
		engine: gin.New(),
		conf:   DefaultHttpConfig()}
	eachFunc := func(each setConfigOpt, i int) { each(srv.conf) }
	lo.ForEach(sets, eachFunc)
	srv.http = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", srv.conf.Host, srv.conf.Port),
		Handler:      srv.engine,
		ReadTimeout:  srv.conf.ReadTimeout,
		WriteTimeout: srv.conf.WriteTimeout}
	gin.DisableBindValidation()
	srv.engine.Use(trace.SetHeaderTrace)
	return srv
}

func (srv *Server) Run() error {
	proc.AddShutdownListener(func() {
		if err := srv.http.Shutdown(context.Background()); err != nil {
			slogx.Default.Error(context.Background(), err.Error())
		}
	})
	return srv.http.ListenAndServe()
}
