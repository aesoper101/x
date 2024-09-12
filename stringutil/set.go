package stringutil

import (
	"encoding/json"
	"fmt"
	"sort"
)

type StringSet map[string]struct{}

func (set StringSet) ToSlice() []string {
	keys := make([]string, 0, len(set))
	for k := range set {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	return keys
}

func (set StringSet) IsEmpty() bool {
	return len(set) == 0
}

func (set StringSet) Len() int {
	return len(set)
}

func (set StringSet) Add(s string) {
	set[s] = struct{}{}
}

func (set StringSet) Remove(s string) {
	delete(set, s)
}

func (set StringSet) Contains(s string) bool {
	_, ok := set[s]
	return ok
}

func (set StringSet) FuncMatch(matchFn func(string, string) bool, matchString string) StringSet {
	nset := NewStringSet()
	for k := range set {
		if matchFn(k, matchString) {
			nset.Add(k)
		}
	}
	return nset
}

func (set StringSet) ApplyFunc(applyFn func(string) string) StringSet {
	nset := NewStringSet()
	for k := range set {
		nset.Add(applyFn(k))
	}
	return nset
}

func (set StringSet) Equals(o StringSet) bool {
	if len(set) != len(o) {
		return false
	}

	for k := range set {
		if !o.Contains(k) {
			return false
		}
	}

	return true
}

func (set StringSet) Intersection(o StringSet) StringSet {
	nset := NewStringSet()
	for k := range set {
		if o.Contains(k) {
			nset.Add(k)
		}
	}
	return nset
}

func (set StringSet) Difference(o StringSet) StringSet {
	nset := NewStringSet()
	for k := range set {
		if !o.Contains(k) {
			nset.Add(k)
		}
	}
	return nset
}

func (set StringSet) Union(o StringSet) StringSet {
	nset := NewStringSet()
	for k := range set {
		nset.Add(k)
	}
	for k := range o {
		nset.Add(k)
	}
	return nset
}

func (set StringSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(set.ToSlice())
}

func (set *StringSet) UnmarshalJSON(data []byte) error {
	var sl []string
	var err error
	if err = json.Unmarshal(data, &sl); err == nil {
		*set = make(StringSet)
		for _, s := range sl {
			set.Add(s)
		}
	} else {
		var s string
		if err = json.Unmarshal(data, &s); err == nil {
			*set = make(StringSet)
			set.Add(s)
		}
	}

	return err
}

func (set StringSet) String() string {
	return fmt.Sprintf("%s", set.ToSlice())
}

func (set StringSet) Clone() StringSet {
	return NewStringSet(set.ToSlice()...)
}

func NewStringSet(sl ...string) StringSet {
	set := make(StringSet)
	for _, s := range sl {
		set.Add(s)
	}

	return set
}
