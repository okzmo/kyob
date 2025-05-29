import type { Window } from '../types/types';

class Windows {
	openWindows = $state<Window[]>([]);
	activeWindow = $state<string | null>();

	getActiveWindow() {
		return this.openWindows.find((w) => w.id === this.activeWindow);
	}

	createWindow({
		id,
		serverId,
		channelId,
		friendId
	}: {
		id: string;
		serverId?: string;
		channelId?: string;
		friendId?: string;
	}) {
		const exist = Boolean(this.openWindows.find((w) => w.id === id));
		if (exist) {
			this.activeWindow = id;
			return;
		}

		this.openWindows.push({ id, serverId, channelId, friendId });
		this.activeWindow = id;
	}

	closeWindow(id: string) {
		this.openWindows = this.openWindows.filter((w) => w.id !== id);
	}

	closeDeadWindow(channelId: string) {
		const exist = this.openWindows.find((w) => w.channelId === channelId);
		if (exist) {
			this.openWindows = this.openWindows.filter((w) => w.id !== exist.id);
		}
	}
}

export const windows = new Windows();
