# DocLoom Templates

This directory contains document templates for use with DocLoom. Templates define the structure, styling, and intelligence prompts for generating different types of technical documentation.

## Template Structure

Each template consists of:
1. **HTML file** - The document structure with data-field placeholders
2. **CSS files** - Styling for the document
3. **Definition file** (optional) - Metadata and AI analysis prompts
4. **Assets** - Images, logos, fonts, etc.

## Available Templates

### 1. Architecture Vision
Location: `architecture-vision/`
- Professional template for architecture documentation
- Includes sections for goals, constraints, decisions, and roadmap
- Print-optimized styling

### 2. Technical Specification
Location: `technical-spec/`
- Detailed technical documentation template
- Sections for requirements, design, implementation details
- Code-friendly formatting

### 3. API Documentation
Location: `api-docs/`
- RESTful API documentation template
- Endpoint listings, request/response examples
- Interactive styling

## Creating a New Template

### Step 1: Create Template Directory
```bash
mkdir templates/my-template
```

### Step 2: Create HTML Structure
Create `templates/my-template/template.html`:

```html
<!DOCTYPE html>
<html>
<head>
    <title>{{title}}</title>
    <link rel="stylesheet" href="style.css">
</head>
<body>
    <header>
        <h1 data-field="title">Document Title</h1>
        <div data-field="author">Author Name</div>
        <div data-field="date">Date</div>
    </header>
    
    <main>
        <section>
            <h2>Summary</h2>
            <div data-field="summary">Executive summary goes here</div>
        </section>
        
        <section>
            <h2>Content</h2>
            <div data-field="content">Main content goes here</div>
        </section>
    </main>
    
    <footer>
        <div data-field="copyright">© 2024</div>
    </footer>
</body>
</html>
```

### Step 3: Add Styling
Create `templates/my-template/style.css`:

```css
body {
    font-family: 'Segoe UI', system-ui, sans-serif;
    line-height: 1.6;
    color: #333;
    max-width: 1200px;
    margin: 0 auto;
    padding: 20px;
}

header {
    border-bottom: 2px solid #0066cc;
    padding-bottom: 20px;
    margin-bottom: 30px;
}

h1 {
    color: #0066cc;
    margin: 0;
}

/* Data field styling */
[data-field] {
    min-height: 1.5em;
    padding: 10px;
    border: 1px dashed #ddd;
    border-radius: 4px;
    background: #f9f9f9;
}

[data-field]:empty::before {
    content: attr(data-field);
    color: #999;
    font-style: italic;
}
```

### Step 4: Define Template Metadata (Optional)
Create `templates/my-template/template.yaml`:

```yaml
name: my-template
display_name: "My Custom Template"
description: "A template for my specific documentation needs"
version: "1.0.0"
author: "Your Name"

# Define the fields and their JSON schema
fields:
  title:
    type: string
    description: "Document title"
    required: true
  author:
    type: string
    description: "Author name"
    required: true
  date:
    type: string
    format: date
    description: "Document date"
  summary:
    type: string
    description: "Executive summary"
    format: html
  content:
    type: string
    description: "Main document content"
    format: html
  copyright:
    type: string
    description: "Copyright notice"
    default: "© 2024 Your Organization"

# AI Analysis Configuration (for agent-driven generation)
analysis:
  system_prompt: |
    You are a technical documentation expert. Your task is to analyze the provided 
    source materials and generate structured content for this document template.
    Focus on clarity, completeness, and technical accuracy.
  
  initial_prompt: |
    Please analyze the following repository/codebase and generate documentation
    that covers:
    1. High-level overview and purpose
    2. Key technical decisions and rationale
    3. Implementation details
    4. Future considerations
    
    Format your response as JSON matching the template field schema.
```

## Using Templates

### Basic Usage
```bash
# Use a built-in template
docloom generate --type architecture-vision --source ./docs --out output.html

# Use a custom template directory
docloom generate --template-dir ./my-templates --type my-template --source ./docs --out output.html
```

### With AI Agent
```bash
# Use with research agent for intelligent content generation
docloom generate --type technical-spec --source ./src --agent csharp-analyzer --out spec.html
```

## Template Best Practices

1. **Clear Field Names**: Use descriptive data-field names that indicate content type
2. **Responsive Design**: Ensure templates work on different screen sizes
3. **Print-Friendly**: Include print-specific CSS for PDF generation
4. **Semantic HTML**: Use proper HTML5 semantic elements
5. **Accessibility**: Include proper ARIA labels and semantic structure
6. **Schema Validation**: Define clear JSON schemas for fields
7. **AI Prompts**: Write specific, goal-oriented analysis prompts

## Examples

See the `examples/` subdirectory for complete template examples:
- `examples/simple/` - Minimal template with basic fields
- `examples/advanced/` - Complex template with nested sections
- `examples/ai-driven/` - Template optimized for AI-generated content

## Contributing Templates

To contribute a template to DocLoom:
1. Create your template following the structure above
2. Include a README.md in your template directory
3. Add example output (generated HTML)
4. Submit a pull request

## Template Registry

Templates can be registered with DocLoom for easy discovery:
- Built-in templates are in `internal/templates/defaults/`
- User templates go in `~/.docloom/templates/`
- Project templates go in `.docloom/templates/`

The template registry automatically discovers templates in these locations.