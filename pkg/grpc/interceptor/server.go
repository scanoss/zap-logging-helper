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

// ContextPropagationUnaryServerInterceptor intercepts the incoming request and checks for a Request ID.
// If none exists, create it, add it to the logging dataset and set the Response ID
func ContextPropagationUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
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
		return handler(ctx, req)
	}
}
