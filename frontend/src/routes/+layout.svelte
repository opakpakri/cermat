<script>
	import '../app.css';
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { browser } from '$app/env';
	import { onMount } from 'svelte';
	import { authStore } from '$lib/stores/authStore';
	import '$lib/stores/themeStore.js';
	import { api } from '$lib/api';
	import Sidebar from '$lib/components/Sidebar.svelte';
	import { startProductTour } from '$lib/tour.js';
	import { Sparkles, Settings, LogOut, HelpCircle } from 'lucide-svelte';

	const publicRoutes = ['/login', '/register'];

	const isPublic = $derived(publicRoutes.includes(page.url.pathname));
	const isSettingsActive = $derived(page.url.pathname.startsWith('/pengaturan'));

	$effect(() => {
		if (!browser) return;
		if (!$authStore.token && !isPublic) {
			goto('/login');
		} else if ($authStore.token && isPublic) {
			goto('/');
		}
	});

	onMount(async () => {
		if ($authStore.token) {
			try {
				const profile = await api.get('/users/profile');
				authStore.setUser(profile);
			} catch {
				// handled by api.js on 401
			}
		}
	});

	function handleLogout() {
		authStore.logout();
		goto('/login');
	}
</script>

{#if isPublic}
	<div class="min-h-screen bg-[#FAFAF9]">
		<slot />
	</div>
{:else if $authStore.token}
	<div class="flex min-h-screen bg-[#FAFAF9]">
		<Sidebar />
		<div class="flex min-w-0 flex-1 flex-col">
			<!-- Mobile Header -->
			<header class="sticky top-0 z-30 flex items-center justify-between border-b border-[#E2E8F0] bg-white px-4 py-3 lg:hidden">
				<a href="/" class="flex items-center gap-2.5">
					<div class="flex h-8 w-8 items-center justify-center rounded-lg bg-[#EA580C] text-white">
						<Sparkles class="h-4 w-4" />
					</div>
					<div>
						<span class="block text-sm font-extrabold leading-tight text-[#0F172A]">
							{$authStore.user?.business_name || 'Cermat'}
						</span>
						<span class="block text-[11px] text-[#64748B]">Halo, {$authStore.user?.name || 'Pengguna'}</span>
					</div>
				</a>

				<div class="flex items-center gap-2">
					<button
						type="button"
						onclick={() => startProductTour(page.url.pathname)}
						title="Panduan Halaman"
						class="rounded-xl border border-[#E2E8F0] p-2 text-[#EA580C] transition hover:bg-[#FFF7ED]"
					>
						<HelpCircle class="h-4 w-4" />
					</button>
					<a
						href="/pengaturan"
						title="Pengaturan"
						class="rounded-xl border border-[#E2E8F0] p-2 transition {isSettingsActive
							? 'bg-[#FFF7ED] text-[#C2410C]'
							: 'text-[#64748B] hover:bg-slate-50 hover:text-[#0F172A]'}"
					>
						<Settings class="h-4 w-4" />
					</a>
					<button
						type="button"
						onclick={handleLogout}
						title="Keluar"
						class="rounded-xl border border-[#E2E8F0] p-2 text-[#64748B] transition hover:bg-[#FEE2E2] hover:text-[#B91C1C]"
					>
						<LogOut class="h-4 w-4" />
					</button>
				</div>
			</header>

			<main class="flex-1 p-4 pb-24 sm:p-6 sm:pb-24 lg:p-8 lg:pb-8">
				<div class="mx-auto max-w-7xl">
					<slot />
				</div>
			</main>
		</div>
	</div>
{/if}
