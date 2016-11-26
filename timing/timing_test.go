package timing

import (
	"testing"

	"github.com/arstd/log"
	"github.com/pborman/uuid"
	"time"
)

func TestTiming(t *testing.T) {
	RemindFunc = func(items ...*Item) {
		log.JSON("remind", items)
	}

	when := uint32(time.Now().Add(time.Second).Unix())

	// load data
	loaded := []*Item{
		{ID: uuid.New(), Timestamp: when + 2, Event: "test", Description: "2"},
		{ID: uuid.New(), Timestamp: when + 10, Event: "test", Description: "10"},
		{ID: uuid.New(), Timestamp: when + 2, Event: "test", Description: "2"},
		{ID: uuid.New(), Timestamp: when + 4, Event: "test", Description: "4"},
	}
	Init(loaded...)

	loaded = []*Item{
		{ID: uuid.New(), Timestamp: when + 3, Event: "test", Description: "3"},
		{ID: uuid.New(), Timestamp: when + 5, Event: "test", Description: "5"},
	}
	Init(loaded...)

	Add(&Item{Timestamp: when + 2, Event: "test", Description: "2"})
	Add(&Item{Timestamp: when + 5, Event: "test", Description: "5"})
	Add(&Item{Timestamp: when + 9, Event: "test", Description: "9"})

	time.Sleep(10 * time.Second)

	Add(&Item{Timestamp: when + 14, Event: "test", Description: "14"})
	Add(&Item{Timestamp: when + 12, Event: "test", Description: "12"})
	Add(&Item{Timestamp: when + 14, Event: "test", Description: "14"})

	time.Sleep(5 * time.Second)
}
