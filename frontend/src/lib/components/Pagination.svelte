<script>
	import { ChevronLeft, ChevronRight } from 'lucide-svelte';

	let { currentPage = $bindable(1), totalItems = 0, pageSize = 10 } = $props();

	const totalPages = $derived(Math.max(1, Math.ceil(totalItems / pageSize)));
	const startIdx = $derived(totalItems === 0 ? 0 : (currentPage - 1) * pageSize + 1);
	const endIdx = $derived(Math.min(currentPage * pageSize, totalItems));
</script>

<div class="flex flex-col items-center justify-between gap-3 border-t border-[#E2E8F0] px-5 py-3 text-xs text-[#64748B] sm:flex-row">
	<span>Menampilkan {startIdx}–{endIdx} dari {totalItems} data</span>
	<div class="flex items-center gap-1.5">
		<button
			type="button"
			disabled={currentPage <= 1}
			onclick={() => (currentPage = Math.max(1, currentPage - 1))}
			class="inline-flex items-center gap-1 rounded-lg border border-[#E2E8F0] bg-white px-2.5 py-1.5 text-xs font-medium text-[#64748B] hover:bg-slate-50 disabled:opacity-40"
		>
			<ChevronLeft class="h-3.5 w-3.5" />
			<span>Sebelumnya</span>
		</button>
		<span class="flex h-7 min-w-[28px] items-center justify-center rounded-lg border border-[#EA580C] bg-[#FFF7ED] px-2 text-xs font-bold text-[#C2410C]">
			{currentPage}
		</span>
		<button
			type="button"
			disabled={currentPage >= totalPages}
			onclick={() => (currentPage = Math.min(totalPages, currentPage + 1))}
			class="inline-flex items-center gap-1 rounded-lg border border-[#E2E8F0] bg-white px-2.5 py-1.5 text-xs font-medium text-[#64748B] hover:bg-slate-50 disabled:opacity-40"
		>
			<span>Berikutnya</span>
			<ChevronRight class="h-3.5 w-3.5" />
		</button>
	</div>
</div>
