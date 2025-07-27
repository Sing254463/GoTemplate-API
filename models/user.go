package models

import (
	"time"
)

// User โครงสร้างหลักสำหรับเก็บข้อมูลผู้ใช้ในฐานข้อมูล
// ใช้ tags สำหรับ JSON serialization, database mapping, และ validation
type User struct {
	ID        int       `json:"id" db:"id"`                                                 // ID ผู้ใช้ (Primary Key)
	Username  string    `json:"username" db:"username" validate:"required,min=3,max=20"`    // ชื่อผู้ใช้ (3-20 ตัวอักษร)
	Email     string    `json:"email" db:"email" validate:"required,email"`                 // อีเมล (ต้องเป็นรูปแบบอีเมล)
	Password  string    `json:"password,omitempty" db:"password" validate:"required,min=6"` // รหัสผ่าน (ขั้นต่ำ 6 ตัวอักษร, omitempty = ไม่แสดงใน JSON)
	Role      string    `json:"role" db:"role"`                                             // สิทธิ์ผู้ใช้ (user/admin)
	CreatedAt time.Time `json:"created_at" db:"created_at"`                                 // วันที่สร้างบัญชี
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`                                 // วันที่อัปเดตล่าสุด
}

// UserLogin โครงสร้างสำหรับรับข้อมูลการเข้าสู่ระบบ
// ใช้เฉพาะการเข้าสู่ระบบ ไม่เก็บในฐานข้อมูล
type UserLogin struct {
	Email    string `json:"email" validate:"required,email"` // อีเมล (จำเป็นต้องมี และเป็นรูปแบบอีเมล)
	Password string `json:"password" validate:"required"`    // รหัสผ่าน (จำเป็นต้องมี)
}

// UserRegister โครงสร้างสำหรับรับข้อมูลการลงทะเบียน
// ใช้เฉพาะการลงทะเบียน ไม่เก็บในฐานข้อมูล
type UserRegister struct {
	Username string `json:"username" validate:"required,min=3,max=20"` // ชื่อผู้ใช้ (จำเป็น, 3-20 ตัวอักษร)
	Email    string `json:"email" validate:"required,email"`           // อีเมล (จำเป็น, รูปแบบอีเมล)
	Password string `json:"password" validate:"required,min=6"`        // รหัสผ่าน (จำเป็น, ขั้นต่ำ 6 ตัวอักษร)
}

// UserResponse โครงสร้างสำหรับส่งข้อมูลผู้ใช้กลับไป
// ไม่รวมข้อมูลที่อ่อนไหว เช่น รหัสผ่าน, เวลาสร้าง/อัปเดต
type UserResponse struct {
	ID       int    `json:"id"`       // ID ผู้ใช้
	Username string `json:"username"` // ชื่อผู้ใช้
	Email    string `json:"email"`    // อีเมล
	Role     string `json:"role"`     // สิทธิ์ผู้ใช้
}

// ConvertToResponse method สำหรับแปลง User เป็น UserResponse
// ใช้เพื่อซ่อนข้อมูลที่อ่อนไหวก่อนส่งกลับไปยัง client
// เช่น รหัสผ่าน, เวลาสร้าง/อัปเดต ที่ไม่จำเป็นต้องแสดง
func (u *User) ConvertToResponse() UserResponse {
	return UserResponse{
		ID:       u.ID,       // คัดลอก ID
		Username: u.Username, // คัดลอกชื่อผู้ใช้
		Email:    u.Email,    // คัดลอกอีเมล
		Role:     u.Role,     // คัดลอกสิทธิ์
		// ไม่รวม Password, CreatedAt, UpdatedAt เพื่อความปลอดภัย
	}
}
