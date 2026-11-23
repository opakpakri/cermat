import { driver } from 'driver.js';
import 'driver.js/dist/driver.css';
import { browser } from '$app/env';

const TOUR_SEEN_PREFIX = 'cermat_tour_seen_';

function isVisible(selector) {
	if (!browser) return false;
	const el = document.querySelector(selector);
	if (!el) return false;
	const rect = el.getBoundingClientRect();
	return rect.width > 0 && rect.height > 0;
}

function getNavStep() {
	if (isVisible('[data-tour="sidebar-nav"]')) {
		return {
			element: '[data-tour="sidebar-nav"]',
			popover: {
				title: 'Menu Navigasi Utama',
				description:
					'Gunakan menu di samping ini untuk berpindah cepat antara Dashboard, Pipeline, Kontak, Katalog, Aktivitas & Layanan, serta Pengaturan.',
				side: 'right',
				align: 'start'
			}
		};
	}
	if (isVisible('[data-tour="mobile-nav"]')) {
		return {
			element: '[data-tour="mobile-nav"]',
			popover: {
				title: 'Navigasi Bawah (Mobile)',
				description:
					'Gunakan bar navigasi di bawah layar untuk berpindah antar menu utama dengan mudah.',
				side: 'top',
				align: 'center'
			}
		};
	}
	return null;
}

function buildStepsForPath(pathname) {
	const steps = [];

	if (pathname === '/') {
		const navStep = getNavStep();
		if (navStep) steps.push(navStep);

		steps.push(
			{
				element: '[data-tour="dash-hero"]',
				popover: {
					title: 'Ringkasan Penjualan Bulan Ini',
					description:
						'Pantau 4 metrik utama workspace Anda: prospek baru bulan ini, estimasi nilai pipeline aktif, total closing/deal, dan tugas follow-up hari ini.',
					side: 'bottom',
					align: 'start'
				}
			},
			{
				element: '[data-tour="dash-quick"]',
				popover: {
					title: 'Pintasan Aksi Cepat',
					description:
						'Klik "Tambah Prospek / Deal Baru" untuk mencatat peluang penjualan baru, atau "Jadwal Follow-up Hari Ini" untuk memantau tugas harian.',
					side: 'bottom',
					align: 'start'
				}
			},
			{
				element: '[data-tour="dash-activity"]',
				popover: {
					title: 'Follow-up & Transaksi Terkini',
					description:
						'Tandai tugas follow-up yang sudah selesai langsung dari Dashboard dan pantau pergerakan transaksi terbaru Anda.',
					side: 'top',
					align: 'start'
				}
			}
		);
	} else if (pathname.startsWith('/pipeline')) {
		steps.push(
			{
				element: '[data-tour="pipeline-add-btn"]',
				popover: {
					title: 'Tambah Prospek / Deal Baru',
					description:
						'Buat transaksi baru dengan memilih kontak pelanggan dan item katalog (nilai transaksi otomatis terisi dari harga katalog).',
					side: 'bottom',
					align: 'end'
				}
			},
			{
				element: '[data-tour="pipeline-filter"]',
				popover: {
					title: 'Filter & Tampilan Kanban / Daftar',
					description:
						'Cari deal berdasarkan nama klien, saring berdasarkan sumber prospek, dan beralih antara tampilan Kanban Board atau Daftar Tabel.',
					side: 'bottom',
					align: 'start'
				}
			},
			{
				element: '[data-tour="pipeline-board"]',
				popover: {
					title: 'Kelola Tahapan Pipeline',
					description:
						'Pindahkan tahapan deal langsung dari dropdown di setiap kartu atau klik judul deal untuk memperbarui detail negosiasi.',
					side: 'top',
					align: 'start'
				}
			}
		);
	} else if (pathname.startsWith('/kontak')) {
		steps.push(
			{
				element: '[data-tour="contacts-add-btn"]',
				popover: {
					title: 'Tambah Kontak Baru',
					description:
						'Simpan data calon pembeli (leads) maupun pelanggan tetap lengkap dengan nomor WhatsApp, perusahaan, dan sumber prospek.',
					side: 'bottom',
					align: 'end'
				}
			},
			{
				element: '[data-tour="contacts-stats"]',
				popover: {
					title: 'Statistik Database Kontak',
					description:
						'Pantau total kontak, calon pembeli aktif, pelanggan yang sudah closing, dan prospek baru bulan ini.',
					side: 'bottom',
					align: 'start'
				}
			},
			{
				element: '[data-tour="contacts-search"]',
				popover: {
					title: 'Cari & Filter Kontak',
					description:
						'Temukan kontak dengan cepat berdasarkan nama, nomor telepon, perusahaan, atau statusnya.',
					side: 'bottom',
					align: 'start'
				}
			},
			{
				element: '[data-tour="contacts-table"]',
				popover: {
					title: 'Daftar Kontak & Pintasan WhatsApp',
					description:
						'Klik tombol "Chat WA" di bawah nomor telepon untuk langsung menghubungi klien via WhatsApp, atau gunakan tombol Edit/Hapus di kanan.',
					side: 'top',
					align: 'start'
				}
			}
		);
	} else if (pathname.startsWith('/katalog')) {
		steps.push(
			{
				element: '[data-tour="catalog-add-btn"]',
				popover: {
					title: 'Tambah Item Katalog',
					description:
						'Tambahkan unit properti atau produk/jasa baru beserta kode, harga dasar, dan spesifikasi khusus.',
					side: 'bottom',
					align: 'end'
				}
			},
			{
				element: '[data-tour="catalog-stats"]',
				popover: {
					title: 'Ringkasan Ketersediaan',
					description:
						'Lihat jumlah item yang masih tersedia (Available), terpesan (Booked), maupun terjual/nonaktif.',
					side: 'bottom',
					align: 'start'
				}
			},
			{
				element: '[data-tour="catalog-search"]',
				popover: {
					title: 'Cari & Saring Katalog',
					description:
						'Cari berdasarkan kode, nama, atau kategori, lalu filter berdasarkan status ketersediaan.',
					side: 'bottom',
					align: 'start'
				}
			},
			{
				element: '[data-tour="catalog-table"]',
				popover: {
					title: 'Tabel Katalog & Spesifikasi Lengkap',
					description:
						'Setiap atribut spesifikasi tampil pada kolomnya masing-masing. Anda juga dapat mengubah status ketersediaan langsung dari tabel.',
					side: 'top',
					align: 'start'
				}
			}
		);
	} else if (pathname.startsWith('/aktivitas')) {
		steps.push(
			{
				element: '[data-tour="activities-add-btn"]',
				popover: {
					title: 'Tambah Follow-up / Tiket Layanan',
					description:
						'Catat jadwal telepon, survei lokasi, atau buat tiket keluhan/layanan purnajual sesuai tab yang sedang aktif.',
					side: 'bottom',
					align: 'end'
				}
			},
			{
				element: '[data-tour="activities-tabs"]',
				popover: {
					title: 'Beralih Tab Follow-up & Tiket',
					description:
						'Angka pada tab menunjukkan jumlah jadwal follow-up yang masih tertunda dan tiket layanan yang belum ditutup (CLOSED).',
					side: 'bottom',
					align: 'start'
				}
			},
			{
				element: '[data-tour="activities-content"]',
				popover: {
					title: 'Tandai Selesai & Update Status',
					description:
						'Centang tugas follow-up yang sudah dilakukan atau perbarui status tiket layanan. Item yang sudah Selesai/CLOSED akan otomatis terkunci.',
					side: 'top',
					align: 'start'
				}
			}
		);
	} else if (pathname.startsWith('/pengguna')) {
		steps.push(
			{
				element: '[data-tour="users-add-btn"]',
				popover: {
					title: 'Tambah Akun Pengguna',
					description:
						'Khusus Admin: buat akun baru untuk pengguna (Tenant) atau sesama Admin beserta Mode Workspace-nya.',
					side: 'left',
					align: 'start'
				}
			},
			{
				element: '[data-tour="users-table"]',
				popover: {
					title: 'Daftar Akun Pengguna',
					description:
						'Pantau daftar akun terdaftar, edit profil/role/mode workspace, atau hapus akun yang sudah tidak aktif.',
					side: 'top',
					align: 'start'
				}
			}
		);
	} else if (pathname.startsWith('/pengaturan')) {
		steps.push(
			{
				element: '[data-tour="settings-theme"]',
				popover: {
					title: 'Pengaturan Tema Tampilan',
					description:
						'Pilih tema Terang, Gelap (Dark Mode), atau Sistem yang otomatis mengikuti pengaturan perangkat Anda.',
					side: 'bottom',
					align: 'start'
				}
			},
			{
				element: '[data-tour="settings-form"]',
				popover: {
					title: 'Profil Akun & Identitas Bisnis',
					description:
						'Perbarui Nama Bisnis, Nomor WhatsApp Bisnis, Alamat Kantor, atau ganti password akun Anda.',
					side: 'top',
					align: 'start'
				}
			}
		);
	}

	return steps.filter((s) => !s.element || isVisible(s.element));
}

export function startProductTour(pathname = '/') {
	if (!browser) return;
	const steps = buildStepsForPath(pathname);
	if (steps.length === 0) return;

	const drv = driver({
		showProgress: true,
		animate: true,
		allowClose: true,
		nextBtnText: 'Lanjut →',
		prevBtnText: '← Kembali',
		doneBtnText: 'Selesai',
		progressText: '{{current}} dari {{total}}',
		popoverClass: 'cermat-driver-popover',
		steps
	});

	drv.drive();
}

export function maybeStartFirstVisitTour(pathname = '/') {
	if (!browser || pathname !== '/') return;
	const key = `${TOUR_SEEN_PREFIX}dashboard`;
	if (localStorage.getItem(key)) return;
	localStorage.setItem(key, '1');
	setTimeout(() => {
		startProductTour('/');
	}, 500);
}
