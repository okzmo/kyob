import * as v from 'valibot';
import { ChannelTypes } from './types';

export const SignUpSchema = v.object({
	email: v.pipe(
		v.string(),
		v.nonEmpty('Please enter your email.'),
		v.email('The email is badly formatted.')
	),
	username: v.pipe(v.string(), v.nonEmpty('Please enter a username.')),
	display_name: v.pipe(v.string(), v.nonEmpty('Please enter a display name.')),
	password: v.pipe(
		v.string(),
		v.nonEmpty('Please enter a password.'),
		v.minLength(8, 'This password is too short.'),
		v.maxLength(254, 'This password is too long.')
	)
});

export const SignInSchema = v.object({
	username: v.pipe(v.string(), v.nonEmpty('Please enter your email or username.')),
	password: v.pipe(v.string(), v.nonEmpty('Please enter your password.'))
});

export const CreateServerSchema = v.object({
	name: v.pipe(
		v.string(),
		v.maxLength(50, 'The length must be equal or below 50 characters.'),
		v.nonEmpty('Please enter a name for your realm.')
	),
	description: v.pipe(
		v.string(),
		v.maxLength(280, 'The length must be equal or below 280 characters.')
	),
	avatar: v.pipe(
		v.file('Please select an image file.'),
		v.mimeType(['image/jpeg', 'image/png'], 'Please select a JPEG or PNG file.'),
		v.maxSize(1024 * 1024 * 10, 'Please select a file smaller than 10 MB.')
	),
	private: v.boolean(),
	crop: v.object({
		height: v.number(),
		width: v.number(),
		x: v.number(),
		y: v.number()
	})
});

export interface CreateServerType extends v.InferInput<typeof CreateServerSchema> {}

export const CreateChannelSchema = v.object({
	name: v.pipe(
		v.string(),
		v.maxLength(50, 'The length must be equal or below 50 characters.'),
		v.nonEmpty('Please enter a name for this channel.')
	),
	type: v.pipe(v.enum(ChannelTypes)),
	x: v.number(),
	y: v.number()
});

export interface CreateChannelType extends v.InferInput<typeof CreateChannelSchema> {}

export const CreateMessageSchema = v.object({
	author_id: v.number(),
	content: v.any(),
	mentions_users: v.optional(v.array(v.number())),
	mentions_channels: v.optional(v.array(v.number()))
});

export interface CreateMessageType extends v.InferInput<typeof CreateMessageSchema> {}
