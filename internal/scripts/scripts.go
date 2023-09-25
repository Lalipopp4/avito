package scripts

import (
	"encoding/json"
	"io"
)

func Decode(d interface{}, r io.Reader) error {
	err := json.NewDecoder(r).Decode(&d)
	if err != nil {
		return err
	}
	return nil
}
