package column

import (
	"reflect"
	"testing"
)

type Happy struct {
	Name          string `column:"name"`
	Value         int    `column:"VALUE"`
	UnlabeledData float64
	Omit          string `column:"-"`
	Sads          []Sad
	// TODO implement omitempty
}

type Sad struct {
	Name string
	Data int
}

func TestMarshal(t *testing.T) {
	tests := []struct {
		v           interface{}
		expected    []byte
		err         error
		description string
	}{
		{
			v:           Happy{Name: "test", Value: 2, UnlabeledData: 3.1, Sads: []Sad{{Name: "test", Data: 3}}},
			expected:    []byte("name: test\nVALUE: 2\nUNLABELED DATA: 3.1\nNAME    DATA\ntest    3"),
			description: "struct",
		},
		{
			v: []Happy{{
				Name: "test", Value: 2,
			}, {
				Name: "test-2", Value: 3,
			}},
			expected:    []byte("name      VALUE    UNLABELED DATA    SADS\ntest      2        0                 []\ntest-2    3        0                 []"),
			description: "slice",
		},
		{
			v:           &Happy{Name: "test", Value: 2},
			expected:    []byte("name: test\nVALUE: 2\nUNLABELED DATA: 0\n"),
			description: "pointer",
		},
		{
			v: []*Happy{{
				Name: "test", Value: 2,
			}, {
				Name: "test-2", Value: 3,
			}},
			expected:    []byte("name      VALUE    UNLABELED DATA    SADS\ntest      2        0                 []\ntest-2    3        0                 []"),
			description: "slice",
		},
		{
			v: &[]Happy{{
				Name: "test", Value: 2,
			}, {
				Name: "test-2", Value: 3,
			}},
			expected:    []byte("name      VALUE    UNLABELED DATA    SADS\ntest      2        0                 []\ntest-2    3        0                 []"),
			description: "slice",
		},
	}
	for _, test := range tests {
		res, err := Marshal(test.v)
		if err != test.err {
			t.Errorf("expected error %v, got %v", test.err, err)
		}
		if !reflect.DeepEqual(res, test.expected) {
			t.Errorf("test: %s...expected \n%s\n got \n%s\n", test.description, string(test.expected), string(res))
			t.Log(test.expected)
			t.Log(res)
		}
	}
}
