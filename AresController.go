// Package ares implements a signaling server based on WebSocket.
package ares

import (
	"crypto/tls"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"strconv"
	"time"
)

const registerTimeoutSec = 10

// This is a temporary solution to avoid holding a zombie connection forever, by
// setting a 1 day timeout on reading from WebSocket connection
const wsReadTimeoutSec = 60 * 60 * 24

type AresController struct {
}

func NewAresController() *AresController {
	return &AresController{}
}

// Run starts the AresController server and blocks the thread until the program exits.
func (p *AresController) Run(port int, useTls bool) {
	http.Handle("/ws", websocket.Handler(p.wsHandler))

	var e error

	pstr := ":" + strconv.Itoa(port)
	if useTls {
		config := &tls.Config{
			// Only allow ciphers that support forward secrecy for iOS9 compatibility

			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			},
			PreferServerCipherSuites: true,
		}
		server := &http.Server{Addr: pstr, Handler: nil, TLSConfig: config}

		e = server.ListenAndServeTLS("/cert/cert.pem", "/cert/key.pem")
	} else {
		e = http.ListenAndServe(pstr, nil)
	}

	if e != nil {
		log.Fatal("Run: " + e.Error())
	}
}

// wsHandler is a WebSocket server that handles requests from the WebSocket client in the form of:
// 1.
// 2.
// Unexpected messages will cause the WebSocket connection to be closed.
func (p *AresController) wsHandler(ws *websocket.Conn) {
	var rid, cid string

	var msg wsClientMsg
loop:
	for {
		err := ws.SetReadDeadline(time.Now().Add(time.Duration(wsReadTimeoutSec) * time.Second))
		if err != nil {
			p.wsError("ws.SetReadDeadline error: "+err.Error(), ws)
			break
		}

		err = websocket.JSON.Receive(ws, &msg)
		if err != nil {
			if err.Error() != "EOF" {
				p.wsError("websocket.JSON.Receive error: "+err.Error(), ws)
			}
			break
		}
	}

}
