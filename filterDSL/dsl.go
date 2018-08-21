package filterDSL

import (
	"github.com/vvampirius/feed-processor/feed"
	"errors"
	"log"
	"regexp"
)

type DSL struct {
	Filters []Filter
	Pipelines []Pipeline
}

func (dsl *DSL) GetPipelineByUrl(url string) *Pipeline {
	for _, pipeline := range dsl.Pipelines {
		if pipeline.Url == url { return &pipeline }
	}
	return nil
}

func (dsl *DSL) GetFilterByName(name string) *Filter {
	for _, filter := range dsl.Filters {
		if filter.Name == name { return &filter }
	}
	return nil
}

func (dsl *DSL) GetFiltersByNames(names []string) []*feed.RegexpFilter {
	filters := make([]*feed.RegexpFilter, 0)
	for _, name := range names {
		filter := dsl.GetFilterByName(name)
		if filter == nil {
			log.Printf("filter '%s' not found!\n", name)
			continue
		}

		rf := feed.RegexpFilter{Name: filter.Name, Type: filter.Type, Regexps: make([]*regexp.Regexp, 0)}
		if filter.Type != `include` && filter.Type != `exclude` {
			log.Printf("Type '%s' in filter '%s' unexpected!\n", filter.Type, filter.Name)
			continue
		}

		//TODO: if not compiled -> don't add to filters
		for _, rs := range filter.Regexps {
			r, err := regexp.Compile(rs)
			if err != nil {
				log.Printf("Can't compile regexp '%s': %s\n", rs, err.Error())
				continue
			}
			rf.Regexps = append(rf.Regexps, r)
		}

		filters = append(filters, &rf)
	}
	return filters
}

func (dsl *DSL) RegexpFilterPipeline(url string) (*feed.RegexpFilterPipeline, error) {
	pipeline := dsl.GetPipelineByUrl(url)
	if pipeline == nil { return nil, errors.New(`Pipeline not found`) }

	rfp := feed.RegexpFilterPipeline{
		Name: pipeline.Name,
		RegexpFilters: dsl.GetFiltersByNames(pipeline.Filters),
	}

	return &rfp, nil
}