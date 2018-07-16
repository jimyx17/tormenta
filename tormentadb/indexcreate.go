package tormentadb

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"strings"

	"github.com/dgraph-io/badger"
	"github.com/jpincas/gouuidv6"
)

// Index format
// i:indexname:root:indexcontent:entityID
// i:order:customer:5:324ds-3werwf-234wef-23wef

const (
	tormentaTag      = "tormenta"
	tormentaTagIndex = "index"
)

func index(txn *badger.Txn, entity Tormentable, keyRoot []byte, id gouuidv6.UUID) error {
	v := reflect.Indirect(reflect.ValueOf(entity))

	for i := 0; i < v.NumField(); i++ {

		// Look for the 'tormenta:index' tag
		fieldType := v.Type().Field(i)
		if idx := fieldType.Tag.Get(tormentaTag); idx == tormentaTagIndex {

			switch fieldType.Type.Kind() {
			case reflect.Slice:
				if err := setMultipleIndexes(txn, v.Field(i), keyRoot, id, fieldType.Name); err != nil {
					return err
				}
			case reflect.Array:
				if err := setMultipleIndexes(txn, v.Field(i), keyRoot, id, fieldType.Name); err != nil {
					return err
				}
			default:
				if err := setIndexKey(txn, keyRoot, id, fieldType.Name, v.Field(i).Interface()); err != nil {
					return err
				}
			}

		}
	}

	return nil
}

func setMultipleIndexes(txn *badger.Txn, v reflect.Value, root []byte, id gouuidv6.UUID, indexName string) error {
	for i := 0; i < v.Len(); i++ {
		if err := setIndexKey(txn, root, id, indexName, v.Index(i).Interface()); err != nil {
			return err
		}
	}

	return nil
}

func setIndexKey(txn *badger.Txn, root []byte, id gouuidv6.UUID, indexName string, indexContent interface{}) error {
	key := makeIndexKey(root, id, indexName, indexContent)

	// Set the index key with no content
	return txn.Set(key, []byte{})
}

// IndexKey returns the bytes of an index key constructed for a particular root, id, index name and index content
func IndexKey(root []byte, id gouuidv6.UUID, indexName string, indexContent interface{}) []byte {
	return makeIndexKey(root, id, indexName, indexContent)
}

func makeIndexKey(root []byte, id gouuidv6.UUID, indexName string, indexContent interface{}) []byte {
	return bytes.Join(
		[][]byte{
			[]byte(indexKeyPrefix),
			root,
			[]byte(strings.ToLower(indexName)),
			interfaceToBytes(indexContent),
			id.Bytes(),
		},
		[]byte(keySeparator),
	)
}

func interfaceToBytes(value interface{}) []byte {
	// Must use BigEndian for correct sorting

	// Empty interface -> empty byte slice
	if value == nil {
		return []byte{}
	}

	// Set up buffer for writing binary values
	buf := new(bytes.Buffer)

	switch value.(type) {
	// For ints, cast the interface to int, convert to uint32 (normal ints don't work)
	case int:
		binary.Write(buf, binary.BigEndian, uint32(value.(int)))
		return buf.Bytes()
	// For floats, write straight to binary
	case float64:
		binary.Write(buf, binary.BigEndian, value.(float64))
		return buf.Bytes()
	}

	// Everything else as a string
	return []byte(fmt.Sprintf("%v", value))
}