package json2map

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type JsonKey struct {
	Key string
	// if Index >= 0, it means that the object of this JsonKey is a json array,
	// and Index is the key's index of the array
	Index int
}

func (self *JsonKey) String() string {
	if self.Index >= 0 {
		return fmt.Sprintf("%s[%d]", self.Key, self.Index)
	}
	return self.Key
}

type Json2csv struct {
	// the leading name representing root, default is empty
	RootKeyName    string
	JsonKey2MapKey func([]JsonKey) string

	result map[string]interface{}
}

func New() *Json2csv {

	self := &Json2csv{
		result: make(map[string]interface{}),
		JsonKey2MapKey: func(keys []JsonKey) string {
			strs := make([]string, len(keys))
			for i, key := range keys {
				strs[i] = key.String()
			}
			return strings.Join(strs, ".")
		},
	}
	return self
}

func NewWithRootKeyName(rootKeyName string) *Json2csv {
	self := New()
	self.RootKeyName = rootKeyName
	return self
}

func (self *Json2csv) Convert(r io.Reader) (map[string]interface{}, error) {

	var object interface{}

	d := json.NewDecoder(r)
	d.UseNumber()
	if err := d.Decode(&object); err != nil {
		return nil, err
	}

	self.doJsonNode([]JsonKey{{Key: self.RootKeyName, Index: -1}}, object)
	return self.result, nil
}

func (self *Json2csv) doJsonNode(lastKeys []JsonKey, object interface{}) {

	switch object.(type) {
	case map[string]interface{}:
		for k, v := range object.(map[string]interface{}) {
			if v == nil {
				continue
			}
			keys := make([]JsonKey, len(lastKeys))
			copy(keys, lastKeys)
			keys = append(keys, JsonKey{Key: k, Index: -1})
			self.doJsonNode(keys, v)
		}
	case []interface{}:
		for i, v := range object.([]interface{}) {
			if v == nil {
				continue
			}
			keys := make([]JsonKey, len(lastKeys))
			copy(keys, lastKeys)
			keys[len(keys)-1].Index = i
			self.doJsonNode(keys, v)
		}
	case json.Number, string, bool:
		mapKey := self.JsonKey2MapKey(lastKeys)
		self.result[mapKey] = object
	case nil:
		// do nothing
	}
}
