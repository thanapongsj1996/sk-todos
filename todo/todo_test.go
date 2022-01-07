package todo

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewTodoSleep(t *testing.T) {
	handler := NewTodoHandler(&gorm.DB{})

	w := httptest.NewRecorder()
	payload := bytes.NewBufferString(`{"text":"sleep"}`)
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8080/todos", payload)
	req.Header.Add("TransactionID", "testIDxxx")

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.NewTask(c)

	want := `{"error":"not allowed"}`

	if want != w.Body.String() {
		t.Errorf("want %s but get %s", want, w.Body.String())
	}
}
