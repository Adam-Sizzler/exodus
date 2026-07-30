const fs = require('fs');
const path = require('path');

const goFiles = {
  nodes: {
    path: path.join(__dirname, '../backend/internal/httpapi/nodes/types.go'),
    schema: require('../frontend/vendor/@exodus/backend-contract/build/backend/models/nodes.schema.js').NodesSchema,
    structNames: ['nodeAPI']
  },
  hosts: {
    path: path.join(__dirname, '../backend/internal/httpapi/hosts/types.go'),
    schema: require('../frontend/vendor/@exodus/backend-contract/build/backend/models/hosts.schema.js').HostsSchema,
    structNames: ['HostAPI']
  },
  users: {
    path: path.join(__dirname, '../backend/internal/httpapi/users/types.go'),
    schema: require('../frontend/vendor/@exodus/backend-contract/build/backend/models/extended-users.schema.js').ExtendedUsersSchema,
    structNames: ['userAPI']
  }
};

function parseGoStructs(filePath) {
  const content = fs.readFileSync(filePath, 'utf8');
  const structs = {};
  
  // Find all structs (match non-indented closing brace at the start of a line to handle nested structs)
  const structRegex = /type\s+(\w+)\s+struct\s*\{([\s\S]*?)\n\}/g;
  let match;
  while ((match = structRegex.exec(content)) !== null) {
    const structName = match[1];
    const structBody = match[2];
    
    // Clean nested anonymous structs to avoid parsing their fields as flat tags
    const cleanedBody = structBody.replace(/struct\s*\{[\s\S]*?\}/g, 'struct_placeholder');
    
    const fields = [];
    // Capture both standard and inline json tags (e.g. `json:"fieldName"` or `json:"fieldName,omitempty"`)
    const fieldRegex = /(\w+)\s+[\w\*\[\]\.\{\}\w\s\:]+\s+`json:"([^",]+)(?:,[^"]*)?"`/g;
    let fieldMatch;
    while ((fieldMatch = fieldRegex.exec(cleanedBody)) !== null) {
      fields.push({
        goName: fieldMatch[1],
        jsonName: fieldMatch[2]
      });
    }
    structs[structName] = fields;
  }
  return structs;
}

let hasError = false;

for (const [domain, config] of Object.entries(goFiles)) {
  console.log(`Checking domain: ${domain}...`);
  if (!fs.existsSync(config.path)) {
    console.error(`Error: File ${config.path} not found`);
    hasError = true;
    continue;
  }
  
  const structs = parseGoStructs(config.path);
  const zodKeys = Object.keys(config.schema.shape);
  
  for (const structName of config.structNames) {
    const structFields = structs[structName];
    if (!structFields) {
      console.warn(`Warning: Struct ${structName} not found in Go file`);
      continue;
    }
    
    console.log(`  Comparing Go struct '${structName}' with Zod schema:`);
    for (const field of structFields) {
      const isPresent = zodKeys.includes(field.jsonName);
      if (!isPresent) {
        console.error(`    [MISSING] Field '${field.jsonName}' (Go: ${field.goName}) is not in Zod schema!`);
        hasError = true;
      } else {
        console.log(`    [OK] Field '${field.jsonName}' matches.`);
      }
    }
  }
}

if (hasError) {
  process.exit(1);
} else {
  console.log("All contracts synchronized successfully!");
}
