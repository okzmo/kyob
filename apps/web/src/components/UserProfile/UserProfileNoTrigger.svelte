<script lang="ts">
	import { Popover } from 'bits-ui';
	import CustomPopoverContent from '../ui/CustomPopoverContent/CustomPopoverContent.svelte';
	import { core } from 'stores/core.svelte';
	import { backend } from 'stores/backend.svelte';
	import UserProfile from './UserProfile.svelte';

	let userProfile = $derived.by(async () => {
		if (core.profileOpen.status) {
			const existingProfile = core.profiles.find((p) => p.id === core.profileOpen.userId);
			if (existingProfile) {
				return existingProfile;
			}

			const res = await backend.getUserProfile(core.profileOpen.userId);

			if (res.isErr()) {
				console.error(res.error);
				core.closeProfile();
				return;
			}

			if (res.isOk()) {
				document.documentElement.style.setProperty('--user-color-85', '#153c45d9');
				document.documentElement.style.setProperty('--user-color-95', '#153c45f2');
				document.documentElement.style.setProperty('--user-color', '#153c45');
				core.profiles.push(res.value);
				return res.value;
			}
		}
	});
</script>

<Popover.Root
	open={core.profileOpen.status}
	onOpenChange={(s) => {
		if (!s) {
			setTimeout(() => {
				core.closeProfile();
			}, 200);
		}
	}}
>
	<Popover.Trigger />
	<CustomPopoverContent
		class="gradient-user-profile relative w-[20rem] rounded-lg p-1"
		align="start"
		side="top"
		sideOffset={10}
		y={10}
		customAnchor={core.profileOpen.element}
	>
		{#await userProfile then user}
			{#if user}
				<UserProfile {user} />
			{/if}
		{/await}
	</CustomPopoverContent>
</Popover.Root>
