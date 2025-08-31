package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaders(t *testing.T) {
	// Test: Valid spacing header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Valid spacing header and lowercase
	headers = NewHeaders()
	data = []byte("hosT: locaLHost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "locaLHost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid character usage
	headers = NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Same header key
	headers = NewHeaders()
	data1 := []byte("Set-Person: lane-loves-go\r\nSet-Person: prime-loves-zig\r\nSet-Person: tj-loves-ocaml\r\n")
	data2 := []byte("Set-Person: prime-loves-zig\r\nSet-Person: tj-loves-ocaml\r\n")
	data3 := []byte("Set-Person: tj-loves-ocaml\r\n")
	n, done, err = headers.Parse(data1)
	require.NoError(t, err)
	require.NotNil(t, headers)

	n, done, err = headers.Parse(data2)
	require.NoError(t, err)
	require.NotNil(t, headers)

	n, done, err = headers.Parse(data3)
	require.NoError(t, err)
	require.NotNil(t, headers)

	assert.Equal(t, "lane-loves-go, prime-loves-zig, tj-loves-ocaml", headers["set-person"])
	assert.False(t, done)

	// Test: Same header key
	headers = map[string]string{"host": "localhost:8000"}
	data = []byte("Host: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:8000, localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)
}
