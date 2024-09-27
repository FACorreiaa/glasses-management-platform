package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	"github.com/FACorreiaa/glasses-management-platform/app"
	"github.com/FACorreiaa/glasses-management-platform/config"
	"github.com/FACorreiaa/glasses-management-platform/db"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}

func handleError(err error, message string) {
	if err != nil {
		log.Printf("%s: %v", message, err)
	}
}

func hasPassword(password string) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hashed Password: %s\n", hash)
}

func run(ctx context.Context) error {
	hasPassword(os.Getenv("ADMIN_PASSWORD"))
	cfg, err := config.NewConfig()

	if err != nil {
		return err
	}

	pprofAddr := os.Getenv("PPROF_ADDR")
	pprofPort := os.Getenv("PPROF_PORT")
	var logHandler slog.Handler

	logHandlerOptions := slog.HandlerOptions{
		AddSource: true,
		Level:     cfg.Log.Level,
	}
	if cfg.Log.Format == "json" {
		logHandler = slog.NewJSONHandler(os.Stdout, &logHandlerOptions)
	} else {
		logHandler = slog.NewTextHandler(os.Stdout, &logHandlerOptions)
	}
	slog.SetDefault(slog.New(logHandler))

	pool, err := db.Init(cfg.Database.ConnectionURL)
	if err != nil {
		log.Println(err)
	}
	defer pool.Close()

	db.WaitForDB(pool)

	if err = db.Migrate(pool); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)

	}

	srv := &http.Server{
		Addr:         cfg.Server.Addr,
		WriteTimeout: cfg.Server.WriteTimeout,
		ReadTimeout:  cfg.Server.ReadTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
		Handler:      app.Router(pool, []byte(cfg.Server.SessionKey)),
	}

	go func() {
		slog.Info("Starting server " + cfg.Server.Addr)
		if err = srv.ListenAndServe(); err != nil {
			slog.Error("ListenAndServe", "error", err)
		}
	}()

	err = config.InitPprof(pprofAddr, pprofPort)
	if err != nil {
		fmt.Printf(" initializing pprof config: %s", err)
		panic(err)
	}

	<-ctx.Done() // Wait for cancellation signal

	// Shutdown server
	ctxShutdown, cancel := context.WithTimeout(context.Background(), cfg.Server.GracefulTimeout)
	defer cancel()

	if err = srv.Shutdown(ctxShutdown); err != nil {
		handleError(err, " shutting down server")
	}

	slog.Info("Shutting down")
	return nil
}

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		cancel()
		os.Exit(1)
	}
}
