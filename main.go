package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"

	"github.com/kainobor/slots/config"
	"github.com/kainobor/slots/src/handler"
	"github.com/kainobor/slots/src/logger"
	"github.com/kainobor/slots/src/processor"
)

func main() {
	err := logger.Init()
	if err != nil {
		panic("can't start logger: " + err.Error())
	}
	logger.Log("Logger initiated")

	conf, err := config.New()
	if err != nil {
		logger.Err("can't read config", "error", err)
		panic("Can't read config")
	}

	proc := processor.New(conf.Processor)

	h := handler.New(conf, proc)

	r := mux.NewRouter()
	r.HandleFunc(conf.Handler.SpinURL, h.Spins).Methods(http.MethodPost)

	logger.Log("Start to listen")

	// Parallelize handling by ports
	for _, port := range conf.Handler.Ports {
		go func(p int) {
			err = http.ListenAndServe(fmt.Sprintf(":%d", p), r)
			if err != nil {
				logger.Err("can't start handler", "error", err)
			}
		}(port)
	}

	// Wait until end
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	select {
	case s := <-c:
		logger.Log("End with signal %s", s.String())
	}
}
