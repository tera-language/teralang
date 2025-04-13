import fs from 'fs/promises';
import path from 'path';
import chalk from 'chalk';

const terraBrandColor = chalk.hex('#1E90FF');
const successColor = chalk.green;
const errorColor = chalk.red;
const infoColor = chalk.cyan;

function logSuccess(msg) {
  console.log(terraBrandColor('[Tera] ') + successColor(msg));
}
function logError(msg) {
  console.error(terraBrandColor('[Tera] ') + errorColor(msg));
}
function logWarning(msg) {
  console.warn(terraBrandColor('[Tera] ') + chalk.yellow(msg));
}
function logInfo(msg) {
  console.log(terraBrandColor('[Tera] ') + infoColor(msg));
}

class TreeNode {
  constructor(path, method) {
    this.path = path;
    this.method = method;
    this.children = {};
    this.data = null;
  }

  addChild(node) {
    this.children[node.path] = node;
  }

  setData(data) {
    this.data = data;
  }

  getData() {
    return this.data;
  }

  getChildren() {
    return this.children;
  }
}

const importedFiles = new Set(); 

export async function parseTeraFile(filePath, basePath = '') {
  const absolutePath = path.resolve(basePath, filePath);
  if (importedFiles.has(absolutePath)) {
    logInfo(`File already imported: ${absolutePath}`); 
    return new TreeNode('root', 'root'); 
  }

  importedFiles.add(absolutePath);

  const content = await fs.readFile(absolutePath, 'utf-8');
  const lines = content.split('\n');
  const rootNode = new TreeNode('root', 'root');

  let collecting = false;
  let braceCount = 0;
  let buffer = '';
  let routeKey = '';
  let lineStart = 0;
  let parentNode = rootNode;
  let routePath = '';
  let routeMethod = '';

  for (let i = 0; i < lines.length; i++) {
    const line = lines[i];
    const trimmed = line.trim();
    if (trimmed.startsWith('import') || trimmed === '') {
      if (trimmed.startsWith('import')) {
        const match = trimmed.match(/^import\s+"(.+?)"\s*$/);
        if (match) {
          let importPath = match[1];
          if (!importPath.endsWith('.tera')) {
            importPath += '.tera';
          }
          logInfo(`Importing file: ${importPath}`);
          try {
            const importedNode = await parseTeraFile(importPath, path.dirname(absolutePath));
            for (const [key, node] of Object.entries(importedNode.getChildren())) {
              rootNode.addChild(node);
            }
          } catch (err) {
            logError(`Failed to import file: ${importPath}`);
            logError(`Error: ${err.message}`);
          }
        } else {
          logWarning(`Skipped invalid import on line ${i + 1}: "${trimmed}"`);
        }
      }
      continue;
    }

    if (!collecting && trimmed.startsWith('route')) {
      const match = trimmed.match(/route\s+"(.+?)"\s+(\w+):\s*{/);
      if (match) {
        [ , routePath, routeMethod ] = match;
        routeKey = `${routeMethod.toUpperCase()} ${routePath}`;
        collecting = true;
        braceCount = 1;
        buffer = '';
        lineStart = i;

        const newNode = new TreeNode(routePath, routeMethod.toUpperCase());
        parentNode.addChild(newNode);
        parentNode = newNode;

        continue;
      } else {
        logWarning(`Skipped invalid route definition on line ${i + 1}: "${trimmed}"`);
        continue;
      }
    }

    if (collecting) {
      const noComment = line.split('//')[0];
      braceCount += (noComment.match(/{/g) || []).length;
      braceCount -= (noComment.match(/}/g) || []).length;
      buffer += noComment + '\n';

      if (braceCount === 0) {
        try {
          const parsed = parseRouteData(buffer.trim());
          parentNode.setData(parsed);
          logSuccess(`Parsed route: ${routeKey}`);
        } catch (err) {
          logError(`Failed to parse route: ${routeKey}`);
          logError(`In file: ${path.basename(absolutePath)} @ line ${lineStart + 1}`);
          logError(`Error: ${err.message}`);
        }

        collecting = false;
        buffer = '';
        parentNode = rootNode;
      }
    }
  }

  return rootNode;
}

function parseRouteData(rawBlock) {
  const routeData = {};
  const lines = rawBlock.replace(/^{|}$/g, '').split('\n');

  let currentKey = null;
  let buffer = '';
  let braceDepth = 0;

  for (let line of lines) {
    const trimmed = line.trim();
    if (!trimmed || trimmed.startsWith('//')) continue;

    if (currentKey) {
      braceDepth += (trimmed.match(/{/g) || []).length;
      braceDepth -= (trimmed.match(/}/g) || []).length;
      buffer += trimmed + '\n';

      if (braceDepth === 0) {
        try {
          if (currentKey === 'json') {
            routeData[currentKey] = parseJsonLike(buffer);
          } else {
            routeData[currentKey] = buffer.trim();
          }
        } catch (err) {
          throw new Error(`Invalid ${currentKey} block: ${err.message}`);
        }
        currentKey = null;
        buffer = '';
      }
      continue;
    }

    const match = trimmed.match(/^(\w+):\s*(.+)?$/);
    if (match) {
      const [ , key, value ] = match;

      if (value === '{') {
        currentKey = key;
        buffer = '{\n';
        braceDepth = 1;
      } else {
        if (key === 'json') {
          routeData[key] = parseJsonLike(value);
        } else {
          routeData[key] = value.trim();
        }
      }
    }
  }

  return routeData;
}

function parseJsonLike(value) {
  const full = value.trim().startsWith('{') ? value : `{${value}}`;
  const jsonified = full
    .replace(/(\w+):/g, '"$1":')
    .replace(/'/g, '"')
    .replace(/,\s*}/g, '}')
    .replace(/,\s*]/g, ']');

  return JSON.parse(jsonified);
}
