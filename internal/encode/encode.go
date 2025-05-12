package encode

import (
	"bytes"
	"encoding/json"
	"log"
	"sort"
	"strings"
	"unsafe"

	// generate/storage import will register all types in the protoregistry.GlobalTypes list
	_ "github.com/stackrox/rox/generated/storage"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

const (
	storagePrefix = "storage."
)

type thinger interface {
	proto.Message
	UnmarshalVTUnsafe([]byte) error
}
type EncodeEntry struct {
	Name      string `json:"name,omitempty"`
	NumFields int    `json:"num_fields,omitempty"`
	SizeOf    uint   `json:"size_of,omitempty"`
	DataLen   int    `json:"data_len,omitempty"`
	ProtoJSON []byte `json:"proto_json,omitempty"`
	RawJSON   []byte `json:"raw_json,omitempty"`
}

// JSONAll attempts to decode the provided bytes into every known proto message.
// The returned slice is sorted (roughly) based on the likelyhood of a successful
// match.
func JSONAll(dataB []byte) ([]*EncodeEntry, error) {
	results := []*EncodeEntry{}
	protoregistry.GlobalTypes.RangeMessages(func(mt protoreflect.MessageType) bool {
		name := string(mt.Descriptor().FullName())
		if skip(name) {
			return true
		}
		name = cleanName(name)

		t, ok := mt.New().Interface().(thinger)
		if !ok {
			return true
		}

		err := t.UnmarshalVTUnsafe(dataB)
		if err != nil {
			return true
		}
		entry := &EncodeEntry{}
		entry.Name = name
		entry.NumFields = mt.Descriptor().Fields().Len()
		entry.SizeOf = uint(unsafe.Sizeof(t))

		b, err := protojson.MarshalOptions{Multiline: true}.Marshal(t)
		if err == nil && len(bytes.TrimSpace(b)) > 0 {
			entry.ProtoJSON = b

			// only try to get raw json if protojson succeeded
			b, err = json.MarshalIndent(t, "", "  ")
			if err == nil && len(bytes.TrimSpace(b)) > 0 {
				entry.RawJSON = b
			}

			// hack
			entry.DataLen = dataLen(entry.RawJSON)

			results = append(results, entry)
		}

		return true
	})

	sort.Slice(results, func(i, j int) bool {
		l := results[i]
		r := results[j]

		if len(l.RawJSON) != len(r.RawJSON) {
			return len(l.RawJSON) > len(r.RawJSON)
		}

		if len(l.ProtoJSON) != len(r.ProtoJSON) {
			return len(l.ProtoJSON) > len(r.ProtoJSON)
		}

		if l.DataLen != r.DataLen {
			return l.DataLen > r.DataLen
		}

		if l.NumFields != r.NumFields {
			return l.NumFields > r.NumFields
		}
		return l.Name < r.Name
	})

	return results, nil
}

func dataLen(b []byte) int {
	d := map[string]interface{}{}
	err := json.Unmarshal(b, &d)
	if err != nil {
		log.Print(err)
		return -1
	}

	return dataLenR(d)
}

func dataLenR(d map[string]interface{}) int {
	var sum int

	queue := []interface{}{}
	for _, v := range d {
		queue = append(queue, v)
	}

	for len(queue) > 0 {
		v := queue[0]
		sum += dataLenField(v)
		queue = queue[1:]
	}

	return sum
}

func dataLenField(v interface{}) int {
	switch t := v.(type) {
	case string:
		return len(t)
	case map[string]interface{}:
		return dataLenR(t)
	case []interface{}:
		sum := 0
		for _, e := range t {
			sum += dataLenField(e)
		}
		return sum
	default:
		// fmt.Printf("huh: %v (%T)\n", v, t)
		return int(unsafe.Sizeof(t))
	}
}

// KnownTypes returns of globally registered proto types that
// are candidated for JSON decoding/encoding.
func KnownTypes() []string {
	types := []string{}
	protoregistry.GlobalTypes.RangeMessages(func(mt protoreflect.MessageType) bool {
		name := string(mt.Descriptor().FullName())
		if skip(name) {
			return true
		}
		name = cleanName(name)
		types = append(types, name)
		return true
	})

	sort.Strings(types)

	return types
}

func skip(name string) bool {
	if !strings.HasPrefix(name, storagePrefix) {
		return true
	}

	name = cleanName(name)
	if strings.Index(name, ".") > 0 {
		// this ensures things like storage.TestStruct.Embedded do not get attempted
		// only 'top level' types are supported
		return true
	}

	return false
}

func cleanName(name string) string {
	if strings.HasPrefix(name, storagePrefix) {
		return strings.TrimPrefix(name, storagePrefix)
	}

	return name
}
