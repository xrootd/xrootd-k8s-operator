package utils

import (
	"errors"
	"strconv"
	"testing"
)

func TestMapIntToString(t *testing.T) {
	input := []int{1, 31, 93}
	t.Errorf("Input: %v", input)
	transformer := func(val int) string {
		return strconv.FormatInt(int64(val), 2)
	}
	output := Map(transformer, input).([]string)
	t.Errorf("Output: %v", output)
}

func TestMapWithErrorIntToString(t *testing.T) {
	input := []int{1, 31, 0}
	t.Errorf("Input: %v", input)
	transformer := func(val int) (string, error) {
		if val == 0 {
			return "", errors.New("Input value is 0")
		}
		return strconv.FormatInt(int64(val), 2), nil
	}
	output, err := MapWithError(transformer, input)
	t.Errorf("Output: %v, Error: %v", output, err)
}
