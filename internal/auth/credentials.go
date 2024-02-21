package auth

type Datastore interface {
	GetAdminUser(username string) (AdminUser, error)
}

type AdminUser struct {
	Username string
	Password string
}

func AuthenticateAdminCredentials(datastore Datastore, username, password string) (bool, error) {
	adminUser, err := datastore.GetAdminUser(username)
	if err != nil {
		return false, err
	}
	// Implement password comparison (use hashed passwords)
	return adminUser.Password == password, nil
}
