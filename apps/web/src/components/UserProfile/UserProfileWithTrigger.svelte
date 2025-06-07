<script lang="ts">
	import { Popover } from 'bits-ui';
	import { type Snippet } from 'svelte';
	import type { User } from 'types/types';
	import CustomPopoverContent from '../ui/CustomPopoverContent/CustomPopoverContent.svelte';
	import UserProfile from './UserProfile.svelte';

	interface Props {
		children: Snippet;
		user: User;
		align?: 'start' | 'center' | 'end';
		side?: 'top' | 'right' | 'bottom' | 'left';
		sideOffset?: number;
		alignOffset?: number;
		y?: number;
		x?: number;
	}

	let {
		children,
		user,
		align = 'start',
		side = 'top',
		sideOffset = 10,
		alignOffset = 0,
		y,
		x
	}: Props = $props();
</script>

<Popover.Root>
	<Popover.Trigger class="shrink-0">
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
		<UserProfile {user} />
	</CustomPopoverContent>
</Popover.Root>
