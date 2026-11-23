<script>
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import { formatDateTime } from '$lib/utils/workspaceLabels';
	import Badge from '$lib/components/Badge.svelte';
	import SearchableSelect from '$lib/components/SearchableSelect.svelte';
	import { CalendarCheck2, LifeBuoy, Plus, Check, X } from 'lucide-svelte';

	let activeTab = $state('activities');
	let activities = $state([]);
	let tickets = $state([]);
	let contacts = $state([]);
	let deals = $state([]);
	let loading = $state(true);

	const pendingActivitiesCount = $derived(activities.filter((a) => !a.is_completed).length);
	const openTicketsCount = $derived(tickets.filter((t) => t.status !== 'CLOSED').length);

	let showActivityModal = $state(false);
	let actForm = $state({
		contact_id: '',
		deal_id: '',
		activity_type: 'CALL',
		title: '',
		description: '',
		scheduled_at: ''
	});

	let showTicketModal = $state(false);
	let ticketForm = $state({
		contact_id: '',
		deal_id: '',
		subject: '',
		description: '',
		priority: 'MEDIUM'
	});

	async function loadAll() {
		loading = true;
		try {
			const [actRes, tktRes, ctRes, dlRes] = await Promise.all([
				api.get('/activities'),
				api.get('/tickets'),
				api.get('/contacts'),
				api.get('/deals')
			]);
			activities = actRes;
			tickets = tktRes;
			contacts = ctRes;
			deals = dlRes;
		} catch (err) {
			console.error(err);
		} finally {
			loading = false;
		}
	}

	function openActivityModal() {
		actForm = {
			contact_id: '',
			deal_id: '',
			activity_type: '',
			title: '',
			description: '',
			scheduled_at: ''
		};
		showActivityModal = true;
	}

	function openTicketModal() {
		ticketForm = {
			contact_id: '',
			deal_id: '',
			subject: '',
			description: '',
			priority: ''
		};
		showTicketModal = true;
	}

	async function submitActivity(e) {
		e.preventDefault();
		if (!actForm.contact_id) {
			alert('Silakan pilih kontak terlebih dahulu');
			return;
		}
		if (!actForm.activity_type) {
			alert('Silakan pilih jenis aktivitas terlebih dahulu');
			return;
		}
		try {
			await api.post('/activities', {
				contact_id: Number(actForm.contact_id),
				deal_id: actForm.deal_id ? Number(actForm.deal_id) : null,
				activity_type: actForm.activity_type,
				title: actForm.title,
				description: actForm.description || null,
				scheduled_at: actForm.scheduled_at || null
			});
			showActivityModal = false;
			await loadAll();
		} catch (err) {
			alert(err.message);
		}
	}

	async function toggleActivity(id) {
		try {
			await api.patch(`/activities/${id}/toggle`, {});
			await loadAll();
		} catch (err) {
			alert(err.message);
		}
	}

	async function submitTicket(e) {
		e.preventDefault();
		if (!ticketForm.contact_id) {
			alert('Silakan pilih pelanggan terlebih dahulu');
			return;
		}
		if (!ticketForm.priority) {
			alert('Silakan pilih prioritas terlebih dahulu');
			return;
		}
		try {
			await api.post('/tickets', {
				contact_id: Number(ticketForm.contact_id),
				deal_id: ticketForm.deal_id ? Number(ticketForm.deal_id) : null,
				subject: ticketForm.subject,
				description: ticketForm.description || null,
				priority: ticketForm.priority
			});
			showTicketModal = false;
			await loadAll();
		} catch (err) {
			alert(err.message);
		}
	}

	async function updateTicketStatus(id, status) {
		try {
			await api.patch(`/tickets/${id}/status`, { status });
			await loadAll();
		} catch (err) {
			alert(err.message);
		}
	}

	onMount(loadAll);
</script>

<svelte:head>
	<title>Aktivitas & Layanan — Cermat</title>
</svelte:head>

<div class="space-y-6">
	<div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
		<div>
			<h1 class="text-2xl font-extrabold text-[#0F172A]">Aktivitas & Layanan</h1>
			<p class="text-sm text-[#64748B]">Jadwal follow-up prospek serta tiket keluhan/layanan purnajual.</p>
		</div>
		<div data-tour="activities-add-btn" class="w-full sm:w-auto">
			{#if activeTab === 'activities'}
				<button
					type="button"
					onclick={openActivityModal}
					class="inline-flex w-full items-center justify-center gap-2 rounded-xl bg-[#EA580C] px-4 py-2.5 text-sm font-semibold text-white shadow-sm transition hover:bg-[#C2410C] sm:w-auto"
				>
					<Plus class="h-4 w-4" />
					<span>Tambah Follow-up</span>
				</button>
			{:else}
				<button
					type="button"
					onclick={openTicketModal}
					class="inline-flex w-full items-center justify-center gap-2 rounded-xl bg-[#EA580C] px-4 py-2.5 text-sm font-semibold text-white shadow-sm transition hover:bg-[#C2410C] sm:w-auto"
				>
					<Plus class="h-4 w-4" />
					<span>Buat Tiket Layanan</span>
				</button>
			{/if}
		</div>
	</div>

	<!-- Segmented Tabs -->
	<div data-tour="activities-tabs" class="grid w-full grid-cols-2 gap-1.5 rounded-xl border border-[#E2E8F0] bg-slate-100/80 p-1">
		<button
			type="button"
			onclick={() => (activeTab = 'activities')}
			class="inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-lg px-3 py-2 text-xs font-semibold transition-all sm:px-4 sm:text-sm {activeTab === 'activities'
				? 'bg-white text-[#EA580C] shadow-sm'
				: 'text-[#64748B] hover:text-[#0F172A]'}"
		>
			<CalendarCheck2 class="h-4 w-4 shrink-0" />
			<span class="sm:hidden">Follow-up</span>
			<span class="hidden sm:inline">Jadwal Follow-up</span>
			<span class="rounded-full px-1.5 py-0.5 text-[11px] font-bold {activeTab === 'activities' ? 'bg-[#FFF7ED] text-[#C2410C]' : 'bg-slate-200/80 text-slate-600'}">
				{pendingActivitiesCount}
			</span>
		</button>
		<button
			type="button"
			onclick={() => (activeTab = 'tickets')}
			class="inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-lg px-3 py-2 text-xs font-semibold transition-all sm:px-4 sm:text-sm {activeTab === 'tickets'
				? 'bg-white text-[#EA580C] shadow-sm'
				: 'text-[#64748B] hover:text-[#0F172A]'}"
		>
			<LifeBuoy class="h-4 w-4 shrink-0" />
			<span class="sm:hidden">Tiket</span>
			<span class="hidden sm:inline">Tiket Layanan / Komplain</span>
			<span class="rounded-full px-1.5 py-0.5 text-[11px] font-bold {activeTab === 'tickets' ? 'bg-[#FFF7ED] text-[#C2410C]' : 'bg-slate-200/80 text-slate-600'}">
				{openTicketsCount}
			</span>
		</button>
	</div>

	<div data-tour="activities-content">
	{#if loading}
		<div class="bg-white border border-[#E2E8F0] rounded-xl p-8 text-center text-sm text-[#64748B]">Memuat data...</div>
	{:else if activeTab === 'activities'}
		<div class="bg-white border border-[#E2E8F0] rounded-xl shadow-sm divide-y divide-[#E2E8F0]">
			{#if activities.length === 0}
				<p class="p-8 text-center text-sm text-[#64748B]">Belum ada jadwal follow-up.</p>
			{:else}
				{#each activities as act}
					<div class="p-4 flex items-start justify-between gap-4 hover:bg-slate-50/70">
						<div class="flex items-start gap-3">
							<button
								type="button"
								disabled={act.is_completed}
								onclick={() => !act.is_completed && toggleActivity(act.id)}
								title={act.is_completed ? 'Sudah selesai' : 'Tandai selesai'}
								class="mt-0.5 w-5 h-5 rounded border flex items-center justify-center shrink-0 {act.is_completed
									? 'bg-[#15803D] border-[#15803D] text-white cursor-not-allowed opacity-90'
									: 'border-slate-300 bg-white hover:border-[#EA580C]'}"
							>
								{#if act.is_completed}
									<Check class="w-3.5 h-3.5" />
								{/if}
							</button>
							<div>
								<p class="text-sm font-semibold text-[#0F172A] {act.is_completed ? 'line-through text-[#64748B]' : ''}">
									{act.title}
								</p>
								<p class="text-xs text-[#64748B] mt-0.5">
									{act.contact_name} ({act.contact_phone}) • <span class="font-medium">{act.activity_type}</span>
									{#if act.deal_title} • Deal: {act.deal_title}{/if}
								</p>
								{#if act.description}
									<p class="text-xs text-[#64748B] mt-1">{act.description}</p>
								{/if}
								{#if act.scheduled_at}
									<p class="text-xs text-[#C2410C] font-medium mt-1">Jadwal: {formatDateTime(act.scheduled_at)}</p>
								{/if}
							</div>
						</div>
						<Badge variant={act.is_completed ? 'won' : 'warning'}>
							{act.is_completed ? 'Selesai' : 'Tertunda'}
						</Badge>
					</div>
				{/each}
			{/if}
		</div>
	{:else}
		<div class="overflow-hidden rounded-2xl border border-[#E2E8F0] bg-white shadow-sm">
			{#if tickets.length === 0}
				<p class="p-8 text-center text-sm text-[#64748B]">Belum ada tiket layanan atau komplain.</p>
			{:else}
				<!-- Mobile Card View -->
				<div class="divide-y divide-[#E2E8F0] md:hidden">
					{#each tickets as tk}
						<div class="space-y-2.5 p-4">
							<div class="flex items-center justify-between gap-2">
								<span class="font-mono text-xs text-[#64748B]">{tk.ticket_code}</span>
								<Badge variant={tk.priority === 'HIGH' ? 'high' : tk.priority === 'MEDIUM' ? 'warning' : 'default'}>
									{tk.priority}
								</Badge>
							</div>
							<div>
								<p class="text-sm font-semibold text-[#0F172A]">{tk.subject}</p>
								{#if tk.description}
									<p class="mt-0.5 text-xs text-[#64748B]">{tk.description}</p>
								{/if}
							</div>
							<div class="flex items-center justify-between gap-3 pt-1">
								<span class="truncate text-xs font-medium text-[#64748B]">{tk.contact_name}</span>
								<select
									value={tk.status}
									disabled={tk.status === 'CLOSED'}
									title={tk.status === 'CLOSED' ? 'Tiket sudah ditutup (CLOSED)' : 'Ubah status tiket'}
									onchange={(e) => tk.status !== 'CLOSED' && updateTicketStatus(tk.id, e.currentTarget.value)}
									class="shrink-0 rounded-lg border border-[#E2E8F0] px-2.5 py-1 text-xs font-medium focus:border-[#EA580C] focus:outline-none {tk.status === 'CLOSED'
										? 'cursor-not-allowed bg-slate-100 text-[#64748B] opacity-80'
										: 'bg-[#FAFAF9]'}"
								>
									<option value="OPEN">OPEN</option>
									<option value="IN_PROGRESS">IN_PROGRESS</option>
									<option value="RESOLVED">RESOLVED</option>
									<option value="CLOSED">CLOSED</option>
								</select>
							</div>
						</div>
					{/each}
				</div>

				<!-- Desktop Table View -->
				<div class="hidden overflow-x-auto md:block">
					<table class="w-full text-left text-sm">
						<thead class="border-b border-[#E2E8F0] bg-[#FAFAF9] text-xs font-semibold text-[#64748B]">
							<tr>
								<th class="px-5 py-3.5">Kode</th>
								<th class="px-5 py-3.5">Subjek Keluhan</th>
								<th class="px-5 py-3.5">Pelanggan</th>
								<th class="px-5 py-3.5">Prioritas</th>
								<th class="px-5 py-3.5">Status</th>
							</tr>
						</thead>
						<tbody class="divide-y divide-[#E2E8F0]">
							{#each tickets as tk}
								<tr class="hover:bg-slate-50/70">
									<td class="whitespace-nowrap px-5 py-4 font-mono text-xs text-[#64748B]">{tk.ticket_code}</td>
									<td class="px-5 py-4">
										<p class="font-semibold text-[#0F172A]">{tk.subject}</p>
										{#if tk.description}
											<p class="mt-0.5 text-xs text-[#64748B]">{tk.description}</p>
										{/if}
									</td>
									<td class="whitespace-nowrap px-5 py-4 text-[#0F172A]">{tk.contact_name}</td>
									<td class="whitespace-nowrap px-5 py-4">
										<Badge variant={tk.priority === 'HIGH' ? 'high' : tk.priority === 'MEDIUM' ? 'warning' : 'default'}>
											{tk.priority}
										</Badge>
									</td>
									<td class="whitespace-nowrap px-5 py-4">
										<select
											value={tk.status}
											disabled={tk.status === 'CLOSED'}
											title={tk.status === 'CLOSED' ? 'Tiket sudah ditutup (CLOSED)' : 'Ubah status tiket'}
											onchange={(e) => tk.status !== 'CLOSED' && updateTicketStatus(tk.id, e.currentTarget.value)}
											class="rounded-lg border border-[#E2E8F0] px-2.5 py-1 text-xs font-medium focus:border-[#EA580C] focus:outline-none {tk.status === 'CLOSED'
												? 'cursor-not-allowed bg-slate-100 text-[#64748B] opacity-80'
												: 'bg-[#FAFAF9]'}"
										>
											<option value="OPEN">OPEN</option>
											<option value="IN_PROGRESS">IN_PROGRESS</option>
											<option value="RESOLVED">RESOLVED</option>
											<option value="CLOSED">CLOSED</option>
										</select>
									</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			{/if}
		</div>
	{/if}
	</div>
</div>

<!-- Modal Tambah Follow-up -->
{#if showActivityModal}
	<div class="fixed inset-0 z-50 bg-slate-900/50 flex items-center justify-center p-4">
		<div class="bg-white rounded-2xl border border-[#E2E8F0] shadow-xl w-full max-w-md">
			<div class="px-6 py-4 border-b border-[#E2E8F0] flex items-center justify-between">
				<h2 class="text-lg font-bold text-[#0F172A]">Tambah Jadwal Follow-up</h2>
				<button type="button" onclick={() => (showActivityModal = false)} class="p-1.5 rounded-lg text-[#64748B] hover:bg-slate-100">
					<X class="w-5 h-5" />
				</button>
			</div>
			<form onsubmit={submitActivity} class="p-6 space-y-4">
				<div>
					<label for="act-contact" class="block text-sm font-medium text-[#0F172A] mb-1">Kontak *</label>
					<SearchableSelect
						bind:value={actForm.contact_id}
						options={contacts.map((ct) => ({ value: String(ct.id), label: `${ct.name} (${ct.phone})` }))}
						placeholder="Pilih Kontak"
						searchPlaceholder="Cari..."
						addHref="/kontak"
					/>
				</div>
				<div>
					<label for="act-deal" class="block text-sm font-medium text-[#0F172A] mb-1">Terkait Deal (Opsional)</label>
					<SearchableSelect
						bind:value={actForm.deal_id}
						options={[
							{ value: '', label: 'Pilih Terkait Deal (Opsional)' },
							...deals.map((dl) => ({ value: String(dl.id), label: `${dl.deal_code} — ${dl.title}` }))
						]}
						placeholder="Pilih Terkait Deal (Opsional)"
						searchPlaceholder="Cari..."
						addHref="/pipeline"
					/>
				</div>
				<div class="grid grid-cols-2 gap-3">
					<div>
						<label for="act-type" class="block text-sm font-medium text-[#0F172A] mb-1">Jenis Aktivitas *</label>
						<SearchableSelect
							bind:value={actForm.activity_type}
							options={[
								{ value: 'CALL', label: 'Telepon' },
								{ value: 'WHATSAPP', label: 'WhatsApp' },
								{ value: 'SITE_VISIT', label: 'Survei Lokasi' },
								{ value: 'MEETING', label: 'Meeting' },
								{ value: 'NOTE', label: 'Catatan' }
							]}
							placeholder="Pilih Jenis Aktivitas"
							searchable={false}
						/>
					</div>
					<div>
						<label for="act-date" class="block text-sm font-medium text-[#0F172A] mb-1">Jadwal</label>
						<input id="act-date" type="datetime-local" bind:value={actForm.scheduled_at} class="w-full px-3 py-2 rounded-lg border border-[#E2E8F0] text-sm" />
					</div>
				</div>
				<div>
					<label for="act-title" class="block text-sm font-medium text-[#0F172A] mb-1">Judul Tugas *</label>
					<input id="act-title" type="text" required bind:value={actForm.title} placeholder="Judul Tugas" class="w-full px-3.5 py-2 rounded-lg border border-[#E2E8F0] text-sm" />
				</div>
				<div>
					<label for="act-desc" class="block text-sm font-medium text-[#0F172A] mb-1">Catatan</label>
					<textarea id="act-desc" rows="2" bind:value={actForm.description} placeholder="Catatan (Opsional)" class="w-full px-3.5 py-2 rounded-lg border border-[#E2E8F0] text-sm"></textarea>
				</div>
				<div class="flex justify-end gap-2">
					<button type="button" onclick={() => (showActivityModal = false)} class="px-4 py-2 rounded-lg border border-[#E2E8F0] text-sm">Batal</button>
					<button type="submit" disabled={contacts.length === 0} class="px-5 py-2 rounded-lg bg-[#EA580C] text-white text-sm font-semibold">Simpan</button>
				</div>
			</form>
		</div>
	</div>
{/if}

<!-- Modal Buat Tiket -->
{#if showTicketModal}
	<div class="fixed inset-0 z-50 bg-slate-900/50 flex items-center justify-center p-4">
		<div class="bg-white rounded-2xl border border-[#E2E8F0] shadow-xl w-full max-w-md">
			<div class="px-6 py-4 border-b border-[#E2E8F0] flex items-center justify-between">
				<h2 class="text-lg font-bold text-[#0F172A]">Buat Tiket Layanan / Komplain</h2>
				<button type="button" onclick={() => (showTicketModal = false)} class="p-1.5 rounded-lg text-[#64748B] hover:bg-slate-100">
					<X class="w-5 h-5" />
				</button>
			</div>
			<form onsubmit={submitTicket} class="p-6 space-y-4">
				<div>
					<label for="tk-contact" class="block text-sm font-medium text-[#0F172A] mb-1">Pelanggan *</label>
					<SearchableSelect
						bind:value={ticketForm.contact_id}
						options={contacts.map((ct) => ({ value: String(ct.id), label: `${ct.name} (${ct.phone})` }))}
						placeholder="Pilih Pelanggan"
						searchPlaceholder="Cari..."
						addHref="/kontak"
					/>
				</div>
				<div>
					<label for="tk-priority" class="block text-sm font-medium text-[#0F172A] mb-1">Prioritas *</label>
					<SearchableSelect
						bind:value={ticketForm.priority}
						options={[
							{ value: 'LOW', label: 'LOW' },
							{ value: 'MEDIUM', label: 'MEDIUM' },
							{ value: 'HIGH', label: 'HIGH' }
						]}
						placeholder="Pilih Prioritas"
						searchable={false}
					/>
				</div>
				<div>
					<label for="tk-subject" class="block text-sm font-medium text-[#0F172A] mb-1">Subjek Keluhan *</label>
					<input id="tk-subject" type="text" required bind:value={ticketForm.subject} placeholder="Subjek Keluhan" class="w-full px-3.5 py-2 rounded-lg border border-[#E2E8F0] text-sm" />
				</div>
				<div>
					<label for="tk-desc" class="block text-sm font-medium text-[#0F172A] mb-1">Detail Keluhan</label>
					<textarea id="tk-desc" rows="3" bind:value={ticketForm.description} placeholder="Detail Keluhan (Opsional)" class="w-full px-3.5 py-2 rounded-lg border border-[#E2E8F0] text-sm"></textarea>
				</div>
				<div class="flex justify-end gap-2">
					<button type="button" onclick={() => (showTicketModal = false)} class="px-4 py-2 rounded-lg border border-[#E2E8F0] text-sm">Batal</button>
					<button type="submit" disabled={contacts.length === 0} class="px-5 py-2 rounded-lg bg-[#EA580C] text-white text-sm font-semibold">Buat Tiket</button>
				</div>
			</form>
		</div>
	</div>
{/if}
