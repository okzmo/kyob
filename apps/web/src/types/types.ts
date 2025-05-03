import * as z from 'zod';

export interface Window {
	id: string;
	channelId: number;
	serverId: number;
}

export const ChannelSchema = z.interface({
	id: z.int(),
	name: z.string(),
	type: z.enum(['voice', 'textual']),
	x: z.int().optional(),
	y: z.int().optional(),
	unread: z.boolean().default(false)
});

export interface Channel extends z.infer<typeof ChannelSchema> {}

export const ServerSchema = z.interface({
	id: z.int(),
	name: z.string(),
	background: z.string(),
	description: z.string().optional(),
	x: z.int().optional(),
	y: z.int().optional()
});

export interface Server extends z.infer<typeof ServerSchema> {}

export const UserSchema = z.interface({
	id: z.int(),
	email: z.email(),
	username: z.string(),
	password: z.string().min(8),
	display_name: z.string(),
	avatar: z.string(),
	about: z.string().optional()
});

export interface User extends z.infer<typeof UserSchema> {}
