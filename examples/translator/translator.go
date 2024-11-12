// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package translator

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/absmach/mproxy/pkg/session"
	"github.com/eclipse/paho.golang/packets"
)

var errSessionMissing = errors.New("session is missing")

var _ session.Handler = (*Translator)(nil)

// Translator implements mqtt.Translator interface
type Translator struct {
	logger *slog.Logger
	topics map[string]string
	revTopics map[string]string
}

// New creates new Event entity
func New(logger *slog.Logger, topics, revTopics map[string]string) *Translator {
	return &Translator{
		logger: logger,
		topics: topics,
		revTopics: revTopics,
	}
}

// AuthConnect is called on device connection,
// prior forwarding to the MQTT broker
func (tr *Translator) AuthConnect(ctx context.Context) error {
	return tr.logAction(ctx, "AuthConnect", nil, nil, nil)
}

// AuthPublish is called on device publish,
// prior forwarding to the MQTT broker
func (tr *Translator) AuthPublish(ctx context.Context, topic *string, payload *[]byte, userProperties *[]packets.User) error {

	newTopic, ok := tr.topics[*topic]
	if ok {
		msg := fmt.Sprintf("Topic %s translated to Topic %s. Publishing...", *topic, newTopic)
		*topic = newTopic
		tr.logger.Info(msg)
	} else {
		msg := fmt.Sprintf("Topic %s could not be translated and thus kept the same. Publishing...", *topic)
		tr.logger.Info(msg)
	}
	
	
	
	return tr.logAction(ctx, "AuthPublish", &[]string{*topic}, payload, userProperties)
}

// AuthSubscribe is called on device publish,
// prior forwarding to the MQTT broker
func (tr *Translator) AuthSubscribe(ctx context.Context, subscriptions *[]packets.SubOptions, userProperties *[]packets.User) error {
	for i, sub := range *subscriptions {
		newTopic, ok := tr.topics[sub.Topic]
		if ok {
			msg := fmt.Sprintf("Topic %s translated to Topic %s. Subscribing...", (*subscriptions)[i].Topic, newTopic)
			(*subscriptions)[i].Topic = newTopic
			tr.logger.Info(msg)
			} else {
				msg := fmt.Sprintf("Topic %s could not be translated and thus kept the same. Subscribing...",  (*subscriptions)[i].Topic)
				tr.logger.Info(msg)
		}
	}

	var topics []string

	for _,x := range *subscriptions {
		topics = append(topics, x.Topic)
	}

	return tr.logAction(ctx, "AuthSubscribe", &topics, nil, userProperties)
}

// Reconvert topics on client going down
// Topics are passed by reference, so that they can be modified
func(tr *Translator) DownSubscribe(ctx context.Context, topic *string, userProperties *[]packets.User) error {
	newTopic, ok := tr.revTopics[*topic]
	if ok {
		msg := fmt.Sprintf("Topic %s translated to Topic %s. Sending...", *topic, newTopic)
		*topic = newTopic
		tr.logger.Info(msg)
		} else {
			msg := fmt.Sprintf("Topic %s could not be translated and thus kept the same. Subscribing...", *topic)
			tr.logger.Info(msg)
	}

	return tr.logAction(ctx, "DownSubscribe",  &[]string{*topic}, nil, nil)
}

// Connect - after client successfully connected
func (tr *Translator) Connect(ctx context.Context) error {
	return tr.logAction(ctx, "Connect", nil, nil, nil)
}

// Publish - after client successfully published
func (tr *Translator) Publish(ctx context.Context, topic *string, payload *[]byte) error {
	return tr.logAction(ctx, "Publish", &[]string{*topic}, payload, nil)
}

// Subscribe - after client successfully subscribed
func (h *Translator) Subscribe(ctx context.Context, subscriptions *[]packets.SubOptions) error {
	var topics []string

	for _,x := range *subscriptions {
		topics = append(topics, x.Topic)
	}
	
	return h.logAction(ctx, "Subscribe", &topics, nil, nil)
}

// Unsubscribe - after client unsubscribed
func (h *Translator) Unsubscribe(ctx context.Context, subscriptions *[]packets.SubOptions) error {
	var topics []string

	for _,x := range *subscriptions {
		topics = append(topics, x.Topic)
	}

	return h.logAction(ctx, "Unsubscribe", &topics, nil, nil)
}

// Disconnect on connection lost
func (tr *Translator) Disconnect(ctx context.Context) error {
	return tr.logAction(ctx, "Disconnect", nil, nil, nil)
}

func (h *Translator) logAction(ctx context.Context, action string, topics *[]string, payload *[]byte, userProperties *[]packets.User) error {
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