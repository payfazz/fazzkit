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
type GRPCDecodeOptions func(ctx context.Context, model interface{}, request interface{})

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
				option(ctx, _model, request)
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
type HTTPDecodeOptions func(ctx context.Context, model interface{}, request *http.Request)

//HTTPDecodeParam decode model with HTTPDecodeOptions
type HTTPDecodeParam struct {
	Model   interface{}
	Options []HTTPDecodeOptions
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
				option(ctx, _model, r)
			}
		} else {
			_model, _ = deepCopy(model)
		}

		getURLParamUsingTag(ctx, _model, r)

		contentType := r.Header["Content-Type"]

		if stringInSlice("application/json", contentType) {
			_model, err = e.ParseHTTPJson(ctx, r, _model)
			if err != nil {
				return nil, err
			}
		}

		err = e.Validate(_model)
		if err != nil {
			return nil, err
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
	return func(ctx context.Context, model interface{}, r *http.Request) {
		getURLParam(ctx, model, r, params)
	}
}

func getURLParamUsingTag(ctx context.Context, model interface{}, r *http.Request) {
	params := []string{}

	typ := reflect.TypeOf(model).Elem()
	for i := 0; i < typ.NumField(); i++ {
		params = append(params, typ.Field(i).Tag.Get("httpurl"))
	}

	getURLParam(ctx, model, r, params)
}

func getURLParam(ctx context.Context, model interface{}, r *http.Request, params []string) {
	typ := reflect.TypeOf(model).Elem()
	val := reflect.ValueOf(model).Elem()

	for i := 0; i < typ.NumField(); i++ {
		name := typ.Field(i).Name
		name = strcase.ToSnake(name)

		value := chi.URLParam(r, name)

		if value == "" {
			continue
		}
		if !stringInSlice(name, params) {
			continue
		}

		valtype := val.Field(i).Type().String()

		switch valtype {
		case "string":
			val.Field(i).SetString(value)
		case "int64":
			v, _ := strconv.ParseInt(value, 10, 64)
			val.Field(i).SetInt(v)
		case "uuid.UUID":
			v, _ := uuid.Parse(value)
			val.Field(i).Set(reflect.ValueOf(v))
		}
	}
}
