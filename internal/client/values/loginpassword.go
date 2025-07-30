package value

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var allowedPasswordChars = regexp.MustCompile(`^[\w!@#\$%\^&\*\(\)_\+\-=]+$`)

const minPasswordSize = 8

type LoginPassword struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (v *LoginPassword) vType() vType {
	return typeLoginPassword
}

func (v *LoginPassword) ToBytes() ([]byte, error) {
	payload, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return append([]byte{byte(typeLoginPassword)}, payload...), nil
}

func (v *LoginPassword) Validate() error {
	if v.Login == "" {
		return errors.New("login is empty")
	}
	if strings.TrimSpace(v.Password) == "" {
		return errors.New("password is empty")
	}
	if len(v.Password) < minPasswordSize {
		return fmt.Errorf("minimal password %d characters", minPasswordSize)
	}
	if !allowedPasswordChars.MatchString(v.Password) {
		return errors.New("password contains invalid characters")
	}
	return nil
}

func (v *LoginPassword) String() string {
	return fmt.Sprintf("Login: %s, Password: %s", v.Login, v.Password)
}
