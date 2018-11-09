package perf

import (
	"bytes"
	"context"
	"fmt"
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

	var (
		buffer = new(bytes.Buffer)
		cmd    = exec.Command("/bin/bash", "-e", "-o", "pipefail", "-c",
			fmt.Sprintf("perf script --input=%s | stackcollapse-perf.pl | flamegraph.pl --width 2000", profile))
	)

	cmd.Stdout = w
	cmd.Stderr = buffer

	err = RunWithContext(ctx, cmd)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to generate flamegraph for profile %s - %s", profile, string(buffer.Bytes()))
		return
	}

	return
}

func (e *execution) Record(ctx context.Context, outputFile string) (err error) {
	argv := []string{"record", "-g", "--freq", strconv.FormatUint(e.frequency, 10)}

	if e.pid != 0 {
		argv = append(argv, "--pid="+strconv.FormatUint(e.pid, 10))
	} else {
		argv = append(argv, "-a")
	}

	argv = append(argv, "--output="+outputFile, "sleep", strconv.FormatUint(e.samplingDuration, 10))

	var (
		buffer = new(bytes.Buffer)
		cmd    = exec.Command("perf", argv...)
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

	fmt.Println(cmd.Args)

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
