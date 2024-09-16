package packet

import "encoding/json"

type UserProperty struct {
	Key string `json:"k"`
	Val string `json:"v"`
}

// value with flag
type FlagV[T any] struct {
	V *T `json:"v"`
}

func NewFlagV[T any](value T) FlagV[T] {
	return FlagV[T]{V: &value}
}

func NewNoFlagV[T any]() FlagV[T] {
	return FlagV[T]{}
}

func (fv FlagV[T]) Flag() bool {
	return fv.V != nil
}
func (fv FlagV[T]) Value() T {
	if fv.Flag() {
		return *fv.V
	}
	return *new(T)
}

func (fv FlagV[T]) MarshalJSON() ([]byte, error) {
	if fv.Flag() {
		return json.Marshal(fv.Value())
	}
	return json.Marshal(nil)
}

// Password
type Password struct {
	FlagV[[]byte]
}

func (p Password) String() string {
	if p.Flag() {
		return string(*p.V)
	}
	return ""
}

func NewPassword(value []byte) Password {
	return Password{FlagV[[]byte]{V: &value}}
}

func NewSPassword(value string) Password {
	pwd := []byte(value)
	return Password{FlagV[[]byte]{V: &pwd}}
}
