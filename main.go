package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	flag "github.com/spf13/pflag"
)

var (
	onetime      bool
	syncInterval time.Duration
	inCluster    bool
	version      bool
)

func main() {
	flags := flag.NewFlagSet("kubeps", flag.ExitOnError)

	flags.Usage = func() {
		flags.PrintDefaults()
	}

	flags.BoolVar(&onetime, "onetime", false, "run one time and exit.")
	flags.BoolVar(&inCluster, "in-cluster", true, `If true, use the built in kubernetes cluster for creating the client`)
	flags.BoolVarP(&version, "version", "v", false, "Print version")
	flags.DurationVar(&syncInterval, "sync-interval", (60 * time.Second), "the time duration between template processing.")

	if err := flags.Parse(os.Args[1:]); err != nil {
		log.Fatalln(err)
	}

	if version {
		printVersion()
		os.Exit(0)
	}

	client := NewClient(inCluster)

	if onetime {
		client.process()
		os.Exit(0)
	}

	log.Println("Starting sltd...")

	var wg sync.WaitGroup
	done := make(chan struct{})

	go func() {
		wg.Add(1)
		for {
			client.process()
			log.Printf("Syncing labels complete. Next sync in %v seconds.", syncInterval.Seconds())
			select {
			case <-time.After(syncInterval):
			case <-done:
				wg.Done()
				return
			}
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan
	log.Printf("Shutdown signal received, exiting...")
	close(done)
	wg.Wait()
	os.Exit(0)
}
