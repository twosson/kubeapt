package content

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStringText(t *testing.T) {
	st := NewStringText("foo")

	data, err := json.Marshal(st)
	require.NoError(t, err)

	expected := `{"text":"foo","type":"string"}`

	assert.Equal(t, expected, string(data))
}

func TestTimeText(t *testing.T) {
	tt := NewTimeText("2018-11-08T17:55:45Z")

	data, err := json.Marshal(tt)
	require.NoError(t, err)

	expected := `{"time":"2018-11-08T17:55:45Z","type":"time"}`

	assert.Equal(t, expected, string(data))
}

func TestLinkText(t *testing.T) {
	lt := NewLinkText("foo", "/bar")

	data, err := json.Marshal(lt)
	require.NoError(t, err)

	expected := `{"ref":"/bar","text":"foo","type":"link"}`

	assert.Equal(t, expected, string(data))
}

func TestLabelsText(t *testing.T) {
	m := map[string]string{
		"foo": "bar",
	}
	lt := NewLabelsText(m)

	data, err := json.Marshal(lt)
	require.NoError(t, err)

	expected := `{"labels":{"foo":"bar"},"type":"labels"}`

	assert.Equal(t, expected, string(data))
}

func TestListText(t *testing.T) {
	list := []string{"foo", "bar"}

	lt := NewListText(list)

	data, err := json.Marshal(lt)
	require.NoError(t, err)

	expected := `{"list":["foo","bar"],"type":"list"}`

	assert.Equal(t, expected, string(data))

}
