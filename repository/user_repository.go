package repository

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"instructor-led-app/config"
	"instructor-led-app/entity"
	"instructor-led-app/entity/dto"
	"instructor-led-app/shared/model"
	"log"
	"math"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Get(id string) (entity.User, error)
	List(page, size int) ([]entity.User, model.Paging, error)
	Created(data entity.User) (entity.User, error)
	CreateByCsv(filePath string) ([]entity.User, error)
	Updated(id string, data entity.User) (entity.User, error)
	Delete(id string) (entity.User, error)
	GetUser(email string) (entity.User, error)
	GetUserByName(name string) (entity.User, error)
	GetUserIDByName(name string) (dto.UserId, error)
}

type userRepository struct {
	db *sql.DB
}

// GetUserIDByName implements UserRepository.
func (t *userRepository) GetUserIDByName(name string) (dto.UserId, error) {
	var userId dto.UserId

	// Query user data based on the email
	err := t.db.QueryRow(config.GetUserIDbyName, name).Scan(&userId.Id)
	if err != nil {
		log.Println("userRepository.GetUser:", err.Error())
		return dto.UserId{}, err
	}

	return userId, nil
}

// GetUserByName implements UserRepository.
func (t *userRepository) GetUserByName(name string) (entity.User, error) {
	var user entity.User

	// Query user data based on the email
	err := t.db.QueryRow(`SELECT id, name, email, username, address, role, hash_password FROM users WHERE name = $1`, name).Scan(&user.Id, &user.Name, &user.Email, &user.Username, &user.Address, &user.Role, &user.Hashpassword)
	if err != nil {
		log.Println("userRepository.GetUser:", err.Error())
		return entity.User{}, err
	}

	return user, nil
}

// GetUser implements UserRepository.
func (t *userRepository) GetUser(email string) (entity.User, error) {
	var user entity.User

	// Query user data based on the email
	err := t.db.QueryRow(`SELECT id, name, email, username, address, role, hash_password FROM users WHERE email = $1`, email).Scan(&user.Id, &user.Name, &user.Email, &user.Username, &user.Address, &user.Role, &user.Hashpassword)
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle case when no rows are returned
			return entity.User{}, fmt.Errorf("user not found with email: %s", email)
		}
		log.Println("userRepository.GetUser:", err.Error())
		return entity.User{}, err
	}

	return user, nil

}

// CreateByCsv implements UserRepository.
func (t *userRepository) CreateByCsv(filePath string) ([]entity.User, error) {
	var users []entity.User

	// Buka file CSV
	file, err := os.Open(filePath)
	if err != nil {
		log.Println("userRepository.CreateByCsv: Error Opening CSV file:", err.Error())
		return nil, err
	}
	defer file.Close()

	// Buat pembaca CSV
	reader := csv.NewReader(file)

	// Baca semua baris CSV
	lines, err := reader.ReadAll()
	if err != nil {
		log.Println("userRepository.CreateByCsv: Error Reading CSV file:", err.Error())
		return nil, err
	}

	// Iterasi semua baris CSV
	for _, line := range lines {
		// Pastikan jumlah kolom sesuai dengan ekspektasi
		if len(line) != 6 {
			log.Println("userRepository.CreateByCsv: Invalid number of columns in CSV line")
			continue
		}

		// Buat user baru dari baris CSV
		user := entity.User{
			Name:         line[0],
			Email:        line[1],
			Username:     line[2],
			Address:      line[3],
			Hashpassword: line[4],
			Role:         line[5],
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		// Save the user
		createdUser, err := t.getCreatedUser(user)
		if err != nil {
			log.Println("userRepository.CreateByCsv: Error creating user:", err.Error())
			// Handle error as needed, for example, log it and continue processing other lines
			continue
		}

		// Append the created user to the slice
		users = append(users, createdUser)

		// Optionally, you can use the createdUser if needed
		fmt.Printf("User created: %v\n", createdUser)
	}

	// Log the total number of users created
	log.Printf("userRepository.CreateByCsv: Total users created: %d\n", len(users))

	return users, nil
}

func (t *userRepository) getCreatedUser(user entity.User) (entity.User, error) {
	//  Hash Password menggunakan Bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Hashpassword), bcrypt.DefaultCost)
	if err != nil {
		log.Println("userRepository.GenerateFromPassword:", err.Error())
		return entity.User{}, err
	}

	var createdUser entity.User
	// Insert and Returning id, role, created_at
	err = t.db.QueryRow(config.InsertAndGetUserRole,
		user.Name,
		user.Email,
		user.Username,
		user.Address,
		string(hashedPassword),
		user.Role,
		user.UpdatedAt).Scan(
		&createdUser.Id,
		&createdUser.CreatedAt,
	)
	if err != nil {
		log.Println("userRepository.QueryRow:", err.Error())
		return entity.User{}, err
	}

	createdUser.Name = user.Name
	createdUser.Email = user.Email
	createdUser.Username = user.Username
	createdUser.Address = user.Address
	createdUser.Hashpassword = string(hashedPassword)
	createdUser.Role = user.Role

	err = t.addParticipant(createdUser)
	if err != nil {
		log.Println("userRepository.getCreatedUser: Error adding Participant:", err.Error())
		return entity.User{}, err
	}

	return createdUser, nil
}

func (t *userRepository) addParticipant(user entity.User) error {
	//  Add user to participant table
	_, err := t.db.Exec(config.InsertUserToParticipant,
		user.Id,
		user.Id,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		log.Println("userRepository.addPArticipant: Error adding Participant:", err.Error())
		return err
	}

	return nil
}

// Created implements UserRepository.
func (t *userRepository) Created(data entity.User) (entity.User, error) {
	var user entity.User

	// Hash password menggunakan bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Hashpassword), bcrypt.DefaultCost)
	if err != nil {
		log.Println("userRepository.GenerateFromPassword:", err.Error())
		return entity.User{}, err
	}

	err = t.db.QueryRow(config.InsertUser,
		data.Name,
		data.Email,
		data.Username,
		data.Address,
		string(hashedPassword),
		data.Role,
		data.UpdatedAt).Scan(
		&user.Id,
		&user.CreatedAt,
	)
	if err != nil {
		log.Println("userRepository.QueryRow:", err.Error())
		return entity.User{}, err
	}

	user.Name = data.Name
	user.Email = data.Email
	user.Username = data.Username
	user.Address = data.Address
	user.Hashpassword = string(hashedPassword)
	user.Role = data.Role

	return user, nil

}

// List implements UserRepository.
func (t *userRepository) List(page, size int) ([]entity.User, model.Paging, error) {
	var users []entity.User
	offset := (page - 1) * size
	row, err := t.db.Query(config.ListUsers, size, offset)

	if err != nil {
		log.Println("UserRepository.Query:", err.Error())
		return nil, model.Paging{}, err
	}

	for row.Next() {
		var user entity.User
		err := row.Scan(
			&user.Id,
			&user.Name,
			&user.Email,
			&user.Username,
			&user.Address,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			log.Println("UserRepository.Rows.Next():", err.Error())
			return nil, model.Paging{}, err
		}

		users = append(users, user)
	}
	totalRows := 0
	if err := t.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&totalRows); err != nil {
		return nil, model.Paging{}, err
	}

	paging := model.Paging{
		Page:        page,
		RowsPerPage: size,
		TotalRows:   totalRows,
		TotalPages:  int(math.Ceil(float64(totalRows) / float64(size))),
	}
	return users, paging, nil

}

// Delete implements UserRepository.
func (t *userRepository) Delete(id string) (entity.User, error) {
	result, err := t.db.Exec(config.DeleteUserByID, id)
	if err != nil {
		log.Println("CustomerRepository.QueryRow:", err.Error())
		return entity.User{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error getting rows affected:", err.Error())
		return entity.User{}, err
	}

	fmt.Printf("Rows affected: %d\n", rowsAffected)

	return entity.User{}, nil
}

// Get implements UserRepository.
func (t *userRepository) Get(id string) (entity.User, error) {
	var user entity.User
	err := t.db.QueryRow(config.GetUserByID, id).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Username,
		&user.Address,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		log.Println("UserRepository.QueryRow:", err.Error())
		return entity.User{}, err
	}
	return user, nil

}

// Updated implements UserRepository.
func (t *userRepository) Updated(id string, data entity.User) (entity.User, error) {
	if data.Name == "" && data.Email == "" && data.Username == "" && data.Address == "" && data.Hashpassword == "" && data.Role == "" {
		return data, nil // Tidak ada pembaruan yang diperlukan jika semua field kosong
	}

	// Pilih jenis pembaruan berdasarkan field yang tidak kosong
	if data.Email != "" {
		_, err := t.db.Exec(config.UpdatedUserByEmail, id, data.Email)
		if err != nil {
			log.Println("UserRepository.Exec:", err.Error())
			return entity.User{}, err
		}
		return data, nil
	} else if data.Username != "" {
		_, err := t.db.Exec(config.UpdatedUserByUsername, id, data.Username)
		if err != nil {
			log.Println("UserRepository.Exec:", err.Error())
			return entity.User{}, err
		}
		return data, nil
	} else if data.Address != "" {
		_, err := t.db.Exec(config.UpdatedUserByAddress, id, data.Address)
		if err != nil {
			log.Println("UserRepository.Exec:", err.Error())
			return entity.User{}, err
		}
		return data, nil
	} else if data.Hashpassword != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Hashpassword), bcrypt.DefaultCost)
		if err != nil {
			log.Println("bcrypt.GenerateFromPassword:", err.Error())
			return entity.User{}, err
		}

		_, err = t.db.Exec(config.UpdatedUserByHashPassword, id, string(hashedPassword))
		if err != nil {
			log.Println("UserRepository.Exec:", err.Error())
			return entity.User{}, err
		}

		data.Hashpassword = string(hashedPassword)
		return data, nil
	} else if data.Role != "" {
		_, err := t.db.Exec(config.UpdatedUserByRole, id, data.Role)
		if err != nil {
			log.Println("UserRepository.Exec:", err.Error())
			return entity.User{}, err
		}
		return data, nil
	}

	// Jika tidak ada field tertentu yang diisi, lakukan pembaruan keseluruhan
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Hashpassword), bcrypt.DefaultCost)
	if err != nil {
		log.Println("bcrypt.GenerateFromPassword:", err.Error())
		return entity.User{}, err
	}

	_, err = t.db.Exec(config.UpdatedUserAll, id, data.Name, data.Email, data.Username, data.Address, string(hashedPassword))
	if err != nil {
		log.Println("UserRepository.Exec:", err.Error())
		return entity.User{}, err
	}

	data.Hashpassword = string(hashedPassword)
	return data, nil
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}
