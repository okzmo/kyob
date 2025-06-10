<script lang="ts">
	import { goto } from '$app/navigation';
	import { backend } from 'stores/backend.svelte';
	import Corners from '../Corners/Corners.svelte';
	import { page } from '$app/state';
	import { windows } from 'stores/windows.svelte';

	interface Props {
		type: 'general' | 'server';
	}

	let { type = 'general' }: Props = $props();

	interface Sections {
		[type: string]: {
			[section: string]: {
				label: string;
				href: string;
				locked: boolean;
			}[];
		};
	}

	const SECTIONS: Sections = {
		general: {
			'Personal Settings': [
				{ label: 'Account', href: '/settings/account', locked: false },
				{ label: 'Profile', href: '/settings/profile', locked: false },
				{ label: 'Data & Privacy', href: '/settings/privacy', locked: false }
			],
			'App Settings': [
				{ label: 'Appearance', href: '/settings/appearance', locked: true },
				{ label: 'Accessibility', href: '/settings/accessibility', locked: true },
				{ label: 'Voice & Video', href: '/settings/call', locked: false },
				{ label: 'Language', href: '/settings/language', locked: true }
			],
			Other: [{ label: 'Changelog', href: '/settings/changelog', locked: true }]
		},
		server: {}
	};
</script>

<aside class="mt-20 flex h-screen w-[20rem] flex-col gap-y-6 overflow-auto select-none">
	<button
		class={[
			'group text-main-300 hocus:bg-main-800 hocus:text-main-50 hocus:inner-main-700 relative block w-full px-3 py-1 text-left transition duration-100'
		]}
		onclick={() => history.back()}
	>
		<Corners color="border-main-300" class="group-hocus:border-main-200" />
		Go back
	</button>
	{#each Object.entries(SECTIONS[type]) as section, idx (idx)}
		<ul>
			<h2 class="text-main-500 mb-1 text-xs font-bold uppercase">{section[0]}</h2>
			{#each section[1] as link, idx (idx)}
				<li class="mt-2.5">
					<a
						data-sveltekit-replacestate
						href={link.locked ? '#' : link.href}
						class={[
							'group relative block w-full px-3 py-1 transition duration-100',
							link.locked && 'opacity-50 hover:cursor-not-allowed',
							page.url.pathname === link.href
								? 'bg-main-800 text-main-50 inner-main-700'
								: 'text-main-300 hocus:bg-main-800 hocus:text-main-50 hocus:inner-main-700 '
						]}
					>
						<Corners color="border-main-300" class="group-hocus:border-main-200" />
						{link.label}
					</a>
				</li>
			{/each}
		</ul>
	{/each}

	{#if type === 'general'}
		<button
			class="group hocus:bg-red-400/20 hocus:inner-red-400/40 hocus:inner-main-700 relative mt-5 block w-full px-3 py-1 text-left text-red-400 transition duration-100 hover:cursor-pointer"
			onclick={async () => {
				await backend.logout();
				goto('/signin');
			}}
		>
			<Corners color="border-red-400" />
			Log out
		</button>
	{/if}
</aside>
