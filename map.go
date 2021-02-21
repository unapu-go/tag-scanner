package tag_scanner

import "strings"

const (
	PreserveKeys ParseFlag = 1 << iota
	ForceTags
	FlagTags
	NotNil
)

type ParseFlag uint8

func (b ParseFlag) Set(flag ParseFlag) ParseFlag    { return b | flag }
func (b ParseFlag) Clear(flag ParseFlag) ParseFlag  { return b &^ flag }
func (b ParseFlag) Toggle(flag ParseFlag) ParseFlag { return b ^ flag }
func (b ParseFlag) Has(flag ParseFlag) bool         { return b&flag != 0 }

func JoinParseFlags(flag ...ParseFlag) (f ParseFlag) {
	for _, flag := range flag {
		f |= flag
	}
	return
}

type Map map[string]string

func (this Map) Clone() (clone Map) {
	if this == nil {
		return
	}
	clone = make(Map)
	for key, value := range this {
		clone[key] = value
	}
	return
}

func (this Map) SetParseFlag(name ...string) Map {
	this = this.Clone()
	if this == nil {
		this = Map{}
	}
	for _, name := range name {
		this[name] = name
	}
	return this
}

func (this Map) Get(key string) (v string) {
	if this != nil {
		v = this[key]
	}
	return
}

func (this Map) GetOk(key string) (v string, ok bool) {
	if this == nil {
		return
	}
	v, ok = this[key]
	return
}

func (this Map) GetString(key string) (v string) {
	if this == nil {
		return
	}
	if v, ok := this[key]; ok {
		return Default.String(v)
	}
	return
}

func (this Map) GetStringAlias(key string, alias ...string) (v string) {
	if this == nil {
		return
	}
	if v, ok := this[key]; ok {
		return Default.String(v)
	}
	for _, key := range alias {
		if v, ok := this[key]; ok {
			return Default.String(v)
		}
	}
	return
}

func (this Map) Empty() bool {
	return this == nil || len(this) == 0
}

func (this Map) Enable(key string) {
	this[key] = key
}

func (this *Map) Set(key, value string) {
	if *this == nil {
		*this = make(Map)
	}
	(*this)[key] = value
}

func (this Map) Flag(name string) bool {
	if this == nil {
		return false
	}
	return this[name] == name
}

func (this Map) Flags() (flags []string) {
	if this != nil {
		for k, v := range this {
			if k == v {
				flags = append(flags, k)
			}
		}
	}
	return
}

func (this Map) String() string {
	var pairs []string
	for k, v := range this {
		if k == v {
			pairs = append(pairs, k)
		} else {
			pairs = append(pairs, k+":"+v)
		}
	}
	return strings.Join(pairs, "; ")
}

func (this Map) Update(setting ...map[string]string) {
	for _, setting := range setting {
		for k, v := range setting {
			this[k] = v
		}
	}
}

func (this Map) GetTags(name string, flags ...ParseFlag) (tags Map) {
	flag := JoinParseFlags(flags...)
	if s := this[name]; s != "" {
		if !this.Scanner().IsTags(s) {
			if !flag.Has(ForceTags) {
				return
			}
			s = this.Scanner().ToTags(s)
		}

		tags = make(Map)
		tags.ParseString(s, flag)
	} else if flag.Has(NotNil) {
		return Map{}
	}
	return
}

func (this Map) TagsOf(value string, flags ...ParseFlag) (tags Map) {
	if value != "" && this.Scanner().IsTags(value) {
		tags = make(Map)
		tags.ParseString(value, JoinParseFlags(flags...))
	}
	return
}

func (this Map) SetFlag(flagName ...string) Map {
	for _, name := range flagName {
		this[name] = name
	}
	return this
}

func (this *Map) ParseDefault(tags StructTag, key string, keyAlias ...string) (ok bool) {
	return this.Parse(tags, key, 0, keyAlias...)
}

func (this *Map) Parse(tags StructTag, key string, flag ParseFlag, keyAlias ...string) (ok bool) {
	return this.ParseCallback(tags, append([]string{key}, keyAlias...), flag)
}

func (this *Map) ParseCallbackDefault(tags StructTag, keys []string, cb ...func(dest map[string]string, n Node)) (ok bool) {
	return this.ParseCallback(tags, keys, 0, cb...)
}

func (this *Map) ParseCallback(tags StructTag, keys []string, flag ParseFlag, cb ...func(dest map[string]string, n Node)) (ok bool) {
	if *this == nil {
		*this = make(Map)
	}
	var tags_ = make([]string, len(keys))
	for i, key := range keys {
		tags_[i] = tags.Get(key)
	}
	for _, str := range tags_ {
		(*this).ParseString(str, flag, cb...)
	}
	return len(*this) > 0
}

func (this *Map) Scanner() Scanner {
	return Default
}

func (this *Map) ParseStringDefault(s string, cb ...func(dest map[string]string, n Node)) {
	this.ParseString(s, 0, cb...)
}

func (this *Map) ParseString(s string, flags ParseFlag, cb ...func(dest map[string]string, n Node)) {
	if *this == nil {
		*this = map[string]string{}
	}
	scanner := this.Scanner()
	scanner.ScanAll(s, func(node Node) {
		for _, cb := range cb {
			cb(*this, node)
		}
		switch node.Type() {
		case Tags:
			if flags.Has(FlagTags) {
				n := node.(NodeTags)
				(*this)[n.Value()] = n.Value()
			}
		case KeyValue:
			kv := node.(NodeKeyValue)
			key := kv.Key
			if !flags.Has(PreserveKeys) {
				key = strings.ToUpper(key)
			}
			(*this)[key] = kv.Value
		case Flag:
			name := node.String()
			if !flags.Has(PreserveKeys) {
				name = strings.ToUpper(name)
			}
			(*this)[name] = name
		}
	})
}
