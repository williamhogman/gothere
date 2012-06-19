package gothere

import (
	"errors"
	"appengine"
	"appengine/datastore"
)

type NamedRedir struct {
	Target string
	Identifier string
}

func get_redirs(c appengine.Context,start_at int,count int) ([]NamedRedir,error) {
	if count > 100 {
		return nil, errors.New("count has to be under 100")
	}
	q := datastore.NewQuery("NamedRedir").Offset(start_at).Limit(count)
	
	redirs := make([]NamedRedir,0, count)
	if _, err := q.GetAll(c,&redirs); err != nil {
		return nil, err
	}
	return redirs,nil
}