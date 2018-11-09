package perf

import (
	"context"
	"os/exec"
	"path"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/xid"
)

// Execution provides the context arguments for the execution of the
// the perf and flamegraph scripts.
type execution struct {
	Uid                  string
	Frequency            uint64
	Pid                  uint64
	DestinationDirectory string
	SamplingDuration     time.Duration
}

func NewExecution(frequency, pid uint64, duration time.Duration) *execution {
	return &execution{
		Uid:                  xid.New().String(),
		Frequency:            frequency,
		Pid:                  pid,
		DestinationDirectory: "/tmp",
		SamplingDuration:     duration,
	}
}

func (e *execution) Record(ctx context.Context) (err error) {
	var (
		destination = path.Join(e.DestinationDirectory, e.Uid)
		cmd         = exec.Command("perf", "record", "-g",
			"--freq="+strconv.FormatUint(e.Frequency, 10),
			"--pid="+strconv.FormatUint(e.Pid, 10),
			"--output="+destination)
		cmdChan = make(chan error, 1)
	)

	err = cmd.Start()
	go func() {
		cmdChan <- cmd.Wait()
	}()

	select {
	case err = <-cmdChan:
		if err != nil {
			err = errors.Wrapf(err,
				"command execution failed")
		}
	case <-ctx.Done():
		err = errors.Wrapf(ctx.Err(), "recording cancelled")
	}

	return
}
