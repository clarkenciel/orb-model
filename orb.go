package orb

import "time"

type PerformerId string

type Word string

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

type PerformerMailbox struct {
	Left  *ReceivedWord
	Right *ReceivedWord
}

func (m *PerformerMailbox) Receive(msg RoutedMessage) {
	switch msg.Address.Slot {
	case Left:
		m.Left = &ReceivedWord{msg.Message, time.Now()}
		break
	case Right:
		m.Right = &ReceivedWord{msg.Message, time.Now()}
		break
	}
}

func (m *PerformerMailbox) Empty() bool {
	return m.Left == nil && m.Right == nil
}

type Script []Word

func (s Script) Current() Word {
	return s[0]
}

func (s Script) Advance() {
	s = s[1:]
}

type ScriptDB map[PerformerId]Script

func Perform(mb *PerformerMailbox, s *Script) Word {
	if mb.Empty() {
		return s.Current()
	}
	
	if mb.Left.Word == mb.Right.Word {
		s.Advance()
		return s.Current()
	}

	if mb.Left.Time.Sub(mb.Right.Time).Nanoseconds() > 0 {
		return mb.Left.Word
	}

	return mb.Right.Word
}

type SpokenWord struct {
	Word Word
	Performer PerformerId
}

func PerformLoop(id PerformerId, router *PerformerRouter, scripts *ScriptDB) (<-chan *SpokenWord, chan<- bool) {
	words := make(chan *SpokenWord)
	stop := make(chan bool)
	
	go func() {
		defer close(words)
		var speakTime time.Time
		var done bool
		for !done {
			select {
			case done <-stop:
			default:
				now = time.Now()
				if now > speakTime {
					word := Perform(router[id], scripts[id])
					speakTime = reschedule(now)
					words <- word
				} else {
					
				}				
			}
		}
	}()

	return words
}

func reschedule(now time.Time) time.Time {
	return now
}
