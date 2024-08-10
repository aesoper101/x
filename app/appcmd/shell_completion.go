package appcmd

import "github.com/spf13/cobra"

var (
	ShellCompDirectiveError         = newShellCompDirective(cobra.ShellCompDirectiveError)
	ShellCompDirectiveNoSpace       = newShellCompDirective(cobra.ShellCompDirectiveNoSpace)
	ShellCompDirectiveNoFileComp    = newShellCompDirective(cobra.ShellCompDirectiveNoFileComp)
	ShellCompDirectiveFilterFileExt = newShellCompDirective(cobra.ShellCompDirectiveFilterFileExt)
	ShellCompDirectiveFilterDirs    = newShellCompDirective(cobra.ShellCompDirectiveFilterDirs)
	ShellCompDirectiveKeepOrder     = newShellCompDirective(cobra.ShellCompDirectiveKeepOrder)
)

type ShellCompDirective struct {
	completion cobra.ShellCompDirective
}

func (s ShellCompDirective) cobra() cobra.ShellCompDirective {
	return s.completion
}

func newShellCompDirective(completion cobra.ShellCompDirective) ShellCompDirective {
	return ShellCompDirective{completion: completion}
}
