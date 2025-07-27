package utils

import (
	"fmt"
	"math/rand"
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

// generateJTI สร้าง JWT ID (JTI) ที่ unique สำหรับแต่ละ token
// JTI ทำให้ token แต่ละครั้งที่ login จะไม่เหมือนกัน แม้ข้อมูลจะเหมือนเดิม
// รูปแบบ: timestamp-nanosecond-random เช่น 1737148167123456789-5847
// วิธีนี้ช่วยให้:
// 1. Token unique ทุกครั้งที่ login
// 2. สามารถติดตาม session แต่ละครั้งได้
// 3. เพิ่มความปลอดภัยโดยไม่ต้องเก็บข้อมูลใน database
func generateJTI() string {
	// สร้าง JTI ที่ประกอบด้วย timestamp และ random number
	// format: timestamp-nanosecond-random
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), rand.Intn(10000))
}

// GenerateJWT ฟังก์ชันสำหรับสร้าง JWT token พร้อม JTI
// รับข้อมูลผู้ใช้และการตั้งค่า แล้วสร้าง token ที่เซ็นแล้ว
//
// Parameters:
// - userID: ID ของผู้ใช้ในฐานข้อมูล
// - username: ชื่อผู้ใช้
// - role: บทบาท (user/admin)
// - secret: กุญแจลับสำหรับเซ็น token
// - expire: ระยะเวลาหมดอายุ (เช่น "24h", "7d")
//
// Returns:
// - string: JWT token ที่เซ็นแล้ว
// - error: ข้อผิดพลาด (ถ้ามี)
func GenerateJWT(userID int, username, role, secret string, expire time.Duration) (string, error) {

	// สร้าง claims ที่มีข้อมูลผู้ใช้และเวลาหมดอายุ
	// JTI (JWT ID) ทำให้ token unique ทุกครั้งที่ login ใหม่
	claims := jwt.MapClaims{
		"user_id":  userID,                        // ID ผู้ใช้ในฐานข้อมูล
		"username": username,                      // ชื่อผู้ใช้
		"role":     role,                          // บทบาท (user/admin)
		"jti":      generateJTI(),                 // JWT ID - ทำให้ token unique
		"exp":      time.Now().Add(expire).Unix(), // เวลาหมดอายุ (Unix timestamp)
		"iat":      time.Now().Unix(),             // เวลาที่สร้าง token (Issued At)
		"iss":      "GoTemplate",                  // ผู้ออก token (Issuer)
		"aud":      "GoTemplate-Users",            // ผู้รับ token (Audience)
	}

	// สร้าง token ด้วย claims และ signing method (HMAC SHA-256)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// เซ็น token ด้วยกุญแจลับและคืนค่า token string
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("ไม่สามารถเซ็น token ได้: %w", err)
	}

	return tokenString, nil
}

// ParseJWT ฟังก์ชันสำหรับตรวจสอบและแยกข้อมูลจาก JWT token
// รับ token string และกุญแจลับ แล้วคืนค่า claims หากถูกต้อง

// Parameters:
// - tokenString: JWT token ที่ต้องการตรวจสอบ
// - secret: กุญแจลับสำหรับตรวจสอบลายเซ็น
//
// Returns:
// - *JWTClaims: ข้อมูล claims ถ้า token ถูกต้อง
// - error: ข้อผิดพลาดถ้า token ไม่ถูกต้องหรือหมดอายุ
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

// ValidateJWT ฟังก์ชันสำหรับตรวจสอบ JWT token และคืนค่า MapClaims
// ใช้สำหรับ middleware ที่ต้องการเข้าถึงข้อมูลใน token แบบ flexible
//
// Parameters:
// - tokenString: JWT token ที่ต้องการตรวจสอบ
// - secret: กุญแจลับสำหรับตรวจสอบลายเซ็น
//
// Returns:
// - jwt.MapClaims: ข้อมูล claims ในรูปแบบ map
// - error: ข้อผิดพลาดถ้า token ไม่ถูกต้อง
func ValidateJWT(tokenString, secret string) (jwt.MapClaims, error) {
	// แปลง token string เป็น JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// ตรวจสอบ algorithm ว่าเป็น HMAC หรือไม่
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// ตรวจสอบว่า token ยังใช้ได้หรือไม่
	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	// แปลง claims เป็น MapClaims
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}
