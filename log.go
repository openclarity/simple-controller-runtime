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

package simplecontrollerruntime

import (
	"context"
	"log/slog"
)

type ContextKeyType string

const LoggerContextKey ContextKeyType = "Logger"

func GetLoggerFromContextOrDefault(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(LoggerContextKey).(*slog.Logger)
	if ok && logger != nil {
		return logger
	}
	return slog.Default()
}

func SetLoggerForContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, LoggerContextKey, logger)
}
