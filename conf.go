package paasport

import (
	"time"
)

// Config paasport configuration
type Config struct {
	// gateway endpoint
	endpoint     string
	ak           string
	sk           string
	headers      map[string]string
	onebox       bool
	ihc          bool
	apiVersion   string
	tenantName   string
	region       string
	appId        int64
	deviceId     string
	userAgent    string
	proxyHost    string
	proxyUser    string
	proxyPwd     string
	HTTPTimeout  HTTPTimeout
	HTTPMaxConns HTTPMaxConns
	logger       Logger
	sslVerify    bool
}

// HTTPTimeout http timeout define
type HTTPTimeout struct {
	ConnectTimeout   time.Duration
	ReadWriteTimeout time.Duration
	HeaderTimeout    time.Duration
	LongTimeout      time.Duration
	IdleConnTimeout  time.Duration
}

// HTTPMaxConns max idle connections
type HTTPMaxConns struct {
	MaxIdleConns        int
	MaxIdleConnsPerHost int
}

// Configurer 配置
type Configurer func(conf *Config)

func (c *Config) withDefault() {
	if c.apiVersion == "" {
		c.apiVersion = DEFAULT_API_VERSION
	}
	if c.HTTPTimeout.ReadWriteTimeout == 0 {
		c.HTTPTimeout.ReadWriteTimeout = 5 * time.Second
	}
}

// WithSslVerify set sslVerify
func WithSslVerify(sslVerify bool) Configurer {
	return func(conf *Config) {
		conf.sslVerify = sslVerify
	}
}

// WithLogger set logger
func WithLogger(logger Logger) Configurer {
	return func(conf *Config) {
		conf.logger = logger
	}
}

// WithTimeout set timeout
func WithTimeout(httpTimeout HTTPTimeout) Configurer {
	return func(conf *Config) {
		conf.HTTPTimeout = httpTimeout
	}
}

// WithProxy set proxy
func WithProxy(proxyHost, proxyUser, proxyPwd string) Configurer {
	return func(conf *Config) {
		conf.proxyHost = proxyHost
		conf.proxyUser = proxyUser
		conf.proxyPwd = proxyPwd
	}
}

// WithMaxConnections set max connections
func WithMaxConnections(httpMaxConns HTTPMaxConns) Configurer {
	return func(conf *Config) {
		conf.HTTPMaxConns = httpMaxConns
	}
}

// WithAppId set appId
func WithAppId(appId int64) Configurer {
	return func(conf *Config) {
		conf.appId = appId
	}
}

// WithDeviceId set deviceId
func WithDeviceId(deviceId string) Configurer {
	return func(conf *Config) {
		conf.deviceId = deviceId
	}
}

// WithTenantName set tenant_name
func WithTenantName(tenantName string) Configurer {
	return func(conf *Config) {
		conf.tenantName = tenantName
	}
}

// WithRegion set region
func WithRegion(region string) Configurer {
	return func(conf *Config) {
		conf.region = region
	}
}

// WithApiVersion set api version
func WithApiVersion(apiVersion string) Configurer {
	return func(conf *Config) {
		conf.apiVersion = apiVersion
	}
}

// WithOnebox set onebox flag
func WithOnebox(onebox bool) Configurer {
	return func(conf *Config) {
		conf.onebox = onebox
	}
}

// WithIhc set ihc flag
func WithIhc(ihc bool) Configurer {
	return func(conf *Config) {
		conf.ihc = ihc
	}
}
