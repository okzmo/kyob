<script lang="ts">
	import { Separator } from 'bits-ui';
	import type { User } from 'types/types';
	import LinkOutside from '../ui/icons/LinkOutside.svelte';
	import Corners from '../ui/Corners/Corners.svelte';
	import { generateHTML } from '@tiptap/core';
	import StarterKit from '@tiptap/starter-kit';
	import { extractFirstNParagraphs, trimEmptyNodes } from 'utils/richInput';
	import { onMount } from 'svelte';

	interface Props {
		user: User;
	}

	let { user }: Props = $props();

	let toggleAbout = $state(false);
	let aboutText = $derived.by(() => {
		if (!user.about) return;

		const html = generateHTML(trimEmptyNodes(user.about), [
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

	onMount(() => {
		document.documentElement.style.setProperty(
			'--user-color-85',
			`rgba(${user?.main_color}, 0.85)`
		);
		document.documentElement.style.setProperty(
			'--user-color-95',
			`rgba(${user?.main_color}, 0.95)`
		);
		document.documentElement.style.setProperty('--user-color', `rgba(${user?.main_color}, 1)`);
	});
</script>

<div class="relative z-[2] h-full overflow-hidden bg-[var(--user-color)] select-none">
	{#if user.avatar}
		<figure class="absolute top-0 left-0 z-[4] h-[14rem] w-full">
			<img
				src={user.banner}
				alt="{user.username}'s banner"
				class="h-full w-full transform-gpu object-cover"
			/>
			<div class="user-profile-gradient"></div>
		</figure>
	{:else}
		<div class="bg-main-700 h-[10rem] w-full"></div>
	{/if}

	<div class="inner-main-50/10 relative z-[4] flex flex-col px-4 pt-[10.25rem] pb-4">
		<Corners color="border-main-50/35" />
		<h3 class="text-xl font-semibold">{user.display_name}</h3>
		<p class="text-main-50/65 text-sm leading-none">{user.username}</p>
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
		{#if user.facts.length > 0 || user.links.length > 0}
			<Separator.Root class="bg-main-50/25 my-5 h-[1px] w-full" />
			{#if user.links}
				<p class="text-main-50/65 mb-2 text-sm font-semibold">Links</p>
				{#each user.links as link, idx (idx)}
					<a
						href={link.url}
						class="hocus:bg-main-50/20 bg-main-50/10 inner-main-50/10 relative flex w-full flex-col px-4 py-2.5 transition-colors duration-100"
						target="_blank"
						rel="noreferrer noopener"
					>
						<span class="font-medium">{link.label}</span>
						<span class="text-main-50/65 text-sm">{link.url}</span>
						<LinkOutside height={20} width={20} class="absolute top-1/2 right-4 -translate-y-1/2" />
					</a>
				{/each}
			{/if}
			{#if user.facts.length > 0}
				<p class={['text-main-50/65 mb-2 text-sm font-semibold', user.links.length > 0 && 'mt-5 ']}>
					Facts
				</p>
				{#each user.facts as link, idx (idx)}
					<div class="flex items-center gap-x-2">
						<span class="text-main-50/50">{link.label}</span>
						<span class="text-main-50 font-semibold">{link.value}</span>
					</div>
				{/each}
			{/if}
		{/if}
	</div>
</div>
