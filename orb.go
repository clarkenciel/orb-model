package orb

type PerformerId int

type Word string

type Slot int

const (
	Left Slot = iota
	Right
	Self
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

type PerformerMailbox struct {
	Left  *Word
	Right *Word
	Self  *Word
}

func (m *PerformerMailbox) Receive(msg RoutedMessage) {
	switch msg.Address.Slot {
	case Left:
		m.Left = &msg.Message
		break
	case Right:
		m.Right = &msg.Message
		break
	case Self:
		m.Self = &msg.Message
		break
	}
}

type Script []Word

type ScriptDB map[PerformerId]Script
