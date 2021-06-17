package column

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"text/tabwriter"
)

type val struct {
	tag   reflect.StructTag
	value reflect.Value
}

type valSet []val

type valSets []valSet

func Marshal(v interface{}) ([]byte, error) {
	return drillDownToStructOrSlice(reflect.ValueOf(v))
}

func drillDownToStructOrSlice(value reflect.Value) ([]byte, error) {
	var buf bytes.Buffer
	switch value.Kind() {
	case reflect.Struct:
		valSet := newValSet(value)
		b, err := valSet.marshal()
		if err != nil {
			return nil, err
		}
		buf.Write(b)

	case reflect.Slice:
		var vSets valSets
		for i := 0; i < reflect.Value.Len(value); i++ {
			valSet := newValSet(value.Index(i))
			vSets = append(vSets, valSet)
		}
		b, err := vSets.marshal()
		if err != nil {
			return nil, err
		}
		buf.Write(b)

	case reflect.Ptr:
		return drillDownToStructOrSlice(value.Elem())

	default:
		return nil, fmt.Errorf("unsupported kind: %s", value.Kind())
	}
	return buf.Bytes(), nil
}

func (v *valSet) marshal() ([]byte, error) {
	var buf bytes.Buffer
	for _, val := range *v {
		if val.value.Kind() == reflect.Slice || val.value.Kind() == reflect.Ptr {
			b, err := drillDownToStructOrSlice(val.value)
			if err != nil {
				return nil, err
			}
			buf.Write(b)
			continue
		}
		tag := val.tag.Get("column")

		// omit?
		if tag == "-" {
			continue
		}
		fmt.Fprintf(&buf, "%s: %v\n", tag, val.value)
	}
	return buf.Bytes(), nil
}

func (vs *valSets) marshal() ([]byte, error) {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 4, 4, ' ', 0)

	if len(*vs) == 0 {
		return nil, nil
	}

	// headers
	for i, val := range (*vs)[0] {
		tag := val.tag.Get("column")
		// omit?
		if tag == "-" {
			continue
		}

		fmt.Fprint(w, tag)
		// tab
		if i+1 < len((*vs)[0]) {
			fmt.Fprint(w, "\t")
		}
	}
	fmt.Fprint(w, "\n")

	// fields
	for i, valSet := range *vs {
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
		if i+1 < len(*vs) {
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
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	var valSet valSet
	for i := 0; i < value.NumField(); i++ {
		tag := value.Type().Field(i).Tag
		// use Field Type as StructTag
		if tag.Get("column") == "" {
			fieldName := value.Type().Field(i).Name
			// split and uppercase field name
			var builder strings.Builder
			for i, r := range fieldName {
				if r > 96 && r < 123 {
					builder.WriteString(string(r - 32))
				} else {
					if i > 0 {
						builder.WriteString(" ")
					}
					builder.WriteString(string(r))
				}
			}
			tag = reflect.StructTag(fmt.Sprintf("column:\"%s\"", builder.String()))
		}
		valSet = append(valSet, val{tag: tag, value: value.Field(i)})
	}
	return valSet
}
