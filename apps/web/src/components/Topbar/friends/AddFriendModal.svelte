<script lang="ts">
	import { Dialog } from 'bits-ui';
	import CustomDialogContent from '../../ui/CustomDialogContent/CustomDialogContent.svelte';
	import { defaults, superForm, setError } from 'sveltekit-superforms';
	import { valibot } from 'sveltekit-superforms/adapters';
	import { core } from '../../../stores/core.svelte';
	import { AddFriendSchema } from '../../../types/schemas';
	import { backend } from '../../../stores/backend.svelte';
	import CloseDialogButton from '../../ui/CustomDialogContent/CloseDialogButton.svelte';
	import FooterDialog from '../../ui/CustomDialogContent/FooterDialog.svelte';
	import SubmitButton from '../../ui/SubmitButton/SubmitButton.svelte';
	import FormInput from '../../ui/FormInput/FormInput.svelte';
	import { delay } from '../../../utils/delay';

	let isSubmitting = $state(false);
	let isSubmitted = $state(false);
	let buttonWidth = $derived(isSubmitted || isSubmitting ? 40 : 100);

	const { form, errors, enhance } = superForm(defaults(valibot(AddFriendSchema)), {
		dataType: 'json',
		SPA: true,
		validators: valibot(AddFriendSchema),
		async onUpdate({ form }) {
			if (form.valid) {
				isSubmitting = true;

				const res = await backend.addFriend({ username: form.data.username });
				if (res.isErr()) {
					if (res.error.code === 'ERR_USER_NOT_FOUND') {
						setError(form, 'username', 'User not found.');
					}

					if (res.error.code === 'ERR_ADDING_ITSELF') {
						setError(form, 'username', 'You cannot add yourself.');
					}

					if (res.error.code === 'ERR_UNKNOWN') {
						console.log(res.error.error);
					}

					isSubmitting = false;
				}

				if (res.isOk()) {
					await delay(1000);
					isSubmitting = false;
					isSubmitted = true;
					await delay(2000);

					core.openAddFriendModal.status = false;
					core.activateMapDragging();
				}
			}
		}
	});

	let isEmpty = $derived(!$form.username);
</script>

<Dialog.Root
	onOpenChange={(s) => {
		if (!s) core.activateMapDragging();
		core.openAddFriendModal.status = s;
	}}
	open={core.openAddFriendModal.status}
>
	<Dialog.Portal>
		<Dialog.Overlay class="fixed inset-0 bg-black/20 transition-opacity" />
		<CustomDialogContent>
			<CloseDialogButton />
			<div class="flex items-center justify-between px-8">
				<div>
					<Dialog.Title class="text-lg font-semibold">Add a new friend</Dialog.Title>
					<Dialog.Description class="text-main-400 max-w-[24rem]  text-sm">
						Who's that friend mmh :)
					</Dialog.Description>
				</div>
			</div>

			<form method="post" use:enhance>
				<FormInput
					title="Friend username"
					id="friend-name"
					type="text"
					bind:error={$errors.username}
					bind:inputValue={$form.username}
					placeholder="bob"
					class="mt-4 px-8"
				/>

				<FooterDialog>
					<SubmitButton
						type="submit"
						{buttonWidth}
						{isEmpty}
						{isSubmitting}
						{isSubmitted}
						class="absolute top-1/2 right-5 -translate-y-1/2"
					>
						Add friend
					</SubmitButton>
				</FooterDialog>
			</form>
		</CustomDialogContent>
	</Dialog.Portal>
</Dialog.Root>
