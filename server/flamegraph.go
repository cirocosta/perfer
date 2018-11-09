package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/cirocosta/perfer/perf"
	"github.com/rs/xid"
)

func (s *Server) HandleFlamegraph(w http.ResponseWriter, r *http.Request) {
	args, err := ProfileArgsFromRequest(r)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintln(w, err)
		return
	}

	execution, err := perf.NewExecution(args.Frequency, args.Pid, args.Seconds)
	if err != nil {
		log.Printf("%+v\n", err)
		w.WriteHeader(500)
		fmt.Fprintln(w, "failed to create execution from args")
		return
	}

	var (
		uid               = xid.New().String()
		outputProfileFile = path.Join(s.assetsDirectory, uid)
		ctx               = context.Background()
	)

	err = execution.Record(context.Background(), outputProfileFile)
	if err != nil {
		log.Printf("%+v\n", err)
		w.WriteHeader(500)
		fmt.Fprintln(w, "failed to record profile from args")
		return
	}

	flamegraphFile, err := os.Create(outputProfileFile + ".flamegraph")
	if err != nil {
		log.Printf("%+v\n", err)
		w.WriteHeader(500)
		fmt.Fprintln(w, "failed to create flamegraph file")
		return
	}

	defer flamegraphFile.Close()

	err = execution.GenerateFlamegraph(ctx, outputProfileFile, flamegraphFile)
	if err != nil {
		log.Printf("%+v\n", err)
		w.WriteHeader(500)
		fmt.Fprintln(w, "failed to generate flamegraph")
		return
	}

	fmt.Fprintln(w, flamegraphFile.Name())
}
