
# ApiV1UserResponse

UserResponse is the user message for responses (excludes password).

## Properties

Name | Type
------------ | -------------
`ids` | [ApiV1UserIDs](ApiV1UserIDs.md)
`email` | string

## Example

```typescript
import type { ApiV1UserResponse } from ''

// TODO: Update the object below with actual values
const example = {
  "ids": null,
  "email": null,
} satisfies ApiV1UserResponse

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ApiV1UserResponse
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


