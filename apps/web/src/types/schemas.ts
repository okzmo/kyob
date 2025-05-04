import * as v from 'valibot';

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
