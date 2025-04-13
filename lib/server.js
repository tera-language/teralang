import http from 'http';
import fs from 'fs';
import { parseTeraFile } from '../parser/parseTera.js';
import mime from 'mime-types';
import chalk from 'chalk';

const terraBrandColor = chalk.hex('#1E90FF');
const successColor = chalk.green;
const errorColor = chalk.red;
const warningColor = chalk.yellow;
const infoColor = chalk.cyan;

// Logging utilities
function logSuccess(message) {
  console.log(terraBrandColor('[Tera] ') + successColor(message));
}
function logError(message) {
  console.error(terraBrandColor('[Tera] ') + errorColor(message));
}
function logWarning(message) {
  console.warn(terraBrandColor('[Tera] ') + warningColor(message));
}
function logInfo(message) {
  console.log(terraBrandColor('[Tera] ') + infoColor(message));
}

export async function serveTeraFile(path, port = 3000) {
  let root;
  try {
    root = await parseTeraFile(path);
  } catch (err) {
    logError(`Failed to parse Tera file: ${err.message}`);
    process.exit(1);
  }

  const routeMap = {};
  try {
    const children = root.getChildren();
    for (const childKey in children) {
      const node = children[childKey];
      const routeKey = `${node.method.toUpperCase()} ${node.path}`;
      routeMap[routeKey] = node.getData();
    }
  } catch (err) {
    logError(`Failed to build route map: ${err.message}`);
    process.exit(1);
  }

  const server = http.createServer(async (req, res) => {
    try {
      const { url, method } = req;
      const routeKey = `${method.toUpperCase()} ${url}`;
      const handler = routeMap[routeKey] || routeMap['GET /*'];

      if (!handler) {
        res.writeHead(404);
        logError(`Route not found for ${method} ${url}`);
        return res.end('Not found');
      }

      if (handler.status) {
        const status = parseInt(handler.status);
        if (!isNaN(status)) {
          res.statusCode = status;
        } else {
          logWarning(`Invalid status code in handler for ${routeKey}`);
        }
      }

      if (handler.headers && typeof handler.headers === 'object') {
        for (const [key, value] of Object.entries(handler.headers)) {
          res.setHeader(key, value);
        }
      }

      if (handler.type) {
        res.setHeader('Content-Type', handler.type);
      }

      if (handler.json) {
        res.setHeader('Content-Type', 'application/json');
        logInfo(`Responding with JSON for ${method} ${url}`);
        return res.end(JSON.stringify(handler.json));
      }

      if (handler.file) {
        const rawFilePath = handler.file;
        const filePath = rawFilePath.replace(/^['"]|['"]$/g, ''); 
      
        const contentType = handler.type || mime.lookup(filePath) || 'text/plain';
        res.setHeader('Content-Type', contentType);
        logInfo(`Serving file for ${method} ${url}: ${filePath}`);
      
        const stream = fs.createReadStream(filePath);
      
        stream.on('error', (err) => {
          logError(`Failed to serve file ${filePath}: ${err.message}`);
          if (!res.headersSent) {
            res.writeHead(500, { 'Content-Type': 'text/plain' });
            res.end(`Error reading file: ${err.message}`);
          }
        });
      
        return stream.pipe(res);
      }
      
      if (handler.html) {
        res.setHeader('Content-Type', 'text/html');
        logInfo(`Responding with HTML for ${method} ${url}`);
        return res.end(handler.html);
      }

      if (handler.response) {
        logInfo(`Responding with custom response for ${method} ${url}`);
        return res.end(handler.response);
      }

      logWarning(`No response specified for route ${routeKey}`);
      res.end();

    } catch (err) {
      logError(`Unhandled error during request processing: ${err.message}`);
      if (!res.headersSent) {
        res.writeHead(500, { 'Content-Type': 'text/plain' });
      }
      res.end('Internal Server Error');
    }
  });

  server.listen(port, () => {
    logSuccess(`Tera server running on http://localhost:${port}`);
  });

  server.on('error', (err) => {
    logError(`Server encountered an error: ${err.message}`);
  });
}
