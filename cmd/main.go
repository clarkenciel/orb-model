package main

import (
	"fmt"
	"strings"

	orb "github.com/clarkenciel/orb-model"
)

func main() {
	ids := []orb.PerformerId{
		"tom",
		"dick",
		"harry",
	}

	// 	script := orb.Script{
	// 		"one",
	// 		"two",
	// p		"three",
	// 	}

	router := orb.PerformerRouter{
		"tom": orb.AddressSet{
			orb.Address{"dick", orb.Left}:   true,
			orb.Address{"harry", orb.Right}: true,
		},
		"dick": orb.AddressSet{
			orb.Address{"harry", orb.Left}: true,
			orb.Address{"tom", orb.Right}:  true,
		},
		"harry": orb.AddressSet{
			orb.Address{"tom", orb.Left}:   true,
			orb.Address{"dick", orb.Right}: true,
		},
	}

	mailRoom := orb.MailRoom{
		"tom":   &orb.Mailbox{},
		"dick":  &orb.Mailbox{},
		"harry": &orb.Mailbox{},
	}

	scripts := orb.ScriptDB{
		"tom":   &orb.Script{"one", "two", "three"},
		"dick":  &orb.Script{"one", "two", "three"},
		"harry": &orb.Script{"one", "two", "three"},
	}

	performers := orb.PerformerDB{
		"tom":   orb.MeteredPerformer{"tom", 1},
		"dick":  orb.MeteredPerformer{"dick", 3},
		"harry": orb.MeteredPerformer{"harry", 7},
	}

	for i := 1; !scripts.AllDone(); i++ { // should not start a 0 for this
		// speaking phase
		var messages []*orb.SentMessage
		for _, id := range ids {
			mb, found := mailRoom[id]
			if !found {
				continue
			}

			script, found := scripts[id]
			if !found {
				continue
			}

			performer, found := performers[id]
			if !found {
				continue
			}

			message, sent := performer.Perform(orb.Time(i), mb, script)
			if sent {
				messages = append(messages, message)
				mb.Clear()
			}
		}

		printMessages(i, messages)

		// routing phase
		var routed []*orb.RoutedMessage
		for _, sent := range messages {
			for _, rm := range router.Route(*sent) {
				routed = append(routed, rm)
			}
		}

		// listening phase
		for _, routedMessage := range routed {
			target := routedMessage.Address.Performer

			mb, found := mailRoom[target]
			if !found {
				continue
			}

			mb.Receive(*routedMessage)
		}
	}
}

func printMessages(g int, ms []*orb.SentMessage) {
	var header strings.Builder
	var row strings.Builder

	_, e := header.WriteString(fmt.Sprintf("%10d. |", g))
	if e != nil {
		return
	}

	_, e = row.WriteString(fmt.Sprintf("%10s |", " "))
	if e != nil {
		return
	}

	for _, m := range ms {
		_, e := header.WriteString(fmt.Sprintf("%10s|", m.Sender))
		if e != nil {
			return
		}

		_, e = row.WriteString(fmt.Sprintf("%10s|", m.Message))
		if e != nil {
			return
		}
	}

	fmt.Println(header.String(), "\n", row.String())
}
