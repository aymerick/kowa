package builder

import (
	"time"

	"github.com/nicksnyder/go-i18n/i18n"

	"github.com/aymerick/kowa/models"
)

// MembersBuilder builds members pages
type MembersBuilder struct {
	*NodeBuilderBase

	// loaded members
	membersVars []*MemberVars
}

// MembersContent represents members node content
type MembersContent struct {
	Members []*MemberVars
}

// MemberVars reprsents member vars
type MemberVars struct {
	Date        time.Time
	Photo       *ImageVars
	Fullname    string
	Role        string
	Description string
}

func init() {
	RegisterNodeBuilder(kindMembers, NewMembersBuilder)
}

// NewMembersBuilder instanciates a new NodeBuilder
func NewMembersBuilder(siteBuilder *SiteBuilder) NodeBuilder {
	return &MembersBuilder{
		NodeBuilderBase: &NodeBuilderBase{
			nodeKind:    kindMembers,
			siteBuilder: siteBuilder,
		},
	}
}

// Load is part of NodeBuilder interface
func (builder *MembersBuilder) Load() {
	// fetch members
	membersVars := builder.members()
	if len(membersVars) == 0 {
		return
	}

	// get page settings
	title, tagline, cover, disabled := builder.pageSettings(models.PageKindMembers)
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
	node.fillURL(slug)

	node.Title = title
	node.Tagline = tagline
	node.Cover = cover

	node.Meta = &NodeMeta{Description: tagline}

	node.InNavBar = true
	node.NavBarOrder = 15

	node.Content = &MembersContent{
		Members: membersVars,
	}

	builder.addNode(node)
}

// Data is part of NodeBuilder interface
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

// NewMemberVars instanciates a new MemberVars
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
