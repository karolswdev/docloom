// Package parser provides a Go-native C# code parser using tree-sitter.
package parser

import (
	"context"
	"fmt"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/csharp"
)

// APISurface represents the extracted public API surface of C# code
type APISurface struct {
	Namespaces []Namespace `json:"namespaces"`
}

// Namespace represents a C# namespace
type Namespace struct {
	Name    string  `json:"name"`
	Classes []Class `json:"classes"`
}

// Class represents a C# class with its members
type Class struct {
	Name        string   `json:"name"`
	DocComment  string   `json:"docComment,omitempty"`
	Methods     []Method `json:"methods"`
	Properties  []Property `json:"properties"`
	IsPublic    bool     `json:"isPublic"`
	IsAbstract  bool     `json:"isAbstract"`
	IsInterface bool     `json:"isInterface"`
}

// Method represents a C# method
type Method struct {
	Name       string   `json:"name"`
	Signature  string   `json:"signature"`
	DocComment string   `json:"docComment,omitempty"`
	IsPublic   bool     `json:"isPublic"`
	IsStatic   bool     `json:"isStatic"`
	Parameters []Parameter `json:"parameters"`
	ReturnType string   `json:"returnType"`
}

// Property represents a C# property
type Property struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	DocComment string `json:"docComment,omitempty"`
	IsPublic   bool   `json:"isPublic"`
	IsStatic   bool   `json:"isStatic"`
}

// Parameter represents a method parameter
type Parameter struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// Parser provides C# code parsing capabilities
type Parser struct {
	parser *sitter.Parser
}

// New creates a new C# parser
func New() *Parser {
	parser := sitter.NewParser()
	parser.SetLanguage(csharp.GetLanguage())
	return &Parser{
		parser: parser,
	}
}

// ExtractAPISurface parses C# source code and extracts the public API surface
func (p *Parser) ExtractAPISurface(ctx context.Context, source string) (*APISurface, error) {
	tree, err := p.parser.ParseCtx(ctx, nil, []byte(source))
	if err != nil {
		return nil, fmt.Errorf("failed to parse C# source: %w", err)
	}
	defer tree.Close()

	root := tree.RootNode()
	
	api := &APISurface{
		Namespaces: make([]Namespace, 0),
	}

	// Extract namespaces and their contents
	p.extractNamespaces(root, source, api)

	// Also extract top-level classes (not in namespaces)
	p.extractTopLevelTypes(root, source, api)

	return api, nil
}

// extractNamespaces extracts all namespaces from the syntax tree
func (p *Parser) extractNamespaces(node *sitter.Node, source string, api *APISurface) {
	if node.Type() == "namespace_declaration" || node.Type() == "file_scoped_namespace_declaration" {
		ns := p.parseNamespace(node, source)
		if ns != nil {
			api.Namespaces = append(api.Namespaces, *ns)
		}
	}

	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		p.extractNamespaces(child, source, api)
	}
}

// extractTopLevelTypes extracts classes not within namespaces
func (p *Parser) extractTopLevelTypes(node *sitter.Node, source string, api *APISurface) {
	// Create a default namespace for top-level types
	defaultNs := Namespace{
		Name:    "<global>",
		Classes: make([]Class, 0),
	}

	p.extractClassesFromNode(node, source, &defaultNs, 0)

	if len(defaultNs.Classes) > 0 {
		api.Namespaces = append(api.Namespaces, defaultNs)
	}
}

// parseNamespace parses a namespace declaration
func (p *Parser) parseNamespace(node *sitter.Node, source string) *Namespace {
	ns := &Namespace{
		Classes: make([]Class, 0),
	}

	// Extract namespace name
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child.Type() == "qualified_name" || child.Type() == "identifier" {
			ns.Name = source[child.StartByte():child.EndByte()]
			break
		}
	}

	// Extract classes within namespace
	p.extractClassesFromNode(node, source, ns, 0)

	return ns
}

// extractClassesFromNode recursively extracts classes from a node
func (p *Parser) extractClassesFromNode(node *sitter.Node, source string, ns *Namespace, depth int) {
	// Prevent infinite recursion for nested namespaces
	if depth > 10 {
		return
	}

	nodeType := node.Type()
	
	if nodeType == "class_declaration" || nodeType == "interface_declaration" || nodeType == "struct_declaration" {
		class := p.parseClass(node, source)
		if class != nil {
			ns.Classes = append(ns.Classes, *class)
		}
		return // Don't recurse into class bodies
	}

	// Skip namespace declarations at this level to avoid duplication
	if nodeType == "namespace_declaration" || nodeType == "file_scoped_namespace_declaration" {
		if depth > 0 {
			return
		}
	}

	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		p.extractClassesFromNode(child, source, ns, depth+1)
	}
}

// parseClass parses a class, interface, or struct declaration
func (p *Parser) parseClass(node *sitter.Node, source string) *Class {
	class := &Class{
		Methods:     make([]Method, 0),
		Properties:  make([]Property, 0),
		IsInterface: node.Type() == "interface_declaration",
	}

	// Extract modifiers
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		childType := child.Type()
		
		if childType == "modifier" {
			modifier := source[child.StartByte():child.EndByte()]
			switch modifier {
			case "public":
				class.IsPublic = true
			case "abstract":
				class.IsAbstract = true
			}
		} else if childType == "identifier" {
			class.Name = source[child.StartByte():child.EndByte()]
		} else if childType == "declaration_list" {
			// Parse class members
			p.parseClassMembers(child, source, class)
		}
	}

	// Extract XML doc comment if present
	class.DocComment = p.extractDocComment(node, source)

	return class
}

// parseClassMembers parses the members of a class
func (p *Parser) parseClassMembers(node *sitter.Node, source string, class *Class) {
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		childType := child.Type()

		switch childType {
		case "method_declaration":
			if method := p.parseMethod(child, source); method != nil {
				class.Methods = append(class.Methods, *method)
			}
		case "property_declaration":
			if prop := p.parseProperty(child, source); prop != nil {
				class.Properties = append(class.Properties, *prop)
			}
		case "field_declaration":
			// Could parse fields as properties if needed
		}
	}
}

// parseMethod parses a method declaration
func (p *Parser) parseMethod(node *sitter.Node, source string) *Method {
	method := &Method{
		Parameters: make([]Parameter, 0),
	}

	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		childType := child.Type()

		switch childType {
		case "modifier":
			modifier := source[child.StartByte():child.EndByte()]
			switch modifier {
			case "public":
				method.IsPublic = true
			case "static":
				method.IsStatic = true
			}
		case "identifier":
			method.Name = source[child.StartByte():child.EndByte()]
		case "predefined_type", "nullable_type", "array_type", "generic_name":
			method.ReturnType = source[child.StartByte():child.EndByte()]
		case "parameter_list":
			p.parseParameters(child, source, method)
		}
	}

	// Build signature
	params := make([]string, 0, len(method.Parameters))
	for _, p := range method.Parameters {
		params = append(params, fmt.Sprintf("%s %s", p.Type, p.Name))
	}
	method.Signature = fmt.Sprintf("%s %s(%s)", method.ReturnType, method.Name, strings.Join(params, ", "))

	// Extract XML doc comment
	method.DocComment = p.extractDocComment(node, source)

	return method
}

// parseProperty parses a property declaration
func (p *Parser) parseProperty(node *sitter.Node, source string) *Property {
	prop := &Property{}

	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		childType := child.Type()

		switch childType {
		case "modifier":
			modifier := source[child.StartByte():child.EndByte()]
			switch modifier {
			case "public":
				prop.IsPublic = true
			case "static":
				prop.IsStatic = true
			}
		case "identifier":
			prop.Name = source[child.StartByte():child.EndByte()]
		case "predefined_type", "nullable_type", "array_type", "generic_name":
			prop.Type = source[child.StartByte():child.EndByte()]
		}
	}

	// Extract XML doc comment
	prop.DocComment = p.extractDocComment(node, source)

	return prop
}

// parseParameters parses method parameters
func (p *Parser) parseParameters(node *sitter.Node, source string, method *Method) {
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child.Type() == "parameter" {
			param := Parameter{}
			
			for j := 0; j < int(child.ChildCount()); j++ {
				grandchild := child.Child(j)
				grandchildType := grandchild.Type()
				
				switch grandchildType {
				case "identifier":
					param.Name = source[grandchild.StartByte():grandchild.EndByte()]
				case "predefined_type", "nullable_type", "array_type", "generic_name":
					param.Type = source[grandchild.StartByte():grandchild.EndByte()]
				}
			}
			
			if param.Name != "" {
				method.Parameters = append(method.Parameters, param)
			}
		}
	}
}

// extractDocComment extracts XML documentation comments preceding a node
func (p *Parser) extractDocComment(node *sitter.Node, source string) string {
	// Look for comment nodes before the current node
	if node.Parent() != nil {
		parent := node.Parent()
		nodeIndex := -1
		
		// Find the index of the current node
		for i := 0; i < int(parent.ChildCount()); i++ {
			if parent.Child(i) == node {
				nodeIndex = i
				break
			}
		}
		
		// Look for comments before this node
		if nodeIndex > 0 {
			prevNode := parent.Child(nodeIndex - 1)
			if prevNode.Type() == "comment" {
				comment := source[prevNode.StartByte():prevNode.EndByte()]
				// Clean up XML doc comment markers
				comment = strings.TrimPrefix(comment, "///")
				comment = strings.TrimSpace(comment)
				return comment
			}
		}
	}
	
	return ""
}