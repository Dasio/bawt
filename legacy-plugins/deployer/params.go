package deployer

import (
	"fmt"
	"strings"

	"github.com/gopherworks/bawt"
)

type DeployParams struct {
	Environment      string
	Branch           string
	Tags             string
	DeploymentBranch string
	InitiatedBy      string
	From             string
	initiatedByChat  *bawt.Message
}

// ParsedTags returns *default* or user-specified tags
func (p *DeployParams) ParsedTags() string {
	tags := strings.Replace(p.Tags, " ", "", -1)
	if tags == "" {
		tags = "updt_streambed"
	}

	return tags
}

// ParsedDeploymentBranch returns the default, or a user-specified branch name
// used in the `deployment/` repo.
func (p *DeployParams) ParsedDeploymentBranch(branch string) string {
	if p.DeploymentBranch == "" {
		return branch
	}

	return p.DeploymentBranch
}

func (p *DeployParams) String() string {
	branch := p.Branch
	if branch == "" {
		branch = "[default]"
	}

	str := fmt.Sprintf("env=%s branch=%s tags=%s", p.Environment, branch, p.ParsedTags())

	if p.DeploymentBranch != "" {
		str = fmt.Sprintf("%s deploy_branch=%s", str, p.DeploymentBranch)
	}

	str = fmt.Sprintf("%s by %s", str, p.InitiatedBy)

	return str
}
