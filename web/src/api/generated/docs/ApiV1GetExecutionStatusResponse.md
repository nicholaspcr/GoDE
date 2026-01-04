
# ApiV1GetExecutionStatusResponse


## Properties

Name | Type
------------ | -------------
`execution` | [ApiV1Execution](ApiV1Execution.md)
`progress` | [ApiV1StreamProgressResponse](ApiV1StreamProgressResponse.md)

## Example

```typescript
import type { ApiV1GetExecutionStatusResponse } from ''

// TODO: Update the object below with actual values
const example = {
  "execution": null,
  "progress": null,
} satisfies ApiV1GetExecutionStatusResponse

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ApiV1GetExecutionStatusResponse
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


