package auth

import (
	"testing"
    "fmt"
	"github.com/alexedwards/argon2id"
)

// TestCreateHash calls argon2id.CreateHash with a password
// for a valid return value.

func compareHashPass(password string) error {

    hashedPassword, err := argon2id.CreateHash(password, argon2id.DefaultParams)
    if err != nil {
        return fmt.Errorf("Error hashing: %s", password)
    }
    _, err = argon2id.ComparePasswordAndHash(password, hashedPassword)
    if err != nil {
        return fmt.Errorf("Hash: %s doesn't match password: %s", hashedPassword, password)
    }

    fmt.Println(password, hashedPassword)

    return nil
}
func TestHashGladys(t *testing.T) {
    err := compareHashPass("Gladys")
    if err != nil {
        t.Error(err)
    }
}

func TestHashLong(t *testing.T) {
    err := compareHashPass("this is a very long password with characters @!ªº`[]2323")
    if err != nil {
        t.Error(err)
    }
}

