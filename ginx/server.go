package ginx

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"github.com/uc1024/f90/core/proc"
	"github.com/uc1024/f90/core/slogx"
	"github.com/uc1024/f90/ginx/tracex"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type Server struct {
	Engine *gin.Engine
	http   *http.Server
	conf   *HttpConfig
}

func NewServer(sets ...setConfigOpt) *Server {
	srv := &Server{
		Engine: gin.New(),
		conf:   DefaultHttpConfig()}
	eachFunc := func(each setConfigOpt, i int) { each(srv.conf) }
	lo.ForEach(sets, eachFunc)
	gin.DisableBindValidation()
	setNoExporterTracerProvider(srv.conf.Host)
	srv.Engine.Use(tracex.Trace(srv.conf.Host))
	srv.http = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", srv.conf.Host, srv.conf.Port),
		Handler:      srv.Engine,
		ReadTimeout:  srv.conf.ReadTimeout,
		WriteTimeout: srv.conf.WriteTimeout}
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

func setTracerProvider(endpoint string, name string) *trace.TracerProvider {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint)))
	if err != nil {
		log.Fatal(err.Error())
	}
	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(exp),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(name),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return tp
}

func setNoExporterTracerProvider(name string) *trace.TracerProvider {
	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(name),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return tp
}
