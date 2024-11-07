// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package injector

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/absmach/mproxy/pkg/session"
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
	return inj.logAction(ctx, "AuthConnect", nil, nil)
}

// AuthPublish is called on device publish,
// prior forwarding to the MQTT broker
func (inj *Injector) AuthPublish(ctx context.Context, topic *string, payload *[]byte) error {

	newPayloadAsString := string((*payload)) + inj.addedPayload
	newPayloadAsByteStream := []byte(newPayloadAsString)
	*payload = newPayloadAsByteStream


	msg := fmt.Sprintf("Payload now is %s ", string((*payload)))
	inj.logger.Info(msg)
	
	
	return inj.logAction(ctx, "AuthPublish", &[]string{*topic}, payload)
}

// AuthSubscribe is called on device publish,
// prior forwarding to the MQTT broker
func (inj *Injector) AuthSubscribe(ctx context.Context, topics *[]string) error {


	return inj.logAction(ctx, "AuthSubscribe", topics, nil)
}

// Connect - after client successfully connected
func (inj *Injector) Connect(ctx context.Context) error {
	return inj.logAction(ctx, "Connect", nil, nil)
}

// Publish - after client successfully published
func (inj *Injector) Publish(ctx context.Context, topic *string, payload *[]byte) error {
	return inj.logAction(ctx, "Publish", &[]string{*topic}, payload)
}

// Subscribe - after client successfully subscribed
func (inj *Injector) Subscribe(ctx context.Context, topics *[]string) error {
	return inj.logAction(ctx, "Subscribe", topics, nil)
}

// Unsubscribe - after client unsubscribed
func (inj *Injector) Unsubscribe(ctx context.Context, topics *[]string) error {
	return inj.logAction(ctx, "Unsubscribe", topics, nil)
}

// Disconnect on connection lost
func (inj *Injector) Disconnect(ctx context.Context) error {
	return inj.logAction(ctx, "Disconnect", nil, nil)
}

func (inj *Injector) logAction(ctx context.Context, action string, topics *[]string, payload *[]byte) error {
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
	if !ok {
		args = append(args, slog.Any("error", errSessionMissing))
		inj.logger.Error(action+"() failed to complete", args...)
		return errSessionMissing
	}
	inj.logger.Info(action+"() completed successfully", args...)

	return nil
}
