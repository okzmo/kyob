<script lang="ts">
	let { props, mentions_users = $bindable() } = $props();

	let selectedIndex = $state(0);
	let scrollableMenu = $state<HTMLDivElement>();

	export function handleKeyDown({ event }: { event: KeyboardEvent }) {
		if (event.key === 'ArrowUp') {
			handleArrowUp();
			return true;
		}

		if (event.key === 'ArrowDown') {
			handleArrowDown();
			return true;
		}

		if (event.key === 'Enter') {
			handleEnter();
			return true;
		}
	}

	function handleArrowUp() {
		selectedIndex = (selectedIndex + props.items.length - 1) % props.items.length;
		scrollToSelectedItem();
	}

	function handleArrowDown() {
		selectedIndex = (selectedIndex + 1) % props.items.length;
		scrollToSelectedItem();
	}

	function scrollToSelectedItem() {
		if (scrollableMenu) {
			const selectedItem = scrollableMenu.children[selectedIndex] as HTMLElement;
			if (selectedItem) {
				const scrollTop = scrollableMenu.scrollTop;
				const scrollBottom = scrollTop + scrollableMenu.clientHeight;
				const elementTop = selectedItem.offsetTop;
				const elementBottom = elementTop + selectedItem.offsetHeight;

				if (elementTop < scrollTop) {
					scrollableMenu.scrollTop = elementTop;
				} else if (elementBottom > scrollBottom) {
					scrollableMenu.scrollTop = elementBottom - scrollableMenu.clientHeight;
				}
			}
		}
	}

	function handleEnter() {
		selectItem(selectedIndex);
	}

	function selectItem(index: number) {
		const item = props.items[index];
		mentions_users.push(item.id);

		if (item) {
			props.command({ id: item.id, avatar: item.avatar, label: item.display_name });
		}
	}
</script>

{#if props.items.length > 0}
	<div
		bind:this={scrollableMenu}
		class="bg-main-900 border-main-800 absolute bottom-[4rem] left-[0.5rem] flex max-h-[20rem] w-[calc(100%-1rem)] flex-col gap-y-1 overflow-y-auto rounded-lg border px-1 py-1"
	>
		{#each props.items as item, idx (idx)}
			<button
				class={[
					'flex w-full items-center gap-x-1.5 rounded-md px-2 py-1 text-left',
					idx === selectedIndex ? 'bg-accent-100/20 text-accent-50' : 'hover:bg-accent-100/20'
				]}
				onclick={() => (selectedIndex = idx)}
			>
				<img src={item.avatar} alt="avatar" class="h-5 w-5 rounded-[50%] object-cover" />
				{item.display_name}
			</button>
		{/each}
	</div>
{/if}
