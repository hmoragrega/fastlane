package httpapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const errBody = `{"error":%q}`

func encode(w http.ResponseWriter, result interface{}, err error) {
	if err != nil {
		errResponse(w, err)
		return
	}
	encodeData(w, result)
}

func encodeData(w http.ResponseWriter, data interface{}) {
	var (
		buf []byte
		err error
	)

	status := http.StatusNoContent
	if data != nil {
		buf, err = json.Marshal(data)
		if err != nil {
			errResponse(w, fmt.Errorf("cannot marshal data %+v: %v", data, err))
			return
		}

		status = http.StatusOK
	}

	response(w, status, buf)
}

func response(w http.ResponseWriter, status int, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if len(data) > 0 {
		_, _ = w.Write(data)
	}
}

func errResponse(w http.ResponseWriter, err error) {
	response(w, http.StatusInternalServerError, encodeErr(err).Bytes())
}

func encodeErr(err error) *bytes.Buffer {
	return bytes.NewBufferString(fmt.Sprintf(errBody, err))
}
