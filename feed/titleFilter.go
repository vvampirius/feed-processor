package feed

import (
	"github.com/mmcdole/gofeed"
	"log"
)

func TitleFilter(items []*gofeed.Item, rfp *RegexpFilterPipeline, logger *log.Logger) []*gofeed.Item {
	newItems := make([]*gofeed.Item, 0)
	for _, item := range items {
		filtered := rfp.Filter(item.Title)
		if logger != nil {
			action := `add`
			if filtered { action = `skip` }
			logger.Printf("TitleFilter: '%s' => %s\n", item.Title, action)
		}
		if !filtered { newItems = append(newItems, item) }
	}
	return newItems
}