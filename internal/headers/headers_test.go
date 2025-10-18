package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaders(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	headers = NewHeaders()
	data = []byte("     Host: localhost:42069     \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 33, n)
	assert.False(t, done)

	headers = Headers{
		"host": "localhost:42069", // existing header
	}

	// Data with two valid headers
	data = []byte("User-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n")

	// First header ("User-Agent")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.False(t, done)
	assert.Equal(t, 25, n) // "User-Agent: curl/7.81.0\r\n" = 27 bytes
	assert.Equal(t, "curl/7.81.0", headers["user-agent"])
	assert.Equal(t, "localhost:42069", headers["host"]) // should still exist

	// Remaining bytes (simulate next read for the next header)
	remaining := data[n:]

	// Second header ("Accept")
	n2, done, err := headers.Parse(remaining)
	require.NoError(t, err)
	assert.False(t, done)
	assert.Equal(t, 13, n2) // "Accept: */*\r\n" = 14 bytes
	assert.Equal(t, "*/*", headers["accept"])
	assert.Equal(t, "curl/7.81.0", headers["user-agent"])
	assert.Equal(t, "localhost:42069", headers["host"])

	// Simulate parsing the final CRLF that marks end of headers
	final := remaining[n2:]
	n3, done, err := headers.Parse(final)
	require.NoError(t, err)
	assert.True(t, done)
	assert.Equal(t, len(final), n3)

	//Invalid field name
	headers = NewHeaders()
	data = []byte("H@st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Multi Value field name

	headers = NewHeaders()
	data = []byte("user: localhost:42069\r\n\r\n")
	_, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.False(t, done)
	data = []byte("user: localhost:42068\r\n\r\n")
	_, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.False(t, done)
	assert.Equal(t, "localhost:42069,localhost:42068", headers["user"])
}
