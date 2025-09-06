package parser

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCSharpParser_ExtractAPISurface(t *testing.T) {
	// Arrange
	sampleCode := `
using System;

namespace SampleNamespace
{
	/// <summary>
	/// Sample class for testing
	/// </summary>
	public class SampleClass
	{
		/// <summary>
		/// Gets or sets the name
		/// </summary>
		public string Name { get; set; }
		
		/// <summary>
		/// Gets the ID
		/// </summary>
		public int Id { get; private set; }
		
		/// <summary>
		/// Calculates the sum of two numbers
		/// </summary>
		public int Add(int a, int b)
		{
			return a + b;
		}
		
		/// <summary>
		/// Static utility method
		/// </summary>
		public static string FormatString(string input)
		{
			return $"Formatted: {input}";
		}
	}
	
	public interface ISampleInterface
	{
		void DoSomething();
		string GetValue();
	}
	
	internal class InternalClass
	{
		public void InternalMethod() { }
	}
}

// Top-level class outside namespace
public class GlobalClass
{
	public void GlobalMethod(string param1, int param2) { }
}`

	// Act
	parser := New()
	api, err := parser.ExtractAPISurface(context.Background(), sampleCode)

	// Assert
	require.NoError(t, err, "Parser should not return an error")
	require.NotNil(t, api, "API surface should not be nil")

	// Verify we have namespaces
	assert.GreaterOrEqual(t, len(api.Namespaces), 1, "Should have at least one namespace")

	// Find the SampleNamespace
	var sampleNs *Namespace
	for i := range api.Namespaces {
		if api.Namespaces[i].Name == "SampleNamespace" {
			sampleNs = &api.Namespaces[i]
			break
		}
	}

	require.NotNil(t, sampleNs, "Should find SampleNamespace")

	// Verify classes in namespace
	assert.GreaterOrEqual(t, len(sampleNs.Classes), 2, "Should have at least 2 types (class and interface)")

	// Find SampleClass
	var sampleClass *Class
	for i := range sampleNs.Classes {
		if sampleNs.Classes[i].Name == "SampleClass" {
			sampleClass = &sampleNs.Classes[i]
			break
		}
	}

	require.NotNil(t, sampleClass, "Should find SampleClass")
	assert.True(t, sampleClass.IsPublic, "SampleClass should be public")
	assert.False(t, sampleClass.IsInterface, "SampleClass should not be an interface")

	// Verify methods
	assert.GreaterOrEqual(t, len(sampleClass.Methods), 2, "Should have at least 2 methods")

	// Find Add method
	var addMethod *Method
	for i := range sampleClass.Methods {
		if sampleClass.Methods[i].Name == "Add" {
			addMethod = &sampleClass.Methods[i]
			break
		}
	}

	require.NotNil(t, addMethod, "Should find Add method")
	assert.True(t, addMethod.IsPublic, "Add method should be public")
	assert.Equal(t, 2, len(addMethod.Parameters), "Add method should have 2 parameters")
	assert.Contains(t, addMethod.Signature, "Add", "Signature should contain method name")
	assert.Contains(t, addMethod.Signature, "int", "Signature should contain return type")

	// Find static method
	var formatMethod *Method
	for i := range sampleClass.Methods {
		if sampleClass.Methods[i].Name == "FormatString" {
			formatMethod = &sampleClass.Methods[i]
			break
		}
	}

	require.NotNil(t, formatMethod, "Should find FormatString method")
	assert.True(t, formatMethod.IsStatic, "FormatString should be static")

	// Verify properties
	assert.GreaterOrEqual(t, len(sampleClass.Properties), 1, "Should have at least 1 property")

	// Find Name property
	var nameProp *Property
	for i := range sampleClass.Properties {
		if sampleClass.Properties[i].Name == "Name" {
			nameProp = &sampleClass.Properties[i]
			break
		}
	}

	require.NotNil(t, nameProp, "Should find Name property")
	assert.True(t, nameProp.IsPublic, "Name property should be public")
	assert.Equal(t, "string", nameProp.Type, "Name property should be of type string")

	// Find interface
	var sampleInterface *Class
	for i := range sampleNs.Classes {
		if sampleNs.Classes[i].Name == "ISampleInterface" {
			sampleInterface = &sampleNs.Classes[i]
			break
		}
	}

	require.NotNil(t, sampleInterface, "Should find ISampleInterface")
	assert.True(t, sampleInterface.IsInterface, "ISampleInterface should be marked as interface")

	// Verify global class (outside namespace)
	var globalNs *Namespace
	for i := range api.Namespaces {
		if api.Namespaces[i].Name == "<global>" {
			globalNs = &api.Namespaces[i]
			break
		}
	}

	if globalNs != nil {
		var globalClass *Class
		for i := range globalNs.Classes {
			if globalNs.Classes[i].Name == "GlobalClass" {
				globalClass = &globalNs.Classes[i]
				break
			}
		}

		if globalClass != nil {
			assert.True(t, globalClass.IsPublic, "GlobalClass should be public")
			assert.GreaterOrEqual(t, len(globalClass.Methods), 1, "GlobalClass should have at least 1 method")
		}
	}
}

func TestCSharpParser_EmptySource(t *testing.T) {
	// Arrange
	parser := New()

	// Act
	api, err := parser.ExtractAPISurface(context.Background(), "")

	// Assert
	require.NoError(t, err, "Should not error on empty source")
	require.NotNil(t, api, "Should return non-nil API surface")
	assert.Empty(t, api.Namespaces, "Should have no namespaces for empty source")
}

func TestCSharpParser_InvalidSyntax(t *testing.T) {
	// Arrange
	parser := New()
	invalidCode := `
	public class { // Missing class name
		public void Method() {
			// Unclosed brace
	}`

	// Act - Parser should still attempt to extract what it can
	api, err := parser.ExtractAPISurface(context.Background(), invalidCode)

	// Assert
	require.NoError(t, err, "Tree-sitter should handle invalid syntax gracefully")
	require.NotNil(t, api, "Should return API surface even for invalid code")
}
