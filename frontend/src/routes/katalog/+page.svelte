<script>
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import { authStore } from '$lib/stores/authStore';
	import { getWorkspaceLabels, formatRupiah } from '$lib/utils/workspaceLabels';
	import StatCard from '$lib/components/StatCard.svelte';
	import Badge from '$lib/components/Badge.svelte';
	import Pagination from '$lib/components/Pagination.svelte';
	import SearchableSelect from '$lib/components/SearchableSelect.svelte';
	import { CheckCircle2, Bookmark, Ban, Plus, Search, Pencil, Trash2, X } from 'lucide-svelte';

	let items = $state([]);
	let loading = $state(true);
	let search = $state('');
	let statusFilter = $state('');
	let currentPage = $state(1);
	const pageSize = 10;

	const paginatedItems = $derived(
		items.slice((currentPage - 1) * pageSize, currentPage * pageSize)
	);

	let showModal = $state(false);
	let editingItem = $state(null);
	let form = $state({
		code: '',
		name: '',
		category: '',
		price: 0,
		status: 'AVAILABLE',
		specs: {}
	});

	const user = $derived($authStore.user);
	const labels = $derived(getWorkspaceLabels(user?.business_mode));

	const availableCount = $derived(items.filter((i) => i.status === 'AVAILABLE').length);
	const bookedCount = $derived(items.filter((i) => i.status === 'BOOKED').length);
	const soldCount = $derived(items.filter((i) => i.status === 'SOLD' || i.status === 'INACTIVE').length);

	async function loadItems() {
		loading = true;
		try {
			const params = new URLSearchParams();
			if (search.trim()) params.set('search', search.trim());
			if (statusFilter) params.set('status', statusFilter);
			const qs = params.toString() ? `?${params.toString()}` : '';
			items = await api.get(`/catalog${qs}`);
		} catch (err) {
			console.error(err);
		} finally {
			loading = false;
		}
	}

	function openCreate() {
		editingItem = null;
		const initialSpecs = {};
		for (const f of labels.specFields) {
			initialSpecs[f.key] = '';
		}
		form = {
			code: '',
			name: '',
			category: '',
			price: 0,
			status: 'AVAILABLE',
			specs: initialSpecs
		};
		showModal = true;
	}

	function openEdit(item) {
		editingItem = item;
		const mergedSpecs = {};
		for (const f of labels.specFields) {
			mergedSpecs[f.key] = item.specs?.[f.key] ?? '';
		}
		form = {
			code: item.code || '',
			name: item.name,
			category: item.category || '',
			price: Number(item.price),
			status: item.status,
			specs: { ...(item.specs || {}), ...mergedSpecs }
		};
		showModal = true;
	}

	async function handleSubmit(e) {
		e.preventDefault();
		try {
			const payload = {
				code: form.code || null,
				name: form.name,
				category: form.category || null,
				price: Number(form.price || 0),
				status: form.status,
				specs: form.specs
			};
			if (editingItem) {
				await api.put(`/catalog/${editingItem.id}`, payload);
			} else {
				await api.post('/catalog', payload);
			}
			showModal = false;
			await loadItems();
		} catch (err) {
			alert(err.message);
		}
	}

	async function handleStatusChange(id, newStatus) {
		try {
			await api.patch(`/catalog/${id}/status`, { status: newStatus });
			await loadItems();
		} catch (err) {
			alert(err.message);
		}
	}

	async function handleDelete(id) {
		if (!confirm('Hapus item katalog ini?')) return;
		try {
			await api.delete(`/catalog/${id}`);
			await loadItems();
		} catch (err) {
			alert(err.message);
		}
	}

	onMount(loadItems);
</script>

<div class="space-y-6">
	<div class="flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
		<div>
			<h1 class="text-2xl font-extrabold text-[#0F172A]">{labels.catalogMenuLabel}</h1>
			<p class="text-sm text-[#64748B]">Kelola daftar {labels.catalogItemSingular.toLowerCase()} dan status ketersediaan.</p>
		</div>
		<button
			type="button"
			data-tour="catalog-add-btn"
			onclick={openCreate}
			class="inline-flex items-center justify-center gap-2 rounded-xl bg-[#EA580C] px-4 py-2.5 text-sm font-semibold text-white shadow-sm transition hover:bg-[#C2410C]"
		>
			<Plus class="h-4 w-4" />
			<span>Tambah {labels.catalogItemSingular}</span>
		</button>
	</div>

	<div data-tour="catalog-stats" class="grid grid-cols-2 gap-3 sm:grid-cols-3 sm:gap-4">
		<StatCard title="Tersedia (Available)" value={availableCount} icon={CheckCircle2} accent="green" />
		<StatCard title="Terpesan (Booked)" value={bookedCount} icon={Bookmark} accent="sky" />
		<StatCard title="Terjual / Nonaktif" value={soldCount} icon={Ban} accent="red" />
	</div>

	<form data-tour="catalog-search" onsubmit={(e) => { e.preventDefault(); currentPage = 1; loadItems(); }} class="flex flex-nowrap items-center gap-2">
		<div class="relative min-w-0 flex-1 sm:max-w-md">
			<Search class="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-[#64748B]" />
			<input
				type="text"
				bind:value={search}
				placeholder="Cari kode, nama, kategori..."
				class="w-full rounded-xl border border-[#E2E8F0] bg-white py-2.5 pl-9 pr-3 text-sm outline-none focus:border-[#EA580C] focus:ring-2 focus:ring-[#FFF7ED]"
			/>
		</div>
		<select
			bind:value={statusFilter}
			onchange={() => { currentPage = 1; loadItems(); }}
			class="max-w-[125px] shrink-0 rounded-xl border border-[#E2E8F0] bg-white px-2.5 py-2.5 text-sm text-[#0F172A] outline-none focus:border-[#EA580C] sm:max-w-none sm:px-3.5"
		>
			<option value="">Semua Status</option>
			<option value="AVAILABLE">Available</option>
			<option value="BOOKED">Booked</option>
			<option value="SOLD">Sold</option>
			<option value="INACTIVE">Inactive</option>
		</select>
		<button
			type="submit"
			class="shrink-0 rounded-xl border border-[#E2E8F0] bg-white px-3.5 py-2.5 text-sm font-semibold text-[#0F172A] hover:bg-slate-50 sm:px-4"
		>
			Cari
		</button>
	</form>

	<div data-tour="catalog-table" class="overflow-hidden rounded-2xl border border-[#E2E8F0] bg-white shadow-sm">
		<div class="overflow-x-auto">
			<table class="w-full text-left text-sm">
				<thead class="border-b border-[#E2E8F0] bg-[#FAFAF9] text-xs font-semibold text-[#64748B]">
					<tr>
						<th class="whitespace-nowrap px-5 py-3.5">Kode</th>
						<th class="min-w-[160px] px-5 py-3.5">Nama {labels.catalogItemSingular}</th>
						<th class="whitespace-nowrap px-5 py-3.5">Kategori</th>
						{#each labels.specFields as field}
							<th class="whitespace-nowrap px-5 py-3.5">{field.label}</th>
						{/each}
						<th class="whitespace-nowrap px-5 py-3.5">Harga</th>
						<th class="whitespace-nowrap px-5 py-3.5">Status</th>
						<th class="whitespace-nowrap px-5 py-3.5 text-right">Aksi</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-[#E2E8F0]">
					{#if loading}
						<tr><td colspan={6 + labels.specFields.length} class="px-5 py-8 text-center text-[#64748B]">Memuat katalog...</td></tr>
					{:else if items.length === 0}
						<tr><td colspan={6 + labels.specFields.length} class="px-5 py-8 text-center text-[#64748B]">Belum ada item katalog.</td></tr>
					{:else}
						{#each paginatedItems as item}
							<tr class="hover:bg-slate-50/70">
								<td class="whitespace-nowrap px-5 py-4 font-mono text-xs text-[#64748B]">{item.code || '-'}</td>
								<td class="px-5 py-4 font-semibold text-[#0F172A]">{item.name}</td>
								<td class="whitespace-nowrap px-5 py-4 text-[#64748B]">{item.category || '-'}</td>
								{#each labels.specFields as field}
									<td class="whitespace-nowrap px-5 py-4 text-[#64748B]">
										{item.specs?.[field.key] || '-'}
									</td>
								{/each}
								<td class="whitespace-nowrap px-5 py-4 font-semibold text-[#0F172A]">{formatRupiah(item.price)}</td>
								<td class="whitespace-nowrap px-5 py-4">
									<select
										value={item.status}
										onchange={(e) => handleStatusChange(item.id, e.currentTarget.value)}
										class="rounded-lg border border-[#E2E8F0] bg-[#FAFAF9] px-2 py-1 text-xs font-medium focus:border-[#EA580C] focus:outline-none"
									>
										<option value="AVAILABLE">AVAILABLE</option>
										<option value="BOOKED">BOOKED</option>
										<option value="SOLD">SOLD</option>
										<option value="INACTIVE">INACTIVE</option>
									</select>
								</td>
								<td class="whitespace-nowrap px-5 py-4 text-right">
									<div class="inline-flex items-center gap-1">
										<button
											type="button"
											onclick={() => openEdit(item)}
											class="rounded-lg p-1.5 text-[#64748B] hover:bg-[#FFF7ED] hover:text-[#EA580C]"
											title="Edit"
										>
											<Pencil class="h-4 w-4" />
										</button>
										<button
											type="button"
											onclick={() => handleDelete(item.id)}
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
		<Pagination bind:currentPage totalItems={items.length} {pageSize} />
	</div>
</div>

{#if showModal}
	<div class="fixed inset-0 z-50 bg-slate-900/50 flex items-center justify-center p-4">
		<div class="bg-white rounded-2xl border border-[#E2E8F0] shadow-xl w-full max-w-lg max-h-[90vh] overflow-y-auto">
			<div class="px-6 py-4 border-b border-[#E2E8F0] flex items-center justify-between">
				<h2 class="text-lg font-bold text-[#0F172A]">
					{editingItem ? `Edit ${labels.catalogItemSingular}` : `Tambah ${labels.catalogItemSingular}`}
				</h2>
				<button type="button" onclick={() => (showModal = false)} class="p-1.5 rounded-lg text-[#64748B] hover:bg-slate-100">
					<X class="w-5 h-5" />
				</button>
			</div>

			<form onsubmit={handleSubmit} class="p-6 space-y-4">
				<div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
					<div>
						<label for="cat-code" class="block text-sm font-medium text-[#0F172A] mb-1">Kode</label>
						<input
							id="cat-code"
							type="text"
							bind:value={form.code}
							placeholder="Kode"
							class="w-full px-3.5 py-2 rounded-lg border border-[#E2E8F0] text-sm focus:outline-none focus:border-[#EA580C]"
						/>
					</div>
					<div>
						<label for="cat-category" class="block text-sm font-medium text-[#0F172A] mb-1">Kategori</label>
						<input
							id="cat-category"
							type="text"
							bind:value={form.category}
							placeholder="Kategori"
							class="w-full px-3.5 py-2 rounded-lg border border-[#E2E8F0] text-sm focus:outline-none focus:border-[#EA580C]"
						/>
					</div>
				</div>

				<div>
					<label for="cat-name" class="block text-sm font-medium text-[#0F172A] mb-1">Nama {labels.catalogItemSingular} *</label>
					<input
						id="cat-name"
						type="text"
						required
						bind:value={form.name}
						placeholder={`Nama ${labels.catalogItemSingular}`}
						class="w-full px-3.5 py-2 rounded-lg border border-[#E2E8F0] text-sm focus:outline-none focus:border-[#EA580C]"
					/>
				</div>

				<div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
					<div>
						<label for="cat-price" class="block text-sm font-medium text-[#0F172A] mb-1">Harga Dasar (Rp) *</label>
						<input
							id="cat-price"
							type="number"
							min="0"
							required
							bind:value={form.price}
							class="w-full px-3.5 py-2 rounded-lg border border-[#E2E8F0] text-sm focus:outline-none focus:border-[#EA580C]"
						/>
					</div>
					<div>
						<label for="cat-status" class="block text-sm font-medium text-[#0F172A] mb-1">Status</label>
						<SearchableSelect
							bind:value={form.status}
							options={[
								{ value: 'AVAILABLE', label: 'AVAILABLE' },
								{ value: 'BOOKED', label: 'BOOKED' },
								{ value: 'SOLD', label: 'SOLD' },
								{ value: 'INACTIVE', label: 'INACTIVE' }
							]}
							placeholder="Pilih Status"
							searchable={false}
						/>
					</div>
				</div>

				<!-- Dynamic JSONB Specs -->
				<div class="pt-2 border-t border-[#E2E8F0]">
					<p class="text-xs font-semibold uppercase tracking-wider text-[#64748B] mb-3">
						Spesifikasi ({labels.modeLabel})
					</p>
					<div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
						{#each labels.specFields as field}
							<div>
								<label for="spec-{field.key}" class="block text-xs font-medium text-[#0F172A] mb-1">
									{field.label}
								</label>
								<input
									id="spec-{field.key}"
									type="text"
									bind:value={form.specs[field.key]}
									placeholder={field.label}
									class="w-full px-3 py-1.5 rounded-lg border border-[#E2E8F0] text-sm focus:outline-none focus:border-[#EA580C]"
								/>
							</div>
						{/each}
					</div>
				</div>

				<div class="flex justify-end gap-2.5 pt-2">
					<button
						type="button"
						onclick={() => (showModal = false)}
						class="px-4 py-2 rounded-lg border border-[#E2E8F0] text-sm font-medium text-[#64748B] hover:bg-slate-50"
					>
						Batal
					</button>
					<button type="submit" class="px-5 py-2 rounded-lg bg-[#EA580C] hover:bg-[#C2410C] text-white text-sm font-semibold">
						Simpan Item
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}
