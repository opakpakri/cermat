<script>
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import { formatDate } from '$lib/utils/workspaceLabels';
	import Badge from '$lib/components/Badge.svelte';
	import StatCard from '$lib/components/StatCard.svelte';
	import Pagination from '$lib/components/Pagination.svelte';
	import SearchableSelect from '$lib/components/SearchableSelect.svelte';
	import { Users, UserPlus, UserCheck, Calendar, Search, Plus, MessageCircle, Pencil, Trash2, Phone, X } from 'lucide-svelte';

	let contacts = $state([]);
	let loading = $state(true);
	let search = $state('');
	let statusFilter = $state('');
	let currentPage = $state(1);
	const pageSize = 10;

	const totalContacts = $derived(contacts.length);
	const leadsCount = $derived(contacts.filter((c) => c.status === 'LEAD' || c.status === 'PROSPECT').length);
	const customersCount = $derived(contacts.filter((c) => c.status === 'CUSTOMER').length);
	const thisMonthCount = $derived(
		contacts.filter((c) => {
			const d = new Date(c.created_at);
			const now = new Date();
			return d.getMonth() === now.getMonth() && d.getFullYear() === now.getFullYear();
		}).length
	);

	const paginatedContacts = $derived(
		contacts.slice((currentPage - 1) * pageSize, currentPage * pageSize)
	);

	let showModal = $state(false);
	let editingContact = $state(null);
	let form = $state({
		name: '',
		phone: '',
		email: '',
		company_name: '',
		address: '',
		lead_source: 'WHATSAPP',
		status: 'LEAD',
		notes: ''
	});

	async function loadContacts() {
		loading = true;
		try {
			const params = new URLSearchParams();
			if (search.trim()) params.set('search', search.trim());
			if (statusFilter) params.set('status', statusFilter);
			const qs = params.toString() ? `?${params.toString()}` : '';
			contacts = await api.get(`/contacts${qs}`);
			currentPage = 1;
		} catch (err) {
			console.error(err);
		} finally {
			loading = false;
		}
	}

	function openCreate() {
		editingContact = null;
		form = {
			name: '',
			phone: '',
			email: '',
			company_name: '',
			address: '',
			lead_source: '',
			status: 'LEAD',
			notes: ''
		};
		showModal = true;
	}

	function openEdit(ct) {
		editingContact = ct;
		form = {
			name: ct.name,
			phone: ct.phone,
			email: ct.email || '',
			company_name: ct.company_name || '',
			address: ct.address || '',
			lead_source: ct.lead_source || 'WHATSAPP',
			status: ct.status || 'LEAD',
			notes: ct.notes || ''
		};
		showModal = true;
	}

	async function handleSubmit(e) {
		e.preventDefault();
		if (!form.lead_source) {
			alert('Silakan pilih sumber prospek terlebih dahulu');
			return;
		}
		try {
			const payload = {
				...form,
				email: form.email || null,
				company_name: form.company_name || null,
				address: form.address || null,
				notes: form.notes || null
			};
			if (editingContact) {
				await api.put(`/contacts/${editingContact.id}`, payload);
			} else {
				await api.post('/contacts', payload);
			}
			showModal = false;
			await loadContacts();
		} catch (err) {
			alert(err.message);
		}
	}

	async function handleDelete(id) {
		if (!confirm('Hapus kontak ini?')) return;
		try {
			await api.delete(`/contacts/${id}`);
			await loadContacts();
		} catch (err) {
			alert(err.message);
		}
	}

	function waLink(phone) {
		const cleaned = String(phone || '').replace(/[^0-9]/g, '');
		const normalized = cleaned.startsWith('0') ? '62' + cleaned.slice(1) : cleaned;
		return `https://wa.me/${normalized}`;
	}

	onMount(loadContacts);
</script>

<div class="space-y-6">
	<div class="flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
		<div>
			<h1 class="text-2xl font-extrabold text-[#0F172A]">Daftar Kontak</h1>
			<p class="text-sm text-[#64748B]">Kelola data kontak, prospek, dan pelanggan aktif Anda.</p>
		</div>
		<button
			type="button"
			data-tour="contacts-add-btn"
			onclick={openCreate}
			class="inline-flex items-center justify-center gap-2 rounded-xl bg-[#EA580C] px-4 py-2.5 text-sm font-semibold text-white shadow-sm transition hover:bg-[#C2410C]"
		>
			<Plus class="h-4 w-4" />
			<span>Tambah Kontak</span>
		</button>
	</div>

	<!-- 4 Stat Cards -->
	<div data-tour="contacts-stats" class="grid grid-cols-2 gap-3 sm:gap-4 lg:grid-cols-4">
		<StatCard title="Total Kontak" value={totalContacts} icon={Users} accent="orange" />
		<StatCard title="Calon Pembeli" value={leadsCount} icon={UserPlus} accent="sky" />
		<StatCard title="Pelanggan Aktif" value={customersCount} icon={UserCheck} accent="green" />
		<StatCard title="Prospek Bulan Ini" value={thisMonthCount} icon={Calendar} accent="slate" />
	</div>

	<!-- Pencarian & Filter -->
	<form data-tour="contacts-search" onsubmit={(e) => { e.preventDefault(); loadContacts(); }} class="flex flex-nowrap items-center gap-2">
		<div class="relative min-w-0 flex-1 sm:max-w-md">
			<Search class="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-[#64748B]" />
			<input
				type="text"
				bind:value={search}
				placeholder="Cari nama, telepon, perusahaan..."
				class="w-full rounded-xl border border-[#E2E8F0] bg-white py-2.5 pl-9 pr-3 text-sm outline-none focus:border-[#EA580C] focus:ring-2 focus:ring-[#FFF7ED]"
			/>
		</div>
		<select
			bind:value={statusFilter}
			onchange={loadContacts}
			class="max-w-[125px] shrink-0 rounded-xl border border-[#E2E8F0] bg-white px-2.5 py-2.5 text-sm text-[#0F172A] outline-none focus:border-[#EA580C] sm:max-w-none sm:px-3.5"
		>
			<option value="">Semua Status</option>
			<option value="LEAD">Lead</option>
			<option value="PROSPECT">Prospect</option>
			<option value="CUSTOMER">Customer</option>
			<option value="INACTIVE">Inactive</option>
		</select>
		<button
			type="submit"
			class="shrink-0 rounded-xl border border-[#E2E8F0] bg-white px-3.5 py-2.5 text-sm font-semibold text-[#0F172A] hover:bg-slate-50 sm:px-4"
		>
			Cari
		</button>
	</form>

	<!-- Tabel Kontak -->
	<div data-tour="contacts-table" class="overflow-hidden rounded-2xl border border-[#E2E8F0] bg-white shadow-sm">
		<div class="overflow-x-auto">
			<table class="w-full text-left text-sm">
				<thead class="border-b border-[#E2E8F0] bg-[#FAFAF9] text-xs font-semibold text-[#64748B]">
					<tr>
						<th class="px-5 py-3.5">Nama</th>
						<th class="px-5 py-3.5">Email</th>
						<th class="whitespace-nowrap px-5 py-3.5">Telepon</th>
						<th class="px-5 py-3.5">Perusahaan</th>
						<th class="min-w-[180px] px-5 py-3.5">Perusahaan</th>
						<th class="whitespace-nowrap px-5 py-3.5">Sumber</th>
						<th class="whitespace-nowrap px-5 py-3.5">Status</th>
						<th class="whitespace-nowrap px-5 py-3.5">Tanggal Dibuat</th>
						<th class="whitespace-nowrap px-5 py-3.5 text-right">Aksi</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-[#E2E8F0]">
					{#if loading}
						<tr><td colspan="8" class="px-5 py-8 text-center text-[#64748B]">Memuat kontak...</td></tr>
					{:else if contacts.length === 0}
						<tr><td colspan="8" class="px-5 py-8 text-center text-[#64748B]">Belum ada kontak ditemukan.</td></tr>
					{:else}
						{#each paginatedContacts as ct}
							<tr class="hover:bg-slate-50/70">
								<td class="px-5 py-4 font-semibold text-[#0F172A]">
									{ct.name}
								</td>
								<td class="px-5 py-4 text-[#64748B]">
									{ct.email || '-'}
								</td>
								<td class="whitespace-nowrap px-5 py-4 text-[#0F172A]">
									<div class="flex flex-col items-start gap-1.5">
										<span class="inline-flex items-center gap-1.5">
											<Phone class="h-3.5 w-3.5 text-[#64748B]" />
											<span>{ct.phone}</span>
										</span>
										<a
											href={waLink(ct.phone)}
											target="_blank"
											rel="noopener noreferrer"
											class="inline-flex items-center gap-1 rounded-md bg-[#DCFCE7] px-2 py-0.5 text-xs font-medium text-[#15803D] hover:opacity-80"
											title="Chat WhatsApp"
										>
											<MessageCircle class="h-3 w-3" />
											<span>Chat WA</span>
										</a>
									</div>
								</td>
								<td class="px-5 py-4 text-[#64748B]">{ct.company_name || '-'}</td>
								<td class="min-w-[180px] px-5 py-4 text-[#64748B]">{ct.company_name || '-'}</td>
								<td class="whitespace-nowrap px-5 py-4 text-xs text-[#64748B]">{ct.lead_source}</td>
								<td class="whitespace-nowrap px-5 py-4">
									<Badge variant={ct.status}>{ct.status}</Badge>
								</td>
								<td class="whitespace-nowrap px-5 py-4 text-xs text-[#64748B]">{formatDate(ct.created_at)}</td>
								<td class="whitespace-nowrap px-5 py-4 text-right">
									<div class="inline-flex items-center gap-1">
										<button
											type="button"
											onclick={() => openEdit(ct)}
											class="rounded-lg p-1.5 text-[#64748B] hover:bg-[#FFF7ED] hover:text-[#EA580C]"
											title="Edit"
										>
											<Pencil class="h-4 w-4" />
										</button>
										<button
											type="button"
											onclick={() => handleDelete(ct.id)}
											class="rounded-lg p-1.5 text-[#64748B] hover:bg-[#FEE2E2] hover:text-[#EF4444]"
											title="Hapus"
										>
											<Trash2 class="h-4 w-4" />
										</button>
									</div>
								</td>
							</tr>
						{/each}
					{/if}
				</tbody>
			</table>
		</div>
		<Pagination bind:currentPage totalItems={contacts.length} {pageSize} />
	</div>
</div>

{#if showModal}
	<div class="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/50 p-4">
		<div class="max-h-[90vh] w-full max-w-lg overflow-y-auto rounded-2xl border border-[#E2E8F0] bg-white shadow-xl">
			<div class="flex items-center justify-between border-b border-[#E2E8F0] px-6 py-4">
				<h2 class="text-lg font-bold text-[#0F172A]">{editingContact ? 'Edit Kontak' : 'Tambah Kontak Baru'}</h2>
				<button type="button" onclick={() => (showModal = false)} class="rounded-lg p-1.5 text-[#64748B] hover:bg-slate-100">
					<X class="h-5 w-5" />
				</button>
			</div>

			<form onsubmit={handleSubmit} class="space-y-4 p-6">
				<div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
					<div>
						<label for="ct-name" class="mb-1 block text-sm font-medium text-[#0F172A]">Nama Lengkap *</label>
						<input
							id="ct-name"
							type="text"
							required
							bind:value={form.name}
							placeholder="Nama Lengkap"
							class="w-full rounded-xl border border-[#E2E8F0] px-3.5 py-2 text-sm outline-none focus:border-[#EA580C]"
						/>
					</div>
					<div>
						<label for="ct-phone" class="mb-1 block text-sm font-medium text-[#0F172A]">No. Telepon / WA *</label>
						<input
							id="ct-phone"
							type="text"
							required
							bind:value={form.phone}
							placeholder="081234567890"
							class="w-full rounded-xl border border-[#E2E8F0] px-3.5 py-2 text-sm outline-none focus:border-[#EA580C]"
						/>
					</div>
				</div>

				<div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
					<div>
						<label for="ct-email" class="mb-1 block text-sm font-medium text-[#0F172A]">Email</label>
						<input
							id="ct-email"
							type="email"
							bind:value={form.email}
							placeholder="Email"
							class="w-full rounded-xl border border-[#E2E8F0] px-3.5 py-2 text-sm outline-none focus:border-[#EA580C]"
						/>
					</div>
					<div>
						<label for="ct-company" class="mb-1 block text-sm font-medium text-[#0F172A]">Perusahaan / Instansi</label>
						<input
							id="ct-company"
							type="text"
							bind:value={form.company_name}
							placeholder="Perusahaan/Instansi"
							class="w-full rounded-xl border border-[#E2E8F0] px-3.5 py-2 text-sm outline-none focus:border-[#EA580C]"
						/>
					</div>
				</div>

				<div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
					<div>
						<label for="ct-source" class="mb-1 block text-sm font-medium text-[#0F172A]">Sumber Prospek *</label>
						<SearchableSelect
							bind:value={form.lead_source}
							options={[
								{ value: 'WHATSAPP', label: 'WhatsApp' },
								{ value: 'INSTAGRAM', label: 'Instagram' },
								{ value: 'WALK_IN', label: 'Walk-in' },
								{ value: 'REFERRAL', label: 'Referensi' },
								{ value: 'WEB', label: 'Website' }
							]}
							placeholder="Pilih Sumber Prospek"
							searchable={false}
						/>
					</div>
					<div>
						<label for="ct-status" class="mb-1 block text-sm font-medium text-[#0F172A]">Status Kontak</label>
						<SearchableSelect
							bind:value={form.status}
							options={[
								{ value: 'LEAD', label: 'Lead' },
								{ value: 'PROSPECT', label: 'Prospect' },
								{ value: 'CUSTOMER', label: 'Customer' },
								{ value: 'INACTIVE', label: 'Inactive' }
							]}
							placeholder="Pilih Status Kontak"
							searchable={false}
						/>
					</div>
				</div>

				<div>
					<label for="ct-address" class="mb-1 block text-sm font-medium text-[#0F172A]">Alamat / Domisili</label>
					<input
						id="ct-address"
						type="text"
						bind:value={form.address}
						placeholder="Alamat/Domisili"
						class="w-full rounded-xl border border-[#E2E8F0] px-3.5 py-2 text-sm outline-none focus:border-[#EA580C]"
					/>
				</div>

				<div>
					<label for="ct-notes" class="mb-1 block text-sm font-medium text-[#0F172A]">Catatan Preferensi</label>
					<textarea
						id="ct-notes"
						rows="2"
						bind:value={form.notes}
						placeholder="Catatan Referensi (Opsional)"
						class="w-full rounded-xl border border-[#E2E8F0] px-3.5 py-2 text-sm outline-none focus:border-[#EA580C]"
					></textarea>
				</div>

				<div class="flex justify-end gap-2.5 pt-2">
					<button
						type="button"
						onclick={() => (showModal = false)}
						class="rounded-xl border border-[#E2E8F0] px-4 py-2 text-sm font-medium text-[#64748B] hover:bg-slate-50"
					>
						Batal
					</button>
					<button type="submit" class="rounded-xl bg-[#EA580C] px-5 py-2 text-sm font-semibold text-white hover:bg-[#C2410C]">
						Simpan Kontak
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}
