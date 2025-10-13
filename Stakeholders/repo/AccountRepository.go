package repo

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"stakeholders.com/dto"
	"stakeholders.com/model"
)

type AccountRepository struct {
	DatabaseConnection *gorm.DB
}

func (repo *AccountRepository) FindAll() ([]model.Account, error) {
	var accounts []model.Account
	dbResult := repo.DatabaseConnection.Select("id", "username", "email", "role", "blocked").Find(&accounts)

	if dbResult.Error != nil {
		return nil, dbResult.Error
	}
	return accounts, nil
}

func (repo *AccountRepository) Create(account *model.Account) error {
	if account.Id == "" {
		account.Id = uuid.New().String()
	}
	dbResult := repo.DatabaseConnection.Create(account)

	if dbResult.Error != nil {
		return dbResult.Error
	}
	fmt.Println("Rows affected: ", dbResult.RowsAffected)
	fmt.Printf("Created account: %+v\n", account)
	return nil
}

func (repo *AccountRepository) FindAccount(accountId string) error {
	var account model.Account
	dbResult := repo.DatabaseConnection.First(&account, "id = ?", accountId)
	if dbResult.Error != nil {
		return dbResult.Error
	}
	return nil
}

func (repo *AccountRepository) FindById(accountId string) (model.Account, error) {
	var account model.Account
	dbResult := repo.DatabaseConnection.First(&account, "id = ?", accountId)

	if dbResult.Error != nil {
		return account, dbResult.Error
	}

	return account, nil
}

func (repo *AccountRepository) Login(loginDetails *dto.LoginDetailsDto) (string, error) {
	var account model.Account
	dbResult := repo.DatabaseConnection.First(&account, "email = ?", loginDetails.Email)
	if dbResult.Error != nil {
		return "Nema", dbResult.Error
	}
	if account.Password != loginDetails.Password {
		return "Lozinka", nil
	}
	if account.Blocked {
		return "Blok", nil
	}
	token, err := CreateToken(&account)

	if err != nil {
		return "Token", err
	}
	return token, nil
}

func CreateToken(account *model.Account) (string, error) {
	secretKey := []byte("sekret_key_12#4")

	claims := &dto.Claims{
		Role:     account.Role,
		Username: account.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "stakeholders.com",
			Subject:   account.Id,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (repo *AccountRepository) BlockAccount(accountId string) error {
	return repo.DatabaseConnection.Model(&model.Account{}).
		Where("id = ?", accountId).
		Update("blocked", true).Error
}

func (repo *AccountRepository) GetUsernameById(accountId string) (string, error) {
	var username string
	dbResult := repo.DatabaseConnection.Table("accounts").
		Select("username").
		Where("id = ?", accountId).
		Scan(&username)

	if dbResult.Error != nil {
		return "", dbResult.Error
	}
	if username == "" {
		return "", nil
	}
	return username, nil
}

func (repo *AccountRepository) ParseToken(token string) *dto.Claims {
	secretKey := []byte("sekret_key_12#4")

	parsedToken, err := jwt.ParseWithClaims(token, &dto.Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Proveri da li je algoritam ispravan
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("nevalidan signing metod: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		log.Fatalf("Greška pri parsiranju tokena: %v", err)
		return nil
	}

	var customClaims *dto.Claims
	if claims, ok := parsedToken.Claims.(*dto.Claims); ok && parsedToken.Valid {
		customClaims = claims
	}
	return customClaims
}
