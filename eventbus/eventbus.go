package eventbus

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	"bitbucket.org/goreorto/sqlhero/logger"
)

var jobTimeout = time.Minute

// HandlerFunc defines a func that is capable of handling an Event
type HandlerFunc func(context.Context, *Event) (interface{}, error)

// Event for the EventBus
type Event struct {
	keep      bool
	name      string
	startTime time.Time
	endTime   time.Time
	done      chan interface{}
	fail      chan error

	Payload interface{}
}

// NewEvent create a new event of `name` type with a payload.
// `keep` events and trigger inmediately after handler registration.
// Only the last event keep is triggered
func NewEvent(name string, payload interface{}, keep bool) *Event {
	return &Event{
		keep:    keep,
		name:    name,
		done:    make(chan interface{}, 1),
		fail:    make(chan error, 1),
		Payload: payload,
	}
}

// Done returns a chan that returns the event result
func (e *Event) Done() chan interface{} {
	return e.done
}

// Fail returns a chan that has an error if the event fails
func (e *Event) Fail() chan error {
	return e.fail
}

func (e *Event) String() string {
	return fmt.Sprintf("name:%s,keep:%t,payload:%s", e.name, e.keep, stringify(e.Payload))
}

func stringify(p interface{}) string {
	switch v := p.(type) {
	case bool:
		return fmt.Sprintf("%t", v)
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%f", v)
	case string:
		return v
	default:
		valueOf := reflect.ValueOf(v)
		if valueOf.Kind() == reflect.Ptr {
			if !valueOf.CanSet() {
				return valueOf.Type().String() + "{}"
			}
			valueOf = valueOf.Elem()
		}
		if valueOf.Kind() == reflect.Struct {
			s := valueOf.Type().String() + "{"
			for i := 0; i < valueOf.NumField(); i++ {
				f := valueOf.Field(i)
				s += stringify(f.Interface())
			}
			s += "}"
			return s
		}

		return fmt.Sprintf("%+v", v)
	}
}

// EventBus is a event bus for issue event and background processes
type EventBus struct {
	workers    chan *worker
	handlers   map[string][]HandlerFunc
	m          sync.Mutex
	keptEvents map[string]*Event
	log        logger.Logger
}

// New create a new EventBus
func New(maxPoolSize int, log logger.Logger) *EventBus {
	p := &EventBus{
		workers:    make(chan *worker, maxPoolSize),
		keptEvents: map[string]*Event{},
		log:        log.SetComponent("BUS"),
	}
	p.handlers = make(map[string][]HandlerFunc)

	for i := cap(p.workers); i > 0; i-- {
		p.workers <- &worker{ID: i, pipe: p}
	}

	return p
}

// RegisterHandler register a func to handle a event type
func (p *EventBus) RegisterHandler(name string, h HandlerFunc) {
	p.m.Lock()
	defer p.m.Unlock()

	if _, ok := p.handlers[name]; !ok {
		p.handlers[name] = make([]HandlerFunc, 0)
	}

	p.handlers[name] = append(p.handlers[name], h)

	if e, ok := p.keptEvents[name]; ok {
		h(context.Background(), e)
	}
}

// Emit issues an event for processing
func (p *EventBus) Emit(e *Event) {
	go func() {
		p.log.WithField("event", e).Debug("event received, waiting for worker")
		w := <-p.workers
		w.Do(e)

		if e.keep {
			p.keptEvents[e.name] = e
		}
	}()
}

type worker struct {
	ID           int
	LastJobEndAt time.Time
	pipe         *EventBus
}

// Do does
func (w *worker) Do(e *Event) {
	e.startTime = time.Now()
	defer close(e.done)
	defer close(e.fail)

	handlers, ok := w.pipe.handlers[e.name]
	if ok {
		w.pipe.log.WithField("event", e).Debugf("processing")
		for _, h := range handlers {
			ctx, cancel := context.WithTimeout(context.Background(), jobTimeout)
			defer cancel()

			result, err := h(ctx, e)
			e.endTime = time.Now()

			if err != nil {
				w.pipe.log.WithField("err", err).Debug("failed processing")
				e.fail <- err
				w.pipe.workers <- w
				return
			}

			e.done <- result
		}
	}
	if !ok {
		w.pipe.log.WithField("event", e).Debugf("no handler")
		e.fail <- fmt.Errorf("no handler defined for %s", e.name)
	}

	w.LastJobEndAt = time.Now()
	w.pipe.workers <- w
}
