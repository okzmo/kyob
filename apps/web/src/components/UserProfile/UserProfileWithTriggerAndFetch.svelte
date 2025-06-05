<script lang="ts">
	import { Popover } from 'bits-ui';
	import CustomPopoverContent from '../ui/CustomPopoverContent/CustomPopoverContent.svelte';
	import { core } from 'stores/core.svelte';
	import { backend } from 'stores/backend.svelte';
	import UserProfile from './UserProfile.svelte';
	import type { Snippet } from 'svelte';

	interface Props {
		children: Snippet;
		userId: string;
		align?: 'start' | 'center' | 'end';
		side?: 'top' | 'right' | 'bottom' | 'left';
		sideOffset?: number;
		alignOffset?: number;
		y?: number;
		x?: number;
	}

	let {
		children,
		userId,
		align = 'start',
		side = 'top',
		sideOffset = 10,
		alignOffset = 0,
		y,
		x
	}: Props = $props();

	let userProfile = $derived.by(async () => {
		const existingProfile = core.profiles.find((p) => p.id === userId);
		if (existingProfile) {
			document.documentElement.style.setProperty(
				'--user-color-85',
				`rgba(${existingProfile.main_color}, 0.85)`
			);
			document.documentElement.style.setProperty(
				'--user-color-95',
				`rgba(${existingProfile.main_color}, 0.95)`
			);
			document.documentElement.style.setProperty(
				'--user-color',
				`rgba(${existingProfile.main_color}, 1)`
			);
			return existingProfile;
		}

		const res = await backend.getUserProfile(userId);

		if (res.isErr()) {
			console.error(res.error);
			return;
		}

		if (res.isOk()) {
			document.documentElement.style.setProperty(
				'--user-color-85',
				`rgba(${res.value.main_color}, 0.85)`
			);
			document.documentElement.style.setProperty(
				'--user-color-95',
				`rgba(${res.value.main_color}, 0.95)`
			);
			document.documentElement.style.setProperty(
				'--user-color',
				`rgba(${res.value.main_color}, 1)`
			);
			//TODO: add a limit, LRU
			core.profiles.push(res.value);
			return res.value;
		}
	});
</script>

<Popover.Root>
	<Popover.Trigger>
		{@render children()}
	</Popover.Trigger>
	<CustomPopoverContent
		class="relative z-[999] w-[20rem] p-0"
		{align}
		{side}
		{sideOffset}
		{alignOffset}
		{y}
		{x}
	>
		{#await userProfile then user}
			{#if user}
				<UserProfile {user} />
			{/if}
		{/await}
	</CustomPopoverContent>
</Popover.Root>
