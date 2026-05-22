package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/example/pulsectl/internal/checker"
	"github.com/example/pulsectl/internal/config"
	"github.com/example/pulsectl/internal/reporter"
	"github.com/example/pulsectl/internal/scheduler"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading config: %v\n", err)
		os.Exit(1)
	}

	chk := checker.New(cfg)
	rep := reporter.New(cfg)
	sched := scheduler.New(cfg, chk)

	results := sched.Start()

	go rep.Consume(results)

	sched.Wait()
	rep.PrintSummary()
}
