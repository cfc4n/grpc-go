/*
 *
 * Copyright 2017 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package backoff implement the backoff strategy for gRPC.
//
// This is kept in internal until the gRPC project decides whether or not to
// allow alternative backoff strategies.
package backoff

import (
	"time"

	"google.golang.org/grpc/internal/grpcrand"
)

// Strategy defines the methodology for backing off after a grpc connection
// failure.
//
type Strategy interface {
	// Backoff returns the amount of time to wait before the next retry given
	// the number of consecutive failures.
	Backoff(retries int) time.Duration
}

const (
	// baseDelay is the amount of time to wait before retrying after the first
	// failure.
	baseDelay = 1.0 * time.Second
	// factor is applied to the backoff after each retry.
	factor = 1.6
	// jitter provides a range to randomize backoff delays.
	jitter = 0.2
)

// Exponential implements exponential backoff algorithm as defined in
// https://github.com/grpc/grpc/blob/master/doc/connection-backoff.md.
type Exponential struct {
	// MaxDelay is the upper bound of backoff delay.
	MaxDelay time.Duration
}

// Backoff returns the amount of time to wait before the next retry given the
// number of retries.
func (bc Exponential) Backoff(retries int) time.Duration {
	var delay uint
	delay = 1 << uint(retries+4)

	//最小delay 时间 30S
	if delay < 30 {
		delay = 30
	}

	//最大delay时间为86400S
	if delay > 86400 {
		delay = 86400
	}

	if delay > 300 {
		delay += uint(grpcrand.Intn(180))
	}

	return time.Duration(delay) * time.Second
}
