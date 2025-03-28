// console.log("estou aqui")

window.onload = function () {

  // input fields detailing
  let elementsdetailed = document.getElementsByClassName("detailed");
  if (elementsdetailed) {
    for (let i=0; i<elementsdetailed.length; i++) {
      let el = elementsdetailed[i];
      let id = el.getAttribute("id");
      el.addEventListener("focusin", displayinfo(id+"info"));
      el.addEventListener("focusout", hideinfo(id+"info"));
    }
  }
}

function selectFile() {
  let filename = document.getElementById('fileudraft').value;
  document.getElementById('fileName').setAttribute('value', filename);
}


// form info functions

function displayinfo(id) {
  return () => {
    let el = document.getElementById(id);
    el.classList.remove("fieldinfohide");
    el.classList.add("fieldinfoshow");
  }
}

function hideinfo(id) {
  return () => {
    let el = document.getElementById(id)
    el.classList.remove("fieldinfoshow");
    el.classList.add("fieldinfohide");
  }
}

// MODAL functions

function closedialog(id) {
  let el = document.getElementById(id);
  el.close();
}

// react 

function dialogreact() {
  // shows dialog element
  let el = document.getElementById("dialogreactel");
  el.showModal();

  // gets reaction from main page into dialog
  let reaction = document.getElementById("reaction");
  let reactionmodal = document.getElementById("reactionmodal");
  reactionmodal.checked = reaction.checked;
  
  // gets outline paragraph to be shown in modal
  let reactpar = document.getElementById("reactionoutline");
  let pagename = document.getElementById("modaloutlinename").innerHTML;

  if (reactionmodal.checked) {
    reactpar.innerHTML = "like "+ pagename;
  } else {
    reactpar.innerHTML = "dislike " + pagename;
  }
}

// leave collective

function dialogleavecollective() {
  // shows dialog element
  let el = document.getElementById("dialogleavecollectiveel");
  el.showModal();  
  
  // gets outline paragraph to be shown in modal
  let leavepar = document.getElementById("leaveoutline");
  let pagename = document.getElementById("modaloutlinename").innerHTML;

  leavepar.innerHTML = "leave "+ pagename + " collective";
}

// join collective

function dialogjoincollective() {
  // shows dialog element
  let el = document.getElementById("dialogjoincollectiveel");
  el.showModal();
  
  // gets outline paragraph to be shown in modal
  let joinpar = document.getElementById("joinoutline");
  let pagename = document.getElementById("modaloutlinename").innerHTML;

  joinpar.innerHTML = "join "+ pagename + " collective";
}

// remove editor from board

function dialogremoveeditor() {
  // shows dialog element
  let el = document.getElementById("dialogremoveeditorel");
  el.showModal();

  // gets editor handle from main page into dialog
  let editor = document.getElementById("editorhandle");
  let editormodal = document.getElementById("modaleditorhandle");
  if (editor) {
    editormodal.value = editor.value;
  }
    
  // gets outline paragraph to be shown in modal
  let removepar = document.getElementById("removeeditoroutline");
  let pagename = document.getElementById("modaloutlinename").innerHTML;

  removepar.innerHTML = "remove editor " + editormodal.value + " from " + pagename + " board";
}

// apply for board editor

function dialogapplyeditor() {
  // shows dialog element
  let el = document.getElementById("dialogapplyeditorel");
  el.showModal(); 
    
  // gets outline paragraph to be shown in modal
  let applypar = document.getElementById("applyeditoroutline");
  let pagename = document.getElementById("modaloutlinename").innerHTML;

  applypar.innerHTML = "apply for editorship of " + pagename + " board";
}

// pin to board

function dialogpintoboard() {
  // shows dialog element
  let el = document.getElementById("dialogpinboardel");
  el.showModal();

  // gets editor handle from main page into dialog
  let board = document.getElementById("boardname");
  let boardmodal = document.getElementById("modalboardname");
  if (board) {
    boardmodal.value = board.value;
  }
    
  // gets outline paragraph to be shown in modal
  let pinpar = document.getElementById("pinboardoutline");
  let pagename = document.getElementById("modaloutlinename").innerHTML;

  pinpar.innerHTML = "pin draft " + pagename + " on " + board.value + " board";
}

// propose stamp

function dialogproposestamp() {
  // shows dialog element
  let el = document.getElementById("dialogproposestampel");
  el.showModal();

  // gets editor handle from main page into dialog
  let collectiverep = document.getElementById("collectiverep");
  let collectiverepmodal = document.getElementById("modalcollectiverep");
  if (collectiverep) {
    collectiverepmodal.value = collectiverep.value;
  }
    
  // gets outline paragraph to be shown in modal
  let stamppar = document.getElementById("propstampoutline");
  let pagename = document.getElementById("modaloutlinename").innerHTML;

  stamppar.innerHTML = "propose " + collectiverep.value + " stamp on " + pagename + " draft";
}

// propose release

function dialogrelease() {
  // shows dialog element
  let el = document.getElementById("dialogreleaseel");
  el.showModal(); 
    
  // gets outline paragraph to be shown in modal
  let releasepar = document.getElementById("releaseoutline");
  let pagename = document.getElementById("modaloutlinename").innerHTML;

  releasepar.innerHTML = "apply for release of " + pagename + " draft";
}


function selectGrammar(view) {
  let actions = document.getElementById("actionsview");
  let descriptions = actions.getElementsByClassName("description");
  const nounview = document.getElementById("nounview")
  const verbview = document.getElementById("verbview")
  if (view === "noun") {
    if (verbview) {
      verbview.classList.remove("selected");
    }
    if (nounview) {
      nounview.classList.add("selected");
    }
    for (const des of descriptions) {
      des.classList.add("noun");
      des.classList.remove("verb");
    }
  } else {
    if (verbview) {
      verbview.classList.add("selected");
    }
    if (nounview) {
      nounview.classList.remove("selected");
    }
    for (const des of descriptions) {
      des.classList.add("verb");
      des.classList.remove("noun");
    }
  }
}

function selectActionKind(kind) {
  const actions = document.getElementById("actionsview");
  const elements = actions.getElementsByClassName("item");
  let lastduration = ""
  for (const el of elements) {    
    if ((kind !== "all") && (!el.classList.contains(kind))) {
      el.classList.add("hide");  
    } else {
      el.classList.remove("hide");  
      let duration = el.getElementsByClassName("duration")[0]
      if (duration.textContent !== lastduration) {
        duration.classList.remove("hideduration");
        lastduration = duration.textContent;
      } else {
        duration.classList.add("hideduration");
      }
    }
  }
  const menu = document.getElementById('selectmenu')
  if (menu) {
    for (el of menu.getElementsByTagName('span')) {
      if (el.getAttribute('id') == 'select_'+kind) {
        el.classList.add('bold')
      } else {
        el.classList.remove('bold')
      }
    }

  }

}

function selectMedia(media) {
  const mydrafts = document.getElementById("mymediadrafts");
  const mydraftsmenu = document.getElementById("mymediadraftsmenu");
  const myedits = document.getElementById("mymediaedits");
  const myeditsmenu = document.getElementById("mymediaeditsmenu");
  if (media === "Draft") {
    mydrafts.classList.remove("none");
    myedits.classList.add("none");
    mydraftsmenu.classList.add("bold");
    myeditsmenu.classList.remove("bold");
  } else {
    mydrafts.classList.add("none");
    myedits.classList.remove("none");
    mydraftsmenu.classList.remove("bold");
    myeditsmenu.classList.add("bold");
  }
}

function selectEventView(view) {
  const attendingView = document.getElementById("attendingView");
  const managingView = document.getElementById("managingView");
  const managing = document.getElementById("managingbulk");
  const attending = document.getElementById("attendingbulk");
  const stats=document.getElementById("statsattending");
  if (view === "attending") {
    attending.classList.remove("hideevent");
    managing.classList.add("hideevent");
    attendingView.classList.add("selected");
    managingView.classList.remove("selected");
    stats.classList.remove("hideevent");
  } else {
    attending.classList.add("hideevent");
    managing.classList.remove("hideevent");
    attendingView.classList.remove("selected");
    managingView.classList.add("selected");
    stats.classList.add("hideevent");
  }
}

function selectToggle(view) {
  const toggle = document.getElementsByClassName("toggle");
  if (toggle) {
    for (el of toggle) {
      if (el.getAttribute('id') === view) {
        el.classList.remove("none"); 
      } else {
        el.classList.add("none");
      }
    }
  } 
  const menu =document.getElementsByClassName("tgmenu");
  if (menu) {
    for (el of menu) {
      if (el.getAttribute('id') === 'tg_'+view) {
        el.classList.add("bold"); 
      } else {
        el.classList.remove("bold");
      }
    }
  } 
}  
