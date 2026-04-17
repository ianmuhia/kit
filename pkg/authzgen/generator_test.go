package authzgen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// writeSchema writes a schema string to a temp file and returns the path.
func writeSchema(t *testing.T, content string) string {
	t.Helper()
	f := filepath.Join(t.TempDir(), "schema.zed")
	require.NoError(t, os.WriteFile(f, []byte(content), 0o644))
	return f
}

// findDef locates a definition by name inside a parsed Schema.
func findDef(t *testing.T, s *Schema, name string) Definition {
	t.Helper()
	for _, d := range s.Definitions {
		if d.Name == name {
			return d
		}
	}
	t.Fatalf("definition %q not found in schema", name)
	return Definition{}
}

// ToPascalCase ───

func TestToPascalCase(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"user", "User"},
		{"doctype", "Doctype"},
		{"direct_member", "DirectMember"},
		{"org_member", "OrgMember"},
		{"my-resource", "MyResource"},
		{"already", "Already"},
		{"UPPER", "Upper"},
		{"", ""},
		{"  spaces  ", "Spaces"},
	}
	for _, tc := range cases {
		assert.Equal(t, tc.want, ToPascalCase(tc.in), "input: %q", tc.in)
	}
}

// splitNamespace ─

func TestSplitNamespace(t *testing.T) {
	cases := []struct {
		ns, wantPkg, wantName string
	}{
		{"user", "authz", "user"},
		{"team", "authz", "team"},
		{"platform/user", "platform", "user"},
		{"org/team/resource", "org/team", "resource"}, // last slash wins
	}
	for _, tc := range cases {
		pkg, name := splitNamespace(tc.ns)
		assert.Equal(t, tc.wantPkg, pkg, "pkg for %q", tc.ns)
		assert.Equal(t, tc.wantName, name, "name for %q", tc.ns)
	}
}

// NewGenerator ───

func TestNewGenerator(t *testing.T) {
	t.Run("missing schema file returns error", func(t *testing.T) {
		_, err := NewGenerator(WithOutputDir(t.TempDir()))
		require.ErrorContains(t, err, "schema file is required")
	})

	t.Run("valid options succeeds", func(t *testing.T) {
		f := writeSchema(t, "definition user {}")
		g, err := NewGenerator(WithSchemaFile(f), WithOutputDir(t.TempDir()))
		require.NoError(t, err)
		assert.NotNil(t, g)
	})

	t.Run("default output dir is dot", func(t *testing.T) {
		f := writeSchema(t, "definition user {}")
		g, err := NewGenerator(WithSchemaFile(f))
		require.NoError(t, err)
		assert.Equal(t, ".", g.outputDir)
	})
}

// parseSchema ────

func TestParseSchema_MissingFile(t *testing.T) {
	g, err := NewGenerator(WithSchemaFile("/nonexistent/schema.zed"), WithOutputDir(t.TempDir()))
	require.NoError(t, err)
	_, err = g.parseSchema()
	require.ErrorContains(t, err, "failed to read schema file")
}

func TestParseSchema_InvalidSchema(t *testing.T) {
	g, err := NewGenerator(
		WithSchemaFile(writeSchema(t, "this is not valid schema content")),
		WithOutputDir(t.TempDir()),
	)
	require.NoError(t, err)
	_, err = g.parseSchema()
	require.ErrorContains(t, err, "failed to compile schema")
}

func TestParseSchema_SimpleDefinition(t *testing.T) {
	schema := `
definition user {}

definition team {
    relation direct_member: user
    permission member = direct_member
}`
	g, err := NewGenerator(WithSchemaFile(writeSchema(t, schema)), WithOutputDir(t.TempDir()))
	require.NoError(t, err)

	s, err := g.parseSchema()
	require.NoError(t, err)
	require.Len(t, s.Definitions, 2)

	user := findDef(t, s, "user")
	assert.Empty(t, user.Relations)
	assert.Empty(t, user.Permissions)
	assert.Equal(t, "authz", user.Package)

	team := findDef(t, s, "team")
	assert.Equal(t, "authz", team.Package)
	require.Len(t, team.Relations, 1)
	assert.Equal(t, "direct_member", team.Relations[0].Name)
	assert.Equal(t, []string{"user"}, team.Relations[0].Types)
	assert.False(t, team.Relations[0].IsUnion)
	require.Len(t, team.Permissions, 1)
	assert.Equal(t, "member", team.Permissions[0].Name)
}

func TestParseSchema_UnionRelation(t *testing.T) {
	schema := `
definition user {}

definition team {
    relation direct_member: user
    permission member = direct_member
}

definition doctype {
    relation admin: user | team#member
    permission view = admin
}`
	g, err := NewGenerator(WithSchemaFile(writeSchema(t, schema)), WithOutputDir(t.TempDir()))
	require.NoError(t, err)

	s, err := g.parseSchema()
	require.NoError(t, err)

	dt := findDef(t, s, "doctype")
	require.Len(t, dt.Relations, 1)

	admin := dt.Relations[0]
	assert.Equal(t, "admin", admin.Name)
	assert.True(t, admin.IsUnion)
	assert.Contains(t, admin.Types, "user")
	assert.Contains(t, admin.Types, "team#member")

	require.Len(t, dt.Permissions, 1)
	assert.Equal(t, "view", dt.Permissions[0].Name)
}

func TestParseSchema_SubjectRelation(t *testing.T) {
	// team#member is a subject relation — type string must include the fragment
	schema := `
definition user {}

definition team {
    relation direct_member: user
    permission member = direct_member
}

definition resource {
    relation viewer: user | team#member
}`
	g, err := NewGenerator(WithSchemaFile(writeSchema(t, schema)), WithOutputDir(t.TempDir()))
	require.NoError(t, err)

	s, err := g.parseSchema()
	require.NoError(t, err)

	res := findDef(t, s, "resource")
	require.Len(t, res.Relations, 1)
	assert.Contains(t, res.Relations[0].Types, "team#member")
}

func TestParseSchema_PrefixedNamespace(t *testing.T) {
	schema := `
definition platform/user {}

definition platform/document {
    relation owner: platform/user
    permission read = owner
}`
	g, err := NewGenerator(
		WithSchemaFile(writeSchema(t, schema)),
		WithOutputDir(t.TempDir()),
	)
	require.NoError(t, err)

	s, err := g.parseSchema()
	require.NoError(t, err)

	for _, d := range s.Definitions {
		assert.Equal(t, "platform", d.Package)
	}

	doc := findDef(t, s, "document")
	require.Len(t, doc.Relations, 1)
	assert.Equal(t, []string{"platform/user"}, doc.Relations[0].Types)
}

func TestParseSchema_ArrowPermission(t *testing.T) {
	// Arrow expressions (team->member) must compile without error and produce a permission.
	schema := `
definition user {}

definition team {
    relation direct_member: user
    permission member = direct_member
}

definition organization {
    relation staff: user | team#member
    relation team: team
    permission org_member = staff + team->member
}`
	g, err := NewGenerator(WithSchemaFile(writeSchema(t, schema)), WithOutputDir(t.TempDir()))
	require.NoError(t, err)

	s, err := g.parseSchema()
	require.NoError(t, err)

	org := findDef(t, s, "organization")
	assert.Len(t, org.Relations, 2)

	permNames := make([]string, len(org.Permissions))
	for i, p := range org.Permissions {
		permNames[i] = p.Name
	}
	assert.Contains(t, permNames, "org_member")
}

// Generate (end-to-end)────

const fullSchema = `
definition user {}

definition team {
    relation direct_member: user

    permission member = direct_member
}

definition organization {
    relation staff: user | team#member
    relation team: team

    permission org_member = staff + team->member
}

definition doctype {
    relation organization: organization

    relation admin: user | team#member
    relation maintainer: user | team#member
    relation triager: user | team#member
    relation writer: user | team#member
    relation reader: user | team#member

    permission create = writer + admin
    permission read = reader + triager + writer + maintainer + admin
    permission edit = maintainer + admin
    permission delete = edit
}`

func TestGenerate_CreatesOutputFile(t *testing.T) {
	outDir := t.TempDir()
	g, err := NewGenerator(
		WithSchemaFile(writeSchema(t, fullSchema)),
		WithOutputDir(outDir),
	)
	require.NoError(t, err)
	require.NoError(t, g.Generate())

	for _, name := range []string{
		"client.gen.go",
		"doctype.gen.go",
		"organization.gen.go",
		"team.gen.go",
		"user.gen.go",
	} {
		_, err = os.Stat(filepath.Join(outDir, name))
		require.NoError(t, err, "%s should exist", name)
	}
}

func TestGenerate_OutputContainsExpectedSymbols(t *testing.T) {
	outDir := t.TempDir()
	g, err := NewGenerator(
		WithSchemaFile(writeSchema(t, fullSchema)),
		WithOutputDir(outDir),
	)
	require.NoError(t, err)
	require.NoError(t, g.Generate())

	readFile := func(name string) string {
		raw, err := os.ReadFile(filepath.Join(outDir, name))
		require.NoError(t, err)
		return string(raw)
	}

	client := readFile("client.gen.go")
	doctype := readFile("doctype.gen.go")
	team := readFile("team.gen.go")
	org := readFile("organization.gen.go")
	user := readFile("user.gen.go")

	// client.gen.go─
	assert.True(t, strings.HasPrefix(strings.TrimSpace(client), "// Code generated"))
	assert.Contains(t, client, "package authz")
	assert.Contains(t, client, "type Type string")
	assert.Contains(t, client, "type Client struct")
	assert.Contains(t, client, "func NewClient(")
	assert.NotContains(t, client, "SetupClient")
	assert.NotContains(t, client, "GetClient")

	// per-definition package declarations
	for _, src := range []string{doctype, team, org, user} {
		assert.Contains(t, src, "package authz")
	}

	// type constants live in their own files─────
	assert.Contains(t, user, `TypeUser Type = "user"`)
	assert.Contains(t, team, `TypeTeam Type = "team"`)
	assert.Contains(t, org, `TypeOrganization Type = "organization"`)
	assert.Contains(t, doctype, `TypeDoctype Type = "doctype"`)

	// doctype relation & permission constants────
	assert.Contains(t, doctype, `DoctypeAdminRel`)
	assert.Contains(t, doctype, `DoctypeMaintainerRel`)
	assert.Contains(t, doctype, `DoctypeTriagerRel`)
	assert.Contains(t, doctype, `DoctypeWriterRel`)
	assert.Contains(t, doctype, `DoctypeReaderRel`)
	assert.Contains(t, doctype, `DoctypeOrganizationRel`)
	assert.Contains(t, doctype, `DoctypeCreatePerm`)
	assert.Contains(t, doctype, `DoctypeReadPerm`)
	assert.Contains(t, doctype, `DoctypeEditPerm`)
	assert.Contains(t, doctype, `DoctypeDeletePerm`)

	// union objects struct has both subject types
	assert.Contains(t, doctype, "DoctypeAdminObjects")
	assert.Contains(t, doctype, "User []User")
	assert.Contains(t, doctype, "Team []Team")

	// store struct, interface, constructor
	assert.Contains(t, doctype, "type DoctypeStore struct")
	assert.Contains(t, doctype, "type DoctypeStoreInterface interface")
	assert.Contains(t, doctype, "func NewDoctypeStore(")

	// CRUD methods are on the store, not on the resource type
	assert.Contains(t, doctype, "func (s *DoctypeStore) CreateAdminRelations(")
	assert.Contains(t, doctype, "func (s *DoctypeStore) DeleteAdminRelations(")
	assert.Contains(t, doctype, "func (s *DoctypeStore) ReadAdminRelations(")
	assert.NotContains(t, doctype, "func (r Doctype) CreateAdminRelations(")

	// permission methods are on the store─
	assert.Contains(t, doctype, "func (s *DoctypeStore) CheckCreate(")
	assert.Contains(t, doctype, "func (s *DoctypeStore) CheckRead(")
	assert.Contains(t, doctype, "func (s *DoctypeStore) CheckEdit(")
	assert.Contains(t, doctype, "func (s *DoctypeStore) CheckDelete(")
	assert.Contains(t, doctype, "func (s *DoctypeStore) LookupReadResources(")
	assert.Contains(t, doctype, "func (s *DoctypeStore) LookupReadSubjects(")

	// Subject type replaces *v1.SubjectReference in the public API
	assert.Contains(t, client, "type Subject struct")
	assert.Contains(t, client, "func NewSubject(")
	assert.Contains(t, client, "func NewSubjectWithRelation(")
	assert.Contains(t, doctype, "subject Subject)")
	assert.NotContains(t, doctype, "*v1.SubjectReference")

	// old patterns must not exist──
	assert.NotContains(t, doctype, "CheckDoctypeCreateInputs")
	assert.NotContains(t, doctype, "func LookupDoctypeReadResources(")
	assert.NotContains(t, doctype, "func LookupDoctypeReadResourcesForUser(")
	assert.NotContains(t, doctype, "GetClient")
}

func TestGenerate_PrefixedNamespaceOutputFile(t *testing.T) {
	schema := `
definition platform/user {}

definition platform/document {
    relation owner: platform/user
    permission read = owner
}`
	outDir := t.TempDir()
	g, err := NewGenerator(
		WithSchemaFile(writeSchema(t, schema)),
		WithOutputDir(outDir),
	)
	require.NoError(t, err)
	require.NoError(t, g.Generate())

	// Package name comes from first definition prefix → "platform".
	// Expect: client.gen.go + one file per definition.
	for _, name := range []string{"client.gen.go", "user.gen.go", "document.gen.go"} {
		_, err = os.Stat(filepath.Join(outDir, name))
		require.NoError(t, err, "%s should exist", name)
		raw, _ := os.ReadFile(filepath.Join(outDir, name))
		assert.Contains(t, string(raw), "package platform", "%s should declare package platform", name)
	}
}

func TestGenerate_InvalidSchemaReturnsError(t *testing.T) {
	g, err := NewGenerator(
		WithSchemaFile(writeSchema(t, "not a valid schema")),
		WithOutputDir(t.TempDir()),
	)
	require.NoError(t, err)
	require.ErrorContains(t, g.Generate(), "failed to parse schema")
}
