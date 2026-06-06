package main

import (
	"context"
	"fmt"
	llm2 "github.com/bookshop/internal/service/ollama/llm"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/bookshop/internal/logger"
	"github.com/bookshop/internal/storage"
	knowledgeStorage "github.com/bookshop/internal/storage/knowledge_base"

	knowledgeService "github.com/bookshop/internal/service/knowledge_base"
	"github.com/bookshop/internal/service/ollama"
	nomicEmbed "github.com/bookshop/internal/service/ollama/embedding"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("Warning: .env file not found:", err.Error())
	}

	l := logger.New(slog.LevelDebug)
	lstor := l.With("Layer", "Storage")
	lserv := l.With("Layer", "Service")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_NAME"),
	)

	db, err := storage.Connection(connStr)
	if err != nil {
		l.Error("failed to connect database", "error", err)
		return
	}
	defer db.Close()

	ollamaURL := os.Getenv("OLLAMA_URL")
	if ollamaURL == "" {
		ollamaURL = "http://localhost:11434"
	}

	//limit := 50
	//if rawLimit := os.Getenv("EMBEDDING_LIMIT"); rawLimit != "" {
	//	parsedLimit, err := strconv.Atoi(rawLimit)
	//	if err == nil && parsedLimit > 0 {
	//		limit = parsedLimit
	//	}
	//}

	ollamaClient := ollama.New(ollamaURL)
	embedder := nomicEmbed.New(ollamaClient)
	aillm := llm2.New(ollamaClient)

	kbst := knowledgeStorage.New(db, lstor)
	kbserv := knowledgeService.New(kbst, embedder, aillm, lserv)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	if err := kbserv.FillWithoutEmbedding(ctx); err != nil {
		l.Error("failed to fill embeddings", "error", err)
		return
	}

	l.Info("embeddings filled successfully")
}
