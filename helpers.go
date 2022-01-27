package tag_scanner

import "strings"

func Strings(v string) (values []string) {
	if v = strings.TrimSpace(v); v == "" {
		return
	}
	if Default.IsTags(v) {
		values = strings.Split(v, string(Default.Fields))
		for i := range values {
			values[i] = strings.TrimSpace(values[i])
		}
		return
	}
	return append(values, v)
}

func KeyValuePairs(s Scanner, v string) (values [][]string) {
	if v = strings.TrimSpace(v); v == "" {
		return
	}
	s.ScanAll(v, func(node Node) {
		switch node.Type() {
		case KeyValue:
			kv := node.(NodeKeyValue)
			values = append(values, []string{kv.Key, kv.Value})
		case Flag:
			f := string(node.(NodeFlag))
			values = append(values, []string{f, f})
		}
	})
	return
}

func Flags(s Scanner, v string) (names []string) {
	s.ScanAll(v, func(node Node) {
		if node.Type() == Flag {
			names = append(names, node.String())
		}
	})
	return
}

func NonFlags(s Scanner, v string) (names []string) {
	s.ScanAll(v, func(node Node) {
		if node.Type() != Flag {
			names = append(names, node.String())
		}
	})
	return
}

func Parse(s Scanner, v string, flag ...ParseFlag) (tags Map, flags Set) {
	s.ScanAllFlags(v, JoinParseFlags(flag...), func(flag bool, key string, value string) {
		if flag {
			flags.Add(key)
		} else {
			tags.Set(key, value)
		}
	})
	return
}
