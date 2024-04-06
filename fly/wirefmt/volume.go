package wirefmt

type ListVolumesOutput []Volume

type Volume struct {
	// "attached_alloc_id": "…",
	// "auto_backup_enabled": true,
	// "block_size": 1,
	// "blocks": 1,
	// "blocks_avail": 1,
	// "blocks_free": 1,
	// "created_at": "…",
	// "encrypted": true,
	// "fstype": "…",
	// "region": "…",
	// "snapshot_retention": 1,
	// "state": "…",
	// "zone": "…"
	ID                string `json:"id"`
	Name              string `json:"name"`
	Region            string `json:"region"`
	SizeGB            int    `json:"size_gb"`
	AttachedMachineID string `json:"attached_machine_id"`
	State             string `json:"state"`
}

var VolumeStatePendingDestroy = "pending_destroy"
var VolumeStateCreated = "created"

type CreateVolumeInput struct {
	Name   string `json:"name"`
	SizeGB int    `json:"size_gb"`
	Region string `json:"region"`
}

type CreateVolumeOutput Volume
