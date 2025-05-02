<script lang="ts">
	import Gear from '../ui/icons/Gear.svelte';
	import Headphone from '../ui/icons/Headphone.svelte';
	import Microphone from '../ui/icons/Microphone.svelte';
	import IconButton from '../ui/IconButton/IconButton.svelte';
</script>

<div
	class="bg-main-900 border-main-800 fixed bottom-5 left-5 z-50 flex items-center gap-x-6 rounded-2xl border py-1 pr-2 pl-1 transition-colors duration-100 hover:cursor-pointer"
>
	<button
		class="group hocus:bg-accent-100/15 flex items-center gap-x-2.5 rounded-xl py-1 pr-4 pl-2 text-left transition-colors hover:cursor-pointer"
	>
		<img
			src="/images/cool-avatar-guy.jpg"
			alt="avatar"
			class="h-[2.75rem] w-[2.75rem] rounded-full"
		/>
		<div class="flex flex-col">
			<p class="group-hocus:text-accent-50 text-sm leading-[1.15rem] font-medium transition-colors">
				Okzmo
			</p>
			<p
				class="text-main-400 group-hocus:text-accent-200 text-sm leading-[1.15rem] transition-colors"
			>
				Connected
			</p>
		</div>
	</button>
	<div class="flex items-center gap-x-1">
		<IconButton Icon={Headphone} label="Headphone" />
		<IconButton Icon={Microphone} label="Microphone" />
		<IconButton Icon={Gear} href="/settings" label="Settings" />
	</div>
</div>
