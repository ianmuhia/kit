package authzgen

import (
	"fmt"
	"log/slog"
)

// Node is the base interface for AST nodes
type Node interface {
	String() string
}

// ObjectType represents an object type in the schema
type ObjectType struct {
	Name   string
	Prefix string
}

// DefinitionNode represents a definition in the AST
type DefinitionNode struct {
	ObjectType  ObjectType
	Relations   []*RelationNode
	Permissions []*PermissionNode
}

// RelationNode represents a relation in the AST
type RelationNode struct {
	Name       string
	Expression RelationExpressionNode
}

// PermissionNode represents a permission in the AST
type PermissionNode struct {
	Name       string
	Expression PermissionExpressionNode
}

// PermissionExpressionNode represents a permission expression
type PermissionExpressionNode interface {
	Node
	permissionExpression()
}

// RelationExpressionNode represents a relation expression
type RelationExpressionNode interface {
	Node
	relationExpression()
}

// IdentifierNode represents an identifier
type IdentifierNode struct {
	Value string
}

func (i *IdentifierNode) String() string        { return i.Value }
func (i *IdentifierNode) permissionExpression() {}

// BinaryOpNode represents a binary operation
type BinaryOpNode struct {
	Operator string
	Left     PermissionExpressionNode
	Right    PermissionExpressionNode
}

func (b *BinaryOpNode) String() string {
	return fmt.Sprintf("%s %s %s", b.Left.String(), b.Operator, b.Right.String())
}
func (b *BinaryOpNode) permissionExpression() {}

// SingleRelationNode represents a single relation
type SingleRelationNode struct {
	Value    string
	Fragment string
}

func (s *SingleRelationNode) String() string {
	if s.Fragment != "" {
		return fmt.Sprintf("%s#%s", s.Value, s.Fragment)
	}
	return s.Value
}
func (s *SingleRelationNode) relationExpression() {}

// UnionRelationNode represents a union of relations
type UnionRelationNode struct {
	Left  RelationExpressionNode
	Right RelationExpressionNode
}

func (u *UnionRelationNode) String() string {
	return fmt.Sprintf("%s | %s", u.Left.String(), u.Right.String())
}
func (u *UnionRelationNode) relationExpression() {}

// Parser parses tokens into an AST
type Parser struct {
	tokens      []Token
	current     int
	Definitions []*DefinitionNode
}

// NewParser creates a new parser for the given tokens
func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens, Definitions: []*DefinitionNode{}}
}

// ParseDefinitions parses all definitions from the tokens
func (p *Parser) ParseDefinitions() ([]*DefinitionNode, error) {
	for p.peek().Type != EOF {
		if p.peek().Type != DEFINITION {
			slog.Error("Expected definition keyword",
				"line", p.peek().Line,
				"position", p.peek().Position,
				"got_literal", p.peek().Literal,
				"got_token", p.peek().Type.String(),
				"expected_token", DEFINITION.String())
			return nil, fmt.Errorf("expected %s, got %s (%s) at line %d",
				DEFINITION.String(), p.peek().Type.String(), p.peek().Literal, p.peek().Line)
		}

		def, err := p.parseDefinition()
		if err != nil {
			slog.Error("Failed to parse definition",
				"line", p.peek().Line,
				"error", err)
			return nil, err
		}
		p.Definitions = append(p.Definitions, def)
	}
	return p.Definitions, nil
}

func (p *Parser) peek() Token {
	if p.current >= len(p.tokens) {
		return Token{Type: EOF}
	}
	return p.tokens[p.current]
}

func (p *Parser) advance() {
	if p.current < len(p.tokens) {
		p.current++
	}
}

func (p *Parser) consume(expected TokenType) (Token, error) {
	token := p.peek()
	if token.Type == expected {
		p.advance()
		return token, nil
	}

	var errorContext string
	switch {
	case token.Type == EOF && expected != EOF:
		errorContext = "unexpected end of file"
	case token.Type == ILLEGAL:
		errorContext = "illegal character encountered"
	case expected == LBRACE && token.Type == IDENTIFIER:
		errorContext = "missing opening brace after object type definition"
	case expected == SLASH && token.Type == IDENTIFIER:
		errorContext = "missing slash in object type definition (expected format: prefix/name)"
	case expected == RBRACE:
		errorContext = "missing closing brace to end definition block"
	default:
		errorContext = "token mismatch"
	}

	slog.Error("Token consumption failed",
		"expected_token", expected.String(),
		"got_token", token.Type.String(),
		"got_literal", token.Literal,
		"line", token.Line,
		"position", token.Position,
		"context", errorContext)
	return token, fmt.Errorf("%s: expected %s, got %s (%s) at line %d",
		errorContext, expected.String(), token.Type.String(), token.Literal, token.Line)
}

func (p *Parser) parseDefinition() (*DefinitionNode, error) {
	defToken, err := p.consume(DEFINITION)
	if err != nil {
		return nil, fmt.Errorf("expected 'definition' keyword: %w", err)
	}

	firstToken, err := p.consume(IDENTIFIER)
	if err != nil {
		slog.Error("Failed to parse object type identifier",
			"line", p.peek().Line,
			"error", err)
		return nil, fmt.Errorf("expected object type identifier after 'definition': %w", err)
	}

	var prefix, name string
	var objectTypeName string

	if p.peek().Type == SLASH {
		_, err = p.consume(SLASH)
		if err != nil {
			slog.Error("Failed to parse slash in object type",
				"prefix", firstToken.Literal,
				"error", err)
			return nil, fmt.Errorf("expected '/' after object type prefix '%s': %w", firstToken.Literal, err)
		}

		nameToken, err := p.consume(IDENTIFIER)
		if err != nil {
			slog.Error("Failed to parse object type name",
				"prefix", firstToken.Literal,
				"error", err)
			return nil, fmt.Errorf("expected object type name after '%s/': %w", firstToken.Literal, err)
		}

		prefix = firstToken.Literal
		name = nameToken.Literal
		objectTypeName = prefix + "/" + name
	} else if p.peek().Type == LBRACE {
		prefix = ""
		name = firstToken.Literal
		objectTypeName = name
		slog.Debug("Using standard AuthZed definition format", "name", name)
	} else {
		slog.Error("Expected either '/' or '{' after object type identifier",
			"identifier", firstToken.Literal,
			"got_token", p.peek().Type.String(),
			"got_literal", p.peek().Literal,
			"line", p.peek().Line)
		return nil, fmt.Errorf("expected either '/' (for prefix/name format) or '{' (for standard format) after identifier '%s', got %s at line %d",
			firstToken.Literal, p.peek().Type.String(), p.peek().Line)
	}

	_, err = p.consume(LBRACE)
	if err != nil {
		slog.Error("Failed to parse opening brace for definition",
			"object_type", objectTypeName,
			"error", err)
		return nil, fmt.Errorf("expected '{' after object type '%s': %w", objectTypeName, err)
	}

	def := &DefinitionNode{
		ObjectType: ObjectType{
			Name:   name,
			Prefix: prefix,
		},
		Relations:   []*RelationNode{},
		Permissions: []*PermissionNode{},
	}

	slog.Debug("Parsing definition body",
		"object_type", objectTypeName,
		"format", func() string {
			if prefix != "" {
				return "prefix/name"
			}
			return "standard"
		}(),
		"line", defToken.Line)

	for p.peek().Type == RELATION || p.peek().Type == PERMISSION {
		if p.peek().Type == RELATION {
			rel, err := p.parseRelation()
			if err != nil {
				slog.Error("Failed to parse relation in definition",
					"object_type", objectTypeName,
					"error", err)
				return nil, err
			}
			def.Relations = append(def.Relations, rel)
		} else if p.peek().Type == PERMISSION {
			perm, err := p.parsePermission()
			if err != nil {
				slog.Error("Failed to parse permission in definition",
					"object_type", objectTypeName,
					"error", err)
				return nil, err
			}
			def.Permissions = append(def.Permissions, perm)
		}
	}

	_, err = p.consume(RBRACE)
	if err != nil {
		slog.Error("Failed to parse closing brace for definition",
			"object_type", objectTypeName,
			"error", err)
		return nil, fmt.Errorf("expected '}' to close definition '%s': %w", objectTypeName, err)
	}

	slog.Debug("Successfully parsed definition",
		"object_type", objectTypeName,
		"relations_count", len(def.Relations),
		"permissions_count", len(def.Permissions))

	return def, nil
}

func (p *Parser) parseRelation() (*RelationNode, error) {
	relationToken, err := p.consume(RELATION)
	if err != nil {
		slog.Error("Failed to consume RELATION token",
			"error", err,
			"line", relationToken.Line)
		return nil, fmt.Errorf("failed to parse relation declaration: %w", err)
	}

	nameToken, err := p.consume(IDENTIFIER)
	if err != nil {
		slog.Error("Failed to get relation name after 'relation' keyword",
			"error", err,
			"line", p.peek().Line)
		return nil, fmt.Errorf("expected relation name after 'relation' keyword: %w", err)
	}

	_, err = p.consume(COLON)
	if err != nil {
		slog.Error("Failed to consume colon after relation name",
			"relation_name", nameToken.Literal,
			"error", err)
		return nil, fmt.Errorf("expected ':' after relation name '%s': %w", nameToken.Literal, err)
	}

	expr, err := p.parseRelationExpression()
	if err != nil {
		slog.Error("Failed to parse relation expression",
			"relation_name", nameToken.Literal,
			"error", err)
		return nil, fmt.Errorf("failed to parse expression for relation '%s': %w", nameToken.Literal, err)
	}

	return &RelationNode{
		Name:       nameToken.Literal,
		Expression: expr,
	}, nil
}

func (p *Parser) parseRelationExpression() (RelationExpressionNode, error) {
	slog.Debug("Starting relation expression parsing", "current_token", p.peek().Type.String())

	left, err := p.parseSingleRelation()
	if err != nil {
		slog.Error("Failed to parse left side of relation expression", "error", err)
		return nil, err
	}

	for p.peek().Type == PIPE {
		slog.Debug("Found pipe in relation expression, parsing union")
		_, err := p.consume(PIPE)
		if err != nil {
			return nil, err
		}

		right, err := p.parseSingleRelation()
		if err != nil {
			slog.Error("Failed to parse right side of relation union", "error", err)
			return nil, err
		}

		left = &UnionRelationNode{Left: left, Right: right}
		slog.Debug("Created union relation node")
	}

	slog.Debug("Completed relation expression parsing")
	return left, nil
}

func (p *Parser) parseSingleRelation() (RelationExpressionNode, error) {
	slog.Debug("Parsing single relation", "current_token", p.peek().Type.String(), "literal", p.peek().Literal)

	identToken, err := p.consume(IDENTIFIER)
	if err != nil {
		slog.Error("Failed to consume identifier in single relation", "error", err)
		return nil, err
	}

	relationType := identToken.Literal
	slog.Debug("Got relation type identifier", "type", relationType)

	if p.peek().Type == SLASH {
		slog.Debug("Found slash, parsing prefixed relation type")
		_, err := p.consume(SLASH)
		if err != nil {
			return nil, err
		}

		nameToken, err := p.consume(IDENTIFIER)
		if err != nil {
			slog.Error("Failed to consume identifier after slash", "prefix", relationType, "error", err)
			return nil, err
		}

		relationType += "/" + nameToken.Literal
		slog.Debug("Updated relation type with prefix", "full_type", relationType)
	}

	node := &SingleRelationNode{Value: relationType}

	if p.peek().Type == HASH {
		slog.Debug("Found hash, parsing subject relation", "base_type", relationType)
		_, err := p.consume(HASH)
		if err != nil {
			slog.Error("Failed to consume hash token", "error", err)
			return nil, err
		}

		fragmentToken, err := p.consume(IDENTIFIER)
		if err != nil {
			slog.Error("Failed to consume fragment identifier after hash", "base_type", relationType, "error", err)
			return nil, err
		}

		node.Fragment = fragmentToken.Literal
		slog.Debug("Set subject relation fragment", "fragment", fragmentToken.Literal, "full_relation", relationType+"#"+fragmentToken.Literal)
	}

	slog.Debug("Completed single relation parsing", "value", node.Value, "fragment", node.Fragment)
	return node, nil
}

func (p *Parser) parsePermission() (*PermissionNode, error) {
	permissionToken, err := p.consume(PERMISSION)
	if err != nil {
		slog.Error("Failed to consume PERMISSION token",
			"error", err,
			"line", permissionToken.Line)
		return nil, fmt.Errorf("failed to parse permission declaration: %w", err)
	}

	nameToken, err := p.consume(IDENTIFIER)
	if err != nil {
		slog.Error("Failed to get permission name after 'permission' keyword",
			"error", err,
			"line", p.peek().Line)
		return nil, fmt.Errorf("expected permission name after 'permission' keyword: %w", err)
	}

	_, err = p.consume(EQUAL)
	if err != nil {
		slog.Error("Failed to consume equals sign after permission name",
			"permission_name", nameToken.Literal,
			"error", err)
		return nil, fmt.Errorf("expected '=' after permission name '%s': %w", nameToken.Literal, err)
	}

	expr, err := p.parsePermissionExpression()
	if err != nil {
		slog.Error("Failed to parse permission expression",
			"permission_name", nameToken.Literal,
			"error", err)
		return nil, fmt.Errorf("failed to parse expression for permission '%s': %w", nameToken.Literal, err)
	}

	return &PermissionNode{
		Name:       nameToken.Literal,
		Expression: expr,
	}, nil
}

func (p *Parser) parsePermissionExpression() (PermissionExpressionNode, error) {
	return p.parseAdditiveExpression()
}

func (p *Parser) parseAdditiveExpression() (PermissionExpressionNode, error) {
	left, err := p.parsePrimaryExpression()
	if err != nil {
		return nil, err
	}

	for p.peek().Type == PLUS {
		opToken, err := p.consume(PLUS)
		if err != nil {
			return nil, err
		}

		right, err := p.parsePrimaryExpression()
		if err != nil {
			return nil, err
		}

		left = &BinaryOpNode{Operator: opToken.Literal, Left: left, Right: right}
	}

	return left, nil
}

func (p *Parser) parsePrimaryExpression() (PermissionExpressionNode, error) {
	return p.parseIdentifierChain()
}

func (p *Parser) parseIdentifierChain() (PermissionExpressionNode, error) {
	identToken, err := p.consume(IDENTIFIER)
	if err != nil {
		return nil, err
	}

	var left PermissionExpressionNode = &IdentifierNode{Value: identToken.Literal}

	for p.peek().Type == MINUS_ARROW {
		opToken, err := p.consume(MINUS_ARROW)
		if err != nil {
			return nil, err
		}

		rightToken, err := p.consume(IDENTIFIER)
		if err != nil {
			return nil, err
		}

		right := &IdentifierNode{Value: rightToken.Literal}
		left = &BinaryOpNode{Operator: opToken.Literal, Left: left, Right: right}
	}

	return left, nil
}
