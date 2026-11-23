# Cermat — Cellix Relationship Management Terpadu

**Cermat** adalah sistem CRM (*Customer Relationship Management*) multi-tenant yang mendukung dua mode bisnis sekaligus (**Properti / Real Estate** dan **Kantor / Jasa / Ritel**), dibangun dengan **Go (Fiber + PostgreSQL)** di backend dan **SvelteKit + Tailwind CSS** di frontend.

## Fitur Utama

- **Dashboard Ringkasan**: Statistik prospek, nilai pipeline aktif, kesepakatan closing (Won), dan jadwal follow-up terdekat.
- **Sales Pipeline (Kanban & Daftar)**: Pengelolaan prospek/deal lintas tahapan secara visual (Kanban) maupun tabel daftar dengan filter sumber prospek, pencarian, dan pagination.
- **Manajemen Kontak**: Pencatatan pelanggan/leads lengkap dengan status (`LEAD`, `PROSPECT`, `CUSTOMER`, `INACTIVE`), sumber prospek, dan riwayat deal.
- **Katalog Dual-Mode**:
  - Mode `PROPERTY`: Manajemen **Unit Properti / Kavling** beserta spesifikasi tipe, luas tanah/bangunan, blok, dan status ketersediaan (`AVAILABLE`, `BOOKED`, `SOLD`, `INACTIVE`).
  - Mode `OFFICE`: Manajemen **Produk & Layanan** beserta SKU, satuan, estimasi pengerjaan, dan status.
- **Aktivitas & Layanan**: Jadwal follow-up (`CALL`, `WHATSAPP`, `MEETING`, `SITE_VISIT`, `NOTE`) dan tiket keluhan pelanggan (`OPEN`, `IN_PROGRESS`, `RESOLVED`, `CLOSED`).
- **Sistem Admin & Manajemen Pengguna**: Hanya `ADMIN` yang dapat membuat, mengedit, dan menghapus akun pengguna (beserta `role` dan `business_mode`).
- **Tampilan Terang & Gelap**: Mendukung tema Terang, Gelap, dan Sistem serta panduan interaktif per halaman (*Product Tour*).

## Teknologi

- **Backend**: Go, Fiber v2, pgx v5 (PostgreSQL), JWT, bcrypt
- **Frontend**: SvelteKit (Svelte 5), Tailwind CSS, Lucide Icons, Driver.js

## Cara Menjalankan

### 1. Database (PostgreSQL)
Buat database PostgreSQL (misalnya `cermat_db`), lalu konfigurasi `backend/.env`:
```env
PORT=8080
DATABASE_URL=postgres://postgres:password@localhost:5432/cermat_db?sslmode=disable
JWT_SECRET=ganti-dengan-secret-key-anda
CORS_ORIGIN=*
```

Migrasi skema ([001_init_multi_tenant.sql](backend/migrations/001_init_multi_tenant.sql)) berjalan otomatis saat backend dijalankan. Jika ingin mengisi data contoh (opsional):
```bash
psql "$DATABASE_URL" -f backend/migrations/seed_dummy.sql
```

### 2. Backend
```bash
cd backend
go mod tidy
go run ./cmd/api
```
Backend akan berjalan di `http://localhost:8080` dan otomatis membuat akun Admin default jika belum ada:
- **Email**: `admin@cermat.id`
- **Password**: `admin123`

### 3. Frontend
```bash
cd frontend
npm install
npm run dev
```
Frontend akan berjalan di `http://localhost:5173`.
