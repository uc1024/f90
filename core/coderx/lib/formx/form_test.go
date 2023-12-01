package formx

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecode(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080?age=18&name=coderx", nil)
	assert.NoError(t, err)
	req.URL.Query()
	user := struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{}
	t.Log(req.URL.Query())
	err = decoder.Decode(&user, req.URL.Query())
	assert.NoError(t, err)
	t.Log(req.URL.Query())
	assert.Equal(t, "coderx", user.Name)
	assert.Equal(t, 18, user.Age)
}

func TestEncode(t *testing.T) {
	user := struct {
		Name string `form:"name"`
		Age  int    `form:"age"`
	}{
		Name: "coderx",
		Age:  18,
	}
	v, err := encoder.Encode(&user)
	assert.NoError(t, err)
	assert.Equal(t, "age=18&name=coderx", v.Encode())
}

func TestMarshal(t *testing.T) {
	user := struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		Name: "coderx",
		Age:  18,
	}
	formUrl := &FormUrl{}
	v, err := formUrl.Marshal(&user)
	assert.NoError(t, err)
	assert.Equal(t, "age=18&name=coderx", string(v))
}

func TestUnmarshal(t *testing.T) {
	user := struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{}
	formUrl := &FormUrl{}
	err := formUrl.Unmarshal([]byte("age=18&name=coderx"), &user)
	assert.NoError(t, err)
	assert.Equal(t, "coderx", user.Name)
	assert.Equal(t, 18, user.Age)
}
