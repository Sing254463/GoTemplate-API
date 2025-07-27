// filepath: utils/response.go
package utils

import (
	"github.com/gofiber/fiber/v2"
)

// Response โครงสร้างสำหรับส่งผลลัพธ์กลับไปยัง client
// ใช้เป็นรูปแบบมาตรฐานสำหรับทุก API response
type Response struct {
	Status  bool        `json:"status"`          // สถานะความสำเร็จ (true/false)
	Message string      `json:"message"`         // ข้อความอธิบาย
	Data    interface{} `json:"data,omitempty"`  // ข้อมูล (จะแสดงเมื่อสำเร็จ)
	Error   string      `json:"error,omitempty"` // ข้อผิดพลาด (จะแสดงเมื่อมีข้อผิดพลาด)
}

// SuccessResponse ฟังก์ชันสำหรับส่ง response เมื่อสำเร็จ
// รับพารามิเตอร์: context, ข้อความ, และข้อมูล
func SuccessResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  true,    // กำหนดสถานะเป็น true (สำเร็จ)
		Message: message, // ข้อความที่ต้องการส่ง
		Data:    data,    // ข้อมูลที่ต้องการส่งกลับ
	})
}

// ErrorResponse ฟังก์ชันสำหรับส่ง response เมื่อเกิดข้อผิดพลาด
// รับพารามิเตอร์: context, รหัสสถานะ HTTP, ข้อความ, และ error object
func ErrorResponse(c *fiber.Ctx, statusCode int, message string, err error) error {
	// สร้าง response object พื้นฐาน
	response := Response{
		Status:  false,   // กำหนดสถานะเป็น false (ไม่สำเร็จ)
		Message: message, // ข้อความอธิบายข้อผิดพลาด
	}

	// หากมี error object ให้เพิ่มรายละเอียดข้อผิดพลาด
	if err != nil {
		response.Error = err.Error()
	}

	// ส่ง response พร้อมรหัสสถานะ HTTP ที่กำหนด
	return c.Status(statusCode).JSON(response)
}

// CreatedResponse ฟังก์ชันสำหรับส่ง response เมื่อสร้างข้อมูลใหม่สำเร็จ
// ใช้สำหรับกรณี POST request ที่สร้างข้อมูลใหม่ (HTTP 201)
func CreatedResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(Response{
		Status:  true,    // กำหนดสถานะเป็น true (สำเร็จ)
		Message: message, // ข้อความยืนยันการสร้าง
		Data:    data,    // ข้อมูลที่สร้างใหม่
	})
}
