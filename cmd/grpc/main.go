package main

import (
	"UserGrpcProj/clickhouse"
	"go.uber.org/zap"

	//appClickHouse "UserGrpcProj/clickhouse"
	appKafka "UserGrpcProj/kafka"
	"UserGrpcProj/kafka/logger"
	"UserGrpcProj/pkg/user/config"
	"UserGrpcProj/pkg/user/server"
	pb "UserGrpcProj/pkg/user/service"
	"context"
	"fmt"
	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	cfg := config.NewUserConfig()
	ctx := context.Background()

	chConn, err := clickhouse.NewClickhouseConf()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("ClickHouse ok")

	go appKafka.StartKafka(chConn)
	kafkaLog := logger.NewLogger()
	fmt.Println("kafka ok")

	db, err := pgxpool.Connect(ctx, cfg.DbConn)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err = MigrateUp("postgres", cfg.DbConn); err != nil {
		kafkaLog.Error("migrations error", zap.Error(err))
		return
	}
	fmt.Println("migrations ok")
	defer db.Close()

	cRds := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	if err := cRds.Ping(ctx).Err(); err != nil {
		kafkaLog.Error("can't connect to redis", zap.Error(err))
		return
	}
	defer cRds.Close()

	userService := server.NewUserService(db, kafkaLog, cRds)
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	srv := grpc.NewServer()
	pb.RegisterUserServer(srv, userService)
	err = srv.Serve(listener)
	if err != nil {
		kafkaLog.Error("error", zap.Error(err))
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
