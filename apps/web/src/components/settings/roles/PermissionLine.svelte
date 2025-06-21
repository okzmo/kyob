<script lang="ts">
	import Switch from 'components/ui/Switch/Switch.svelte';
	import type { AbilitiesType } from 'types/types';

	interface Props {
		label: string;
		description: string;
		ability: AbilitiesType;
		input: string[];
		active: boolean;
	}

	let { active, label, description, ability, input = $bindable() }: Props = $props();

	function addAbility(ability: AbilitiesType) {
		input = [...input, ability];
	}

	function removeAbility(ability: AbilitiesType) {
		input = input.filter((a) => a !== ability);
	}
</script>

<li
	class="border-main-800 flex w-full items-center justify-between border-b py-4 first:pt-5 last:border-b-transparent"
>
	<div class="flex flex-col gap-y-1">
		<p>{label}</p>
		<p class="text-main-500">{description}</p>
	</div>
	<Switch {active} action={() => addAbility(ability)} reverse={() => removeAbility(ability)} />
</li>
