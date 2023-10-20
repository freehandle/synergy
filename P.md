This is the official implementation of the ***Synergy Protocol*** by the 
***Synergy Protocol Collective***. In order to learn how to use it and deploy a 
running instance of the interface click 
*[here](https://github.com/freehandle/synergy/blob/main/how.md)*

## Synergy Protocol Phisolophy

Throughout history, there have been numerous instances where individuals have 
assembled together, united by a common cause, to challenge the status quo. They 
invented new media for coordination of action and sharing of ideas. One such 
example is the Republic of Letters, a collective of intellectuals in the 16th 
and 17th centuries who revolted against the dominant scholasticism of the time.  
By engaging in an extensive correspondence network, they sought to promote the 
exchange of ideas and foster intellectual freedom. In due time it converged into
scientific societies, their periodicals, and the birth of the scientific 
revolution.

Similarly, the emergence of the Internet Task Force demonstrated a collective 
effort to combat the dominance of closed or proprietary standards in the digital 
communication realm. This group, thorugh the RFC mechanism, championed open 
standards, advocating for a more inclusive and accessible internet. Furthermore, 
the Riot Grrrl zine movement in the 1990s played a crucial role in shaping the 
ideas behind the third wave of feminism. By using self-published zines, they 
challenged traditional gender roles and provided a platform for common people 
voices to be heard. These examples illustrate the power of collective action and 
how people coming together can challenge established norms and drive societal 
change.

All of these share a common theme: experimentalism rather than dogma, reputation 
rather than authority, autonomy rather than control. 

In the present era, there is an urgent need for individuals to assemble together 
and reconsider the role of the internet, social media, and technology platforms 
at large. With the growing concerns over privacy breaches, misinformation, 
censorship and algorithmic biases, it has become imperative to experiment with 
technologies that empower autonomous and free digital interactions between 
individuals. By fostering interdisciplinary collaborations among technologists, 
intelectuals, activists, and users, we can work towards building transparent, 
inclusive, and user-centric digital environments. Assembling together now, we 
can harness our collective wisdom to shape the future of the internet and ensure 
that it remains a force for positive societal transformation, preserving the 
values of autonomy, freedom, and equality in the digital realm.

We strongly believe that any social media worthy of its name must be invented in 
itself, and must be invented socially. We cannot continue to accept that the 
terms of our digital experince are to be dictate by a tiny group of individuals 
with their peculiar cultural biases, worldview and motivations. We all must build 
a new, personal, internet!

The purpose of a social protocol is to define a pool of actions that can be 
performed by individuals or groups of individuals and lay down the rules 
governing those actions. 

Consider twitter for example. The basic pool of actions are: *tweet*, *reply*, 
*retweet*, *like* and *follow*. There are minimal rules governing those actions.
Users might, for example, restrict reply funcionality only to people followed by
them. 

It is a firm belief behind ***Synergy Protocol*** underlying philosophy that one of
the key reasons for the toxic environment within social media is the lack of 
experimentation of alternative forms of social protocols. The modes of 
interaction that we have grown accostumed to are hardly functional. We must 
begin to try different modes of interaction, different interfaces, mutiple ranking
algorithms and so on. 

In this sense, ***Synergy Protocol*** is a call for action. Let´s build a 
preliminary functionality 


## Synergy Architecture

***Synergy Protocol*** is conceived as a strategy to bootstrap a new personal internet 
based on user autonomy. It runs on top of axé social protocol, that might be
conceived as some sort of TCP/IP for humans. Axé itself runs on top of breeze, 
an innovative crypto network designed to scale at orders close to those required
by a global social network. 

The architecture of the whole system might be understood as:

***Breeze*** provides a decentralized gateway to submit individual social actions
that implements one or more social protocols. Beeing a crypto network it offers
two basic funcionalities: an approximate proof-of-timestamp (reliable on the 
order of a few seconds) and a global consensual ordering of actions received by 
different nodes of the network. 

Besides processing instructions that transfer or stake its fungible token,
that governs the economics of the network, ***Breeze*** has a single general purpose 
void instruction. It is used to instruct actions of social protocols. 

One of such social protocols is the ***Axé*** social protocol. It offers a general 
purpose identity management by associating a crypto key to a human readable
unique handle. For example, @synergy handle can be associated to a assymetric
cryptography key pair. With this functionality anyone in posession of that key
can sign messages on behalf of the handle. Because it is naïve to expect that 
end users can be trusted to safekeep and use judiciously cryptographic key ***Axé***
also implements a functionality to appoint other keys as attorneys of the 
identity key. 

***Axé*** is not conceived as an end-use protocol, but only as a layer to provide 
proof-of-authorship funcionality. Besides several actions to manage that 
functionality it also offers an unique general purpose authored void action. 
This void action can be used for other social protocols that leverages on 
***Breeze*** proof-of-timestamp and action ordering and on ***Axé*** proof-of-timestamp. 
Anyone listening to a relible source of ***Axé*** validated blocks full of void 
actions will have a global consensual sequence of authored actions to be 
processed by more specialized protocols. 

One such protocol is the ***Aiyeh***  , that provides functionaty of digital
stages. These can be considered as managed digital venues for social 
interaction. Co-Owners of those stages can define them as public or private, 
open or moderated. And define which social protocols it can be used for, among
other things. 

Another such protocol is the ***Synergy Protocol***. 

