package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	cmn "github.com/sunao-uehara/go-restapi-sample/common"
	handler "github.com/sunao-uehara/go-restapi-sample/handlers"
	r "github.com/sunao-uehara/go-restapi-sample/router"
	"github.com/sunao-uehara/go-restapi-sample/storages/mysql"
	"go.uber.org/zap"
)

func main() {
	// utilize multicore CPUs. enable this if go version is under 1.5
	// runtime.GOMAXPROCS(runtime.NumCPU())

	// initialize logger
	// logger, _ := zap.NewProduction()
	logger, _ := zap.NewDevelopment()
	defer logger.Sync() // flushes buffer, if any
	log := logger.Sugar()

	// initialize WaitGroup
	wg := &sync.WaitGroup{}

	// initialize mysql
	db, err := mysql.Initialize(os.Getenv(cmn.MYSQL_URL))
	if err != nil {
		log.Fatal("unable to initialize mysql", err)
	}
	defer db.Close()

	// initialize redis

	// initialize handler
	h := handler.NewHandler(&handler.HandlerOptions{
		Log:   log,
		Wg:    wg,
		Mysql: db,
	})

	srv := &http.Server{
		Addr:    ":" + os.Getenv(cmn.PORT),
		Handler: r.NewRouter(h),
	}
	// start up http server
	go func() {
		log.Info("listen and serve")
		if err := srv.ListenAndServe(); err != nil {
			log.Info(err.Error())
		}
	}()

	// wait for SIGTERM signal
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	s := <-sig
	log.Infof("signal %s received\n", s)
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Info("shut down the server gracefully...")
	if err := srv.Shutdown(ctxTimeout); err != nil {
		log.Fatalf("could not gracefully shut down server, %s", err.Error())
	}

	log.Info("waiting all goroutines are finished...")
	wg.Wait()
	log.Info("all done, really closing")
}
