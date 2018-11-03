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
		"tom":   &orb.Script{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten"},
		"dick":  &orb.Script{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten"},
		"harry": &orb.Script{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten"},
	}

	performers := orb.PerformerDB{
		"tom":   orb.Performer{"tom", orb.Meter(2)},
		"dick":  orb.Performer{"dick", orb.Meter(3)},
		"harry": orb.Performer{"harry", orb.Meter(7)},
	}

	printHeader(ids)
	for i := 1; !scripts.AllDone(); i++ {
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

		printMessages(i, ids, indexMessages(messages))

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

func indexMessages(ms []*orb.SentMessage) map[orb.PerformerId]orb.Word {
	index := make(map[orb.PerformerId]orb.Word)
	for _, m := range ms {
		index[m.Sender] = m.Message
	}
	return index
}

const RowItemSize = 10

func padded(padding int, formatSigil string) string {
	return fmt.Sprintf("%%%d%s", padding, formatSigil)
}

func printHeader(ids []orb.PerformerId) {
	fmtStr := padded(RowItemSize, "s |")

	var header strings.Builder

	_, e := header.WriteString(fmt.Sprintf(fmtStr, " "))
	if e != nil {
		return
	}

	for _, id := range ids {
		_, e := header.WriteString(fmt.Sprintf(fmtStr, id))
		if e != nil {
			return
		}
	}

	fmt.Println(header.String())
}

func printMessages(g int, ids []orb.PerformerId, ms map[orb.PerformerId]orb.Word) {
	var row strings.Builder

	_, e := row.WriteString(fmt.Sprintf(padded(RowItemSize, "d |"), g))
	if e != nil {
		return
	}

	fmtStr := padded(RowItemSize, "s |")
	for _, id := range ids {
		m, found := ms[id]
		if found {
			_, e = row.WriteString(fmt.Sprintf(fmtStr, m))
			if e != nil {
				return
			}

			continue
		}

		_, e = row.WriteString(fmt.Sprintf(fmtStr, " "))
		if e != nil {
			return
		}
	}

	fmt.Println(row.String())
}
