# Simple Blog Platform API

RESTful API untuk platform blog sederhana menggunakan Golang, Gin Web Framework, GORM, dan MySQL 8.0 (Dockerized).

---

## 🚀 Fitur Utama

- **User Authentication**: Register & Login menggunakan `bcrypt` password hashing dan stateless `JWT Token`.
- **Blog Posts**: CRUD artikel blog (Create, Read All, Read Detail, Update, Delete) dengan proteksi otorisasi khusus pembuat artikel.
- **Comments**: Penambahan dan pencarian komentar pada artikel tertentu.
- **Clean Architecture**: Pemisahan layer yang jelas (`domain`, `repository`, `service`, `handler`, `middleware`, `router`).
- **Dockerized**: Siap dijalankan langsung menggunakan `docker-compose`.

---

## 🛠️ Prasyarat

- [Docker](https://www.docker.com/) & Docker Compose
- [DBeaver](https://dbeaver.io/) atau MySQL Client pilihanmu
- [Go 1.21+](https://golang.org/) *(opsional, hanya jika ingin running tanpa Docker)*

---

## ⚡ Cara Menjalankan Aplikasi

### 1. Clone Repository

```bash
git clone https://github.com/muzaim/go-simple-blog.git
cd go-simple-blog
```

---

### 2. Setup Environment Variables

Copy file `.env.example` di dalam folder `app` menjadi `.env`:

```bash
cp app/.env.example app/.env
```

Isi file `app/.env`:

```env
DB_USER=root
DB_PASS=abc123
DB_HOST=db
DB_PORT=3306
DB_NAME=appdb

PORT=8080
JWT_SECRET=super_secret_jwt_key_takehome
```

---

### 3. Eksekusi Skema Database (DBeaver / MySQL Client)

Buka DBeaver, hubungkan ke MySQL dengan konfigurasi:
- **Host**: `localhost` (atau `127.0.0.1`)
- **Port**: `3333`
- **User**: `root`
- **Password**: `abc123`
- **Database**: `appdb`

> **Catatan DBeaver**: Jika mengalami error `Public Key Retrieval is not allowed`, buka *Edit Connection* -> *Driver Properties* -> set `allowPublicKeyRetrieval=true` dan `useSSL=false`.

Jalankan query DDL berikut di SQL Editor untuk membuat tabel dan indeksnya:

```sql
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(150) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_users_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS posts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    author_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_posts_author FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_posts_author (author_id),
    INDEX idx_posts_created_at (created_at DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS comments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    post_id INT NOT NULL,
    author_name VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_comments_post FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    INDEX idx_comments_post (post_id),
    INDEX idx_comments_created_at (created_at DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

---

### 4. Jalankan Aplikasi dengan Docker

Jalankan perintah ini di root direktori project:

```bash
docker-compose up --build
```

Setelah log terminal menampilkan `Connected to database successfully!`, API siap di-hit di:
👉 `http://localhost:8080`

---

## 🧪 Pengujian API dengan Postman

File collection siap pakai sudah tersedia di root project: [`postman_collection.json`](postman_collection.json).

1. Buka aplikasi **Postman**.
2. Klik **Import** -> pilih file `postman_collection.json`.
3. Alur Pengujian:
   - Hit **Register User** (`POST /register`).
   - Hit **Login User** (`POST /login`) -> *Token JWT akan otomatis tersimpan di variabel collection*.
   - Hit **Create Post** (`POST /posts`).
   - Hit **Get All Posts** (`GET /posts`).
   - Hit **Add Comment** (`POST /posts/1/comments`).

---

## 📂 Struktur Folder Project

```text
.
├── app/
│   ├── main.go               # Entrypoint utama aplikasi
│   ├── .env                  # File konfigurasi environment (ignored by git)
│   ├── .env.example          # Template contoh environment
│   ├── internal/
│   │   ├── config/           # Koneksi DB & retry mechanism
│   │   ├── domain/           # Struct entity model & DTOs
│   │   ├── handler/          # HTTP Controllers
│   │   ├── middleware/       # Middleware JWT Auth
│   │   ├── repository/       # Layer Data Access (GORM)
│   │   ├── router/           # Registrasi endpoint routing
│   │   └── service/          # Layer Logika Bisnis & Validasi
│   └── pkg/
│       └── utils/            # Helper Bcrypt & JWT
├── Dockerfile
├── docker-compose.yml
├── postman_collection.json   # Collection Postman untuk testing
└── readme.md
```