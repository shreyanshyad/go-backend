package models

import (
	"encoding/json"
	"io"

	"github.com/google/uuid"
)

type View struct {
	ID          uuid.UUID `json:"id"`
	DashID      uuid.UUID `json:"dash_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   string    `json:"-"`
}

func (v *View) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(v)
}

func (v *View) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(v)
}
