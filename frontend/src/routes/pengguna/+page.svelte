<script>
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { api } from '$lib/api';
	import { authStore } from '$lib/stores/authStore';
	import { formatDate } from '$lib/utils/workspaceLabels';
	import Badge from '$lib/components/Badge.svelte';
	import Pagination from '$lib/components/Pagination.svelte';
	import SearchableSelect from '$lib/components/SearchableSelect.svelte';
	import {
		Plus,
		Search,
		Pencil,
		Trash2,
		X,
		ShieldCheck,
		User as UserIcon,
		AlertTriangle
	} from 'lucide-svelte';

	let users = $state([]);
	let loading = $state(true);
	let search = $state('');
	let currentPage = $state(1);
	const pageSize = 10;

	let modalOpen = $state(false);
	let editing = $state(null);
	let saving = $state(false);
	let errorMsg = $state('');

	let form = $state({
		name: '',
		email: '',
		password: '',
		role: 'USER',
		business_name: '',
		business_mode: 'PROPERTY',
		phone: '',
		address: ''
	});

	let deleteTarget = $state(null);
	let deleting = $state(false);

	const roleOptions = [
		{ value: 'USER', label: 'User (Pemilik Bisnis)' },
		{ value: 'ADMIN', label: 'Admin (Pengelola Sistem)' }
	];

	const modeOptions = [
		{ value: 'PROPERTY', label: 'Properti / Real Estate' },
		{ value: 'OFFICE', label: 'Kantor / Jasa / Ritel' }
	];

	async function load() {
		loading = true;
		try {
			const qs = search.trim() ? `?search=${encodeURIComponent(search.trim())}` : '';
			const res = await api.get(`/admin/users${qs}`);
			users = Array.isArray(res) ? res : [];
			currentPage = 1;
		} catch (e) {
			users = [];
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		if ($authStore.user && $authStore.user.role !== 'ADMIN') {
			goto('/');
			return;
		}
		load();
	});

	const totalItems = $derived(users.length);
	const pagedUsers = $derived(users.slice((currentPage - 1) * pageSize, currentPage * pageSize));

	function openAdd() {
		editing = null;
		errorMsg = '';
		form = {
			name: '',
			email: '',
			password: '',
			role: 'USER',
			business_name: '',
			business_mode: 'PROPERTY',
			phone: '',
			address: ''
		};
		modalOpen = true;
	}

	function openEdit(u) {
		editing = u;
		errorMsg = '';
		form = {
			name: u.name || '',
			email: u.email || '',
			password: '',
			role: u.role || 'USER',
			business_name: u.business_name || '',
			business_mode: u.business_mode || 'PROPERTY',
			phone: u.phone || '',
			address: u.address || ''
		};
		modalOpen = true;
	}

	async function handleSave(e) {
		e.preventDefault();
		errorMsg = '';
		saving = true;
		try {
			const payload = {
				name: form.name.trim(),
				email: form.email.trim(),
				role: form.role,
				business_name: form.business_name.trim(),
				business_mode: form.business_mode,
				phone: form.phone.trim() ? form.phone.trim() : null,
				address: form.address.trim() ? form.address.trim() : null
			};
			if (!editing || form.password.trim()) {
				payload.password = form.password;
			}
			if (editing) {
				await api.put(`/admin/users/${editing.id}`, payload);
			} else {
				await api.post('/admin/users', payload);
			}
			modalOpen = false;
			await load();
		} catch (err) {
			errorMsg = err?.message || 'Gagal menyimpan data pengguna.';
		} finally {
			saving = false;
		}
	}

	async function confirmDelete() {
		if (!deleteTarget) return;
		deleting = true;
		try {
			await api.delete(`/admin/users/${deleteTarget.id}`);
			deleteTarget = null;
			await load();
		} catch (err) {
			alert(err?.message || 'Gagal menghapus pengguna.');
		} finally {
			deleting = false;
		}
	}
</script>

<svelte:head>
	<title>Manajemen Pengguna — Cermat</title>
</svelte:head>

<div class="space-y-5">
	<!-- Header -->
	<div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
		<div>
			<h1 class="text-xl font-bold text-[#0F172A] sm:text-2xl">
				Manajemen Pengguna
			</h1>
			<p class="text-xs text-[#64748B]">
				Menu khusus Admin untuk membuat dan mengelola akun bisnis pengguna.
			</p>
		</div>
		<button
			type="button"
			data-tour="users-add-btn"
			onclick={openAdd}
			class="inline-flex items-center justify-center gap-2 rounded-xl bg-[#EA580C] px-4 py-2.5 text-sm font-semibold text-white shadow-sm transition hover:bg-[#C2410C]"
		>
			<Plus class="h-4 w-4" />
			<span>Tambah Pengguna</span>
		</button>
	</div>

	<!-- Search Bar -->
	<form
		onsubmit={(e) => {
			e.preventDefault();
			load();
		}}
		class="flex flex-wrap items-center gap-2"
	>
		<div class="relative flex-1 min-w-[220px] max-w-md">
			<Search class="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-[#94A3B8]" />
			<input
				type="text"
				bind:value={search}
				placeholder="Cari nama, email, atau nama bisnis..."
				class="w-full rounded-xl border border-[#E2E8F0] bg-white py-2.5 pl-9 pr-3 text-sm text-[#0F172A] focus:border-[#EA580C] focus:outline-none"
			/>
		</div>
		<button
			type="submit"
			class="rounded-xl border border-[#E2E8F0] bg-white px-4 py-2.5 text-sm font-semibold text-[#0F172A] hover:bg-slate-50"
		>
			Cari
		</button>
	</form>

	<!-- Users Table -->
	<div
		data-tour="users-table"
		class="overflow-hidden rounded-2xl border border-[#E2E8F0] bg-white shadow-sm"
	>
		{#if loading}
			<div class="space-y-3 p-5">
				{#each Array(4) as _}
					<div class="h-12 animate-pulse rounded-xl bg-slate-100"></div>
				{/each}
			</div>
		{:else if users.length === 0}
			<div class="py-12 text-center">
				<p class="text-sm font-semibold text-[#0F172A]">
					Belum ada pengguna ditemukan
				</p>
				<p class="mt-1 text-xs text-[#64748B]">
					Klik tombol "Tambah Pengguna" untuk membuat akun baru.
				</p>
			</div>
		{:else}
			<div class="overflow-x-auto">
				<table class="w-full text-left text-sm">
					<thead class="border-b border-[#E2E8F0] bg-slate-50 text-[11px] font-bold uppercase tracking-wider text-[#64748B]">
						<tr>
							<th class="px-4 py-3">Nama</th>
							<th class="px-4 py-3">Email</th>
							<th class="px-4 py-3">Nama Bisnis</th>
							<th class="px-4 py-3">Mode</th>
							<th class="px-4 py-3">Role</th>
							<th class="px-4 py-3">Dibuat</th>
							<th class="px-4 py-3 text-right">Aksi</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-[#E2E8F0]">
						{#each pagedUsers as u}
							<tr class="hover:bg-slate-50">
								<td class="px-4 py-3 font-semibold text-[#0F172A]">
									{u.name}
								</td>
								<td class="px-4 py-3 text-xs text-[#64748B]">
									{u.email}
								</td>
								<td class="px-4 py-3 text-xs font-medium text-[#0F172A]">
									{u.business_name}
								</td>
								<td class="px-4 py-3">
									<Badge variant={u.business_mode === 'PROPERTY' ? 'info' : 'purple'}>
										{u.business_mode === 'PROPERTY' ? 'Properti' : 'Kantor/Jasa'}
									</Badge>
								</td>
								<td class="px-4 py-3">
									<Badge variant={u.role === 'ADMIN' ? 'orange' : 'slate'}>
										{#if u.role === 'ADMIN'}
											<ShieldCheck class="mr-1 inline h-3 w-3" />ADMIN
										{:else}
											<UserIcon class="mr-1 inline h-3 w-3" />USER
										{/if}
									</Badge>
								</td>
								<td class="px-4 py-3 text-xs text-[#64748B]">
									{formatDate(u.created_at)}
								</td>
								<td class="px-4 py-3 text-right">
									<div class="inline-flex items-center gap-1">
										<button
											type="button"
											onclick={() => openEdit(u)}
											title="Edit Pengguna"
											class="rounded-lg p-1.5 text-[#64748B] hover:bg-slate-100 hover:text-[#EA580C]"
										>
											<Pencil class="h-4 w-4" />
										</button>
										{#if u.id !== $authStore.user?.id}
											<button
												type="button"
												onclick={() => (deleteTarget = u)}
												title="Hapus Pengguna"
												class="rounded-lg p-1.5 text-[#64748B] hover:bg-slate-100 hover:text-red-600"
											>
												<Trash2 class="h-4 w-4" />
											</button>
										{/if}
									</div>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
			<Pagination bind:currentPage {totalItems} {pageSize} />
		{/if}
	</div>
</div>

<!-- Add/Edit Modal -->
{#if modalOpen}
	<div class="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/50 p-4 backdrop-blur-xs">
		<div class="w-full max-w-lg rounded-2xl border border-[#E2E8F0] bg-white p-6 shadow-xl">
			<div class="mb-4 flex items-center justify-between">
				<h2 class="text-lg font-bold text-[#0F172A]">
					{editing ? 'Edit Akun Pengguna' : 'Tambah Pengguna Baru'}
				</h2>
				<button
					type="button"
					onclick={() => (modalOpen = false)}
					class="rounded-lg p-1 text-[#64748B] hover:bg-slate-100"
				>
					<X class="h-5 w-5" />
				</button>
			</div>

			{#if errorMsg}
				<div class="mb-4 rounded-xl bg-[#FEE2E2] px-3 py-2 text-xs font-medium text-[#B91C1C]">
					{errorMsg}
				</div>
			{/if}

			<form onsubmit={handleSave} class="space-y-3.5">
				<div class="grid gap-3.5 sm:grid-cols-2">
					<div>
						<label for="user-name" class="mb-1 block text-xs font-semibold text-[#0F172A]">
							Nama Lengkap <span class="text-red-500">*</span>
						</label>
						<input
							id="user-name"
							type="text"
							required
							bind:value={form.name}
							placeholder="Nama Lengkap"
							class="w-full rounded-xl border border-[#E2E8F0] bg-white px-3 py-2 text-sm text-[#0F172A] focus:border-[#EA580C] focus:outline-none"
						/>
					</div>
					<div>
						<span class="mb-1 block text-xs font-semibold text-[#0F172A]">
							Role <span class="text-red-500">*</span>
						</span>
						<SearchableSelect
							options={roleOptions}
							bind:value={form.role}
							searchable={false}
							placeholder="Pilih Role"
						/>
					</div>
				</div>

				<div class="grid gap-3.5 sm:grid-cols-2">
					<div>
						<label for="user-email" class="mb-1 block text-xs font-semibold text-[#0F172A]">
							Email Login <span class="text-red-500">*</span>
						</label>
						<input
							id="user-email"
							type="email"
							required
							bind:value={form.email}
							placeholder="Email Pengguna"
							class="w-full rounded-xl border border-[#E2E8F0] bg-white px-3 py-2 text-sm text-[#0F172A] focus:border-[#EA580C] focus:outline-none"
						/>
					</div>
					<div>
						<label for="user-biz" class="mb-1 block text-xs font-semibold text-[#0F172A]">
							Nama Bisnis <span class="text-red-500">*</span>
						</label>
						<input
							id="user-biz"
							type="text"
							required
							bind:value={form.business_name}
							placeholder="Nama Bisnis / Perusahaan"
							class="w-full rounded-xl border border-[#E2E8F0] bg-white px-3 py-2 text-sm text-[#0F172A] focus:border-[#EA580C] focus:outline-none"
						/>
					</div>
				</div>

				<div class="grid gap-3.5 sm:grid-cols-2">
					<div>
						<span class="mb-1 block text-xs font-semibold text-[#0F172A]">
							Mode Workspace <span class="text-red-500">*</span>
						</span>
						<SearchableSelect
							options={modeOptions}
							bind:value={form.business_mode}
							searchable={false}
							placeholder="Pilih Mode Workspace"
						/>
					</div>
					<div>
						<label for="user-phone" class="mb-1 block text-xs font-semibold text-[#0F172A]">
							No. WhatsApp
						</label>
						<input
							id="user-phone"
							type="text"
							bind:value={form.phone}
							placeholder="0812xxxxxxxx"
							class="w-full rounded-xl border border-[#E2E8F0] bg-white px-3 py-2 text-sm text-[#0F172A] focus:border-[#EA580C] focus:outline-none"
						/>
					</div>
				</div>

				<div>
					<label for="user-pass" class="mb-1 block text-xs font-semibold text-[#0F172A]">
						{editing ? 'Password Baru (kosongkan jika tidak diubah)' : 'Password (min. 6 karakter)'}
						{#if !editing}<span class="text-red-500">*</span>{/if}
					</label>
					<input
						id="user-pass"
						type="password"
						required={!editing}
						minlength={editing ? undefined : 6}
						bind:value={form.password}
						placeholder="••••••••"
						class="w-full rounded-xl border border-[#E2E8F0] bg-white px-3 py-2 text-sm text-[#0F172A] focus:border-[#EA580C] focus:outline-none"
					/>
				</div>

				<div>
					<label for="user-addr" class="mb-1 block text-xs font-semibold text-[#0F172A]">
						Alamat Bisnis
					</label>
					<input
						id="user-addr"
						type="text"
						bind:value={form.address}
						placeholder="Alamat Bisnis (Opsional)"
						class="w-full rounded-xl border border-[#E2E8F0] bg-white px-3 py-2 text-sm text-[#0F172A] focus:border-[#EA580C] focus:outline-none"
					/>
				</div>

				<div class="flex justify-end gap-2 pt-2">
					<button
						type="button"
						onclick={() => (modalOpen = false)}
						class="rounded-xl border border-[#E2E8F0] px-4 py-2 text-sm font-semibold text-[#64748B] hover:bg-slate-50"
					>
						Batal
					</button>
					<button
						type="submit"
						disabled={saving}
						class="rounded-xl bg-[#EA580C] px-4 py-2 text-sm font-semibold text-white hover:bg-[#C2410C] disabled:opacity-60"
					>
						{saving ? 'Menyimpan...' : 'Simpan Pengguna'}
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}

<!-- Delete Confirmation Modal -->
{#if deleteTarget}
	<div class="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/50 p-4 backdrop-blur-xs">
		<div class="w-full max-w-sm rounded-2xl border border-[#E2E8F0] bg-white p-6 shadow-xl">
			<div class="flex items-start gap-3">
				<div class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-red-50 text-red-600">
					<AlertTriangle class="h-5 w-5" />
				</div>
				<div>
					<h3 class="text-base font-bold text-[#0F172A]">Hapus Pengguna?</h3>
					<p class="mt-1 text-xs text-[#64748B]">
						Akun <strong class="text-[#0F172A]">{deleteTarget.name}</strong> ({deleteTarget.email}) dan seluruh datanya akan dihapus permanen.
					</p>
				</div>
			</div>
			<div class="mt-5 flex justify-end gap-2">
				<button
					type="button"
					onclick={() => (deleteTarget = null)}
					class="rounded-xl border border-[#E2E8F0] px-4 py-2 text-sm font-semibold text-[#64748B] hover:bg-slate-50"
				>
					Batal
				</button>
				<button
					type="button"
					disabled={deleting}
					onclick={confirmDelete}
					class="rounded-xl bg-red-600 px-4 py-2 text-sm font-semibold text-white hover:bg-red-700 disabled:opacity-60"
				>
					{deleting ? 'Menghapus...' : 'Hapus'}
				</button>
			</div>
		</div>
	</div>
{/if}
