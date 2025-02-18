package network

import (
	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/handles/attorney"
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

type HandlesDB struct {
	TokenToHandle map[crypto.Token]UserInfo
	HandleToToken map[string]crypto.Token
	Attorneys     map[crypto.Token]struct{}
	SynergyApp    crypto.Token
}

func (a *HandlesDB) Handle(token crypto.Token) *state.UserInfo {
	if user, ok := a.TokenToHandle[token]; ok {
		return &state.UserInfo{
			Handle: user.Handle,
		}
	}
	return nil
}

func (a *HandlesDB) Token(handle string) *crypto.Token {
	if token, ok := a.HandleToToken[handle]; ok {
		return &token
	}
	return nil
}

func (a *HandlesDB) IncorporateJoin(action []byte) {
	join := attorney.ParseJoinNetwork(action)
	if join == nil {
		return
	}
	a.TokenToHandle[join.Author] = UserInfo{
		Handle:  join.Handle,
		Details: join.Details,
	}
	a.HandleToToken[join.Handle] = join.Author
}

func (a *HandlesDB) IncorporateUpdate(action []byte) {
	update := attorney.ParseUpdateInfo(action)
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

func (a *HandlesDB) IncorporateGrant(action []byte) {
	grant := attorney.ParseGrantPowerOfAttorney(action)
	if grant == nil {
		return
	}
	if grant.Attorney.Equal(a.SynergyApp) {
		a.Attorneys[grant.Author] = struct{}{}
	}
}

func (a *HandlesDB) IncorporateRevoke(action []byte) {
	revoke := attorney.ParseRevokePowerOfAttorney(action)
	if revoke == nil {
		return
	}
	if revoke.Attorney.Equal(a.SynergyApp) {
		delete(a.Attorneys, revoke.Author)
	}
}

func (a *HandlesDB) Incorporate(action []byte) []byte {
	switch attorney.Kind(action) {
	case attorney.VoidType:
		//if FilterSynergyProtocolCode(action) {
		//fmt.Println("Synergy protocol code")
		return action
		//}
	case attorney.JoinNetworkType:
		a.IncorporateJoin(action)
		return nil
	case attorney.UpdateInfoType:
		a.IncorporateUpdate(action)
		return nil
	case attorney.GrantPowerOfAttorneyType:
		a.IncorporateGrant(action)
		return nil
	case attorney.RevokePowerOfAttorneyType:
		a.IncorporateRevoke(action)
		return nil
	}
	return action
}
