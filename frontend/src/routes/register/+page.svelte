<script>
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';

	onMount(() => {
		goto('/login');
	});
</script>

<div class="flex min-h-screen items-center justify-center bg-[#FAFAF9] p-4">
	<p class="text-sm text-[#64748B]">Mengalihkan ke halaman login...</p>
</div>
