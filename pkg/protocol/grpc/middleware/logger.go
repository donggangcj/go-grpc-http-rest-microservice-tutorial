package middleware

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// codeToLevel redirects OK to DEBUG level logging instead of INFO
// This is example how you can log several gRPC code results
func codeToLevel(code codes.Code) zapcore.Level {
	if code == codes.OK {
		// It is debug
		return zapcore.DebugLevel
	}
	return grpc_zap.DefaultClientCodeToLevel(code)
}

// AddLogging return grpc.Server config option that return to logging
func AddLogging(logger *zap.Logger, opts []grpc.ServerOption) []grpc.ServerOption {
	// shared options for the logger,with a custom gRPC code to log level function
	o := []grpc_zap.Option{grpc_zap.WithLevels(codeToLevel)}

	// Make sure that log statements internal to gRPC library are logged using tha zapLogger as well
	grpc_zap.ReplaceGrpcLoggerV2(logger)

	// Add unary interceptor
	opts = append(opts, grpc_middleware.WithUnaryServerChain(
		grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_zap.UnaryServerInterceptor(logger, o...)))

	// Add stream interceptor
	opts = append(opts, grpc_middleware.WithStreamServerChain(
		grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_zap.StreamServerInterceptor(logger, o...)))


	return opts
}
