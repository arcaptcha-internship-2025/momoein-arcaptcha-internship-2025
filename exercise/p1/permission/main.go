package main

const (
	SeeMessages = 1 << iota
	DeleteMessages
	EditMessages
	KickMembers
	MakeMembersAdmin
	AddMembers
)

type Permissions struct {
	canSeeMessages      bool
	canDeleteMessages   bool
	canEditMessages     bool
	canKickMembers      bool
	canMakeMembersAdmin bool
	canAddMembers       bool
}

func SetUserPermissions(permissions *Permissions) (perm int8) {
	if permissions.canSeeMessages {
		perm |= SeeMessages
	}
	if permissions.canDeleteMessages {
		perm |= DeleteMessages
	}
	if permissions.canEditMessages {
		perm |= EditMessages
	}
	if permissions.canKickMembers {
		perm |= KickMembers
	}
	if permissions.canMakeMembersAdmin {
		perm |= MakeMembersAdmin
	}
	if permissions.canAddMembers {
		perm |= AddMembers
	}
	return
}

func GetUserPermissions(permissionsMask int8) *Permissions {
	return &Permissions{
		canSeeMessages:      permissionsMask&SeeMessages != 0,
		canDeleteMessages:   permissionsMask&DeleteMessages != 0,
		canEditMessages:     permissionsMask&EditMessages != 0,
		canKickMembers:      permissionsMask&KickMembers != 0,
		canMakeMembersAdmin: permissionsMask&MakeMembersAdmin != 0,
		canAddMembers:       permissionsMask&AddMembers != 0,
	}
}
