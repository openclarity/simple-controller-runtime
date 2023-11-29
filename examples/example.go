// Copyright © 2023 Cisco Systems, Inc. and its affiliates.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"time"
	"context"
	"os/signal"
	"syscall"
	"log/slog"

	"github.com/openclarity/simple-controller-runtime"
)

type ExampleResource struct {
	DesiredValue int
	CurrentValue int
}

var DataStore map[string]ExampleResource = make(map[string]ExampleResource)

type ExampleEvent struct {
	ID string
}

func (e ExampleEvent) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("ID", e.ID),
	)
}

func (e ExampleEvent) String() string {
	return e.ID
}

func (e ExampleEvent) Hash() string {
	return e.ID
}

func GetItems(ctx context.Context) ([]ExampleEvent, error) {
	needsReconcile := make([]ExampleEvent, 0, 0)
	for id, item := range DataStore {
		if item.DesiredValue != item.CurrentValue {
			needsReconcile = append(needsReconcile, ExampleEvent{ID: id})
		}
	}
	return needsReconcile, nil
}

func Reconcile(ctx context.Context, e ExampleEvent) error {
	resource := DataStore[e.ID]

	logger := simplecontrollerruntime.GetLoggerFromContextOrDefault(ctx)

	if resource.DesiredValue != resource.CurrentValue {
		logger.InfoContext(ctx, "Values are out of sync", "desired", resource.DesiredValue, "current", resource.CurrentValue)
		resource.CurrentValue = resource.DesiredValue
		DataStore[e.ID] = resource
		return simplecontrollerruntime.NewRequeueAfterError(time.Second * 5, "requeue to demo this feature")
	}

	logger.InfoContext(ctx, "Nothing to change")

	return nil
}

func StartController(ctx context.Context) {
	queue := simplecontrollerruntime.NewQueue[ExampleEvent]()

	poller := simplecontrollerruntime.Poller[ExampleEvent]{
		PollPeriod: time.Second * 30,
		Queue: queue,
		GetItems: GetItems,
	}
	poller.Start(ctx)

	reconciler := simplecontrollerruntime.Reconciler[ExampleEvent]{
		ReconcileTimeout: time.Minute * 5,
		Queue: queue,
		ReconcileFunction: Reconcile,
	}
	reconciler.Start(ctx)
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	StartController(ctx)

	DataStore["foo"] = ExampleResource{
		DesiredValue: 2,
		CurrentValue: 1,
	}

	<-ctx.Done()
}
