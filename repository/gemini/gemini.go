package gemini

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

type geminiRepo struct {
	logger *slog.Logger
	APIKey string
	client *http.Client
}

type Repository interface {
	GenerateExercises(req GenerateRequest) (string, error)
}

func NewGeminiRepository(logger *slog.Logger, apiKey string) Repository {
	return &geminiRepo{
		logger: logger,
		APIKey: apiKey,
		client: &http.Client{},
	}
}

type GenerateRequest struct {
	WorkoutID int64    `json:"workout_id"`
	Target    string   `json:"target"`
	Goal      string   `json:"goal"`
	Equipment []string `json:"equipment"`
	Save      bool     `json:"save"`
}

type respBody struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func (r *geminiRepo) GenerateExercises(req GenerateRequest) (string, error) {
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent"

	prompt := fmt.Sprintf(`
		Generate ONLY **3 to 5** workout exercises as JSON array.

		RULES (IMPORTANT):
		- JSON ONLY, no explanation.
		- NO descriptive sentences.
		- "sets" and "reps" MUST be short values, example: "3-4", "10-12", "30 sec".
		- "equipment" MUST be "none" if empty.
		- Each object MUST follow exactly:

		[
		{ "name": "Push-ups", "sets": "3-4", "reps": "10-15", "equipment": "none" }
		]

		If you are unsure, guess the values. Do NOT include any descriptive text.
		Do NOT add any extra keys.
		Do NOT write paragraphs or instructions.

		Goal: %s
		Target: %s
		Equipment Available: %v
		`, req.Goal, req.Target, req.Equipment)

	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": prompt},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"responseMimeType": "application/json",
			"responseJsonSchema": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"name":      map[string]string{"type": "string"},
						"sets":      map[string]string{"type": "string"},
						"reps":      map[string]string{"type": "string"},
						"equipment": map[string]string{"type": "string"},
					},
					"required": []string{"name", "sets", "reps"},
				},
			},
		},
	}

	body, _ := json.Marshal(payload)

	reqHTTP, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	reqHTTP.Header.Set("x-goog-api-key", r.APIKey)
	reqHTTP.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(reqHTTP)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)

	var data respBody
	if err := json.Unmarshal(respBytes, &data); err != nil {
		return "", err
	}

	if len(data.Candidates) == 0 || len(data.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no valid response from Gemini")
	}

	return strings.TrimSpace(data.Candidates[0].Content.Parts[0].Text), nil
}
