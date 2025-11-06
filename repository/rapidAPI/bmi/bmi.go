package bmi

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

type apiBMIRepository struct {
	logger *slog.Logger
	APIKey string
}

type Repository interface {
	CalculateBMI(weight, height float64, callback func(weight, height float64) float64) (float64, error)
}

func NewRapidAPIRepository(logger *slog.Logger, apiKey string) *apiBMIRepository {
	return &apiBMIRepository{logger: logger, APIKey: apiKey}
}

// CalculateBMI tries the API, calls callback if API fails
func (r *apiBMIRepository) CalculateBMI(weight, height float64, callbackBMI func(weight, height float64) float64) (float64, error) {
	baseURL := "https://body-mass-index-bmi-calculator.p.rapidapi.com/metric"

	u, err := url.Parse(baseURL)
	if err != nil {
		r.logger.Error("failed to parse BMI API URL", "error", err)
		if callbackBMI != nil {
			return callbackBMI(weight, height), nil
		}
		return 0, err
	}

	q := u.Query()
	q.Set("weight", fmt.Sprintf("%f", weight))
	q.Set("height", fmt.Sprintf("%f", height))
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		r.logger.Error("failed to create BMI API request", "error", err)
		if callbackBMI != nil {
			return callbackBMI(weight, height), nil
		}
		return 0, err
	}

	req.Header.Add("x-rapidapi-key", r.APIKey)
	req.Header.Add("x-rapidapi-host", "body-mass-index-bmi-calculator.p.rapidapi.com")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		r.logger.Error("failed to call BMI API", "error", err)
		if callbackBMI != nil {
			return callbackBMI(weight, height), nil
		}
		return 0, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		r.logger.Error("failed to read BMI API response", "error", err)
		if callbackBMI != nil {
			return callbackBMI(weight, height), nil
		}
		return 0, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		r.logger.Error("failed to unmarshal BMI API response", "error", err)
		if callbackBMI != nil {
			return callbackBMI(weight, height), nil
		}
		return 0, err
	}

	bmiVal, ok := result["bmi"].(float64)
	if !ok {
		r.logger.Error("invalid BMI API response format", "response", string(body))
		if callbackBMI != nil {
			return callbackBMI(weight, height), nil
		}
		return 0, fmt.Errorf("invalid response from BMI API")
	}

	return bmiVal, nil
}

func DefaultBMICallback(weight, height float64) float64 {
	if height == 0 {
		return 0
	}
	return weight / (height * height)
}
