package parser

type SelectorStrategy struct {
	TitleSelectors   []string
	ContentSelectors []string
	AuthorSelectors  []string
	DateSelectors    []string
}

func DefaultStrategy() *SelectorStrategy {
	return &SelectorStrategy{
		TitleSelectors: []string{

			"h1.post-title",
			"h1.title",
			"div.post-header h1",
			".entry-title",
			"article h1",
			"h1",
			"h2",
			"title",
		},

		ContentSelectors: []string{
			"div.post-content",
			"div.entry-content",
			"article.post",
			".post-body",
			".content",
			"div.message-body",
			"article",
			"main",
			"body",
		},

		AuthorSelectors: []string{
			".author",
			".post-author",
			".username",
			"span.user",
			"a.author-link",
			"div.author-info",
		},

		DateSelectors: []string{
			"time",
			".post-date",
			".published",
			".date",
			"span.timestamp",
			".entry-date",
		},
	}
}
