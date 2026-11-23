<script>
	let { title, value, subtitle = '', icon: Icon = null, accent = 'orange' } = $props();

	const accentClasses = {
		orange: 'bg-[#FFF7ED] text-[#EA580C]',
		green: 'bg-[#DCFCE7] text-[#15803D]',
		sky: 'bg-[#E0F2FE] text-[#0369A1]',
		red: 'bg-[#FEE2E2] text-[#B91C1C]',
		slate: 'bg-slate-100 text-slate-700'
	};
</script>

<div class="flex items-start justify-between gap-2.5 rounded-2xl border border-[#E2E8F0] bg-white p-4 shadow-sm transition hover:border-slate-300 sm:gap-3 sm:p-5">
	<div class="min-w-0">
		<p class="text-[11px] font-medium uppercase leading-snug tracking-wider text-[#64748B] sm:text-xs">{title}</p>
		<p class="mt-1.5 truncate text-xl font-bold text-[#0F172A] sm:mt-2 sm:text-2xl">{value}</p>
		{#if subtitle}
			<p class="mt-1 truncate text-xs text-[#64748B]">{subtitle}</p>
		{/if}
	</div>
	{#if Icon}
		<div class="flex h-9 w-9 shrink-0 items-center justify-center rounded-xl sm:h-10 sm:w-10 {accentClasses[accent] || accentClasses.orange}">
			<Icon class="h-4 w-4 sm:h-5 sm:w-5" strokeWidth={1.75} />
		</div>
	{/if}
</div>
