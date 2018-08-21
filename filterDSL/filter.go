package filterDSL

//TODO: add expire date
type Filter struct {
	Name string
	Type string
	Regexps []string
}

