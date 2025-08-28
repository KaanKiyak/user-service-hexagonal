package domain

import "testing"

func BenchmarkValidateBusinessRules(b *testing.B) {
	user := &User{
		UUID:     "uuid-123",
		ID:       1,
		Name:     "John Doe",
		Email:    "john.doe@gmail.com",
		Age:      30,
		Password: "1234",
		Role:     "user",
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := user.ValidateBusinessRules(); err != nil {
			b.Fatal(err)
		}
	}
}
