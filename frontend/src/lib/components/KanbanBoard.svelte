<script>
	import { formatRupiah, formatDateTime } from '$lib/utils/workspaceLabels';
	import { Phone, Package, ArrowRight, Clock } from 'lucide-svelte';

	let { deals = [], stages = [], onStageChange, onSelectDeal } = $props();

	function dealsForStage(stageKey) {
		return deals.filter((d) => d.stage === stageKey);
	}

	function stageTotal(stageKey) {
		return dealsForStage(stageKey).reduce((sum, d) => sum + Number(d.deal_value || 0), 0);
	}
</script>

<div class="flex gap-4 overflow-x-auto pb-4">
	{#each stages as stage}
		{@const stageDeals = dealsForStage(stage.key)}
		<div class="w-72 shrink-0 bg-slate-50/80 border border-[#E2E8F0] rounded-xl flex flex-col max-h-[72vh]">
			<!-- Column Header -->
			<div class="p-3.5 border-b border-[#E2E8F0] bg-white rounded-t-xl">
				<div class="flex items-center justify-between gap-2">
					<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-semibold {stage.badge}">
						{stage.label}
					</span>
					<span class="text-xs font-semibold text-[#64748B]">{stageDeals.length}</span>
				</div>
				<p class="mt-1.5 text-xs font-medium text-[#64748B]">
					{formatRupiah(stageTotal(stage.key))}
				</p>
			</div>

			<!-- Cards -->
			<div class="p-2.5 space-y-2.5 overflow-y-auto flex-1">
				{#if stageDeals.length === 0}
					<div class="text-center py-8 text-xs text-[#64748B]">Belum ada prospek</div>
				{:else}
					{#each stageDeals as deal}
						<div class="bg-white border border-[#E2E8F0] rounded-lg p-3 shadow-sm hover:border-[#EA580C]/50 transition-colors">
							<div class="flex items-start justify-between gap-2">
								<button
									type="button"
									onclick={() => onSelectDeal?.(deal)}
									class="text-left text-sm font-semibold text-[#0F172A] hover:text-[#EA580C] line-clamp-2"
								>
									{deal.title}
								</button>
							</div>

							<p class="mt-1 text-[11px] font-mono text-[#64748B]">{deal.deal_code}</p>

							<div class="mt-2.5 space-y-1 text-xs text-[#64748B]">
								<div class="flex items-center gap-1.5">
									<Phone class="w-3.5 h-3.5 text-[#EA580C]" />
									<span class="truncate font-medium text-[#0F172A]">{deal.contact_name}</span>
								</div>
								{#if deal.catalog_item_name}
									<div class="flex items-center gap-1.5">
										<Package class="w-3.5 h-3.5 text-slate-400" />
										<span class="truncate">{deal.catalog_item_name}</span>
									</div>
								{/if}
							</div>

							<div class="mt-2.5 space-y-0.5 text-[11px] text-[#64748B]">
								<div class="flex items-center gap-1.5">
									<Clock class="w-3 h-3 shrink-0 text-slate-400" />
									<span class="truncate">Dibuat: {formatDateTime(deal.created_at)}</span>
								</div>
								<div class="flex items-center gap-1.5">
									<Clock class="w-3 h-3 shrink-0 text-slate-400" />
									<span class="truncate">Diubah: {formatDateTime(deal.updated_at)}</span>
								</div>
							</div>

							<div class="mt-3 pt-2.5 border-t border-slate-100 flex items-center justify-between gap-2">
								<span class="text-xs font-bold text-[#0F172A]">
									{formatRupiah(deal.deal_value)}
								</span>

								<select
									aria-label="Ubah Tahapan"
									value={deal.stage}
									onchange={(e) => onStageChange?.(deal.id, e.currentTarget.value)}
									class="text-[11px] border border-[#E2E8F0] rounded-md px-1.5 py-1 bg-[#FAFAF9] text-[#0F172A] focus:outline-none focus:border-[#EA580C]"
								>
									{#each stages as opt}
										<option value={opt.key}>{opt.label}</option>
									{/each}
								</select>
							</div>
						</div>
					{/each}
				{/if}
			</div>
		</div>
	{/each}
</div>
