import type { Window } from '../types/types';

class Windows {
	openWindows = $state<Window[]>([]);
	activeWindow = $state<string | null>();

	createWindow(id: string, serverId: number, channelId: number) {
		this.openWindows.push({ id, serverId, channelId });
	}

	closeWindow(id: string) {
		this.openWindows = this.openWindows.filter((w) => w.id !== id);
	}
}

export const windows = new Windows();
