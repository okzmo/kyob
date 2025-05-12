class Core {
	canDragMap = $state(true);
	offsetServerMap = $state({ x: 0, y: 0 });
	totalOffsetServerMap = $state({ x: 0, y: 0 });
	openCreateChannelModal = $state({ status: false, x: 0, y: 0 });
	openCreateServerModal = $state({ status: false, x: 0, y: 0 });
	openJoinServerModal = $state({ status: false, x: 0, y: 0 });

	deactivateMapDragging() {
		this.canDragMap = false;
	}

	activateMapDragging() {
		this.canDragMap = true;
	}
}

export const core = new Core();
