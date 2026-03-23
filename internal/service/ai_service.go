package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// DTO untuk input dari user
type AITaskBreakdownRequest struct {
	TaskTitle string `json:"task_title" binding:"required"`
}

// DTO untuk struktur balasan tunggal dari AI
type SubTaskSuggestion struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
}

// Interface untuk AI Service
type AIService interface {
	GenerateTaskBreakdown(ctx context.Context, taskTitle string) ([]SubTaskSuggestion, error)
}

type aiService struct{}

func NewAIService() AIService {
	return &aiService{}
}

func (s *aiService) GenerateTaskBreakdown(ctx context.Context, taskTitle string) ([]SubTaskSuggestion, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, errors.New("GEMINI_API_KEY belum dikonfigurasi di file .env")
	}

	// 1. Inisialisasi Client Gemini API
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("gagal menginisialisasi client AI: %v", err)
	}
	defer client.Close()

	// 2. Pilih Model: Menggunakan alias -latest untuk kompatibilitas tertinggi di v1beta
	model := client.GenerativeModel("gemini-3-flash-preview")
	
	// Konfigurasi tambahan agar AI merespons dengan format JSON yang ketat
	model.ResponseMIMEType = "application/json"

	// 3. Merakit Prompt (System Prompt & User Prompt)
	prompt := fmt.Sprintf(`Kamu adalah seorang asisten produktivitas dan manajer proyek yang ahli.
Tugas utama kamu adalah memecah sebuah tugas besar menjadi 3 sampai 5 langkah kecil (sub-tasks) yang bisa langsung dieksekusi (actionable).

Tugas yang diberikan user: "%s"

Kembalikan jawaban HANYA dalam format JSON Array mentah seperti struktur berikut, tanpa tambahan markdown (jangan gunakan blok kode backtick):
[
  {
    "title": "judul langkah (maks 50 karakter)",
    "description": "deskripsi singkat apa yang harus dilakukan",
    "priority": "low | medium | high"
  }
]`, taskTitle)

	// 4. Eksekusi Request ke Gemini
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		// Menambahkan log spesifik jika terjadi error dari server Google
		return nil, fmt.Errorf("gagal mendapatkan respons dari server AI: %v", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, errors.New("AI mengembalikan respons kosong")
	}

	// 5. Ekstrak teks dari respons AI
	var aiTextResponse string
	if part, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
		aiTextResponse = string(part)
	} else {
		return nil, errors.New("format respons AI tidak dapat dibaca oleh sistem")
	}

	// Membersihkan backticks markdown jika AI masih mengirimkannya terlepas dari MIME Type
	aiTextResponse = strings.TrimSpace(aiTextResponse)
	aiTextResponse = strings.TrimPrefix(aiTextResponse, "```json")
	aiTextResponse = strings.TrimPrefix(aiTextResponse, "```")
	aiTextResponse = strings.TrimSuffix(aiTextResponse, "```")
	aiTextResponse = strings.TrimSpace(aiTextResponse)

	// 6. Parsing JSON String dari AI menjadi Struct Golang
	var suggestions []SubTaskSuggestion
	err = json.Unmarshal([]byte(aiTextResponse), &suggestions)
	if err != nil {
		return nil, fmt.Errorf("gagal mem-parsing JSON dari AI: %v. Raw Response: %s", err, aiTextResponse)
	}

	return suggestions, nil
}