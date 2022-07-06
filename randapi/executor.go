package randapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const randAPIPath = "https://api.random.org/json-rpc/4/invoke"

func RandAPIExecute(ctx context.Context, randReq *RandomRequest) (RandResponseResult, error) {
	var (
		result RandResponseResult
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

	var randResp RandResponse

	if err := json.Unmarshal(data, &randResp); err != nil {
		return result, fmt.Errorf("decode response: %w", err)
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
