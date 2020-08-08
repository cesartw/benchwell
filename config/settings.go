package config

import (
	"strconv"
	"strings"
	"sync"
)

type Setting struct {
	id    int64
	name  string
	value string

	m    sync.Mutex
	l    uint
	subs map[uint]*SettingUpdater
}

type SettingUpdater struct {
	p *Setting
	l uint
	f func(interface{})
	c chan interface{}
}

func (s *Setting) notify() {
	for _, updater := range s.subs {
		go func(c chan interface{}) {
			c <- s.value
		}(updater.c)
	}
}

func (s *SettingUpdater) Unsubscribe() {
	close(s.c)
	s.p.unsubscribe(s.l)
}

func (s *Setting) SetBool(b bool) {
	if b {
		s.value = "1"
	} else {
		s.value = "0"
	}
	s.notify()
}

func (s *Setting) SetString(v string) {
	s.value = v
	s.notify()
}

func (s *Setting) Bool() bool {
	return s.value == "1" || strings.EqualFold(s.value, "1")
}

func (s *Setting) String() string {
	return s.value
}

func (s *Setting) Int() int {
	i, _ := strconv.ParseInt(s.value, 10, 64)
	return int(i)
}

func (s *Setting) Int64() int64 {
	i, _ := strconv.ParseInt(s.value, 10, 64)
	return i
}

func (s *Setting) unsubscribe(l uint) {
	s.m.Lock()
	defer s.m.Unlock()

	delete(s.subs, l)
}

func (s *Setting) Subscribe(f func(interface{})) *SettingUpdater {
	s.m.Lock()
	defer s.m.Unlock()

	s.l++
	u := &SettingUpdater{
		l: s.l,
		f: f,
		c: make(chan interface{}, 1),
	}

	go func() {
		for {
			select {
			case v := <-u.c:
				f(v)
			}
		}
	}()

	s.subs[s.l] = u
	return u
}

func (p *Setting) Settinglisher(v interface{}) {
	p.m.Lock()
	defer p.m.Unlock()

	for _, s := range p.subs {
		s.c <- v
	}
}
