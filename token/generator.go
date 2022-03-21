//  *@createTime    2022/3/21 16:04
//  *@author        hay&object
//  *@version       v1.0.0

package token

import "time"

// Generator generate token and verify
type Generator interface {
	// GenerateToken generates a token with given username and expiredTime
	GenerateToken(username string, duration time.Duration) (string, error)

	//VerifyToken verify  if the token is valid or not
	VerifyToken(token string) (*Payload, error)
}
