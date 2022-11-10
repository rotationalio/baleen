package baleen

import (
	"errors"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/rotationalio/watermill-ensign/pkg/ensign"
)

var (
	ErrUnhandledType = errors.New("ensign type not handled")
	ErrUnhandledMIME = errors.New("ensign mimetype not handled")
)

func TypeFilter(mime string, etypes ...string) message.HandlerMiddleware {
	typeFilter := make(map[string]struct{}, len(etypes))
	for _, etype := range etypes {
		typeFilter[etype] = struct{}{}
	}

	return func(h message.HandlerFunc) message.HandlerFunc {
		return func(msg *message.Message) ([]*message.Message, error) {
			if _, ok := typeFilter[msg.Metadata.Get(ensign.TypeNameKey)]; !ok {
				// TODO: when ensign has topics return ErrUnhandledType
				// return nil, ErrUnhandledType

				// HACK: to prevent tons of error logs we're just returning nil.
				msg.Nack()
				return nil, nil
			}

			if msg.Metadata.Get(ensign.MIMEKey) != mime {
				return nil, ErrUnhandledMIME
			}

			return h(msg)
		}
	}
}
