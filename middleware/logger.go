package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Logger() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\n",
		TimeFormat: time.RFC3339,
		TimeZone:   "Local",
	})
}

// Logger เป็น middleware สำหรับบันทึก log ของการร้องขอ HTTP
// โดยจะบันทึกข้อมูลต่าง ๆ เช่น เวลา, สถานะ, ระยะเวลาในการประมวลผล, IP ของผู้ใช้, วิธีการร้องขอ (GET, POST, ฯลฯ), เส้นทางที่ร้องขอ และข้อผิดพลาด (ถ้ามี)
// การใช้งาน Logger จะช่วยให้เราสามารถติดตามและวิเคราะห์การทำงานของแอปพลิเคชันได้ง่ายขึ้น
// โดยสามารถนำไปใช้ใน main.go ได้ดังนี้:
// app.Use(middleware.Logger())
// ซึ่งจะทำให้ทุกครั้งที่มีการร้องขอ HTTP จะมีการบันทึก log ตามรูปแบบที่กำหนดไว้ใน logger.Config
// โดยสามารถปรับแต่งรูปแบบการบันทึกได้ตามต้องการ เช่น การเปลี่ยนรูปแบบเวลา, การเพิ่มข้อมูลอื่น ๆ ที่ต้องการบันทึก, หรือการเปลี่ยนแปลงโซนเวลา
// นอกจากนี้ยังสามารถใช้ middleware นี้ร่วมกับ middleware อื่น ๆ ได้ เช่น การตรวจสอบ JWT, การจัดการ CORS, หรือการจัดการ session
// เพื่อให้การจัดการการร้องขอ HTTP มีความยืดหยุ่นและมีประสิทธิภาพมากยิ่งขึ้น
// นอกจากนี้ยังสามารถปรับแต่งการบันทึก log ให้เหมาะสมกับความต้องการของแอปพลิเคชัน เช่น การบันทึก log ในรูปแบบ JSON หรือการส่ง log ไปยังระบบจัดการ log อื่น ๆ ได้
// การใช้ middleware.Logger() จะช่วยให้เราสามารถติดตามและวิเคราะห์การทำงานของแอปพลิเคชันได้อย่างมีประสิทธิภาพ
// นอกจากนี้ยังสามารถปรับแต่งการบันทึก log ให้เหมาะสมกับความต้องการของแอปพลิเคชัน เช่น การบันทึก log ในรูปแบบ JSON หรือการส่ง log ไปยังระบบจัดการ log อื่น ๆ ได้
