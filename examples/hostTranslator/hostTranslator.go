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

var _ session.Handler = (*Translator)(nil)

// Translator implements mqtt.Translator interface
type Translator struct {
	logger *slog.Logger
}

// New creates new Event entity
func New(logger *slog.Logger) *Translator {
	return &Translator{
		logger: logger,
	}
}

// AuthConnect is called on device connection,
// prior forwarding to the MQTT broker
func (tr *Translator) AuthConnect(ctx context.Context) error {
	return tr.logAction(ctx, "AuthConnect", nil, nil)
}

// AuthPublish is called on device publish,
// prior forwarding to the MQTT broker
func (tr *Translator) AuthPublish(ctx context.Context, topic *string, payload *[]byte) error {
	
	hostName, err := os.Hostname()
	if err == nil {
		msg := fmt.Sprintf("Topic %s is in device %s, Adding namespace and publishing...", *topic, hostName)
		*topic = hostName + *topic
		tr.logger.Info(msg)
	} else {
		msg := fmt.Sprintf("Topic %s could not be translated and thus kept the same. Publishing...", *topic)
		tr.logger.Info(msg)
	}
	
	
	
	return tr.logAction(ctx, "AuthPublish", &[]string{*topic}, payload)
}

// AuthSubscribe is called on device publish,
// prior forwarding to the MQTT broker
func (tr *Translator) AuthSubscribe(ctx context.Context, topics *[]string) error {
	for i, _ := range *topics {
		hostName, err := os.Hostname()
		if err == nil {
			msg := fmt.Sprintf("Topic %s is in device %s, Adding namespace and subscribing...", (*topics)[i], hostName)
			(*topics)[i] = hostName + (*topics)[i]
			tr.logger.Info(msg)
			} else {
				msg := fmt.Sprintf("Topic %s could not be translated and thus kept the same. Subscribing...", (*topics)[i])
				tr.logger.Info(msg)
		}
	}

	return tr.logAction(ctx, "AuthSubscribe", topics, nil)
}

// Connect - after client successfully connected
func (tr *Translator) Connect(ctx context.Context) error {
	return tr.logAction(ctx, "Connect", nil, nil)
}

// Publish - after client successfully published
func (tr *Translator) Publish(ctx context.Context, topic *string, payload *[]byte) error {
	return tr.logAction(ctx, "Publish", &[]string{*topic}, payload)
}

// Subscribe - after client successfully subscribed
func (tr *Translator) Subscribe(ctx context.Context, topics *[]string) error {
	return tr.logAction(ctx, "Subscribe", topics, nil)
}

// Unsubscribe - after client unsubscribed
func (tr *Translator) Unsubscribe(ctx context.Context, topics *[]string) error {
	return tr.logAction(ctx, "Unsubscribe", topics, nil)
}

// Disconnect on connection lost
func (tr *Translator) Disconnect(ctx context.Context) error {
	return tr.logAction(ctx, "Disconnect", nil, nil)
}

func (tr *Translator) logAction(ctx context.Context, action string, topics *[]string, payload *[]byte) error {
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
