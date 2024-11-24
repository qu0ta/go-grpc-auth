package grpcapp

import (
	"fmt"
	authgrpc "github.com/qu0ta/go-grpc-auth/internal/grpc/auth"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

// New creates a new App instance.
//
// Parameters:
// - log: a pointer to a slog.Logger instance for logging.
// - port: an integer representing the port number to listen on.
//
// Returns:
// - a pointer to an App instance.
func New(log *slog.Logger, port int) *App {
	gRPCServer := grpc.NewServer()

	authgrpc.Register(gRPCServer)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

// Run starts the gRPC server for the application and listens on the specified port.
//
// This method initializes a TCP listener on the configured port, logs the server
// startup details, and begins serving requests using the gRPC server. If any errors
// occur during setup or while serving, they are returned as wrapped errors for
// improved context and debugging.
//
// Returns:
//
//	error: An error if the server fails to start or encounters an issue during operation.
//
// Logging:
//
//	Logs the following details at the INFO level:
//	  - The operation name ("grpcapp.Run")
//	  - The configured port
//	  - The server's address after successfully starting
//
// Error Handling:
//   - Returns an error if the TCP listener fails to initialize.
//   - Returns an error if the gRPC server fails to start serving requests.
func (a *App) Run() error {
	const op = "grpcapp.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Starting gRPC server", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// MustRun starts the Run() method and panics if an error is encountered.
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}

}

// Stop gracefully shuts down the gRPC server for the application.
//
// This method logs the shutdown operation and stops the gRPC server gracefully,
// ensuring that in-progress requests are completed before the server halts.
//
// Logging:
//
//	Logs the following details at the INFO level:
//	  - The operation name ("grpcapp.Stop")
//	  - The configured port
//
// Example Usage:
//
//	app.Stop()
//
// Notes:
//
//	This method should be called during application shutdown to ensure a clean
//	termination of the gRPC server.
func (a *App) Stop() {
	const op = "grpcapp.Stop"
	a.log.With(slog.String("op", op)).Info("Stopping gRPC server", slog.Int("port", a.port))
	a.gRPCServer.GracefulStop()
}
