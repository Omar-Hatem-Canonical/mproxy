// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package hostTranslator

import (
	"context"
	"errors"
	"os"
	"fmt"
	"log/slog"

	"github.com/absmach/mproxy/pkg/session"
)

var errSessionMissing = errors.New("session is missing")

var _ session.Handler = (*hostTranslator)(nil)

// hostTranslator implements mqtt.hostTranslator interface
type hostTranslator struct {
	logger *slog.Logger
}

// New creates new Event entity
func New(logger *slog.Logger) *hostTranslator {
	return &hostTranslator{
		logger: logger,
	}
}

// AuthConnect is called on device connection,
// prior forwarding to the MQTT broker
func (tr *hostTranslator) AuthConnect(ctx context.Context) error {
	return tr.logAction(ctx, "AuthConnect", nil, nil)
}

// AuthPublish is called on device publish,
// prior forwarding to the MQTT broker
func (tr *hostTranslator) AuthPublish(ctx context.Context, topic *string, payload *[]byte) error {
	
	hostName, err := os.Hostname()
	if err == nil {
		msg := fmt.Sprintf("Topic %s is in device %s, Adding namespace and publishing...", *topic, hostName)
		*topic = hostName + "/" + *topic
		tr.logger.Info(msg)
	} else {
		msg := fmt.Sprintf("Topic %s could not be translated and thus kept the same. Publishing...", *topic)
		tr.logger.Info(msg)
	}
	
	
	
	return tr.logAction(ctx, "AuthPublish", &[]string{*topic}, payload)
}

// Reconvert topics on client going down
// Topics are passed by reference, so that they can be modified
func (tr *hostTranslator) DownSubscribe(ctx context.Context, topics *[]string) error {
	return tr.logAction(ctx, "DownSubscribe", topics, nil)
}

// AuthSubscribe is called on device publish,
// prior forwarding to the MQTT broker
func (tr *hostTranslator) AuthSubscribe(ctx context.Context, topics *[]string) error {
	for i := range *topics {
		hostName, err := os.Hostname()
		if err == nil {
			msg := fmt.Sprintf("Topic %s is in device %s, Adding namespace and subscribing...", (*topics)[i], hostName)
			(*topics)[i] = hostName + "/" + (*topics)[i]
			tr.logger.Info(msg)
			} else {
				msg := fmt.Sprintf("Topic %s could not be translated and thus kept the same. Subscribing...", (*topics)[i])
				tr.logger.Info(msg)
		}
	}

	return tr.logAction(ctx, "AuthSubscribe", topics, nil)
}

// Connect - after client successfully connected
func (tr *hostTranslator) Connect(ctx context.Context) error {
	return tr.logAction(ctx, "Connect", nil, nil)
}

// Publish - after client successfully published
func (tr *hostTranslator) Publish(ctx context.Context, topic *string, payload *[]byte) error {
	return tr.logAction(ctx, "Publish", &[]string{*topic}, payload)
}

// Subscribe - after client successfully subscribed
func (tr *hostTranslator) Subscribe(ctx context.Context, topics *[]string) error {
	return tr.logAction(ctx, "Subscribe", topics, nil)
}

// Unsubscribe - after client unsubscribed
func (tr *hostTranslator) Unsubscribe(ctx context.Context, topics *[]string) error {
	return tr.logAction(ctx, "Unsubscribe", topics, nil)
}

// Disconnect on connection lost
func (tr *hostTranslator) Disconnect(ctx context.Context) error {
	return tr.logAction(ctx, "Disconnect", nil, nil)
}

func (tr *hostTranslator) logAction(ctx context.Context, action string, topics *[]string, payload *[]byte) error {
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
		tr.logger.Error(action+"() failed to complete", args...)
		return errSessionMissing
	}
	tr.logger.Info(action+"() completed successfully", args...)

	return nil
}
