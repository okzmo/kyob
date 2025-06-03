<script lang="ts">
	import { Dialog } from 'bits-ui';
	import { defaults, superForm } from 'sveltekit-superforms';
	import { valibot } from 'sveltekit-superforms/adapters';
	import { CreateChannelSchema } from 'types/schemas';
	import { core } from 'stores/core.svelte';
	import { backend } from 'stores/backend.svelte';
	import { page } from '$app/state';
	import CustomDialogContent from '../../CustomDialogContent/CustomDialogContent.svelte';
	import { delay } from '../../../../utils/delay';
	import CloseDialogButton from '../../CustomDialogContent/CloseDialogButton.svelte';
	import FooterDialog from '../../CustomDialogContent/FooterDialog.svelte';
	import SubmitButton from '../../SubmitButton/SubmitButton.svelte';
	import FormInput from '../../FormInput/FormInput.svelte';

	let isSubmitting = $state(false);
	let isSubmitted = $state(false);
	let buttonWidth = $derived(isSubmitted || isSubmitting ? 40 : 130);

	const { form, errors, enhance } = superForm(defaults(valibot(CreateChannelSchema)), {
		dataType: 'json',
		SPA: true,
		validators: valibot(CreateChannelSchema),
		async onUpdate({ form }) {
			if (form.valid) {
				const serverId = page.params.server_id;
				form.data.x = core.openCreateChannelModal.x;
				form.data.y = core.openCreateChannelModal.y;

				isSubmitting = true;
				const res = await backend.createChannel(serverId, form.data);
				if (res.isErr()) {
					if (res.error.code === 'ERR_VALIDATION_FAILED') {
						console.log(res.error.error);
					}
					if (res.error.code === 'ERR_UNAUTHORIZED') {
						console.log(res.error.error);
					}
					isSubmitting = false;
				}

				if (res.isOk()) {
					await delay(500);
					isSubmitting = false;
					isSubmitted = true;
					await delay(1000);

					core.openCreateChannelModal.status = false;
					core.activateMapDragging();

					isSubmitted = false;
				}
			}
		}
	});

	let isEmpty = $derived(!$form.name);
</script>

<Dialog.Root
	onOpenChange={(s) => {
		if (!s) core.activateMapDragging();
		core.openCreateChannelModal.status = s;
	}}
	open={core.openCreateChannelModal.status}
>
	<Dialog.Portal>
		<Dialog.Overlay class="fixed inset-0 bg-black/20" />
		<CustomDialogContent>
			<CloseDialogButton />
			<div class="flex items-center justify-between px-8">
				<div>
					<Dialog.Title class="text-lg font-semibold">Create a new channel</Dialog.Title>
					<Dialog.Description class="text-main-400 max-w-[24rem]  text-sm">
						Create channels to interact with peoples based on topics. What you like, dislike, memes,
						you choose.
					</Dialog.Description>
				</div>
			</div>

			<form method="post" use:enhance>
				<FormInput
					title="Channel name"
					id="channel-name"
					bind:error={$errors.name}
					bind:inputValue={$form.name}
					placeholder="General"
					type="text"
					class="mt-4 px-8"
				/>

				<FooterDialog>
					<SubmitButton
						type="submit"
						{buttonWidth}
						{isEmpty}
						{isSubmitting}
						{isSubmitted}
						class="relative"
					>
						Create channel
					</SubmitButton>
				</FooterDialog>
			</form>
		</CustomDialogContent>
	</Dialog.Portal>
</Dialog.Root>
