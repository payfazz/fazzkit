package exactlyonce

import "github.com/payfazz/fazzkit/event/stan/message"

//Opt exactly once options
type Opt struct {
	DbName     string
	Repository message.Repository
}

//NewOpt create exactly once options
func NewOpt(opt Opt) *Opt {
	newOpt := &Opt{
		DbName:     opt.DbName,
		Repository: opt.Repository,
	}
	return newOpt
}
