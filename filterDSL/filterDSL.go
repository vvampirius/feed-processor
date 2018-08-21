package filterDSL

import (
	"time"
	"gopkg.in/yaml.v2"
	"os"
	"github.com/vvampirius/feed-processor/feed"
	"errors"
	"log"
)

type FilterDSL struct {
	FileName string
	FileTimestamp time.Time
	DSL *DSL
}

func (filterDSL *FilterDSL) Reload() error {
	f, err := os.Open(filterDSL.FileName)
	if err != nil { return err }
	defer f.Close()

	dsl := DSL{}

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&dsl)
	if err != nil { return err }

	filterDSL.DSL = &dsl

	return nil
}

func (filterDSL *FilterDSL) CheckUpdate() {
	for {
		if i, err := os.Stat(filterDSL.FileName); err==nil {
			if i.ModTime().After(filterDSL.FileTimestamp) {
				log.Println(`Reloading due to DSL update`)
				if err := filterDSL.Reload(); err==nil {
					filterDSL.FileTimestamp = time.Now()
				} else { log.Println(err) }
			}
		}
		time.Sleep(30 * time.Second)
	}
}

func (filterDSL *FilterDSL) RegexpFilterPipeline(url string) (*feed.RegexpFilterPipeline, error) {
	if filterDSL.DSL == nil { return nil, errors.New(`FilterDSL is not loaded!`) }
	return filterDSL.DSL.RegexpFilterPipeline(url)
}