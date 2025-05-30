<script lang="ts">
	import { Popover, Separator } from 'bits-ui';
	import { type Snippet } from 'svelte';
	import type { User } from '../../types/types';
	import LinkOutside from '../ui/icons/LinkOutside.svelte';
	import CustomPopoverContent from '../ui/CustomPopoverContent/CustomPopoverContent.svelte';
	import Corners from '../ui/Corners/Corners.svelte';

	interface Props {
		children: Snippet;
		user: User;
	}

	let { children, user }: Props = $props();
</script>

<Popover.Root>
	<Popover.Trigger>
		{@render children()}
	</Popover.Trigger>
	<CustomPopoverContent
		class="relative z-[999] w-[20rem] p-0"
		align="start"
		side="top"
		sideOffset={10}
		y={10}
	>
		<div class="relative z-[2] h-full overflow-hidden bg-[#153c45]">
			{#if user.avatar}
				<figure class="absolute top-0 left-0 z-[4] h-[14rem] w-full">
					<img
						src={user.avatar}
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
				{#if user?.about}
					<p class="text-main-50/80 mt-2">CEO of my own delusion</p>
				{/if}
				{#if user.facts || user.links}
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
								<LinkOutside
									height={20}
									width={20}
									class="absolute top-1/2 right-4 -translate-y-1/2"
								/>
							</a>
						{/each}
					{/if}
				{/if}
				{#if user.facts}
					<p class="text-main-50/65 mt-5 mb-2 text-sm font-semibold">Facts</p>
					{#each user.facts as link, idx (idx)}
						<div class="flex items-center gap-x-2">
							<span class="text-main-50/40">{link.label}</span>
							<span class="text-main-50 font-semibold">{link.value}</span>
						</div>
					{/each}
				{/if}
			</div>
		</div>
	</CustomPopoverContent>
</Popover.Root>

<style>
	:global(.user-profile-gradient) {
		content: '';
		position: absolute;
		inset: 0;
		z-index: 2;
		background: linear-gradient(
			180deg,
			rgba(21, 60, 69, 0) 50%,
			rgba(21, 60, 69, 0.85) 80%,
			rgba(21, 60, 69, 0.95) 90%,
			rgba(21, 60, 69, 1) 14rem,
			#153c45 100%
		);
	}
</style>
