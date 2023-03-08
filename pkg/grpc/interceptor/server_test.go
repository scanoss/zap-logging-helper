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

package interceptor

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

var (
	unaryInfo = &grpc.UnaryServerInfo{
		FullMethod: "TestService.UnaryMethod",
	}
	streamInfo = &grpc.StreamServerInfo{
		FullMethod:     "TestService.StreamMethod",
		IsServerStream: true,
	}
	unaryHandler = func(ctx context.Context, req interface{}) (interface{}, error) {
		return "output", nil
	}
	streamHandler = func(srv interface{}, stream grpc.ServerStream) error {
		return nil
	}
)

type testServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (ss *testServerStream) Context() context.Context {
	return ss.ctx
}

func TestGetSetRequestID(t *testing.T) {
	ctx := context.Background()
	md := metadata.Pairs("x-real-ip", "222.25.118.1")
	ctx = metadata.NewIncomingContext(ctx, md)
	newCtx := getSetRequestID(ctx)
	requestID := RequestIDFromContext(ctx) // old context should be empty
	assert.Falsef(t, len(requestID) > 0, "Request ID should be empty")
	requestID = RequestIDFromContext(newCtx) // new context should have a request ID
	assert.True(t, len(requestID) > 0, "Request ID should not be empty")
}

func TestCtxPropUnaryServerNoReqID(t *testing.T) {
	ctx := context.Background()
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8080")
	if err != nil {
		t.Fatalf("failed to parse TCP Addr: %v", err)
	}
	pd := &peer.Peer{Addr: addr}
	ctx = peer.NewContext(ctx, pd)
	md := metadata.Pairs("x-real-ip", "222.25.118.1")
	ctx = metadata.NewIncomingContext(ctx, md)
	_, err = ContextPropagationUnaryServerInterceptor()(ctx, "xyz", unaryInfo, unaryHandler)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCtxPropUnaryServerReqID(t *testing.T) {
	ctx := context.Background()
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8080")
	if err != nil {
		t.Fatalf("failed to parse TCP Addr: %v", err)
	}
	pd := &peer.Peer{Addr: addr}
	ctx = peer.NewContext(ctx, pd)
	md := metadata.Pairs("x-request-id", "444444")
	ctx = metadata.NewIncomingContext(ctx, md)
	_, err = ContextPropagationUnaryServerInterceptor()(ctx, "xyz", unaryInfo, unaryHandler)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCtxPropStreamServerNoReqID(t *testing.T) {
	ctx := context.Background()
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8080")
	if err != nil {
		t.Fatalf("failed to parse TCP Addr: %v", err)
	}
	pd := &peer.Peer{Addr: addr}
	ctx = peer.NewContext(ctx, pd)
	md := metadata.Pairs("x-real-ip", "222.25.118.1")
	ctx = metadata.NewIncomingContext(ctx, md)
	testService := struct{}{}
	testStream := &testServerStream{ctx: ctx}
	err = ContextPropagationStreamServerInterceptor()(testService, testStream, streamInfo, streamHandler)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
