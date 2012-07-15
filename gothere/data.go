package gothere

import (
	"appengine"
	"appengine/datastore"
	"errors"
	"net/url"
	"regexp"
)

type NamedRedir struct {
	Target     string
	Identifier string
}

func (s NamedRedir) RenderHtml() string {
	return "it is murder time!!"
}

type RedirCollection struct {
	Redirs []NamedRedir
}

func (s RedirCollection) RenderHtml() string {
	return "collection of redirs"
}

func (s NamedRedir) Validate() error {
	NRIdent, _ := regexp.Compile("[\\w-]+")
	u, err := url.Parse(s.Target)
	if err != nil {
		return err
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("Must have http or https scheme")
	}

	if !NRIdent.MatchString(s.Identifier) {
		return errors.New("Identifier must be letters (a-z), numbers (0-9)" +
			" underscores (_) and dashes (-) only")
	}
	return nil
}

func (s NamedRedir) SaveNew(c appengine.Context) error {
	//k := datastore.NewIncompleteKey(c, "NamedRedir", nil)
	k := datastore.NewKey(c,"NamedRedir",s.Identifier,0,nil)
	_, err := datastore.Put(c, k, &s)
	return err
}

func (s NamedRedir) Location() string {
	return "/urls/" + s.Identifier
}

func GetRedirs(c appengine.Context, start_at int, count int) ([]NamedRedir, error) {
	if count > 100 {
		return nil, errors.New("count has to be under 100")
	}
	q := datastore.NewQuery("NamedRedir").Offset(start_at).Limit(count)

	redirs := make([]NamedRedir, 0, count)
	if _, err := q.GetAll(c, &redirs); err != nil {
		return nil, err
	}
	return redirs, nil
}

func GetRedir(c appengine.Context,ident string) (NamedRedir, error) {
	k := datastore.NewKey(c, "NamedRedir", ident, 0, nil)
	var redir NamedRedir
	if err := datastore.Get(c,k,&redir); err != nil {
		return redir, err
	}
	return redir, nil
}