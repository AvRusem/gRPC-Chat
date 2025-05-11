package interceptors

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
)

func LoggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	log.Printf("-> Unary call: %s", info.FullMethod)
	start := time.Now()
	resp, err := handler(ctx, req)
	dur := time.Since(start)
	if err != nil {
		log.Printf("<- Completed with error: %v (method %s, %v)", err, info.FullMethod, dur)
	} else {
		log.Printf("<- Completed: method %s, duration=%v", info.FullMethod, dur)
	}
	return resp, err
}

func AuditStreamInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	log.Printf("<-> Stream started: %s (IsClientStream:%v, IsServerStream:%v)",
		info.FullMethod, info.IsClientStream, info.IsServerStream)
	err := handler(srv, ss)
	if err != nil {
		log.Printf("<-> Stream %s finished with error: %v", info.FullMethod, err)
	} else {
		log.Printf("<-> Stream %s finished successfully", info.FullMethod)
	}
	return err
}
