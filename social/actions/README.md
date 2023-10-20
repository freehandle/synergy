# Actions in the Synergy Protocol


## Consensus

Synergy is built around the concept of collective consensus. This is encapsulated
on the Policy which defines a majority and a supermajority. Majority is a % of 
the numbers of participants required for consensus over ordinary actions and
supermajority is the % number of participants required for extraordinary actions
like change the policy itself. For example a collective of 10 individuals with
majority of 50 can take action if any 6 of those 10 individuals agree on the
action.

```
Policy {
    Majority        0-100 int 
    Supermajority   0-100 int
}

```

Every action within Synergy has a basic template

```
{
    Epoch           64bit uint
    Author          Token
    Reasons         string (optional)
}
```

Epoch and Author provides timestamp and authorship, as verified by the base
social protocol, and the reasons for the submission of any action to the 
protocol can be laid out in the Reasons field. Some actions might be invoked
not on the Author capacity but on Behalf of a collective, in which case, the 
template includes another field

```
{
    OnBehalfOf      string
}
```

In order to exert their role within the collective, individuals cast 

```
VoteAction {
	Epoch           64bit uint
	Author          Token
	Reasons         string (optional)
	Hash            Hash
	Approve         bool
}
```
These will be processed by the protocol and once consensus is achieved action 
is performed accordingly. 


## Collectives

OnBehalfOf must refer the name of a previously created collective. In order to
start a new collective an action must be provided

```
CreateCollectiveAction {
	Epoch           64bit uint
	Author          Token
	Reasons         string (optional)
	Name            string
	Description     string (optional)
	Policy          Policy   
}
```

Collectives cannot be created by other collectives. So there is no scope for a
OnBehalfOf field. The collective is obviously created with a single individual,
namely the Author of the instruction. 

Interested individuals might apply to become members of the collective by 
submiting and membership action

```
RequestMembershipAction {
	Epoch           64bit uint
	Author          Token
	Reasons         string (optional)
    Collective      string 
	Include         bool 
}
```

The same action is used to request membership (include = true) and to quit from
the collective (include = false).

The collective itself might decide to remove one of its members, in which case
a remove member action must be sent

```
RemoveMemberAction  {
	Epoch           64bit uint
	Author          Token
	Reasons         string (optional)
    OnBehalfOf      string
	Member          Token
}
```

Finally in order to update details about the collective, one might submit a

```
UpdateCollectiveAction {
	Epoch           64bit uint
	Author          Token
	Reasons         string (optional)
    OnBehalfOf      string
	Description     string (optional)
	Policy          Policy (optional)
}
```

## Draft

In order to submit a new draft to the protocol 

```
DraftAction  {
	Epoch           64bit uint
	Author          Token
	Reasons         string (optional)
    OnBehalfOf      string (optional)
	CoAuthors       []Token (optional)
	Policy          Policy (optional)
	Title           string
	Keywords        []string
	Description     string (optional)
	ContentType     string
	ContentHash     Hash 
	NumberOfParts   8bit uint
	Content         []byte 
	PreviousDraft   Hash (optional)
	References      []Hash (optional)
}
```

A draft can be authored by a single individual, a list of CoAuthors, or a named
collective. A draft action on behalf of a collective cannot contain coauthors
or policy. A default (unanimous) policy is adopted in case a non collective draft
is submitted without explicity policy.
Draft can be sent in multiple parts, in which case Number of Parts > 1. The 
following parts are submitted in 

```
MultipartMediaAction {
	Epoch           64bit uint
	Author          Token
    Hash            Hash
	Part            8bit uint
	Of              8bit uiny
	Data            []byte
}
```
The Hash provided in the DraftAction must be the hash of the entire content (the
concatenation of all the parts). The draft will only be valid after all the parts
are processed by the protocol and the hash matches.

## Board

