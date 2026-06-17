package auth

import (
	"testing"
	"time"
	uuid "github.com/google/uuid"
)

func generateJWT(t *testing.T) error {

	id := uuid.New()

	_, err := MakeJWT(id, "secret-string", time.Second)

    return err
}

// ADD MORE