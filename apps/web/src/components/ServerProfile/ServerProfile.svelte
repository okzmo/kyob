<script lang="ts">
	import type { Server } from 'types/types';
	import Corners from '../ui/Corners/Corners.svelte';
	import { generateHTML } from '@tiptap/core';
	import StarterKit from '@tiptap/starter-kit';
	import { extractFirstNParagraphs, trimEmptyNodes } from 'utils/richInput';
	import UserProfileWithTriggerAndFetch from 'components/UserProfile/UserProfileWithTriggerAndFetch.svelte';
	import { goto } from '$app/navigation';
	import Gear from 'components/ui/icons/Gear.svelte';
	import { isColorLight } from 'utils/colors';
	import { userStore } from 'stores/user.svelte';

	interface Props {
		server: Server;
	}

	let { server }: Props = $props();

	let isOwner = $derived(server.owner_id === userStore.user?.id);
	let needDarkFontColor = $state(isColorLight(`rgb(${server?.main_color})`));
	let toggleDescription = $state(false);
	let descriptionText = $derived.by(() => {
		if (!server.description) return;

		const html = generateHTML(trimEmptyNodes(server.description), [
			StarterKit.configure({
				gapcursor: false,
				dropcursor: false,
				heading: false,
				orderedList: false,
				bulletList: false,
				blockquote: false
			})
		]);

		const { paragraphs, enoughMatches } = extractFirstNParagraphs(html, 2);

		if (!toggleDescription && enoughMatches) {
			return { content: paragraphs, enoughMatches };
		}

		return { content: html, enoughMatches };
	});
</script>

<div class="relative z-[2] h-full overflow-hidden bg-[var(--server-color)] select-none">
	{#if server.member_count > 6}
		<div
			role="presentation"
			class="pointer-events-none absolute bottom-[1px] left-[1px] z-10 h-[5rem] w-[calc(100%-2px)] bg-gradient-to-t from-[var(--server-color)] to-transparent"
		></div>
	{/if}
	{#if server.avatar}
		<figure class="absolute top-0 left-0 z-[4] h-[14rem] w-full">
			<img
				src={server.banner || server.avatar}
				alt="{server.name}'s banner"
				class="h-full w-full transform-gpu object-cover"
			/>
			<div class="server-profile-gradient"></div>
		</figure>
	{:else}
		<div class="bg-main-700 h-[10rem] w-full"></div>
	{/if}

	{#if isOwner}
		<button
			class="bg-main-900/50 absolute top-2 right-2 z-10 flex h-[2rem] w-[2rem] items-center justify-center backdrop-blur-lg hover:cursor-pointer"
			aria-label="Server settings"
			onclick={() => goto(`/server-settings/${server.id}/profile`)}
		>
			<Gear height={20} width={20} />
		</button>
	{/if}

	<div class="inner-main-50/10 relative z-[4] flex flex-col px-4 pt-[10.25rem]">
		<Corners color="border-main-50/35" />
		<h3 class={['text-xl font-semibold', needDarkFontColor && 'text-main-900']}>
			{server.name}
		</h3>
		{#if descriptionText}
			<div
				class={[
					'mt-2 [&>p]:min-h-[24px]',
					needDarkFontColor ? 'text-main-900/80' : 'text-main-50/80'
				]}
			>
				{@html descriptionText.content}
			</div>
			{#if descriptionText.enoughMatches}
				{#if !toggleDescription}
					<span>...</span>
				{/if}
				<button
					class={[
						'w-fit text-left text-sm transition-colors hover:cursor-pointer',
						needDarkFontColor ? 'hocus:text-main-900/75 ' : 'hocus:text-main-50/75'
					]}
					onclick={() => (toggleDescription = !toggleDescription)}
				>
					{toggleDescription ? 'Hide' : 'Show more'}
				</button>
			{/if}
		{/if}
		<div class="mt-4 flex flex-col">
			<div class="mb-1 flex w-full items-center gap-x-2">
				<p class={['text-sm', needDarkFontColor ? 'text-main-900/50' : 'text-main-50/50']}>
					Members
				</p>
				<span
					class={['block h-[1px] w-full', needDarkFontColor ? 'bg-main-900/35' : 'bg-main-50/35']}
				>
				</span>
			</div>
			<div class="flex max-h-[20rem] flex-col gap-y-0.5 overflow-auto pb-3">
				{#each server.members as member (member.id)}
					<UserProfileWithTriggerAndFetch
						userId={member.id!}
						side="left"
						align="end"
						alignOffset={-20}
						x={5}
					>
						<div
							class={[
								'border-main-50/0 relative flex w-full items-center gap-x-3 border py-2 pr-4 pl-2 transition-colors duration-100 hover:cursor-pointer',
								needDarkFontColor
									? 'hocus:bg-main-900/10 hocus:border-main-900/20'
									: 'hocus:bg-main-50/10 hocus:border-main-50/20'
							]}
						>
							<img
								src={member.avatar}
								alt="{member.username}'s avatar"
								class="h-8 w-8 transform-gpu object-cover"
							/>
							<p class={needDarkFontColor ? 'text-main-900' : 'text-main-50'}>
								{member.display_name}
							</p>
						</div>
					</UserProfileWithTriggerAndFetch>
				{/each}
			</div>
		</div>
	</div>
</div>
