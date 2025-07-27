# GoTemplate API

เทมเพลต API ที่สมบูรณ์สำหรับ Go โดยใช้ Fiber framework พร้อมระบบยืนยันตัวตนและการจัดการผู้ใช้

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![Fiber](https://img.shields.io/badge/Fiber-v2.52.0-00ADD8?style=for-the-badge&logo=fiber)](https://gofiber.io/)
[![MySQL](https://img.shields.io/badge/MySQL-8.0+-4479A1?style=for-the-badge&logo=mysql&logoColor=white)](https://www.mysql.com/)
[![JWT](https://img.shields.io/badge/JWT-000000?style=for-the-badge&logo=JSON%20web%20tokens)](https://jwt.io/)
[![Swagger](https://img.shields.io/badge/Swagger-85EA2D?style=for-the-badge&logo=swagger&logoColor=black)](https://swagger.io/)
[![Air](https://img.shields.io/badge/Air-Hot%20Reload-FF6B6B?style=for-the-badge&logo=go&logoColor=white)](https://github.com/cosmtrek/air)
[![bcrypt](https://img.shields.io/badge/bcrypt-Password%20Hash-4B8BBE?style=for-the-badge&logo=security&logoColor=white)](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
[![SQLX](https://img.shields.io/badge/SQLX-Database-336791?style=for-the-badge&logo=postgresql&logoColor=white)](https://github.com/jmoiron/sqlx)
## 📋 สารบัญ

- [คุณสมบัติ](#-คุณสมบัติ)
- [โครงสร้างโปรเจค](#-โครงสร้างโปรเจค)
- [เทคโนโลยีที่ใช้](#-เทคโนโลยีที่ใช้)
- [การติดตั้ง](#-การติดตั้ง)
- [การตั้งค่า](#-การตั้งค่า)
- [การรันโปรแกรม](#-การรันโปรแกรม)
- [API Endpoints](#-api-endpoints)
- [การป้องกันและความปลอดภัย](#-การป้องกันและความปลอดภัย)
- [การใช้งาน Swagger](#-การใช้งาน-swagger)
- [การพัฒนา](#-การพัฒนา)
- [ลิขสิทธิ์](#-ลิขสิทธิ์)

## 🚀 คุณสมบัติ

### ✨ คุณสมบัติหลัก
- 🔐 **ระบบยืนยันตัวตน JWT** - การเข้าสู่ระบบและลงทะเบียนที่ปลอดภัย
- 👥 **การจัดการผู้ใช้** - CRUD operations สำหรับ Admin
- 🛡️ **การควบคุมสิทธิ์** - Role-based access control (User/Admin)
- 🔒 **เข้ารหัสรหัสผ่าน** - bcrypt hashing สำหรับความปลอดภัย
- 📝 **เอกสาร API อัตโนมัติ** - Swagger/OpenAPI documentation
- 🏥 **Health Check** - ตรวจสอบสถานะเซิร์ฟเวอร์
- 📊 **การบันทึก Log** - ติดตามการใช้งาน API
- 🌐 **CORS Support** - รองรับการเชื่อมต่อข้าม domain

### 🔧 คุณสมบัติทางเทคนิค
- ⚡ **Hot Reload** - พัฒนาได้อย่างรวดเร็วด้วย Air
- 🗄️ **Database Connection Pooling** - จัดการการเชื่อมต่อฐานข้อมูลอย่างมีประสิทธิภาพ
- ✅ **Data Validation** - ตรวจสอบความถูกต้องของข้อมูล
- 🚫 **Error Handling** - จัดการข้อผิดพลาดแบบ Global
- 🔄 **Middleware Pipeline** - ระบบ middleware ที่ยืดหยุ่น

## 📁 โครงสร้างโปรเจค

```
GoTemplate/
├── 📄 main.go                 # จุดเริ่มต้นของแอปพลิเคชัน
├── 📄 go.mod                  # Go modules และ dependencies
├── 📄 go.sum                  # Lock file สำหรับ dependencies
├── 📄 .env                    # ตัวแปรสภาพแวดล้อม
├── 📄 .air.toml               # การตั้งค่า hot reload
├── 📄 README.md               # เอกสารโปรเจค
│
├── 📁 config/                 # การตั้งค่าระบบ
│   └── 📄 config.go           # จัดการการตั้งค่าและการเชื่อมต่อ DB
│
├── 📁 controllers/            # ตัวควบคุม API handlers
│   ├── 📄 auth_controller.go  # การจัดการยืนยันตัวตน
│   └── 📄 user_controller.go  # การจัดการผู้ใช้
│
├── 📁 middleware/             # ตัวกลางประมวลผล
│   ├── 📄 jwt_middleware.go   # ตรวจสอบ JWT และสิทธิ์
│   └── 📄 logger.go           # บันทึก log การใช้งาน
│
├── 📁 models/                 # โครงสร้างข้อมูล
│   └── 📄 user.go             # โมเดลผู้ใช้และโครงสร้างข้อมูล
│
├── 📁 routes/                 # เส้นทาง API
│   └── 📄 routes.go           # กำหนดเส้นทางและ middleware
│
├── 📁 utils/                  # ฟังก์ชันช่วยเหลือ
│   ├── 📄 hash.go             # เข้ารหัสและตรวจสอบรหัสผ่าน
│   ├── 📄 jwt.go              # จัดการ JWT tokens
│   └── 📄 response.go         # รูปแบบการตอบกลับมาตรฐาน
│
├── 📁 docs/                   # เอกสาร API
│   ├── 📄 docs.go             # Generated Swagger docs
│   ├── 📄 swagger.json        # Swagger specification (JSON)
│   └── 📄 swagger.yaml        # Swagger specification (YAML)
│
├── 📁 bin/                    # ไฟล์ executable
│   └── 📄 main.exe            # ไฟล์โปรแกรมที่คอมไพล์แล้ว
│
└── 📁 tmp/                    # ไฟล์ชั่วคระหว่างพัฒนา
    └── 📄 main.exe            # ไฟล์ temporary สำหรับ hot reload
```

## 🛠️ เทคโนโลยีที่ใช้

### Backend Framework & Libraries
| เทคโนโลยี | เวอร์ชัน | วัตถุประสงค์ |
|-----------|---------|-------------|
| **Go** | 1.21+ | ภาษาโปรแกรมหลัก |
| **Fiber** | v2.52.0 | Web framework ประสิทธิภาพสูง |
| **SQLX** | v1.3.5 | Database driver ที่ขยายจาก database/sql |
| **MySQL Driver** | v1.7.1 | เชื่อมต่อกับฐานข้อมูล MySQL |

### Authentication & Security
| เทคโนโลยี | เวอร์ชัน | วัตถุประสงค์ |
|-----------|---------|-------------|
| **JWT** | v5.2.0 | JSON Web Tokens สำหรับการยืนยันตัวตน |
| **bcrypt** | crypto/v0.18.0 | เข้ารหัสรหัสผ่านแบบ one-way |
| **Validator** | v10.16.0 | ตรวจสอบความถูกต้องของข้อมูล |

### Development Tools
| เทคโนโลยี | เวอร์ชัน | วัตถุประสงค์ |
|-----------|---------|-------------|
| **Air** | Latest | Hot reload สำหรับการพัฒนา |
| **Swagger** | v1.0.0 | สร้างเอกสาร API อัตโนมัติ |
| **Godotenv** | v1.5.1 | โหลดตัวแปรจากไฟล์ .env |

### Database
| เทคโนโลยี | เวอร์ชัน | วัตถุประสงค์ |
|-----------|---------|-------------|
| **MySQL** | 8.0+ | ฐานข้อมูลหลัก |

## ⚙️ การติดตั้ง

### ข้อกำหนดของระบบ
- Go 1.21 หรือสูงกว่า
- MySQL 8.0 หรือสูงกว่า
- Git

### 1. Clone โปรเจค
```bash
git clone https://github.com/Sing254463/GoTemplate-API.git
cd GoTemplate-API
```

### 2. ติดตั้ง Dependencies
```bash
go mod download
```

### 3. ติดตั้ง Development Tools
```bash
# ติดตั้ง Air สำหรับ hot reload
go install github.com/cosmtrek/air@latest

# ติดตั้ง Swag สำหรับ generate documentation
go install github.com/swaggo/swag/cmd/swag@latest
```

### 4. สร้างฐานข้อมูล
```sql
CREATE DATABASE apitest;
USE apitest;

CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role ENUM('user', 'admin') DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- สร้าง admin user เริ่มต้น (password: admin123)
INSERT INTO users (username, email, password, role) VALUES 
('admin', 'admin@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin');
```

## 🔧 การตั้งค่า

### สร้างไฟล์ .env
สร้างไฟล์ `.env` ในโฟลเดอร์หลักและกำหนดค่าดังนี้:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=your_DB

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-here
JWT_EXPIRE=24h

# Server Configuration
PORT=8080
ENVIRONMENT=development

# Application Configuration
APP_NAME=GoTemplate API
APP_VERSION=1.0.0
APP_DESCRIPTION=Go API template with authentication and user management
```

### คำอธิบายการตั้งค่า

| ตัวแปร | คำอธิบาย | ค่าเริ่มต้น |
|--------|----------|------------|
| `DB_HOST` | ที่อยู่เซิร์ฟเวอร์ฐานข้อมูล | localhost |
| `DB_PORT` | พอร์ตฐานข้อมูล MySQL | 3306 |
| `DB_USER` | ชื่อผู้ใช้ฐานข้อมูล | root |
| `DB_PASSWORD` | รหัสผ่านฐานข้อมูล | - |
| `DB_NAME` | ชื่อฐานข้อมูล | apitest |
| `JWT_SECRET` | กุญแจลับสำหรับ JWT | - |
| `JWT_EXPIRE` | ระยะเวลาหมดอายุ JWT | 24h |
| `PORT` | พอร์ตเซิร์ฟเวอร์ | 8080 |
| `ENVIRONMENT` | สภาพแวดล้อม | development |

## 🚀 การรันโปรแกรม

### 🔥 Development Mode (สำหรับนักพัฒนา)
```bash
air
```

### 🏭 Production Mode (สำหรับผู้ใช้ทั่วไป)
**กดที่ไฟล์ `run-server.bat`**

หรือรันแบบ manual:
```bash
# Windows
go build -o tmp/main.exe .
tmp\main.exe

# Linux/Mac  
go build -o tmp/main .
./tmp/main
```

### Generate Swagger Documentation
```bash
swag init -g main.go --output docs/
```

## 📡 API Endpoints

### 🔓 Public Endpoints (ไม่ต้องเข้าสู่ระบบ)

| Method | Endpoint | คำอธิบาย |
|--------|----------|----------|
| `GET` | `/api/v1/health` | ตรวจสอบสถานะเซิร์ฟเวอร์ |
| `GET` | `/api/v1/version` | ข้อมูลเวอร์ชันแอปพลิเคชัน |
| `POST` | `/api/v1/auth/register` | ลงทะเบียนผู้ใช้ใหม่ |
| `POST` | `/api/v1/auth/login` | เข้าสู่ระบบ |
| `GET` | `/swagger/*` | เอกสาร API |

### 🔒 Protected Endpoints (ต้องเข้าสู่ระบบ)

| Method | Endpoint | คำอธิบาย | สิทธิ์ |
|--------|----------|----------|-------|
| `GET` | `/api/v1/auth/profile` | ดูข้อมูลโปรไฟล์ | User/Admin |

### 👑 Admin Only Endpoints (เฉพาะ Admin)

| Method | Endpoint | คำอธิบาย |
|--------|----------|----------|
| `GET` | `/api/v1/users` | ดูรายชื่อผู้ใช้ทั้งหมด |
| `GET` | `/api/v1/users/{id}` | ดูข้อมูลผู้ใช้ตาม ID |
| `DELETE` | `/api/v1/users/{id}` | ลบผู้ใช้ตาม ID |

### ตัวอย่างการใช้งาน

#### ลงทะเบียนผู้ใช้ใหม่
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com", 
    "password": "password123"
  }'
```

#### เข้าสู่ระบบ
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

#### ดูข้อมูลโปรไฟล์ (ใช้ JWT Token)
```bash
curl -X GET http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 🛡️ การป้องกันและความปลอดภัย

### 🔐 การยืนยันตัวตนและการควบคุมสิทธิ์

#### JWT (JSON Web Tokens)
- **การทำงาน**: ใช้ JWT สำหรับการยืนยันตัวตนแบบ stateless
- **ข้อมูลใน Token**: User ID, Username, Role, เวลาหมดอายุ
- **การเซ็น**: ใช้ HMAC SHA-256 กับ secret key
- **หมดอายุ**: กำหนดได้ในไฟล์ .env (ค่าเริ่มต้น 24 ชั่วโมง)

#### Role-Based Access Control (RBAC)
- **User Role**: สามารถเข้าถึงข้อมูลส่วนตัวและ API พื้นฐาน
- **Admin Role**: เข้าถึงการจัดการผู้ใช้และ API ทั้งหมด
- **Middleware**: `AdminMiddleware()` ตรวจสอบสิทธิ์ Admin

### 🔒 การเข้ารหัสรหัสผ่าน

#### bcrypt Hashing
- **อัลกอริธึม**: bcrypt with default cost (10)
- **Salt**: สุ่มอัตโนมัติสำหรับแต่ละรหัสผ่าน
- **One-way**: ไม่สามารถถอดรหัสกลับได้
- **การตรวจสอบ**: เปรียบเทียบ hash แทนการถอดรหัส

```go
// ตัวอย่างการเข้ารหัส
hashedPassword, err := utils.HashPassword("password123")

// ตัวอย่างการตรวจสอบ
err := utils.CheckPassword(hashedPassword, "password123")
```

### 🛡️ การป้องกันข้อมูล

#### Input Validation
- **Library**: go-playground/validator/v10
- **การตรวจสอบ**: Email format, ความยาวรหัสผ่าน, ชื่อผู้ใช้
- **Custom Tags**: รองรับการตรวจสอบแบบกำหนดเอง

```go
type UserRegister struct {
    Username string `validate:"required,min=3,max=20"`
    Email    string `validate:"required,email"`
    Password string `validate:"required,min=6"`
}
```

#### SQL Injection Prevention
- **Prepared Statements**: ใช้ placeholder (?) ใน SQL queries
- **SQLX Library**: ป้องกัน SQL injection อัตโนมัติ
- **Parameter Binding**: ไม่ concatenate ข้อมูลเข้ากับ SQL string

#### CORS (Cross-Origin Resource Sharing)
- **การตั้งค่า**: รองรับการเชื่อมต่อข้าม domain
- **Headers**: กำหนด allowed headers และ methods
- **Origins**: ควบคุมได้ว่า domain ไหนเข้าถึงได้

### 🚨 Error Handling และ Security Headers

#### Global Error Handler
- **การจัดการ**: จัดการ error แบบรวมศูนย์
- **ไม่เปิดเผยข้อมูล**: ซ่อนข้อมูลที่อ่อนไหวในกรณี error
- **Standard Response**: รูปแบบการตอบกลับที่สม่ำเสมอ

#### Security Best Practices
- **Password Policy**: ขั้นต่ำ 6 ตัวอักษร (สามารถปรับได้)
- **JWT Expiration**: กำหนดเวลาหมดอายุป้องกันการใช้งานไม่สิ้นสุด
- **Database Connection Pool**: จำกัดจำนวนการเชื่อมต่อ
- **Environment Variables**: เก็บข้อมูลสำคัญในไฟล์ .env

### 🔍 การติดตาม (Monitoring)

#### Request Logging
- **ข้อมูลที่บันทึก**: เวลา, สถานะ, ระยะเวลา, IP, Method, Path
- **รูปแบบ**: RFC3339 timestamp format
- **การใช้งาน**: สามารถนำไปวิเคราะห์และติดตามปัญหาได้

#### Health Monitoring
- **Health Check Endpoint**: `/api/v1/health`
- **ข้อมูล**: สถานะเซิร์ฟเวอร์, ข้อมูลแอปพลิเคชัน, เวลาปัจจุบัน
- **Database Ping**: ตรวจสอบการเชื่อมต่อฐานข้อมูล

## 📚 การใช้งาน Swagger

### เข้าถึง Swagger UI
เมื่อเซิร์ฟเวอร์ทำงานแล้ว สามารถเข้าถึงเอกสาร API ได้ที่:
- **Swagger UI**: http://localhost:8080/swagger/
- **JSON Spec**: http://localhost:8080/swagger/doc.json

### ฟีเจอร์ของ Swagger
- 📖 **เอกสารครบถ้วน**: รายละเอียด API ทั้งหมด
- 🧪 **ทดสอบ API**: ทดสอบ endpoint ได้ในหน้าเว็บ
- 🔐 **การยืนยันตัวตน**: รองรับการใส่ JWT Token
- 📋 **Model Schema**: แสดงโครงสร้างข้อมูลที่ชัดเจน

### การใช้งาน Authorization ใน Swagger
1. เข้าสู่ระบบผ่าน `/auth/login` เพื่อรับ JWT token
2. คลิกปุ่ม "Authorize" ใน Swagger UI
3. ใส่ token ในรูปแบบ: `Bearer your_jwt_token_here`
4. ทดสอบ protected endpoints ได้

## 🔧 การพัฒนา

### การเพิ่ม API Endpoint ใหม่

1. **สร้าง Model** ใน `models/`
2. **เพิ่ม Controller** ใน `controllers/`
3. **เพิ่ม Route** ใน `routes/routes.go`
4. **เพิ่ม Swagger Comments** สำหรับเอกสาร
5. **Generate Swagger**: `swag init -g main.go --output docs/`

### ตัวอย่างการเพิ่ม Endpoint

```go
// @Summary Get products
// @Description Get all products
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response
// @Router /products [get]
func (pc *ProductController) GetProducts(c *fiber.Ctx) error {
    // Implementation here
    return utils.SuccessResponse(c, "Success", products)
}
```

### การใช้งาน Hot Reload
```bash
# รันแบบ development mode
air

# ไฟล์จะถูก rebuild อัตโนมัติเมื่อมีการเปลี่ยนแปลง
```

### การ Debug
- ตรวจสอบ logs ในคอนโซล
- ใช้ `/api/v1/health` เพื่อตรวจสอบสถานะ
- ตรวจสอบไฟล์ `build-errors.log` หากมีข้อผิดพลาด

## 📊 Performance และ Optimization

### Database Connection Pooling
```go
db.SetMaxOpenConns(25)                 // จำนวน connection สูงสุด
db.SetMaxIdleConns(5)                  // จำนวน idle connection
db.SetConnMaxLifetime(5 * time.Minute) // อายุของ connection
```

### Middleware Pipeline
Middleware ทำงานตามลำดับ:
1. **Recover** - กู้คืนจาก panic
2. **Logger** - บันทึก request logs
3. **CORS** - จัดการ cross-origin requests
4. **JWT** - ตรวจสอบ authentication (สำหรับ protected routes)
5. **Admin** - ตรวจสอบสิทธิ์ admin (สำหรับ admin routes)

## 🧪 การทดสอบ

### ทดสอบด้วย cURL
```bash
# Health check
curl http://localhost:8080/api/v1/health

# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@test.com","password":"123456"}'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"123456"}'
```

### ทดสอบด้วย Postman
สามารถ import Swagger specification เข้า Postman ได้จาก:
http://localhost:8080/swagger/doc.json

## 🚀 Deployment

### Build สำหรับ Production
```bash
# Build optimized binary
go build -ldflags="-s -w" -o bin/main .

# หรือใช้ script ใน .air.toml
swag init -g main.go --output docs/ && go build -ldflags="-s -w" -o ./bin/main.exe .
```

### Environment Variables สำหรับ Production
```env
ENVIRONMENT=production
JWT_SECRET=very-strong-secret-key-for-production
DB_PASSWORD=strong-database-password
PORT=80
```

### การตั้งค่า CORS สำหรับ Production
```go
// แก้ไขใน main.go สำหรับ production
AllowOrigins: "https://yourdomain.com,https://api.yourdomain.com"
```

## 🤝 การมีส่วนร่วม

1. Fork โปรเจค
2. สร้าง feature branch (`git checkout -b feature/amazing-feature`)
3. Commit การเปลี่ยนแปลง (`git commit -m 'Add amazing feature'`)
4. Push ไปยัง branch (`git push origin feature/amazing-feature`)
5. สร้าง Pull Request

## 📄 ลิขสิทธิ์

โปรเจคนี้เป็น Template สำหรับการเรียนรู้ สามารถนำไปใช้ได้อย่างอิสระ

สร้างโดย: **Sing254463**

## 👨‍💻 ผู้พัฒนา

- **Sing254463** - *ผู้พัฒนาหลัก* - [GitHub](https://github.com/Sing254463)

## 🙏 กิตติกรรมประกาศ

- [Fiber Framework](https://gofiber.io/) - Web framework ที่รวดเร็ว
- [Go Community](https://golang.org/) - ภาษา Go และชุมชน
- [JWT.io](https://jwt.io/) - ข้อมูลเกี่ยวกับ JSON Web Tokens
- [Swagger](https://swagger.io/) - เครื่องมือสำหรับเอกสาร API
- [SQLX](https://github.com/jmoiron/sqlx) - Extensions สำหรับ database/sql
- [Validator](https://github.com/go-playground/validator) - ตรวจสอบความถูกต้องของข้อมูล
- [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) - เข้ารหัสรหัสผ่านแบบปลอดภัย
- [godotenv](https://github.com/joho/godotenv) - โหลดตัวแปรจากไฟล์ .env
- [Air](https://github.com/cosmtrek/air) - Hot reload สำหรับการพัฒนา Go
- [MySQL](https://www.mysql.com/) - ระบบจัดการฐานข้อมูล

---


**🌟 หากโปรเจคนี้มีประโยชน์ อย่าลืม Star ให้ด้วยนะครับ!**