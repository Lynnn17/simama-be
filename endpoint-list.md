# Daftar Endpoint Project

Dokumen ini merangkum endpoint HTTP yang terdaftar di project ini berdasarkan router dan handler yang benar-benar dipakai di codebase.

## Struktur Tabel

| Kolom | Isi |
|---|---|
| Method | HTTP method yang dipakai |
| Path | Path endpoint lengkap |
| Handler | Nama handler/method yang menangani request |
| Purpose | Fungsi endpoint |
| Auth | Status auth / middleware |
| Request DTO | Payload yang dikirim client |
| Response DTO | Bentuk data respons utama |
| Notes | Catatan tambahan |

## Endpoint Infrastruktur

| Method | Path | Handler | Purpose | Auth | Request DTO | Response DTO | Notes |
|---|---|---|---|---|---|---|---|
| GET | `/health` | `HTTP.HealthCheck` | Health check service | Public | - | - | Didefinisikan di `transport/http/http.go` |
| GET | `/swagger/*` | Swagger handler | Dokumentasi API | Public | - | - | Umumnya dev-only |
| ANY | `/socket.io/` | Socket.IO handler | Transport realtime | Public | - | - | Bukan endpoint REST biasa |

## Endpoint User

| Method | Path | Handler | Purpose | Auth | Request DTO | Response DTO | Notes |
|---|---|---|---|---|---|---|---|
| POST | `/v1/user/login` | `UserHandler.Login` | Login user | Public | `auth.InputLogin` | `auth.ResponseLogin` | Autentikasi dasar |
| POST | `/v1/user/login-web` | `UserHandler.LoginWeb` | Login khusus web | Public | `auth.InputLoginWeb` | `auth.ResponseLogin` | Jalur login alternatif |
| POST | `/v1/user/validasi-login` | `UserHandler.ValidasiLogin` | Validasi login | Public | - | - | Dipakai sebelum sesi aktif |
| GET | `/v1/user/captcha` | `UserHandler.GetCaptcha` | Generate captcha | Public | - | `captchaId` | - |
| GET | `/v1/user/captcha/image/{captchaId}` | `UserHandler.GetCaptchaImage` | Ambil gambar captcha | Public | Path param | Image | - |
| GET | `/v1/user/captcha/refresh/{captchaId}` | `UserHandler.RefreshCaptcha` | Refresh captcha | Public | Path param | - | - |
| POST | `/v1/user/` | `UserHandler.CreateUser` | Buat user | Protected | `auth.InputUser` | `auth.UserDTO` | `VerifyToken` |
| PUT | `/v1/user/{id}` | `UserHandler.UpdateUser` | Update user | Protected | `auth.InputUser` | `auth.UserDTO` | `VerifyToken` |
| PUT | `/v1/user/fcm-token/{id}` | `UserHandler.UpdateFcmToken` | Update FCM token | Protected | `auth.UserUpdateFcmTokenFormat` | - | `VerifyToken` |
| DELETE | `/v1/user/{id}` | `UserHandler.DeleteUser` | Hapus user | Protected | Path param | - | `VerifyToken` |
| PUT | `/v1/user/active-status/{id}` | `UserHandler.UpdateActiveStatus` | Ubah status aktif | Protected | Path param | - | `VerifyToken` |
| GET | `/v1/user/all` | `UserHandler.GetAllUsers` | List semua user | Protected | Query params | `[]auth.UserDTO` | `VerifyToken` |
| GET | `/v1/user/` | `UserHandler.GetUsers` | List user paginated | Protected | Query params | `[]auth.UserDTO` | `VerifyToken` |
| GET | `/v1/user/{id}` | `UserHandler.GetUserByID` | Detail user | Protected | Path param | `auth.UserDTO` | `VerifyToken` |
| PUT | `/v1/user/password/{id}` | `UserHandler.ChangePassword` | Ganti password | Protected | `auth.InputChangePassword` | - | `VerifyToken` |
| PUT | `/v1/user/password/pw/{id}` | `UserHandler.ChangePasswordByAdmin` | Reset password via admin | Protected | `auth.InputChangePassword` | - | `VerifyToken` |
| PUT | `/v1/user/password/reset/{id}` | `UserHandler.ResetPassword` | Reset password | Protected | Path param | - | `VerifyToken` |
| POST | `/v1/user/update-foto` | `UserHandler.UpdateFoto` | Update foto profil | Protected | `auth.UpdateFotoRequest` | - | `VerifyToken` |

## Endpoint Role

| Method | Path | Handler | Purpose | Auth | Request DTO | Response DTO | Notes |
|---|---|---|---|---|---|---|---|
| GET | `/v1/roles/` | `RoleHandler.GetRoles` | List role | Protected | Query params | `[]auth.Role` | `VerifyToken` |
| GET | `/v1/roles/all` | `RoleHandler.GetAllRoles` | List semua role | Protected | - | `[]auth.Role` | `VerifyToken` |
| GET | `/v1/roles/{id}` | `RoleHandler.GetRoleByID` | Detail role | Protected | Path param | `auth.Role` | `VerifyToken` |
| POST | `/v1/roles/` | `RoleHandler.CreateRole` | Buat role | Protected | `auth.RequestRole` | `auth.Role` | `VerifyToken` |
| PUT | `/v1/roles/{id}` | `RoleHandler.UpdateRole` | Update role | Protected | `auth.RequestRole` | `auth.Role` | `VerifyToken` |
| DELETE | `/v1/roles/{id}` | `RoleHandler.DeleteRole` | Hapus role | Protected | Path param | - | `VerifyToken` |

## Endpoint Menu & Permission

| Method | Path | Handler | Purpose | Auth | Request DTO | Response DTO | Notes |
|---|---|---|---|---|---|---|---|
| GET | `/v1/menu-role/` | `MenuHandler.GetMenuRole` | List permission menu per role | Protected | Query params | `[]auth.MenuResponse` | `VerifyToken` |
| GET | `/v1/menu-role/trx` | `MenuHandler.GetMenuRoleTrx` | Ambil data transaksi menu-role | Protected | Query params | `[]auth.MenuResponse` | `VerifyToken` |
| POST | `/v1/menu-role/bulk` | `MenuHandler.CreateBulkMenuRole` | Simpan permission banyak sekaligus | Protected | `auth.RequestBulkMenuRole` | - | `VerifyToken` |
| PUT | `/v1/menu-role/update-permission` | `MenuHandler.UpdateMenuPermission` | Update permission menu | Protected | `auth.RequestMenuPermissionFormat` | - | `VerifyToken` |
| GET | `/v1/menu/` | `MenuHandler.GetMenus` | List menu | Protected | Query params | `[]auth.MenuResponse` | `VerifyToken` |
| GET | `/v1/menu/all` | `MenuHandler.GetAllMenus` | List semua menu | Protected | - | `[]auth.MenuResponse` | `VerifyToken` |
| GET | `/v1/menu/{id}` | `MenuHandler.GetMenuByID` | Detail menu | Protected | Path param | `auth.MenuResponse` | `VerifyToken` |
| POST | `/v1/menu/` | `MenuHandler.CreateMenu` | Buat menu | Protected | `auth.RequestMenuFormat` | `auth.MenuResponse` | `VerifyToken` |
| PUT | `/v1/menu/{id}` | `MenuHandler.UpdateMenu` | Update menu | Protected | `auth.RequestMenuFormat` | `auth.MenuResponse` | `VerifyToken` |
| DELETE | `/v1/menu/{id}` | `MenuHandler.DeleteMenu` | Hapus menu | Protected | Path param | - | `VerifyToken` |

## Endpoint App Config

| Method | Path | Handler | Purpose | Auth | Request DTO | Response DTO | Notes |
|---|---|---|---|---|---|---|---|
| GET | `/v1/app-config/public/{id}` | `AppConfigHandler.GetPublicAppConfigByID` | Ambil config publik | Public | Path param | `auth.AppConfigDTOByID` | - |
| PUT | `/v1/app-config/{id}` | `AppConfigHandler.UpdateAppConfig` | Update config | Protected | `auth.RequestAppConfigFormat` | `auth.AppConfigDTO` | `VerifyToken` |
| GET | `/v1/app-config/{id}` | `AppConfigHandler.GetAppConfigByID` | Detail config | Protected | Path param | `auth.AppConfigDTOByID` | `VerifyToken` |
| POST | `/v1/app-config/upload` | `AppConfigHandler.UploadAppConfig` | Upload config/file | Protected | Multipart/form-data | - | `VerifyToken` |

> Status scope: **out of scope** untuk FR P1 / SRS yang sedang dipakai.

## Endpoint Master

### Academic Year

| Method | Path | Handler | Purpose | Auth | Request DTO | Response DTO | Notes |
|---|---|---|---|---|---|---|---|
| GET | `/v1/master/academic-year/` | `AcademicYearHandler.GetAcademicYears` | List academic year | Protected | Query params | `[]master.AcademicYearDTO` | `VerifyToken` |
| GET | `/v1/master/academic-year/all` | `AcademicYearHandler.GetAllAcademicYears` | List semua academic year | Protected | - | `[]master.AcademicYearDTO` | `VerifyToken` |
| POST | `/v1/master/academic-year/` | `AcademicYearHandler.CreateAcademicYear` | Buat academic year | Protected | `master.RequestAcademicYearFormat` | `master.AcademicYearDTO` | `VerifyToken` |
| POST | `/v1/master/academic-year/import/preview` | `AcademicYearHandler.PreviewImportAcademicYear` | Preview import Excel | Protected | File upload | `master.PreviewAcademicYearResult` | `VerifyToken` |
| POST | `/v1/master/academic-year/import` | `AcademicYearHandler.ImportAcademicYear` | Import dari preview | Protected | `master.ImportFromPreviewRequest` | - | `VerifyToken` |
| PUT | `/v1/master/academic-year/{id}` | `AcademicYearHandler.UpdateAcademicYear` | Update academic year | Protected | `master.RequestAcademicYearFormat` | `master.AcademicYearDTO` | `VerifyToken` |
| GET | `/v1/master/academic-year/{id}` | `AcademicYearHandler.GetAcademicYearByID` | Detail academic year | Protected | Path param | `master.AcademicYearDTO` | `VerifyToken` |
| DELETE | `/v1/master/academic-year/{id}` | `AcademicYearHandler.DeleteAcademicYear` | Hapus academic year | Protected | Path param | - | `VerifyToken` |

> Status scope: **out of scope** untuk FR P1 / SRS. Field seperti "Asal Perguruan Tinggi" dan "Program Studi" tetap berupa text field, bukan referensi master data.

### Company

| Method | Path | Handler | Purpose | Auth | Request DTO | Response DTO | Notes |
|---|---|---|---|---|---|---|---|
| GET | `/v1/master/company/` | `CompanyHandler.GetCompanies` | List company | Protected | Query params | `[]master.CompanyDTO` | `VerifyToken` |
| GET | `/v1/master/company/all` | `CompanyHandler.GetAllCompanies` | List semua company | Protected | - | `[]master.CompanyDTO` | `VerifyToken` |
| POST | `/v1/master/company/` | `CompanyHandler.CreateCompany` | Buat company | Protected | `master.RequestCompanyFormat` | `master.CompanyDTO` | `VerifyToken` |
| POST | `/v1/master/company/import/preview` | `CompanyHandler.PreviewImportCompany` | Preview import Excel | Protected | File upload | `master.PreviewCompanyResult` | `VerifyToken` |
| POST | `/v1/master/company/import` | `CompanyHandler.ImportCompany` | Import company | Protected | `master.ImportFromPreviewCompanyRequest` | - | `VerifyToken` |
| PUT | `/v1/master/company/{id}` | `CompanyHandler.UpdateCompany` | Update company | Protected | `master.RequestCompanyFormat` | `master.CompanyDTO` | `VerifyToken` |
| GET | `/v1/master/company/{id}` | `CompanyHandler.GetCompanyByID` | Detail company | Protected | Path param | `master.CompanyDTO` | `VerifyToken` |
| DELETE | `/v1/master/company/{id}` | `CompanyHandler.DeleteCompany` | Hapus company | Protected | Path param | - | `VerifyToken` |

> Status scope: **out of scope** untuk FR P1 / SRS.

### Department

| Method | Path | Handler | Purpose | Auth | Request DTO | Response DTO | Notes |
|---|---|---|---|---|---|---|---|
| GET | `/v1/master/department/` | `DepartmentHandler.GetDepartments` | List department | Protected | Query params | `[]master.DepartmentDTO` | `VerifyToken` |
| GET | `/v1/master/department/all` | `DepartmentHandler.GetAllDepartments` | List semua department | Protected | - | `[]master.DepartmentDTO` | `VerifyToken` |
| POST | `/v1/master/department/` | `DepartmentHandler.CreateDepartment` | Buat department | Protected | `master.RequestDepartmentFormat` | `master.DepartmentDTO` | `VerifyToken` |
| POST | `/v1/master/department/import/preview` | `DepartmentHandler.PreviewImportDepartment` | Preview import Excel | Protected | File upload | `master.PreviewDepartmentResult` | `VerifyToken` |
| POST | `/v1/master/department/import` | `DepartmentHandler.ImportDepartment` | Import department | Protected | `master.ImportFromPreviewDepartmentRequest` | - | `VerifyToken` |
| PUT | `/v1/master/department/{id}` | `DepartmentHandler.UpdateDepartment` | Update department | Protected | `master.RequestDepartmentFormat` | `master.DepartmentDTO` | `VerifyToken` |
| GET | `/v1/master/department/{id}` | `DepartmentHandler.GetDepartmentByID` | Detail department | Protected | Path param | `master.DepartmentDTO` | `VerifyToken` |
| DELETE | `/v1/master/department/{id}` | `DepartmentHandler.DeleteDepartment` | Hapus department | Protected | Path param | - | `VerifyToken` |

> Status scope: **out of scope** untuk FR P1 / SRS.

## Endpoint Internship

### Registration

| Method | Path | Handler | Purpose | Auth | Request DTO | Response DTO | Notes |
|---|---|---|---|---|---|---|---|
| GET | `/v1/internship/registration/` | `RegistrationHandler.GetRegistrations` | List pendaftaran | Protected | Query params | `[]internship.RegistrationDTO` | `VerifyToken` |
| GET | `/v1/internship/registration/all` | `RegistrationHandler.GetAllRegistrations` | List semua pendaftaran | Protected | - | `[]internship.RegistrationDTO` | `VerifyToken` |
| PUT | `/v1/internship/registration/status/{id}` | `RegistrationHandler.UpdateRegistrationStatus` | Ubah status pendaftaran | Protected | `internship.RequestUpdateRegistrationStatusFormat` | - | `VerifyToken` |
| POST | `/v1/internship/registration/` | `RegistrationHandler.CreateRegistration` | Buat pendaftaran | Protected | `internship.RequestRegistrationFormat` | `internship.RegistrationDTO` | `VerifyToken` |

### Mentor Assignment

| Method | Path | Handler | Purpose | Auth | Request DTO | Response DTO | Notes |
|---|---|---|---|---|---|---|---|
| POST | `/v1/internship/mentor-assignment/` | `MentorAssignmentHandler.CreateMentorAssignment` | Assign mentor | Protected | `internship.RequestMentorAssignmentFormat` | `internship.MentorAssignmentDTO` | `VerifyToken` |
| GET | `/v1/internship/mentor-assignment/mentor/{mentorId}/students` | `MentorAssignmentHandler.GetStudentsByMentor` | List student by mentor | Protected | Path param | `[]internship.MentorAssignmentDTO` | `VerifyToken` |
| GET | `/v1/internship/mentor-assignment/student/{studentId}` | `MentorAssignmentHandler.GetMentorByStudent` | Ambil mentor dari student | Protected | Path param | `internship.MentorAssignmentDTO` | `VerifyToken` |

### Logbook

| Method | Path | Handler | Purpose | Auth | Request DTO | Response DTO | Notes |
|---|---|---|---|---|---|---|---|
| POST | `/v1/internship/logbook/` | `LogbookHandler.CreateLogbook` | Submit logbook | Protected | `internship.RequestLogbookFormat` | `internship.LogbookDTO` | `VerifyToken` |
| GET | `/v1/internship/logbook/student/{studentId}` | `LogbookHandler.GetLogbooksByStudent` | Logbook student | Protected | Path param | `[]internship.LogbookDTO` | `VerifyToken` |
| GET | `/v1/internship/logbook/mentor/{mentorId}` | `LogbookHandler.GetLogbooksByMentor` | Logbook mentor | Protected | Path param | `[]internship.LogbookDTO` | `VerifyToken` |
| PUT | `/v1/internship/logbook/status/{id}` | `LogbookHandler.UpdateLogbookStatus` | Review / ubah status logbook | Protected | `internship.RequestUpdateLogbookStatusFormat` | - | `VerifyToken` |

### Task

| Method | Path | Handler | Purpose | Auth | Request DTO | Response DTO | Notes |
|---|---|---|---|---|---|---|---|
| POST | `/v1/internship/task/` | `TaskHandler.CreateTask` | Buat task | Protected | `internship.RequestTaskFormat` | `internship.TaskDTO` | `VerifyToken` |
| GET | `/v1/internship/task/student/{studentId}` | `TaskHandler.GetTasksByStudent` | List task student | Protected | Path param | `[]internship.TaskDTO` | `VerifyToken` |
| GET | `/v1/internship/task/mentor/{mentorId}` | `TaskHandler.GetTasksByMentor` | List task mentor | Protected | Path param | `[]internship.TaskDTO` | `VerifyToken` |
| POST | `/v1/internship/task/submit` | `TaskHandler.SubmitTask` | Submit file task | Protected | `internship.RequestSubmitTaskFileFormat` | `internship.TaskDTO` | `VerifyToken` |
| PUT | `/v1/internship/task/grade/{id}` | `TaskHandler.GradeTask` | Beri nilai task | Protected | `internship.RequestGradeTaskFormat` | `internship.TaskDTO` | `VerifyToken` |

## Endpoint File & Import Template

| Method | Path | Handler | Purpose | Auth | Request DTO | Response DTO | Notes |
|---|---|---|---|---|---|---|---|
| GET | `/v1/files/` | `FileHandler.GetFiles` | List file | Unclear | Query params | File metadata | Router tidak menambahkan middleware auth |
| POST | `/v1/files/upload` | `FileHandler.UploadFile` | Upload file | Unclear | Multipart/form-data | File metadata | `files.go` tidak memasang `VerifyToken` |
| GET | `/v1/files/image` | `FileHandler.GetImage` | Ambil file gambar | Unclear | Query params | Image | - |
| GET | `/v1/import/template/{domain}` | `ImportTemplateHandler.DownloadTemplate` | Download template import | Protected | Path param | File | `VerifyToken` |

## Catatan

- Semua route domain utama dimount di `/v1` dari `transport/http/router/router.go`.
- `files.go` adalah satu-satunya group yang terlihat tidak memasang middleware `VerifyToken` di router.
- Beberapa nama DTO di atas diambil dari file domain/handler yang terkait dan disusun sebagai referensi dokumentasi.
- Untuk antrean coding FR P1, jangan implementasikan `master/*` dan `app-config/*` karena sudah dinyatakan di luar cakupan SRS.
