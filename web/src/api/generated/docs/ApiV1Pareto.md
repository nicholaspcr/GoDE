
# ApiV1Pareto


## Properties

Name | Type
------------ | -------------
`ids` | [ApiV1ParetoIDs](ApiV1ParetoIDs.md)
`vectors` | [Array&lt;ApiV1Vector&gt;](ApiV1Vector.md)
`maxObjs` | Array&lt;number&gt;

## Example

```typescript
import type { ApiV1Pareto } from ''

// TODO: Update the object below with actual values
const example = {
  "ids": null,
  "vectors": null,
  "maxObjs": null,
} satisfies ApiV1Pareto

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ApiV1Pareto
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


