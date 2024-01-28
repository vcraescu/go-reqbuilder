package reqbuilder

import (
	"encoding/json"

	"github.com/vcraescu/go-urlvalues"
)

type builderOptions struct {
	bodyMarshaler      BodyMarshaler
	urlValuesMarshaler URLValuesMarshaler
}

type BuilderOption interface {
	apply(cfg *builderOptions)
}

var _ BuilderOption = builderOptionFunc(nil)

type builderOptionFunc func(opts *builderOptions)

func (fn builderOptionFunc) apply(cfg *builderOptions) {
	fn(cfg)
}

func WithBodyMarshaler(marshaler BodyMarshaler) BuilderOption {
	return builderOptionFunc(func(opts *builderOptions) {
		opts.bodyMarshaler = marshaler
	})
}

func WithURLValuesMarshaler(marshaler URLValuesMarshaler) BuilderOption {
	return builderOptionFunc(func(opts *builderOptions) {
		opts.urlValuesMarshaler = marshaler
	})
}

func defaultBuilderOptions(opts ...BuilderOption) *builderOptions {
	options := &builderOptions{
		bodyMarshaler:      BodyMarshalerFunc(json.Marshal),
		urlValuesMarshaler: URLValuesMarshalerFunc(urlvalues.Marshal),
	}

	for _, opt := range opts {
		opt.apply(options)
	}

	return options
}
