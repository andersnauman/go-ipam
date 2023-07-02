package misc

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

func PopString(s []string) (ss []string, err error) {
	if len(s) < 2 {
		err = errors.New("not enough values to parse")
		return
	}
	ss = s[1:]
	return
}

type ErrorStatus struct {
	Error []ErrorFields `json:"errors"`
}

type ErrorFields struct {
	Status  string `json:"status,omitempty"`
	Details string `json:"details,omitempty"`
}

func ReturnError(w http.ResponseWriter, errorCode int, msg string) (err error) {
	errMsg := ErrorStatus{
		Error: []ErrorFields{
			{
				Status:  strconv.Itoa(errorCode),
				Details: msg,
			},
		},
	}
	var data []byte
	data, err = json.Marshal(&errMsg)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(data)
	return
}
