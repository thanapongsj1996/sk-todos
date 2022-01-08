package todo

import (
	"testing"
)

func TestNewTodoSleep(t *testing.T) {
	handler := NewTodoHandler(&TestDB{})
	c := &TextContext{}

	handler.NewTask(c)

	want := "not allowed"

	if want != c.v["error"] {
		t.Errorf("want %s but get %s", want, c.v["error"])
	}
}

type TestDB struct{}

func (TestDB) New(*Todo) error {
	return nil
}

type TextContext struct {
	v map[string]interface{}
}

func (TextContext) Bind(v interface{}) error {
	*v.(*Todo) = Todo{
		Title: "sleep",
	}
	return nil
}
func (c *TextContext) JSON(code int, v interface{}) {
	c.v = v.(map[string]interface{})
}
func (TextContext) TransactionID() string {
	return "TransactionID"
}
func (TextContext) Audience() string {
	return "Unit test"
}
