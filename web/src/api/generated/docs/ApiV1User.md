
# ApiV1User

User is the standard user message for requests (includes password).

## Properties

Name | Type
------------ | -------------
`ids` | [ApiV1UserIDs](ApiV1UserIDs.md)
`email` | string
`password` | string

## Example

```typescript
import type { ApiV1User } from ''

// TODO: Update the object below with actual values
const example = {
  "ids": null,
  "email": null,
  "password": null,
} satisfies ApiV1User

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ApiV1User
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


