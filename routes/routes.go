package routes

import (
	"runtime"
	"time"

	"github.com/Sing254463/GoTemplate/Backend/config"
	"github.com/Sing254463/GoTemplate/Backend/controllers"
	"github.com/Sing254463/GoTemplate/Backend/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// SetupRoutes ฟังก์ชันสำหรับตั้งค่าเส้นทาง (routes) ทั้งหมดของ API
// รับพารามิเตอร์ app (Fiber app) และ cfg (configuration)
func SetupRoutes(app *fiber.App, cfg *config.Config) {
	// สร้างและเตรียมคอนโทรลเลอร์สำหรับจัดการคำร้องขอ
	// authController จัดการเรื่องการลงทะเบียน, เข้าสู่ระบบ, และโปรไฟล์
	authController := controllers.NewAuthController(cfg)
	// userController จัดการเรื่องข้อมูลผู้ใช้ (สำหรับ admin เท่านั้น)
	userController := controllers.NewUserController(cfg)

	// ตั้งค่าเส้นทางสำหรับ Swagger UI (เอกสาร API)
	// เส้นทาง /swagger แสดงหน้า Swagger UI หลัก
	app.Get("/swagger", swagger.New(swagger.Config{
		URL:         "/swagger/doc.json", // ที่อยู่ของไฟล์ JSON ที่มีข้อมูล API
		DeepLinking: false,               // ปิดการใช้งาน deep linking
	}))
	// เส้นทาง /swagger/* สำหรับไฟล์ static ของ Swagger
	app.Get("/swagger/*", swagger.HandlerDefault)

	// สร้างกลุ่มเส้นทางหลักสำหรับ API version 1
	// ทุกเส้นทาง API จะเริ่มต้นด้วย /api/v1
	api := app.Group("/api/v1")

	// เส้นทางตรวจสอบสถานะเซิร์ฟเวอร์ (Health Check)
	// ใช้เพื่อตรวจสอบว่าเซิร์ฟเวอร์ทำงานปกติหรือไม่
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":      "ok",
			"message":     "GoTemplate API is running!",
			"app_name":    cfg.App.Name,
			"version":     cfg.App.Version,
			"description": cfg.App.Description,
			"environment": cfg.Server.Environment,
			"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
		})
	})

	// เส้นทางสำหรับแสดงข้อมูลเวอร์ชัน
	api.Get("/version", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"app_name":    cfg.App.Name,
			"version":     cfg.App.Version,
			"description": cfg.App.Description,
			"go_version":  runtime.Version(),
			"build_time":  time.Now().Format("2006-01-02 15:04:05"),
		})
	})

	// กลุ่มเส้นทางสำหรับการจัดการการยืนยันตัวตน (Authentication)
	// เส้นทางเหล่านี้เปิดให้สาธารณะเข้าถึงได้ (ไม่ต้องเข้าสู่ระบบ)
	auth := api.Group("/auth")
	auth.Post("/register", authController.Register) // ลงทะเบียนผู้ใช้ใหม่
	auth.Post("/login", authController.Login)       // เข้าสู่ระบบ

	// กลุ่มเส้นทางที่ต้องมีการยืนยันตัวตน (Protected Routes)
	// ต้องส่ง JWT Token ใน Authorization header จึงจะเข้าถึงได้
	protected := api.Group("")
	protected.Use(middleware.JWTMiddleware(cfg.JWT.Secret)) // ใช้ middleware ตรวจสอบ JWT

	// เส้นทางที่ต้องเข้าสู่ระบบสำหรับข้อมูลส่วนตัว
	authProtected := protected.Group("/auth")
	authProtected.Get("/profile", authController.GetProfile) // ดูข้อมูลโปรไฟล์ตนเอง

	// กลุ่มเส้นทางสำหรับจัดการผู้ใช้ (User Management)
	// ต้องเป็น Admin เท่านั้นถึงจะเข้าถึงได้
	users := protected.Group("/users")
	users.Use(middleware.AdminMiddleware())         // ใช้ middleware ตรวจสอบสิทธิ์ Admin
	users.Get("/", userController.GetAllUsers)      // ดูรายชื่อผู้ใช้ทั้งหมด
	users.Get("/:id", userController.GetUserByID)   // ดูข้อมูลผู้ใช้ตาม ID
	users.Delete("/:id", userController.DeleteUser) // ลบผู้ใช้ตาม ID
}
