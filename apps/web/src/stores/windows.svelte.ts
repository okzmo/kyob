import type { Window } from '../types/types';

class Windows {
	openWindows = $state<Window[]>([]);
	activeWindow = $state<string | null>();
	lastActiveWindow = $state<string | null>();

	setActiveWindow(window: string | null) {
		if (this.activeWindow) this.lastActiveWindow = this.activeWindow;
		this.activeWindow = window;
	}

	reuseLastWindow() {
		this.activeWindow = this.lastActiveWindow;
	}

	getActiveWindow() {
		return this.openWindows.find((w) => w.id === this.activeWindow);
	}

	getWindow({ id, channelId, friendId }: { id?: string; channelId?: string; friendId?: string }) {
		return this.openWindows.find((w) => {
			if (id) {
				return w.id === id;
			}

			if (channelId) {
				return w.channelId === channelId;
			}

			if (friendId) {
				return w.friendId === friendId;
			}
		});
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

	closeDeadWindow(id: string) {
		this.openWindows = this.openWindows.filter((w) => w.id !== id);
	}
}

export const windows = new Windows();
