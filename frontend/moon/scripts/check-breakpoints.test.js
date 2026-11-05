#!/usr/bin/env node
'use strict';

/**
 * Checks check-breakpoints.js itself. Plain node, no dependency:
 *   node scripts/check-breakpoints.test.js
 */

const assert = require('assert');
const fs = require('fs');
const os = require('os');
const path = require('path');
const { spawnSync } = require('child_process');

const APP_ROOT = path.resolve(__dirname, '..');
const CHECKER = path.join(__dirname, 'check-breakpoints.js');

const GOOD_FILES = {
  'src/common/env-util.ts': 'export const MOBILE_MAX_WIDTH = 768;\n',
  'src/styles.css': [
    '@media only screen and (max-width: 768px) {',
    '  .only-wide { display: none !important; }',
    '}',
    '@media only screen and (min-width: 769px) {',
    '  .only-compact { display: none !important; }',
    '}',
    '',
  ].join('\n'),
  'src/app/gallery/gallery.component.css': '@media (max-width: 768px) { .cell { color: red; } }\n',
  'src/app/gallery/gallery.component.ts': "const label = 'ok';\n",
  'src/app/gallery/gallery.component.html': '<p>ok</p>\n',
};

function writeFixture(files) {
  const dir = fs.mkdtempSync(path.join(os.tmpdir(), 'breakpoints-'));
  Object.keys(files).forEach((relative) => {
    const target = path.join(dir, relative);
    fs.mkdirSync(path.dirname(target), { recursive: true });
    fs.writeFileSync(target, files[relative]);
  });
  return dir;
}

function run(root) {
  return spawnSync(process.execPath, [CHECKER, root], { encoding: 'utf8' });
}

function check(label, files, expectedStatus, expectedFragment) {
  const dir = writeFixture(files);
  try {
    const result = run(dir);
    const output = `${result.stdout}\n${result.stderr}`;
    if (expectedStatus === 0) {
      assert.strictEqual(result.status, 0, `${label}: expected exit 0, got ${result.status}\n${output}`);
    } else {
      assert.notStrictEqual(result.status, 0, `${label}: expected a non-zero exit\n${output}`);
    }
    if (expectedFragment) {
      assert.ok(output.indexOf(expectedFragment) >= 0,
        `${label}: expected output to name ${expectedFragment}\n${output}`);
    }
    console.log(`ok - ${label}`);
  } finally {
    fs.rmSync(dir, { recursive: true, force: true });
  }
}

function withFile(files, relative, contents) {
  const copy = Object.assign({}, files);
  copy[relative] = contents;
  return copy;
}

// The repo as shipped must satisfy its own rule.
const repo = run(APP_ROOT);
assert.strictEqual(repo.status, 0, `this repo: expected exit 0\n${repo.stdout}\n${repo.stderr}`);
console.log('ok - this repo passes');

// A good fixture: canonical blocks, a component query at 768, no window read.
check('good fixture', GOOD_FILES, 0);

// A component query at a width that is neither 768 nor 769.
check('bad fixture: 600px media query', withFile(GOOD_FILES,
  'src/app/gallery/gallery.component.css',
  '@media (max-width: 600px) { .cell { color: red; } }\n'), 1,
'src/app/gallery/gallery.component.css:1');

// A decimal width is not a width the tiers allow, so it is checked like any other.
check('bad fixture: decimal media query', withFile(GOOD_FILES,
  'src/app/gallery/gallery.component.css',
  '@media (max-width: 767.98px) { .cell { color: red; } }\n'), 1,
'src/app/gallery/gallery.component.css:1');

// The offending width is reported on its own line, not one shifted by the prelude.
check('bad fixture: width wrapped onto its own line', withFile(GOOD_FILES,
  'src/app/gallery/gallery.component.css',
  '/* header */\n@media\n(max-width: 600px) { .cell { color: red; } }\n'), 1,
'src/app/gallery/gallery.component.css:3');

// The stylesheet wide block disagreeing with the service boundary.
check('bad fixture: stylesheet says 767', withFile(GOOD_FILES,
  'src/styles.css',
  '@media only screen and (max-width: 767px) { .only-wide { display: none !important; } }\n' +
  '@media only screen and (min-width: 768px) { .only-compact { display: none !important; } }\n'),
1, 'src/styles.css:1');

// A direct window read outside the service.
check('bad fixture: window.innerWidth outside the service', withFile(GOOD_FILES,
  'src/app/gallery/gallery.component.ts',
  "const compact = window.innerWidth <= 768;\n"), 1,
'src/app/gallery/gallery.component.ts:1');

console.log('\ncheck-breakpoints: all cases passed');
