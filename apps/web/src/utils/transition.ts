export function animateCoordinates(
	targetObject: { x: number; y: number },
	startCoordinates: { x: number; y: number },
	endCoordinates: { x: number; y: number },
	duration = 750,
	easingFunction = null
) {
	// Default easing function (easeInOutCubic)
	const easing =
		easingFunction ||
		function (t: number) {
			return t < 0.5 ? 4 * t * t * t : 1 - Math.pow(-2 * t + 2, 3) / 2;
		};

	const startTime = performance.now();

	return new Promise<void>((resolve) => {
		function animate(currentTime: number) {
			const elapsed = currentTime - startTime;
			const progress = Math.min(elapsed / duration, 1);

			const easedProgress = easing(progress);

			targetObject.x = startCoordinates.x + (endCoordinates.x - startCoordinates.x) * easedProgress;
			targetObject.y = startCoordinates.y + (endCoordinates.y - startCoordinates.y) * easedProgress;

			if (progress < 1) {
				requestAnimationFrame(animate);
			} else {
				resolve();
			}
		}

		requestAnimationFrame(animate);
	});
}
