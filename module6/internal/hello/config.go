package hello

type Config struct {
	MethodConfig        []*MethodConfig        `json:"methodConfig"`
	LoadBalancingConfig []*LoadBalancingConfig `json:"loadBalancingConfig"`
}

type MethodConfig struct {
	Names       []*NameConfig `json:"name"`
	Timeout     string        `json:"timeout"`
	RetryPolicy *RetryPolicy  `json:"retryPolicy"`
}

type NameConfig struct {
	Service string `json:"service"`
	Method  string `json:"method"`
}

type RetryPolicy struct {
	MaxAttempts          int      `json:"maxAttempts"`
	InitialBackoff       string   `json:"initialBackoff"`
	MaxBackoff           string   `json:"maxBackoff"`
	BackoffMultiplier    int      `json:"backoffMultiplier"`
	RetryableStatusCodes []string `json:"retryableStatusCodes"`
}

type LoadBalancingConfig struct {
	RoundRobin string `json:"round_robin"`
}
