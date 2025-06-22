<script lang="ts">
	let { children } = $props();
</script>

<div
	class="border-t-main-800 relative mt-8 flex w-full items-center justify-end gap-x-2 border-t px-5 py-5"
>
	{@render children()}
</div>
