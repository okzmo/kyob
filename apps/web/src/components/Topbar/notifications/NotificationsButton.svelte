<script lang="ts">
	import { Popover } from 'bits-ui';
	import CustomPopoverContent from '../../ui/CustomPopoverContent/CustomPopoverContent.svelte';
	import Corners from '../../ui/Corners/Corners.svelte';
	import Bell from '../../ui/icons/Bell.svelte';
	import Tooltip from 'components/ui/Tooltip/Tooltip.svelte';

	let isOpen = $state(false);
</script>

<Popover.Root open={isOpen} onOpenChange={(s) => (isOpen = s)}>
	<Popover.Trigger
		aria-label="Notifications"
		class="top-bar-button group text-main-400 hocus:text-accent-50 hocus:bg-accent-100/15 relative flex h-[2.25rem] w-[2.25rem] items-center justify-center transition-colors hover:cursor-pointer"
	>
		<Tooltip text="Friends" y={-5} />
		<Corners color="border-main-300" class="group-hocus:border-accent-100 duration-100" />
		<Bell height={22} width={22} />
	</Popover.Trigger>
	<CustomPopoverContent
		class="bg-main-900 border-main-800 relative z-30 w-[20rem] rounded-2xl border p-2"
		align="end"
		side="bottom"
		sideOffset={10}
		y={-10}
	>
		Notifications :D
	</CustomPopoverContent>
</Popover.Root>
