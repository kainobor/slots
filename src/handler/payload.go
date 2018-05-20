package handler

import (
    "fmt"
    "github.com/dgrijalva/jwt-go"
    "github.com/pkg/errors"
    "io/ioutil"
    "net/http"
)

type PayloadClaims struct {
    Uid   string `json:"uid"`
    Chips int64  `json:"chips"`
    Bet   int64  `json:"bet"`
    jwt.StandardClaims
}

func (p *PayloadClaims) Valid() error {
    if p.Bet > p.Chips {
        return errors.New("Bet is too big")
    }

    return nil
}

func newPayload(r *http.Request, secret []byte) (*PayloadClaims, error) {
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        return nil, fmt.Errorf("can't read body: %#v", err)
    }

    token, err := jwt.ParseWithClaims(
        string(body),
        &PayloadClaims{},
        func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }

            return secret, nil
        })
    if err != nil {
        return nil, fmt.Errorf("errors while try to parse payload: %v", err)
    }

    pld, ok := token.Claims.(*PayloadClaims)
    if !ok || !token.Valid {
        return nil, fmt.Errorf("not valid token: %#v", token)
    }

    return pld, nil
}
