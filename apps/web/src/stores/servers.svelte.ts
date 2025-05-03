import type { Server } from '../types/types';

const SERVERS = [
	{
		id: 1,
		name: 'Minecraft community',
		background: '/images/cool-avatar-guy.jpg',
		x: 1200,
		y: 748,
		description: ''
	},
	{
		id: 2,
		name: 'Memes',
		background: '/images/cool-avatar-guy.jpg',
		x: 1200,
		y: 748,
		description: ''
	},
	{
		id: 3,
		name: 'League of legends',
		background: '/images/cool-avatar-guy.jpg',
		x: 200,
		y: 248,
		description: ''
	},
	{
		id: 4,
		name: 'WAAAAZAAAA',
		background: '/images/cool-avatar-guy.jpg',
		x: 1000,
		y: 500,
		description: ''
	},
	{
		id: 5,
		name: 'Primeagen',
		background: '/images/cool-avatar-guy.jpg',
		x: 1800,
		y: 108,
		description: ''
	},
	{
		id: 6,
		name: 'Coding Lab',
		background: '/images/cool-avatar-guy.jpg',
		x: 500,
		y: 848,
		description: ''
	},
	{
		id: 7,
		name: 'Web world',
		background: '/images/cool-avatar-guy.jpg',
		x: 1200,
		y: 48,
		description: ''
	}
];

class Servers {
	servers = $state<Server[]>(SERVERS);

	getServer(id: number) {
		return this.servers.find((s) => s.id === id);
	}
}

export const serversStore = new Servers();
