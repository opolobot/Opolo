package lib

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

// MsgCollector collects messages.
type MsgCollector struct {
	ID        uint16
	ChannelID string
	Msgs      chan *discordgo.Message
}

// Collect collects messages in the designated channel for the supplied amount of time.
func (w *Whiskey) Collect(channelID string, sleepTime time.Duration, amnt int) *MsgCollector {
	msgs := (func() chan *discordgo.Message {
		if amnt == 0 {
			return make(chan *discordgo.Message)
		}

		return make(chan *discordgo.Message, amnt)
	})()

	cols, ok := w.Collectors[channelID]
	var id uint16
	if !ok {
		cols = []*MsgCollector{}
	} else {
		id = w.findCollectorID(channelID)
	}

	collector := &MsgCollector{
		ID:        id,
		ChannelID: channelID,
		Msgs:      msgs,
	}

	cols = append(cols, collector)

	w.Collectors[channelID] = cols

	if sleepTime > 0 {
		go (func() {
			time.Sleep(sleepTime)
			w.CancelCollector(channelID, id)
		})()
	}

	return collector
}

// CancelCollector cancels a message collector.
func (w *Whiskey) CancelCollector(channelID string, id uint16) error {
	cols, ok := w.Collectors[channelID]
	if ok {
		found := false
		for i, col := range cols {
			if col.ID == id {
				// Close the message collection channel.
				close(col.Msgs)

				// Delete the element from the slice.
				// https://yourbasic.org/golang/delete-element-slice/
				// Copy last element to index i.
				cols[i] = cols[len(cols)-1]
				// Erase last element (write zero value).
				cols[len(cols)-1] = nil
				cols = cols[:len(cols)-1]

				// Recommit the slice to the whiskey instance.
				w.Collectors[channelID] = cols // <-- forgot this
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("No collector for channel %v with ID %v was found", channelID, id)
		}
	} else {
		return fmt.Errorf("No collectors running for channel with ID %v", channelID)
	}

	return nil
}

func (w *Whiskey) findCollectorID(channelID string) uint16 {
	var currentID uint16
	for _, col := range w.Collectors[channelID] {
		if col.ID == currentID {
			currentID++
		} else {
			break
		}
	}

	return currentID
}
