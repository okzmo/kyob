import * as v from 'valibot';
import { ChannelTypes } from './types';
import { userStore } from 'stores/user.svelte';

export const SignUpSchema = v.object({
	email: v.pipe(
		v.string(),
		v.nonEmpty('Please enter your email.'),
		v.email('The email is badly formatted.')
	),
	username: v.pipe(
		v.string(),
		v.minLength(1, 'The length must be equal or above 1 character.'),
		v.maxLength(20, 'The length must be equal or below 20 characters.'),
		v.nonEmpty('Please enter a username.')
	),
	display_name: v.pipe(
		v.string(),
		v.minLength(1, 'The length must be equal or above 1 character.'),
		v.maxLength(20, 'The length must be equal or below 20 characters.'),
		v.nonEmpty('Please enter a display name.')
	),
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
		v.maxLength(20, 'The length must be equal or below 20 characters.'),
		v.nonEmpty('Please enter a name for your realm.')
	),
	description: v.any(),
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
	}),
	x: v.number(),
	y: v.number()
});

export interface CreateServerType extends v.InferInput<typeof CreateServerSchema> {}

export const CreateChannelSchema = v.object({
	name: v.pipe(
		v.string(),
		v.minLength(1, 'The length must be equal or above 1 character.'),
		v.maxLength(20, 'The length must be equal or below 20 characters.'),
		v.nonEmpty('Please enter a name for this channel.')
	),
	type: v.pipe(v.enum(ChannelTypes)),
	x: v.number(),
	y: v.number()
});

export interface CreateChannelType extends v.InferInput<typeof CreateChannelSchema> {}

export const CreateMessageSchema = v.object({
	content: v.any(),
	everyone: v.optional(v.boolean()),
	mentions_users: v.optional(v.array(v.string())),
	mentions_channels: v.optional(v.array(v.string())),
	attachments: v.optional(v.array(v.pipe(v.file('Please select a valid file.'))))
});

export interface CreateMessageType extends v.InferInput<typeof CreateMessageSchema> {}

export const EditMessageSchema = v.object({
	content: v.any(),
	mentions_users: v.optional(v.array(v.string())),
	mentions_channels: v.optional(v.array(v.string()))
});

export interface EditMessageType extends v.InferInput<typeof EditMessageSchema> {}

export const JoinServerSchema = v.object({
	invite_url: v.string(),
	x: v.number(),
	y: v.number()
});

export interface JoinServerType extends v.InferInput<typeof JoinServerSchema> {}

export const AddFriendSchema = v.object({
	username: v.string()
});

export interface AddFriendType extends v.InferInput<typeof AddFriendSchema> {}

export const AcceptFriendSchema = v.object({
	friendship_id: v.string(),
	user_id: v.string(),
	friend_id: v.string()
});

export interface AcceptFriendType extends v.InferInput<typeof AcceptFriendSchema> {}

export const DeleteFriendSchema = v.object({
	friendship_id: v.string(),
	friend_id: v.string(),
	user_id: v.string()
});

export interface DeleteFriendType extends v.InferInput<typeof DeleteFriendSchema> {}

export const UpdateAccountSchema = v.object({
	email: v.optional(v.pipe(v.string(), v.email('The email is badly formatted.'))),
	username: v.optional(
		v.pipe(v.string(), v.maxLength(20, 'The length must be equal or below 20 characters.'))
	)
});

export interface UpdateAccountType extends v.InferInput<typeof UpdateAccountSchema> {}

export const UpdateProfileSchema = v.object({
	display_name: v.optional(
		v.pipe(
			v.string(),
			v.minLength(1, 'Minimum 1 character.'),
			v.maxLength(20, 'Maximum 20 characters.')
		),
		''
	),
	about: v.any(),
	links: v.optional(
		v.array(
			v.object({
				id: v.string(),
				label: v.pipe(v.string(), v.maxLength(20, 'Maximum 20 characters.')),
				url: v.pipe(v.string(), v.url())
			})
		),
		[]
	),
	facts: v.optional(
		v.array(
			v.object({
				id: v.string(),
				label: v.pipe(v.string()),
				value: v.pipe(v.string(), v.maxLength(20, 'Maximum 20 characters.'))
			})
		),
		[]
	)
});

export interface UpdateProfileType extends v.InferInput<typeof UpdateProfileSchema> {}

export const UpdateAvatarSchema = v.object({
	avatar: v.pipe(
		v.file('Please select an image file.'),
		v.mimeType(
			['image/jpeg', 'image/jpg', 'image/png', 'image/gif'],
			'Please select a JPEG, PNG or GIF file.'
		),
		v.maxSize(1024 * 1024 * 10, 'Please select a file smaller than 10 MB.')
	),
	crop_avatar: v.object({
		height: v.number(),
		width: v.number(),
		x: v.number(),
		y: v.number()
	}),
	crop_banner: v.object({
		height: v.number(),
		width: v.number(),
		x: v.number(),
		y: v.number()
	}),
	main_color: v.optional(v.string())
});

export interface UpdateAvatarType extends v.InferInput<typeof UpdateAvatarSchema> {}

export const UpdatePasswordSchema = v.object({
	current_password: v.pipe(v.string(), v.nonEmpty('Please enter your old password.')),
	new_password: v.pipe(
		v.string(),
		v.nonEmpty('Please enter a new password.'),
		v.minLength(8, 'This password is too short.'),
		v.maxLength(254, 'This password is too long.')
	),
	confirm: v.pipe(v.string(), v.nonEmpty('Please the same password as the new one.'))
});

export interface UpdatePasswordType extends v.InferInput<typeof UpdatePasswordSchema> {}
