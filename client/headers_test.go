package client

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeaders(t *testing.T) {

	t.Run("not empty init state", func(t *testing.T) {
		headerName1 := "header1"
		headerValue1 := "value1"
		initState := map[string]string{
			headerName1: headerValue1,
		}
		headers := NewHeaders(initState)
		headerName2 := "test-header"
		headerValue2 := "header-value"
		headers.Set(headerName2, headerValue2)
		assert.Equal(t, headerValue1, headers.Get(headerName1))
		assert.Equal(t, headerValue2, headers.Get(headerName2))
		keys := headers.Keys()
		sort.Strings(keys)
		assert.Equal(t, []string{headerName1, headerName2}, keys)
		assert.Equal(t, map[string]string{headerName1: headerValue1, headerName2: headerValue2}, headers.Value())
	})

	t.Run("empty init state", func(t *testing.T) {
		headers := NewHeaders(nil)
		headerName := "test-header"
		headerValue := "header-value"
		headers.Set(headerName, headerValue)
		assert.Equal(t, headerValue, headers.Get(headerName))
		assert.Equal(t, []string{headerName}, headers.Keys())
		assert.Equal(t, map[string]string{headerName: headerValue}, headers.Value())
	})
}
