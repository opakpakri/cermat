<script>
	let { variant = 'default', children } = $props();

	const variantMap = {
		won: 'bg-[#DCFCE7] text-[#15803D]',
		available: 'bg-[#DCFCE7] text-[#15803D]',
		resolved: 'bg-[#DCFCE7] text-[#15803D]',
		customer: 'bg-[#DCFCE7] text-[#15803D]',
		info: 'bg-[#E0F2FE] text-[#0369A1]',
		booked: 'bg-[#E0F2FE] text-[#0369A1]',
		prospect: 'bg-[#E0F2FE] text-[#0369A1]',
		in_progress: 'bg-[#E0F2FE] text-[#0369A1]',
		warning: 'bg-[#FFEDD5] text-[#C2410C]',
		lead: 'bg-[#FFEDD5] text-[#C2410C]',
		open: 'bg-[#FFEDD5] text-[#C2410C]',
		danger: 'bg-[#FEE2E2] text-[#B91C1C]',
		lost: 'bg-[#FEE2E2] text-[#B91C1C]',
		sold: 'bg-[#FEE2E2] text-[#B91C1C]',
		high: 'bg-[#FEE2E2] text-[#B91C1C]',
		orange: 'bg-[#FFEDD5] text-[#C2410C]',
		purple: 'bg-[#F3E8FF] text-[#7E22CE]',
		slate: 'bg-[#F1F5F9] text-[#475569]',
		default: 'bg-slate-100 text-slate-700'
	};

	const classes = $derived(variantMap[variant?.toLowerCase()] || variantMap.default);
</script>

<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium {classes}">
	{@render children?.()}
</span>
