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

// Package interceptor provides helpers to capture/set request/response id in gRPC servers,
// as well as added the request id to the zap logging context.
package interceptor

import (
	"context"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"strings"
)

const RequestIDKey = "x-request-id"
const ResponseIDKey = "x-response-id"

// ContextPropagationUnaryServerInterceptor intercepts the incoming unary request and checks for a Request ID.
// If none exists, create it, add it to the logging dataset and set the Response ID
// It also adds the Request ID to any new outgoing (downstream) requests
func ContextPropagationUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		ctx = getSetRequestID(ctx)
		return handler(ctx, req)
	}
}

// ContextPropagationStreamServerInterceptor intercepts the incoming stream request and checks for a Request ID.
// If none exists, create it, add it to the logging dataset and set the Response ID
// It also adds the Request ID to any new outgoing (downstream) requests
func ContextPropagationStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) (err error) {
		ctx := stream.Context()
		ctx = getSetRequestID(ctx)
		stream = newServerStreamWithContext(stream, ctx)
		return handler(srv, stream)
	}
}

type serverStreamWithContext struct {
	grpc.ServerStream
	ctx context.Context
}

// newServerStreamWithContext returns a new Server Stream with context
func newServerStreamWithContext(stream grpc.ServerStream, ctx context.Context) grpc.ServerStream {
	return serverStreamWithContext{
		ServerStream: stream,
		ctx:          ctx,
	}
}

// getSetRequestID looks for a request ID from incoming metadata
// If none exists, create it, add it to the logging dataset and set the Response ID
// It also adds the Request ID to any new outgoing (downstream) requests
func getSetRequestID(ctx context.Context) context.Context {

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		s := ctxzap.Extract(ctx).Sugar()
		var reqId string
		xrId := md[RequestIDKey] // Check if we have a request ID. If not create one
		if len(xrId) > 0 {
			reqId = strings.Trim(xrId[0], " ")
		}
		if len(reqId) == 0 { // No Request ID, create one
			reqId = uuid.New().String()
			md.Set(RequestIDKey, reqId)
			s.Debugf("Creating Request ID: %v", reqId)
			ctx = metadata.NewIncomingContext(ctx, md) // Add the Request ID to the incoming metadata
		}
		ctxzap.AddFields(ctx, zap.String(RequestIDKey, reqId)) // Add Request ID to the logging
		ctx = context.WithValue(ctx, RequestIDKey, reqId)      // Add Request ID to current context
		ctx = metadata.NewOutgoingContext(ctx, md)             // Add the incoming metadata to any outgoing requests

		header := metadata.New(map[string]string{ResponseIDKey: reqId}) // Set the Response ID
		if err := grpc.SendHeader(ctx, header); err != nil {
			s.Debugf("Warning: Unable to set response header '%v' %v: %v", ResponseIDKey, reqId, err)
		}
	}
	return ctx
}
