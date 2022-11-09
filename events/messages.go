package events

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	mime "github.com/rotationalio/ensign/pkg/mimetype/v1beta1"
	"github.com/rotationalio/watermill-ensign/pkg/ensign"
)

var mimetype = mime.ApplicationMsgPack.MimeType()

func Marshal(event TypedEvent, uuid string) (msg *message.Message, err error) {
	// Marshal the message into msgpack bytes
	var payload []byte
	if payload, err = event.MarshalMsg(nil); err != nil {
		return nil, err
	}

	// Get the type and version of the event
	etype := event.Type()

	// Create the watermill message and set the ensign metadata
	msg = message.NewMessage(uuid, payload)
	msg.Metadata.Set(ensign.MIMEKey, mimetype)
	msg.Metadata.Set(ensign.TypeNameKey, etype.Name)
	msg.Metadata.Set(ensign.TypeVersionKey, strconv.FormatUint(uint64(etype.Version), 10))
	msg.Metadata.Set(ensign.CreatedKey, time.Now().Format(time.RFC3339Nano))

	return msg, nil
}

func Unmarshal(msg *message.Message) (_ TypedEvent, err error) {
	switch t := msg.Metadata.Get(ensign.TypeNameKey); t {
	case TypeSubscription:
		sub := &Subscription{}
		if _, err = sub.UnmarshalMsg(msg.Payload); err != nil {
			return nil, fmt.Errorf("cannot unmarshal %s: %w", t, err)
		}
		return sub, nil
	case TypeFeedSync:
		fs := &FeedSync{}
		if _, err = fs.UnmarshalMsg(msg.Payload); err != nil {
			return nil, fmt.Errorf("cannot unmarshal %s: %w", t, err)
		}
		return fs, nil
	case TypeFeedItem:
		fi := &FeedItem{}
		if _, err = fi.UnmarshalMsg(msg.Payload); err != nil {
			return nil, fmt.Errorf("cannot unmarshal %s: %w", t, err)
		}
		return fi, nil
	case TypeDocument:
		doc := &Document{}
		if _, err = doc.UnmarshalMsg(msg.Payload); err != nil {
			return nil, fmt.Errorf("cannot unmarshal %s: %w", t, err)
		}
		return doc, nil
	default:
		return nil, fmt.Errorf("cannot unmarshal message type %q", t)
	}
}

func UnmarshalSubscription(msg *message.Message) (e *Subscription, err error) {
	var event TypedEvent
	if event, err = Unmarshal(msg); err != nil {
		return nil, err
	}

	var ok bool
	if e, ok = event.(*Subscription); !ok {
		return nil, errors.New("message does not contain a Subscription event")
	}
	return e, nil
}

func UnmarshalFeedSync(msg *message.Message) (e *FeedSync, err error) {
	var event TypedEvent
	if event, err = Unmarshal(msg); err != nil {
		return nil, err
	}

	var ok bool
	if e, ok = event.(*FeedSync); !ok {
		return nil, errors.New("message does not contain a FeedSync event")
	}
	return e, nil
}

func UnmarshalFeedItem(msg *message.Message) (e *FeedItem, err error) {
	var event TypedEvent
	if event, err = Unmarshal(msg); err != nil {
		return nil, err
	}

	var ok bool
	if e, ok = event.(*FeedItem); !ok {
		return nil, errors.New("message does not contain a FeedItem event")
	}
	return e, nil
}

func UnmarshalDocument(msg *message.Message) (e *Document, err error) {
	var event TypedEvent
	if event, err = Unmarshal(msg); err != nil {
		return nil, err
	}

	var ok bool
	if e, ok = event.(*Document); !ok {
		return nil, errors.New("message does not contain a Document event")
	}
	return e, nil
}
