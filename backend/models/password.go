package models

import "golang.org/x/crypto/bcrypt"

type PasswordType struct {
	Text *string
	Hash []byte
}

// Set sets the password to the hash of the password_txt
func (p *PasswordType) Set(password_txt string) error {
	var err error
	p.Hash, err = bcrypt.GenerateFromPassword([]byte(password_txt), bcrypt.MinCost)
	p.Text = &password_txt
	if err != nil {
		return err
	}
	return nil
}

func (p *PasswordType) Compare(pass string) error {
	if err := bcrypt.CompareHashAndPassword(p.Hash, []byte(pass)); err != nil {
		return err
	}
	return nil
}
