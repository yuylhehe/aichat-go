package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

// HashPassword 密码哈希
func HashPassword(password string, salt string) string {
	hash := sha256.Sum256([]byte(password + salt))
	return hex.EncodeToString(hash[:])
}

// GenerateSalt 生成盐值
func GenerateSalt() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// VerifyPassword 验证密码
func VerifyPassword(password, salt, hashedPassword string) bool {
	return HashPassword(password, salt) == hashedPassword
}
