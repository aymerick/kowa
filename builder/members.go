package builder

import (
	"time"

	"github.com/nicksnyder/go-i18n/i18n"

	"github.com/aymerick/kowa/models"
)

// Members node builder
type MembersBuilder struct {
	*NodeBuilderBase

	// loaded members
	membersVars []*MemberVars
}

// Members node content
type MembersContent struct {
	Node *Node

	Members []*MemberVars
}

// Member vars
type MemberVars struct {
	Date        time.Time
	Photo       *ImageVars
	Fullname    string
	Role        string
	Description string
}

func init() {
	RegisterNodeBuilder(KIND_MEMBERS, NewMembersBuilder)
}

func NewMembersBuilder(siteBuilder *SiteBuilder) NodeBuilder {
	return &MembersBuilder{
		NodeBuilderBase: &NodeBuilderBase{
			nodeKind:    KIND_MEMBERS,
			siteBuilder: siteBuilder,
		},
	}
}

// NodeBuilder
func (builder *MembersBuilder) Load() {
	// fetch members
	membersVars := builder.members()
	if len(membersVars) == 0 {
	}

	// get page settings
	title, tagline, cover, disabled := builder.pageSettings(models.PAGE_KIND_MEMBERS)
	if disabled {
		return
	}

	T := i18n.MustTfunc(builder.siteLang())
	slug := T("members")

	if title == "" {
		title = slug
	}

	// build node
	node := builder.newNode()
	node.fillUrl(slug)

	node.Title = title
	node.Tagline = tagline
	node.Cover = cover

	node.Meta = &NodeMeta{Description: tagline}

	node.InNavBar = true
	node.NavBarOrder = 15

	node.Content = &MembersContent{
		Node:    node,
		Members: membersVars,
	}

	builder.addNode(node)
}

// NodeBuilder
func (builder *MembersBuilder) Data(name string) interface{} {
	switch name {
	case "members":
		return builder.members()
	}

	return nil
}

// returns members contents
func (builder *MembersBuilder) members() []*MemberVars {
	if len(builder.membersVars) == 0 {
		// fetch members
		for _, member := range *builder.site().FindAllMembers() {
			memberVars := builder.NewMemberVars(member)

			builder.membersVars = append(builder.membersVars, memberVars)
		}
	}

	return builder.membersVars
}

func (builder *MembersBuilder) NewMemberVars(member *models.Member) *MemberVars {
	result := &MemberVars{
		Date:     member.CreatedAt,
		Fullname: member.Fullname,
		Role:     member.Role,
	}

	photo := member.FindPhoto()
	if photo != nil {
		result.Photo = builder.addImage(photo)
	}

	result.Description = member.Description

	return result
}
