package main

import "github.com/goph/serverz"

// newServerQueue returns a new server queue with all the registered servers.
func newServerQueue(a *application) *serverz.Queue {
	queue := serverz.NewQueue()

	debugServer := newDebugServer(a)
	queue.Prepend(debugServer, nil)

	daemonServer := newDaemonServer(appCtx)
	queue.Append(daemonServer, nil)

	return queue
}
