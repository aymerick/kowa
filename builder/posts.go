package builder

type PostsBuilder struct {
	*NodeBuilder
}

func NewPostsBuilder(site *Site) *PostsBuilder {
	return &PostsBuilder{
		&NodeBuilder{
			Site:     site,
			NodeKind: KIND_POST,
		},
	}
}

func (builder *PostsBuilder) Load() {
	builder.BuildPostsLists()
	builder.BuildPosts()
}

func (builder *PostsBuilder) BuildPostsLists() {
	node := builder.NewNodeForKind(KIND_POSTS)

	node.Title = "Posts"

	node.Meta = &NodeMeta{
		Description: "Posts list",
	}

	node.Content = "Soon"

	builder.AddNode(node)
}

func (builder *PostsBuilder) BuildPosts() {
	node := builder.NewNode()

	node.Path = "post-1.html"

	node.Title = "Post #1"

	node.Meta = &NodeMeta{
		Description: "Post test #1",
	}

	node.Content = "Soon"

	builder.AddNode(node)
}
