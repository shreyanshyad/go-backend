package models

import (
	"encoding/json"
	"io"

	"github.com/google/uuid"
)

type Dash struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Views       []*View   `json:"views,omitempty"`
	CreatedAt   string    `json:"-"`
}

func (d *Dash) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(d)
}

func (d *Dash) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(d)
}
