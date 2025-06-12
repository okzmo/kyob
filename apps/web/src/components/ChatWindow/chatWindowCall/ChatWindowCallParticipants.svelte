<script lang="ts">
	import { serversStore } from 'stores/servers.svelte';
	import type { Channel, Server } from 'types/types';

	interface Props {
		server: Server;
		channel: Channel;
	}

	let { server, channel }: Props = $props();

	let mainParticipant = $state<string | undefined>();

	function toggleMainParticipant(participantId?: string) {
		if (channel.voice_users.length <= 1) return;
		mainParticipant = participantId;
	}
</script>

<div
	class="@container relative my-auto flex w-full flex-wrap items-center justify-center gap-2 px-4 pt-4 pb-20"
>
	{#if mainParticipant}
		{@const participantInfos = serversStore.getMemberById(server.id, mainParticipant)}
		<button
			class="attachment relative aspect-[4/3] min-h-0 max-w-full @lg:max-w-[45rem]"
			onclick={() => toggleMainParticipant()}
		>
			<img
				src={participantInfos?.banner || participantInfos?.avatar}
				alt=""
				class="h-full w-full object-cover select-none"
			/>
		</button>
	{:else}
		{#each channel.voice_users as participant (participant.user_id)}
			{@const participantInfos = serversStore.getMemberById(server.id, participant.user_id)}
			<button
				class="attachment relative aspect-[4/3] min-h-0 max-w-full @lg:max-w-[20rem]"
				onclick={() => toggleMainParticipant(participant.user_id)}
			>
				<img
					src={participantInfos?.banner || participantInfos?.avatar}
					alt=""
					class="h-full w-full object-cover select-none"
				/>
			</button>
		{/each}
	{/if}
</div>
