RATIONALE
	pprof[1] is great for a user to understand what's going under
	the hood of a particular golang application, but it limits
	its reach to the userspace.

	perf[2] + flamegraph[3] gives us a deeper view into how the
	kernel is seeing the execution of a given process (or set
	of them).

	By providing an agent that can perform such analysis, we're
	able to mimic the greatness of pprof: providing a known
	endpoint which a remote user can use to gather insights
	about the runtime behavior of an application.

	[1]: https://golang.org/pkg/runtime/pprof/
	[2]: http://man7.org/linux/man-pages/man1/perf.1.html
	[3]: https://github.com/brendangregg/FlameGraph


DEPENDENCIES
	- perf
	- flamegrapudo 
	- kallsyms available


API
	GET /profile?pid=$PID&freq=$FREQ&seconds=$SECONDS

		Performs the capture and streams back in tgz the `perf script`
		results (after a record happens) so that further processing can 
		be performed by the user.

		If no `PID` specified, all of processes get sampled.


	GET /flamegraph?pid=$PID&freq=$FREQ&seconds=$SECONDS

		Performs the capture and then returns a URL where a given
		flamegraph.svg is served, allowing the user to share it.

		If no `PID` specified, all of processes get sampled.


	GET /static/<content> 	

		Retrieves the contents stored under the configured assets
		directory (this is where profiles and framegraphs get saved).

