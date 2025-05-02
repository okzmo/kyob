import type { Server } from '../types/types';

const SERVERS = [
	{
		id: 1,
		name: 'Minecraft community',
		serverBg: '/images/cool-avatar-guy.jpg',
		x: 1200,
		y: 748
	},
	{ id: 2, name: 'Memes', serverBg: '/images/cool-avatar-guy.jpg', x: 1200, y: 748 },
	{ id: 3, name: 'League of legends', serverBg: '/images/cool-avatar-guy.jpg', x: 200, y: 248 },
	{ id: 4, name: 'WAAAAZAAAA', serverBg: '/images/cool-avatar-guy.jpg', x: 1000, y: 500 },
	{ id: 5, name: 'Primeagen', serverBg: '/images/cool-avatar-guy.jpg', x: 1800, y: 108 },
	{ id: 6, name: 'Coding Lab', serverBg: '/images/cool-avatar-guy.jpg', x: 500, y: 848 },
	{ id: 7, name: 'Web world', serverBg: '/images/cool-avatar-guy.jpg', x: 1200, y: 48 }
];

class Servers {
	servers = $state<Server[]>(SERVERS);

	getServer(id: number) {
		return this.servers.find((s) => s.id === id);
	}
}

export const serversStore = new Servers();
