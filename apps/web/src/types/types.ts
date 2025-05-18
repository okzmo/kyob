export const ChannelTypes = {
	Textual: 'textual',
	Voice: 'voice'
} as const;
export type ChannelTypes = (typeof ChannelTypes)[keyof typeof ChannelTypes];

export const contextMenuTargets = [
	'serverButton',
	'channelButton',
	'message',
	'inServer',
	'mainMap'
] as const;
export type ContextMenuTarget = (typeof contextMenuTargets)[number];

export interface Window {
	id: string;
	channelId: number;
	serverId: number;
}

export interface Channel {
	id: number;
	name: string;
	type: ChannelTypes;
	x: number;
	y: number;
	unread: boolean;
	messages?: Message[];
}

export interface Server {
	id: number;
	owner_id: number;
	name: string;
	avatar: string;
	banner: string;
	description?: string;
	x: number;
	y: number;
	channels: Record<number, Channel>;
	active_count: number[];
	member_count: number;
	members: Partial<User>[];
}

export interface User {
	id: number;
	email: string;
	username: string;
	display_name: string;
	avatar: string;
	banner: string;
	gradient_top?: string;
	gradient_bottom?: string;
	about?: string;
	facts?: {
		label: string;
		value: string;
	}[];
	links?: {
		label: string;
		url: string;
	}[];
}

export interface Setup {
	user: User;
	servers: Record<number, Server>;
}

export interface DefaultResponse {
	message: string;
}

export interface Message {
	id: number;
	author: Partial<User>;
	server_id: number;
	channel_id: number;
	content: any;
	mentions_users: number[];
	mentions_channels: number[];
	created_at: string;
}

export type ActorMessageTypes = 'channel:message';
