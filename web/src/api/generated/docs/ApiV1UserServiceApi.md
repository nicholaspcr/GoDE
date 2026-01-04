# ApiV1UserServiceApi

All URIs are relative to *http://localhost*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**userServiceCreate**](ApiV1UserServiceApi.md#userservicecreate) | **POST** /v1/user |  |
| [**userServiceDelete**](ApiV1UserServiceApi.md#userservicedelete) | **DELETE** /v1/user/{userIds.username} |  |
| [**userServiceGet**](ApiV1UserServiceApi.md#userserviceget) | **GET** /v1/user/{userIds.username} |  |
| [**userServiceUpdate**](ApiV1UserServiceApi.md#userserviceupdate) | **PUT** /v1/user |  |



## userServiceCreate

> object userServiceCreate(body)



### Example

```ts
import {
  Configuration,
  ApiV1UserServiceApi,
} from '';
import type { UserServiceCreateRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new ApiV1UserServiceApi();

  const body = {
    // ApiV1UserServiceCreateRequest
    body: ...,
  } satisfies UserServiceCreateRequest;

  try {
    const data = await api.userServiceCreate(body);
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
| **body** | [ApiV1UserServiceCreateRequest](ApiV1UserServiceCreateRequest.md) |  | |

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


## userServiceDelete

> object userServiceDelete(userIdsUsername)



### Example

```ts
import {
  Configuration,
  ApiV1UserServiceApi,
} from '';
import type { UserServiceDeleteRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new ApiV1UserServiceApi();

  const body = {
    // string
    userIdsUsername: userIdsUsername_example,
  } satisfies UserServiceDeleteRequest;

  try {
    const data = await api.userServiceDelete(body);
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


## userServiceGet

> ApiV1UserServiceGetResponse userServiceGet(userIdsUsername)



### Example

```ts
import {
  Configuration,
  ApiV1UserServiceApi,
} from '';
import type { UserServiceGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new ApiV1UserServiceApi();

  const body = {
    // string
    userIdsUsername: userIdsUsername_example,
  } satisfies UserServiceGetRequest;

  try {
    const data = await api.userServiceGet(body);
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

### Return type

[**ApiV1UserServiceGetResponse**](ApiV1UserServiceGetResponse.md)

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


## userServiceUpdate

> object userServiceUpdate(body)



### Example

```ts
import {
  Configuration,
  ApiV1UserServiceApi,
} from '';
import type { UserServiceUpdateRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new ApiV1UserServiceApi();

  const body = {
    // ApiV1UserServiceUpdateRequest
    body: ...,
  } satisfies UserServiceUpdateRequest;

  try {
    const data = await api.userServiceUpdate(body);
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
| **body** | [ApiV1UserServiceUpdateRequest](ApiV1UserServiceUpdateRequest.md) |  | |

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

