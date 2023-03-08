// SPDX-License-Identifier: MIT
/*
 * Copyright (c) 2022, SCANOSS
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 */

// Package logger simplifies the setup of a zap logger.
// It provides helpers for creating Dev & Prod loggers, including setting a logging level.
// These loggers are stored in package level global variables to aid calling them from other packages.
// There is also support for Atomic Levels, which enables the modification of a logging level while the system is running.
package logger

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var L *zap.Logger                      // Global Logger
var S *zap.SugaredLogger               // Global Sugared Logger
var atomicLevel = zap.NewAtomicLevel() // Atomic logging level

// NewDevLogger creates a new Development logger.
func NewDevLogger() error {
	return NewDevLoggerLevel(zapcore.DebugLevel)
}

// NewProdLogger creates a new Production logger.
func NewProdLogger() error {
	return NewProdLoggerLevel(zapcore.InfoLevel)
}

// NewDevLoggerLevel creates a Dev logger at the specified logging level.
func NewDevLoggerLevel(lvl zapcore.Level) error {
	atomicLevel = zap.NewAtomicLevelAt(lvl)
	pc := zap.NewDevelopmentConfig()
	pc.Level = atomicLevel
	var err error
	L, err = pc.Build()
	if err != nil {
		return fmt.Errorf("failed to load dev logger: %v", err)
	}
	return nil
}

// NewProdLoggerLevel creates a Prod logger at the specified logging level.
func NewProdLoggerLevel(lvl zapcore.Level) error {
	atomicLevel = zap.NewAtomicLevelAt(lvl)
	pc := zap.NewProductionConfig()
	pc.Level = atomicLevel
	var err error
	L, err = pc.Build()
	if err != nil {
		return fmt.Errorf("failed to load prod logger: %v", err)
	}
	return nil
}

// NewSugaredDevLogger creates a new Development Sugared logger.
func NewSugaredDevLogger() error {
	if err := NewDevLogger(); err != nil {
		return err
	}
	S = L.Sugar()
	return nil
}

// NewSugaredProdLogger creates a new Production Sugared logger.
func NewSugaredProdLogger() error {
	if err := NewProdLogger(); err != nil {
		return err
	}
	S = L.Sugar()
	return nil
}

// NewSugaredProdLoggerLevel creates a new Production Sugared logger at the specified logging level.
func NewSugaredProdLoggerLevel(lvl zapcore.Level) error {
	if err := NewProdLoggerLevel(lvl); err != nil {
		return err
	}
	S = L.Sugar()
	return nil
}

// NewLoggerFromFile created a logger from the supplied JSON config file
// Details for the fields can be found here: https://github.com/uber-go/zap/blob/master/config.go
func NewLoggerFromFile(filename string) error {
	if filename == "" {
		return fmt.Errorf("no logging config filename provided")
	}
	byteArray, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read logging config file '%v': %v", filename, err)
	}
	var cfg zap.Config
	if err = json.Unmarshal(byteArray, &cfg); err != nil {
		return fmt.Errorf("failed to parse logging config json file '%v': %v", filename, err)
	}
	atomicLevel = cfg.Level // Assign the atomic level parsed from the config file
	L, err = cfg.Build()
	if err != nil {
		return fmt.Errorf("failed to load prod logger: %v", err)
	}
	return nil
}

// NewSugaredLoggerFromFile created a sugared logger from the supplied JSON config file.
func NewSugaredLoggerFromFile(filename string) error {
	if err := NewLoggerFromFile(filename); err != nil {
		return err
	}
	S = L.Sugar()
	return nil
}

// SetLevel enables the setting of the logging level while the system is still running.
func SetLevel(level string) {
	if len(level) > 0 {
		l, err := zapcore.ParseLevel(level)
		if err != nil {
			logMsg(zapcore.WarnLevel, fmt.Sprintf("Failed to set level '%v': %v. Ignoring.", level, err))
		} else {
			logMsg(zapcore.InfoLevel, fmt.Sprintf("Setting logging level to %v.", l.String()))
			atomicLevel.SetLevel(l)
		}
	} else {
		logMsg(zapcore.WarnLevel, "No level supplied to set")
	}
}

// SyncZap flushes the buffered logs and captures any sync issues.
func SyncZap() {
	// Sync the Sugared logger if it's set
	if S != nil {
		err := S.Sync()
		if err != nil {
			fmt.Printf("Warning: Failed to sync zap: %v\n", err)
		}
	} else if L != nil { // Otherwise, sync the Logger if it's set
		err := L.Sync()
		if err != nil {
			fmt.Printf("Warning: Failed to sync zap: %v\n", err)
		}
	}
}

// SetupDynamicLogging enables the ability to modify logging levels on the fly
// Details on how to call the endpoint can be found here: https://pkg.go.dev/go.uber.org/zap#section-readme
// To get debug status run: curl -X GET localhost:1065/log/level
// To set debug status run: curl -X PUT localhost:1065/log/level -d level=debug.
func SetupDynamicLogging(addr string) {
	if len(addr) > 0 {
		mux := http.NewServeMux()
		mux.Handle("/log/level", atomicLevel)
		server := &http.Server{
			Addr:              addr,
			ReadHeaderTimeout: 3 * time.Second,
			Handler:           mux,
		}
		go func() {
			err := server.ListenAndServe()
			if err != nil {
				logMsg(zapcore.ErrorLevel, fmt.Sprintf("Failed to start dynamic logging interface on '%v': %v", addr, err))
			}
		}()
	} else {
		logMsg(zapcore.WarnLevel, "No port/address supplied to enable dynamic logging.")
	}
}

// logMsg logs the given message to the default logger if available, otherwise standard error.
func logMsg(level zapcore.Level, msg string) {
	if len(msg) > 0 {
		if L != nil {
			switch level {
			case zapcore.DebugLevel, zapcore.InfoLevel, zapcore.WarnLevel, zapcore.ErrorLevel:
				L.Log(level, msg)
			case zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel, zapcore.InvalidLevel:
				L.Log(level, msg)
			default:
				L.Warn(msg)
			}
		} else {
			message := fmt.Sprintf("%v: %v\n", level.String(), msg)
			_, _ = fmt.Fprintln(os.Stderr, message)
		}
	}
}

// SetupAppLogger creates a zap logger based on the application configuration options.
func SetupAppLogger(appMode, configFile string, appDebug bool) error {
	var err error
	switch strings.ToLower(appMode) {
	case "prod":
		if len(configFile) > 0 {
			err = NewSugaredLoggerFromFile(configFile)
		} else {
			err = NewSugaredProdLogger()
		}
	default:
		if len(configFile) > 0 {
			err = NewSugaredLoggerFromFile(configFile)
		} else {
			err = NewSugaredDevLogger()
		}
	}
	if err != nil {
		return fmt.Errorf("failed to load logger: %v", err)
	}
	if appDebug {
		SetLevel("debug")
	}
	L.Debug("Running with debug enabled")
	return nil
}

// SetupAppDynamicLogging enables dynamic app logging if requested.
func SetupAppDynamicLogging(dynamicPort string, dynamicLogging bool) {
	if dynamicLogging && len(dynamicPort) > 0 {
		S.Infof("Setting up dynamic logging level on %v.", dynamicPort)
		SetupDynamicLogging(dynamicPort)
		S.Infof("Use the following to get the current status: curl -X GET %v/log/level", dynamicPort)
		S.Infof("Use the following to set the current status: curl -X PUT %v/log/level -d level=debug", dynamicPort)
	}
}
