package main

import (
	"context"
	"math/rand"
	"os"
	"runtime"
	"transactions/task"

	randomtransaction "transactions/task/worker/randomTransaction"

	"github.com/jessevdk/go-flags"
)

var opts struct {
	LowerLimitMinValue float32 `long:"llmin" description:"Lower limit min value" required:"true"`
	LowerLimitMaxValue float32 `long:"llmax" description:"Lower limit max value" required:"true"`
	UpperLimitMinValue float32 `long:"ulmin" description:"Upper limit min value" required:"true"`
	UpperLimitMaxValue float32 `long:"ulmax" description:"Upper limit max value" required:"true"`
	Percentage         float32 `short:"p" long:"percentage" description:"Percentage for lower limit" required:"true"`
	TotalLines         int     `short:"t" long:"totallines" description:"Total lines" required:"true"`
	File               string  `short:"f" long:"file" description:"Output file" required:"true"`
}

func run(args []string) error {
	flags.ParseArgs(&opts, args)
	ctx := context.Background()
	maxGoRoutines := runtime.GOMAXPROCS(0)
	pool := task.New(ctx, maxGoRoutines)
	lowerLimit := float32(opts.TotalLines) * opts.Percentage
	remaining := float32(opts.TotalLines) - lowerLimit
	workers := make([]task.Worker, opts.TotalLines)
	for i := 0; i < int(lowerLimit); i++ {
		workers[i] = &randomtransaction.Worker{FilePath: opts.File, MinAmount: opts.LowerLimitMinValue, MaxAmount: opts.LowerLimitMaxValue}
	}
	for i := int(remaining); i < opts.TotalLines; i++ {
		workers[i] = &randomtransaction.Worker{FilePath: opts.File, MinAmount: opts.UpperLimitMinValue, MaxAmount: opts.UpperLimitMaxValue}
	}
	rand.Shuffle(len(workers), func(i, j int) { workers[i], workers[j] = workers[j], workers[i] })
	for _, w := range workers {
		pool.Do(w)
	}
	pool.ShutDown()
	return nil
}

func main() {
	run(os.Args)
}
