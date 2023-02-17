package models

import (
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// MyString model
type MyString string

// MyStringSlice model
type MyStringSlice []string

// CombineWhitespaces combines all whitespaces into one and removes leading/trailing whitespaces
func (str MyString) CombineWhitespaces() string {
	reLeadcloseWhtsp := regexp.MustCompile(`^[\s\p{Zs}]+|[\s\p{Zs}]+$`)
	reInsideWhtsp := regexp.MustCompile(`[\s\p{Zs}]{2,}`)
	result := reLeadcloseWhtsp.ReplaceAllString(string(str), "")
	result = reInsideWhtsp.ReplaceAllString(result, " ")

	return result
}

// ToNullString invalidates a sql.NullString if empty, validates if not empty
func (str MyString) ToNullString() sql.NullString {
	s := string(str)

	return sql.NullString{String: s, Valid: s != ""}
}

// ContainsString checks string for string
func (sl MyStringSlice) ContainsString(str string) bool {
	for _, v := range sl {
		if v == str {
			return true
		}
	}

	return false
}

// UniqueString returns map of unique strings
func (sl MyStringSlice) UniqueString() []string {
	keys := make(map[string]bool)
	var list []string

	for _, item := range sl {
		if _, value := keys[item]; !value {
			keys[item] = true
			list = append(list, item)
		}
	}

	return list
}

// SplitAny split all fields
func (str MyString) SplitAny(seps string) []string {
	splitter := func(r rune) bool {
		return strings.ContainsRune(seps, r)
	}

	return strings.FieldsFunc(string(str), splitter)
}

// GetStringFromURL calls given url and returns content as string
func (str MyString) GetStringFromURL() (string, error) {
	resp, err := http.Get(string(str)) //nolint:bodyclose,gosec
	if err != nil {
		return "", fmt.Errorf("GET error: %v", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read body: %w", err)
	}

	return string(data), nil
}

// CastToString cast variable to string
func (str MyString) CastToString(x interface{}) string {
	var result string

	switch x.(type) {
	case bool:
		result = strconv.FormatBool(x.(bool))
	case int:
		result = strconv.Itoa(x.(int))
	case uint8:
		result = strconv.Itoa(int(x.(uint8)))
	case uint16:
		result = strconv.Itoa(int(x.(uint16)))
	case uint32:
		result = strconv.Itoa(int(x.(uint32)))
	case uint64:
		result = strconv.Itoa(int(x.(uint64)))
	case float32:
		result = fmt.Sprintf("%f", x.(float32))
	case float64:
		result = fmt.Sprintf("%f", x.(float64))
	case string:
		result = x.(string)
	case interface{}:
		result = x.(string)
	case reflect.Value:
		result = fmt.Sprintf("%#v", x.(reflect.Value).Interface())
	}

	return result
}
