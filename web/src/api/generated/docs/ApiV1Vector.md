
# ApiV1Vector

Vector is an element used in the Differential Evoluition algorithm.

## Properties

Name | Type
------------ | -------------
`ids` | [ApiV1VectorIDs](ApiV1VectorIDs.md)
`elements` | Array&lt;number&gt;
`objectives` | Array&lt;number&gt;
`crowdingDistance` | number

## Example

```typescript
import type { ApiV1Vector } from ''

// TODO: Update the object below with actual values
const example = {
  "ids": null,
  "elements": null,
  "objectives": null,
  "crowdingDistance": null,
} satisfies ApiV1Vector

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ApiV1Vector
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


