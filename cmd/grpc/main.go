package main

import (
	"UserGrpcProj/pkg/user/config"
	"UserGrpcProj/pkg/user/server"
	pb "UserGrpcProj/pkg/user/service"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {



	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		os.Stdout,
		zap.DebugLevel,
	))

	cfg := config.NewUserConfig()
	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, cfg.DbConn)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err = MigrateUp("postgres", cfg.DbConn); err != nil{
		fmt.Println(err)
		return
	}
	fmt.Println("migrations ok")

	defer db.Close()
	userService := server.NewUserService(db, logger)
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	srv := grpc.NewServer()
	pb.RegisterUserServer(srv, userService)
	err = srv.Serve(listener)
	if err != nil{
		fmt.Println(err)
	}
	//userHandler := server.NewGin(userService)
	//userServer := &http.Server{
	//	Addr:    fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
	//	Handler: userHandler,
	//}
	//go func() {
	//	// service connections
	//	if err := userServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	//	}
	//}()
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	// catching ctx.Done(). timeout of 5 seconds.
	<-ctx.Done()

	cancel()
}