import type { User } from '../types/types';

class UserStore {
	user = $state<User>();
}

export const userStore = new UserStore();
