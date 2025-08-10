package headers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHeaders(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	require.Equal(t, "localhost:42069", headers["host"])
	require.Equal(t, 23, n)
	require.False(t, done)

	// Test: Valid single header with extra whitespace
	headers = NewHeaders()
	data = []byte("   Host:   localhost:42069   \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	require.Equal(t, "localhost:42069", headers["host"])
	require.Equal(t, 31, n)
	require.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	require.Equal(t, 0, n)
	require.False(t, done)

	// Test: Valid 2 headers
	headers = NewHeaders()
	data = []byte("Host:localhost:42069\r\n Content-Type: Application/Json  \r\n")
	n, done, err = headers.Parse(data)
	require.Nil(t, err)
	require.NotNil(t, headers)
	require.Equal(t, "localhost:42069", headers["host"])
	require.Equal(t, 22, n)
	require.False(t, done)
	n, done, err = headers.Parse(data[n:])
	require.Nil(t, err)
	require.NotNil(t, headers)
	require.Equal(t, "Application/Json", headers["content-type"])
	require.Equal(t, 35, n)
	require.False(t, done)

	// Test: Valid done
	headers = NewHeaders()
	data = []byte("\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	require.Equal(t, 2, n)
	require.True(t, done)


	// Test: Invalid header value, containing space
	headers = NewHeaders()
	data = []byte("Host: localhost with a space\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	require.Equal(t, 0, n)
	require.False(t, done)

	// Test: invalid characters in header key
	headers = NewHeaders()
	data = []byte("MySuperCoolHeader@Me: localhost\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	require.Equal(t, 0, n)
	require.False(t, done)
}
