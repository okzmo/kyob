<script lang="ts">
	import { Dialog } from 'bits-ui';
	import CustomDialogContent from 'components/ui/CustomDialogContent/CustomDialogContent.svelte';
	import { errorsStore } from 'stores/errors.svelte';
</script>

<Dialog.Root
	open={errorsStore.attachmentError}
	onOpenChange={(s) => (errorsStore.attachmentError = s)}
>
	<Dialog.Overlay class="absolute inset-0 z-50 bg-black/20" />
	<CustomDialogContent class="max-w-[350px]">
		<div class="mb-8 flex items-center justify-between px-8">
			<div>
				<Dialog.Title class="text-lg font-semibold">File exceeds size limit!</Dialog.Title>
				<Dialog.Description class="text-main-400 max-w-[24rem] text-sm">
					Be careful, the size limit at your level is only 15mb. However the more active you are the
					more this limit will grow up to 500mb!
				</Dialog.Description>
			</div>
		</div>
	</CustomDialogContent>
</Dialog.Root>
