package perf

import (
	"bytes"
	"context"
	"io"
	"os/exec"
	"strconv"

	"github.com/pkg/errors"
)

// Execution provides the context arguments for the execution of the
// the perf and flamegraph scripts.
type execution struct {
	frequency uint64
	pid       uint64

	// samplingDuration describes how long the sampling should
	// take in seconds.
	samplingDuration uint64
}

func NewExecution(frequency, pid, duration uint64) (e *execution, err error) {
	e = &execution{
		frequency:        frequency,
		pid:              pid,
		samplingDuration: duration,
	}

	return
}

func (e *execution) GenerateFlamegraph(ctx context.Context, profile string, w io.Writer) (err error) {
	if profile == "" {
		err = errors.Errorf("a profile must be specified")
		return
	}

	cmd := exec.Command("/bin/bash", "-e", "-o", "pipefail", "-c",
		"perf script | stackcollapse-perf.pl | flamegraph.pl")
	cmd.Stdout = w

	err = RunWithContext(ctx, cmd)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to generate flamegraph for profile %s", profile)
		return
	}

	return
}

func (e *execution) Record(ctx context.Context, outputFile string) (err error) {
	var (
		cmd = exec.Command("perf", "record", "-g",
			"--freq", strconv.FormatUint(e.frequency, 10),
			"--pid", strconv.FormatUint(e.pid, 10),
			"--output="+outputFile,
			"sleep", strconv.FormatUint(e.samplingDuration, 10))
		buffer = new(bytes.Buffer)
	)

	cmd.Stderr = buffer

	err = RunWithContext(ctx, cmd)
	if err != nil {
		err = errors.Wrapf(err,
			"'perf record' execution failed - %s", string(buffer.Bytes()))
		return
	}

	return
}

func RunWithContext(ctx context.Context, cmd *exec.Cmd) (err error) {
	var cmdChan = make(chan error, 1)

	err = cmd.Start()
	go func() {
		cmdChan <- cmd.Wait()
	}()

	select {
	case err = <-cmdChan:
	case <-ctx.Done():
		err = ctx.Err()
	}

	return
}
