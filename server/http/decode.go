package http

import (
	"context"
	"net/http"
	"reflect"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/iancoleman/strcase"

	"github.com/payfazz/fazzkit/server/common"
	"github.com/payfazz/fazzkit/server/validator"
)

//DecodeOptions executed before decode process
type DecodeOptions func(ctx context.Context, model interface{}, request *http.Request) error

//DecodeParam decode model with DecodeOptions
type DecodeParam struct {
	Model   interface{}
	Options []DecodeOptions
}

//ErrorWithStatusCode error with http status code
type ErrorWithStatusCode struct {
	err        error
	StatusCode int
}

func (e *ErrorWithStatusCode) Error() string {
	return e.err.Error()
}

//Decode generate a decode function to decode request body (json) to model
func Decode(model interface{}) func(context.Context, *http.Request) (request interface{}, err error) {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		if model == nil {
			return nil, nil
		}

		var _model interface{}
		var err error

		param, ok := model.(DecodeParam)
		if ok {
			_model, _ = common.DeepCopy(param.Model)
			for _, option := range param.Options {
				err = option(ctx, _model, r)
				if err != nil {
					return nil, &ErrorWithStatusCode{err, http.StatusUnprocessableEntity}
				}
			}
		} else {
			_model, _ = common.DeepCopy(model)
		}

		err = getURLParamUsingTag(ctx, _model, r)
		if err != nil {
			return nil, &ErrorWithStatusCode{err, http.StatusUnprocessableEntity}
		}

		contentType := r.Header["Content-Type"]

		if common.StringInSlice("application/json", contentType) {
			_model, err = ParseJSON(ctx, r, _model)
			if err != nil {
				return nil, &ErrorWithStatusCode{err, http.StatusUnprocessableEntity}
			}
		}

		err = validator.DefaultValidator()(_model)
		if err != nil {
			return nil, &ErrorWithStatusCode{err, http.StatusUnprocessableEntity}
		}

		return _model, nil
	}
}

//GetURLParam built-in DecodeOptions for decode using url params
func GetURLParam(params []string) DecodeOptions {
	return func(ctx context.Context, model interface{}, r *http.Request) error {
		var err error
		typ := reflect.TypeOf(model).Elem()
		for i := 0; i < typ.NumField(); i++ {
			name := typ.Field(i).Name
			name = strcase.ToSnake(name)
			if common.StringInSlice(name, params) {
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
