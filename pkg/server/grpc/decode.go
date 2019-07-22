package grpc

import (
	"context"

	"github.com/payfazz/fazzkit/pkg/server/common"
	"github.com/payfazz/fazzkit/pkg/server/validator"
)

//DecodeOptions executed before decode process
type DecodeOptions func(ctx context.Context, model interface{}, request interface{}) error

//DecodeParam decode model with DecodeOptions
type DecodeParam struct {
	Model   interface{}
	Options []DecodeOptions
}

//Decode generate a decode function to decode proto message to model
func Decode(model interface{}) func(context.Context, interface{}) (request interface{}, err error) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if model == nil {
			return nil, nil
		}

		var _model interface{}
		var err error

		param, ok := model.(DecodeParam)
		if ok {
			_model, _ = common.DeepCopy(param.Model)
			for _, option := range param.Options {
				err = option(ctx, _model, request)
				if err != nil {
					return nil, err
				}
			}
		} else {
			_model, _ = common.DeepCopy(model)
		}

		_model, err = Parse(ctx, request, _model)

		if err != nil {
			return nil, err
		}

		err = validator.DefaultValidator()(_model)

		if err != nil {
			return nil, err
		}

		return _model, nil
	}
}
