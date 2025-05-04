export const contextMenuTargets = ['serverButton', 'channelButton', 'message'] as const;
export type ContextMenuTarget = (typeof contextMenuTargets)[number];

export interface Window {
	id: string;
	channelId: number;
	serverId: number;
}

export interface Channel {
	id: number;
	name: string;
	type: 'textual' | 'voice';
	x: number;
	y: number;
	unread: boolean;
}

export interface Server {
	id: number;
	owner_id: number;
	name: string;
	background: string;
	description?: string;
	x: number;
	y: number;
	is_member: boolean;
	channels?: Channel[];
}

export interface User {
	id: number;
	email: string;
	username: string;
	display_name: string;
	avatar: string;
	about?: string;
}

export interface Setup {
	user: User;
	servers: Server[];
}
