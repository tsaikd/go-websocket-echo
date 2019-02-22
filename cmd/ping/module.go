package ping

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
	"github.com/tsaikd/KDGoLib/cliutil/cobrather"
	"github.com/tsaikd/KDGoLib/errutil"
	"github.com/tsaikd/go-websocket-echo/logger"
	"golang.org/x/sync/errgroup"
)

var flagURL = &cobrather.StringFlag{
	Name:    "ping.url",
	Default: "ws://localhost:8080",
	Usage:   "websocket server URL",
}

var flagMessage = &cobrather.StringFlag{
	Name:    "ping.message",
	Default: "",
	Usage:   "ping message",
}

var flagNumber = &cobrather.Int64Flag{
	Name:    "ping.number",
	Default: 1,
	Usage:   "ping concurrent threads",
}

var flagKeep = &cobrather.BoolFlag{
	Name:    "ping.keep",
	Default: false,
	Usage:   "keep ping connection until user interrupt",
}

// Module info of package
var Module = &cobrather.Module{
	Use:   "ping",
	Short: "ping grpc echo server",
	Flags: []cobrather.Flag{
		flagURL,
		flagMessage,
		flagNumber,
		flagKeep,
	},
	RunE: func(ctx context.Context, cmd *cobra.Command, args []string) (err error) {
		url := flagURL.String()
		message := flagMessage.String()
		number := flagNumber.Int64()
		keep := flagKeep.Bool()
		ctx = catchInterruptSignal(ctx)

		eg, ctx := errgroup.WithContext(ctx)
		for i := int64(0); i < number; i++ {
			eg.Go(func() error {
				return ping(ctx, url, message, nil, keep)
			})
		}
		return eg.Wait()
	},
}

func catchInterruptSignal(ctx context.Context) context.Context {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	cancelCtx, cancel := context.WithCancel(ctx)
	go func() {
		select {
		case <-interrupt:
			cancel()
		}
	}()
	return cancelCtx
}

var connectionCount int64

func ping(ctx context.Context, url string, message string, header http.Header, keep bool) (err error) {
	logger := logger.Logger()
	dialer := websocket.Dialer{
		Proxy: http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	conn, _, err := dialer.Dial(url, header)
	if err != nil {
		return
	}
	defer func() {
		errutil.Trace(conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")))
		errutil.Trace(conn.Close())
	}()

	atomic.AddInt64(&connectionCount, 1)
	defer atomic.AddInt64(&connectionCount, -1)

	for {
		logger.Printf("(%d) connected to %q", connectionCount, url)
		timestamp := time.Now().Format(time.RFC3339Nano)
		if err = conn.WriteMessage(websocket.TextMessage, []byte(timestamp+": "+message)); err != nil {
			return
		}
		if _, _, err = conn.ReadMessage(); err != nil {
			return
		}
		if !keep {
			return
		}
		select {
		case <-ctx.Done():
			return
		}
	}
}
