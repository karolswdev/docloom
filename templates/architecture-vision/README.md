Architecture Vision HTML Template
=================================

This template mirrors the "Architecture Vision" PDF layout with:

- Fixed header with brand, title, project, author
- Fixed footer with website, copyright, and document code
- Section structure matching the PDF
- Letter-sized print styles with page breaks between sections
- Simple JSON-based filler for LLM or automation

Files

- `architecture-vision.html`: The HTML skeleton with data-field placeholders
- `terumo.css`: Brand tokens and theme (from Terumo BCT guidelines)
- `style.css`: Document layout styles for screen and print
- `fill.js`: Helper to populate placeholders from a JSON object

Usage

1) Open `architecture-vision.html` in a browser directly, or serve the folder with any static server.
2) To auto-fill, define `window.DOC_DATA` before the `fill.js` script loads, or call `window.DocFill.fill(data)` at runtime.

Example

```html
<script>
  window.DOC_DATA = {
    project_name: "NextGen EHR Integration",
    author: "Jane Doe",
    summary: "<p>This initiative unifies...</p>",
    introduction: "<p>We aim to...</p>",
    problem_description: "...",
    stakeholders: "<ul><li>Clinical Ops</li><li>IT</li></ul>",
    business_scenarios: ["Clinician orders", "Lab response"],
    business_alignment: "...",
    scope: "...",
    dependencies: "...",
    requirements: "...",
    architecture_objective: "...",
    critical_issues_and_risks: "...",
    constraints: "...",
    assumptions: "...",
    architecture_approach: "...",
    current_architecture: "...",
    to_be_architecture: "...",
    references: ["HL7 FHIR R4", "Internal SSO Docs"],
    copyright_year: 2025,
    doc_code: "AV-2025-00012"
  };
</script>
<script src="fill.js"></script>
```

Branding

- Replace the placeholder square in the header with your logo by styling `.logo` with a background-image, or by placing an `<img>` inside `.logo` and adjusting CSS.
- Colors and typography come from `terumo.css`. The document layout in `style.css` maps to those tokens (e.g., uses `--terumo-color-green` for accents).

Printing

- The template is optimized for Letter size (612×792 pt). Use browser print to export to PDF; ensure margins are set to default and background graphics are enabled for the header/footers to show.

Notes

- Corporate address blocks from the PDF are intentionally omitted, per guidelines.
