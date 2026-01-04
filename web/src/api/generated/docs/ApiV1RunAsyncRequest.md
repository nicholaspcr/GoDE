
# ApiV1RunAsyncRequest


## Properties

Name | Type
------------ | -------------
`algorithm` | string
`variant` | string
`problem` | string
`deConfig` | [ApiV1DEConfig](ApiV1DEConfig.md)

## Example

```typescript
import type { ApiV1RunAsyncRequest } from ''

// TODO: Update the object below with actual values
const example = {
  "algorithm": null,
  "variant": null,
  "problem": null,
  "deConfig": null,
} satisfies ApiV1RunAsyncRequest

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ApiV1RunAsyncRequest
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


