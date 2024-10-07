package event_consumer

import (
	"log"
	"strings"
	"sync"
	"time"

	"app/events"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c *Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())

			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		if err := c.handleEvents(gotEvents); err != nil {
			log.Print(err)

			continue
		}
	}
}

func (c *Consumer) handleEvents(eventList []events.Event) error {
	var wg sync.WaitGroup

	for _, event := range eventList {
		log.Printf("got new event: %s", strings.Fields(event.Text)[0])

		// Увеличиваем счетчик WaitGroup
		wg.Add(1)

		go func(evt events.Event) {
			defer wg.Done() // Уменьшаем счетчик по завершении горутины

			if err := c.processor.Process(evt); err != nil {
				log.Printf("can't handle event: %s", err.Error())
			}
		}(event)
	}

	// Ждем завершения всех горутин
	wg.Wait()

	return nil
}
