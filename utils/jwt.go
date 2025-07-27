package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims โครงสร้างสำหรับเก็บข้อมูลใน JWT token
// ประกอบด้วยข้อมูลผู้ใช้และ claims มาตรฐานของ JWT
type JWTClaims struct {
	UserID               int    `json:"user_id"`  // ID ของผู้ใช้
	Username             string `json:"username"` // ชื่อผู้ใช้
	Role                 string `json:"role"`     // สิทธิ์ของผู้ใช้ (user/admin)
	jwt.RegisteredClaims        // Claims มาตรฐาน (เวลาหมดอายุ, เวลาออก, ฯลฯ)
}

// GenerateJWT ฟังก์ชันสำหรับสร้าง JWT token
// รับข้อมูลผู้ใช้และการตั้งค่า แล้วสร้าง token ที่เซ็นแล้ว
func GenerateJWT(userID int, username, role, secret string, expire time.Duration) (string, error) {
	// สร้าง claims ที่มีข้อมูลผู้ใช้และเวลาหมดอายุ
	claims := &JWTClaims{
		UserID:   userID,   // ID ของผู้ใช้
		Username: username, // ชื่อผู้ใช้
		Role:     role,     // สิทธิ์ของผู้ใช้
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)), // เวลาหมดอายุ
			IssuedAt:  jwt.NewNumericDate(time.Now()),             // เวลาที่ออก token
		},
	}

	// สร้าง token ด้วย claims และ signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// เซ็น token ด้วยกุญแจลับและคืนค่า token string
	return token.SignedString([]byte(secret))
}

// ParseJWT ฟังก์ชันสำหรับตรวจสอบและแยกข้อมูลจาก JWT token
// รับ token string และกุญแจลับ แล้วคืนค่า claims หากถูกต้อง
func ParseJWT(tokenString, secret string) (*JWTClaims, error) {
	// แยกและตรวจสอบ token ด้วยกุญแจลับ
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil // คืนค่ากุญแจลับสำหรับการตรวจสอบ
	})

	// หากเกิดข้อผิดพลาดในการแยก token
	if err != nil {
		return nil, err
	}

	// ตรวจสอบว่า claims ถูกต้องและ token ยังใช้งานได้
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil // คืนค่า claims หากถูกต้อง
	}

	// หาก token ไม่ถูกต้อง
	return nil, jwt.ErrInvalidKey
}
