package handler

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    "github.com/dgrijalva/jwt-go"
    "github.com/kainobor/slots/config"
    "github.com/kainobor/slots/src/logger"
    "github.com/kainobor/slots/src/processor"
)

type Handler struct {
    config *config.Config
    proc   *processor.Processor
    secret []byte
}

type SpinResponse struct {
    Total int64                   `json:"total"`
    Spins []*processor.SpinResult `json:"spins"`
    JWT   string                  `json:"jwt"`
}

const (
    SpinTypeMain = "main"
    SpinTypeFree = "free"
)

func New(conf *config.Config, proc *processor.Processor) *Handler {
    h := new(Handler)
    h.config = conf
    h.proc = proc
    h.secret = []byte(conf.Handler.Secret)

    return h
}

func (h *Handler) Spins(w http.ResponseWriter, r *http.Request) {
    start := time.Now().UnixNano()

    pld, err := newPayload(r, h.secret)
    if err != nil {
        t := time.Now().UnixNano()
        logger.Err("error while getting payload", "time", t, "error", err)
        http.Error(w, fmt.Sprintf("Internal error at %d", t), 500)
        return
    }

    newBalance := pld.Chips - pld.Bet*int64(len(h.config.Processor.Lines))
    if newBalance < 0 {
        http.Error(w, fmt.Sprintf("You're too poor"), 500)
        return
    }

    resp := new(SpinResponse)

    for i := 0; i <= 0; i++ {
        outcome := h.proc.GenerateOutcome()
        spin, err := h.proc.CheckWins(outcome, pld.Bet)
        if err != nil {
            t := time.Now().UnixNano()
            logger.Err("error while checking wins", "time", t, "error", err, "outcome", outcome)
            http.Error(w, fmt.Sprintf("Internal error at %d", time.Now().UnixNano()), 500)
            return
        }

        spin.Type = SpinTypeMain
        if spin.IsFree {
            spin.Type = SpinTypeFree
            i = i - h.config.Handler.FreeRollsByOneAddiction
        }

        resp.Total = resp.Total + spin.Total
        resp.Spins = append(resp.Spins, spin)

        logger.Log("Result", "outcome", outcome)
    }

    pld.Chips = newBalance + resp.Total
    token := jwt.NewWithClaims(jwt.SigningMethodHS512, pld)
    resp.JWT, err = token.SignedString([]byte(h.config.Handler.Secret))
    if err != nil {
        t := time.Now().UnixNano()
        logger.Err("error while creating token", "time", t, "error", err, "response", resp)
        http.Error(w, fmt.Sprintf("Internal error at %d", time.Now().UnixNano()), 500)
        return
    }

    json.NewEncoder(w).Encode(resp)

    logger.Log("Time", "duration", float64(time.Now().UnixNano()-start)/1000000000)
}
