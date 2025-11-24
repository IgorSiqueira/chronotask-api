package port

// HasherService defines the interface for password hashing (Port)
// This follows DIP - the application layer defines the interface,
// and the infrastructure layer provides the implementation
type HasherService interface {
	// Hash generates a hashed version of the password
	Hash(password string) (string, error)

	// Compare compares a hashed password with a plain text password
	Compare(hashedPassword, password string) error
}
