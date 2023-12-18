package pokt_v0

import "os-gateway/pkg/pokt/pokt_v0/models"

// getSessionFromRequest obtains a session from a relay request.
// Parameters:
//   - req: SendRelayRequest instance containing the relay request parameters.
//
// Returns:
//   - (*GetSessionResponse): Session response.
//   - (error): Error, if any.
func GetSessionFromRequest(pocketService PocketService, req *models.SendRelayRequest) (*models.Session, error) {
	if req.Session != nil {
		return req.Session, nil
	}
	sessionResp, err := pocketService.GetSession(&models.GetSessionRequest{
		AppPubKey: req.Signer.PublicKey,
		Chain:     req.Chain,
	})
	if err != nil {
		return nil, err
	}
	return sessionResp.Session, nil
}
