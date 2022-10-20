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

package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"testing"
)

func TestZapDevSugar(t *testing.T) {
	err := NewSugaredDevLogger()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a sugared logger", err)
	}
	defer SyncZap()
	S.Debug("Debug test statement.")
}

func TestZapProdSugar(t *testing.T) {
	err := NewSugaredProdLogger()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a sugared logger", err)
	}
	defer SyncZap()
	S.Info("Info test statement.")
}

func TestZapProdSugarLevel(t *testing.T) {
	err := NewSugaredProdLoggerLevel(zap.DebugLevel)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a sugared logger", err)
	}
	defer SyncZap()
	S.Info("Info test statement.")
}

func TestZapPro(t *testing.T) {
	S = nil
	err := NewProdLoggerLevel(zap.DebugLevel)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a prod logger", err)
	}
	defer SyncZap()
	L.Info("Info test statement.")
}

func TestZapLocalLog(t *testing.T) {
	S = nil
	L = nil
	logMsg(zapcore.InfoLevel, "Local Info message")
	logMsg(zapcore.WarnLevel, "Local Warning message")
	err := NewDevLogger()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a prod logger", err)
	}
	defer SyncZap()
	logMsg(zapcore.DebugLevel, "Local Debug message")
	logMsg(zapcore.InfoLevel, "Local Info message")
	logMsg(zapcore.WarnLevel, "Local Warning message")
	logMsg(zapcore.ErrorLevel, "Local Error message")
	logMsg(zapcore.PanicLevel, "Local Panic message")
}

func TestZapSetLevel(t *testing.T) {
	S = nil
	err := NewProdLogger()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a prod logger", err)
	}
	defer SyncZap()
	L.Info("Info message should appear")
	L.Debug("Debug Message should not appear")
	SetLevel("debug")
	L.Debug("Debug Message should appear")
	SetLevel("error")
	L.Debug("Debug Message should not appear")
	L.Warn("Warn Message should not appear")
	L.Error("Error Message should appear")
	SetLevel("debug")
	SetLevel("")
	SetLevel("random")
}

func TestZapHttpLogSet(t *testing.T) {
	S = nil
	err := NewProdLogger()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a prod logger", err)
	}
	defer SyncZap()

	SetupDynamicLogging("")
	SetupDynamicLogging("doesnotexist")
	SetupDynamicLogging("localhost:9999")
}

func TestZapFromLogFile(t *testing.T) {

	err := NewSugaredLoggerFromFile("")
	if err == nil {
		t.Fatalf("expected to get an error from unsupplied config file")
	}
	err = NewSugaredLoggerFromFile("./tests/does-not-exist.json")
	if err == nil {
		t.Fatalf("expected to get an error from non-existant config file")
	}
	fmt.Printf("Got expected error message: %v\n", err)
	err = NewSugaredLoggerFromFile("./tests/zap_config-broken.json")
	if err == nil {
		t.Fatalf("expected to get an error from non-existant config file")
	}
	fmt.Printf("Got expected error message: %v\n", err)
	err = NewSugaredLoggerFromFile("./tests/zap_config.json")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a json logger", err)
	}
	S.Info("Successful JSON logger config loaded")
}
