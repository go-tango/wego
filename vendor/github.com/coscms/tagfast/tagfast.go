package tagfast

import (
	"reflect"
	"strconv"
	"sync"
)

//[struct_name][field_name]
var CachedStructTags map[string]map[string]*TagFast = make(map[string]map[string]*TagFast)
var lock *sync.RWMutex = new(sync.RWMutex)

func CacheTag(struct_name string, field_name string, value *TagFast) {
	lock.Lock()
	defer lock.Unlock()
	if _, ok := CachedStructTags[struct_name]; !ok {
		CachedStructTags[struct_name] = make(map[string]*TagFast)
	}
	CachedStructTags[struct_name][field_name] = value
}

func GetTag(struct_name string, field_name string) (r *TagFast, ok bool) {
	lock.RLock()
	defer lock.RUnlock()
	var v map[string]*TagFast
	v, ok = CachedStructTags[struct_name]
	if !ok {
		return
	}
	r, ok = v[field_name]
	return
}

//usage: Tag1(t, i, "form")
func Tag1(t reflect.Type, field_no int, key string) (tag string) {
	f := t.Field(field_no)
	tag = Tag(t, f, key)
	return
}

//usage: Tag2(t, "Id", "form")
func Tag2(t reflect.Type, field_name string, key string) (tag string) {
	f, ok := t.FieldByName(field_name)
	if !ok {
		return ""
	}
	tag = Tag(t, f, key)
	return
}

func Tag(t reflect.Type, f reflect.StructField, key string) (tag string) {
	if f.Tag == "" {
		return ""
	}
	if v, ok := GetTag(t.String(), f.Name); ok {
		tag = v.Get(key)
	} else {
		v := TagFast{Tag: f.Tag}
		tag = v.Get(key)
		CacheTag(t.String(), f.Name, &v)
	}
	return
}

func Tago(t reflect.Type, f reflect.StructField, key string) (tag string, tf *TagFast) {
	if f.Tag == "" {
		return "", nil
	}
	if v, ok := GetTag(t.String(), f.Name); ok {
		tag = v.Get(key)
		tf = v
	} else {
		tf = &TagFast{Tag: f.Tag}
		tag = tf.Get(key)
		CacheTag(t.String(), f.Name, tf)
	}
	return
}

func ClearTag() {
	CachedStructTags = make(map[string]map[string]*TagFast)
}

func ParseStructTag(tag string) map[string]string {
	lock.Lock()
	defer lock.Unlock()
	var tagsArray map[string]string = make(map[string]string)
	for tag != "" {
		// skip leading space
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		// scan to colon.
		// a space or a quote is a syntax error
		i = 0
		for i < len(tag) && tag[i] != ' ' && tag[i] != ':' && tag[i] != '"' {
			i++
		}
		if i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		// scan quoted string to find value
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qvalue := string(tag[:i+1])
		tag = tag[i+1:]

		value, _ := strconv.Unquote(qvalue)
		tagsArray[name] = value
	}
	return tagsArray
}

type TagFast struct {
	Tag    reflect.StructTag
	Cached map[string]string
	Parsed map[string]interface{}
}

func (a *TagFast) Get(key string) string {
	if a.Cached == nil {
		a.Cached = ParseStructTag(string(a.Tag))
	}
	lock.RLock()
	defer lock.RUnlock()
	if v, ok := a.Cached[key]; ok {
		return v
	}
	return ""
}

func (a *TagFast) GetParsed(key string, fns ...func() interface{}) interface{} {
	if a.Parsed == nil {
		a.Parsed = make(map[string]interface{})
	}
	lock.RLock()
	if v, ok := a.Parsed[key]; ok {
		lock.RUnlock()
		return v
	}
	lock.RUnlock()
	if len(fns) > 0 {
		fn := fns[0]
		if fn != nil {
			v := fn()
			a.SetParsed(key, v)
			return v
		}
	}
	return nil
}

func (a *TagFast) SetParsed(key string, value interface{}) bool {
	if a.Parsed == nil {
		a.Parsed = make(map[string]interface{})
	}
	lock.Lock()
	defer lock.Unlock()
	a.Parsed[key] = value
	return true
}
