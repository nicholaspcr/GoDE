
# ApiV1ListExecutionsResponse


## Properties

Name | Type
------------ | -------------
`executions` | [Array&lt;ApiV1Execution&gt;](ApiV1Execution.md)
`totalCount` | number
`limit` | number
`offset` | number
`hasMore` | boolean

## Example

```typescript
import type { ApiV1ListExecutionsResponse } from ''

// TODO: Update the object below with actual values
const example = {
  "executions": null,
  "totalCount": null,
  "limit": null,
  "offset": null,
  "hasMore": null,
} satisfies ApiV1ListExecutionsResponse

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ApiV1ListExecutionsResponse
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


