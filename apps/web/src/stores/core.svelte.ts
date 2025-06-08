import type { User } from '../types/types';

class Core {
	canDragMap = $state(true);
	offsetServerMap = $state({ x: 0, y: 0 });
	totalOffsetServerMap = $state({ x: 0, y: 0 });
	openCreateChannelModal = $state({ status: false, x: 0, y: 0 });
	openCreateServerModal = $state({ status: false, x: 0, y: 0 });
	openAddFriendModal = $state({ status: false });
	openJoinServerModal = $state({ status: false, x: 0, y: 0 });
	openAttachmentsModal = $state({ status: false, idx: 0, attachments: [] });
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
}

export const core = new Core();
