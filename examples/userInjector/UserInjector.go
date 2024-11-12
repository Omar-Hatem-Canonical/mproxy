// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package UserInjector

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/absmach/mproxy/pkg/session"
	"github.com/eclipse/paho.golang/packets"
)

var errSessionMissing = errors.New("session is missing")

var _ session.Handler = (*UserInjector)(nil)

// UserInjector implements mqtt.UserInjector interface
type UserInjector struct {
	logger *slog.Logger
}

// New creates new Event entity
func New(logger *slog.Logger) *UserInjector {
	return &UserInjector{
		logger: logger,
	}
}

// AuthConnect is called on device connection,
// prior forwarding to the MQTT broker
func (h *UserInjector) AuthConnect(ctx context.Context) error {
	return h.logAction(ctx, "AuthConnect", nil, nil, nil)
}

// AuthPublish is called on device publish,
// prior forwarding to the MQTT broker
func (h *UserInjector) AuthPublish(ctx context.Context, topic *string, payload *[]byte, userProperties *[]packets.User) error {
	deviceName, err := os.Hostname()
	if err == nil {
		newUserProperty := packets.User{
			Key: "Device",
			Value: deviceName,
		}

		*userProperties = append(*userProperties, newUserProperty)

		msg := fmt.Sprintf("New user property added: %s", newUserProperty)
		h.logger.Info(msg)
	} 

	return h.logAction(ctx, "AuthPublish", &[]string{*topic}, payload, userProperties)
}

// AuthSubscribe is called on device publish,
// prior forwarding to the MQTT broker
func (h *UserInjector) AuthSubscribe(ctx context.Context, subscriptions *[]packets.SubOptions, userProperties *[]packets.User) error {
	
	var topics []string

	for _,x := range *subscriptions {
		topics = append(topics, x.Topic)
	}

	return h.logAction(ctx, "AuthSubscribe", &topics, nil, userProperties)
}

// Reconvert topics on client going down
// Topics are passed by reference, so that they can be modified
func(h *UserInjector) DownSubscribe(ctx context.Context, topics *[]string, userProperties *[]packets.User) error {

	return h.logAction(ctx, "DownSubscribe", topics, nil, userProperties)
}


// Connect - after client successfully connected
func (h *UserInjector) Connect(ctx context.Context) error {
	return h.logAction(ctx, "Connect", nil, nil, nil)
}

// Publish - after client successfully published
func (h *UserInjector) Publish(ctx context.Context, topic *string, payload *[]byte) error {
	return h.logAction(ctx, "Publish", &[]string{*topic}, payload, nil)
}

// Subscribe - after client successfully subscribed
func (h *UserInjector) Subscribe(ctx context.Context, subscriptions *[]packets.SubOptions) error {
	var topics []string

	for _,x := range *subscriptions {
		topics = append(topics, x.Topic)
	}
	
	return h.logAction(ctx, "Subscribe", &topics, nil, nil)
}

// Unsubscribe - after client unsubscribed
func (h *UserInjector) Unsubscribe(ctx context.Context, subscriptions *[]packets.SubOptions) error {
	var topics []string

	for _,x := range *subscriptions {
		topics = append(topics, x.Topic)
	}

	return h.logAction(ctx, "Unsubscribe", &topics, nil, nil)
}

// Disconnect on connection lost
func (h *UserInjector) Disconnect(ctx context.Context) error {
	return h.logAction(ctx, "Disconnect", nil, nil, nil)
}

func (h *UserInjector) logAction(ctx context.Context, action string, topics *[]string, payload *[]byte, userProperties *[]packets.User) error {
	s, ok := session.FromContext(ctx)
	args := []interface{}{
		slog.Group("session", slog.String("id", s.ID), slog.String("username", s.Username)),
	}
	if s.Cert.Subject.CommonName != "" {
		args = append(args, slog.Group("cert", slog.String("cn", s.Cert.Subject.CommonName)))
	}
	if topics != nil {
		args = append(args, slog.Any("topics", *topics))
	}
	if payload != nil {
		args = append(args, slog.Any("payload", *payload))
	}
	if userProperties != nil {
		args = append(args, slog.Any("user_properties", *userProperties))
	}
	if !ok {
		args = append(args, slog.Any("error", errSessionMissing))
		h.logger.Error(action+"() failed to complete", args...)
		return errSessionMissing
	}
	h.logger.Info(action+"() completed successfully", args...)

	return nil
}
