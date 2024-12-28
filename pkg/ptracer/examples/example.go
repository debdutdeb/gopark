package main

import (
	"context"
	"flag"
	"fmt"
	"syscall"

	"github.com/debdutdeb/gopark/pkg/ptracer"
)

func main() {
	pid := flag.Int("pid", -1, "--pid [pid]")

	flag.Parse()

	controller, err := ptracer.Attach(*pid)
	if err != nil {
		panic(err)
	}

	controller.Trace(context.Background(), syscall.SYS_WRITE)

	for {
		call, err := controller.Syscall()
		if err != nil {
			panic(err)
		}

		write, ok := call.(*ptracer.WriteSyscall)
		if ok {
			content, err := write.Content()
			if err != nil {
				panic(err)
			}

			fmt.Print(content)
		}
	}
}
