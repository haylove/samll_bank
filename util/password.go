//  *@createTime    2022/3/21 5:06
//  *@author        hay&object
//  *@version       v1.0.0

package util

import "golang.org/x/crypto/bcrypt"

//HashPassword returns the bcrypt hash of password
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), err
}

//CheckPassword checks if the provided password is correct or not
func CheckPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
