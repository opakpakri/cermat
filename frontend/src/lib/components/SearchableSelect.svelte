<script>
	import { Search, Plus, ChevronDown, ChevronUp } from 'lucide-svelte';

	let {
		options = [],
		value = $bindable(''),
		placeholder = 'Pilih opsi',
		searchPlaceholder = 'Cari...',
		searchable = true,
		addHref = '',
		onAdd = null,
		onChange = null
	} = $props();

	let open = $state(false);
	let query = $state('');
	let containerEl = $state(null);

	const selectedOption = $derived(options.find((opt) => String(opt.value) === String(value)));

	const filteredOptions = $derived(
		query.trim() === ''
			? options
			: options.filter((opt) =>
					String(opt.label || '')
						.toLowerCase()
						.includes(query.trim().toLowerCase())
				)
	);

	function selectOption(optValue) {
		value = String(optValue);
		open = false;
		query = '';
		onChange?.(value);
	}

	function handleWindowClick(e) {
		if (open && containerEl && !containerEl.contains(e.target)) {
			open = false;
		}
	}
</script>

<svelte:window onclick={handleWindowClick} />

<div bind:this={containerEl} class="relative w-full">
	<button
		type="button"
		onclick={() => (open = !open)}
		class="flex w-full items-center justify-between rounded-xl border border-[#E2E8F0] bg-white px-3.5 py-2 text-left text-sm transition focus:border-[#EA580C] focus:outline-none focus:ring-2 focus:ring-[#FFF7ED]"
	>
		<span class="truncate {selectedOption ? 'text-[#0F172A]' : 'text-[#64748B]'}">
			{selectedOption ? selectedOption.label : placeholder}
		</span>
		{#if open}
			<ChevronUp class="h-4 w-4 shrink-0 text-[#64748B]" />
		{:else}
			<ChevronDown class="h-4 w-4 shrink-0 text-[#64748B]" />
		{/if}
	</button>

	{#if open}
		<div class="absolute left-0 right-0 z-50 mt-1.5 rounded-2xl border border-[#E2E8F0] bg-white p-2.5 shadow-xl">
			{#if searchable}
				<div class="mb-2 flex items-center gap-2">
					<div class="relative flex-1">
						<Search class="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-[#64748B]" />
						<input
							type="text"
							bind:value={query}
							placeholder={searchPlaceholder}
							class="w-full rounded-xl border border-[#EA580C] bg-white py-2 pl-9 pr-3 text-sm text-[#0F172A] outline-none ring-2 ring-[#FFF7ED]"
						/>
					</div>
					{#if addHref}
						<a
							href={addHref}
							title="Tambah baru"
							class="flex h-9 w-9 shrink-0 items-center justify-center rounded-xl border border-[#E2E8F0] bg-white text-[#0F172A] transition hover:border-[#EA580C] hover:bg-[#FFF7ED] hover:text-[#EA580C]"
						>
							<Plus class="h-4 w-4" />
						</a>
					{:else if onAdd}
						<button
							type="button"
							onclick={() => {
								open = false;
								onAdd();
							}}
							title="Tambah baru"
							class="flex h-9 w-9 shrink-0 items-center justify-center rounded-xl border border-[#E2E8F0] bg-white text-[#0F172A] transition hover:border-[#EA580C] hover:bg-[#FFF7ED] hover:text-[#EA580C]"
						>
							<Plus class="h-4 w-4" />
						</button>
					{/if}
				</div>
			{/if}

			<div class="max-h-48 overflow-y-auto">
				<p class="px-2.5 py-1.5 text-xs font-medium text-[#64748B]">{placeholder}</p>
				{#if filteredOptions.length === 0}
					<p class="px-2.5 py-3 text-center text-xs text-[#64748B]">Data tidak ditemukan.</p>
				{:else}
					{#each filteredOptions as opt}
						<button
							type="button"
							onclick={() => selectOption(opt.value)}
							class="w-full rounded-lg px-2.5 py-2 text-left text-sm transition {String(opt.value) === String(value)
								? 'bg-[#FFF7ED] font-semibold text-[#C2410C]'
								: 'text-[#0F172A] hover:bg-slate-50'}"
						>
							{opt.label}
						</button>
					{/each}
				{/if}
			</div>
		</div>
	{/if}
</div>
