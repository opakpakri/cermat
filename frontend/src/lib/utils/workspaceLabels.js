export function getWorkspaceLabels(businessMode) {
	const isProperty = businessMode === 'PROPERTY';

	return {
		mode: isProperty ? 'PROPERTY' : 'OFFICE',
		modeLabel: isProperty ? 'Properti' : 'Kantor Usaha',
		catalogMenuLabel: isProperty ? 'Unit Properti' : 'Produk & Layanan',
		catalogItemSingular: isProperty ? 'Unit / Kavling' : 'Produk / Jasa',
		catalogCodePlaceholder: isProperty ? 'Contoh: BLOK-A12' : 'Contoh: PKT-WEB-01',
		catalogCategoryPlaceholder: isProperty ? 'Contoh: Rumah Tipe 45 / Kavling' : 'Contoh: Jasa Konsultan / Paket SaaS',
		pipelineStages: isProperty
			? [
					{ key: 'NEW', label: 'Leads Masuk', badge: 'bg-[#FFEDD5] text-[#C2410C]' },
					{ key: 'QUALIFIED', label: 'Kualifikasi', badge: 'bg-[#E0F2FE] text-[#0369A1]' },
					{ key: 'PROPOSAL_SURVEY', label: 'Survei Lokasi', badge: 'bg-[#E0F2FE] text-[#0369A1]' },
					{ key: 'NEGOTIATION_BOOKING', label: 'Booking Fee', badge: 'bg-[#E0F2FE] text-[#0369A1]' },
					{ key: 'WON', label: 'Akad (Won)', badge: 'bg-[#DCFCE7] text-[#15803D]' },
					{ key: 'LOST', label: 'Batal (Lost)', badge: 'bg-[#FEE2E2] text-[#B91C1C]' }
				]
			: [
					{ key: 'NEW', label: 'Prospek Baru', badge: 'bg-[#FFEDD5] text-[#C2410C]' },
					{ key: 'QUALIFIED', label: 'Kualifikasi', badge: 'bg-[#E0F2FE] text-[#0369A1]' },
					{ key: 'PROPOSAL_SURVEY', label: 'Kirim Penawaran', badge: 'bg-[#E0F2FE] text-[#0369A1]' },
					{ key: 'NEGOTIATION_BOOKING', label: 'Negosiasi', badge: 'bg-[#E0F2FE] text-[#0369A1]' },
					{ key: 'WON', label: 'Deal / PO (Won)', badge: 'bg-[#DCFCE7] text-[#15803D]' },
					{ key: 'LOST', label: 'Batal (Lost)', badge: 'bg-[#FEE2E2] text-[#B91C1C]' }
				],
		specFields: isProperty
			? [
					{ key: 'blok_unit', label: 'Blok / No. Unit', placeholder: 'A-12' },
					{ key: 'tipe', label: 'Tipe', placeholder: '36/72' },
					{ key: 'luas_tanah', label: 'Luas Tanah (m²)', placeholder: '72' },
					{ key: 'luas_bangunan', label: 'Luas Bangunan (m²)', placeholder: '45' },
					{ key: 'hadap', label: 'Hadap', placeholder: 'Timur' }
				]
			: [
					{ key: 'sku', label: 'SKU / Kode Paket', placeholder: 'SVC-01' },
					{ key: 'satuan', label: 'Satuan', placeholder: 'Paket / Bulan / Unit' },
					{ key: 'durasi', label: 'Estimasi Pengerjaan', placeholder: '14 Hari Kerja' },
					{ key: 'deskripsi', label: 'Ringkasan Spesifikasi', placeholder: 'Termasuk maintenance 3 bulan' }
				],
		paymentPlans: isProperty
			? ['KPR', 'CASH_KERAS', 'CASH_BERTAHAP']
			: ['LUNAS', 'TERMIN', 'BERLANGGANAN']
	};
}

export function formatRupiah(value) {
	const num = Number(value || 0);
	return new Intl.NumberFormat('id-ID', {
		style: 'currency',
		currency: 'IDR',
		maximumFractionDigits: 0
	}).format(num);
}

export function formatDate(dateStr) {
	if (!dateStr) return '-';
	const d = new Date(dateStr);
	if (Number.isNaN(d.getTime())) return '-';
	return d.toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' });
}

export function formatDateTime(dateStr) {
	if (!dateStr) return '-';
	const d = new Date(dateStr);
	if (Number.isNaN(d.getTime())) return '-';
	return d.toLocaleString('id-ID', {
		day: '2-digit',
		month: 'short',
		year: 'numeric',
		hour: '2-digit',
		minute: '2-digit'
	});
}

export function formatTime(dateStr) {
	if (!dateStr) return '-';
	const d = new Date(dateStr);
	if (Number.isNaN(d.getTime())) return '-';
	return d.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' });
}
