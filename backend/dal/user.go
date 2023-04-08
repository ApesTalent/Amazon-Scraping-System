package dal

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// PasswordPayload for password change payload
type PasswordPayload struct {
	OldPass string `json:"old_pass"`
	NewPass string `json:"new_pass"`
}

// CreateUser creates new user
func CreateUser(u User) error {
	c := db.Collection(ColUser)
	ctx := context.Background()

	hp, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
	u.Password = string(hp)
	u.Role = "member"
	_, err := c.InsertOne(ctx, u)
	fmt.Println("user create", err)
	return err
}

// Login login
func Login(u User) (User, error) {
	var user User
	c := db.Collection(ColUser)
	ctx := context.Background()

	err := c.FindOne(ctx, bson.M{"username": u.Username}).Decode(&user)
	if err != nil {
		return user, err
	}
	compare := comparePasswords(user.Password, u.Password)
	if compare != true {
		return user, errors.New("Wrong Password")
	}
	return user, nil
}

// GetUserByID get user by id
func GetUserByID(id primitive.ObjectID) (User, error) {
	var user User
	c := db.Collection(ColUser)
	ctx := context.Background()

	err := c.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return user, err
	}
	return deletePwd(user), nil
}

// UpdateRole updates user's role
func UpdateRole(u User) error {
	c := db.Collection(ColUser)
	ctx := context.Background()

	if u.Password != "prime_admin_pwd" {
		return errors.New("permission denied")
	}
	_, err := c.UpdateOne(ctx, bson.M{"_id": u.ID}, bson.M{"$set": bson.M{
		"role": u.Role,
	}})
	return err
}

// ChangePassword updates user's password
func ChangePassword(p PasswordPayload, id primitive.ObjectID) error {
	var user User
	c := db.Collection(ColUser)
	ctx := context.Background()

	err := c.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return err
	}

	compare := comparePasswords(user.Password, p.OldPass)
	if compare != true {
		return errors.New("Wrong Password")
	}
	hp, _ := bcrypt.GenerateFromPassword([]byte(p.NewPass), bcrypt.MinCost)
	user.Password = string(hp)

	_, err = c.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{
		"password": string(hp),
	}})
	return err
}

func deletePwd(u User) User {
	u.Password = ""
	return u
}

func comparePasswords(hashedPwd, plainPwd string) bool {
	byteHash := []byte(hashedPwd)
	bytePlain := []byte(plainPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePlain)
	if err != nil {
		return false
	}
	return true
}
