<script lang="ts">
	import { page } from '$app/state';
	import { serversStore } from '../../stores/servers.svelte';
	import type { Server } from '../../types/types';
	import ServerbarMembers from './ServerbarMembers.svelte';

	let server = $state<Server>();

	$effect(() => {
		server = serversStore.getServer(Number(page.params.server_id));
	});
</script>

<div
	class="bg-main-900 border-main-800 fixed right-5 bottom-5 z-50 flex items-center gap-x-1 rounded-2xl border px-1 py-1 transition-colors duration-100 hover:cursor-pointer"
>
	<button
		class="group hocus:bg-accent-100/15 flex items-center gap-x-2.5 rounded-xl py-1 pr-2 pl-4 text-left transition-colors hover:cursor-pointer"
	>
		<div class="flex flex-col">
			<p
				class="group-hocus:text-accent-50 text-right text-sm leading-[1.15rem] font-medium transition-colors"
			>
				{server?.name}
			</p>
			<ServerbarMembers />
		</div>
		<img src={server?.background} alt="avatar" class="h-[2.75rem] w-[2.75rem] rounded-full" />
	</button>
</div>
