# ApiV1ParetoServiceApi

All URIs are relative to *http://localhost*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**paretoServiceDelete**](ApiV1ParetoServiceApi.md#paretoservicedelete) | **DELETE** /v1/pareto/{paretoIds.id} |  |
| [**paretoServiceGet**](ApiV1ParetoServiceApi.md#paretoserviceget) | **GET** /v1/pareto/{paretoIds.id} |  |
| [**paretoServiceListByUser**](ApiV1ParetoServiceApi.md#paretoservicelistbyuser) | **GET** /v1/paretos/{userIds.username} |  |



## paretoServiceDelete

> object paretoServiceDelete(paretoIdsId, paretoIdsUserId)



### Example

```ts
import {
  Configuration,
  ApiV1ParetoServiceApi,
} from '';
import type { ParetoServiceDeleteRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new ApiV1ParetoServiceApi();

  const body = {
    // string
    paretoIdsId: paretoIdsId_example,
    // string (optional)
    paretoIdsUserId: paretoIdsUserId_example,
  } satisfies ParetoServiceDeleteRequest;

  try {
    const data = await api.paretoServiceDelete(body);
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
| **paretoIdsId** | `string` |  | [Defaults to `undefined`] |
| **paretoIdsUserId** | `string` |  | [Optional] [Defaults to `undefined`] |

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


## paretoServiceGet

> ApiV1ParetoServiceGetResponse paretoServiceGet(paretoIdsId, paretoIdsUserId)



### Example

```ts
import {
  Configuration,
  ApiV1ParetoServiceApi,
} from '';
import type { ParetoServiceGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new ApiV1ParetoServiceApi();

  const body = {
    // string
    paretoIdsId: paretoIdsId_example,
    // string (optional)
    paretoIdsUserId: paretoIdsUserId_example,
  } satisfies ParetoServiceGetRequest;

  try {
    const data = await api.paretoServiceGet(body);
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
| **paretoIdsId** | `string` |  | [Defaults to `undefined`] |
| **paretoIdsUserId** | `string` |  | [Optional] [Defaults to `undefined`] |

### Return type

[**ApiV1ParetoServiceGetResponse**](ApiV1ParetoServiceGetResponse.md)

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


## paretoServiceListByUser

> StreamResultOfApiV1ParetoServiceListByUserResponse paretoServiceListByUser(userIdsUsername, limit, offset)



### Example

```ts
import {
  Configuration,
  ApiV1ParetoServiceApi,
} from '';
import type { ParetoServiceListByUserRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new ApiV1ParetoServiceApi();

  const body = {
    // string
    userIdsUsername: userIdsUsername_example,
    // number | Page size (default: 50, max: 100) (optional)
    limit: 56,
    // number | Starting position (default: 0) (optional)
    offset: 56,
  } satisfies ParetoServiceListByUserRequest;

  try {
    const data = await api.paretoServiceListByUser(body);
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
| **userIdsUsername** | `string` |  | [Defaults to `undefined`] |
| **limit** | `number` | Page size (default: 50, max: 100) | [Optional] [Defaults to `undefined`] |
| **offset** | `number` | Starting position (default: 0) | [Optional] [Defaults to `undefined`] |

### Return type

[**StreamResultOfApiV1ParetoServiceListByUserResponse**](StreamResultOfApiV1ParetoServiceListByUserResponse.md)

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

