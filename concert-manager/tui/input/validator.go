package input

import (
	"concert-manager/util"
	"errors"
	"unicode"
)

func NoValidation(_ string) error {
	return nil
}

func OnlyLettersValidation(in string) error {
	for _, c := range in {
		if !unicode.IsLetter(c) {
			return errors.New("all characters must be letters")
		}
	}
	return nil
}

func OnlyLettersOrSpacesValidation(in string) error {
	for _, c := range in {
		if !unicode.IsLetter(c) && !unicode.IsSpace(c) {
			return errors.New("only letters and spaces are allowed")
		}
	}
	return nil
}

func StateValidation(in string) error {
	if len(in) != 2 {
		return errors.New("state code must be two characters")
	}
	return OnlyLettersValidation(in)
}

func DateValidation(date string) error {
	if !util.ValidDate(date) {
		return errors.New("expected date format is mm/dd/yyyy")
	}
	return nil
}
