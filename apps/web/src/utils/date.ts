export function formatMessageTime(time: string) {
	const now = new Date();
	const timestamp = new Date(time);

	const day = timestamp.getDate();
	const month = timestamp.getMonth();
	const year = timestamp.getFullYear().toString().slice(2);
	const hour = timestamp.getHours();
	const minutes = timestamp.getMinutes();

	if (now.toDateString() === timestamp.toDateString()) {
		return `${hour < 10 ? '0' + hour : hour}:${minutes < 10 ? '0' + minutes : minutes}`;
	}

	return `${day}/${month}/${year}, ${hour < 10 ? '0' + hour : hour}:${minutes < 10 ? '0' + minutes : minutes}`;
}
