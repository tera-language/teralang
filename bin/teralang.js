#!/usr/bin/env node
import { serveTeraFile } from '../lib/server.js';

const args = process.argv.slice(2);
const command = args[0];
const file = args[1];
const portFlagIndex = args.indexOf('--port');
const port = portFlagIndex !== -1 ? parseInt(args[portFlagIndex + 1], 10) : 3000;

if (command === 'serve' && file) {
  serveTeraFile(file, port);
} else {
  console.error('Usage: npx teralang serve <file.tera> [--port <number>]');
  process.exit(1);
}
