// Code is generated. Do not edit

package there

import (
	"strconv"
	"strings"
)

func (reader MapReader) GetInt(key string) (int, error) {
	list, ok := reader.GetSlice(key)
	if !ok {
		return 0, ErrorParameterNotPresent
	}
	return strconv.Atoi(list[0])
}

func (reader MapReader) GetIntDefault(key string, defaultValue int) int {
	list, ok := reader.GetSlice(key)
	if !ok {
		return defaultValue
	}
	v, err := strconv.Atoi(list[0])
	if err != nil {
		return defaultValue
	}
	return v
}

func (reader MapReader) GetIntSlice(key, delimiter string) ([]int, error) {
	list, ok := reader.GetSlice(key)
	if !ok {
		return nil, ErrorParameterNotPresent
	}
	split := strings.Split(list[0], delimiter)
	ints := make([]int, len(split))
	var err error
	for i, s := range split {
		ints[i], err = strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
	}
	return ints, nil
}

func (reader MapReader) GetIntSliceDefault(key, delimiter string, defaultValue []int) []int {
	list, ok := reader.GetSlice(key)
	if !ok {
		return defaultValue
	}
	split := strings.Split(list[0], delimiter)
	ints := make([]int, len(split))
	var err error
	for i, s := range split {
		ints[i], err = strconv.Atoi(s)
		if err != nil {
			return defaultValue
		}
	}
	return ints
}
