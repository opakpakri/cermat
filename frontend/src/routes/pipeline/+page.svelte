<script>
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import { authStore } from '$lib/stores/authStore';
	import { getWorkspaceLabels, formatRupiah, formatDate, formatDateTime } from '$lib/utils/workspaceLabels';
	import KanbanBoard from '$lib/components/KanbanBoard.svelte';
	import Badge from '$lib/components/Badge.svelte';
	import Pagination from '$lib/components/Pagination.svelte';
	import SearchableSelect from '$lib/components/SearchableSelect.svelte';
	import { Plus, Search, LayoutGrid, List, X } from 'lucide-svelte';

	let deals = $state([]);
	let contacts = $state([]);
	let catalogItems = $state([]);
	let loading = $state(true);
	let viewMode = $state('kanban');
	let search = $state('');
	let leadSourceFilter = $state('');
	let currentPage = $state(1);
	const pageSize = 10;
	const totalItems = $derived(deals.length);
	const pagedDeals = $derived(deals.slice((currentPage - 1) * pageSize, currentPage * pageSize));

	let showModal = $state(false);
	let editingDeal = $state(null);
	let form = $state({
		title: '',
		contact_id: '',
		catalog_item_id: '',
		stage: 'NEW',
		deal_value: 0,
		payment_plan: '',
		expected_close_date: '',
		notes: ''
	});

	const user = $derived($authStore.user);
	const labels = $derived(getWorkspaceLabels(user?.business_mode));

	async function loadAll() {
		loading = true;
		try {
			const params = new URLSearchParams();
			if (search.trim()) params.set('search', search.trim());
			if (leadSourceFilter) params.set('lead_source', leadSourceFilter);
			const qs = params.toString() ? `?${params.toString()}` : '';

			const [dealsRes, contactsRes, catalogRes] = await Promise.all([
				api.get(`/deals${qs}`),
				api.get('/contacts'),
				api.get('/catalog')
			]);
			deals = dealsRes;
			contacts = contactsRes;
			catalogItems = catalogRes;
		} catch (err) {
			console.error(err);
		} finally {
			loading = false;
		}
	}

	function openCreate() {
		editingDeal = null;
		form = {
			title: '',
			contact_id: '',
			catalog_item_id: '',
			stage: 'NEW',
			deal_value: 0,
			payment_plan: labels.paymentPlans[0],
			expected_close_date: '',
			notes: ''
		};
		showModal = true;
	}

	function openEdit(deal) {
		editingDeal = deal;
		form = {
			title: deal.title,
			contact_id: String(deal.contact_id),
			catalog_item_id: deal.catalog_item_id ? String(deal.catalog_item_id) : '',
			stage: deal.stage,
			deal_value: deal.deal_value,
			payment_plan: deal.payment_plan || labels.paymentPlans[0],
			expected_close_date: deal.expected_close_date ? deal.expected_close_date.slice(0, 10) : '',
			notes: deal.notes || ''
		};
		showModal = true;
	}

	function handleCatalogSelect(id) {
		form.catalog_item_id = id;
		const item = catalogItems.find((c) => String(c.id) === id);
		if (item) {
			form.deal_value = Number(item.price);
			if (!form.title) {
				form.title = `Penawaran ${item.name}`;
			}
		}
	}

	async function handleSubmit(e) {
		e.preventDefault();
		if (!form.contact_id) {
			alert('Silakan pilih pelanggan terlebih dahulu');
			return;
		}
		try {
			const payload = {
				title: form.title,
				contact_id: Number(form.contact_id),
				catalog_item_id: form.catalog_item_id ? Number(form.catalog_item_id) : null,
				stage: form.stage,
				deal_value: Number(form.deal_value || 0),
				payment_plan: form.payment_plan || null,
				expected_close_date: form.expected_close_date || null,
				notes: form.notes || null
			};

			if (editingDeal) {
				await api.put(`/deals/${editingDeal.id}`, payload);
			} else {
				await api.post('/deals', payload);
			}
			showModal = false;
			await loadAll();
		} catch (err) {
			alert(err.message);
		}
	}

	async function handleStageChange(dealId, newStage) {
		try {
			await api.patch(`/deals/${dealId}/stage`, { stage: newStage });
			await loadAll();
		} catch (err) {
			alert(err.message);
		}
	}

	function stageLabel(key) {
		return labels.pipelineStages.find((s) => s.key === key)?.label || key;
	}

	onMount(loadAll);
</script>

<div class="space-y-6">
	<div class="flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
		<div>
			<h1 class="text-2xl font-extrabold text-[#0F172A]">Pipeline Penjualan</h1>
			<p class="text-sm text-[#64748B]">Pantau pergerakan prospek hingga closing ({labels.modeLabel}).</p>
		</div>

		<button
			type="button"
			data-tour="pipeline-add-btn"
			onclick={openCreate}
			class="inline-flex items-center justify-center gap-2 rounded-xl bg-[#EA580C] px-4 py-2.5 text-sm font-semibold text-white shadow-sm transition hover:bg-[#C2410C]"
		>
			<Plus class="h-4 w-4" />
			<span>Tambah Deal Baru</span>
		</button>
	</div>

	<!-- Filter & View Switcher -->
	<div data-tour="pipeline-filter" class="flex flex-col gap-2.5 rounded-2xl border border-[#E2E8F0] bg-white p-3.5 shadow-sm lg:flex-row lg:items-center lg:justify-between">
		<!-- Searchbar: full width on mobile, compact on desktop -->
		<div class="relative w-full lg:max-w-xs lg:flex-1">
			<Search class="pointer-events-none absolute left-3.5 top-1/2 h-4 w-4 -translate-y-1/2 text-[#64748B]" />
			<input
				type="text"
				bind:value={search}
				onkeydown={(e) => e.key === 'Enter' && loadAll()}
				placeholder="Cari nama klien atau judul deal..."
				class="w-full rounded-xl border border-[#E2E8F0] bg-[#FAFAF9] py-2.5 pl-10 pr-3.5 text-sm outline-none focus:border-[#EA580C] focus:bg-white focus:ring-2 focus:ring-[#FFF7ED]"
			/>
		</div>

		<!-- Row 2 on mobile (Kategori + Filter + Switch Kanban/Daftar sejajar) -->
		<div class="flex flex-nowrap items-center justify-between gap-2 lg:flex-1">
			<div class="flex min-w-0 flex-1 items-center gap-2">
				<select
					bind:value={leadSourceFilter}
					onchange={loadAll}
					class="min-w-0 flex-1 rounded-xl border border-[#E2E8F0] bg-white px-2.5 py-2 text-xs font-medium text-[#0F172A] outline-none focus:border-[#EA580C] sm:flex-none sm:px-3.5 sm:py-2.5 sm:text-sm"
				>
					<option value="">Semua Sumber</option>
					<option value="WHATSAPP">WhatsApp</option>
					<option value="INSTAGRAM">Instagram</option>
					<option value="WALK_IN">Walk-in</option>
					<option value="REFERRAL">Referensi</option>
					<option value="WEB">Website</option>
				</select>

				<button
					type="button"
					onclick={loadAll}
					class="shrink-0 rounded-xl border border-[#E2E8F0] bg-white px-3 py-2 text-xs font-semibold text-[#0F172A] hover:bg-slate-50 sm:px-4 sm:py-2.5 sm:text-sm"
				>
					Filter
				</button>
			</div>

			<div class="inline-flex shrink-0 rounded-xl border border-[#E2E8F0] bg-[#FAFAF9] p-0.5">
				<button
					type="button"
					onclick={() => (viewMode = 'kanban')}
					class="inline-flex items-center gap-1 rounded-lg px-2.5 py-1.5 text-xs font-semibold transition-colors sm:px-3 {viewMode === 'kanban'
						? 'bg-white text-[#EA580C] shadow-sm'
						: 'text-[#64748B]'}"
				>
					<LayoutGrid class="h-3.5 w-3.5" />
					<span>Kanban</span>
				</button>
				<button
					type="button"
					onclick={() => (viewMode = 'list')}
					class="inline-flex items-center gap-1 rounded-lg px-2.5 py-1.5 text-xs font-semibold transition-colors sm:px-3 {viewMode === 'list'
						? 'bg-white text-[#EA580C] shadow-sm'
						: 'text-[#64748B]'}"
				>
					<List class="h-3.5 w-3.5" />
					<span>Daftar</span>
				</button>
			</div>
		</div>
	</div>

	<div data-tour="pipeline-board">
	{#if loading}
		<div class="bg-white border border-[#E2E8F0] rounded-xl p-10 text-center text-sm text-[#64748B]">
			Memuat pipeline...
		</div>
	{:else if viewMode === 'kanban'}
		<KanbanBoard
			{deals}
			stages={labels.pipelineStages}
			onStageChange={handleStageChange}
			onSelectDeal={openEdit}
		/>
	{:else}
		<div class="bg-white border border-[#E2E8F0] rounded-xl overflow-hidden shadow-sm">
			<div class="overflow-x-auto">
				<table class="w-full min-w-[1100px] text-left text-sm">
					<thead class="bg-[#FAFAF9] border-b border-[#E2E8F0] text-xs uppercase text-[#64748B] whitespace-nowrap">
						<tr>
							<th class="px-4 py-3">Kode</th>
							<th class="px-4 py-3">Judul Deal</th>
							<th class="px-4 py-3">Pelanggan</th>
							<th class="px-4 py-3">{labels.catalogItemSingular}</th>
							<th class="px-4 py-3">Nilai</th>
							<th class="px-4 py-3">Tahapan</th>
							<th class="px-4 py-3">Target Closing</th>
							<th class="px-4 py-3">Dibuat</th>
							<th class="px-4 py-3">Diperbarui</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-[#E2E8F0]">
						{#each pagedDeals as deal}
							<tr class="hover:bg-slate-50/70">
								<td class="px-4 py-3 font-mono text-xs text-[#64748B] whitespace-nowrap">{deal.deal_code}</td>
								<td class="px-4 py-3 font-semibold text-[#0F172A] whitespace-nowrap">
									<button type="button" onclick={() => openEdit(deal)} class="hover:text-[#EA580C] text-left">
										{deal.title}
									</button>
								</td>
								<td class="px-4 py-3 text-[#0F172A] whitespace-nowrap">{deal.contact_name}</td>
								<td class="px-4 py-3 text-[#64748B] whitespace-nowrap">{deal.catalog_item_name || '-'}</td>
								<td class="px-4 py-3 font-semibold text-[#0F172A] whitespace-nowrap">{formatRupiah(deal.deal_value)}</td>
								<td class="px-4 py-3 whitespace-nowrap">
									<Badge variant={deal.stage === 'WON' ? 'won' : deal.stage === 'LOST' ? 'lost' : 'info'}>
										{stageLabel(deal.stage)}
									</Badge>
								</td>
								<td class="px-4 py-3 text-[#64748B] whitespace-nowrap">{formatDate(deal.expected_close_date)}</td>
								<td class="px-4 py-3 text-xs text-[#64748B] whitespace-nowrap">{formatDateTime(deal.created_at)}</td>
								<td class="px-4 py-3 text-xs text-[#64748B] whitespace-nowrap">{formatDateTime(deal.updated_at)}</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
			<Pagination bind:currentPage {totalItems} {pageSize} />
		</div>
	{/if}
	</div>
</div>

<!-- Modal Tambah / Edit Deal -->
{#if showModal}
	<div class="fixed inset-0 z-50 bg-slate-900/50 flex items-center justify-center p-4">
		<div class="bg-white rounded-2xl border border-[#E2E8F0] shadow-xl w-full max-w-lg max-h-[90vh] overflow-y-auto">
			<div class="px-6 py-4 border-b border-[#E2E8F0] flex items-center justify-between">
				<h2 class="text-lg font-bold text-[#0F172A]">{editingDeal ? 'Detail & Edit Deal' : 'Tambah Prospek / Deal Baru'}</h2>
				<button type="button" onclick={() => (showModal = false)} class="p-1.5 rounded-lg text-[#64748B] hover:bg-slate-100">
					<X class="w-5 h-5" />
				</button>
			</div>

			<form onsubmit={handleSubmit} class="p-6 space-y-4">
				<div>
					<label for="deal-contact" class="block text-sm font-medium text-[#0F172A] mb-1">Kontak Pelanggan *</label>
					<SearchableSelect
						bind:value={form.contact_id}
						options={contacts.map((ct) => ({ value: String(ct.id), label: `${ct.name} (${ct.phone})` }))}
						placeholder="Pilih pelanggan"
						searchPlaceholder="Cari..."
						addHref="/kontak"
					/>
				</div>

				<div>
					<label for="deal-catalog" class="block text-sm font-medium text-[#0F172A] mb-1">{labels.catalogItemSingular} (Opsional)</label>
					<SearchableSelect
						bind:value={form.catalog_item_id}
						options={[
							{ value: '', label: 'Tidak memilih item' },
							...catalogItems.map((item) => ({
								value: String(item.id),
								label: `${item.code ? `[${item.code}] ` : ''}${item.name} — ${formatRupiah(item.price)} (${item.status})`
							}))
						]}
						placeholder="Tidak memilih item"
						searchPlaceholder="Cari..."
						addHref="/katalog"
						onChange={handleCatalogSelect}
					/>
				</div>

				<div>
					<label for="deal-title" class="block text-sm font-medium text-[#0F172A] mb-1">Judul Deal *</label>
					<input
						id="deal-title"
						type="text"
						required
						bind:value={form.title}
						placeholder="Judul Deal"
						class="w-full px-3.5 py-2 rounded-lg border border-[#E2E8F0] text-sm focus:outline-none focus:border-[#EA580C]"
					/>
				</div>

				<div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
					<div>
						<label for="deal-value" class="block text-sm font-medium text-[#0F172A] mb-1">Estimasi Nilai (Rp)</label>
						<input
							id="deal-value"
							type="number"
							min="0"
							bind:value={form.deal_value}
							class="w-full px-3.5 py-2 rounded-lg border border-[#E2E8F0] text-sm focus:outline-none focus:border-[#EA580C]"
						/>
					</div>
					<div>
						<label for="deal-stage" class="block text-sm font-medium text-[#0F172A] mb-1">Tahapan</label>
						<SearchableSelect
							bind:value={form.stage}
							options={labels.pipelineStages.map((st) => ({ value: st.key, label: st.label }))}
							placeholder="Pilih Tahapan"
							searchable={false}
						/>
					</div>
				</div>

				<div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
					<div>
						<label for="deal-plan" class="block text-sm font-medium text-[#0F172A] mb-1">Rencana Pembayaran</label>
						<SearchableSelect
							bind:value={form.payment_plan}
							options={labels.paymentPlans.map((plan) => ({ value: plan, label: plan.replace('_', ' ') }))}
							placeholder="Pilih Rencana Pembayaran"
							searchable={false}
						/>
					</div>
					<div>
						<label for="deal-close" class="block text-sm font-medium text-[#0F172A] mb-1">Target Closing</label>
						<input
							id="deal-close"
							type="date"
							bind:value={form.expected_close_date}
							class="w-full px-3.5 py-2 rounded-lg border border-[#E2E8F0] text-sm focus:outline-none focus:border-[#EA580C]"
						/>
					</div>
				</div>

				<div>
					<label for="deal-notes" class="block text-sm font-medium text-[#0F172A] mb-1">Catatan Progres</label>
					<textarea
						id="deal-notes"
						rows="2"
						bind:value={form.notes}
						placeholder="Catatan Progress (Opsional)"
						class="w-full px-3.5 py-2 rounded-lg border border-[#E2E8F0] text-sm focus:outline-none focus:border-[#EA580C]"
					></textarea>
				</div>

				<div class="flex justify-end gap-2.5 pt-2">
					<button
						type="button"
						onclick={() => (showModal = false)}
						class="px-4 py-2 rounded-lg border border-[#E2E8F0] text-sm font-medium text-[#64748B] hover:bg-slate-50"
					>
						Batal
					</button>
					<button
						type="submit"
						disabled={contacts.length === 0}
						class="px-5 py-2 rounded-lg bg-[#EA580C] hover:bg-[#C2410C] text-white text-sm font-semibold disabled:opacity-50"
					>
						{editingDeal ? 'Simpan Perubahan' : 'Simpan Deal'}
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}
