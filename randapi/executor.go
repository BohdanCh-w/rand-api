package randapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/bohdanch-w/rand-api/entities"
	"github.com/bohdanch-w/rand-api/services"
	"github.com/google/uuid"
)

const jsonRPCVersion = "2.0"

const (
	ErrUnexpectedStatusCode     = entities.Error("unexpected status code")
	ErrRequestResponseMissmatch = entities.Error("request and response id mismatch")
	ErrUnexpectedJSONRPSVersion = entities.Error("unexpected json rpc version")
	ErrErrorInResponse          = entities.Error("error in response")
)

var _ services.RandRetiever = (*RandomOrgRetriever)(nil)

func NewRandomOrgRetriever(randAPIPath string, client *http.Client, signed bool) *RandomOrgRetriever {
	return &RandomOrgRetriever{
		apiPath:        randAPIPath,
		jsonRPCVersion: jsonRPCVersion,
		client:         client,
	}
}

type RandomOrgRetriever struct {
	apiPath        string
	jsonRPCVersion string
	client         *http.Client
}

func (svc *RandomOrgRetriever) ExecuteRequest( // nolint: funlen
	ctx context.Context,
	randReq *entities.RandomRequest,
) (entities.RandResponseResult, error) {
	var (
		result entities.RandResponseResult
		buf    = bytes.NewBuffer(nil)
		enc    = json.NewEncoder(buf)
	)

	enc.SetEscapeHTML(false)

	if err := enc.Encode(randReq); err != nil {
		return result, fmt.Errorf("encode payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, svc.apiPath, buf)
	if err != nil {
		return result, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("User-Agent", "rand-api/0.1")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return result, fmt.Errorf("execute request: %w", err)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, fmt.Errorf("read body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("%w: %d", ErrUnexpectedStatusCode, resp.StatusCode)
	}

	var randResp randResponse

	if err := randResp.parse(data); err != nil {
		return result, fmt.Errorf("invalid response: %w", err)
	}

	if randResp.ID != randReq.ID {
		return result, fmt.Errorf("%w: %s != %s", ErrRequestResponseMissmatch, randResp.ID.String(), randReq.ID.String())
	}

	result = *randResp.Result

	return result, nil
}

type randResponse struct {
	ID             uuid.UUID                    `json:"id"`
	JsonrpcVersion string                       `json:"jsonrpc"`
	Result         *entities.RandResponseResult `json:"result"`
	Error          *entities.ErrorResponse      `json:"error"`
}

func (resp *randResponse) parse(data []byte) error {
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
