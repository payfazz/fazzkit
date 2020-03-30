package http

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"net/http"
)

type err interface {
	error() error
}

type EncodeFunc func() func(ctx context.Context, w http.ResponseWriter, response interface{}) error

//Encode generate a encode function to encode response to json
func Encode() func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
		if e, ok := response.(err); ok && e.error() != nil {
			return e.error()
		}

		if response == nil {
			w.WriteHeader(http.StatusNoContent)
			json.NewEncoder(w).Encode("")
			return nil
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(response)
		return nil
	}
}

type CSVResponse struct {
	Filename string
	Data     [][]string
}

//EncodeCSV generate a encode function to encode csv in attachment
func EncodeCSV() func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
		if e, ok := response.(err); ok && e.error() != nil {
			return e.error()
		}

		if response == nil {
			w.WriteHeader(http.StatusNoContent)
			_ = json.NewEncoder(w).Encode("")
			return nil
		}

		resp, ok := response.(CSVResponse)
		if !ok {
			return errors.New("response is not CSVResponse")
		}

		w.Header().Set("Content-Disposition", "attachment; filename="+resp.Filename)
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Transfer-Encoding", "chunked")

		writer := csv.NewWriter(w)
		err := writer.WriteAll(resp.Data)
		writer.Flush()

		if nil != err {
			return err
		}

		return nil
	}
}
