<script>
	import { goto } from '$app/navigation';
	import { api } from '$lib/api';
	import { authStore } from '$lib/stores/authStore';
	import {
		Sparkles,
		Mail,
		Lock,
		Eye,
		EyeOff,
		ArrowRight,
		TrendingUp,
		CalendarCheck2,
		ShieldCheck
	} from 'lucide-svelte';

	let email = $state('');
	let password = $state('');
	let showPassword = $state(false);
	let loading = $state(false);
	let errorMsg = $state('');

	async function handleSubmit(e) {
		e.preventDefault();
		errorMsg = '';
		loading = true;
		try {
			const res = await api.post('/auth/login', { email, password });
			authStore.login(res.token, res.user);
			goto('/');
		} catch (err) {
			errorMsg = err.message || 'Gagal masuk';
		} finally {
			loading = false;
		}
	}
</script>

<svelte:head>
	<title>Masuk — Cermat</title>
</svelte:head>

<div class="min-h-screen grid lg:grid-cols-2">
	<!-- Left Branding Panel -->
	<div class="hidden lg:flex lg:flex-col lg:justify-center bg-gradient-to-br from-[#EA580C] via-[#DC4E0A] to-[#C2410C] text-white px-12 xl:px-20">
		<div class="max-w-md">
			<div class="flex items-center gap-3 mb-8">
				<div class="flex h-10 w-10 items-center justify-center rounded-xl bg-white/15 backdrop-blur-xs">
					<Sparkles class="h-5 w-5 text-white" />
				</div>
				<span class="text-base font-bold tracking-tight">CERMAT - Cellix Relationship Management Terpadu</span>
			</div>

			<h1 class="text-4xl font-extrabold leading-tight tracking-tight">
				Relasi bisnis Anda, rapi di mana saja.
			</h1>

			<ul class="mt-8 space-y-4 text-sm font-medium text-orange-50">
				<li class="flex items-center gap-3">
					<TrendingUp class="h-4 w-4 shrink-0 text-orange-200" />
					<span>Pantau pipeline penjualan secara langsung</span>
				</li>
				<li class="flex items-center gap-3">
					<CalendarCheck2 class="h-4 w-4 shrink-0 text-orange-200" />
					<span>Kelola kontak, follow-up, dan tiket layanan terpadu</span>
				</li>
				<li class="flex items-center gap-3">
					<ShieldCheck class="h-4 w-4 shrink-0 text-orange-200" />
					<span>Data setiap akun terpisah dan aman</span>
				</li>
			</ul>
		</div>
	</div>

	<!-- Right Login Form Panel -->
	<div class="flex flex-col items-center justify-center bg-[#FAFAF9] p-6 sm:p-10">
		<div class="w-full max-w-md rounded-2xl border border-[#E2E8F0] bg-white p-7 shadow-sm sm:p-8">
			<div class="mb-6 flex items-center gap-3">
				<div class="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl bg-[#EA580C] text-white shadow-sm">
					<Sparkles class="h-5 w-5" />
				</div>
				<div>
					<div class="text-xl font-extrabold leading-tight text-[#EA580C]">Cermat</div>
					<p class="text-xs text-[#64748B]">Cellix Relationship Management</p>
				</div>
			</div>

			<h2 class="text-2xl font-extrabold text-[#0F172A]">Masuk ke Workspace</h2>
			<p class="mt-1 text-sm text-[#64748B]">Kelola prospek, pipeline, dan layanan pelanggan Anda.</p>

			{#if errorMsg}
				<div class="mt-4 rounded-xl bg-[#FEE2E2] px-3.5 py-2.5 text-xs font-medium text-[#B91C1C]">
					{errorMsg}
				</div>
			{/if}

			<form onsubmit={handleSubmit} class="mt-6 space-y-4">
				<div>
					<label for="email" class="mb-1.5 block text-sm font-semibold text-[#0F172A]">Email</label>
					<div class="relative">
						<Mail class="pointer-events-none absolute left-3.5 top-1/2 h-4 w-4 -translate-y-1/2 text-[#94A3B8]" />
						<input
							id="email"
							type="email"
							required
							bind:value={email}
							placeholder="nama@perusahaan.com"
							class="w-full rounded-xl border border-[#E2E8F0] bg-white py-2.5 pl-10 pr-3.5 text-sm text-[#0F172A] placeholder:text-[#94A3B8] focus:border-[#EA580C] focus:outline-none focus:ring-1 focus:ring-[#EA580C]"
						/>
					</div>
				</div>

				<div>
					<label for="password" class="mb-1.5 block text-sm font-semibold text-[#0F172A]">Password</label>
					<div class="relative">
						<Lock class="pointer-events-none absolute left-3.5 top-1/2 h-4 w-4 -translate-y-1/2 text-[#94A3B8]" />
						<input
							id="password"
							type={showPassword ? 'text' : 'password'}
							required
							bind:value={password}
							placeholder="••••••••"
							class="w-full rounded-xl border border-[#E2E8F0] bg-white py-2.5 pl-10 pr-10 text-sm text-[#0F172A] placeholder:text-[#94A3B8] focus:border-[#EA580C] focus:outline-none focus:ring-1 focus:ring-[#EA580C]"
						/>
						<button
							type="button"
							onclick={() => (showPassword = !showPassword)}
							aria-label={showPassword ? 'Sembunyikan password' : 'Lihat password'}
							class="absolute right-3 top-1/2 -translate-y-1/2 text-[#64748B] hover:text-[#0F172A]"
						>
							{#if showPassword}
								<Eye class="h-4 w-4" />
							{:else}
								<EyeOff class="h-4 w-4" />
							{/if}
						</button>
					</div>
				</div>

				<button
					type="submit"
					disabled={loading}
					class="flex w-full items-center justify-center gap-2 rounded-xl bg-[#EA580C] py-2.5 px-4 text-sm font-semibold text-white shadow-sm transition-colors hover:bg-[#C2410C] disabled:opacity-60"
				>
					<span>{loading ? 'Memproses...' : 'Masuk'}</span>
					<ArrowRight class="h-4 w-4" />
				</button>
			</form>

			<p class="mt-6 text-center text-xs text-[#64748B]">
				Pendaftaran akun baru hanya dapat dibuat oleh Admin.
			</p>
		</div>
	</div>
</div>
