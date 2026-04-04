package judge0

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"wetalk-academy/config"
	"wetalk-academy/package/logger"
)

type SubmissionRequest struct {
	SourceCode     string `json:"source_code"`
	LanguageID     int    `json:"language_id"`
	Stdin          string `json:"stdin,omitempty"`
	ExpectedOutput string `json:"expected_output,omitempty"`
}

type StatusResponse struct {
	Stdout        *string `json:"stdout"`
	Stderr        *string `json:"stderr"`
	CompileOutput *string `json:"compile_output"`
	Message       *string `json:"message"`
	Time          *string `json:"time"`
	Memory        *int    `json:"memory"`
	Status        struct {
		ID          int    `json:"id"`
		Description string `json:"description"`
	} `json:"status"`
}

type Client struct {
	baseURL string
	http    *http.Client
}

func NewClient(conf *config.Config) *Client {
	return &Client{
		baseURL: conf.Judge0.BaseURL,
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) Submit(ctx context.Context, req *SubmissionRequest) (*StatusResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Judge0Client: failed to marshal request: %v", err)
		return nil, fmt.Errorf("failed to marshal request")
	}

	url := fmt.Sprintf("%s/submissions?base64_encoded=false&wait=true", c.baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Judge0Client: failed to create http request: %v", err)
		return nil, fmt.Errorf("failed to create request")
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(httpReq)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Judge0Client: failed to send request: %v", err)
		return nil, fmt.Errorf("failed to submit code")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Judge0Client: failed to read response body: %v", err)
		return nil, fmt.Errorf("failed to read response")
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		logger.ErrorfWithCtx(ctx, "[Err] Judge0Client: unexpected status %d: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("judge0 error: %s", string(body))
	}

	var result StatusResponse
	if err := json.Unmarshal(body, &result); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Judge0Client: failed to unmarshal response: %v", err)
		return nil, fmt.Errorf("failed to parse response")
	}

	logger.InfofWithCtx(ctx, "[Info] Judge0Client: submission completed successfully")
	return &result, nil
}
