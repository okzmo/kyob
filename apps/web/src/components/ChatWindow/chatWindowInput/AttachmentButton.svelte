<script lang="ts">
	import Plus from 'components/ui/icons/Plus.svelte';

	interface Props {
		attachments: File[];
	}

	let { attachments = $bindable() }: Props = $props();

	function onFile(e: Event) {
		const target = e.target as HTMLInputElement;
		if (target.files) {
			attachments = [...attachments, ...target.files];
		}
	}
</script>

<label
	class="text-main-600 hocus:text-main-200 absolute top-4.5 left-4 z-[1] transition-colors duration-100 hover:cursor-pointer"
	for="file-attachement"
>
	<input
		id="file-attachement"
		type="file"
		class="absolute h-0 w-0 text-transparent opacity-0"
		onchange={onFile}
		multiple
	/>
	<Plus height={20} width={20} />
</label>
