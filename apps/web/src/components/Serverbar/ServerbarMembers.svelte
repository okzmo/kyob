<div class="flex items-center gap-x-2.5">
	<p class="active-member flex items-center gap-x-1.5 text-sm text-green-300">123 active</p>
	<p
		class="global-member text-main-400 group-hocus:text-accent-200 before:bg-main-400 group-hocus:before:bg-accent-200! flex items-center gap-x-1.5 text-sm transition-colors duration-100 before:transition-colors before:duration-100"
	>
		2,283 members
	</p>
</div>

<style lang="postcss">
	.active-member::before {
		content: '';
		flex-shrink: 0;
		height: 5px;
		width: 5px;
		border-radius: 50%;
		background-color: var(--color-green-300);
	}

	.global-member::before {
		content: '';
		flex-shrink: 0;
		height: 5px;
		width: 5px;
		border-radius: 50%;
		background-color: var(--color-main-500);
	}
</style>
