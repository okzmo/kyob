<script lang="ts">
	import PersonalForm from '../../../../components/settings/PersonalForm.svelte';
	import Corners from '../../../../components/ui/Corners/Corners.svelte';
	import SettingsSection from '../../../../components/ui/SettingsSection/SettingsSection.svelte';
</script>

<h1 class="text-2xl font-bold">Account</h1>

<div class="mt-10 flex flex-col gap-y-8">
	<SettingsSection
		title="Personal informations"
		description="The display name is what others see, username is what they use to add you as friend."
	>
		<PersonalForm />
	</SettingsSection>

	<SettingsSection title="Password">
		<button
			class="group inner-accent/15 hocus:bg-accent-100/25 hocus:inner-accent-no-shadow/25 bg-accent-100/15 text-accent-50 relative mt-2.5 block w-fit px-4 py-1 text-left transition duration-100 hover:cursor-pointer"
		>
			<Corners color="border-accent-100" />
			Change password
		</button>
	</SettingsSection>

	<SettingsSection
		title="Account removal"
		description="All your personal data including any messages you've ever sent will be deleted."
	>
		<button
			class="group inner-red-400/20 hocus:bg-red-400/25 hocus:inner-red-400/40 relative mt-3 block w-fit bg-red-400/15 px-4 py-1 text-left text-red-400 transition duration-100 hover:cursor-pointer"
		>
			<Corners color="border-red-400" />
			Delete account
		</button>
	</SettingsSection>
</div>
