package generator

import "strings"

type CodeTemplate string

func (g CodeTemplate) String() string {
	return string(g)
}

func (g CodeTemplate) Replace(s, r string, n int) CodeTemplate {
	return CodeTemplate(strings.Replace(string(g), s, r, n))
}

func resolveRepTag(g CodeTemplate) CodeTemplate {
	return g.Replace("<Ptr>", "{{.Ptr}}", -1).
		Replace("<E>", "{{.E}}", -1).
		Replace("<E2>", "{{.E2}}", -1).
		Replace("<R>", "{{.R}}", -1).
		Replace("<R2>", "{{.R2}}", -1).
		Replace("<PK>", "{{.PK}}", -1)
}

func init() {
	TPL_ENTITY_REP = resolveRepTag(TPL_ENTITY_REP)
	TPL_ENTITY_REP_INTERFACE = resolveRepTag(TPL_ENTITY_REP_INTERFACE)
	TPL_REPO_FACTORY = resolveRepTag(TPL_REPO_FACTORY)
}
