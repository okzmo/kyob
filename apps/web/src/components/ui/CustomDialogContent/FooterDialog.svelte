<script lang="ts">
	let { children } = $props();
</script>

<div class="border-t-main-800 relative mt-8 w-full border-t py-9">
	{@render children()}
</div>
