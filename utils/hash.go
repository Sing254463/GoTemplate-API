package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword ฟังก์ชันสำหรับเข้ารหัสรหัสผ่านด้วย bcrypt
// รับรหัสผ่านธรรมดาและคืนค่ารหัสผ่านที่เข้ารหัสแล้ว
// bcrypt เป็นอัลกอริธึมที่ปลอดภัยและแนะนำสำหรับการเข้ารหัสรหัสผ่าน
func HashPassword(password string) (string, error) {
	// ใช้ bcrypt.DefaultCost (ระดับ 10) ในการเข้ารหัส
	// ยิ่งค่า cost สูง ยิ่งใช้เวลานานแต่ปลอดภัยขึ้น
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err // คืนค่าข้อผิดพลาดหากเข้ารหัสไม่สำเร็จ
	}
	return string(hashedPassword), nil // คืนค่ารหัสผ่านที่เข้ารหัสแล้ว
}

// CheckPassword ฟังก์ชันสำหรับตรวจสอบรหัสผ่าน
// เปรียบเทียบรหัสผ่านที่เข้ารหัสแล้วกับรหัสผ่านธรรมดา
// คืนค่า nil หากรหัสผ่านถูกต้อง, คืนค่า error หากไม่ถูกต้อง
func CheckPassword(hashedPassword, password string) error {
	// ใช้ bcrypt.CompareHashAndPassword เพื่อเปรียบเทียบ
	// ฟังก์ชันนี้จะจัดการการเปรียบเทียบอย่างปลอดภัย
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
