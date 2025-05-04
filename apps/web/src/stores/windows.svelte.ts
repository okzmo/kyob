import type { Window } from '../types/types';

class Windows {
	openWindows = $state<Window[]>([]);
	activeWindow = $state<string | null>();

	createWindow(id: string, serverId: number, channelId: number) {
		const exist = Boolean(this.openWindows.find((w) => w.id === id));
		if (exist) {
			this.activeWindow = id;
			return;
		}

		this.openWindows.push({ id, serverId, channelId });
		this.activeWindow = id;
	}

	closeWindow(id: string) {
		this.openWindows = this.openWindows.filter((w) => w.id !== id);
	}
}

export const windows = new Windows();
