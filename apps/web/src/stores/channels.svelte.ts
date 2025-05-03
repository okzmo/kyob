import type { Channel } from '../types/types';

const CHANNELS: Channel[] = [
	{ id: 1, name: 'memes', type: 'textual', x: 1200, y: 748, unread: false },
	{ id: 2, name: 'Général', type: 'textual', x: 200, y: 248, unread: false },
	{ id: 3, name: 'News', type: 'textual', x: 1000, y: 500, unread: true }
];

class Channels {
	channels = $state<Channel[]>(CHANNELS);

	getChannel(id: number) {
		return this.channels.find((c) => c.id === id);
	}
}

export const channelsStore = new Channels();
