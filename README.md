# Cermat — Cellix Relationship Management Terpadu

Sistem CRM multi-tenant dengan dua mode bisnis (**Properti / Real Estate** dan **Kantor / Jasa / Ritel**). Backend Go + PostgreSQL, frontend SvelteKit + Tailwind CSS.

```
cermat/
  backend/
    cmd/api/              entry point dan routing
    internal/             auth, kontak, pipeline, katalog, aktivitas, user admin
    migrations/           migrasi SQL (001_init_multi_tenant.sql, seed_dummy.sql)
  frontend/src/           aplikasi SvelteKit (dashboard, pipeline, kontak, katalog, aktivitas, pengguna)
```

## Fitur
- Dashboard: statistik prospek, nilai pipeline aktif, deal Won, follow-up terdekat.
- Sales pipeline: Kanban dan daftar (filter sumber, pencarian, pagination).
- Kontak: status `LEAD`, `PROSPECT`, `CUSTOMER`, `INACTIVE`, sumber prospek, riwayat deal.
- Katalog dual-mode: `PROPERTY` (unit/kavling) dan `OFFICE` (produk & layanan).
- Aktivitas follow-up (`CALL`, `WHATSAPP`, `MEETING`, `SITE_VISIT`, `NOTE`) dan tiket keluhan.
- Manajemen pengguna: hanya `ADMIN` yang dapat membuat, mengubah, dan menghapus akun (`role`, `business_mode`).
- Tema Terang/Gelap/Sistem dan panduan interaktif per halaman (Driver.js).

## Prasyarat
Go 1.22+, PostgreSQL 14+, Node.js 20+.

## Setup sekali jalan
1. Buat database: `CREATE DATABASE cermat_db;`
2. Di `backend`: buat `.env` (lihat Konfigurasi di bawah) dan sesuaikan `DATABASE_URL`.
3. Di `backend`: `go mod tidy`. Di `frontend`: `npm install`.
4. Migrasi skema berjalan otomatis saat backend pertama kali dijalankan, tidak perlu menjalankan tool migrasi manual.
5. Opsional, isi data contoh: `psql "$DATABASE_URL" -f backend/migrations/seed_dummy.sql`

## Jalankan
Gunakan dua terminal:
- `cd backend` lalu `go run ./cmd/api` (http://localhost:8080)
- `cd frontend` lalu `npm run dev` (http://localhost:5173)

Saat pertama kali jalan, backend membuat akun Admin default jika belum ada:
- Email: `admin@cermat.id`
- Password: `admin123`

Segera ganti password ini di lingkungan selain development.

## Konfigurasi
`backend/.env`:
```env
PORT=8080
DATABASE_URL=postgres://postgres:password@localhost:5432/cermat_db?sslmode=disable
JWT_SECRET=ganti-dengan-secret-key-anda
CORS_ORIGIN=*
```
Untuk produksi, ganti `JWT_SECRET` dengan nilai acak yang kuat dan batasi `CORS_ORIGIN` ke domain frontend.

## Tests
<!-- Sesuaikan jika repo punya test -->
- `cd backend` lalu `go test ./...` menjalankan unit test.
- Jika ada integration test, gunakan database sekali pakai (throwaway) terpisah, bukan database utama.

## Teknologi
- Backend: Go, Fiber v2, pgx v5, JWT, bcrypt
- Frontend: SvelteKit (Svelte 5), Tailwind CSS, Lucide Icons, Driver.js

## Produksi
<!-- Sesuaikan jika repo punya docker-compose -->
Jalankan `npm run build` di `frontend`, lalu sajikan hasil build lewat web server (mis. Caddy/Nginx) dan proxy `/api` ke backend. Pastikan `JWT_SECRET` dan `CORS_ORIGIN` sudah diatur untuk produksi.
