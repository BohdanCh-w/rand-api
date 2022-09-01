package randapi_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bohdanch-w/rand-api/randapi"
	"github.com/bohdanch-w/rand-api/services"
)

func TestRandomOrgRetrieverNewRequest(t *testing.T) {
	const (
		method = "random"
	)

	params := struct {
		A int
		B int
	}{3, 4}

	svc := randapi.NewRandomOrgRetriever(
		"https://random.org/api",
		&http.Client{},
		false,
	)

	req, err := svc.NewRequest(method, params)
	require.NoError(t, err)

	require.Equal(t, "2.0", req.JsonrpcVersion)
	require.Equal(t, method, req.Method)
	require.JSONEq(t, `{"A": 3, "B": 4}`, string(req.Params))
}

func TestRandomOrgRetrieverNewRequest_InvalidParams(t *testing.T) {
	const (
		method = "random"
	)

	params := (services.RandParameters)(nil)

	svc := randapi.NewRandomOrgRetriever(
		"https://random.org/api",
		&http.Client{},
		false,
	)

	req, err := svc.NewRequest(method, params)
	require.Zero(t, req)
	require.EqualError(t, err, "invalid parameters")
}
