package server

import (
	"context"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
	"github.com/tsaikd/KDGoLib/cliutil/cobrather"
	"github.com/tsaikd/KDGoLib/errutil"
	"github.com/tsaikd/go-websocket-echo/logger"
)

var flagAddr = &cobrather.StringFlag{
	Name:    "server.addr",
	Default: ":8080",
	Usage:   "websocket echo server listen port",
}

// Module info of package
var Module = &cobrather.Module{
	Use:   "server",
	Short: "websocket echo server",
	Flags: []cobrather.Flag{
		flagAddr,
	},
	RunE: func(ctx context.Context, cmd *cobra.Command, args []string) (err error) {
		addr := flagAddr.String()
		go connectionStatistics(ctx, 5*time.Second)
		return listen(addr)
	},
}

var upgrader = websocket.Upgrader{}

func listen(addr string) (err error) {
	logger := logger.Logger()
	http.HandleFunc("/", echo)
	logger.Printf("Listen websocket: %q", addr)
	return http.ListenAndServe(addr, nil)
}

var connectionCount int64

func echo(w http.ResponseWriter, r *http.Request) {
	logger := logger.Logger()
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Print("upgrade:", err)
		return
	}
	defer func() {
		errutil.Trace(conn.Close())
	}()

	atomic.AddInt64(&connectionCount, 1)
	defer atomic.AddInt64(&connectionCount, -1)

	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			switch e := err.(type) {
			case *websocket.CloseError:
				if e.Code == 1000 {
					return
				}
			}
			logger.Println("read:", err)
			break
		}
		logger.Printf("recv: %s", message)
		err = conn.WriteMessage(mt, message)
		if err != nil {
			logger.Println("write:", err)
			break
		}
	}
}

func connectionStatistics(ctx context.Context, duration time.Duration) {
	logger := logger.Logger()
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			logger.Printf("Connection: %d", connectionCount)
		}
	}
}
