// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package injector

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

var _ session.Handler = (*Injector)(nil)

// Translator implements mqtt.Translator interface
type Injector struct {
	logger *slog.Logger
	addedPayload string
}

// New creates new Event entity
func New(logger *slog.Logger, addedPayload string) *Injector {
	return &Injector{
		logger: logger,
		addedPayload: addedPayload,
	}
}

// AuthConnect is called on device connection,
// prior forwarding to the MQTT broker
func (inj *Injector) AuthConnect(ctx context.Context) error {
	return inj.logAction(ctx, "AuthConnect", nil, nil, nil)
}

// AuthPublish is called on device publish,
// prior forwarding to the MQTT broker
func (inj *Injector) AuthPublish(ctx context.Context, topic *string, payload *[]byte, userProperties *[]packets.User) error {

	
	
	hostName, _ := os.Hostname()
	addedInfo := fmt.Sprintf(", Device: %s", hostName)
	newPayloadAsString := string((*payload)) + inj.addedPayload + addedInfo
	newPayloadAsByteStream := []byte(newPayloadAsString)
	*payload = newPayloadAsByteStream

	msg := fmt.Sprintf("Payload now is %s", string((*payload)))
	inj.logger.Info(msg)
	
	
	return inj.logAction(ctx, "AuthPublish", &[]string{*topic}, payload, userProperties)
}

// AuthSubscribe is called on device publish,
// prior forwarding to the MQTT broker
func (h *Injector) AuthSubscribe(ctx context.Context, subscriptions *[]packets.SubOptions, userProperties *[]packets.User) error {
	
	var topics []string

	for _,x := range *subscriptions {
		topics = append(topics, x.Topic)
	}

	return h.logAction(ctx, "AuthSubscribe", &topics, nil, userProperties)
}

// Reconvert topics on client going down
// Topics are passed by reference, so that they can be modified
func(inj *Injector) DownSubscribe(ctx context.Context, topics *[]string, userProperties *[]packets.User) error {

	return inj.logAction(ctx, "DownSubscribe", topics, nil, userProperties)
}

// Connect - after client successfully connected
func (inj *Injector) Connect(ctx context.Context) error {
	return inj.logAction(ctx, "Connect", nil, nil, nil)
}

// Publish - after client successfully published
func (inj *Injector) Publish(ctx context.Context, topic *string, payload *[]byte) error {
	return inj.logAction(ctx, "Publish", &[]string{*topic}, payload, nil)
}

// Subscribe - after client successfully subscribed
func (h *Injector) Subscribe(ctx context.Context, subscriptions *[]packets.SubOptions) error {
	var topics []string

	for _,x := range *subscriptions {
		topics = append(topics, x.Topic)
	}
	
	return h.logAction(ctx, "Subscribe", &topics, nil, nil)
}

// Unsubscribe - after client unsubscribed
func (h *Injector) Unsubscribe(ctx context.Context, subscriptions *[]packets.SubOptions) error {
	var topics []string

	for _,x := range *subscriptions {
		topics = append(topics, x.Topic)
	}

	return h.logAction(ctx, "Unsubscribe", &topics, nil, nil)
}

// Disconnect on connection lost
func (inj *Injector) Disconnect(ctx context.Context) error {
	return inj.logAction(ctx, "Disconnect", nil, nil, nil)
}

func (inj *Injector) logAction(ctx context.Context, action string, topics *[]string, payload *[]byte,  userProperties *[]packets.User) error {
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
		inj.logger.Error(action+"() failed to complete", args...)
		return errSessionMissing
	}
	inj.logger.Info(action+"() completed successfully", args...)

	return nil
}
