
# ApiV1ParetoServiceListByUserResponse


## Properties

Name | Type
------------ | -------------
`pareto` | [ApiV1Pareto](ApiV1Pareto.md)
`totalCount` | number
`limit` | number
`offset` | number
`hasMore` | boolean

## Example

```typescript
import type { ApiV1ParetoServiceListByUserResponse } from ''

// TODO: Update the object below with actual values
const example = {
  "pareto": null,
  "totalCount": null,
  "limit": null,
  "offset": null,
  "hasMore": null,
} satisfies ApiV1ParetoServiceListByUserResponse

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ApiV1ParetoServiceListByUserResponse
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


