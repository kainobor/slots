package main

import (
    "fmt"
    "github.com/gorilla/mux"
    "net/http"

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

    logger.Log("CONF", "CO", conf.Handler)
    proc := processor.New(conf.Processor)

    h := handler.New(conf, proc)

    r := mux.NewRouter()
    r.HandleFunc(conf.Handler.SpinURL, h.Spins).Methods(http.MethodPost)

    logger.Log("Start to listen")
    err = http.ListenAndServe(fmt.Sprintf(":%d", conf.Handler.Port), r)
    if err != nil {
        logger.Err("can't start handler", "error", err)
    }
}
