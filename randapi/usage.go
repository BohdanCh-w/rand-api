package randapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"

	"github.com/bohdanch-w/rand-api/entities"
)

func RandAPIUsageExecute(ctx context.Context, apiKey string) (entities.UsageStatus, error) {
	var (
		usage entities.UsageStatus
		buf   = bytes.NewBuffer(nil)
		enc   = json.NewEncoder(buf)
	)

	enc.SetEscapeHTML(false)

	randReq, err := NewRandomRequest("getUsage", usageStatusParams{ApiKey: apiKey})
	if err != nil {
		return usage, fmt.Errorf("create request: %v", err)
	}

	if err := enc.Encode(randReq); err != nil {
		return usage, fmt.Errorf("encode payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, randAPIPath, buf)
	if err != nil {
		return usage, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("User-Agent", "rand-api/0.1")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return usage, fmt.Errorf("execute request: %w", err)
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return usage, fmt.Errorf("read body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return usage, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var randResp usageStatusResponse

	if err := json.Unmarshal(data, &randResp); err != nil {
		return usage, fmt.Errorf("decode response: %w", err)
	}

	if randResp.ID != randReq.ID {
		return usage, fmt.Errorf("response id mismatch request: %s != %s", randResp.ID.String(), randReq.ID.String())
	}

	if randResp.JsonrpcVersion != jsonRPCVersion {
		return usage, fmt.Errorf("unexpected jsonrpc version: %s", randResp.JsonrpcVersion)
	}

	usage = randResp.Result
	usage.APIKey = apiKey

	return usage, nil
}

type usageStatusParams struct {
	ApiKey string `json:"apiKey"`
}

type usageStatusResponse struct {
	ID             uuid.UUID            `json:"id"`
	JsonrpcVersion string               `json:"jsonrpc"`
	Result         entities.UsageStatus `json:"result"`
}
