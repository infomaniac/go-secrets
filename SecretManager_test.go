package secman

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const VALID_PROJECT_ID = "<project-id>"

func TestNewSetGetDeleteClose(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s, err := New(ctx, VALID_PROJECT_ID, "europe-west1", "europe-west6")
	assert.NoError(t, err)
	assert.NotNil(t, s)
	defer s.Close()

	secretName := fmt.Sprintf("test-secret-%d", time.Now().UnixNano())
	secret := &bytes.Buffer{}
	secret.WriteString(fmt.Sprintf("secret-%d", time.Now().UnixNano()))

	err = s.Set(secretName, secret.Bytes())
	assert.NoError(t, err)

	got, err := s.Get(secretName)
	assert.NoError(t, err)
	assert.EqualValues(t, secret.Bytes(), got)

	err = s.Delete(secretName)
	assert.NoError(t, err)
}

func TestNewInvalidContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	s, err := New(ctx, "")
	assert.Error(t, err)
	assert.Nil(t, s)
}

func TestGetInvalid(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s, err := New(ctx, "non-existing-test-project-12342")
	assert.NoError(t, err)
	assert.NotNil(t, s)
	defer s.Close()

	sec, err := s.Get("non-existing-secret")
	assert.Error(t, err)
	assert.Nil(t, sec)
}

func TestSetTwice(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	secretName := fmt.Sprintf("test-secret-%d", time.Now().UnixNano())

	s, err := New(ctx, VALID_PROJECT_ID, "europe-west1", "europe-west6")
	assert.NoError(t, err)
	assert.NotNil(t, s)
	defer s.Close()

	secret := &bytes.Buffer{}
	secret.WriteString(fmt.Sprintf("secret-%d", 1))

	err = s.Set(secretName, secret.Bytes())
	assert.NoError(t, err)

	got1, err := s.Get(secretName)
	assert.NoError(t, err)
	assert.EqualValues(t, secret.Bytes(), got1)

	secret.Reset()
	secret.WriteString(fmt.Sprintf("secret-%d", 2))

	err = s.Set(secretName, secret.Bytes())
	assert.NoError(t, err)

	got2, err := s.Get(secretName)
	assert.NoError(t, err)
	assert.EqualValues(t, secret.Bytes(), got2)

	err = s.Delete(secretName)
	assert.NoError(t, err)
}

func TestCreateInvalid(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s, err := New(ctx, "non-existing-test-project-12342")
	assert.NoError(t, err)
	assert.NotNil(t, s)
	defer s.Close()

	err = s.Set("non-existing-secret", []byte("secret"))
	assert.Error(t, err)
}
