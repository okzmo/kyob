import type { User } from '../types/types';

class Core {
	canDragMap = $state(true);
	offsetServerMap = $state({ x: 0, y: 0 });
	totalOffsetServerMap = $state({ x: 0, y: 0 });
	openCreateChannelModal = $state({ status: false, x: 0, y: 0 });
	openCreateServerModal = $state({ status: false, x: 0, y: 0 });
	openAddFriendModal = $state({ status: false });
	openJoinServerModal = $state({ status: false, x: 0, y: 0 });
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
}

export const core = new Core();
