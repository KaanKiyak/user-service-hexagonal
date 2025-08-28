package services

import (
	"testing"

	"user-service-hexagonal/internal/core/domain"

	"golang.org/x/crypto/bcrypt"
)

// mockRepo, ports.UserRepository'i implemente eder
type mockRepo struct{}

func (m *mockRepo) CreateUser(user *domain.User) (*domain.User, error) {
	if user.ID == 0 {
		user.ID = 1
	}
	if user.UUID == "" {
		user.UUID = "mock-uuid-1"
	}
	return user, nil
}

func (m *mockRepo) ReadUser(id int) (*domain.User, error) {
	return &domain.User{
		ID:    id,
		Email: "test@gmail.com",
		Role:  "user",
		Age:   25,
	}, nil
}

func (m *mockRepo) ReadUsers() ([]*domain.User, error) { return []*domain.User{}, nil }

func (m *mockRepo) UpdateUser(id, email, password string) error { return nil }

func (m *mockRepo) DeleteUser(id string) error { return nil }

func (m *mockRepo) LoginUser(email string) (*domain.User, error) {
	// Burada her çağrıda hash üretildiği için şifre uyumsuzluğu asla olmaz.
	hash, _ := bcrypt.GenerateFromPassword([]byte("1234"), bcrypt.DefaultCost)

	return &domain.User{
		Email:    email,
		Password: string(hash), // bcrypt("1234")
		Role:     "user",
		Age:      25,
	}, nil
}

// Benchmark: sadece LoginUser performansı (bcrypt compare + JWT üretimi)
func BenchmarkUserServiceLogin(b *testing.B) {

	service := NewUserService(&mockRepo{})

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := service.LoginUser("test@gmail.com", "1234"); err != nil {
			b.Fatal(err)
		}
	}
}
