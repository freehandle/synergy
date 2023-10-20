From time to time groups of worried and bothered individuals take action for 
themselves to establish new media for the coordination of action and sharing ideas on a 
commitment to face the status quo.

Examples abound, of greater or lesser consequence. For example, the so called 
Republic of Letters gave birth, through the culmination of scientific societies
and their periodicals, to the scientific revolution. ARPANET gave birth, through 
the internet task force and the RFC publication, to the open standards internet.
More recently, riot grrrl gave birth, through their zine, to the third wave
of feminism. All of these share a common theme: experimentalism rather than dogma,
reputation rather than authority, autonomy rather than control.  

We believe that time has come for the worried and bothered of our era to stand 
against technology plataforms. And we are proposing a new media, the Synergy 
Social Protocol, in order to facilitate coordination of action and sharing of 
ideas.

We strongly believe that any social media worthy of its name must be 
invented in itself, and must be invented socially. We cannot continue to accept
that the terms of our digital experince are to be dictate by a tiny group of 
individuals with their peculiar cultural biases, worldview and motivations. 
We all must build a new, personal, internet!

# Synergy Protocol

## Overview

Synergy protocol was designed as a digital framework for collaboration and
collective construction. Protocol functionalities were inspired by the - arguably - most successful form of human cooperation: scientific publishing.

Synergy was built on top of axé social protocol, which provides primitives
for identity and stage management.

## General structure

The most fundamental functionalities for collaborative project development 
derive from two basic constructions: groups of individuals with a common goal,
which are called COLLECTIVE, and content creation dynamics, that this protocol
enables by the evolution of DRAFTs. 

A COLLECTIVE's ability to perform actions as a single entity makes it easy for
groups of people to act as a unity. For every action taken on behalf of a 
COLLECTIVE, the protocol automatically triggers a voting mechanism.

A DRAFT's evolution is acconted for by its EDIT, RELEASE and STAMP, this last action
is performed in the sense of a peer review, and can be taken on behalf of a COLLECTIVE.

## Collective actions

Groups of people can form the so called COLLECTIVE. They can perform any possible actions within the protocol as a unit. The only action that can be performed  exclusively on behalf of a COLLECTIVE is the creation of an EVENT.

### COLLECTIVE

COLLECTIVE are an association of members that have a common goal. Each
COLLECTIVE is created with a name that is unique within the network and
must not be a handle (cannot be an @), and with a description of its goal. 

Besides the name and description, upon being created a COLLECTIVE must provide 
its choice policy for pooling, which will dictate the voting mechanism 
for the actions that will be performed in its name. It must also provide a "super" 
policty for policy update, meaning the policy for changing the COLLECTIVE voting policy.

Every action taken on behalf of the collective triggers a voting mechanism
according to the COLLECTIVE policy. Whenever an instruction is posted on behalf 
of a COLLECTIVE the pooling mechanism is automatically triggered and its result 
is accepted as the decision of the COLLECTIVE.

If pool results in acceptance of the instruction, instruction in considered valid
and on behalf of the pointed COLLECTIVE. If an instruction is posted on behalf of 
a COLLECTIVE and the instruction’s pool results in non-acceptance of the instruction,
the instruction is discarded. 

The pooling mechanism of a COLLECTIVE may be updated.

Any of the actions prescribed by the protocol can be submitted as
being on behalf of a COLLECTIVE. That means COLLECTIVE can act as a
unit in the network. 

Everyone can apply to join a COLLECTIVE. Upon applying, a voting amongst COLLECTIVE 
members is taken to either accept or deny the request.

### EVENT

EVENTs can only be proposed on behalf of a COLLECTIVE. They are a means for members to interact for a period of time with a previously defined goal. EVENTs can be either public of private and can happen either on-chain or not. 

Upon creating an EVENT on behalf of the COLLECTIVE, the instruction author may specify a group of members as the EVENT's managers and these members will be able to accept participation requests. If managers are not appointed to an EVENT, all members of the COLLECTIVE can accept participation requests.


## Information dynamics

To account for basic elements for information exchange, the protocol provides
the following constructions.

### DRAFT

It is the basic element used for publishing ideas and contributions that are in
progress. DRAFTs are to be used for public idea elaboration and collaboration,
amongst network members.

By creating a DRAFT the member, or group of members, is sharing a starting proposition with
the community, and anyone who wishes to contribute to the forming of it can
apply for an EDIT to a DRAFT. New versions of the DRAFT may include the EDIT request, and DRAFT authors may include EDIT author as co-author to the new version. 
Once a new version of the DRAF is published, it references it's old version.

Any member can propose a new DRAFT as an individual contribution or
on behalf of a COLLECTIVE.

Besides the actual DRAFT content, DRAFT instructions must include a
title for easy identification (it does not need to be unique), a brief description
of the content, and a list of keywords the content is related to. It may, or
may not include a list of internal references used (content previously posted on
the network), if the DRAFT posted is a new version of a previously posted DRAFT, 
it must contain the hash of the previous DRAFT as its predecessor.

When a DRAFT is published on behalf of a COLLECTIVE, no co-authors can be named, even if the DRAFT includes alterations proposed by EDITs. Also, new version of it will automatically trigger COLLECTIVE policy for pooling it's acceptance by the COLLECTIVE. 

If the DRAFT was not created on behalf of a COLLECTIVE, it must provide a policy for updates. Should any new versions of this DRAFT be proposed, their acceptance will trigger a pool according to the policy specified at creation. Co-authors will need to approve new version's acceptance according to policy specified.

DRAFT instructions are necessarily public. All information published as a
DRAFT can be viewed and revised by the whole community.

### EDIT

The whole community has access to the DRAFTs created. Anyone can propose EDITs to existing DRAFTs. To propose an EDIT, the EDIT author references the DRAFT that's being edited. 

EDITs can be proposed on behalf of a COLLECTIVE. If so, the proposal will automatically trigger a pool according to the COLLECTIVE's policy. The EDIT instruction will only become valid once it is approved by the COLLECTIVE, otherwise it is discarded.

### RELEASE

If a DRAFT reaches a final form, it can be "promoted" as a RELEASE. To do so, any of it's authors (or anyone from the COLLECTIVE, should the DRAFT be on behalf of a COLLECTIVE) can create an instruction to point that DRAFT version as a RELEASE.

Acceptance of the instruction will follow either DRAFT policy, if written by a group of authors, or COLLECTIVE policy, if on behalf of a COLLECTIVE.

### STAMP

All RELEASEs can be stamped by a person or a COLLECTIVE. STAMPs are a way of endorsing content that is thought to be in accordance with the person or COLLECTIVE's criteria. They are also a way of promoting RELEASEs as being peer reviewed. 

That makes it easier for community members to know which content has been reviewed and is being
vouched by peers. Anyone can propose a STAMP on it's own behalf, or on behalf of a COLLECTIVE. 

When a STAMP is proposed on behalf of a COLLECTIVE, it automatically triggers the voting mechanism according to the COLLECTIVE policy.

### REACTION

To be used as both a ranking tool for content and a means
for the community to express its interest.

Members can react either positively or negatively to any instructions within the protocol.

REACTIONs can be either signed by a single member, as an individual, or
on behalf of a COLLECTIVE. If signed on behalf of a COLLECTIVE, a pool
according to the COLLECTIVE’s pooling mechanism is automatically triggered. 

The REACTION is only considered valid if the pool results in acceptance
of the instruction.

### BOARD

BOARDs are instructions that provide a keyword, or group of keywords,
to index DRAFTs that were posted referencing the BOARD’s dedicated keyword
or group of keywords.

BOARDs must have a name, a description and a dedicated policy for pin acceptance. They can be created either by a member or on behalf of a COLLECTIVE. If created on behalf of a COLLECTIVE, the BOARD may appoint editors, which are members of the COLLECTIVE that have power to propose the instruction to pin a DRAFT to the BOARD, and that are allowed to vote if proposed pins are to be accepted. 

Members chosen as editors can have their editing priviledges revoked by the COLLECTIVE. More editors can be appointed by the COLLECTIVE. 

Content indexed by a BOARD has not necessarily been reviewed by board’s editors.


