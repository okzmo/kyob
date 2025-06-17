<script lang="ts">
	import type { Fact, Link, User } from 'types/types';
	import Corners from '../ui/Corners/Corners.svelte';
	import { Dialog, Separator } from 'bits-ui';
	import LinkOutside from '../ui/icons/LinkOutside.svelte';
	import CustomDialogContent from '../ui/CustomDialogContent/CustomDialogContent.svelte';
	import Pen from '../ui/icons/Pen.svelte';
	import Cropper from 'svelte-easy-crop';
	import FooterDialog from '../ui/CustomDialogContent/FooterDialog.svelte';
	import SubmitButton from '../ui/SubmitButton/SubmitButton.svelte';
	import { backend } from 'stores/backend.svelte';
	import { userStore } from 'stores/user.svelte';
	import { delay } from 'utils/time';
	import { UpdateAvatarSchema } from 'types/schemas';
	import ColorThief, { type RGBColor } from 'colorthief';
	import * as v from 'valibot';
	import { generateHTML, type JSONContent } from '@tiptap/core';
	import StarterKit from '@tiptap/starter-kit';
	import { extractFirstNParagraphs, trimEmptyNodes } from 'utils/richInput';
	import Button from 'components/ui/Button/Button.svelte';
	import { isColorLight } from 'utils/colors';

	interface Props {
		user: User;
		displayName?: string;
		about?: JSONContent;
		links: Link[];
		facts: Fact[];
	}

	let {
		user,
		about = $bindable(),
		displayName = $bindable(),
		links = $bindable(),
		facts = $bindable()
	}: Props = $props();

	let isSubmitting = $state(false);
	let isSubmitted = $state(false);
	let buttonWidth = $derived(isSubmitted || isSubmitting ? 40 : 135);
	let toggleAbout = $state(false);
	let aboutText = $derived.by(() => {
		if (!about) return;

		const html = generateHTML(trimEmptyNodes(about || user.about), [
			StarterKit.configure({
				gapcursor: false,
				dropcursor: false,
				heading: false,
				orderedList: false,
				bulletList: false,
				blockquote: false
			})
		]);

		const { paragraphs, enoughMatches } = extractFirstNParagraphs(html, 2);

		if (!toggleAbout && enoughMatches) {
			return { content: paragraphs, enoughMatches };
		}

		return { content: html, enoughMatches };
	});

	let openImageModal = $state(false);
	let avatar = $state<string | undefined>();
	let image = $state<File | undefined>();
	let colors = $state<RGBColor[]>([]);
	let mainColor = $state<RGBColor | null>();
	let needDarkFontColor = $derived(isColorLight(`rgb(${userStore.user?.main_color})`));

	let cropBanner = $state({ x: 0, y: 0 });
	let cropAvatar = $state({ x: 0, y: 0 });
	let cropBannerPixels = $state({ x: 0, y: 0, height: 0, width: 0 });
	let cropAvatarPixels = $state({ x: 0, y: 0, height: 0, width: 0 });

	let zoomAvatar = $state(1);
	let zoomBanner = $state(1);

	let minZoomAvatar = $state(3);
	let minZoomBanner = $state(3);
	let maxZoomAvatar = $state(5);
	let maxZoomBanner = $state(5);

	function onFile(e: Event) {
		const target = e.target as HTMLInputElement;
		image = target.files?.[0];

		const colorThief = new ColorThief();

		if (image) {
			const dataUrl = URL.createObjectURL(image);
			const img = new Image();

			img.onload = () => {
				const aspectAvatar = 1;
				const aspectBanner = 320 / 224;
				const aspectImage = img.naturalWidth / img.naturalHeight;

				const importantColor = colorThief.getColor(img);
				const palette = colorThief.getPalette(img);
				if (importantColor) {
					colors.push(importantColor);
					mainColor = importantColor;
				}
				if (palette) colors.push(...palette);

				minZoomAvatar = Math.max(aspectAvatar / aspectImage, aspectImage / aspectAvatar);
				minZoomBanner = Math.max(aspectBanner / aspectImage, aspectImage / aspectBanner);

				zoomAvatar = minZoomAvatar;
				zoomBanner = minZoomBanner;

				URL.revokeObjectURL(dataUrl);
			};

			img.src = dataUrl;
			avatar = dataUrl;
		}
	}

	async function handleNewAvatar() {
		if (!image) return;
		const parsedData = v.parse(UpdateAvatarSchema, {
			avatar: image,
			crop_banner: cropBannerPixels,
			crop_avatar: cropAvatarPixels,
			main_color: mainColor?.join(',')
		});

		isSubmitting = true;
		const res = await backend.updateAvatar(parsedData);

		if (res.isErr()) {
			console.error(res.error.error);
			isSubmitting = false;
			return;
		}

		if (res.isOk()) {
			await delay(400);
			isSubmitting = false;
			isSubmitted = true;
			await delay(800);

			userStore.user!.avatar = res.value.avatar;
			userStore.user!.banner = res.value.banner;
			userStore.user!.main_color = res.value.main_color;
			openImageModal = false;

			isSubmitted = false;
		}
	}

	let isEmpty = $derived(!avatar || !image);
</script>

<Dialog.Root open={openImageModal} onOpenChange={(s) => (openImageModal = s)}>
	<Dialog.Overlay class="fixed inset-0 z-50 bg-black/20" />
	<CustomDialogContent>
		<div class="flex items-center justify-between px-8">
			<div>
				<Dialog.Title class="text-lg font-semibold">Change your avatar</Dialog.Title>
				<Dialog.Description class="text-main-400 max-w-[24rem] text-sm">
					It will be used as your profile banner and avatar.
				</Dialog.Description>
			</div>
		</div>

		<div class="mt-4 px-8">
			{#if avatar}
				<Button
					variants="danger"
					onclick={() => {
						avatar = undefined;
						image = undefined;
						mainColor = null;
						colors = [];
					}}
				>
					Remove avatar
				</Button>
			{:else}
				<label
					for="avatar-profile"
					class="group inner-accent/15 hocus:inner-accent-no-shadow/25 bg-accent-100/15 hover:bg-accent-100/25 text-accent-50 relative flex w-fit items-center justify-center overflow-hidden px-2 py-1 whitespace-nowrap transition duration-100"
				>
					<input
						type="file"
						id="avatar-profile"
						name="avatar-profile"
						aria-label="Profile avatar and banner"
						class="absolute h-full w-full text-transparent hover:cursor-pointer"
						onchange={onFile}
					/>
					<p>Choose an image</p>
				</label>
			{/if}

			{#if avatar}
				<div class="mt-3 flex gap-x-2">
					<div class="relative h-[224px] w-[320px]">
						<Cropper
							image={avatar}
							cropSize={{ height: 224, width: 320 }}
							cropShape="rect"
							showGrid={false}
							bind:crop={cropBanner}
							bind:zoom={zoomBanner}
							minZoom={minZoomBanner}
							maxZoom={maxZoomBanner}
							oncropcomplete={(e) => {
								cropBannerPixels = e.pixels;
							}}
						/>
					</div>

					<div class="relative h-[85px] w-[85px]">
						<Cropper
							image={avatar}
							cropSize={{ height: 85, width: 85 }}
							cropShape="rect"
							showGrid={false}
							bind:crop={cropAvatar}
							bind:zoom={zoomAvatar}
							minZoom={minZoomAvatar}
							maxZoom={maxZoomAvatar}
							oncropcomplete={(e) => {
								cropAvatarPixels = e.pixels;
							}}
						/>
					</div>
				</div>
				{#if colors.length > 0 && mainColor}
					<p class="mt-3 text-lg font-semibold">Main color</p>
					<div class="mt-1 flex gap-x-1">
						{#each colors as color, idx (idx)}
							<button
								class={[
									'h-7 w-10 border transition-colors hover:cursor-pointer',
									mainColor.join(',') === color.join(',')
										? 'border-main-50'
										: 'border-main-800 hocus:border-main-500'
								]}
								aria-label="color"
								style="background-color: rgb({color.join(',')});"
								onclick={() => (mainColor = color)}
							></button>
						{/each}
					</div>
				{/if}
			{/if}
		</div>

		<FooterDialog>
			<SubmitButton
				type="button"
				{buttonWidth}
				{isEmpty}
				{isSubmitting}
				{isSubmitted}
				class="relative"
				onclick={handleNewAvatar}
			>
				Use new avatar
			</SubmitButton>
		</FooterDialog>
	</CustomDialogContent>
</Dialog.Root>

<div
	class="relative z-[2] h-full w-[20rem] shrink-0 overflow-hidden bg-[var(--user-color)] select-none"
>
	<button
		class="bg-main-900/70 hocus:bg-main-800/70 text-main-400 hocus:text-main-50 absolute top-2 right-2 z-10 p-1.5 backdrop-blur-2xl transition duration-100 hover:cursor-pointer"
		onclick={() => (openImageModal = true)}
	>
		<Pen height={18} width={18} />
	</button>

	{#if user.avatar}
		<figure class="absolute top-0 left-0 z-[4] h-[14rem] w-full">
			<img
				src={user.banner}
				alt="{user.username}'s banner"
				class="h-full w-full transform-gpu object-cover"
			/>
			<div class="user-profile-gradient"></div>
		</figure>
	{:else}
		<div class="bg-main-700 h-[10rem] w-full"></div>
	{/if}

	<div class="inner-main-50/10 relative z-[4] flex flex-col px-4 pt-[10.25rem] pb-4">
		<Corners color="border-main-50/35" />
		<h3 class={['text-xl font-semibold', needDarkFontColor && 'text-main-900']}>{displayName}</h3>
		<p class={['text-sm leading-none', needDarkFontColor ? 'text-main-900/65' : 'text-main-50/65']}>
			{user.username}
		</p>
		{#if aboutText}
			<div
				class={[
					'mt-2 [&>p]:min-h-[24px]',
					needDarkFontColor ? 'text-main-900/80' : 'text-main-50/80'
				]}
			>
				{@html aboutText.content}
			</div>
			{#if aboutText.enoughMatches}
				{#if !toggleAbout}
					<span>...</span>
				{/if}
				<button
					class={[
						'w-fit text-left text-sm transition-colors hover:cursor-pointer',
						needDarkFontColor ? 'hocus:text-main-900/75 ' : 'hocus:text-main-50/75'
					]}
					onclick={() => (toggleAbout = !toggleAbout)}
				>
					{toggleAbout ? 'Hide' : 'Show more'}
				</button>
			{/if}
		{/if}
		{#if facts.length > 0 || links.length > 0}
			<Separator.Root
				class={['my-5 h-[1px] w-full', needDarkFontColor ? 'bg-main-900/25' : 'bg-main-50/25']}
			/>
			{#if links.length > 0}
				<p
					class={[
						'mb-2 text-sm font-semibold',
						needDarkFontColor ? 'text-main-900/65' : 'text-main-50/65'
					]}
				>
					Links
				</p>
				{#each links as link, idx (idx)}
					<a
						href={link.url}
						class={[
							'relative flex w-full flex-col px-4 py-2.5 transition-colors duration-100',
							needDarkFontColor
								? 'hocus:bg-main-900/20 bg-main-900/10 inner-main-900/10'
								: 'hocus:bg-main-50/20 bg-main-50/10 inner-main-50/10'
						]}
						target="_blank"
						rel="noreferrer noopener"
					>
						<span class={['font-medium', needDarkFontColor && 'text-main-900']}>{link.label}</span>
						<span class={['text-sm', needDarkFontColor ? 'text-main-900/65' : 'text-main-50/65']}
							>{link.url}</span
						>
						<LinkOutside
							height={20}
							width={20}
							class={[
								'absolute top-1/2 right-4 -translate-y-1/2',
								needDarkFontColor ? 'text-main-900' : ''
							]}
						/>
					</a>
				{/each}
			{/if}
			{#if facts.length > 0}
				<p
					class={[
						'mb-2 text-sm font-semibold',
						links.length > 0 && 'mt-5 ',
						needDarkFontColor ? 'text-main-900/65' : 'text-main-50/65'
					]}
				>
					Facts
				</p>
				{#each facts as link, idx (idx)}
					<div class="flex items-center gap-x-1">
						<span class={needDarkFontColor ? 'text-main-900/50' : 'text-main-50/50'}>
							{link.label}
						</span>
						<span class={['font-semibold', needDarkFontColor ? 'text-main-900' : 'text-main-50']}>
							{link.value}
						</span>
					</div>
				{/each}
			{/if}
		{/if}
	</div>
</div>
