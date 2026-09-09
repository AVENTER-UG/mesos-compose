package scheduler

import (
	"context"
	"testing"

	cfg "github.com/m3scluster/mesos-compose/types"
)

func TestEnqueueCommandDoesNotBlockWhenQueueIsFull(t *testing.T) {
	commands := make(chan cfg.Command, 1)
	commands <- cfg.Command{TaskID: "existing"}

	if enqueueCommand(commands, cfg.Command{TaskID: "new"}) {
		t.Fatal("enqueueCommand reported success for a full queue")
	}

	if got := (<-commands).TaskID; got != "existing" {
		t.Fatalf("full queue was modified: got task %q", got)
	}
}

func TestEnqueueCommandQueuesAvailableCommand(t *testing.T) {
	commands := make(chan cfg.Command, 1)
	command := cfg.Command{TaskID: "new"}

	if !enqueueCommand(commands, command) {
		t.Fatal("enqueueCommand rejected an available queue")
	}
	if got := (<-commands).TaskID; got != command.TaskID {
		t.Fatalf("queued task = %q, want %q", got, command.TaskID)
	}
}

func TestRunAfterSubscriptionWaitsForSubscription(t *testing.T) {
	e := &Scheduler{Subscribed: make(chan struct{})}
	called := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go e.RunAfterSubscription(ctx, func(context.Context) { close(called) })
	select {
	case <-called:
		t.Fatal("loop ran before subscription")
	default:
	}

	close(e.Subscribed)
	select {
	case <-called:
	case <-ctx.Done():
		t.Fatal("loop did not run after subscription")
	}
}
