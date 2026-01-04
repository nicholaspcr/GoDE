# ApiV1AuthServiceApi

All URIs are relative to *http://localhost*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**authServiceLogin**](ApiV1AuthServiceApi.md#authservicelogin) | **POST** /v1/auth/login |  |
| [**authServiceLogout**](ApiV1AuthServiceApi.md#authservicelogout) | **POST** /v1/auth/logout |  |
| [**authServiceRefreshToken**](ApiV1AuthServiceApi.md#authservicerefreshtoken) | **POST** /v1/auth/refresh |  |
| [**authServiceRegister**](ApiV1AuthServiceApi.md#authserviceregister) | **POST** /v1/auth/register |  |



## authServiceLogin

> ApiV1AuthServiceLoginResponse authServiceLogin(body)



### Example

```ts
import {
  Configuration,
  ApiV1AuthServiceApi,
} from '';
import type { AuthServiceLoginRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new ApiV1AuthServiceApi();

  const body = {
    // ApiV1AuthServiceLoginRequest
    body: ...,
  } satisfies AuthServiceLoginRequest;

  try {
    const data = await api.authServiceLogin(body);
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
| **body** | [ApiV1AuthServiceLoginRequest](ApiV1AuthServiceLoginRequest.md) |  | |

### Return type

[**ApiV1AuthServiceLoginResponse**](ApiV1AuthServiceLoginResponse.md)

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


## authServiceLogout

> object authServiceLogout(body)



### Example

```ts
import {
  Configuration,
  ApiV1AuthServiceApi,
} from '';
import type { AuthServiceLogoutRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new ApiV1AuthServiceApi();

  const body = {
    // ApiV1AuthServiceLogoutRequest
    body: ...,
  } satisfies AuthServiceLogoutRequest;

  try {
    const data = await api.authServiceLogout(body);
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
| **body** | [ApiV1AuthServiceLogoutRequest](ApiV1AuthServiceLogoutRequest.md) |  | |

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


## authServiceRefreshToken

> ApiV1AuthServiceRefreshTokenResponse authServiceRefreshToken(body)



### Example

```ts
import {
  Configuration,
  ApiV1AuthServiceApi,
} from '';
import type { AuthServiceRefreshTokenRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new ApiV1AuthServiceApi();

  const body = {
    // ApiV1AuthServiceRefreshTokenRequest
    body: ...,
  } satisfies AuthServiceRefreshTokenRequest;

  try {
    const data = await api.authServiceRefreshToken(body);
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
| **body** | [ApiV1AuthServiceRefreshTokenRequest](ApiV1AuthServiceRefreshTokenRequest.md) |  | |

### Return type

[**ApiV1AuthServiceRefreshTokenResponse**](ApiV1AuthServiceRefreshTokenResponse.md)

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


## authServiceRegister

> object authServiceRegister(body)



### Example

```ts
import {
  Configuration,
  ApiV1AuthServiceApi,
} from '';
import type { AuthServiceRegisterRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new ApiV1AuthServiceApi();

  const body = {
    // ApiV1AuthServiceRegisterRequest
    body: ...,
  } satisfies AuthServiceRegisterRequest;

  try {
    const data = await api.authServiceRegister(body);
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
| **body** | [ApiV1AuthServiceRegisterRequest](ApiV1AuthServiceRegisterRequest.md) |  | |

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

