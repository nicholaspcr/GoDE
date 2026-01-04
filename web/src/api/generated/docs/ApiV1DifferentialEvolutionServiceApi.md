# ApiV1DifferentialEvolutionServiceApi

All URIs are relative to *http://localhost*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**differentialEvolutionServiceCancelExecution**](ApiV1DifferentialEvolutionServiceApi.md#differentialevolutionservicecancelexecution) | **POST** /v1/de/executions/{executionId}/cancel |  |
| [**differentialEvolutionServiceDeleteExecution**](ApiV1DifferentialEvolutionServiceApi.md#differentialevolutionservicedeleteexecution) | **DELETE** /v1/de/executions/{executionId} |  |
| [**differentialEvolutionServiceGetExecutionResults**](ApiV1DifferentialEvolutionServiceApi.md#differentialevolutionservicegetexecutionresults) | **GET** /v1/de/executions/{executionId}/results |  |
| [**differentialEvolutionServiceGetExecutionStatus**](ApiV1DifferentialEvolutionServiceApi.md#differentialevolutionservicegetexecutionstatus) | **GET** /v1/de/executions/{executionId} |  |
| [**differentialEvolutionServiceListExecutions**](ApiV1DifferentialEvolutionServiceApi.md#differentialevolutionservicelistexecutions) | **GET** /v1/de/executions |  |
| [**differentialEvolutionServiceListSupportedAlgorithms**](ApiV1DifferentialEvolutionServiceApi.md#differentialevolutionservicelistsupportedalgorithms) | **GET** /v1/de/supported/algorithms |  |
| [**differentialEvolutionServiceListSupportedProblems**](ApiV1DifferentialEvolutionServiceApi.md#differentialevolutionservicelistsupportedproblems) | **GET** /v1/de/supported/problems |  |
| [**differentialEvolutionServiceListSupportedVariants**](ApiV1DifferentialEvolutionServiceApi.md#differentialevolutionservicelistsupportedvariants) | **GET** /v1/de/supported/variants |  |
| [**differentialEvolutionServiceRunAsync**](ApiV1DifferentialEvolutionServiceApi.md#differentialevolutionservicerunasync) | **POST** /v1/de/run | Async execution RPCs |
| [**differentialEvolutionServiceStreamProgress**](ApiV1DifferentialEvolutionServiceApi.md#differentialevolutionservicestreamprogress) | **GET** /v1/de/executions/{executionId}/progress |  |



## differentialEvolutionServiceCancelExecution

> object differentialEvolutionServiceCancelExecution(executionId, body)



### Example

```ts
import {
  Configuration,
  ApiV1DifferentialEvolutionServiceApi,
} from '';
import type { DifferentialEvolutionServiceCancelExecutionRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new ApiV1DifferentialEvolutionServiceApi();

  const body = {
    // string
    executionId: executionId_example,
    // object
    body: Object,
  } satisfies DifferentialEvolutionServiceCancelExecutionRequest;

  try {
    const data = await api.differentialEvolutionServiceCancelExecution(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **executionId** | `string` |  | [Defaults to `undefined`] |
| **body** | `object` |  | |

### Return type

**object**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | A successful response. |  -  |
| **0** | An unexpected error response. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## differentialEvolutionServiceDeleteExecution

> object differentialEvolutionServiceDeleteExecution(executionId)



### Example

```ts
import {
  Configuration,
  ApiV1DifferentialEvolutionServiceApi,
} from '';
import type { DifferentialEvolutionServiceDeleteExecutionRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new ApiV1DifferentialEvolutionServiceApi();

  const body = {
    // string
    executionId: executionId_example,
  } satisfies DifferentialEvolutionServiceDeleteExecutionRequest;

  try {
    const data = await api.differentialEvolutionServiceDeleteExecution(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **executionId** | `string` |  | [Defaults to `undefined`] |

### Return type

**object**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | A successful response. |  -  |
| **0** | An unexpected error response. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## differentialEvolutionServiceGetExecutionResults

> ApiV1GetExecutionResultsResponse differentialEvolutionServiceGetExecutionResults(executionId)



### Example

```ts
import {
  Configuration,
  ApiV1DifferentialEvolutionServiceApi,
} from '';
import type { DifferentialEvolutionServiceGetExecutionResultsRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new ApiV1DifferentialEvolutionServiceApi();

  const body = {
    // string
    executionId: executionId_example,
  } satisfies DifferentialEvolutionServiceGetExecutionResultsRequest;

  try {
    const data = await api.differentialEvolutionServiceGetExecutionResults(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **executionId** | `string` |  | [Defaults to `undefined`] |

### Return type

[**ApiV1GetExecutionResultsResponse**](ApiV1GetExecutionResultsResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | A successful response. |  -  |
| **0** | An unexpected error response. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## differentialEvolutionServiceGetExecutionStatus

> ApiV1GetExecutionStatusResponse differentialEvolutionServiceGetExecutionStatus(executionId)



### Example

```ts
import {
  Configuration,
  ApiV1DifferentialEvolutionServiceApi,
} from '';
import type { DifferentialEvolutionServiceGetExecutionStatusRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new ApiV1DifferentialEvolutionServiceApi();

  const body = {
    // string
    executionId: executionId_example,
  } satisfies DifferentialEvolutionServiceGetExecutionStatusRequest;

  try {
    const data = await api.differentialEvolutionServiceGetExecutionStatus(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **executionId** | `string` |  | [Defaults to `undefined`] |

### Return type

[**ApiV1GetExecutionStatusResponse**](ApiV1GetExecutionStatusResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | A successful response. |  -  |
| **0** | An unexpected error response. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## differentialEvolutionServiceListExecutions

> ApiV1ListExecutionsResponse differentialEvolutionServiceListExecutions(status, limit, offset)



### Example

```ts
import {
  Configuration,
  ApiV1DifferentialEvolutionServiceApi,
} from '';
import type { DifferentialEvolutionServiceListExecutionsRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new ApiV1DifferentialEvolutionServiceApi();

  const body = {
    // 'EXECUTION_STATUS_UNSPECIFIED' | 'EXECUTION_STATUS_PENDING' | 'EXECUTION_STATUS_RUNNING' | 'EXECUTION_STATUS_COMPLETED' | 'EXECUTION_STATUS_FAILED' | 'EXECUTION_STATUS_CANCELLED' | Optional filter (optional)
    status: status_example,
    // number | Page size (default: 50, max: 100) (optional)
    limit: 56,
    // number | Starting position (default: 0) (optional)
    offset: 56,
  } satisfies DifferentialEvolutionServiceListExecutionsRequest;

  try {
    const data = await api.differentialEvolutionServiceListExecutions(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **status** | `EXECUTION_STATUS_UNSPECIFIED`, `EXECUTION_STATUS_PENDING`, `EXECUTION_STATUS_RUNNING`, `EXECUTION_STATUS_COMPLETED`, `EXECUTION_STATUS_FAILED`, `EXECUTION_STATUS_CANCELLED` | Optional filter | [Optional] [Defaults to `&#39;EXECUTION_STATUS_UNSPECIFIED&#39;`] [Enum: EXECUTION_STATUS_UNSPECIFIED, EXECUTION_STATUS_PENDING, EXECUTION_STATUS_RUNNING, EXECUTION_STATUS_COMPLETED, EXECUTION_STATUS_FAILED, EXECUTION_STATUS_CANCELLED] |
| **limit** | `number` | Page size (default: 50, max: 100) | [Optional] [Defaults to `undefined`] |
| **offset** | `number` | Starting position (default: 0) | [Optional] [Defaults to `undefined`] |

### Return type

[**ApiV1ListExecutionsResponse**](ApiV1ListExecutionsResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | A successful response. |  -  |
| **0** | An unexpected error response. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## differentialEvolutionServiceListSupportedAlgorithms

> ApiV1ListSupportedAlgorithmsResponse differentialEvolutionServiceListSupportedAlgorithms()



### Example

```ts
import {
  Configuration,
  ApiV1DifferentialEvolutionServiceApi,
} from '';
import type { DifferentialEvolutionServiceListSupportedAlgorithmsRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new ApiV1DifferentialEvolutionServiceApi();

  try {
    const data = await api.differentialEvolutionServiceListSupportedAlgorithms();
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

This endpoint does not need any parameter.

### Return type

[**ApiV1ListSupportedAlgorithmsResponse**](ApiV1ListSupportedAlgorithmsResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | A successful response. |  -  |
| **0** | An unexpected error response. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## differentialEvolutionServiceListSupportedProblems

> ApiV1ListSupportedProblemsResponse differentialEvolutionServiceListSupportedProblems()



### Example

```ts
import {
  Configuration,
  ApiV1DifferentialEvolutionServiceApi,
} from '';
import type { DifferentialEvolutionServiceListSupportedProblemsRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new ApiV1DifferentialEvolutionServiceApi();

  try {
    const data = await api.differentialEvolutionServiceListSupportedProblems();
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

This endpoint does not need any parameter.

### Return type

[**ApiV1ListSupportedProblemsResponse**](ApiV1ListSupportedProblemsResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | A successful response. |  -  |
| **0** | An unexpected error response. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## differentialEvolutionServiceListSupportedVariants

> ApiV1ListSupportedVariantsResponse differentialEvolutionServiceListSupportedVariants()



### Example

```ts
import {
  Configuration,
  ApiV1DifferentialEvolutionServiceApi,
} from '';
import type { DifferentialEvolutionServiceListSupportedVariantsRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new ApiV1DifferentialEvolutionServiceApi();

  try {
    const data = await api.differentialEvolutionServiceListSupportedVariants();
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

This endpoint does not need any parameter.

### Return type

[**ApiV1ListSupportedVariantsResponse**](ApiV1ListSupportedVariantsResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | A successful response. |  -  |
| **0** | An unexpected error response. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## differentialEvolutionServiceRunAsync

> ApiV1RunAsyncResponse differentialEvolutionServiceRunAsync(body)

Async execution RPCs

### Example

```ts
import {
  Configuration,
  ApiV1DifferentialEvolutionServiceApi,
} from '';
import type { DifferentialEvolutionServiceRunAsyncRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new ApiV1DifferentialEvolutionServiceApi();

  const body = {
    // ApiV1RunAsyncRequest
    body: ...,
  } satisfies DifferentialEvolutionServiceRunAsyncRequest;

  try {
    const data = await api.differentialEvolutionServiceRunAsync(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **body** | [ApiV1RunAsyncRequest](ApiV1RunAsyncRequest.md) |  | |

### Return type

[**ApiV1RunAsyncResponse**](ApiV1RunAsyncResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | A successful response. |  -  |
| **0** | An unexpected error response. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## differentialEvolutionServiceStreamProgress

> StreamResultOfApiV1StreamProgressResponse differentialEvolutionServiceStreamProgress(executionId)



### Example

```ts
import {
  Configuration,
  ApiV1DifferentialEvolutionServiceApi,
} from '';
import type { DifferentialEvolutionServiceStreamProgressRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new ApiV1DifferentialEvolutionServiceApi();

  const body = {
    // string
    executionId: executionId_example,
  } satisfies DifferentialEvolutionServiceStreamProgressRequest;

  try {
    const data = await api.differentialEvolutionServiceStreamProgress(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **executionId** | `string` |  | [Defaults to `undefined`] |

### Return type

[**StreamResultOfApiV1StreamProgressResponse**](StreamResultOfApiV1StreamProgressResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | A successful response.(streaming responses) |  -  |
| **0** | An unexpected error response. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

