<script lang="ts">
	import Gear from '../ui/icons/Gear.svelte';
	import Headphone from '../ui/icons/Headphone.svelte';
	import Microphone from '../ui/icons/Microphone.svelte';
	import { userStore } from 'stores/user.svelte';
	import Corners from '../ui/Corners/Corners.svelte';
	import UserProfileWithTrigger from '../UserProfile/UserProfileWithTrigger.svelte';
	import Button from 'components/ui/Button/Button.svelte';
	import { goto } from '$app/navigation';
</script>

<div
	class="bg-main-900/80 inner-shadow-main-800 fixed bottom-5 left-5 z-50 flex items-center gap-x-6 py-1 pr-2 pl-1 backdrop-blur-2xl transition-colors duration-100"
>
	<Corners color="border-main-700" />
	<UserProfileWithTrigger user={userStore.user!} y={10} alignOffset={-4}>
		<button
			class="group hocus:bg-accent-100/15 hocus:inner-accent/15 relative flex items-center gap-x-2.5 py-1 pr-4 pl-1 text-left transition select-none hover:cursor-pointer"
		>
			<img src={userStore.user?.avatar} alt="avatar" class="h-[2.75rem] w-[2.75rem] object-cover" />
			<div class="flex flex-col">
				<p
					class="group-hocus:text-accent-50 text-sm leading-[1.15rem] font-medium transition-colors"
				>
					{userStore.user?.display_name}
				</p>
				<p
					class="text-main-400 group-hocus:text-accent-200 text-sm leading-[1.15rem] transition-colors"
				>
					Connected
				</p>
			</div>
		</button>
	</UserProfileWithTrigger>
	<div class="flex items-center gap-x-1">
		<Button
			variants="icon"
			class={[
				'h-[2.25rem] w-[2.25rem]',
				userStore.mute
					? 'hocus:bg-red-400/25 text-red-400'
					: 'text-main-400 hocus:text-accent-50 hocus:bg-accent-100/15'
			]}
			onclick={() => userStore.toggleMute()}
			tooltip={userStore.mute ? 'Unmute' : 'Mute'}
			corners
			cornerColor="border-transparent"
			cornerClass={userStore.mute ? 'group-hocus:border-red-400' : 'group-hocus:border-accent-100'}
		>
			<Microphone height={22} width={22} mute={userStore.mute} />
		</Button>

		<Button
			variants="icon"
			class={[
				'h-[2.25rem] w-[2.25rem]',
				userStore.deafen
					? 'hocus:bg-red-400/25 text-red-400'
					: 'text-main-400 hocus:text-accent-50 hocus:bg-accent-100/15'
			]}
			onclick={() => userStore.toggleDeafen()}
			tooltip={userStore.deafen ? 'Undeafen' : 'Deafen'}
			corners
			cornerColor="border-transparent"
			cornerClass={userStore.deafen
				? 'group-hocus:border-red-400'
				: 'group-hocus:border-accent-100'}
		>
			<Headphone height={22} width={22} deafen={userStore.deafen} />
		</Button>

		<Button
			variants="icon"
			class="text-main-400 hocus:text-accent-50 hocus:bg-accent-100/15 h-[2.25rem] w-[2.25rem] !p-0"
			onclick={() => goto('/settings')}
			tooltip="Settings"
			corners
			cornerColor="border-transparent"
			cornerClass="group-hocus:border-accent-100"
		>
			<Gear height={22} width={22} />
		</Button>
	</div>
</div>
