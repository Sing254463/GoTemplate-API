package middleware

import (
	"strings"

	"github.com/Sing254463/GoTemplate/Backend/utils"
	"github.com/gofiber/fiber/v2"
)

// JWTMiddleware ฟังก์ชันสร้าง middleware สำหรับตรวจสอบ JWT token
// รับพารามิเตอร์ secret (กุญแจลับ) และคืนค่า fiber.Handler
func JWTMiddleware(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// ดึง token จาก Authorization header
		// รูปแบบที่คาดหวัง: "Bearer <token>"
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "ไม่พบ authorization header", nil)
		}

		// แยก header เป็นส่วนๆ และตรวจสอบรูปแบบ
		// ต้องเป็น "Bearer <token>" เท่านั้น
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "รูปแบบ authorization header ไม่ถูกต้อง", nil)
		}

		// ดึง token string จากส่วนที่ 2
		tokenString := parts[1]

		// ตรวจสอบและแยกข้อมูลจาก JWT token
		claims, err := utils.ParseJWT(tokenString, secret)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Token ไม่ถูกต้องหรือหมดอายุ", err)
		}

		// เก็บข้อมูลผู้ใช้ใน context เพื่อให้ handler ต่อไปใช้งานได้
		// ข้อมูลที่เก็บ: user_id, username, role
		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("role", claims.Role)

		// ส่งต่อไปยัง handler ถัดไป
		return c.Next()
	}
}

// AdminMiddleware ฟังก์ชันสร้าง middleware สำหรับตรวจสอบสิทธิ์ Admin
// ใช้ร่วมกับ JWTMiddleware เพื่อให้มั่นใจว่าผู้ใช้เป็น Admin
func AdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// ดึงข้อมูล role จาก context (ที่เก็บไว้โดย JWTMiddleware)
		role, ok := c.Locals("role").(string)

		// ตรวจสอบว่า role เป็น "admin" หรือไม่
		if !ok || role != "admin" {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "ต้องมีสิทธิ์ Admin เท่านั้น", nil)
		}

		// หากเป็น Admin ให้ส่งต่อไปยัง handler ถัดไป
		return c.Next()
	}
}
