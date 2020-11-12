package utils

import "github.com/bwmarrin/discordgo"

// HasPermission checks if the bot has a permission.
func HasPermission(state *discordgo.State, chanID string, permMask int) (bool, error) {
	perm, err := state.UserChannelPermissions(state.User.ID, chanID)
	if err != nil {
		return false, err
	}

	return perm & permMask != 0, nil
}