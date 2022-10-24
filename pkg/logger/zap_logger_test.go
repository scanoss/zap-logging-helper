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
