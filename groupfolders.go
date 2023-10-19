package gonextcloud

type groupFolderBadFormatIDAndGroups struct {
	ID         int            `json:"id"`
	MountPoint string         `json:"mount_point"`
	Groups     map[string]int `json:"groups"`
	Quota      int            `json:"quota"`
	Size       int            `json:"size"`
}

type groupFolderBadFormatGroups struct {
	ID         int            `json:"id"`
	MountPoint string         `json:"mount_point"`
	Groups     map[string]int `json:"groups"`
	Quota      int            `json:"quota"`
	Size       int            `json:"size"`
}

// GroupFolder is group shared folder from groupfolders application
type GroupFolder struct {
	ID         int                        `json:"id"`
	MountPoint string                     `json:"mount_point"`
	Groups     map[string]SharePermission `json:"groups"`
	Quota      int                        `json:"quota"`
	Size       int                        `json:"size"`
}

func (gf *groupFolderBadFormatGroups) FormatGroupFolder() GroupFolder {
	g := GroupFolder{}
	g.ID = gf.ID
	g.MountPoint = gf.MountPoint
	g.Groups = map[string]SharePermission{}
	for k, v := range gf.Groups {
		g.Groups[k] = SharePermission(v)
	}
	g.Quota = gf.Quota
	g.Size = gf.Size
	return g
}

func (gf *groupFolderBadFormatIDAndGroups) FormatGroupFolder() GroupFolder {
	g := GroupFolder{}
	g.ID = gf.ID
	g.MountPoint = gf.MountPoint
	g.Groups = map[string]SharePermission{}
	for k, v := range gf.Groups {
		g.Groups[k] = SharePermission(v)
	}
	g.Quota = gf.Quota
	g.Size = gf.Size
	return g
}
