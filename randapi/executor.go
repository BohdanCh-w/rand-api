package randapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/bohdanch-w/rand-api/entities"
	"github.com/bohdanch-w/rand-api/services"
)

const (
	randAPIPath    = "https://api.random.org/json-rpc/4/invoke"
	jsonRPCVersion = "2.0"
)

var _ services.RandRetiever = (*RandomOrgRetriever)(nil)

func NewRandomOrgRetriever(client *http.Client, signed bool) *RandomOrgRetriever {
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

func (svc *RandomOrgRetriever) ExecuteRequest(
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

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, randAPIPath, buf)
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

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, fmt.Errorf("read body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var randResp entities.RandResponse

	if err := json.Unmarshal(data, &randResp); err != nil {
		return result, fmt.Errorf("decode response: %w", err)
	}

	if randResp.Result.Random.Data == nil {
		msg, err := parseErrorResponse(data)
		if err != nil {
			return result, fmt.Errorf("decode error response: %w", err)
		}

		return result, fmt.Errorf("random.org: request failed: %w", msg)
	}

	if randResp.ID != randReq.ID {
		return result, fmt.Errorf("response id mismatch request: %s != %s", randResp.ID.String(), randReq.ID.String())
	}

	if randResp.JsonrpcVersion != jsonRPCVersion {
		return result, fmt.Errorf("unexpected jsonrpc version: %s", randResp.JsonrpcVersion)
	}

	result = randResp.Result

	return result, nil
}

func parseErrorResponse(data []byte) (error, error) {
	var errResp entities.ErrorResponse

	if err := json.Unmarshal(data, &errResp); err != nil {
		return nil, err
	}

	return fmt.Errorf("%d - %s", errResp.Error.Code, errResp.Error.Message), nil
}
