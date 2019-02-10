package types

import (
	"errors"
)

type DBConfig struct {
	Host string
	User     string
	Password string
	DBName   string
}

func (config DBConfig) ValidateBasic() error {
	if len(config.Host) == 0 {
		return errors.New("DB Host not set")
	}
	if len(config.User) == 0 {
		return errors.New("DB User not set")
	}
	if len(config.Password) == 0 {
		return errors.New("DB Password not set")
	}
	if len(config.DBName) == 0 {
		return errors.New("DB Name not set")
	}
	return nil
}
