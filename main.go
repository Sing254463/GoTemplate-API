package main

import (
	"fmt"
	"log"

	"github.com/Sing254463/GoTemplate/Backend/config"
	_ "github.com/Sing254463/GoTemplate/Backend/docs"
	"github.com/Sing254463/GoTemplate/Backend/middleware"
	"github.com/Sing254463/GoTemplate/Backend/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// ============================================
// Swagger API Documentation Annotations
// ============================================
// คำอธิบาย API สำหรับสร้างเอกสาร Swagger อัตโนมัติ

// @title GoTemplate API
// ชื่อของ API ที่จะแสดงใน Swagger UI

// @version 1.0.0
// เวอร์ชันของ API

// @description A complete Go template API with authentication and user management
// คำอธิบายของ API - เป็น API สำหรับจัดการผู้ใช้และการยืนยันตัวตน

// @termsOfService http://swagger.io/terms/
// ลิงก์ไปยังข้อกำหนดการใช้งาน

// @contact.name API Support
// ชื่อของทีมสนับสนุน API

// @contact.email support@example.com
// อีเมลติดต่อสำหรับการสนับสนุน

// @license.name MIT
// ใบอนุญาตการใช้งาน

// @license.url https://opensource.org/licenses/MIT
// ลิงก์ไปยังใบอนุญาต MIT

// @host localhost:8080
// ที่อยู่และพอร์ตของเซิร์ฟเวอร์

// @BasePath /api/v1
// เส้นทางหลักของ API

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Enter: Bearer {token}
// การตั้งค่าความปลอดภัยสำหรับ JWT Token

// ============================================
// ฟังก์ชันหลักของแอปพลิเคชัน
// ============================================
func main() {
	// ============================================
	// 1. โหลดการตั้งค่าจากไฟล์ .env
	// ============================================
	// cfg := config.LoadConfig() ทำหน้าที่:
	// - อ่านค่าตัวแปรจากไฟล์ .env (DB_HOST, DB_PORT, JWT_SECRET ฯลฯ)
	// - สร้างการเชื่อมต่อกับฐานข้อมูล MySQL
	// - ตั้งค่า JWT และ Server configuration
	// - ตรวจสอบการเชื่อมต่อฐานข้อมูลและแสดงผลลัพธ์
	fmt.Println("🔧 กำลังโหลดการตั้งค่าระบบ...")
	cfg := config.LoadConfig()
	fmt.Println("✅ โหลดการตั้งค่าสำเร็จ")

	// ============================================
	// 2. สร้างแอปพลิเคชัน Fiber
	// ============================================
	// Fiber เป็น web framework ที่รวดเร็วสำหรับ Go (คล้าย Express.js)
	// การสร้าง app พร้อมกับการตั้งค่า ErrorHandler สำหรับจัดการข้อผิดพลาด
	fmt.Println("🚀 กำลังสร้างแอปพลิเคชัน Fiber...")

	app := fiber.New(fiber.Config{
		// ErrorHandler: ฟังก์ชันสำหรับจัดการข้อผิดพลาดที่เกิดขึ้นในแอปพลิเคชัน
		// จะถูกเรียกใช้เมื่อมี error เกิดขึ้นในระหว่างการประมวลผล request
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// กำหนดรหัสสถานะเริ่มต้นเป็น 500 (Internal Server Error)
			code := fiber.StatusInternalServerError

			// ตรวจสอบว่า error เป็นประเภท fiber.Error หรือไม่
			// หากใช่ จะใช้รหัสสถานะที่กำหนดไว้ใน error นั้น
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			// ส่งกลับ JSON response ที่มีข้อมูลข้อผิดพลาด
			return c.Status(code).JSON(fiber.Map{
				"status":  false,                   // สถานะการทำงาน (ไม่สำเร็จ)
				"message": "Internal Server Error", // ข้อความแสดงข้อผิดพลาด
				"error":   err.Error(),             // รายละเอียดข้อผิดพลาด
			})
		},
	})
	fmt.Println("✅ สร้างแอปพลิเคชัน Fiber สำเร็จ")

	// ============================================
	// 3. การติดตั้ง Middleware (ตัวกลางประมวลผล)
	// ============================================
	// Middleware คือฟังก์ชันที่ทำงานระหว่าง request และ response
	// จะทำงานตามลำดับที่กำหนดไว้

	fmt.Println("🔌 กำลังติดตั้ง Middleware...")

	// 3.1 Recover Middleware - กู้คืนจาก panic
	// หากเกิด panic ในแอปพลิเคชัน จะไม่ให้เซิร์ฟเวอร์หยุดทำงาน
	// แต่จะแสดงข้อผิดพลาดและทำงานต่อไป
	app.Use(recover.New())
	fmt.Println("   ✅ ติดตั้ง Recover Middleware (กู้คืนจาก panic)")

	// 3.2 Logger Middleware - บันทึกข้อมูลการร้องขอ
	// จะบันทึกข้อมูลการร้องขอทุกครั้ง เช่น:
	// - เวลาที่ใช้ในการประมวลผล
	// - รหัสสถานะของการตอบสนอง (200, 404, 500 ฯลฯ)
	// - URL ที่ร้องขอ
	// - HTTP method (GET, POST, PUT, DELETE)
	app.Use(middleware.Logger())
	fmt.Println("   ✅ ติดตั้ง Logger Middleware (บันทึกข้อมูลการร้องขอ)")

	// 3.3 CORS Middleware - จัดการ Cross-Origin Resource Sharing
	// อนุญาตให้เว็บไซต์จากโดเมนอื่นสามารถเรียกใช้ API ได้
	// เช่น หากมี Frontend ที่รันบนพอร์ต 3000 ต้องการเรียก API บนพอร์ต 8080
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // อนุญาตให้ทุกโดเมนเข้าถึง (ใช้สำหรับการพัฒนา)
		// ในการใช้งานจริงควรระบุโดเมนที่แน่นอน เช่น "https://mywebsite.com"

		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		// HTTP methods ที่อนุญาตให้ใช้:
		// - GET: ดึงข้อมูล
		// - POST: สร้างข้อมูลใหม่
		// - PUT: อัปเดตข้อมูล
		// - DELETE: ลบข้อมูล
		// - PATCH: อัปเดตข้อมูลบางส่วน
		// - OPTIONS: ตรวจสอบสิทธิ์ CORS
		// - HEAD: ดึงข้อมูล header เท่านั้น

		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		// Headers ที่อนุญาตให้ส่งมาด้วย:
		// - Origin: ที่มาของการร้องขอ
		// - Content-Type: ประเภทข้อมูลที่ส่ง (เช่น application/json)
		// - Accept: ประเภทข้อมูลที่ต้องการรับ
		// - Authorization: JWT Token สำหรับการยืนยันตัวตน
	}))
	fmt.Println("   ✅ ติดตั้ง CORS Middleware (จัดการ Cross-Origin)")

	// ============================================
	// 4. ตั้งค่าเส้นทาง API (Routes)
	// ============================================
	// routes.SetupRoutes() จะกำหนดเส้นทาง URL และ handler ต่างๆ:
	// - /api/v1/health (GET) - ตรวจสอบสถานะเซิร์ฟเวอร์
	// - /api/v1/auth/register (POST) - ลงทะเบียนผู้ใช้ใหม่
	// - /api/v1/auth/login (POST) - เข้าสู่ระบบ
	// - /api/v1/auth/profile (GET) - ดูข้อมูลโปรไฟล์ (ต้องเข้าสู่ระบบ)
	// - /api/v1/users/* (GET/DELETE) - จัดการผู้ใช้ (ต้องเป็น Admin)
	// - /swagger/* - เอกสาร API
	fmt.Println("🛣️  กำลังตั้งค่าเส้นทาง API...")
	routes.SetupRoutes(app, cfg)
	fmt.Println("✅ ตั้งค่าเส้นทาง API สำเร็จ")

	// ============================================
	// 5. เริ่มต้นเซิร์ฟเวอร์
	// ============================================
	// แสดงข้อความยินดีต้อนรับ
	fmt.Println("\n🎉 ยินดีต้อนรับสู่ GoTemplate API!")
	fmt.Println("Hello, World!")
	fmt.Println("ยินดีต้อนรับสู่โลกของ GoTemplate ของฉัน!")
	fmt.Println("Welcome to my GoTemplate world!")

	// แสดงข้อมูลการเริ่มต้นเซิร์ฟเวอร์
	fmt.Println("\n📋 ข้อมูลเซิร์ฟเวอร์:")
	log.Printf("📱 App Name: %s", cfg.App.Name)
	log.Printf("🔢 Version: %s", cfg.App.Version)
	log.Printf("📝 Description: %s", cfg.App.Description)
	log.Printf("🚀 Server starting on port %s", cfg.Server.Port)
	log.Printf("🌐 API Base URL: http://localhost:%s/api/v1", cfg.Server.Port)
	log.Printf("📚 Swagger docs available at: http://localhost:%s/swagger/", cfg.Server.Port)
	log.Printf("💾 Database: %s@%s:%s/%s", cfg.Database.User, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)
	log.Printf("🔐 JWT Expire: %s", cfg.JWT.Expire)
	log.Printf("🏗️  Environment: %s", cfg.Server.Environment)

	fmt.Println("\n🔗 เส้นทาง API ที่สำคัญ:")
	fmt.Println("   Health Check: http://localhost:%s/api/v1/health", cfg.Server.Port)
	fmt.Println("   Version: http://localhost:%s/api/v1/version", cfg.Server.Port)
	fmt.Println("   Swagger UI: http://localhost:%s/swagger/", cfg.Server.Port)

	fmt.Println("\n⚡ เซิร์ฟเวอร์พร้อมใช้งาน! กด Ctrl+C เพื่อหยุด")
	fmt.Println("=====================================")

	// เริ่มต้นเซิร์ฟเวอร์และฟังการร้องขอที่พอร์ตที่กำหนด
	// log.Fatal จะหยุดโปรแกรมหากเกิดข้อผิดพลาดในการเริ่มเซิร์ฟเวอร์
	log.Fatal(app.Listen(":" + cfg.Server.Port))

	// บรรทัดข้างล่างนี้จะไม่ทำงาน เพราะ log.Fatal จะหยุดโปรแกรมไปแล้ว
	// หากต้องการให้ข้อความเหล่านี้แสดง ต้องย้ายไปไว้ก่อน log.Fatal
}
