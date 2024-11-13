// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	// "strings"
	"syscall"

	"github.com/absmach/mproxy/examples/translator"
	"github.com/absmach/mproxy/examples/injector"
	"github.com/absmach/mproxy/examples/hostTranslator"


	"github.com/absmach/mproxy"
	"github.com/absmach/mproxy/examples/simple"
	"github.com/absmach/mproxy/pkg/mqtt"
	"github.com/absmach/mproxy/pkg/session"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
)

const (
	mqttWithoutTLS = "MPROXY_MQTT_WITHOUT_TLS_"
	mqttWithTLS    = "MPROXY_MQTT_WITH_TLS_"
	mqttWithmTLS   = "MPROXY_MQTT_WITH_MTLS_"

)

func main() {
	
	ctx, cancel := context.WithCancel(context.Background())
	g, ctx := errgroup.WithContext(ctx)
	
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := slog.New(logHandler)
	
	topicTranslation := make(map[string]string)
	revTopicTranslation := make(map[string]string)

	topicTranslation["test/topic"] = "test/foo"
	revTopicTranslation["test/foo"] = "test/topic"

	handlerType := flag.Int("hand", 0, "selects the handler that inspects the packets")

	
	var interceptor session.Interceptor
	
	
	pathPtr := flag.String("env", "", "The .env path")
	flag.Parse()
	
	var handler session.Handler

	switch *handlerType {
		case 0: handler = simple.New(logger)
		case 1: handler = translator.New(logger, topicTranslation, revTopicTranslation)
		case 2: handler = injector.New(logger, "Hello")
		case 3: handler = hostTranslator.New(logger)
		default: simple.New(logger)
	}
	
	
	if (*pathPtr == "") {
		// Loading .env file to environment
		err := godotenv.Load()
		if err != nil {
			panic(err)
		}
		
		} else {
			// Loading specified file to environment
			err := godotenv.Load(*pathPtr)
			if err != nil {
				panic(err)
			}
	}

	// mProxy server Configuration for MQTT with TLS
	mqttTLSConfig, err := mproxy.NewConfig(env.Options{Prefix: mqttWithTLS})
	if err != nil {
		panic(err)
	}

	// mProxy server for MQTT with TLS
	mqttTLSProxy := mqtt.New(mqttTLSConfig, handler, interceptor, logger)
	g.Go(func() error {
		return mqttTLSProxy.Listen(ctx)
	})

	g.Go(func() error {
		return StopSignalHandler(ctx, cancel, logger)
	})

	if err := g.Wait(); err != nil {
		logger.Error(fmt.Sprintf("mProxy service terminated with error: %s", err))
	} else {
		logger.Info("mProxy service stopped")
	}
}

func StopSignalHandler(ctx context.Context, cancel context.CancelFunc, logger *slog.Logger) error {
	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGABRT)
	select {
	case <-c:
		cancel()
		return nil
	case <-ctx.Done():
		return nil
	}
}
