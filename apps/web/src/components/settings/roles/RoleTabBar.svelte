<script lang="ts">
	const TABS = ['display', 'permissions', 'members'];

	let { activeTab = $bindable(), activeRole } = $props();
</script>

<div class="flex gap-x-1.5">
	{#each TABS as tab, idx (idx)}
		<button
			type="button"
			class={[
				'border px-3 py-1.5 capitalize transition-colors duration-100 hover:cursor-pointer',
				!activeRole && tab === 'members' ? 'hidden' : '',
				activeTab === tab
					? 'bg-accent-100/15 border-accent-100 text-accent-50'
					: 'bg-main-900 border-main-800 text-main-50/50 hocus:bg-main-800 hocus:border-main-700'
			]}
			onclick={() => (activeTab = tab)}
		>
			{tab}
		</button>
	{/each}
</div>
