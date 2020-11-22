package ocl

// PermissionLevel is a permission enum for ocl.
type PermissionLevel int

const (
	// PermissionMember is the member level permission.
	PermissionMember PermissionLevel = iota

	// PermissionModerator is the moderator level permission.
	PermissionModerator

	// PermissionAdministrator is the administrator level permission.
	PermissionAdministrator

	// PermissionServerOwner is the server owner administration level permission.
	PermissionServerOwner

	// PermissionMaintainer is the permissions for the bot maintainers.
	PermissionMaintainer
)

var permNames = [...]string{"Member", "Moderator", "Administrator", "Owner", "Maintainer"}

func (perm PermissionLevel) String() string {
	return permNames[perm]
}
