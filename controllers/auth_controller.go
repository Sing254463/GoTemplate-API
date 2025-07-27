package controllers

import (
	"database/sql"
	"time"

	"github.com/Sing254463/GoTemplate/Backend/config"
	"github.com/Sing254463/GoTemplate/Backend/models"
	"github.com/Sing254463/GoTemplate/Backend/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// AuthController โครงสร้างสำหรับจัดการการยืนยันตัวตน
// ประกอบด้วย Config (การตั้งค่า) และ Validator (ตัวตรวจสอบข้อมูล)
type AuthController struct {
	Config    *config.Config      // การตั้งค่าระบบ (ฐานข้อมูล, JWT, เซิร์ฟเวอร์)
	Validator *validator.Validate // ตัวตรวจสอบความถูกต้องของข้อมูล
}

// NewAuthController ฟังก์ชันสร้าง AuthController ใหม่
// รับพารามิเตอร์ cfg (การตั้งค่า) และคืนค่า pointer ของ AuthController
func NewAuthController(cfg *config.Config) *AuthController {
	return &AuthController{
		Config:    cfg,             // เก็บการตั้งค่าที่ได้รับ
		Validator: validator.New(), // สร้างตัวตรวจสอบข้อมูลใหม่
	}
}

// Register ฟังก์ชันสำหรับลงทะเบียนผู้ใช้ใหม่
// รับข้อมูล username, email, password และสร้างบัญชีผู้ใช้ใหม่
// @Summary Register a new user
// @Description Register a new user with username, email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.UserRegister true "User registration data"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 409 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/register [post]
func (ac *AuthController) Register(c *fiber.Ctx) error {
	var userRegister models.UserRegister

	// แปลงข้อมูล JSON จาก request body เป็น struct
	if err := c.BodyParser(&userRegister); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "ข้อมูลที่ส่งมาไม่ถูกต้อง", err)
	}

	// ตรวจสอบความถูกต้องของข้อมูล (validation)
	// เช่น email ต้องเป็นรูปแบบอีเมล, username ต้องมีความยาว 3-20 ตัวอักษร
	if err := ac.Validator.Struct(&userRegister); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "ข้อมูลไม่ผ่านการตรวจสอบ", err)
	}

	// ตรวจสอบว่ามีผู้ใช้ที่มี email หรือ username นี้อยู่แล้วหรือไม่
	var existingUser models.User
	query := "SELECT id FROM users WHERE email = ? OR username = ?"
	err := ac.Config.Database.DB.Get(&existingUser, query, userRegister.Email, userRegister.Username)
	if err == nil {
		// หากพบผู้ใช้ที่มีข้อมูลซ้ำ ให้ส่งข้อผิดพลาดกลับ
		return utils.ErrorResponse(c, fiber.StatusConflict, "มีผู้ใช้นี้อยู่แล้ว", nil)
	} else if err != sql.ErrNoRows {
		// หากเกิดข้อผิดพลาดอื่นๆ ในฐานข้อมูล
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "เกิดข้อผิดพลาดในฐานข้อมูล", err)
	}

	// เข้ารหัสรหัสผ่านด้วย bcrypt เพื่ือความปลอดภัย
	hashedPassword, err := utils.HashPassword(userRegister.Password)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "ไม่สามารถเข้ารหัสรหัสผ่านได้", err)
	}

	// สร้าง struct ผู้ใช้ใหม่พร้อมข้อมูลที่จำเป็น
	user := models.User{
		Username:  userRegister.Username, // ชื่อผู้ใช้
		Email:     userRegister.Email,    // อีเมล
		Password:  hashedPassword,        // รหัสผ่านที่เข้ารหัสแล้ว
		Role:      "user",                // สิทธิ์เริ่มต้นเป็น user (ไม่ใช่ admin)
		CreatedAt: time.Now(),            // เวลาที่สร้างบัญชี
		UpdatedAt: time.Now(),            // เวลาที่อัปเดตล่าสุด
	}

	// บันทึกข้อมูลผู้ใช้ใหม่ลงในฐานข้อมูล
	query = `INSERT INTO users (username, email, password, role, created_at, updated_at) 
             VALUES (?, ?, ?, ?, ?, ?)`
	result, err := ac.Config.Database.DB.Exec(query, user.Username, user.Email, user.Password, user.Role, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "ไม่สามารถสร้างผู้ใช้ได้", err)
	}

	// ดึง ID ของผู้ใช้ที่เพิ่งสร้างจากฐานข้อมูล
	userID, _ := result.LastInsertId()
	user.ID = int(userID)

	// ส่งผลลัพธ์การลงทะเบียนสำเร็จกลับไป (ไม่รวมรหัสผ่าน)
	return utils.CreatedResponse(c, "ลงทะเบียนผู้ใช้สำเร็จ", user.ConvertToResponse())
}

// Login ฟังก์ชันสำหรับเข้าสู่ระบบ
// รับ email และ password แล้วตรวจสอบความถูกต้อง
// หากถูกต้องจะสร้าง JWT token ให้
// @Summary Login user
// @Description Login with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.UserLogin true "User login data"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/login [post]
func (ac *AuthController) Login(c *fiber.Ctx) error {
	var userLogin models.UserLogin

	// แปลงข้อมูล JSON จาก request body เป็น struct
	if err := c.BodyParser(&userLogin); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "ข้อมูลที่ส่งมาไม่ถูกต้อง", err)
	}

	// ตรวจสอบความถูกต้องของข้อมูล (email และ password จำเป็นต้องมี)
	if err := ac.Validator.Struct(&userLogin); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "ข้อมูลไม่ผ่านการตรวจสอบ", err)
	}

	// ค้นหาผู้ใช้ในฐานข้อมูลด้วย email
	var user models.User
	query := "SELECT id, username, email, password, role FROM users WHERE email = ?"
	err := ac.Config.Database.DB.Get(&user, query, userLogin.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			// หากไม่พบผู้ใช้ที่มี email นี้
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "อีเมลหรือรหัสผ่านไม่ถูกต้อง", nil)
		}
		// หากเกิดข้อผิดพลาดอื่นๆ ในฐานข้อมูล
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "เกิดข้อผิดพลาดในฐานข้อมูล", err)
	}

	// ตรวจสอบรหัสผ่าน โดยเปรียบเทียบกับรหัสผ่านที่เข้ารหัสไว้ในฐานข้อมูล
	if err := utils.CheckPassword(user.Password, userLogin.Password); err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "อีเมลหรือรหัสผ่านไม่ถูกต้อง", nil)
	}

	// สร้าง JWT token สำหรับผู้ใช้ที่เข้าสู่ระบบสำเร็จ
	// token จะมีข้อมูล user ID, username, role และจะหมดอายุตามที่กำหนดใน config
	token, err := utils.GenerateJWT(user.ID, user.Username, user.Role, ac.Config.JWT.Secret, ac.Config.JWT.Expire)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "ไม่สามารถสร้าง token ได้", err)
	}

	// ส่งผลลัพธ์การเข้าสู่ระบบสำเร็จพร้อม token และข้อมูลผู้ใช้
	return utils.SuccessResponse(c, "เข้าสู่ระบบสำเร็จ", fiber.Map{
		"token": token,                    // JWT token สำหรับการยืนยันตัวตน
		"user":  user.ConvertToResponse(), // ข้อมูลผู้ใช้ (ไม่รวมรหัสผ่าน)
	})
}

// GetProfile ฟังก์ชันสำหรับดูข้อมูลโปรไฟล์ของผู้ใช้ที่เข้าสู่ระบบ
// ต้องส่ง JWT token มาด้วยจึงจะใช้งานได้
// @Summary Get user profile
// @Description Get current user profile
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/profile [get]
func (ac *AuthController) GetProfile(c *fiber.Ctx) error {
	// ดึง user ID จาก JWT token ที่ได้รับการตรวจสอบแล้วโดย middleware
	userID := c.Locals("user_id").(int)

	// ค้นหาข้อมูลผู้ใช้ในฐานข้อมูลด้วย ID
	var user models.User
	query := "SELECT id, username, email, role, created_at, updated_at FROM users WHERE id = ?"
	err := ac.Config.Database.DB.Get(&user, query, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "ไม่สามารถดึงข้อมูลโปรไฟล์ได้", err)
	}

	// ส่งข้อมูลโปรไฟล์กลับไป (ไม่รวมรหัสผ่าน)
	return utils.SuccessResponse(c, "ดึงข้อมูลโปรไฟล์สำเร็จ", user.ConvertToResponse())
}
