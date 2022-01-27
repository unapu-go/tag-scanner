package tag_scanner

import "strings"

var Default = Scanner{';', ':', '{', '}'}

type Scanner struct {
	Fields,
	Key,
	Start,
	End uint8
}

func (this Scanner) ScanValue(s string, keyArgs *[]string) (value, news string) {
	s = strings.TrimSpace(s)
	if s == "" {
		return
	}
start:
	if s[0] == this.Start {
		s = s[1:]
		var count = 1
		for i, r := range s {
			switch uint8(r) {
			case this.Start:
				count++
			case this.End:
				count--
				if count == 0 {
					return string(this.Start) + s[:i] + string(this.End), s[i+1:]
				}
			}
		}
	} else {
		for i, r := range s {
			switch uint8(r) {
			case this.Key:
				*keyArgs = append(*keyArgs, s[:i])
				s = s[i+1:]
				goto start
			case this.Fields:
				return s[:i], s[i+1:]
			}
		}
	}
	return s, ""
}

func (this Scanner) Scan(s string) (key, value, news string, keyArgs []string) {
	s = strings.TrimPrefix(strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(s), string(this.Fields))), string(this.Fields))
	if s == "" {
		return
	}
	if s[0] == this.Start {
		value, news = this.ScanValue(s, &keyArgs)
		return
	}
	for i, r := range s {
		switch uint8(r) {
		case this.Fields:
			key = s[:i]
			news = s[i+1:]
			return
		case this.Key:
			key = s[:i]
			value, news = this.ScanValue(s[i+1:], &keyArgs)
			return
		}
	}
	key = s
	news = ""
	return
}

func (this Scanner) ScanAll(s string, cb func(node Node)) {
	if this.Start == 0 {
		panic("undef scanner")
	}
	if s = strings.TrimSpace(s); s == "" {
		return
	}
	if this.IsTags(s) {
		if s = strings.TrimSpace(s[1 : len(s)-1]); s == "" {
			return
		}
	}
	if s = strings.TrimPrefix(s, string(this.Fields)); s == "" {
		return
	}

	var (
		key, value string
		keyArgs    []string
	)
	for s != "" {
		key, value, s, keyArgs = this.Scan(s)
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if key != "" {
			if value == "" {
				cb(NodeFlag(key))
			} else {
				cb(NodeKeyValue{key, value, keyArgs})
			}
		} else if value != "" {
			cb(NodeTags(value))
		}
		s = strings.Trim(s, string(this.Fields))
	}
}

func (this Scanner) ScanAllFlags(s string, flags ParseFlag, cb func(flag bool, key string, value string)) {
	this.ScanAll(s, func(node Node) {
		switch node.Type() {
		case Tags:
			if flags.Has(FlagTags) {
				n := node.(NodeTags)
				cb(false, n.Value(), n.Value())
			}
		case KeyValue:
			kv := node.(NodeKeyValue)
			key := kv.Key
			if !flags.Has(FlagPreserveKeys) {
				key = strings.ToUpper(key)
			}
			cb(false, key, kv.Value)
		case Flag:
			name := node.String()
			if !flags.Has(FlagPreserveKeys) {
				name = strings.ToUpper(name)
			}
			cb(true, name, "")
		}
	})
}

func (this Scanner) IsTags(value string) bool {
	value = strings.TrimSpace(value)
	if value == "" {
		return false
	}
	return value[0] == this.Start && value[len(value)-1] == this.End
}

func (this Scanner) String(value string) string {
	if this.IsTags(value) {
		return strings.ReplaceAll(value[1:len(value)-1], `\"`, `"`)
	}
	return value
}

func (this Scanner) ToTags(value string) string {
	return string(this.Start) + value + string(this.End)
}
