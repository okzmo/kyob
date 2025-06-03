<script lang="ts">
	import { Dialog } from 'bits-ui';
	import { defaults, superForm } from 'sveltekit-superforms';
	import { valibot } from 'sveltekit-superforms/adapters';
	import { JoinServerSchema } from 'types/schemas';
	import { core } from 'stores/core.svelte';
	import { backend } from 'stores/backend.svelte';
	import CustomDialogContent from '../../CustomDialogContent/CustomDialogContent.svelte';
	import { serversStore } from 'stores/servers.svelte';
	import type { Server } from 'types/types';
	import { delay } from '../../../../utils/delay';
	import FooterDialog from '../../CustomDialogContent/FooterDialog.svelte';
	import SubmitButton from '../../SubmitButton/SubmitButton.svelte';
	import FormInput from '../../FormInput/FormInput.svelte';
	import CloseDialogButton from '../../CustomDialogContent/CloseDialogButton.svelte';

	let isSubmitting = $state(false);
	let isSubmitted = $state(false);
	let buttonWidth = $derived(isSubmitted || isSubmitting ? 40 : 100);

	const { form, errors, enhance } = superForm(defaults(valibot(JoinServerSchema)), {
		dataType: 'json',
		SPA: true,
		validators: valibot(JoinServerSchema),
		async onUpdate({ form }) {
			if (form.valid) {
				form.data.x = core.openJoinServerModal.x;
				form.data.y = core.openJoinServerModal.y;

				isSubmitting = true;
				const res = await backend.joinServer(form.data);
				if (res.isErr()) {
					if (res.error.code === 'ERR_VALIDATION_FAILED') {
						console.log(res.error.error);
					}
					isSubmitting = false;
				}

				if (res.isOk()) {
					const server: Server = {
						...res.value
					};
					serversStore.addServer(server);

					await delay(1000);
					isSubmitting = false;
					isSubmitted = true;
					await delay(2000);

					core.openJoinServerModal.status = false;
					core.activateMapDragging();
				}
			}
		}
	});

	let isEmpty = $derived(!$form.invite_url);
</script>

<Dialog.Root
	onOpenChange={(s) => {
		if (!s) core.activateMapDragging();
		core.openJoinServerModal.status = s;
	}}
	open={core.openJoinServerModal.status}
>
	<Dialog.Portal>
		<Dialog.Overlay class="fixed inset-0 bg-black/20" />
		<CustomDialogContent>
			<CloseDialogButton />

			<div class="flex items-center justify-between px-8">
				<div>
					<Dialog.Title class="text-lg font-semibold">Join a realm</Dialog.Title>
					<Dialog.Description class="text-main-400 max-w-[24rem]  text-sm">
						Join a realm to interact with peoples based on topics. What you like, dislike, memes,
						you choose.
					</Dialog.Description>
				</div>
			</div>

			<form method="post" use:enhance>
				<FormInput
					title="Invitation link"
					type="text"
					id="invite-url"
					bind:error={$errors.invite_url}
					bind:inputValue={$form.invite_url}
					placeholder="https://kyob.app/invite/123"
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
						Join realm
					</SubmitButton>
				</FooterDialog>
			</form>
		</CustomDialogContent>
	</Dialog.Portal>
</Dialog.Root>
