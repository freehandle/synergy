package network

import (
	axe "github.com/freehandle/axe/attorney"
	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/synergy/social/state"
)

// Synergy protocol code in binary is
// 00000000 00000000 00000001 00000010
func FilterSynergyProtocolCode(action []byte) bool {
	if len(action) < 14 {
		return false
	}
	return action[10] == 0 && action[11] == 0 && action[12] == 1 && action[13] == 2
}

type UserInfo struct {
	Handle  string
	Details string
}

type AxeDB struct {
	TokenToHandle map[crypto.Token]UserInfo
	HandleToToken map[string]crypto.Token
	Attorneys     map[crypto.Token]struct{}
	SynergyApp    crypto.Token
}

func (a *AxeDB) Handle(token crypto.Token) *state.UserInfo {
	if user, ok := a.TokenToHandle[token]; ok {
		return &state.UserInfo{
			Handle: user.Handle,
		}
	}
	return nil
}

func (a *AxeDB) Token(handle string) *crypto.Token {
	if token, ok := a.HandleToToken[handle]; ok {
		return &token
	}
	return nil
}

func (a *AxeDB) IncorporateJoin(action []byte) {
	join := axe.ParseJoinNetwork(action)
	if join == nil {
		return
	}
	a.TokenToHandle[join.Author] = UserInfo{
		Handle:  join.Handle,
		Details: join.Details,
	}
	a.HandleToToken[join.Handle] = join.Author
}

func (a *AxeDB) IncorporateUpdate(action []byte) {
	update := axe.ParseUpdateInfo(action)
	if update == nil {
		return
	}
	handle, ok := a.TokenToHandle[update.Author]
	if !ok {
		return
	}
	a.TokenToHandle[update.Author] = UserInfo{
		Handle:  handle.Handle,
		Details: update.Details,
	}
}

func (a *AxeDB) IncorporateGrant(action []byte) {
	grant := axe.ParseGrantPowerOfAttorney(action)
	if grant == nil {
		return
	}
	if grant.Attorney.Equal(a.SynergyApp) {
		a.Attorneys[grant.Author] = struct{}{}
	}
}

func (a *AxeDB) IncorporateRevoke(action []byte) {
	revoke := axe.ParseRevokePowerOfAttorney(action)
	if revoke == nil {
		return
	}
	if revoke.Attorney.Equal(a.SynergyApp) {
		delete(a.Attorneys, revoke.Author)
	}
}

func (a *AxeDB) Incorporate(action []byte) []byte {
	switch axe.Kind(action) {
	case axe.VoidType:
		//if FilterSynergyProtocolCode(action) {
		return action
		//}
	case axe.JoinNetworkType:
		a.IncorporateJoin(action)
		return nil
	case axe.UpdateInfoType:
		a.IncorporateUpdate(action)
		return nil
	case axe.GrantPowerOfAttorneyType:
		a.IncorporateGrant(action)
		return nil
	case axe.RevokePowerOfAttorneyType:
		a.IncorporateRevoke(action)
		return nil
	}
	return action
}
