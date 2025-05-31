<script lang="ts">
	import { Popover, Separator } from 'bits-ui';
	import LinkOutside from '../ui/icons/LinkOutside.svelte';
	import CustomPopoverContent from '../ui/CustomPopoverContent/CustomPopoverContent.svelte';
	import { core } from '../../stores/core.svelte';
	import { backend } from '../../stores/backend.svelte';

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
				<div role="presentation" class="user-popover"></div>
				<div class="relative z-[2] overflow-hidden rounded-[6px]">
					{#if user.banner}
						<figure class="h-[10rem] w-full">
							<img
								src={user.banner}
								alt="{user.username}'s banner"
								class="h-full w-full object-cover"
							/>
						</figure>
					{:else}
						<div class="bg-main-900 h-[10rem] w-full"></div>
					{/if}

					<figure
						class="shadow-user-avatar absolute top-[7.25rem] left-6 h-[4.5rem] w-[4.5rem] overflow-hidden rounded-[50%]"
					>
						<img
							src={user.avatar}
							alt="{user.username}'s avatar"
							class="relative z-[1] h-full w-full rounded-[50%] object-cover"
						/>
					</figure>

					<div class="flex flex-col px-4 pt-10 pb-4">
						<h3 class="text-xl font-semibold">{user.display_name}</h3>
						<p class="text-sm leading-none text-white/65">{user.username}</p>
						{#if user.about}
							<p class="mt-2 text-white/80">CEO of my own delusion</p>
						{/if}
						{#if user.facts || user.links}
							<Separator.Root class="my-5 h-[1px] w-full bg-white/55" />
							{#if user.links}
								<p class="mb-2 text-sm font-semibold text-white/65">Links</p>
								{#each user.links as link, idx (idx)}
									<a
										href={link.url}
										class="hocus:bg-white/30 relative flex w-full flex-col rounded-2xl bg-white/15 px-4 py-2.5 transition-colors duration-100"
										target="_blank"
										rel="noreferrer noopener"
									>
										<span class="font-medium">{link.label}</span>
										<span class="text-sm text-white/65">{link.url}</span>
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
							<p class="text-sm font-semibold text-white/65">Facts</p>
							{#each user.facts as link, idx (idx)}
								<div>
									<span>{link.label}</span>
									<span>{link.value}</span>
								</div>
							{/each}
						{/if}
					</div>
				</div>
			{/if}
		{/await}
	</CustomPopoverContent>
</Popover.Root>

<style>
	:global(.gradient-user-profile) {
		background: linear-gradient(
			180deg,
			var(--user-gradient-top) 0%,
			var(--user-gradient-top) 66.35%,
			var(--user-gradient-bottom) 100%
		);
	}

	:global(.shadow-user-avatar) {
		box-shadow:
			0px 0px 0px 5px rgba(0, 0, 0, 0.25),
			0px 0px 0px 5px var(--user-gradient-top);
	}

	:global(.user-popover) {
		content: '';
		position: absolute;
		width: calc(100% - 8px);
		height: calc(100% - 8px);
		background-color: rgba(0, 0, 0, 0.25);
		box-shadow: inset 0px 0px 24px rgba(0, 0, 0, 0.25);
		border-radius: 6px;
	}
</style>
