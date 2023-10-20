package actions

import (
	"reflect"
	"testing"
	"time"

	"github.com/freehandle/breeze/crypto"
)

var (
	event = &CreateEvent{
		Epoch:        21,
		Author:       crypto.Token{},
		Reasons:      "create event test",
		OnBehalfOf:   "first_collective",
		StartAt:      time.Date(2021, 8, 15, 14, 30, 45, 100, time.Local),
		EstimatedEnd: time.Date(2022, 8, 15, 14, 30, 45, 100, time.Local),
		Description:  "create event test",
		Venue:        "first_venue",
		Open:         true,
		Public:       true,
		Managers:     []crypto.Token{},
	}

	cancel = &CancelEvent{
		Epoch:   22,
		Author:  crypto.Token{},
		Reasons: "test cancel event",
		Hash:    crypto.Hash{},
	}

	updtSTr = "test update event"

	uEvent = &UpdateEvent{
		Epoch:       23,
		Author:      crypto.Token{},
		Reasons:     "test update event",
		EventHash:   crypto.Hash{},
		Description: &updtSTr,
		Venue:       &updtSTr,
		Open:        nil,
		Public:      nil,
		Managers:    nil,
	}
)

func TestCreateEvent(t *testing.T) {
	e := ParseCreateEvent(event.Serialize())
	if e == nil {
		t.Error("Could not parse actions CreateEvent")
		return
	}
	if !reflect.DeepEqual(e, event) {
		t.Error("Parse and Serialize not working for actions CreateEvent")
	}
}

func TestCancelEvent(t *testing.T) {
	c := ParseCancelEvent(cancel.Serialize())
	if c == nil {
		t.Error("Could not parse actions CancelEvent")
		return
	}
	if !reflect.DeepEqual(c, cancel) {
		t.Error("Parse and Serialize not working for actions CancelEvent")
	}
}

func TestUpdateEvent(t *testing.T) {
	u := ParseUpdateEvent(uEvent.Serialize())
	if u == nil {
		t.Error("Could not parse actions UpdateEvent")
		return
	}
	if !reflect.DeepEqual(u, uEvent) {
		t.Error("Parse and Serialize not working for actions UpdateEvent")
	}
}
