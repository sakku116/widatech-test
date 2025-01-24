package model

import "github.com/op/go-logging"

var logger = logging.MustGetLogger("model")

type ModelProps struct {
	QueriableFields []string
	SortableFields  []string
}

type IModel interface {
	GetProps() ModelProps
}
