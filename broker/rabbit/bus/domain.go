package bus

import (
	"context"
	"github.com/boostgo/lite/log"
)

/*
	Message Bus (or Event bus) is a combination of a common data model, a common command set,
	and a messaging infrastructure to allow different systems to communicate through a shared set of interfaces.

	There are 2 instances for Message Bus pattern: Dispatcher and Listener
*/

type ListenerAction func(event EventContext) error

// Dispatcher event sender by "message bus" pattern (read above).
//
// To create dispatcher use boost.NewDispatcher(MessageBusConnector)
type Dispatcher interface {
	// Dispatch sends event
	Dispatch(ctx context.Context, event any) error
}

// Listener listen for events by "message bus" pattern (read above).
//
// To create listener call app.Listener()
type Listener interface {
	// Run Listener
	Run() error
	// Bind bind Listener handler for provided event object
	Bind(event any, action func(ctx EventContext) error)
	// EventsCount returns count of events.
	// Need for "greeting" text
	EventsCount() int
	// Close closes Listener connection.
	// If Listener was called by di package (Dependency Injection) Close will be called automatically
	Close() error
}

// EventContext context object which comes to Listener handler
type EventContext interface {
	// Body returns incoming event in bytes (JSON)
	Body() []byte
	// Parse incoming event body (JSON) to provided "export"
	Parse(export any) error
	// Context returns context.Context.
	// If for listener set up setting "msgbus_listener__timeout" (int) context will wrap with WithTimeout
	Context() context.Context
	// Cancel returns context cancel method
	Cancel()
	// TraceID returns trace id if it is there in headers
	TraceID() string
	// Logger returns logger with namespace
	Logger() log.Logger
}

type Event struct {
	Name   string
	Action ListenerAction
	Object any
}
