package orb

import (
	"fmt"
	"time"
)

type PerformerId string

type Word string

const (
	Done Word = "---"
)

type Slot int

const (
	Left Slot = iota
	Right
)

type Address struct {
	Performer PerformerId
	Slot      Slot
}

type AddressSet map[Address]bool

func (p AddressSet) Add(id Address) {
	p[id] = true
}

func (p AddressSet) Remove(id Address) {
	if p.Contains(id) {
		delete(p, id)
	}
}

func (p AddressSet) Slice() []Address {
	var out []Address
	for p, exists := range p {
		if exists {
			out = append(out, p)
		}
	}
	return out
}

func (p AddressSet) Contains(id Address) bool {
	_, found := p[id]
	return found
}

// PerformerRouter is some mapping of performer id to a slice of performer ids,
// e.g. an index of performer id to that performer's listeners
type PerformerRouter map[PerformerId]AddressSet

func (r PerformerRouter) Route(m SentMessage) []*RoutedMessage {
	receivers, found := r[m.Sender]
	if !found {
		return []*RoutedMessage{}
	}

	messages := make([]*RoutedMessage, len(receivers))
	for i, receiver := range receivers.Slice() {
		messages[i] = &RoutedMessage{
			Address: receiver,
			Message: m.Message,
		}
	}

	return messages
}

type SentMessage struct {
	Sender  PerformerId
	Message Word
}

type RoutedMessage struct {
	Address Address
	Message Word
}

type ReceivedWord struct {
	Word Word
	Time time.Time
}

type Mailbox struct {
	Left  *ReceivedWord
	Right *ReceivedWord
}

func (m *Mailbox) Receive(msg RoutedMessage) {
	switch msg.Address.Slot {
	case Left:
		m.Left = &ReceivedWord{msg.Message, time.Now()}
		break
	case Right:
		m.Right = &ReceivedWord{msg.Message, time.Now()}
		break
	}
}

func (m *Mailbox) Empty() bool {
	return m.Left == nil && m.Right == nil
}

func (m *Mailbox) Clear() {
	m.Left = nil
	m.Right = nil
}

type MailRoom map[PerformerId]*Mailbox

type Script []Word

func (s Script) Current() Word {
	if s.Done() {
		return Done
	}

	return s[0]
}

func (s *Script) Advance() {
	if !s.Done() {
		*s = (*s)[1:]
	}
}

func (s Script) Done() bool {
	return len(s) <= 0
}

func (s Script) Copy() *Script {
	out := make(Script, len(s))
	for i, s := range s {
		out[i] = s
	}
	return &out
}

type ScriptDB map[PerformerId]*Script

func (d ScriptDB) AllDone() bool {
	for _, s := range d {
		if !s.Done() {
			return false
		}
	}

	return true
}

type Time int

type Performer interface {
	Perform(Time, *Mailbox, *Script) (*SentMessage, bool)
}

type PerformerDB map[PerformerId]Performer

// maybe this "Performer" interface should just be performer checks
type MeteredPerformer struct {
	Id    PerformerId
	Meter Time
}

func (p MeteredPerformer) Perform(t Time, mb *Mailbox, s *Script) (*SentMessage, bool) {
	if t%p.Meter != 0 {
		return nil, false
	}

	if s.Done() {
		return &SentMessage{p.Id, Done}, true
	}

	if mb.Empty() {
		return &SentMessage{p.Id, s.Current()}, true
	}

	if mb.Left == nil {
		if mb.Right.Word == Done {
			return &SentMessage{p.Id, s.Current()}, true
		}

		return &SentMessage{p.Id, mb.Right.Word}, true
	}

	if mb.Right == nil {
		if mb.Left.Word == Done {
			return &SentMessage{p.Id, s.Current()}, true
		}

		return &SentMessage{p.Id, mb.Left.Word}, true
	}

	if mb.Left.Word == mb.Right.Word {
		fmt.Printf("%s advances\n", p.Id)
		s.Advance()
		return &SentMessage{p.Id, s.Current()}, true
	}

	if mb.Left.Word == Done {
		return &SentMessage{p.Id, mb.Right.Word}, true
	}

	if mb.Right.Word == Done {
		return &SentMessage{p.Id, mb.Left.Word}, true
	}

	if mb.Left.Time.Sub(mb.Right.Time).Nanoseconds() > 0 {
		return &SentMessage{p.Id, mb.Left.Word}, true
	}

	return &SentMessage{p.Id, mb.Right.Word}, true
}

// type PerformanceTimer struct {
// 	ticks <-chan Time
// 	clients []chan Time
// 	stop chan<-bool
// }

// func NewPerformanceTimer(rate time.Duration) *PerformanceTimer {
// 	tickChan := make(chan Time)
// 	stopChan := make(chan bool)
// 	var clients []chan Time

// 	go func() {
// 		defer close(tickChan)
// 		defer close(stopChan)

// 		timer := time.NewTicker(rate)
// 		timeChan := timer.C
// 		var tick Time
// 		for {
// 			select {
// 			case <-stopChan:
// 				timer.Stop()
// 				for _, client := range clients {
// 					close(client)
// 				}
// 				close(stopChan)
// 				close(tickChan)
// 				return
// 			case <-timeChan:
// 				tick++
// 				for _, client := range clients {
// 					go func(client chan<- Time) { client <- tick }(client)
// 				}
// 			}
// 		}
// 	}()

// 	return &PerformanceTimer{
// 		ticks: tickChan,
// 		stop: stopChan,
// 		clients: clients,
// 	}
// }

// func (t *PerformanceTimer) Stop() {
// 	t.stop <- true
// }

// func (t *PerformanceTimer) GetTicks() <-chan Time {
// 	out := make(chan Time)
// 	t.clients = append(t.clients, out)
// 	return out
// }

// type PerformerTimer struct {
// 	performer PerformerId
// 	speakRate int
// 	nextSpeak int
// 	stopChan  chan<-bool
// }

// func NewPerformerTimer(id PerformerId, rate int, performanceTime <-chan Time) *PerformerTimer {
// 	return &PerformerTimer{
// 		performer: id,
// 		speakRate: rate,
// 		nextSpeak: rate,
// 		stopChan: make(chan bool),
// 	}
// }

// // func (t *PerformerTimer) Start() <-chan Time {
// // 	ticks := make(chan
// // }

// func (t *PerformerTimer) Stop() {
// 	t.stopChan <- true
// }

// func PerformLoop(id PerformerId, router PerformerRouter, scripts ScriptDB, ticks <-chan Time) (<-chan *SentMessage, chan<- bool) {
// 	words := make(chan *SentMessage)
// 	stop := make(chan bool)

// 	go func() {
// 		defer close(words)
// 		var speakTime time.Time
// 		var done bool
// 		for !done {
// 			select {
// 			case done := <-stop:
// 			case <-ticks:
// 				word := Perform(router[id], scripts[id])
// 				words <- &SentMessage{id, word}
// 			}
// 		}
// 	}()

// 	return words
// }
