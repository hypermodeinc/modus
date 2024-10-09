// webpack.config.js
import path from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

export default {
  entry: './index.js', // Path to your entry file
  output: {
    filename: 'bundle.js', // Output bundle filename
    path: path.resolve(__dirname, 'dist'), // Output directory for the bundle
    clean: true, // Clean the output directory before each build
  },
  target: 'node', // Target Node.js environment
  mode: 'production', // Set mode to production for optimization
  resolve: {
    extensions: ['.js'], // Resolve JavaScript files
  },
  externals: {
    v8CompileCache: 'commonjs v8-compile-cache'
  }
};
