// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package session

import "context"
import "github.com/eclipse/paho.golang/packets"

// Handler is an interface for mProxy hooks.
type Handler interface {
	// Authorization on client `CONNECT`
	// Each of the params are passed by reference, so that it can be changed
	AuthConnect(ctx context.Context) error

	// Authorization on client `PUBLISH`
	// Topic is passed by reference, so that it can be modified
	AuthPublish(ctx context.Context, topic *string, payload *[]byte, userProperties *[]packets.User) error

	// Authorization on client `SUBSCRIBE`
	// Topics are passed by reference, so that they can be modified
	AuthSubscribe(ctx context.Context, subscriptions *[]packets.SubOptions, userProperties *[]packets.User) error

	// Reconvert topics on client going down
	// Topics are passed by reference, so that they can be modified
	DownSubscribe(ctx context.Context, topics *[]string, userProperties *[]packets.User) error

	// After client successfully connected
	Connect(ctx context.Context) error

	// After client successfully published
	Publish(ctx context.Context, topic *string, payload *[]byte) error

	// After client successfully subscribed
	Subscribe(ctx context.Context, subscriptions *[]packets.SubOptions) error

	// After client unsubscribed
	Unsubscribe(ctx context.Context, subscriptions *[]packets.SubOptions) error

	// Disconnect on connection with client lost
	Disconnect(ctx context.Context) error
}

