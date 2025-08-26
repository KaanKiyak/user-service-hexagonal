package domain

import "fmt"

type User struct {
	UUID     string
	ID       int
	Name     string
	Email    string
	Age      int
	Password string
	Role     string
}

// ValidateBusinessRules domain kurallarını kontrol eder
func (u *User) ValidateBusinessRules() error {
	// Name sadece harf ve boşluk içerebilir
	for _, c := range u.Name {
		if (c < 'A' || c > 'Z') && (c < 'a' || c > 'z') && c != ' ' {
			return fmt.Errorf("isim sadece harf içerebilir boşluk olamaz")
		}
	}

	// Email sadece gmail olabilir (örnek)
	if len(u.Email) < 11 || u.Email[len(u.Email)-10:] != "@gmail.com" {
		return fmt.Errorf("email gmail olmalıdır")
	}

	// Age 18-100 arası olmalı
	if u.Age < 18 || u.Age > 100 {
		return fmt.Errorf("yaş 18 ile 100 arasında olmalıdır")
	}

	// Password minimum 4 karakter olmalı
	if len(u.Password) < 4 {
		return fmt.Errorf("şifre en az 4 karakter olmalıdır")
	}

	// Role sadece "user" olabilir
	if u.Role != "user" {
		return fmt.Errorf("rol sadece 'user' olabilir")
	}

	return nil
}
