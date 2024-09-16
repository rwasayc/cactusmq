package packet

import (
	"encoding/json"
	"fmt"
)

func JSON(v any) string {
	bytes, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Errorf("failed to marshal(%+v) to json: %w", v, err))
	}
	return string(bytes)
}
