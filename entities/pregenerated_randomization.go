package entities

import "encoding/json"

type PregenRand struct {
	ID   *string
	Date *string
}

func (pr PregenRand) MarshalJSON() ([]byte, error) {
	const (
		nullString             = "null"
		errBothValuesSpecified = Error("only one of date or id is allowed")
	)

	if pr.ID == nil && pr.Date == nil {
		return []byte(nullString), nil
	}

	if pr.ID != nil && pr.Date != nil {
		return nil, errBothValuesSpecified
	}

	var s string

	if pr.ID != nil {
		s = *pr.ID
	}

	if pr.Date != nil {
		s = *pr.Date
	}

	return json.Marshal(s)
}
