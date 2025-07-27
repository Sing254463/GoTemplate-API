// filepath: config/config.go
package config

import (
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

// Config struct คือโครงสร้างหลักที่เก็บการตั้งค่าทั้งหมดของแอปพลิเคชัน
// ประกอบด้วย 3 ส่วนหลัก: Database, JWT, และ Server configuration
type Config struct {
	Database *DatabaseConfig // การตั้งค่าเกี่ยวกับฐานข้อมูล
	JWT      *JWTConfig      // การตั้งค่าเกี่ยวกับ JSON Web Token
	Server   *ServerConfig   // การตั้งค่าเกี่ยวกับเซิร์ฟเวอร์
	App      *AppConfig      // เพิ่มการตั้งค่าเกี่ยวกับแอปพลิเคชัน
}

// DatabaseConfig struct เก็บข้อมูลการเชื่อมต่อฐานข้อมูล MySQL
type DatabaseConfig struct {
	Host     string   // ที่อยู่ของเซิร์ฟเวอร์ฐานข้อมูล (เช่น localhost)
	Port     string   // พอร์ตที่ใช้เชื่อมต่อ (เช่น 3306)
	User     string   // ชื่อผู้ใช้สำหรับเข้าสู่ฐานข้อมูล
	Password string   // รหัสผ่านสำหรับเข้าสู่ฐานข้อมูล
	DBName   string   // ชื่อของฐานข้อมูลที่จะใช้งาน
	DB       *sqlx.DB // ตัวแปรเก็บการเชื่อมต่อฐานข้อมูลจริง
}

// JWTConfig struct เก็บการตั้งค่าเกี่ยวกับ JWT (JSON Web Token)
type JWTConfig struct {
	Secret string        // กุญแจลับสำหรับเซ็น JWT token
	Expire time.Duration // ระยะเวลาที่ token จะหมดอายุ
}

// ServerConfig struct เก็บการตั้งค่าเกี่ยวกับเซิร์ฟเวอร์
type ServerConfig struct {
	Port        string // พอร์ตที่เซิร์ฟเวอร์จะรัน (เช่น 8080)
	Environment string // สภาพแวดล้อมการทำงาน (development, production)
}

// เพิ่ม AppConfig struct สำหรับข้อมูลแอปพลิเคชัน
type AppConfig struct {
	Name        string // ชื่อแอปพลิเคชัน
	Version     string // เวอร์ชันของแอปพลิเคชัน
	Description string // คำอธิบายแอปพลิเคชัน
}

// LoadConfig ฟังก์ชันหลักสำหรับโหลดการตั้งค่าทั้งหมด
// จะอ่านค่าจากไฟล์ .env และสร้าง Config object พร้อมเชื่อมต่อฐานข้อมูล
func LoadConfig() *Config {
	// โหลดไฟล์ .env เพื่ออ่านตัวแปร environment variables
	// หากไม่พบไฟล์ .env จะแสดงข้อความแจ้งเตือน แต่จะไม่หยุดการทำงาน
	if err := godotenv.Load(); err != nil {
		log.Println("ไม่พบไฟล์ .env / No .env file found")
	}

	// สร้าง Config object ใหม่โดยอ่านค่าจากตัวแปร environment
	// หากไม่มีค่าใน environment จะใช้ค่า default ที่กำหนดไว้
	config := &Config{
		Database: &DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),  // ค่าเริ่มต้น: localhost
			Port:     getEnv("DB_PORT", "3306"),       // ค่าเริ่มต้น: 3306 (MySQL default port)
			User:     getEnv("DB_USER", "root"),       // ค่าเริ่มต้น: root
			Password: getEnv("DB_PASSWORD", ""),       // ค่าเริ่มต้น: ว่าง (ไม่มีรหัสผ่าน)
			DBName:   getEnv("DB_NAME", "gotemplate"), // ค่าเริ่มต้น: gotemplate
		},
		JWT: &JWTConfig{
			Secret: getEnv("JWT_SECRET", "default-secret"),     // ค่าเริ่มต้น: default-secret
			Expire: parseDuration(getEnv("JWT_EXPIRE", "24h")), // ค่าเริ่มต้น: 24 ชั่วโมง
		},
		Server: &ServerConfig{
			Port:        getEnv("PORT", "8080"),               // ค่าเริ่มต้น: 8080
			Environment: getEnv("ENVIRONMENT", "development"), // ค่าเริ่มต้น: development
		},
		App: &AppConfig{
			Name:        getEnv("APP_NAME", "GoTemplate API"),                                                            // ค่าเริ่มต้น: GoTemplate API
			Version:     getEnv("APP_VERSION", "1.0.0"),                                                                  // ค่าเริ่มต้น: 1.0.0
			Description: getEnv("APP_DESCRIPTION", "A complete Go template API with authentication and user management"), // ค่าเริ่มต้น: คำอธิบายแอปพลิเคชัน
		},
	}

	// เริ่มการเชื่อมต่อกับฐานข้อมูล
	config.Database.ConnectDB()

	return config
}

// ConnectDB method สำหรับสร้างการเชื่อมต่อกับฐานข้อมูล MySQL
// ใช้กับ DatabaseConfig struct
func (d *DatabaseConfig) ConnectDB() {
	// สร้าง DSN (Data Source Name) string สำหรับเชื่อมต่อ MySQL
	// รูปแบบ: user:password@tcp(host:port)/database?parameters
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		d.User,     // ชื่อผู้ใช้
		d.Password, // รหัสผ่าน
		d.Host,     // ที่อยู่เซิร์ฟเวอร์
		d.Port,     // พอร์ต
		d.DBName)   // ชื่อฐานข้อมูล

	// พารามิเตอร์เพิ่มเติม:
	// charset=utf8mb4: รองรับการเข้ารหัส UTF-8 แบบเต็ม (รวม emoji)
	// parseTime=True: แปลง MySQL time/date เป็น Go time.Time อัตโนมัติ
	// loc=Local: ใช้ timezone ของเครื่องท้องถิ่น

	// เชื่อมต่อกับฐานข้อมูลโดยใช้ sqlx library
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		// หากเชื่อมต่อไม่สำเร็จ จะหยุดการทำงานของโปรแกรมทันที
		log.Fatal("ไม่สามารถเชื่อมต่อกับฐานข้อมูล/Failed to connect to database:", err)
	}

	// ตั้งค่า Connection Pool เพื่อจัดการการเชื่อมต่อแบบมีประสิทธิภาพ
	db.SetMaxOpenConns(25)                 // จำนวนการเชื่อมต่อสูงสุดที่เปิดพร้อมกัน
	db.SetMaxIdleConns(5)                  // จำนวนการเชื่อมต่อที่เก็บไว้ในสถานะรอ
	db.SetConnMaxLifetime(5 * time.Minute) // ระยะเวลาที่การเชื่อมต่อหนึ่งจะมีชีวิตอยู่สูงสุด

	// เก็บการเชื่อมต่อไว้ใน struct เพื่อใช้งานต่อไป
	d.DB = db

	// แสดงข้อความยืนยันการเชื่อมต่อสำเร็จ
	log.Println("Database connected successfully")
	log.Println("ฐานข้อมูลเชื่อมต่อสำเร็จแล้ว")
}

// getEnv ฟังก์ชันช่วยสำหรับอ่านค่าจากตัวแปร environment
// หากไม่พบค่าที่ต้องการ จะคืนค่า default ที่กำหนดไว้
func getEnv(key, defaultValue string) string {
	// ตรวจสอบว่ามีตัวแปร environment ชื่อ 'key' หรือไม่
	if value := os.Getenv(key); value != "" {
		return value // หากมีค่า ให้คืนค่าที่พบ
	}
	return defaultValue // หากไม่มีค่า ให้คืนค่า default
}

// parseDuration ฟังก์ชันช่วยสำหรับแปลง string เป็น time.Duration
// ใช้สำหรับการตั้งค่าระยะเวลาหมดอายุของ JWT token
func parseDuration(s string) time.Duration {
	// พยายามแปลง string เป็น Duration (เช่น "24h", "30m", "1h30m")
	duration, err := time.ParseDuration(s)
	if err != nil {
		// หากแปลงไม่สำเร็จ ให้คืนค่า default เป็น 24 ชั่วโมง
		return 24 * time.Hour
	}
	return duration // หากแปลงสำเร็จ ให้คืนค่าที่แปลงได้
}
