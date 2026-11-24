<script>
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { authStore } from '$lib/stores/authStore';
	import { getWorkspaceLabels } from '$lib/utils/workspaceLabels';
	import { startProductTour } from '$lib/tour.js';
	import {
		LayoutDashboard,
		Kanban,
		Users,
		Building2,
		Package,
		CalendarCheck2,
		Settings,
		UserCog,
		LogOut,
		Store,
		Sparkles,
		HelpCircle
	} from 'lucide-svelte';

	const user = $derived($authStore.user);
	const isAdmin = $derived(user?.role === 'ADMIN');
	const labels = $derived(getWorkspaceLabels(user?.business_mode));
	const currentPath = $derived(page.url.pathname);

	const navItems = $derived(
		[
			{ href: '/', label: 'Dashboard', mobileLabel: 'Dashboard', icon: LayoutDashboard },
			{ href: '/pipeline', label: 'Pipeline', mobileLabel: 'Pipeline', icon: Kanban },
			{ href: '/kontak', label: 'Kontak', mobileLabel: 'Kontak', icon: Users },
			{
				href: '/katalog',
				label: labels.catalogMenuLabel,
				mobileLabel: labels.mode === 'PROPERTY' ? 'Properti' : 'Produk',
				icon: labels.mode === 'PROPERTY' ? Building2 : Package
			},
			{ href: '/aktivitas', label: 'Aktivitas & Layanan', mobileLabel: 'Aktivitas', icon: CalendarCheck2 },
			{ href: '/pengguna', label: 'Pengguna', mobileLabel: 'Pengguna', icon: UserCog, adminOnly: true },
			{ href: '/pengaturan', label: 'Pengaturan', mobileLabel: 'Pengaturan', icon: Settings }
		].filter((item) => !item.adminOnly || isAdmin)
	);

	const mobileNavItems = $derived(navItems.slice(0, 5));

	function handleLogout() {
		authStore.logout();
		goto('/login');
	}

	function isActive(href) {
		if (href === '/') return currentPath === '/';
		return currentPath.startsWith(href);
	}
</script>

<!-- Desktop Sticky Sidebar (Orange Gradient with White Active Pill) -->
<aside
	class="hidden w-60 shrink-0 flex-col bg-gradient-to-b from-[#EA580C] via-[#C2410C] to-[#D97706] text-white shadow-sm lg:sticky lg:top-0 lg:flex lg:h-screen"
>
	<!-- Branding Header -->
	<div class="flex items-center justify-between border-b border-white/15 px-4 py-4">
		<a href="/" class="flex items-center gap-2.5">
			<div class="flex h-9 w-9 shrink-0 items-center justify-center rounded-xl bg-white/20 text-white backdrop-blur-sm">
				<Sparkles class="h-5 w-5" />
			</div>
			<div class="min-w-0">
				<span class="block text-base font-extrabold tracking-tight text-white">Cermat</span>
				<span class="block truncate text-[10px] font-medium text-orange-100">Cellix Relationship Management<br>Penjualan</span>
			</div>
		</a>
	</div>

	<!-- Active Tenant Business Card -->
	<div class="px-3 pt-3">
		<div class="rounded-xl border border-white/15 bg-white/10 p-3 backdrop-blur-sm">
			<div class="flex items-center gap-2">
				<div class="flex h-7 w-7 shrink-0 items-center justify-center rounded-lg bg-white/20 text-white">
					<Store class="h-3.5 w-3.5" />
				</div>
				<div class="min-w-0">
					<p class="truncate text-xs font-bold text-white">
						{user?.business_name || 'Cermat Workspace'}
					</p>
					<p class="truncate text-[11px] text-orange-100">
						{user?.name || 'Pengguna'}
					</p>
				</div>
			</div>
		</div>
	</div>

	<!-- Main Navigation Items -->
	<nav data-tour="sidebar-nav" class="flex-1 space-y-1 overflow-y-auto px-3 py-3">
		{#each navItems as item}
			{@const active = isActive(item.href)}
			{@const Icon = item.icon}
			<a
				href={item.href}
				class="flex items-center gap-2.5 rounded-xl px-3 py-2 text-sm font-semibold transition {active
					? 'bg-white text-[#EA580C] shadow-sm'
					: 'text-orange-100 hover:bg-white/10 hover:text-white'}"
			>
				<Icon class="h-4 w-4 shrink-0" />
				<span class="truncate">{item.label}</span>
			</a>
		{/each}
	</nav>

	<!-- Panduan & Logout Footer -->
	<div class="space-y-1 border-t border-white/15 p-3">
		<button
			type="button"
			onclick={() => startProductTour(currentPath)}
			class="flex w-full items-center gap-2.5 rounded-xl px-3 py-2 text-sm font-semibold text-orange-100 transition hover:bg-white/10 hover:text-white"
		>
			<HelpCircle class="h-4 w-4 shrink-0" />
			<span>Panduan Halaman</span>
		</button>
		<button
			type="button"
			onclick={handleLogout}
			class="flex w-full items-center gap-2.5 rounded-xl px-3 py-2 text-sm font-semibold text-orange-100 transition hover:bg-red-500/20 hover:text-white"
		>
			<LogOut class="h-4 w-4 shrink-0" />
			<span>Keluar</span>
		</button>
	</div>
</aside>

<!-- Mobile Bottom Navbar -->
<nav
	data-tour="mobile-nav"
	class="fixed inset-x-0 bottom-0 z-40 grid grid-cols-5 gap-1.5 rounded-t-2xl bg-[#EA580C] px-2.5 py-2.5 shadow-[0_-4px_20px_rgba(234,88,12,0.25)] lg:hidden"
>
	{#each mobileNavItems as item}
		{@const active = isActive(item.href)}
		{@const Icon = item.icon}
		<a
			href={item.href}
			class="flex flex-col items-center justify-center gap-1 rounded-xl py-2 text-[11px] font-semibold transition {active
				? 'bg-white text-[#EA580C] shadow-sm'
				: 'text-orange-100 hover:bg-white/10 hover:text-white'}"
		>
			<Icon class="h-5 w-5 shrink-0" />
			<span class="max-w-full truncate px-1">{item.mobileLabel || item.label}</span>
		</a>
	{/each}
</nav>
