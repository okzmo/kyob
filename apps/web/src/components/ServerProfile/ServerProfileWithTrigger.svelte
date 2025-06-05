<script lang="ts">
	import { Popover } from 'bits-ui';
	import { type Snippet } from 'svelte';
	import type { Server } from 'types/types';
	import CustomPopoverContent from '../ui/CustomPopoverContent/CustomPopoverContent.svelte';
	import ServerProfile from './ServerProfile.svelte';

	interface Props {
		children: Snippet;
		server: Server;
		align?: 'start' | 'center' | 'end';
		side?: 'top' | 'right' | 'bottom' | 'left';
		sideOffset?: number;
		alignOffset?: number;
		y?: number;
		x?: number;
	}

	let {
		children,
		server,
		alignOffset = 0,
		align = 'start',
		side = 'top',
		sideOffset = 10,
		y,
		x
	}: Props = $props();
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
		{y}
		{x}
		{alignOffset}
	>
		<ServerProfile {server} />
	</CustomPopoverContent>
</Popover.Root>
