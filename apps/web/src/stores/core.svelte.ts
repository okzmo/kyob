import type { User } from '../types/types';

class Core {
	canDragMap = $state(true);
	offsetServerMap = $state({ x: 0, y: 0 });
	totalOffsetServerMap = $state({ x: 0, y: 0 });
	createChannelModal = $state({ status: false, x: 0, y: 0 });
	createServerModal = $state({ status: false, x: 0, y: 0 });
	addFriendModal = $state({ status: false });
	joinServerModal = $state({ status: false, x: 0, y: 0 });
	attachmentsModal = $state({ status: false, idx: 0, attachments: [] });
	editingMessage = $state({ id: '' });
	profiles = $state<User[]>([]);
	profileOpen = $state<{
		status: boolean;
		userId: string;
		element: HTMLElement | null;
	}>({
		status: false,
		userId: '',
		element: null
	});

	deactivateMapDragging() {
		this.canDragMap = false;
	}

	activateMapDragging() {
		this.canDragMap = true;
	}

	openProfile(userId: string, element: HTMLElement) {
		if (!this.profileOpen.status && this.profileOpen.userId !== userId) {
			this.profileOpen = { status: true, userId, element };
		}
	}

	closeProfile() {
		this.profileOpen = { status: false, userId: '', element: null };
	}

	startEditingMessage(messageId: string) {
		this.editingMessage.id = messageId;
	}

	stopEditingMessage() {
		this.editingMessage.id = '';
	}

	modifyProfile(userId: string, user: Partial<User>) {
		const profile = this.profiles.find((p) => p.id === userId);
		if (!profile) return;

		if (user.avatar) profile.avatar = user.avatar;
		if (user.banner) profile.banner = user.banner;
		if (user.username) profile.username = user.username;
		if (user.display_name) profile.display_name = user.display_name;
		if (user.facts) profile.facts = user.facts;
		if (user.links) profile.links = user.links;
		if (user.about) profile.about = user.about;
		if (user.main_color) profile.main_color = user.main_color;
	}

	openAttachmentsModal(attachments: any, idx: number) {
		core.attachmentsModal.status = true;
		core.attachmentsModal.attachments = attachments;
		core.attachmentsModal.idx = idx;
	}

	closeAttachmentsModal() {
		core.attachmentsModal.status = false;
		core.attachmentsModal.idx = 0;
	}
}

export const core = new Core();
