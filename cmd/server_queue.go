package main

import "github.com/goph/serverz"

// newServerQueue returns a new server queue with all the registered servers.
func newServerQueue(appCtx *application) *serverz.Queue {
	queue := serverz.NewQueue()

	debugServer := newDebugServer(appCtx)
	queue.Prepend(debugServer, nil)

	daemonServer := newDaemonServer(appCtx)
	queue.Append(daemonServer, nil)

	return queue
}
