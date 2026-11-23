<script>
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import { authStore } from '$lib/stores/authStore';
	import { getWorkspaceLabels, formatRupiah, formatDateTime } from '$lib/utils/workspaceLabels';
	import { maybeStartFirstVisitTour } from '$lib/tour.js';
	import Badge from '$lib/components/Badge.svelte';
	import {
		Users,
		TrendingUp,
		CheckCircle2,
		CalendarClock,
		PlusCircle,
		ArrowRight,
		Clock,
		Check
	} from 'lucide-svelte';

	let summary = $state(null);
	let loading = $state(true);

	const user = $derived($authStore.user);
	const labels = $derived(getWorkspaceLabels(user?.business_mode));

	async function loadSummary() {
		loading = true;
		try {
			summary = await api.get('/dashboard/summary');
		} catch (err) {
			console.error(err);
		} finally {
			loading = false;
		}
	}

	async function toggleActivity(id) {
		try {
			await api.patch(`/activities/${id}/toggle`, {});
			await loadSummary();
		} catch (err) {
			alert(err.message);
		}
	}

	function stageLabel(stageKey) {
		return labels.pipelineStages.find((s) => s.key === stageKey)?.label || stageKey;
	}

	onMount(async () => {
		await loadSummary();
		maybeStartFirstVisitTour('/');
	});
</script>

<svelte:head>
	<title>Dashboard — Cermat</title>
</svelte:head>

<div class="space-y-6">
	<!-- Hero Card Gradien (Oranye-Amber) -->
	<section data-tour="dash-hero" class="rounded-2xl bg-gradient-to-r from-[#EA580C] via-[#F97316] to-[#D97706] p-6 text-white shadow-sm sm:p-8">
		<div class="flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
			<div class="min-w-0">
				<p class="text-xs font-semibold uppercase tracking-wider text-orange-100">
					{user?.business_name || 'Cermat'} ({labels.modeLabel})
				</p>
				<h1 class="mt-1 text-2xl font-extrabold sm:text-3xl">
					Halo, {user?.name || 'Tim Sales'}!
				</h1>
				<p class="mt-1 text-sm text-orange-100">
					Pantau ringkasan prospek, nilai pipeline, dan tugas follow-up bulan ini.
				</p>
			</div>

			<a
				href="/pipeline"
				class="inline-flex shrink-0 items-center justify-center gap-2 whitespace-nowrap rounded-xl bg-white px-4 py-2.5 text-sm font-bold text-[#EA580C] shadow-sm transition hover:bg-orange-50"
			>
				<PlusCircle class="h-4 w-4 shrink-0" />
				<span>Tambah Prospek / Deal</span>
			</a>
		</div>

		<!-- 4 Ringkasan Metrik Bulan Ini -->
		<div class="mt-6 grid grid-cols-2 gap-3 sm:grid-cols-4 sm:gap-4">
			<div class="rounded-xl bg-white/15 p-4 backdrop-blur-sm">
				<div class="flex items-center gap-2 text-xs font-medium text-orange-100">
					<Users class="h-4 w-4 shrink-0" />
					<span class="truncate">Prospek Baru Bulan Ini</span>
				</div>
				<p class="mt-2 text-2xl font-extrabold">{loading ? '...' : summary?.new_leads_this_month ?? 0}</p>
			</div>

			<div class="rounded-xl bg-white/15 p-4 backdrop-blur-sm">
				<div class="flex items-center gap-2 text-xs font-medium text-orange-100">
					<TrendingUp class="h-4 w-4 shrink-0" />
					<span class="truncate">Estimasi Nilai Pipeline</span>
				</div>
				<p class="mt-2 text-lg font-extrabold sm:text-xl">
					{loading ? '...' : formatRupiah(summary?.pipeline_value ?? 0)}
				</p>
			</div>

			<div class="rounded-xl bg-white/15 p-4 backdrop-blur-sm">
				<div class="flex items-center gap-2 text-xs font-medium text-orange-100">
					<CheckCircle2 class="h-4 w-4 shrink-0" />
					<span class="truncate">Total Closing / Deal</span>
				</div>
				<p class="mt-2 text-2xl font-extrabold">
					{loading ? '...' : `${summary?.won_deals_count ?? 0} Deal`}
				</p>
				<p class="text-xs text-orange-100">
					{loading ? '' : formatRupiah(summary?.won_deals_value ?? 0)}
				</p>
			</div>

			<div class="rounded-xl bg-white/15 p-4 backdrop-blur-sm">
				<div class="flex items-center gap-2 text-xs font-medium text-orange-100">
					<CalendarClock class="h-4 w-4 shrink-0" />
					<span class="truncate">Tugas Follow-up Hari Ini</span>
				</div>
				<p class="mt-2 text-2xl font-extrabold">{loading ? '...' : summary?.followups_today ?? 0}</p>
			</div>
		</div>
	</section>

	<!-- Aksi Cepat (Quick Actions) -->
	<section data-tour="dash-quick" class="grid gap-4 sm:grid-cols-2">
		<a
			href="/pipeline"
			class="group flex min-w-0 items-center justify-between gap-3 overflow-hidden rounded-2xl border border-[#E2E8F0] bg-white p-4 shadow-sm transition hover:border-[#EA580C] sm:p-5"
		>
			<div class="flex min-w-0 flex-1 items-center gap-3.5 sm:gap-4">
				<div class="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl bg-[#FFF7ED] text-[#EA580C] sm:h-12 sm:w-12">
					<PlusCircle class="h-5 w-5 sm:h-6 sm:w-6" />
				</div>
				<div class="min-w-0 flex-1">
					<h2 class="truncate text-sm font-bold text-[#0F172A] sm:text-base">Tambah Prospek / Deal Baru</h2>
					<p class="truncate text-xs text-[#64748B] sm:text-sm">Catat calon pembeli & estimasi nilai transaksi</p>
				</div>
			</div>
			<ArrowRight class="h-5 w-5 shrink-0 text-[#64748B] transition group-hover:translate-x-1 group-hover:text-[#EA580C]" />
		</a>

		<a
			href="/aktivitas"
			class="group flex min-w-0 items-center justify-between gap-3 overflow-hidden rounded-2xl border border-[#E2E8F0] bg-white p-4 shadow-sm transition hover:border-[#EA580C] sm:p-5"
		>
			<div class="flex min-w-0 flex-1 items-center gap-3.5 sm:gap-4">
				<div class="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl bg-[#E0F2FE] text-[#0369A1] sm:h-12 sm:w-12">
					<CalendarClock class="h-5 w-5 sm:h-6 sm:w-6" />
				</div>
				<div class="min-w-0 flex-1">
					<h2 class="truncate text-sm font-bold text-[#0F172A] sm:text-base">Jadwal Follow-up Hari Ini</h2>
					<p class="truncate text-xs text-[#64748B] sm:text-sm">Pantau tugas telepon, survei & tiket layanan</p>
				</div>
			</div>
			<ArrowRight class="h-5 w-5 shrink-0 text-[#64748B] transition group-hover:translate-x-1 group-hover:text-[#EA580C]" />
		</a>
	</section>

	<!-- Aktivitas & Pengingat Terbaru + Progres Transaksi Terkini -->
	<div data-tour="dash-activity" class="grid gap-4 lg:grid-cols-2">
		<section class="rounded-2xl border border-[#E2E8F0] bg-white shadow-sm">
			<div class="flex items-center justify-between border-b border-[#E2E8F0] px-6 py-4">
				<div class="flex items-center gap-2">
					<Clock class="h-5 w-5 text-[#EA580C]" />
					<h2 class="text-base font-bold text-[#0F172A]">Jadwal Follow-up & Pengingat</h2>
				</div>
				<a href="/aktivitas" class="text-sm font-semibold text-[#EA580C] hover:underline">Lihat semua</a>
			</div>

			{#if loading}
				<div class="p-8 text-center text-sm text-[#64748B]">Memuat aktivitas...</div>
			{:else if !summary?.recent_activities?.length}
				<div class="p-8 text-center text-sm text-[#64748B]">Belum ada jadwal follow-up.</div>
			{:else}
				<div class="divide-y divide-[#E2E8F0]">
					{#each summary.recent_activities as act}
						<div class="flex items-center justify-between gap-3 px-6 py-4">
							<div class="flex items-start gap-3 min-w-0">
								<button
									type="button"
									disabled={act.is_completed}
									onclick={() => !act.is_completed && toggleActivity(act.id)}
									class="mt-0.5 flex h-5 w-5 shrink-0 items-center justify-center rounded border transition-colors {act.is_completed
										? 'cursor-not-allowed border-[#15803D] bg-[#15803D] text-white opacity-90'
										: 'border-slate-300 bg-white hover:border-[#EA580C]'}"
									title={act.is_completed ? 'Sudah selesai' : 'Tandai selesai'}
								>
									{#if act.is_completed}
										<Check class="h-3.5 w-3.5" />
									{/if}
								</button>
								<div class="min-w-0">
									<p class="truncate text-sm font-bold text-[#0F172A] {act.is_completed ? 'line-through text-[#64748B]' : ''}">
										{act.title}
									</p>
									<p class="mt-0.5 truncate text-xs text-[#64748B]">
										{act.contact_name} • {act.activity_type}
										{#if act.scheduled_at} • {formatDateTime(act.scheduled_at)}{/if}
									</p>
								</div>
							</div>
							<Badge variant={act.is_completed ? 'won' : 'warning'}>
								{act.is_completed ? 'Selesai' : 'Tertunda'}
							</Badge>
						</div>
					{/each}
				</div>
			{/if}
		</section>

		<section class="rounded-2xl border border-[#E2E8F0] bg-white shadow-sm">
			<div class="flex items-center justify-between border-b border-[#E2E8F0] px-6 py-4">
				<div class="flex items-center gap-2">
					<TrendingUp class="h-5 w-5 text-[#EA580C]" />
					<h2 class="text-base font-bold text-[#0F172A]">Progres Transaksi Terkini</h2>
				</div>
				<a href="/pipeline" class="text-sm font-semibold text-[#EA580C] hover:underline">Buka Pipeline</a>
			</div>

			{#if loading}
				<div class="p-8 text-center text-sm text-[#64748B]">Memuat transaksi...</div>
			{:else if !summary?.recent_deals?.length}
				<div class="p-8 text-center text-sm text-[#64748B]">Belum ada transaksi di pipeline.</div>
			{:else}
				<div class="divide-y divide-[#E2E8F0]">
					{#each summary.recent_deals as deal}
						<div class="flex items-center justify-between gap-3 px-6 py-4">
							<div class="min-w-0">
								<p class="truncate text-sm font-bold text-[#0F172A]">{deal.title}</p>
								<p class="mt-0.5 truncate text-xs text-[#64748B]">
									{deal.contact_name} • <span class="font-semibold text-[#0F172A]">{formatRupiah(deal.deal_value)}</span>
								</p>
							</div>
							<Badge variant={deal.stage === 'WON' ? 'won' : deal.stage === 'LOST' ? 'lost' : 'info'}>
								{stageLabel(deal.stage)}
							</Badge>
						</div>
					{/each}
				</div>
			{/if}
		</section>
	</div>
</div>
