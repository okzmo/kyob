import type { User } from '../types/types';

class Core {
	canDragMap = $state(true);
	offsetServerMap = $state({ x: 0, y: 0 });
	totalOffsetServerMap = $state({ x: 0, y: 0 });
	openCreateChannelModal = $state({ status: false, x: 0, y: 0 });
	openCreateServerModal = $state({ status: false, x: 0, y: 0 });
	openJoinServerModal = $state({ status: false, x: 0, y: 0 });
	editingMessage = $state({ id: -1 });
	profiles = $state<User[]>([]);
	profileOpen = $state<{
		status: boolean;
		userId: number;
		element: HTMLElement | null;
	}>({
		status: false,
		userId: -1,
		element: null
	});

	deactivateMapDragging() {
		this.canDragMap = false;
	}

	activateMapDragging() {
		this.canDragMap = true;
	}

	openProfile(userId: number, element: HTMLElement) {
		if (!this.profileOpen.status && this.profileOpen.userId !== userId) {
			this.profileOpen = { status: true, userId, element };
		}
	}

	closeProfile() {
		this.profileOpen = { status: false, userId: -1, element: null };
	}

	startEditingMessage(messageId: number) {
		this.editingMessage.id = messageId;
	}

	stopEditingMessage() {
		this.editingMessage.id = -1;
	}
}

export const core = new Core();
