class Core {
	canDragMap = $state(true);
	offsetServerMap = $state({ x: 0, y: 0 });
	totalOffsetServerMap = $state({ x: 0, y: 0 });

	deactivateMapDragging() {
		this.canDragMap = false;
	}

	activateMapDragging() {
		this.canDragMap = true;
	}
}

export const core = new Core();
