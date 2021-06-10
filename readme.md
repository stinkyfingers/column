### Column marshaler

Print tabbed columns like:

        name    VALUE    Unlabeled
        test    2        3.0

#### Quickstart:

```
type MyStruct struct {
	Name      string `column:"name"`
	Value     int    `column:"VALUE"`
	Unlabeled float64
	Omit      string `column:"-"`
}

myStruct := MyStruct {
	Name: "test",
	Value: 2,
	Unlabeled: 3.1
}

c, err := Marshal(myStruct)
if err != nil {
	return err
}

fmt.Println(string(c))

```
