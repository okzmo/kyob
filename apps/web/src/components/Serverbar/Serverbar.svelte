<script lang="ts">
	import { page } from '$app/state';
	import { serversStore } from 'stores/servers.svelte';
	import type { Server } from 'types/types';
	import Corners from '../ui/Corners/Corners.svelte';
	import ServerbarMembers from './ServerbarMembers.svelte';
	import ServerProfileWithTrigger from 'components/ServerProfile/ServerProfileWithTrigger.svelte';
	import Button from 'components/ui/Button/Button.svelte';
	import { windows } from 'stores/windows.svelte';
	import Planet from 'components/ui/icons/Planet.svelte';

	let server = $derived<Server>(serversStore.getServer(page.params.server_id));
	let activeMembers = $derived(serversStore.getActiveMembers(page.params.server_id));
</script>

<div class="fixed right-5 bottom-5 flex items-center gap-x-2">
	<Button
		variants="nostyle"
		onclick={() => {
			windows.createWindow({
				id: `window-${server.id}-world`,
				serverId: server.id,
				type: 'world'
			});
		}}
		class="group bg-main-900/80 inner-shadow-main-800 hocus:bg-main-800/80 hocus:inner-shadow-main-500 text-main-500 hocus:text-main-300 relative flex aspect-square h-[3.75rem] items-center justify-center transition hover:cursor-pointer"
		tooltip="Join world"
	>
		<Corners color="border-main-700" class="group-hocus:border-main-300" />
		<Planet height={32} width={32} />
	</Button>
	<div
		class="bg-main-900/80 inner-shadow-main-800 z-50 flex items-center gap-x-6 p-1 backdrop-blur-2xl transition-colors duration-100 select-none"
	>
		<Corners color="border-main-700" />
		<ServerProfileWithTrigger {server} y={10} align="end" alignOffset={-4}>
			<button
				class="group hocus:bg-accent-100/15 hocus:inner-accent/15 relative flex items-center gap-x-2.5 py-1 pr-1 pl-4 text-left transition hover:cursor-pointer"
			>
				<div class="flex flex-col">
					<p
						class="group-hocus:text-accent-50 text-right text-sm leading-[1.15rem] font-medium transition-colors"
					>
						{server?.name}
					</p>
					<ServerbarMembers totalMembers={server?.member_count || 0} {activeMembers} />
				</div>
				<img src={server?.avatar} alt="avatar" class="h-[2.75rem] w-[2.75rem]" />
			</button>
		</ServerProfileWithTrigger>
	</div>
</div>
