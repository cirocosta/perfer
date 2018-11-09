package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path"

	"github.com/cirocosta/perfer/perf"
	"github.com/gorilla/schema"
	"github.com/pkg/errors"
	"github.com/rs/xid"
)

type profileArgs struct {
	Pid       uint64 `schema:"pid"`
	Frequency uint64 `schema:"freq"`
	Seconds   uint64 `schema:"seconds"`
}

func ProfileArgsFromRequest(r *http.Request) (args *profileArgs, err error) {
	args = &profileArgs{}

	err = r.ParseForm()
	if err != nil {
		err = errors.Wrapf(err, "couldn't parse malformed query")
		return
	}

	err = schema.NewDecoder().Decode(args, r.Form)
	if err != nil {
		err = errors.Wrapf(err, "failed to decode query")
		return
	}

	if args.Seconds == 0 {
		err = errors.Wrapf(err, "sampling duration kmust be non-zero")
		return
	}

	if args.Frequency == 0 {
		err = errors.Wrapf(err, "frequency must be non-zero")
		return
	}

	if args.Pid == 0 {
		err = errors.Wrapf(err, "pid must be non-zero")
		return
	}

	return
}

func (s *Server) HandleProfile(w http.ResponseWriter, r *http.Request) {
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
		uid        = xid.New().String()
		outputFile = path.Join(s.assetsDirectory, uid)
	)
	err = execution.Record(context.Background(), outputFile)
	if err != nil {
		log.Printf("%+v\n", err)
		w.WriteHeader(500)
		fmt.Fprintln(w, "failed to record profile from args")
		return
	}

	fmt.Fprintln(w, uid)
}
