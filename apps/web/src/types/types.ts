export interface Server {
	id: number;
	name: string;
	serverBg: string;
	x: number;
	y: number;
}

export interface Channel {
	id: number;
	name: string;
	type: 'voice' | 'text';
	x: number;
	y: number;
	unread: boolean;
}

export interface Window {
	id: string;
	channelId: number;
	serverId: number;
}
