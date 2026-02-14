import yaml from "js-yaml";

// Types derived from OpenAPI spec
export interface OpenAPISpec {
  openapi: string;
  info: {
    title: string;
    description: string;
    version: string;
  };
  servers: Array<{
    url: string;
    description: string;
  }>;
  paths: Record<string, PathItem>;
  components: {
    schemas: Record<string, SchemaObject>;
    securitySchemes?: Record<string, SecurityScheme>;
  };
  tags: Array<{
    name: string;
    description: string;
  }>;
}

export interface PathItem {
  get?: Operation;
  post?: Operation;
  put?: Operation;
  patch?: Operation;
  delete?: Operation;
}

export interface Operation {
  operationId: string;
  summary: string;
  description: string;
  tags: string[];
  parameters?: Parameter[];
  requestBody?: RequestBody;
  responses: Record<string, Response>;
}

export interface Parameter {
  name: string;
  in: "query" | "path" | "header" | "cookie";
  description?: string;
  required?: boolean;
  schema: SchemaObject;
  example?: unknown;
}

export interface RequestBody {
  required?: boolean;
  content: {
    "application/json": {
      schema: SchemaObject | RefObject;
      example?: unknown;
    };
  };
}

export interface Response {
  description: string;
  content?: {
    "application/json": {
      schema: SchemaObject | RefObject;
      example?: unknown;
    };
  };
}

export interface SchemaObject {
  type?: string;
  format?: string;
  description?: string;
  properties?: Record<string, SchemaObject | RefObject>;
  items?: SchemaObject | RefObject;
  required?: string[];
  enum?: string[];
  example?: unknown;
  nullable?: boolean;
  $ref?: string;
}

export interface RefObject {
  $ref: string;
}

export interface SecurityScheme {
  type: string;
  in?: string;
  name?: string;
  description?: string;
}

// Transformed types for the UI
export interface ApiCategory {
  id: string;
  title: string;
  description?: string;
  endpoints: ApiEndpoint[];
}

export interface ApiEndpoint {
  id: string;
  operationId: string;
  title: string;
  method: string;
  path: string;
  description: string;
  authentication: boolean;
  parameters: ApiParameter[];
  requestBody?: {
    example: string;
    schema?: SchemaObject | RefObject;
  };
  response?: {
    example: string;
    schema?: SchemaObject | RefObject;
  };
}

export interface ApiParameter {
  name: string;
  type: string;
  required: boolean;
  description: string;
  location: "query" | "path" | "body";
  example?: unknown;
}

function isRefObject(schema: SchemaObject | RefObject): schema is RefObject {
  return "$ref" in schema && typeof (schema as RefObject).$ref === "string";
}

function resolveRef(spec: OpenAPISpec, ref: string): SchemaObject {
  const path = ref.replace("#/components/schemas/", "");
  return spec.components.schemas[path] || {};
}

function resolveSchema(
  spec: OpenAPISpec,
  schema: SchemaObject | RefObject,
): SchemaObject {
  if (isRefObject(schema)) {
    return resolveRef(spec, schema.$ref);
  }
  return schema;
}

function getSchemaType(
  schema: SchemaObject | RefObject,
  spec: OpenAPISpec,
): string {
  const resolved = resolveSchema(spec, schema);
  if (resolved.type === "array" && resolved.items) {
    const itemType = getSchemaType(resolved.items, spec);
    return `${itemType}[]`;
  }
  return resolved.type || "object";
}

function generateSlug(text: string): string {
  return text
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/(^-|-$)/g, "");
}

export function parseOpenAPISpec(spec: OpenAPISpec): ApiCategory[] {
  const categoriesMap = new Map<string, ApiCategory>();

  // Initialize categories from tags
  for (const tag of spec.tags || []) {
    categoriesMap.set(tag.name, {
      id: generateSlug(tag.name),
      title: tag.name,
      description: tag.description,
      endpoints: [],
    });
  }

  // Parse paths into endpoints
  for (const [path, pathItem] of Object.entries(spec.paths)) {
    const methods: Array<{ method: string; operation: Operation }> = [];

    if (pathItem.get) methods.push({ method: "GET", operation: pathItem.get });
    if (pathItem.post)
      methods.push({ method: "POST", operation: pathItem.post });
    if (pathItem.put) methods.push({ method: "PUT", operation: pathItem.put });
    if (pathItem.patch)
      methods.push({ method: "PATCH", operation: pathItem.patch });
    if (pathItem.delete)
      methods.push({ method: "DELETE", operation: pathItem.delete });

    for (const { method, operation } of methods) {
      const tagName = operation.tags?.[0] || "Other";

      // Ensure category exists
      if (!categoriesMap.has(tagName)) {
        categoriesMap.set(tagName, {
          id: generateSlug(tagName),
          title: tagName,
          endpoints: [],
        });
      }

      const category = categoriesMap.get(tagName)!;

      // Parse parameters
      const parameters: ApiParameter[] = [];

      // Path and query parameters
      for (const param of operation.parameters || []) {
        parameters.push({
          name: param.name,
          type: getSchemaType(param.schema, spec),
          required: param.required || false,
          description: param.description || "",
          location: param.in as "query" | "path",
          example: param.example,
        });
      }

      // Request body parameters
      if (operation.requestBody?.content?.["application/json"]) {
        const bodyContent = operation.requestBody.content["application/json"];
        const bodySchema = resolveSchema(spec, bodyContent.schema);

        if (bodySchema.properties) {
          for (const [propName, propSchemaRef] of Object.entries(
            bodySchema.properties,
          )) {
            const propSchema = resolveSchema(spec, propSchemaRef);
            parameters.push({
              name: propName,
              type: getSchemaType(propSchemaRef, spec),
              required: bodySchema.required?.includes(propName) || false,
              description: propSchema.description || "",
              location: "body",
              example: propSchema.example,
            });
          }
        }
      }

      // Get request example
      let requestExample: string | undefined;
      if (operation.requestBody?.content?.["application/json"]?.example) {
        requestExample = JSON.stringify(
          operation.requestBody.content["application/json"].example,
          null,
          2,
        );
      }

      // Get response example (from 200 response)
      let responseExample: string | undefined;
      const successResponse = operation.responses["200"];
      if (successResponse?.content?.["application/json"]?.example) {
        responseExample = JSON.stringify(
          successResponse.content["application/json"].example,
          null,
          2,
        );
      }

      const endpoint: ApiEndpoint = {
        id: generateSlug(operation.summary),
        operationId: operation.operationId,
        title: operation.summary,
        method,
        path,
        description: operation.description || "",
        authentication: true, // All endpoints require auth based on security scheme
        parameters,
        requestBody: requestExample
          ? {
              example: requestExample,
              schema:
                operation.requestBody?.content?.["application/json"]?.schema,
            }
          : undefined,
        response: responseExample
          ? {
              example: responseExample,
              schema: successResponse?.content?.["application/json"]?.schema,
            }
          : undefined,
      };

      category.endpoints.push(endpoint);
    }
  }

  // Convert map to array and filter out empty categories
  return Array.from(categoriesMap.values()).filter(
    (cat) => cat.endpoints.length > 0,
  );
}

// Cache the spec to avoid refetching
let cachedSpec: OpenAPISpec | null = null;
let cachedCategories: ApiCategory[] | null = null;

export async function fetchOpenAPISpec(): Promise<OpenAPISpec> {
  if (cachedSpec) {
    return cachedSpec;
  }

  // Try to fetch from the API server first, fall back to local
  const urls = ["/openapi.yaml"];

  for (const url of urls) {
    try {
      const response = await fetch(url);
      if (response.ok) {
        const yamlText = await response.text();
        // Parse YAML - we'll use a simple approach since we control the format
        cachedSpec = parseYAML(yamlText) as OpenAPISpec;
        return cachedSpec;
      }
    } catch {
      // Continue to next URL
    }
  }

  throw new Error("Failed to fetch OpenAPI spec");
}

export async function getApiEndpoints(): Promise<ApiCategory[]> {
  if (cachedCategories) {
    return cachedCategories;
  }

  const spec = await fetchOpenAPISpec();
  cachedCategories = parseOpenAPISpec(spec);
  return cachedCategories;
}

// Simple YAML parser for our specific OpenAPI format
function parseYAML(yamlText: string): unknown {
  return yaml.load(yamlText);
}
