<script lang="ts">
	interface Props {
		id: number;
		username: string;
		displayName: string;
		time: string;
	}
	let { id, displayName, username, time }: Props = $props();
</script>

<div
	id="message-{id}"
	class="hocus:bg-main-800/50 flex gap-x-3 px-4 py-2 transition-colors duration-100"
>
	<img
		src="/images/cool-avatar-guy.jpg"
		alt="{username}'s avatar"
		class="h-[3rem] w-[3rem] rounded-full"
	/>
	<div class="pt-1">
		<div class="flex items-baseline gap-x-3">
			<p class="text-sm font-semibold">{displayName}</p>
			<time class="text-main-600 text-xs">{time}</time>
		</div>
		<div class="mt-1 flex flex-col gap-y-1">
			<p>this is a text</p>
			<p>this is a text</p>
		</div>
	</div>
</div>
