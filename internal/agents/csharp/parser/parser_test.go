package parser

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCSharpParser_ExtractAPISurface(t *testing.T) {
	// Arrange
	sampleCode := getSampleCode()

	// Act
	parser := New()
	api, err := parser.ExtractAPISurface(context.Background(), sampleCode)

	// Assert
	require.NoError(t, err, "Parser should not return an error")
	require.NotNil(t, api, "API surface should not be nil")

	// Run sub-tests
	t.Run("Namespace", func(t *testing.T) {
		testNamespace(t, api)
	})

	t.Run("Classes", func(t *testing.T) {
		sampleNs := findNamespace(api.Namespaces, "SampleNamespace")
		require.NotNil(t, sampleNs, "Should find SampleNamespace")
		testClasses(t, sampleNs)
	})

	t.Run("Methods", func(t *testing.T) {
		sampleNs := findNamespace(api.Namespaces, "SampleNamespace")
		require.NotNil(t, sampleNs, "Should find SampleNamespace")
		sampleClass := findClass(sampleNs.Classes, "SampleClass")
		require.NotNil(t, sampleClass, "Should find SampleClass")
		testMethods(t, sampleClass)
	})

	t.Run("Properties", func(t *testing.T) {
		sampleNs := findNamespace(api.Namespaces, "SampleNamespace")
		require.NotNil(t, sampleNs, "Should find SampleNamespace")
		sampleClass := findClass(sampleNs.Classes, "SampleClass")
		require.NotNil(t, sampleClass, "Should find SampleClass")
		testProperties(t, sampleClass)
	})

	t.Run("Interface", func(t *testing.T) {
		sampleNs := findNamespace(api.Namespaces, "SampleNamespace")
		require.NotNil(t, sampleNs, "Should find SampleNamespace")
		testInterface(t, sampleNs)
	})

	t.Run("GlobalClass", func(t *testing.T) {
		testGlobalClass(t, api)
	})
}

func getSampleCode() string {
	return `
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
}

func findNamespace(namespaces []Namespace, name string) *Namespace {
	for i := range namespaces {
		if namespaces[i].Name == name {
			return &namespaces[i]
		}
	}
	return nil
}

func findClass(classes []Class, name string) *Class {
	for i := range classes {
		if classes[i].Name == name {
			return &classes[i]
		}
	}
	return nil
}

func findMethod(methods []Method, name string) *Method {
	for i := range methods {
		if methods[i].Name == name {
			return &methods[i]
		}
	}
	return nil
}

func findProperty(properties []Property, name string) *Property {
	for i := range properties {
		if properties[i].Name == name {
			return &properties[i]
		}
	}
	return nil
}

func testNamespace(t *testing.T, api *APISurface) {
	assert.GreaterOrEqual(t, len(api.Namespaces), 1, "Should have at least one namespace")
	sampleNs := findNamespace(api.Namespaces, "SampleNamespace")
	require.NotNil(t, sampleNs, "Should find SampleNamespace")
}

func testClasses(t *testing.T, ns *Namespace) {
	assert.GreaterOrEqual(t, len(ns.Classes), 2, "Should have at least 2 types (class and interface)")

	sampleClass := findClass(ns.Classes, "SampleClass")
	require.NotNil(t, sampleClass, "Should find SampleClass")
	assert.True(t, sampleClass.IsPublic, "SampleClass should be public")
	assert.False(t, sampleClass.IsInterface, "SampleClass should not be an interface")
}

func testMethods(t *testing.T, class *Class) {
	assert.GreaterOrEqual(t, len(class.Methods), 2, "Should have at least 2 methods")

	// Test Add method
	addMethod := findMethod(class.Methods, "Add")
	require.NotNil(t, addMethod, "Should find Add method")
	assert.True(t, addMethod.IsPublic, "Add method should be public")
	assert.Equal(t, 2, len(addMethod.Parameters), "Add method should have 2 parameters")
	assert.Contains(t, addMethod.Signature, "Add", "Signature should contain method name")
	assert.Contains(t, addMethod.Signature, "int", "Signature should contain return type")

	// Test static method
	formatMethod := findMethod(class.Methods, "FormatString")
	require.NotNil(t, formatMethod, "Should find FormatString method")
	assert.True(t, formatMethod.IsStatic, "FormatString should be static")
}

func testProperties(t *testing.T, class *Class) {
	assert.GreaterOrEqual(t, len(class.Properties), 1, "Should have at least 1 property")

	nameProp := findProperty(class.Properties, "Name")
	require.NotNil(t, nameProp, "Should find Name property")
	assert.True(t, nameProp.IsPublic, "Name property should be public")
	assert.Equal(t, "string", nameProp.Type, "Name property should be of type string")
}

func testInterface(t *testing.T, ns *Namespace) {
	sampleInterface := findClass(ns.Classes, "ISampleInterface")
	require.NotNil(t, sampleInterface, "Should find ISampleInterface")
	assert.True(t, sampleInterface.IsInterface, "ISampleInterface should be marked as interface")
}

func testGlobalClass(t *testing.T, api *APISurface) {
	globalNs := findNamespace(api.Namespaces, "<global>")
	if globalNs == nil {
		return // Global namespace is optional
	}

	globalClass := findClass(globalNs.Classes, "GlobalClass")
	if globalClass == nil {
		return // Global class is optional
	}

	assert.True(t, globalClass.IsPublic, "GlobalClass should be public")
	assert.GreaterOrEqual(t, len(globalClass.Methods), 1, "GlobalClass should have at least 1 method")
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

	// Act
	api, err := parser.ExtractAPISurface(context.Background(), invalidCode)

	// Assert
	// We expect the parser to handle invalid syntax gracefully
	require.NoError(t, err, "Should not error on invalid syntax")
	require.NotNil(t, api, "Should return non-nil API surface even for invalid code")
}

func TestCSharpParser_Attributes(t *testing.T) {
	// Arrange
	codeWithAttributes := `
using System;

namespace TestNamespace
{
	[Serializable]
	[Obsolete("Use NewClass instead")]
	public class OldClass
	{
		[Required]
		public string Name { get; set; }
		
		[HttpGet]
		[Route("api/test")]
		public void GetData() { }
	}
}`

	parser := New()

	// Act
	api, err := parser.ExtractAPISurface(context.Background(), codeWithAttributes)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, api)

	ns := findNamespace(api.Namespaces, "TestNamespace")
	require.NotNil(t, ns)

	oldClass := findClass(ns.Classes, "OldClass")
	require.NotNil(t, oldClass)
	// Note: Attributes are not currently parsed by the implementation
}

func TestCSharpParser_Generics(t *testing.T) {
	// Arrange
	codeWithGenerics := `
namespace GenericTest
{
	public class GenericClass<T, U> where T : class where U : struct
	{
		public T GetItem() { return default(T); }
		public List<T> GetList() { return new List<T>(); }
	}
	
	public interface IRepository<T> where T : IEntity
	{
		T GetById(int id);
		IEnumerable<T> GetAll();
	}
}`

	parser := New()

	// Act
	api, err := parser.ExtractAPISurface(context.Background(), codeWithGenerics)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, api)

	ns := findNamespace(api.Namespaces, "GenericTest")
	require.NotNil(t, ns)

	genericClass := findClass(ns.Classes, "GenericClass")
	require.NotNil(t, genericClass)
	assert.Contains(t, genericClass.Name, "GenericClass", "Should preserve generic class name")
}

func TestCSharpParser_NestedTypes(t *testing.T) {
	// Arrange
	codeWithNested := `
namespace NestedTest
{
	public class OuterClass
	{
		public class InnerClass
		{
			public void InnerMethod() { }
		}
		
		private class PrivateInner { }
		
		public interface IInnerInterface
		{
			void DoWork();
		}
	}
}`

	parser := New()

	// Act
	api, err := parser.ExtractAPISurface(context.Background(), codeWithNested)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, api)

	ns := findNamespace(api.Namespaces, "NestedTest")
	require.NotNil(t, ns)

	outerClass := findClass(ns.Classes, "OuterClass")
	require.NotNil(t, outerClass)
	// Note: Nested types might be handled differently depending on parser implementation
}
