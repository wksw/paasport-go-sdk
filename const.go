package paasport

const (
	// HTTPHeaderAppId the app id name in request header
	HTTPHeaderAppId = "paasport-app-id"
	// HTTPHeaderSubAppId the sub app id name in request header
	HTTPHeaderSubAppId = "paasport-sub-app-id"
	// HTTPHeaderRegion the region name in request header
	HTTPHeaderRegion = "paasport-region"
	// HTTPHeaderDeviceId the device id name in request header
	HTTPHeaderDeviceId = "paasport-device-id"
	// HTTPHeaderSubDeviceId the sub device id name in request header
	HTTPHeaderSubDeviceId = "paasport-sub-device-id"
	// HTTPHeaderTenant the tenant name in request header
	HTTPHeaderTenant = "paasport-tenant-name"
	// HTTPHeaderSubTenant the sub tenant name in request header
	HTTPHeaderSubTenant = "paasport-sub-tenant-name"
	// HTTPHeaderUserAgent the user agent name in request header
	HTTPHeaderUserAgent = "User-Agent"
	// HTTPHeaderToken the token name in request header
	HTTPHeaderToken = "x-auth-token"
	// HTTPHeaderSubToken the sub token name in request header
	HTTPHeaderSubToken = "x-sub-token"
	// HTTPHeaderTerminalType the terminal type name in request header
	HTTPHeaderTerminalType = "paasport-terminal-type"
	// HTTPHeaderContentType the content-type name in request header
	HTTPHeaderContentType = "content-type"
	// HTTPQueryTime the time name in request query
	HTTPQueryTime = "time"
	// HTTPQuerySignNonce the sign_nonce name in request query
	HTTPQuerySignNonce = "sign_nonce"
)

const (
	// SDK_VERSION sdk version
	SDK_VERSION = "v0.0.1"
	// PROJECT_NAME project name
	PROJECT_NAME = "PAASPORT"
	// DEFAULT_API_VERSION defautl api version
	DEFAULT_API_VERSION = "a1"
)
