package main

import (
	"fmt"
	"log"
	"net/http"

	"go-minio/internal/config"
	"go-minio/internal/handler"
	"go-minio/internal/route"
	"go-minio/internal/service"
)

func main() {
	// 1. Load Config
	cfg := config.LoadConfig()

	// 2. Inisialisasi Service
	storageSvc, err := service.NewStorageService(cfg)
	if err != nil {
		log.Fatalln("Gagal inisialisasi Storage Service:", err)
	}

	// 3. Inisialisasi Handler
	storageHandler := handler.NewStorageHandler(storageSvc)

	// 4. Setup Router
	r := route.Setup(cfg, storageHandler)

	// 5. Start Server
	fmt.Printf("Server berjalan di port :%s...\n", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
