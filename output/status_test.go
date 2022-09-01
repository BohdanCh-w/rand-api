package output_test

import (
	"strings"
	"testing"
	"time"

	"github.com/bohdanch-w/rand-api/entities"
	"github.com/bohdanch-w/rand-api/output"
	"github.com/stretchr/testify/require"
)

func TestGenerateUsageOutput(t *testing.T) {
	status := entities.UsageStatus{
		APIKey:        "b959bdf4-8a46-480d-a72b-f5d1be7711a6",
		Status:        "running",
		CreationTime:  entities.RandTime(time.Date(2022, 8, 25, 12, 0, 0, 0, time.UTC)),
		TotalRequests: 200,
		TotalBits:     200_000,
		RequestsLeft:  1000,
		BitsLeft:      250_000,
	}

	expected := `
Usage statistic for API key b959bdf4-8a46-480d-a72b-f5d1be7711a6:
  Status:        running
  CreationTime:  12:00:00 25-08-2022
  TotalRequests: 200
  TotalBits:     200000
  RequestsLeft:  1000
  BitsLeft:      250000`

	rr := &Recorder{}

	outputer := output.NewOutputProcessor(true, false, "", rr)

	err := outputer.GenerateUsageOutput(status)
	require.NoError(t, err)

	require.Equal(t, strings.TrimSpace(expected), strings.TrimSpace(rr.String()))
}

func TestFailedGenerateUsageOutput(t *testing.T) {
	rr := &ErrRecorder{err: entities.Error("test error")}

	outputer := output.NewOutputProcessor(true, false, "", rr)

	err := outputer.GenerateUsageOutput(entities.UsageStatus{})
	require.EqualError(t, err, "write output: test error")
}

type Recorder struct {
	data []byte
}

func (r *Recorder) Write(b []byte) (int, error) {
	r.data = append(r.data, b...)

	return len(b), nil
}

func (r *Recorder) String() string {
	return string(r.data)
}

type ErrRecorder struct {
	err error
}

func (r *ErrRecorder) Write(b []byte) (int, error) {
	return 0, r.err
}
