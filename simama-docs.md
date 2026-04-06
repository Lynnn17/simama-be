# SYSTEM CONTEXT: SIMAMA (Sistem Informasi Manajemen Program Magang)

**Client:** PT Greatsoft Solusi Indonesia
**Tech Stack:** Golang (Backend), PostgreSQL (Database).

## 1. PROJECT OVERVIEW

SIMAMA memfasilitasi pendaftaran magang, penugasan mentor, logbook harian, dan manajemen tugas. Fokus pada MVP (Prioritas P1) yang divalidasi.

## 2. ACTORS & ROLES (Boilerplate Auth & JWT)

Sistem menggunakan boilerplate Golang eksisting untuk Autentikasi (`auth_user`, `auth_role`).

1. **Admin (HRD):** Kelola pendaftaran, validasi data, penugasan mentor.
2. **Mentor:** Pantau logbook, buat tugas, beri nilai.
3. **Mahasiswa:** Isi pendaftaran, isi logbook (wajib isi blockers), kumpul tugas.

## 3. DATABASE SCHEMA (ERD)

Catatan: Semua tabel berelasi (FK) ke `auth_user.id` (tipe data varchar 36).

- **`auth_user` & `auth_role`**: (Tabel eksisting boilerplate).
- **`internship_registrations`**: `id` (UUID, PK), `user_id` (varchar 36, FK), `university`, `major`, `cv_file_path`, `status` (pending, accepted, rejected), `created_at`.
- **`mentor_assignments`**: `id` (UUID, PK), `mentor_id` (varchar 36, FK), `student_id` (varchar 36, FK).
- **`logbooks`**: `id` (UUID, PK), `student_id` (varchar 36, FK), `log_date`, `activities`, `status`, `blockers`, `plan_tomorrow`, `evidence_url`, `submitted_at`.
- **`tasks`**: `id` (UUID, PK), `mentor_id` (varchar 36, FK), `student_id` (varchar 36, FK), `title`, `description`, `deadline`, `status` (assigned, submitted, graded), `grade`, `feedback`.
- **`task_files`**: `id` (UUID, PK), `task_id` (UUID, FK), `file_url`, `uploaded_by` (varchar 36, FK).

## 4. PROJECT STRUCTURE & ARCHITECTURE PATTERN

Sistem WAJIB menggunakan struktur folder eksisting (boilerplate) yang berada di dalam folder `internal/`. Setiap pembuatan fitur baru HANYA BOLEH meniru pola CRUD yang ada di `internal/domain/master` dan `internal/handlers`.

Aturan hierarki folder & penamaan file untuk entitas baru (misal: fitur `logbook` bisa dimasukkan ke folder `internal/domain/internship/` atau digabung ke `master`):

- `internal/domain/internship/logbook.model.go` -> Struct DB (GORM/SQL) & Struct Payload/DTO.
- `internal/domain/internship/logbook.repository.go` -> Interface & implementasi query database.
- `internal/domain/internship/logbook.service.go` -> Business logic.
- `internal/handlers/logbook.go` -> HTTP Handler (Echo/Gin/Fiber) untuk parsing request, memanggil service, dan return JSON.

## 5. STRICT RULES FOR AI AGENT

1. **Clone the Boilerplate Style:** Sebelum men-generate kode, pelajari gaya penulisan fungsi, interface, dan dependency injection yang ada di `internal/domain/master/company.repository.go` atau file lainnya. Tiru persis gaya tersebut.
2. **Strict File Naming:** Gunakan pemisah titik untuk layer di dalam domain (contoh: `nama_tabel.model.go`, `nama_tabel.service.go`) sesuai screenshot arsitektur yang diberikan pengguna.
3. **No New Architectures:** Jangan memperkenalkan struktur folder baru di luar `internal/domain/` dan `internal/handlers/`.
