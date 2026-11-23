<script>
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import { authStore } from '$lib/stores/authStore';
	import { themeStore } from '$lib/stores/themeStore.js';
	import { Save, Sun, Moon, Monitor, UserCog, ArrowRight } from 'lucide-svelte';

	let name = $state('');
	let email = $state('');
	let business_name = $state('');
	let business_mode = $state('OFFICE');
	let phone = $state('');
	let address = $state('');
	let new_password = $state('');
	let loading = $state(true);
	let saving = $state(false);
	let message = $state('');
	let errorMsg = $state('');

	onMount(async () => {
		try {
			const profile = await api.get('/users/profile');
			name = profile.name || '';
			email = profile.email || '';
			business_name = profile.business_name || '';
			business_mode = profile.business_mode || 'OFFICE';
			phone = profile.phone || '';
			address = profile.address || '';
			authStore.setUser(profile);
		} catch (err) {
			errorMsg = err.message;
		} finally {
			loading = false;
		}
	});

	async function handleSave(e) {
		e.preventDefault();
		message = '';
		errorMsg = '';
		saving = true;
		try {
			const payload = {
				name,
				business_name,
				business_mode,
				phone: phone || null,
				address: address || null,
				new_password: new_password || null
			};
			const updated = await api.put('/users/profile', payload);
			authStore.setUser(updated);
			new_password = '';
			message = 'Profil akun dan bisnis berhasil diperbarui.';
		} catch (err) {
			errorMsg = err.message;
		} finally {
			saving = false;
		}
	}
</script>

<div class="space-y-6">
	<div>
		<h1 class="text-2xl font-extrabold text-[#0F172A]">Pengaturan Akun & Bisnis</h1>
		<p class="text-sm text-[#64748B]">Kelola informasi pengguna, tampilan aplikasi, dan detail profil bisnis Anda.</p>
	</div>

	<!-- Kartu Pengaturan Tema (Light / Dark Mode) -->
	<div data-tour="settings-theme" class="flex flex-col justify-between gap-4 rounded-2xl border border-[#E2E8F0] bg-white p-6 shadow-sm sm:flex-row sm:items-center">
		<div>
			<h2 class="text-base font-bold text-[#0F172A]">Tampilan Aplikasi</h2>
			<p class="text-sm text-[#64748B]">Pilih tema terang atau gelap sesuai kenyamanan Anda.</p>
		</div>
		<div class="grid grid-cols-3 rounded-xl border border-[#E2E8F0] bg-[#FAFAF9] p-1 sm:w-80">
			<button
				type="button"
				onclick={() => themeStore.setTheme('light')}
				class="flex items-center justify-center gap-1.5 rounded-lg py-2 text-xs font-semibold transition {$themeStore === 'light'
					? 'bg-white text-[#EA580C] shadow-sm'
					: 'text-[#64748B] hover:text-[#0F172A]'}"
			>
				<Sun class="h-4 w-4" />
				<span>Terang</span>
			</button>
			<button
				type="button"
				onclick={() => themeStore.setTheme('dark')}
				class="flex items-center justify-center gap-1.5 rounded-lg py-2 text-xs font-semibold transition {$themeStore === 'dark'
					? 'bg-white text-[#EA580C] shadow-sm'
					: 'text-[#64748B] hover:text-[#0F172A]'}"
			>
				<Moon class="h-4 w-4" />
				<span>Gelap</span>
			</button>
			<button
				type="button"
				onclick={() => themeStore.setTheme('system')}
				class="flex items-center justify-center gap-1.5 rounded-lg py-2 text-xs font-semibold transition {$themeStore === 'system'
					? 'bg-white text-[#EA580C] shadow-sm'
					: 'text-[#64748B] hover:text-[#0F172A]'}"
			>
				<Monitor class="h-4 w-4" />
				<span>Sistem</span>
			</button>
		</div>
	</div>

	{#if $authStore.user?.role === 'ADMIN'}
		<a
			href="/pengguna"
			class="flex items-center justify-between gap-4 rounded-2xl border border-[#E2E8F0] bg-white p-6 shadow-sm transition hover:border-[#EA580C] hover:shadow-md"
		>
			<div class="flex items-center gap-4">
				<div class="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl bg-[#FFF7ED] text-[#EA580C]">
					<UserCog class="h-5 w-5" />
				</div>
				<div>
					<h2 class="text-base font-bold text-[#0F172A]">Manajemen Pengguna</h2>
					<p class="text-sm text-[#64748B]">Buat, edit, dan kelola akun bisnis pengguna (khusus Admin).</p>
				</div>
			</div>
			<ArrowRight class="h-5 w-5 shrink-0 text-[#64748B]" />
		</a>
	{/if}

	{#if loading}
		<div class="rounded-2xl border border-[#E2E8F0] bg-white p-8 text-center text-sm text-[#64748B] shadow-sm">
			Memuat pengaturan...
		</div>
	{:else}
		<div data-tour="settings-form" class="rounded-2xl border border-[#E2E8F0] bg-white p-6 shadow-sm sm:p-8">
			<form onsubmit={handleSave} class="space-y-6">
				<div>
					<h2 class="text-base font-bold text-[#0F172A]">Pengaturan Akun</h2>
					<p class="text-sm text-[#64748B]">Informasi identitas pengguna dan email login akun Anda.</p>
				</div>

				{#if message}
					<div class="rounded-xl bg-[#DCFCE7] px-4 py-3 text-sm font-medium text-[#15803D]">{message}</div>
				{/if}
				{#if errorMsg}
					<div class="rounded-xl bg-[#FEE2E2] px-4 py-3 text-sm font-medium text-[#B91C1C]">{errorMsg}</div>
				{/if}

				<div class="grid gap-4 sm:grid-cols-2">
					<div>
						<label for="set-name" class="block text-sm font-semibold text-[#0F172A]">Nama Pengguna</label>
						<input
							id="set-name"
							type="text"
							value={name}
							readonly
							class="mt-1.5 w-full cursor-not-allowed rounded-xl border border-[#E2E8F0] bg-[#FAFAF9] px-3.5 py-2.5 text-sm text-[#64748B] outline-none"
						/>
					</div>
					<div>
						<label for="set-email" class="block text-sm font-semibold text-[#0F172A]">Email Login</label>
						<input
							id="set-email"
							type="email"
							value={email}
							readonly
							class="mt-1.5 w-full cursor-not-allowed rounded-xl border border-[#E2E8F0] bg-[#FAFAF9] px-3.5 py-2.5 text-sm text-[#64748B] outline-none"
						/>
					</div>
				</div>

				<div class="border-t border-[#E2E8F0] pt-5">
					<h2 class="text-base font-bold text-[#0F172A]">Profil Bisnis / Kantor</h2>
					<p class="text-sm text-[#64748B]">Detail identitas bisnis yang tampil pada workspace Anda.</p>
				</div>

				<div class="grid gap-4 sm:grid-cols-2">
					<div>
						<label for="set-business" class="block text-sm font-semibold text-[#0F172A]">Nama Bisnis / Perusahaan</label>
						<input
							id="set-business"
							type="text"
							required
							bind:value={business_name}
							class="mt-1.5 w-full rounded-xl border border-[#E2E8F0] bg-white px-3.5 py-2.5 text-sm text-[#0F172A] outline-none focus:border-[#EA580C] focus:ring-2 focus:ring-[#FFF7ED]"
						/>
					</div>
					<div>
						<label for="set-phone" class="block text-sm font-semibold text-[#0F172A]">Nomor WhatsApp Bisnis</label>
						<input
							id="set-phone"
							type="text"
							bind:value={phone}
							class="mt-1.5 w-full rounded-xl border border-[#E2E8F0] bg-white px-3.5 py-2.5 text-sm text-[#0F172A] outline-none focus:border-[#EA580C] focus:ring-2 focus:ring-[#FFF7ED]"
						/>
					</div>
				</div>

				<div>
					<label for="set-address" class="block text-sm font-semibold text-[#0F172A]">Alamat Kantor</label>
					<textarea
						id="set-address"
						rows="3"
						bind:value={address}
						class="mt-1.5 w-full rounded-xl border border-[#E2E8F0] bg-white px-3.5 py-2.5 text-sm text-[#0F172A] outline-none focus:border-[#EA580C] focus:ring-2 focus:ring-[#FFF7ED]"
					></textarea>
				</div>

				<div>
					<label for="set-password" class="block text-sm font-semibold text-[#0F172A]">Password Baru (Opsional)</label>
					<input
						id="set-password"
						type="password"
						bind:value={new_password}
						placeholder="Kosongkan jika tidak ingin mengubah password"
						class="mt-1.5 w-full rounded-xl border border-[#E2E8F0] bg-white px-3.5 py-2.5 text-sm text-[#0F172A] outline-none focus:border-[#EA580C] focus:ring-2 focus:ring-[#FFF7ED]"
					/>
				</div>

				<div class="flex justify-end pt-2">
					<button
						type="submit"
						disabled={saving}
						class="inline-flex items-center gap-2 rounded-xl bg-[#EA580C] px-5 py-2.5 text-sm font-semibold text-white shadow-sm transition hover:bg-[#C2410C] disabled:opacity-60"
					>
						<Save class="h-4 w-4" />
						<span>{saving ? 'Menyimpan...' : 'Simpan Pengaturan'}</span>
					</button>
				</div>
			</form>
		</div>
	{/if}
</div>
