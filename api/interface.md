Páginas

backend

/api (POST method)

    actions
        AcceptCheckinEvent (separado)
        CreateBoard (separado)
        CreateCollective (separado)
        CreateEvent (separado)
        Draft (separado)
        Edit (separado)
        UpdateBoard (separado)
        UpdateCollective (separado)
        UpdateEvent (separado)
        
        ImprintStamp (incorporado)
        Pin (incorporado)
        React (incorporado)
        ReleaseDraft (incorporado)
        RemoveMember (incorporado)
        RequestMembership (incorporado)
        BoardEditor (incorporado)
        CancelEvent (incorporado)
        CheckinEvent (incorporado)
        
        Vote 

Templates:

/collectives 
    botões: create collective 
/collective/ 
    botões: create board, create event, 
    forms: react

/boards
    botões: create board
/board 
    botão: update board, 
    forms: unpin, board editor

/drafts 
    botões: new draft
/draft/ 
    botões: new version, edit, 
    forms: pin, stamp, react, release, 
    votes forms: authorship, pin, stamp, release

/events 
    botões: create event
/event/ 
    botões: update event, cancel event, accept checkin (editor only)
    form: check in
    votes forms: cancel event, create event
    votes links: update event 

/members
/member/

/votes
    votes forms:
        create board, request membership, remove member,  
    votes links:
        create event, update event, cancel event, 
        pin, stamp, authorship, release  


/static/...