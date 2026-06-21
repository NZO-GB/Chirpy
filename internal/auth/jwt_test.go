package auth

import (
	"testing"
	"time"
	"fmt"
	uuid "github.com/google/uuid"
)

func generateJWT(t *testing.T) error {

	id := uuid.New()

	_, err := MakeJWT(id, "secret-string", time.Second)

    return err
}

func generateValidateJWT (t *testing.T) error {

	id := uuid.New()

	secretString := "secret-string"

	tokenString, err := MakeJWT(id, secretString, 3 * time.Second)
	if err != nil {
		return fmt.Errorf("Error making JWT: %v", err)
	}

	validatedID, err := ValidateJWT(tokenString, secretString)
    if err != nil {
		return fmt.Errorf("Error validating JWT: %v", err)
	}

	if validatedID != id {
		return fmt.Errorf("Error comapring ids: %v", err)
	}

	return nil

}