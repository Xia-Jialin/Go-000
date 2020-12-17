package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello Word!")
}

//1. 基于 errgroup 实现一个 http server 的启动和关闭 ，
// 以及 linux signal 信号的注册和处理，要保证能够一个退出，全部注销退出。
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 监听sig信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func(c chan (os.Signal)) {
		for {
			select {
			case s := <-c:
				switch s {
				case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
					cancel()
				default:
				}
			}
		}
	}(quit)

	g, _ := errgroup.WithContext(ctx)
	// 启动http1
	http1 := http.Server{
		Addr:    ":8080",
		Handler: http.DefaultServeMux,
	}
	http.HandleFunc("/", helloHandler)
	g.Go(func() error {
		if err := http1.ListenAndServe(); err != nil {
			return err
		}
		return nil
	})

	// context取消后，关闭http server
	go func(ctx context.Context) {
		select {
		case <-ctx.Done():
			log.Println("服务退出")
			http1.Shutdown(ctx)
		}
	}(ctx)

	if err := g.Wait(); err != nil {
		log.Println("all exit: ", err)
	}

}
