package access

// Effect type for policy effects
type Effect string

const (
	EffectAllow = Effect("allow")
	EffectDeny  = Effect("deny")
)
