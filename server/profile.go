package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cirocosta/perfer/perf"
	"github.com/gorilla/schema"
)

type profileArgs struct {
	Pid       uint64        `schema:"pid"`
	Frequency uint64        `schema:"freq"`
	Seconds   time.Duration `schema:"seconds"`
}

func HandleProfile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Fprintf(w, "couldn't parse malformed query")
		w.WriteHeader(400)
		return
	}

	args := &profileArgs{}
	err = schema.NewDecoder().Decode(args, r.Form)
	if err != nil {
		fmt.Fprintf(w, "malformed query")
		w.WriteHeader(400)
		return
	}

	execution := perf.NewExecution(args.Frequency, args.Pid, args.Seconds)
	fmt.Println("hey", execution)
}
