package main

import (
	"context"
	"finance-dashboard/engine/api"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

func main() {
	dataPath := findDataPath()
	log.Printf("Finance Dashboard Engine - datos: %s", dataPath)

	server := api.NewServer(dataPath)

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	} else if port[0] != ':' {
		port = ":" + port
	}

	go func() {
		log.Printf("Engine corriendo en http://localhost%s", port)
		if err := server.Start(port); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error en servidor: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Apagando servidor...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Error apagando servidor: %v", err)
	}
	log.Println("Servidor detenido correctamente.")
}

func findDataPath() string {
	execPath, err := os.Executable()
	if err == nil {
		projectRoot := filepath.Join(filepath.Dir(execPath), "..")
		dataPath := filepath.Join(projectRoot, "data")
		if fi, err := os.Stat(dataPath); err == nil && fi.IsDir() {
			return dataPath
		}
	}

	for _, candidate := range []string{"../data", "data"} {
		if fi, err := os.Stat(candidate); err == nil && fi.IsDir() {
			abs, _ := filepath.Abs(candidate)
			return abs
		}
	}

	log.Println("No se encontró la carpeta data/. Ejecuta el collector de Python primero:")
	log.Println("   cd ../collector && python main.py")
	os.Exit(1)
	return ""
}
