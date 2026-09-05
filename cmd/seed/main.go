package main

import (
	"log"

	"github.com/l0ng7h0r/golang/internal/repository"
	"github.com/l0ng7h0r/golang/internal/usecase"
	"github.com/l0ng7h0r/golang/pkg/config"
	"github.com/l0ng7h0r/golang/pkg/database"
	"github.com/l0ng7h0r/golang/pkg/security"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()

	db, err := database.NewPostgres(cfg.DBDsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// 1. Ensure standard roles exist in 'roles' table
	roles := []string{"admin", "seller", "user"}
	for _, role := range roles {
		_, err := db.Exec(`INSERT INTO roles (name) VALUES ($1) ON CONFLICT (name) DO NOTHING`, role)
		if err != nil {
			// If ON CONFLICT (name) isn't unique constraint, fallback to NOT EXISTS query
			_, err = db.Exec(`INSERT INTO roles (name) SELECT $1 WHERE NOT EXISTS (SELECT 1 FROM roles WHERE name = $1)`, role)
			if err != nil {
				log.Printf("Warning inserting role %s: %v", role, err)
			}
		}
	}

	userRepo := repository.NewUserRepository(db)
	authUsecase := usecase.NewAuthUsecase(userRepo)

	adminEmail := "admin@gmail.com"
	adminPassword := "admin@123"

	// 2. Check if admin user already exists
	existingUser, err := userRepo.FindByEmail(adminEmail)
	if err == nil && existingUser != nil {
		log.Println("Admin user already exists. Updating password & roles...")
		hashed, err := security.HashPassword(adminPassword)
		if err != nil {
			log.Fatal("Failed to hash password:", err)
		}
		existingUser.Password = hashed
		if err := userRepo.UpdateUser(existingUser); err != nil {
			log.Fatal("Failed to update admin user:", err)
		}

		// Ensure admin role is assigned
		var count int
		_ = db.QueryRow(`SELECT COUNT(*) FROM user_roles ur JOIN roles r ON ur.role_id = r.id WHERE ur.user_id = $1 AND r.name = 'admin'`, existingUser.ID).Scan(&count)
		if count == 0 {
			_, _ = db.Exec(`INSERT INTO user_roles (user_id, role_id) SELECT $1, id FROM roles WHERE name = 'admin'`, existingUser.ID)
		}
		log.Println("Admin user successfully updated!")
		return
	}

	// 3. Create admin user if not exists
	err = authUsecase.CreateUser(adminEmail, adminPassword, []string{"admin"})
	if err != nil {
		log.Fatal("Failed to create admin user:", err)
	}

	log.Println("Admin user successfully seeded!")
	log.Printf("Email: %s", adminEmail)
	log.Printf("Password: %s", adminPassword)
	log.Println("Roles: [admin]")
}
