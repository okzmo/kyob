<script lang="ts">
	import { Dialog } from 'bits-ui';
	import { backend } from 'stores/backend.svelte';
	import { core } from 'stores/core.svelte';
	import { defaults, setError, superForm } from 'sveltekit-superforms';
	import { valibot } from 'sveltekit-superforms/adapters';
	import { AddFriendSchema } from 'types/schemas';
	import { delay } from 'utils/delay';
	import CustomDialogContent from 'components/ui/CustomDialogContent/CustomDialogContent.svelte';
	import FooterDialog from 'components/ui/CustomDialogContent/FooterDialog.svelte';
	import FormInput from 'components/ui/FormInput/FormInput.svelte';
	import SubmitButton from 'components/ui/SubmitButton/SubmitButton.svelte';

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
						console.error(res.error.error);
					}

					isSubmitting = false;
				}

				if (res.isOk()) {
					await delay(400);
					isSubmitting = false;
					isSubmitted = true;
					await delay(800);

					core.addFriendModal.status = false;
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
		core.addFriendModal.status = s;
	}}
	open={core.addFriendModal.status}
>
	<Dialog.Portal>
		<Dialog.Overlay class="fixed inset-0 bg-black/20 transition-opacity" />
		<CustomDialogContent>
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
						class="relative"
					>
						Add friend
					</SubmitButton>
				</FooterDialog>
			</form>
		</CustomDialogContent>
	</Dialog.Portal>
</Dialog.Root>
