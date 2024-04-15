package pretender

import (
	"encoding/json"
	"time"
)

type response struct {
	StatusCode uint              `json:"status_code"`
	Body       string            `json:"body"`
	Headers    map[string]string `json:"headers"`
	Delay      time.Duration     `json:"delay_ms"`
}

func (r *response) UnmarshalJSON(data []byte) error {
	type alias response
	tmp := (*alias)(r)

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	if r.StatusCode == 0 {
		r.StatusCode = 200
	}

	r.Delay = time.Duration(r.Delay) * time.Millisecond

	return nil
}
