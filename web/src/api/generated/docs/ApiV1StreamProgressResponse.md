
# ApiV1StreamProgressResponse


## Properties

Name | Type
------------ | -------------
`executionId` | string
`currentGeneration` | number
`totalGenerations` | number
`completedExecutions` | number
`totalExecutions` | number
`partialPareto` | [Array&lt;ApiV1Vector&gt;](ApiV1Vector.md)
`updatedAt` | Date

## Example

```typescript
import type { ApiV1StreamProgressResponse } from ''

// TODO: Update the object below with actual values
const example = {
  "executionId": null,
  "currentGeneration": null,
  "totalGenerations": null,
  "completedExecutions": null,
  "totalExecutions": null,
  "partialPareto": null,
  "updatedAt": null,
} satisfies ApiV1StreamProgressResponse

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ApiV1StreamProgressResponse
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


