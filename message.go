package easy_amqp

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"strconv"
)

type message struct {
	raw         interface{}
	ContentType ContentType
	Exchange    string
	RouterKey   string
	UUID        string
	Type        string
	AppId       string
	Mode        uint8
	Priority    uint8
}

var (
	ErrorInvalidContentType = "invalid content type"
)

func (m *message) Body() (bytes []byte, err error) {
	switch m.ContentType {
	case Json:
		bytes, err = json.Marshal(m.raw)
		return
	case Text:
		defer func() {
			e := recover()
			if e != nil {
				err = e.(error)
			}
		}()
		bytes = []byte(fmt.Sprint(m.raw))
		return
	}
	return nil, ErrorInvalidContentType
}

type ContentType uint8

const (
	Json ContentType = iota
	Text
)

var contentTypes = map[ContentType]string{Json: "application/json", Text: "application/text"}

func (c ContentType) String() string {
	if s, ok := contentTypes[c]; ok {
		return s
	}
	return strconv.FormatUint(uint64(c), 10)
}

func WithContentType(contentType ContentType) MessageOption {
	return func(i *message) {
		i.ContentType = contentType
	}
}

func WithExchange(exchange string) MessageOption {
	return func(i *message) {
		i.Exchange = exchange
	}
}

func WithRouterKey(routerkey string) MessageOption {
	return func(i *message) {
		i.RouterKey = routerkey
	}
}

func WithMode(mode uint8) MessageOption {
	return func(i *message) {
		i.Mode = mode
	}
}

func WithAppId(appId string) MessageOption {
	return func(i *message) {
		i.AppId = appId
	}
}

func WithType(tp string) MessageOption {
	return func(i *message) {
		i.Type = tp
	}
}

type MessageOption func(*message)

func NewMessage(raw interface{}, contentType ContentType, opts ...MessageOption) *message {
	m := message{UUID: uuid.New().String(), Mode: amqp.Persistent}
	for _, o := range opts {
		o(&m)
	}
	return &m
}
