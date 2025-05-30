<script lang="ts">
	import Gear from '../ui/icons/Gear.svelte';
	import Headphone from '../ui/icons/Headphone.svelte';
	import Microphone from '../ui/icons/Microphone.svelte';
	import IconButton from '../ui/IconButton/IconButton.svelte';
	import { userStore } from '../../stores/user.svelte';
	import UserProfile from '../UserProfile/UserProfile.svelte';
	import Corners from '../ui/Corners/Corners.svelte';
</script>

<div
	class="bg-main-900/80 inner-shadow-main-800 fixed bottom-5 left-5 z-50 flex items-center gap-x-6 py-1 pr-2 pl-1 backdrop-blur-2xl transition-colors duration-100"
>
	<Corners color="border-main-700" />
	<UserProfile user={userStore.user!}>
		<button
			class="group hocus:bg-accent-100/15 hocus:inner-accent/15 relative flex items-center gap-x-2.5 py-1 pr-4 pl-1 text-left transition hover:cursor-pointer"
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
	</UserProfile>
	<div class="flex items-center gap-x-1">
		<IconButton Icon={Headphone} label="Headphone" />
		<IconButton Icon={Microphone} label="Microphone" />
		<IconButton Icon={Gear} href="/settings" label="Settings" />
	</div>
</div>
