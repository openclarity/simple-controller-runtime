# Simple Controller Runtime

A simple implementation similar to the [Kubernetes controller runtime](https://github.com/kubernetes-sigs/controller-runtime)
which can be used to create [controllers](https://kubernetes.io/docs/concepts/architecture/controller) for
both Kubernetes and non-Kubernetes resources.

## Architecture

The library has 3 basic constructs:

### Feeders/Detector

This is responsible for detecting resources that need to be reconciled from any
API/datasource and pushes them into the work queue.

The library currently provides a simple Poller feeder which will poll a
GetItems function on a configurable interval and push all returned items into
the queue.

Custom feeders (for example one listening to an event stream) can be written as
long as they can push into the Enqueuer interface implemented by the work
queue.

### Work Queue

This is a FIFO queue which keeps track of resources which need to be
reconciled. It handles duplicate enqueues of resources already in the queue and
keeps track of resources which are currently being processed to allow parallel
reconciliation without race conditions.

It also provides an ability to enqueue/re-queue an item after a certain time
interval, during this period the queue keeps track of the item so that
additional enqueues (triggered by a feeder for example) do not supersede the
requested time interval.

### Reconciler

This is the control logic loop, it reads events from the Work queue and runs
the reconcile function provided for each event.

The reconcile function provided should be idempotent and attempt to move the
resource towards the desired state.

If an error occurs during reconciliation then it is logged by the reconciler
framework and that iteration will be marked as completed by the framework and
the item removed from the queue.

If during the reconciliation it is known that the item will need to be
re-queued for further processing (without waiting for the feeder, or after a
specific amount of time) then there is a RequeueAfter error that can be raised
and the reconciler will re-queue the item.

It is safe to run multiple Reconcilers consuming off the same work queue. Items
can only be in the queue once and items will not be marked as done and removed
from the queue until the reconciler has finished with them; this ensures
parallel reconcilers never end up reconciling the same item so there are no
race conditions.

## Example

A example implementation is available in the examples/ folder, which can be run like:

```
go run example.go
```
