WALLET command provides basic functionality for self-hosted key management for 
the axe social protocol.


wallet create-user user.json

user.json must be of the form

{
    "handle": "user-handle",
    "


}

on creation the user will be asked to provide a passphrase for encryption 

wallet show-secret [all|kid]

wallet show-keys

wallet grant-attorney attorney_token

wallet revoke-attorney attorney_token

wallet create-stage stage.json

stage.json

{
    "name": "stage name",
    "description": "stage description",
    "public": bool,
    "open": bool,
    "moderated": bool
}

wallet show-own-stages

wallet show-stages

drum update-user user.json
drum create-user user.json
drum grant-attorney attorney-token
drum revoke-attorney attorney-token
drum show-attorneys
drum create-stage stage.json
drum update-stage stage.json
drum rotate-keys rotation.json
drum show-stage-requests
drum aprove-stage-request #
drum post post.json
drum show-stages
drum moderation


