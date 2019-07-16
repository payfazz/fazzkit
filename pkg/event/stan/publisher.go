package event

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	stan "github.com/nats-io/stan.go"
)

//MetaBuilder wraping and building meta data before send it to nats
type MetaBuilder func(data interface{}) interface{}

// Publisher wraps a URL and provides a method that implements endpoint.Endpoint.
type Publisher struct {
	publisher stan.Conn
	logger    log.Logger
}

//NewPublisher to create new Publisher
func NewPublisher(conn stan.Conn, logger log.Logger) *Publisher {

	return &Publisher{
		publisher: conn,
		logger:    logger,
	}
}

//Store for publish event (begin and commit) to nats and data wrapping as a middleware
func (p *Publisher) Store(domain, model, eventType, subject, eventSource string, f endpoint.Endpoint, metabuilder MetaBuilder) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, errResponse error) {
		var requestData map[string]interface{}
		var requestBundle = make(map[string]interface{})
		data, err := json.Marshal(request)
		if err != nil {
			return nil, err
		}
		if err = json.Unmarshal(data, &requestData); err != nil {
			return nil, err
		}

		requestBundle["domain"] = domain
		requestBundle["model"] = model
		requestBundle["status"] = "begin"
		requestBundle["event_type"] = eventType
		requestBundle["event_source"] = eventSource
		requestBundle["data"] = requestData
		dataBundle, err := json.Marshal(requestBundle)
		if err != nil {
			return nil, err
		}
		p.publisher.Publish(subject, dataBundle)
		p.logger.Log("nats", "Published message on channel: "+subject)
		p.logger.Log("nats", fmt.Sprintf("data : %s", requestBundle))

		defer func() {
			if errResponse == nil {
				var resultData map[string]interface{}
				var resultBundle = make(map[string]interface{})
				buildData := metabuilder(response)
				dataResult, err := json.Marshal(buildData)
				if err != nil {
					p.logger.Log("error_publish_commit", err)
				}

				if err = json.Unmarshal(dataResult, &resultData); err != nil {
					p.logger.Log("error_publish_commit", err)
				}

				resultBundle["domain"] = domain
				resultBundle["model"] = model
				resultBundle["status"] = "commit"
				resultBundle["event_type"] = eventType
				requestBundle["event_source"] = eventSource
				resultBundle["data"] = resultData

				dataBundle, err := json.Marshal(resultBundle)
				if err != nil {
					p.logger.Log("error_publish_commit", err)
				}
				p.publisher.Publish(subject, dataBundle)
				p.logger.Log("nats", "Published message on channel: "+subject)
				p.logger.Log("nats", fmt.Sprintf("data : %s", resultBundle))
			} else {
				var resultBundle = make(map[string]interface{})

				resultBundle["domain"] = domain
				resultBundle["model"] = model
				resultBundle["status"] = "error"
				resultBundle["event_type"] = eventType
				requestBundle["event_source"] = eventSource
				resultBundle["data"] = errResponse.Error()

				dataBundle, err := json.Marshal(resultBundle)
				if err != nil {
					p.logger.Log("error_publish_event_error", err)
				}
				p.publisher.Publish(subject, dataBundle)

				p.logger.Log("nats", "Published message on channel: "+subject)
				p.logger.Log("nats", fmt.Sprintf("data : %s", resultBundle))

			}
		}()

		return f(ctx, request)
	}
}
