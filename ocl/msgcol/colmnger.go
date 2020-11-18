package msgcol

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

var instance *CollectionManager

// CollectionManager is a mediator for message collectors.
type CollectionManager struct {
	Collectors map[string][]*MessageCollector
}


// GetCollectionManager returns CollectionManager instance.
func GetCollectionManager() *CollectionManager {
	if instance == nil {
		instance = &CollectionManager{
			Collectors: make(map[string][]*MessageCollector),
		}
	}

	return instance
}

// Dispatch dispatches a message to the appropriate collectors.
func (colMnger *CollectionManager) Dispatch(msg *discordgo.Message) {
	for chanID, cols := range colMnger.Collectors {
		if chanID == msg.ChannelID {
			for _, col := range cols {
				col.Msgs <- msg
				if len(col.Msgs) == cap(col.Msgs) {
					colMnger.CancelCollector(chanID, col.ID)
				}
			}
		}
	}
}

// NewCollector collects messages in the designated channel for the supplied amount of time.
func (colMnger *CollectionManager) NewCollector(channelID string, sleepTime time.Duration, amnt int) *MessageCollector {
	msgs := (func() chan *discordgo.Message {
		if amnt == 0 {
			return make(chan *discordgo.Message)
		}

		return make(chan *discordgo.Message, amnt)
	})()

	cols, ok := colMnger.Collectors[channelID]
	var id uint16
	if !ok {
		cols = []*MessageCollector{}
	} else {
		id = colMnger.findCollectorID(channelID)
	}

	collector := &MessageCollector{
		ID:        id,
		ChannelID: channelID,
		Msgs:      msgs,
	}

	cols = append(cols, collector)

	colMnger.Collectors[channelID] = cols

	if sleepTime > 0 {
		go (func() {
			time.Sleep(sleepTime)
			colMnger.CancelCollector(channelID, id)
		})()
	}

	return collector
}

// CancelCollector cancels a message collector.
// TODO(@zorbyte): Make this part of the collector itself.
func (colMnger *CollectionManager) CancelCollector(channelID string, id uint16) error {
	cols, ok := colMnger.Collectors[channelID]
	if ok {
		found := false
		for i, col := range cols {
			if col.ID == id {
				// Close the message collection channel.
				close(col.Msgs)

				// Delete the element from the slice.
				lastIdx := len(cols) - 1
				cols[i] = cols[lastIdx]
				cols[lastIdx] = nil
				cols = cols[:lastIdx]

				// Recommit the slice to the whiskey instance.
				colMnger.Collectors[channelID] = cols
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


func (colMnger *CollectionManager) findCollectorID(channelID string) uint16 {
	var currentID uint16
	for _, col := range colMnger.Collectors[channelID] {
		if col.ID == currentID {
			currentID++
		} else {
			break
		}
	}

	return currentID
}