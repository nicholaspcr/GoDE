
# ApiV1Execution


## Properties

Name | Type
------------ | -------------
`id` | string
`userId` | string
`status` | [ApiV1ExecutionStatus](ApiV1ExecutionStatus.md)
`config` | [ApiV1DEConfig](ApiV1DEConfig.md)
`createdAt` | Date
`updatedAt` | Date
`completedAt` | Date
`error` | string
`paretoId` | string
`algorithm` | string
`variant` | string
`problem` | string

## Example

```typescript
import type { ApiV1Execution } from ''

// TODO: Update the object below with actual values
const example = {
  "id": null,
  "userId": null,
  "status": null,
  "config": null,
  "createdAt": null,
  "updatedAt": null,
  "completedAt": null,
  "error": null,
  "paretoId": null,
  "algorithm": null,
  "variant": null,
  "problem": null,
} satisfies ApiV1Execution

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ApiV1Execution
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


