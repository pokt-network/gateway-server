package altruist_registry

type AltruistRegistryService interface {
	GetAltruistURL(chainId string) (string, bool)
}
