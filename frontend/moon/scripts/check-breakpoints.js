#!/usr/bin/env node
'use strict';

/**
 * Breakpoint checker. Plain node, no dependencies.
 *
 *   node scripts/check-breakpoints.js [app-root]
 *
 * The rules:
 *   1. src/common/env-util.ts owns the boundary; src/styles.css must open a
 *      compact block at that width and a wide block one pixel above it.
 *   2. No stylesheet under src/ may declare a width media query at any other
 *      value. Compact is 768, wide is 769, and nothing else exists.
 *   3. Nothing outside the service may read the viewport width directly.
 *
 * Exits non-zero and names every violation by file and line.
 */

const fs = require('fs');
const path = require('path');

const DEFAULT_ROOT = path.resolve(__dirname, '..');
const ROOT = path.resolve(process.argv[2] || DEFAULT_ROOT);

const SERVICE_FILE = 'src/common/env-util.ts';
const STYLES_FILE = 'src/styles.css';

const violations = [];

function report(file, line, message) {
  violations.push({ file, line, message });
}

function lineAt(contents, index) {
  return contents.slice(0, index).split('\n').length;
}

function relativePath(absolute) {
  return path.relative(ROOT, absolute).split(path.sep).join('/');
}

function read(relative) {
  const absolute = path.join(ROOT, relative);
  if (!fs.existsSync(absolute)) {
    report(relative, 1, 'file is missing, so the rule cannot be checked');
    return null;
  }
  return fs.readFileSync(absolute, 'utf8');
}

function walk(directory, collected) {
  if (!fs.existsSync(directory)) {
    return collected;
  }
  fs.readdirSync(directory, { withFileTypes: true }).forEach((entry) => {
    const absolute = path.join(directory, entry.name);
    if (entry.isDirectory()) {
      walk(absolute, collected);
    } else if (entry.isFile()) {
      collected.push(absolute);
    }
  });
  return collected;
}

function filesWithExtensions(extensions) {
  const src = path.join(ROOT, 'src');
  return walk(src, []).filter((absolute) => extensions.indexOf(path.extname(absolute)) >= 0);
}

function boundaryFromService() {
  const contents = read(SERVICE_FILE);
  if (contents === null) {
    return null;
  }
  const match = /\bMOBILE_MAX_WIDTH\s*=\s*(\d+)\b/.exec(contents);
  if (!match) {
    report(SERVICE_FILE, 1, 'MOBILE_MAX_WIDTH is not declared as a number');
    return null;
  }
  return Number(match[1]);
}

/** The compact max-width and wide min-width blocks in the global stylesheet. */
function tiersFromStyles() {
  const compact = [];
  const wide = [];
  const contents = read(STYLES_FILE);
  if (contents === null) {
    return { compact, wide };
  }

  const media = /@media\b([^{]*)\{/g;
  let match;
  while ((match = media.exec(contents)) !== null) {
    const prelude = match[1];
    const line = lineAt(contents, match.index);
    const max = /\(\s*max-width\s*:\s*(\d+(?:\.\d+)?)px\s*\)/.exec(prelude);
    const min = /\(\s*min-width\s*:\s*(\d+(?:\.\d+)?)px\s*\)/.exec(prelude);
    if (max) {
      compact.push({ width: Number(max[1]), line });
    }
    if (min) {
      wide.push({ width: Number(min[1]), line });
    }
  }

  return { compact, wide };
}

function checkTiers(boundary) {
  const tiers = tiersFromStyles();

  if (tiers.compact.length !== 1) {
    report(STYLES_FILE, tiers.compact.length ? tiers.compact[0].line : 1,
      `expected one compact (max-width) block, found ${tiers.compact.length}`);
  }
  if (tiers.wide.length !== 1) {
    report(STYLES_FILE, tiers.wide.length ? tiers.wide[0].line : 1,
      `expected one wide (min-width) block, found ${tiers.wide.length}`);
  }
  if (boundary === null) {
    return;
  }
  const compact = tiers.compact[0];
  const wide = tiers.wide[0];
  if (compact && compact.width !== boundary) {
    report(STYLES_FILE, compact.line,
      `compact block is ${compact.width}px but MOBILE_MAX_WIDTH is ${boundary}`);
  }
  if (wide && wide.width !== boundary + 1) {
    report(STYLES_FILE, wide.line,
      `wide block is ${wide.width}px but must start at ${boundary + 1}px, one above MOBILE_MAX_WIDTH`);
  }
}

function checkMediaWidths(boundary) {
  if (boundary === null) {
    return;
  }
  const canonical = [boundary, boundary + 1];

  filesWithExtensions(['.css']).forEach((absolute) => {
    const relative = relativePath(absolute);
    const contents = fs.readFileSync(absolute, 'utf8');
    const media = /@media\b([^{]*)\{/g;
    let match;
    while ((match = media.exec(contents)) !== null) {
      const prelude = match[1];
      // match[0] spells the prelude with '@media' in front and '{' behind, so
      // its own offset in the file is match.index plus that lead-in.
      const preludeStart = match.index + (match[0].length - prelude.length - 1);
      const width = /\bwidth\s*[:<>=]+\s*(\d+(?:\.\d+)?)px/g;
      let value;
      while ((value = width.exec(prelude)) !== null) {
        const pixels = Number(value[1]);
        if (canonical.indexOf(pixels) < 0) {
          report(relative, lineAt(contents, preludeStart + value.index),
            `width media query ${value[1]}px is not ${canonical[0]}px or ${canonical[1]}px`);
        }
      }
    }
  });
}

function checkWindowReads() {
  const windowRead = /\bwindow\.innerWidth\b|\bdocument\.documentElement\.clientWidth\b/g;

  filesWithExtensions(['.ts', '.html']).forEach((absolute) => {
    const relative = relativePath(absolute);
    if (relative === SERVICE_FILE) {
      return;
    }
    const contents = fs.readFileSync(absolute, 'utf8');
    let match;
    while ((match = windowRead.exec(contents)) !== null) {
      report(relative, lineAt(contents, match.index),
        `${match[0]} reads the viewport directly; ask Env.isMobile() or Env.mobile$ instead`);
    }
  });
}

const boundary = boundaryFromService();
checkTiers(boundary);
checkMediaWidths(boundary);
checkWindowReads();

if (violations.length > 0) {
  violations
    .sort((a, b) => (a.file === b.file ? a.line - b.line : a.file.localeCompare(b.file)))
    .forEach((violation) => {
      console.error(`${violation.file}:${violation.line}  ${violation.message}`);
    });
  console.error(`\n${violations.length} breakpoint violation(s).`);
  process.exit(1);
}

console.log(`Breakpoints OK: compact <= ${boundary}px, wide >= ${boundary + 1}px.`);
