/*
  Simple filler: populates elements with data-field attributes from window.DOC_DATA if present.
  You can set window.DOC_DATA before this script or call fill(data) manually.
*/
(function () {
  function setNodeContent(node, value) {
    if (node == null) return;
    if (value == null) return;
    if (Array.isArray(value)) {
      if (node.tagName === 'OL' || node.querySelector('ol')) {
        const ol = node.tagName === 'OL' ? node : (node.querySelector('ol') || node);
        ol.innerHTML = value.map(v => `<li>${escapeHtml(v)}</li>`).join('');
      } else {
        node.innerHTML = value.map(escapeHtml).join('<br/>');
      }
      return;
    }
    if (typeof value === 'object') {
      node.textContent = JSON.stringify(value, null, 2);
      return;
    }
    // Default: treat as HTML if looks like HTML; else text
    if (/(<\w+[^>]*>)/.test(String(value))) {
      node.innerHTML = String(value);
    } else {
      node.textContent = String(value);
    }
  }
  }

  function escapeHtml(s) {
    return String(s)
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
      .replace(/"/g, '&quot;')
      .replace(/'/g, '&#39;');
  }

  function fill(data) {
    if (!data) return;
    // Simple mapping: find [data-field] nodes and set content
    document.querySelectorAll('[data-field]').forEach(node => {
      const key = node.getAttribute('data-field');
      if (!key) return;
      const value = data[key];
      if (value == null) return;
      setNodeContent(node, value);
    });
  }

  // Expose
  window.DocFill = { fill };

  // Auto-run
  if (window.DOC_DATA) {
    if (document.readyState === 'loading') {
      document.addEventListener('DOMContentLoaded', () => fill(window.DOC_DATA));
    } else {
      fill(window.DOC_DATA);
    }
  }
})();
