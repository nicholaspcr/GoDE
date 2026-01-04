
# ApiV1DEConfig


## Properties

Name | Type
------------ | -------------
`executions` | string
`generations` | string
`populationSize` | string
`dimensionsSize` | string
`objectivesSize` | string
`floorLimiter` | number
`ceilLimiter` | number
`gde3` | [ApiV1GDE3Config](ApiV1GDE3Config.md)

## Example

```typescript
import type { ApiV1DEConfig } from ''

// TODO: Update the object below with actual values
const example = {
  "executions": null,
  "generations": null,
  "populationSize": null,
  "dimensionsSize": null,
  "objectivesSize": null,
  "floorLimiter": null,
  "ceilLimiter": null,
  "gde3": null,
} satisfies ApiV1DEConfig

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ApiV1DEConfig
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


