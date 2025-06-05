<script lang="ts">
	import type { Server } from 'types/types';
	import Corners from '../ui/Corners/Corners.svelte';
	import { generateHTML } from '@tiptap/core';
	import StarterKit from '@tiptap/starter-kit';
	import { extractFirstNParagraphs, trimEmptyNodes } from 'utils/richInput';
	import { userStore } from 'stores/user.svelte';
	import UserProfileWithTriggerAndFetch from 'components/UserProfile/UserProfileWithTriggerAndFetch.svelte';

	interface Props {
		server: Server;
	}

	let { server }: Props = $props();

	let toggleAbout = $state(false);
	let aboutText = $derived.by(() => {
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

		if (!toggleAbout && enoughMatches) {
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

	<div class="inner-main-50/10 relative z-[4] flex flex-col px-4 pt-[10.25rem]">
		<Corners color="border-main-50/35" />
		<h3 class="text-xl font-semibold">{server.name}</h3>
		{#if aboutText}
			<div class="text-main-50/80 mt-2 [&>p]:min-h-[24px]">
				{@html aboutText.content}
			</div>
			{#if aboutText.enoughMatches}
				{#if !toggleAbout}
					<span>...</span>
				{/if}
				<button
					class="hocus:text-main-50/75 w-fit text-left text-sm transition-colors hover:cursor-pointer"
					onclick={() => (toggleAbout = !toggleAbout)}
				>
					{toggleAbout ? 'Hide' : 'Show more'}
				</button>
			{/if}
		{/if}
		<div class="mt-4 flex flex-col">
			<div class="mb-1 flex w-full items-center gap-x-2">
				<p class="text-main-50/50 text-sm">Members</p>
				<span class="bg-main-50/35 block h-[1px] w-full"></span>
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
							class="hocus:bg-main-50/10 border-main-50/0 hocus:border-main-50/20 relative flex w-full items-center gap-x-3 border py-2 pr-4 pl-2 transition-colors duration-100 hover:cursor-pointer"
						>
							<img
								src={member.avatar}
								alt="{member.username}'s avatar"
								class="h-8 w-8 transform-gpu object-cover"
							/>
							<p>{member.display_name}</p>
						</div>
					</UserProfileWithTriggerAndFetch>
				{/each}
			</div>
		</div>
	</div>
</div>
