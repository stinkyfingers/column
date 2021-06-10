package column

import (
	"bytes"
	"fmt"
	"reflect"
	"text/tabwriter"
)

type val struct {
	tag   reflect.StructTag
	value reflect.Value
}

type valSet []val

func Marshal(v interface{}) ([]byte, error) {
	var valSets []valSet

	value := reflect.ValueOf(v)
	switch value.Kind() {
	case reflect.Struct:
		valSet := newValSet(value)
		valSets = append(valSets, valSet)
	case reflect.Slice:
		for i := 0; i < reflect.Value.Len(value); i++ {
			valSet := newValSet(value.Index(i))
			valSets = append(valSets, valSet)
		}
	case reflect.Ptr:
		valSet := newValSet(value.Elem())
		valSets = append(valSets, valSet)
	default:
		return nil, fmt.Errorf("unsupported kind: %s", value.Kind())
	}

	// no data
	if len(valSets) < 1 {
		return nil, nil
	}

	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 4, 4, ' ', 0)

	// headers
	for i, val := range valSets[0] {
		tag := val.tag.Get("column")
		// omit?
		if tag == "-" {
			continue
		}

		fmt.Fprint(w, tag)
		// tab
		if i+1 < len(valSets[0]) {
			fmt.Fprint(w, "\t")
		}
	}
	fmt.Fprint(w, "\n")

	// fields
	for i, valSet := range valSets {
		for j, val := range valSet {
			// omit?
			if val.tag.Get("column") == "-" {
				continue
			}

			fmt.Fprint(w, val.value)

			// tab
			if j+1 < len(valSet) {
				fmt.Fprint(w, "\t")

			}
		}
		// newline
		if i+1 < len(valSets) {
			fmt.Fprint(w, "\n")
		}
	}
	err := w.Flush()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// newValSet iterates over a value's fields and assigns the Value and StructTag to a valSet
func newValSet(value reflect.Value) valSet {
	var valSet valSet
	for i := 0; i < value.NumField(); i++ {
		tag := value.Type().Field(i).Tag
		// use Field Type as StructTag
		if tag.Get("column") == "" {
			tag = reflect.StructTag(fmt.Sprintf("column:\"%s\"", value.Type().Field(i).Name))
		}
		valSet = append(valSet, val{tag: tag, value: value.Field(i)})
	}
	return valSet
}
