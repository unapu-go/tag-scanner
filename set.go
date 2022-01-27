package tag_scanner

type (
	void struct{}
	Set  map[string]void
)

var setValue void

func (s *Set) Add(key ...string) *Set {
	if *s == nil {
		*s = map[string]void{}
	}
	for _, key := range key {
		(*s)[key] = setValue
	}
	return s
}

func (s *Set) Del(key ...string) *Set {
	if *s == nil {
		return s
	}
	for _, key := range key {
		delete(*s, key)
	}
	return s
}

func (s Set) Has(key string) (ok bool) {
	if s == nil {
		return
	}
	_, ok = s[key]
	return
}

func (s Set) Strings() (keys []string) {
	if s == nil {
		return
	}
	keys = make([]string, len(s))
	var i int
	for key := range s {
		keys[i] = key
		i++
	}
	return
}
