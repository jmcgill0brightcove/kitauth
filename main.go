package main // import "github.com/jmc-audio/kitauth"

import (
	"fmt"
	stdlog "log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmc-audio/kitauth/bindings"
	"github.com/jmc-audio/kitauth/consts"
	"github.com/jmc-audio/kitauth/log"

	"golang.org/x/net/context"

	kitlog "github.com/go-kit/kit/log"
	levlog "github.com/go-kit/kit/log/levels"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	ctx := context.Background()

	var logger kitlog.Logger
	logger = kitlog.NewLogfmtLogger(os.Stderr)
	logger = kitlog.NewContext(logger).With(log.LogTimestamp, kitlog.DefaultTimestampUTC)
	stdlog.SetOutput(kitlog.NewStdlibAdapter(logger)) // redirect stdlib logging to us
	stdlog.SetFlags(0)

	ctx = context.WithValue(ctx, consts.ContextBaseLogger, logger)
	ctx = context.WithValue(ctx, consts.ContextLogger, levlog.New(logger))

	log.Debug(ctx, "init", "kitsession")

	errc := make(chan error)
	go func() {
		errc <- interrupt()
	}()

	ctx = context.WithValue(ctx, consts.ContextErrorChannel, errc)

	bindings.StartHTTPListener(ctx)

	log.Debug(ctx, "signal", <-errc)
}

func interrupt() error {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	return fmt.Errorf("%s", <-c)
}
