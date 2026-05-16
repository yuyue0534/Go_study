package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"school-student-system/internal/config"
	"school-student-system/internal/database"
	"school-student-system/internal/handler"
	"school-student-system/internal/repository"
	"school-student-system/internal/router"
	"school-student-system/internal/service"
)

func main() {
	cfg := config.Load()

	db, err := database.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("database open failed: %v", err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		log.Fatalf("database migration failed: %v", err)
	}

	studentRepo := repository.NewStudentRepository(db)
	studentService := service.NewStudentService(studentRepo)
	studentHandler := handler.NewStudentHandler(studentService)

	httpRouter := router.New(studentHandler)
	server := &http.Server{
		Addr:              cfg.ListenAddr(),
		Handler:           httpRouter,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("student system started on http://%s", cfg.ListenAddr())
		log.Printf("database path: %s", cfg.DBPath)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http server failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown failed: %v", err)
		return
	}

	log.Println("student system stopped gracefully")
}
