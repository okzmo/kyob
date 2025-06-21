<script lang="ts">
	import { serversStore } from 'stores/servers.svelte';
	import { userStore } from 'stores/user.svelte';
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

	function getParticipant(id: string) {
		if (userStore.user?.id === id) return userStore.user;

		const friend = userStore.getFriend(id);
		if (friend) return friend;

		const member = serversStore.getMemberById(server.id, id);
		if (member) return member;

		return undefined;
	}
</script>

<div
	class="@container relative my-auto flex w-full flex-wrap items-center justify-center gap-2 px-4 pt-4 pb-20"
>
	{#if mainParticipant}
		{@const participantInfos = getParticipant(mainParticipant)}

		<button
			class="attachment relative aspect-[16/9] w-full"
			onclick={() => toggleMainParticipant()}
		>
			<img
				src={participantInfos?.banner}
				alt={participantInfos?.username}
				class="h-full w-full object-cover select-none"
			/>
		</button>
	{:else if channel.voice_users.length > 0}
		{#each channel.voice_users as participant (participant.user_id)}
			{@const participantInfos = getParticipant(participant.user_id)}

			<button
				class="attachment relative h-[15rem] w-[20rem]"
				onclick={() => toggleMainParticipant(participant.user_id)}
			>
				<img
					src={participantInfos?.banner}
					alt={participantInfos?.username}
					class="h-full w-full object-cover select-none"
				/>
			</button>
		{/each}
	{:else}
		<p class="text-main-500 text-lg font-bold">Nobody is in this channel</p>
	{/if}
</div>
