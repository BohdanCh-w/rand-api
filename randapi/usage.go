package randapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/bohdanch-w/rand-api/entities"
	"github.com/google/uuid"
)

func (svc *RandomOrgRetriever) GetUsage(ctx context.Context, apiKey string) (entities.UsageStatus, error) {
	var (
		usage entities.UsageStatus
		buf   = bytes.NewBuffer(nil)
		enc   = json.NewEncoder(buf)
	)

	enc.SetEscapeHTML(false)

	randReq, err := svc.NewRequest("getUsage", usageStatusParams{APIKey: apiKey})
	if err != nil {
		return usage, fmt.Errorf("create request: %w", err)
	}

	if err := enc.Encode(randReq); err != nil {
		return usage, fmt.Errorf("encode payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, svc.apiPath, buf)
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

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return usage, fmt.Errorf("read body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return usage, fmt.Errorf("%w: %d", ErrUnexpectedStatusCode, resp.StatusCode)
	}

	var randResp UsageStatusResponse

	if err := randResp.parse(data); err != nil {
		return usage, fmt.Errorf("invalid response: %w", err)
	}

	if randResp.ID != randReq.ID {
		return usage, fmt.Errorf("%w: %s != %s", ErrRequestResponseMissmatch, randResp.ID.String(), randReq.ID.String())
	}

	usage = *randResp.Result
	usage.APIKey = apiKey

	return usage, nil
}

type usageStatusParams struct {
	APIKey string `json:"apiKey"`
}

type UsageStatusResponse struct {
	ID             uuid.UUID               `json:"id"`
	JsonrpcVersion string                  `json:"jsonrpc"`
	Result         *entities.UsageStatus   `json:"result"`
	Error          *entities.ErrorResponse `json:"error"`
}

func (resp *UsageStatusResponse) parse(data []byte) error {
	if err := json.Unmarshal(data, resp); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	if resp.Error != nil {
		return fmt.Errorf("%w: %d - %s", ErrErrorInResponse, resp.Error.Code, resp.Error.Message)
	}

	if resp.Result == nil {
		return entities.Error("missing result in response")
	}

	if resp.JsonrpcVersion != jsonRPCVersion {
		return fmt.Errorf("%w: %s", ErrUnexpectedJSONRPSVersion, resp.JsonrpcVersion)
	}

	return nil
}
