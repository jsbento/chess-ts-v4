package behavior

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
	t "github.com/jsbento/chess-server-v4/cmd/services/users/types"
)

func (s *Store) CreateUser(req *t.CreateUser) (*t.User, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	user := &t.User{
		Id:       uuid.New().String(),
		Username: req.Username,
		Email:    req.Email,
		Password: hashPassword(req.Password),
	}
	if err := s.pg.GetDB().Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Store) UpdateUser(req *t.UpdateUser) (*t.User, error) {
	updatedUser := &t.User{}
	txn := s.pg.GetDB().Model(&t.User{}).Where("id = ?", req.Id).Omit("Password")
	if req.Username == nil {
		txn = txn.Omit("Username")
	} else {
		updatedUser.Username = *req.Username
	}
	if req.Email == nil {
		txn = txn.Omit("Email")
	} else {
		updatedUser.Email = *req.Email
	}

	if err := txn.Updates(&updatedUser).Error; err != nil {
		return nil, err
	}
	return updatedUser, nil
}

func (s *Store) GetUser(id string) (*t.User, error) {
	user := &t.User{}
	if err := s.pg.GetDB().Model(&t.User{}).Where("id = ?", id).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func hashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hashedPassword)
}
