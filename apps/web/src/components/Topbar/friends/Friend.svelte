<script lang="ts">
	import Check from '../../ui/icons/Check.svelte';
	import Close from '../../ui/icons/Close.svelte';
	import MoreIcon from '../../ui/icons/MoreIcon.svelte';

	let { id, displayName, avatar, accepted = false } = $props();
</script>

<button
	class="group hocus:bg-accent-100/15 flex w-full items-center gap-x-2.5 rounded-xl py-1.5 pr-4 pl-2 text-left transition-colors hover:cursor-pointer"
>
	<img src={avatar} alt="avatar" class="h-[2.75rem] w-[2.75rem] rounded-full object-cover" />
	<div class="flex flex-col">
		<p class="group-hocus:text-accent-50 leading-[1.15rem] font-medium transition-colors">
			{displayName}
		</p>
		<p
			class="text-main-400 group-hocus:text-accent-200 text-sm leading-[1.15rem] transition-colors"
		>
			Connected
		</p>
	</div>
</button>

{#if accepted}
	<button
		class="bg-main-900 hocus:bg-main-800 hocus:text-main-50 text-main-600 absolute top-1/2 right-2 -translate-y-1/2 rounded-full p-1.5 transition-colors duration-100 hover:cursor-pointer"
	>
		<MoreIcon height={20} width={20} />
	</button>
{:else}
	<div class="absolute top-1/2 right-2 flex items-center gap-x-2">
		<button
			class="hocus:bg-green-400/50 flex h-[2rem] w-[2rem] -translate-y-1/2 items-center justify-center rounded-lg border border-green-400 bg-green-400/20 text-green-400 transition-colors duration-100 hover:cursor-pointer"
		>
			<Check height={16} width={16} />
		</button>
		<button
			class="hocus:bg-red-400/50 flex h-[2rem] w-[2rem] -translate-y-1/2 items-center justify-center rounded-lg border border-red-400 bg-red-400/20 text-red-400 transition-colors duration-100 hover:cursor-pointer"
		>
			<Close height={16} width={16} />
		</button>
	</div>
{/if}
