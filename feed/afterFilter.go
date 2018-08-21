package feed

import (
	"time"
	"github.com/mmcdole/gofeed"
	"log"
)

func AfterFilter(items []*gofeed.Item, after time.Time, logger *log.Logger) []*gofeed.Item {
	filteredItems := make([]*gofeed.Item, 0)
	if logger != nil { logger.Printf("Filter items before '%v'\n", after)}
	for _, item := range items {
		action := `skip`
		if item.PublishedParsed.After(after) {
			filteredItems = append(filteredItems, item)
			action = `add`
		}
		if logger != nil { logger.Printf("AfterFilter: (%v :: %s) => %s\n", item.PublishedParsed, item.Title, action)}
	}
	return filteredItems
}