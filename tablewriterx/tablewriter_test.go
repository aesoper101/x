package tablewriterx

import "testing"

type TestStruct struct {
	S struct {
		Name string
		Age  int
	} `tablewriter:"-"`
	Name string `tablewriter:"n测试ame"`
	Age  int    `json:"age"`
}

func TestSetStructs(t *testing.T) {
	tt := NewWriter()
	tt.SetStructs(
		[]TestStruct{
			TestStruct{
				S: struct {
					Name string
					Age  int
				}{
					Name: "test",
					Age:  10,
				},
				Name: "222",
				Age:  222,
			},
		},
	)
	tt.Render()
}
