# docloom

Beautiful, template-driven technical documentation — fast.

Docloom is a system for technical folks to generate high‑quality documents by combining structured templates with source materials and model‑assisted content. The aim is consistent, branded, and reviewable outputs that you can print, share, and iterate on quickly.

Docloom is currently in early development. The Software Requirements Specification (SRS) lives in `docs/` and a working HTML template prototype is in `source-example-template-that-spawned-the-idea/`.

## What It Is

- **Templates:** Curated, branded skeletons (e.g., Architecture Vision) with `data-field` placeholders for structured content.
- **Fields & Schema:** Each template defines a field schema; AI generations and manual edits populate these fields.
- **Sources:** Local Markdown, text, and other artifacts provide the factual basis for generation.
- **Renderers:** Produce attractive HTML ready for print/PDF, plus a JSON sidecar of the filled fields for traceability.

## Current Status

- The SRS defines the planned CLI, behaviors, and extensibility. See `docs/SRS.md`.
- A sample HTML template (Architecture Vision) demonstrates the look, layout, and simple JSON-based filling.
- CLI and end‑to‑end pipeline are under active design per the SRS.

## Repository Layout

- `docs/` — Authoritative SRS: scope, terminology, and requirements.
- `source-example-template-that-spawned-the-idea/` — Prototype HTML template, styles, and a simple filler script:
  - `architecture-vision.html` — HTML skeleton with `data-field` placeholders
  - `style.css`, `terumo.css`, `logo.svg` — Layout and brand styling
  - `fill.js` — Minimal helper to populate placeholders from JSON

## Try the Sample Template

You can open the template directly in a browser or serve it locally to test auto‑filling.

1) Change into the sample directory:
   - `cd source-example-template-that-spawned-the-idea`
2) Serve the folder (pick one):
   - `python3 -m http.server 8080` (then open `http://localhost:8080/architecture-vision.html`)
   - or simply double‑click `architecture-vision.html` to open it directly
3) To auto‑fill fields, define `window.DOC_DATA` before `fill.js` or call `window.DocFill.fill(data)` in the console.

Example snippet to embed before `fill.js` in the HTML (or paste in DevTools console and then call `DocFill.fill`):

```html
<script>
  window.DOC_DATA = {
    project_name: "NextGen EHR Integration",
    author: "Jane Doe",
    summary: "<p>This initiative unifies...</p>",
    introduction: "<p>We aim to...</p>",
    // ...other fields...
    copyright_year: 2025,
    doc_code: "AV-2025-00012"
  };
</script>
<script src="fill.js"></script>
```

Printing tips:
- Use the browser’s Print dialog to export to PDF.
- Enable printing of backgrounds so header/footer visuals appear.

## Usage

To see available commands and options:

```bash
docloom --help
docloom generate --help
docloom templates list
```

### Generating Documents

To generate a document from your source materials:

```bash
# Basic usage with API key via environment variable
export OPENAI_API_KEY="your-api-key-here"
docloom generate \
  --type architecture-vision \
  --source ./docs \
  --out output.html

# Or provide API key directly
docloom generate \
  --type architecture-vision \
  --source ./docs \
  --out output.html \
  --api-key "your-api-key-here"

# Using multiple source paths
docloom generate \
  --type technical-debt-summary \
  --source ./docs --source ./notes.md \
  --out debt-report.html

# Dry-run mode to preview without API calls
docloom generate \
  --type architecture-vision \
  --source ./docs \
  --out output.html \
  --dry-run

# Force overwrite existing output
docloom generate \
  --type architecture-vision \
  --source ./docs \
  --out output.html \
  --force
```

## Configuration

Docloom supports configuration through multiple sources with the following precedence order (highest to lowest):

1. **CLI flags** - Command-line arguments passed directly
2. **Environment variables** - Variables prefixed with `DOCLOOM_`
3. **Configuration file** - YAML file specified with `--config`
4. **Defaults** - Built-in default values

### Example Configuration File (docloom.yaml)

```yaml
# Model configuration
model: gpt-4
base_url: https://api.openai.com/v1
temperature: 0.7
max_retries: 3

# Template configuration
template_dir: ./templates

# Output configuration
force: false

# Operational configuration
verbose: false
dry_run: false
```

### Environment Variables

- `OPENAI_API_KEY` - API key for OpenAI or compatible services (required for generation)
- `DOCLOOM_MODEL` - AI model to use (default: gpt-4)
- `DOCLOOM_BASE_URL` - Base URL for OpenAI-compatible API (default: https://api.openai.com/v1)
- `DOCLOOM_API_KEY` - Alternative to `OPENAI_API_KEY`
- `DOCLOOM_TEMPERATURE` - Temperature for model generation (0.0-1.0, default: 0.7)
- `DOCLOOM_TEMPLATE_DIR` - Directory containing custom templates
- `DOCLOOM_VERBOSE` - Enable verbose logging (shows detailed progress)
- `DOCLOOM_DRY_RUN` - Enable dry-run mode (preview without API calls)

**Note:** API keys are never logged or included in output files. The system automatically redacts them from any debug output.

## Planned CLI (from SRS)

The SRS outlines a CLI along these lines:

- `docloom generate --type <template> --source <paths...> --out <file>`
- Template registry (e.g., `architecture-vision`, `technical-debt-summary`, `reference-architecture`).
- HTML output plus a JSON sidecar of filled fields.
- Config via `docloom.yaml`, env vars, and flags; verbose/dry‑run modes; schema validation and repair.

For full details and requirement IDs, see `docs/SRS.md`.

## Container Usage

### Building the Docker Image

To build the docloom Docker image locally:

```bash
docker build -t docloom:latest .
```

### Running with Docker

To generate documents using the Docker container:

```bash
# Mount your source directory and run docloom
docker run --rm \
  -v $(pwd)/sources:/workspace/sources \
  -v $(pwd)/output:/workspace/output \
  docloom:latest generate \
  --type architecture-vision \
  --source /workspace/sources \
  --out /workspace/output/document.html
```

### Using Pre-built Images

Pre-built images will be available from GitHub Container Registry:

```bash
docker pull ghcr.io/karolswdev/docloom:latest
```

## Contributing

- Start by reading `docs/SRS.md` to understand scope and constraints.
- Explore the template prototype to get a feel for structure and fields.
- Proposals, templates, and small improvements are welcome as issues or PRs.

## License

TBD.
