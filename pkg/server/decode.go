package server

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/iancoleman/strcase"
)

//GRPCDecodeOptions executed before decode process
type GRPCDecodeOptions func(ctx context.Context, model interface{}, request interface{}) error

//GRPCDecodeParam decode model with GRPCDecodeOptions
type GRPCDecodeParam struct {
	Model   interface{}
	Options []GRPCDecodeOptions
}

//DecodeGRPC generate a decode function to decode proto message to model
func (e *Endpoint) DecodeGRPC(model interface{}) func(context.Context, interface{}) (request interface{}, err error) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if model == nil {
			return nil, nil
		}

		var _model interface{}
		var err error

		param, ok := model.(GRPCDecodeParam)
		if ok {
			_model, _ = deepCopy(param.Model)
			for _, option := range param.Options {
				err = option(ctx, _model, request)
				if err != nil {
					return nil, err
				}
			}
		} else {
			_model, _ = deepCopy(model)
		}

		_model, err = e.ParseGRPC(ctx, request, _model)

		if err != nil {
			return nil, err
		}

		err = e.Validate(_model)

		if err != nil {
			return nil, err
		}

		return _model, nil
	}
}

//HTTPDecodeOptions executed before decode process
type HTTPDecodeOptions func(ctx context.Context, model interface{}, request *http.Request) error

//HTTPDecodeParam decode model with HTTPDecodeOptions
type HTTPDecodeParam struct {
	Model   interface{}
	Options []HTTPDecodeOptions
}

//ErrorWithStatusCode error with http status code
type ErrorWithStatusCode struct {
	err        error
	statusCode int
}

func (e *ErrorWithStatusCode) Error() string {
	return e.err.Error()
}

//DecodeHTTP generate a decode function to decode request body (json) to model
func (e *Endpoint) DecodeHTTP(model interface{}) func(context.Context, *http.Request) (request interface{}, err error) {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		if model == nil {
			return nil, nil
		}

		var _model interface{}
		var err error

		param, ok := model.(HTTPDecodeParam)
		if ok {
			_model, _ = deepCopy(param.Model)
			for _, option := range param.Options {
				err = option(ctx, _model, r)
				if err != nil {
					return nil, &ErrorWithStatusCode{err, http.StatusUnprocessableEntity}
				}
			}
		} else {
			_model, _ = deepCopy(model)
		}

		err = getURLParamUsingTag(ctx, _model, r)
		if err != nil {
			return nil, &ErrorWithStatusCode{err, http.StatusUnprocessableEntity}
		}

		contentType := r.Header["Content-Type"]

		if stringInSlice("application/json", contentType) {
			_model, err = e.ParseHTTPJson(ctx, r, _model)
			if err != nil {
				return nil, &ErrorWithStatusCode{err, http.StatusUnprocessableEntity}
			}
		}

		err = e.Validate(_model)
		if err != nil {
			return nil, &ErrorWithStatusCode{err, http.StatusUnprocessableEntity}
		}

		return _model, nil
	}
}

func deepCopy(v interface{}) (interface{}, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	vptr := reflect.New(reflect.TypeOf(v))
	err = json.Unmarshal(data, vptr.Interface())
	if err != nil {
		return nil, err
	}
	return vptr.Elem().Interface(), err
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

//GetURLParam built-in HTTPDecodeOptions for decode using url params
func GetURLParam(params []string) HTTPDecodeOptions {
	return func(ctx context.Context, model interface{}, r *http.Request) error {
		var err error
		typ := reflect.TypeOf(model).Elem()
		for i := 0; i < typ.NumField(); i++ {
			name := typ.Field(i).Name
			name = strcase.ToSnake(name)
			if stringInSlice(name, params) {
				err = getURLParam(ctx, model, r, name, i)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}
}

func getURLParamUsingTag(ctx context.Context, model interface{}, r *http.Request) error {
	var err error
	typ := reflect.TypeOf(model).Elem()
	for i := 0; i < typ.NumField(); i++ {
		tag := typ.Field(i).Tag.Get("httpurl")
		if tag != "" {
			err = getURLParam(ctx, model, r, tag, i)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func getURLParam(ctx context.Context, model interface{}, r *http.Request, param string, valIdx int) error {
	value := chi.URLParam(r, param)
	if value == "" {
		return nil
	}

	val := reflect.ValueOf(model).Elem()

	switch valtype := val.Field(valIdx).Type().String(); valtype {
	case "string":
		val.Field(valIdx).SetString(value)
	case "int64":
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		val.Field(valIdx).SetInt(v)
	case "uuid.UUID":
		v, err := uuid.Parse(value)
		if err != nil {
			return err
		}
		val.Field(valIdx).Set(reflect.ValueOf(v))
	}

	return nil
}
