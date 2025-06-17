<script lang="ts">
	import { Separator } from 'bits-ui';
	import type { User } from 'types/types';
	import LinkOutside from '../ui/icons/LinkOutside.svelte';
	import Corners from '../ui/Corners/Corners.svelte';
	import { generateHTML } from '@tiptap/core';
	import StarterKit from '@tiptap/starter-kit';
	import { extractFirstNParagraphs, trimEmptyNodes } from 'utils/richInput';
	import { onMount } from 'svelte';
	import { isColorLight } from 'utils/colors';

	interface Props {
		user: User;
	}

	let { user }: Props = $props();

	let needDarkFontColor = $state(isColorLight(`rgb(${user?.main_color})`));
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

	<div
		class={[
			'relative z-[4] flex flex-col px-4 pt-[10.25rem] pb-4',
			needDarkFontColor ? 'inner-main-900/10' : 'inner-main-50/10'
		]}
	>
		<Corners color={needDarkFontColor ? 'border-main-900/35' : 'border-main-50/35'} />
		<h3 class={['text-xl font-semibold', needDarkFontColor && 'text-main-900']}>
			{user.display_name}
		</h3>
		<p class={['text-sm leading-none', needDarkFontColor ? 'text-main-900/65' : 'text-main-50/65']}>
			{user.username}
		</p>
		{#if aboutText}
			<div
				class={[
					'mt-2 [&>p]:min-h-[24px]',
					needDarkFontColor ? 'text-main-900/80' : 'text-main-50/80'
				]}
			>
				{@html aboutText.content}
			</div>
			{#if aboutText.enoughMatches}
				{#if !toggleAbout}
					<span>...</span>
				{/if}
				<button
					class={[
						'w-fit text-left text-sm transition-colors hover:cursor-pointer',
						needDarkFontColor ? 'hocus:text-main-900/75 ' : 'hocus:text-main-50/75'
					]}
					onclick={() => (toggleAbout = !toggleAbout)}
				>
					{toggleAbout ? 'Hide' : 'Show more'}
				</button>
			{/if}
		{/if}
		{#if user.facts.length > 0 || user.links.length > 0}
			<Separator.Root
				class={['my-5 h-[1px] w-full', needDarkFontColor ? 'bg-main-900/25' : 'bg-main-50/25']}
			/>
			{#if user.links}
				<p
					class={[
						'mb-2 text-sm font-semibold',
						needDarkFontColor ? 'text-main-900/65' : 'text-main-50/65'
					]}
				>
					Links
				</p>
				{#each user.links as link, idx (idx)}
					<a
						href={link.url}
						class={[
							'relative flex w-full flex-col px-4 py-2.5 transition-colors duration-100',
							needDarkFontColor
								? 'hocus:bg-main-900/20 bg-main-900/10 inner-main-900/10'
								: 'hocus:bg-main-50/20 bg-main-50/10 inner-main-50/10'
						]}
						target="_blank"
						rel="noreferrer noopener"
					>
						<span class={['font-medium', needDarkFontColor && 'text-main-900']}>{link.label}</span>
						<span class={['text-sm', needDarkFontColor ? 'text-main-900/65' : 'text-main-50/65']}
							>{link.url}</span
						>
						<LinkOutside
							height={20}
							width={20}
							class={[
								'absolute top-1/2 right-4 -translate-y-1/2',
								needDarkFontColor ? 'text-main-900' : ''
							]}
						/>
					</a>
				{/each}
			{/if}
			{#if user.facts.length > 0}
				<p
					class={[
						'mb-2 text-sm font-semibold',
						user.links.length > 0 && 'mt-5 ',
						needDarkFontColor ? 'text-main-900/65' : 'text-main-50/65'
					]}
				>
					Facts
				</p>
				{#each user.facts as link, idx (idx)}
					<div class="flex items-center gap-x-1">
						<span class={needDarkFontColor ? 'text-main-900/50' : 'text-main-50/50'}>
							{link.label}
						</span>
						<span class={['font-semibold', needDarkFontColor ? 'text-main-900' : 'text-main-50']}>
							{link.value}
						</span>
					</div>
				{/each}
			{/if}
		{/if}
	</div>
</div>
