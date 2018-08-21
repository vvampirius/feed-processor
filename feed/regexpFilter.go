package feed

import (
	"regexp"
	"errors"
	"fmt"
	"os"
)

type RegexpFilter struct {
	Name string
	Type string
	Regexps []*regexp.Regexp
}

func (regexpFilter *RegexpFilter) Filter(s string) bool {
	for _, r := range regexpFilter.Regexps {
		if r.MatchString(s) {
			if regexpFilter.Type == `include` { return false }  // не фильтровать
			if regexpFilter.Type == `exclude` { return true }   // фильтровать
		}
	}
	if regexpFilter.Type == `include` { return true } // фильтровать
	return false  // не фильровать
}

func NewRegexpFilter(name string, type_ string, regexps []string) (*RegexpFilter, error) {
	if type_ != `include` && type_ != `exclude` {
		return nil, errors.New(fmt.Sprintf("%s != 'include' or 'exclude'!", type_))
	}
	regexpFilter := RegexpFilter{ Name: name, Type: type_, Regexps: make([]*regexp.Regexp, 0) }
	for _, stringRegexp := range regexps {
		r, err := regexp.Compile(stringRegexp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can't compile regexp '%s': %v", stringRegexp, err)
			continue
		}
		regexpFilter.Regexps = append(regexpFilter.Regexps, r)
	}
	return &regexpFilter, nil
}


type RegexpFilterPipeline struct {
	Name string
	RegexpFilters []*RegexpFilter
}

func (regexpFilterPipeline *RegexpFilterPipeline) Filter(s string) bool {
	for _, filter := range regexpFilterPipeline.RegexpFilters {
		filtered := filter.Filter(s)
		//fmt.Printf("RegexpFilterPipeline: '%s' '%s' %v\n", s, filter.Name, filtered)
		if filtered { return true }
	}
	return false
}