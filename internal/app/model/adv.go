package model

import (
	"fmt"
	"unicode/utf8"
)

type Adv struct {
	ID          int
	Name        string
	Description string
	Price       int
	Date        string
	Ref         []Ref
}

type PageAdv struct {
	ListAdv  []Adv
	NextPage bool
}

func (a *Adv) Check() error {
	if utf8.RuneCountInString(a.Name) > 200 {
		return fmt.Errorf("maximum name length 200 characters")
	}
	if utf8.RuneCountInString(a.Description) > 1000 {
		return fmt.Errorf("maximum description length 1000 characters")
	}
	return nil
}
