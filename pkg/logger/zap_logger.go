// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * Copyright (C) 2018-2022 SCANOSS.COM
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

// Package logger handles logging for everything in the dependency system
// It uses zap to achieve this
package logger

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io/ioutil"
	"net/http"
	"os"
)

var L *zap.Logger
var S *zap.SugaredLogger
var atomicLevel = zap.NewAtomicLevel()

// NewDevLogger creates a new Development logger
func NewDevLogger() error {
	return NewDevLoggerLevel(zapcore.DebugLevel)
}

// NewProdLogger creates a new Production logger
func NewProdLogger() error {
	return NewProdLoggerLevel(zapcore.InfoLevel)
}

// NewDevLoggerLevel creates a Dev logger at the specified logging level
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

// NewProdLoggerLevel creates a Prod logger at the specified logging level
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

// NewSugaredDevLogger creates a new Development Sugared logger
func NewSugaredDevLogger() error {
	if err := NewDevLogger(); err != nil {
		return err
	}
	S = L.Sugar()
	return nil
}

// NewSugaredProdLogger creates a new Production Sugared logger
func NewSugaredProdLogger() error {
	if err := NewProdLogger(); err != nil {
		return err
	}
	S = L.Sugar()
	return nil
}

// NewSugaredProdLoggerLevel creates a new Production Sugared logger at the specified logging level
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
	byteArray, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read logging config file '%v': %v", filename, err)
	}
	var cfg zap.Config
	if err := json.Unmarshal(byteArray, &cfg); err != nil {
		return fmt.Errorf("failed to parse logging config json file '%v': %v", filename, err)
	}
	L, err = cfg.Build()
	if err != nil {
		return fmt.Errorf("failed to load prod logger: %v", err)
	}
	return nil
}

// NewSugaredLoggerFromFile created a sugared logger from the supplied JSON config file
func NewSugaredLoggerFromFile(filename string) error {
	if err := NewLoggerFromFile(filename); err != nil {
		return err
	}
	S = L.Sugar()
	return nil
}

// SetLevel enables the setting of the logging level while the system is still running
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

// SyncZap flushes the buffered logs and captures any sync issues
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
// To set debug status run: curl -X PUT localhost:1065/log/level -d level=debug
func SetupDynamicLogging(addr string) {
	if len(addr) > 0 {
		mux := http.NewServeMux()
		mux.Handle("/log/level", atomicLevel)
		go func() {
			err := http.ListenAndServe(addr, mux)
			if err != nil {
				logMsg(zapcore.ErrorLevel, fmt.Sprintf("Failed to start dynamic logging interface on '%v': %v", addr, err))
			}
		}()
	} else {
		logMsg(zapcore.WarnLevel, "No port/address supplied to enable dynamic logging.")
	}
}

// logMsg logs the given message to the default logger if available, otherwise standard error
func logMsg(level zapcore.Level, msg string) {
	if len(msg) > 0 {
		if L != nil {
			switch level {
			case zapcore.DebugLevel, zapcore.InfoLevel, zapcore.WarnLevel, zapcore.ErrorLevel:
				L.Log(level, msg)
			default:
				L.Warn(msg)
			}
		} else {
			_, _ = fmt.Fprintln(os.Stderr, fmt.Sprintf("%v: %v\n", level.String(), msg))
		}
	}
}
