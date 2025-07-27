package controllers

import (
	"strconv"

	"github.com/Sing254463/GoTemplate/Backend/config"
	"github.com/Sing254463/GoTemplate/Backend/models"
	"github.com/Sing254463/GoTemplate/Backend/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// UserController โครงสร้างสำหรับจัดการข้อมูลผู้ใช้
// ใช้สำหรับผู้ดูแลระบบ (Admin) ในการจัดการผู้ใช้ต่างๆ
type UserController struct {
	Config    *config.Config      // การตั้งค่าระบบ
	Validator *validator.Validate // ตัวตรวจสอบความถูกต้องของข้อมูล
}

// NewUserController ฟังก์ชันสร้าง UserController ใหม่
func NewUserController(cfg *config.Config) *UserController {
	return &UserController{
		Config:    cfg,
		Validator: validator.New(),
	}
}

// GetAllUsers ฟังก์ชันสำหรับดูรายชื่อผู้ใช้ทั้งหมด (เฉพาะ Admin)
// @Summary Get all users
// @Description Get all users (Admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /users [get]
func (uc *UserController) GetAllUsers(c *fiber.Ctx) error {
	// สร้างตัวแปรสำหรับเก็บรายชื่อผู้ใช้ทั้งหมด
	var users []models.User

	// Query ดึงข้อมูลผู้ใช้ทั้งหมดจากฐานข้อมูล (ไม่รวมรหัสผ่าน)
	query := "SELECT id, username, email, role, created_at, updated_at FROM users"
	err := uc.Config.Database.DB.Select(&users, query)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "ไม่สามารถดึงข้อมูลผู้ใช้ได้", err)
	}

	// แปลงข้อมูลผู้ใช้เป็นรูปแบบที่จะส่งกลับ (ซ่อนข้อมูลที่ไม่จำเป็น)
	var userResponses []models.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, user.ConvertToResponse())
	}

	// ส่งรายชื่อผู้ใช้ทั้งหมดกลับไป
	return utils.SuccessResponse(c, "ดึงข้อมูลผู้ใช้ทั้งหมดสำเร็จ", userResponses)
}

// GetUserByID ฟังก์ชันสำหรับดูข้อมูลผู้ใช้ตาม ID (เฉพาะ Admin)
// @Summary Get user by ID
// @Description Get user by ID (Admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "User ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /users/{id} [get]
func (uc *UserController) GetUserByID(c *fiber.Ctx) error {
	// แปลงพารามิเตอร์ id จาก string เป็น integer
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "ID ผู้ใช้ไม่ถูกต้อง", err)
	}

	// ค้นหาผู้ใช้ในฐานข้อมูลด้วย ID
	var user models.User
	query := "SELECT id, username, email, role, created_at, updated_at FROM users WHERE id = ?"
	err = uc.Config.Database.DB.Get(&user, query, id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "ไม่พบผู้ใช้", err)
	}

	// ส่งข้อมูลผู้ใช้ที่พบกลับไป
	return utils.SuccessResponse(c, "ดึงข้อมูลผู้ใช้สำเร็จ", user.ConvertToResponse())
}

// DeleteUser ฟังก์ชันสำหรับลบผู้ใช้ตาม ID (เฉพาะ Admin)
// @Summary Delete user
// @Description Delete user by ID (Admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "User ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /users/{id} [delete]
func (uc *UserController) DeleteUser(c *fiber.Ctx) error {
	// แปลงพารามิเตอร์ id จาก string เป็น integer
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "ID ผู้ใช้ไม่ถูกต้อง", err)
	}

	// ตรวจสอบว่ามีผู้ใช้ที่มี ID นี้อยู่หรือไม่
	var user models.User
	query := "SELECT id FROM users WHERE id = ?"
	err = uc.Config.Database.DB.Get(&user, query, id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "ไม่พบผู้ใช้", err)
	}

	// ลบผู้ใช้จากฐานข้อมูล
	deleteQuery := "DELETE FROM users WHERE id = ?"
	_, err = uc.Config.Database.DB.Exec(deleteQuery, id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "ไม่สามารถลบผู้ใช้ได้", err)
	}

	// ส่งผลลัพธ์การลบสำเร็จกลับไป
	return utils.SuccessResponse(c, "ลบผู้ใช้สำเร็จ", nil)
}
