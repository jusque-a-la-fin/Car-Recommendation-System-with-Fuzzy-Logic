var btns = document.getElementsByClassName("answer");

for (var i = 0; i < btns.length; i++) {
    btns[i].addEventListener("click", changeStyleWhenClick);
}

function changeStyleWhenClick() {
    var parentElement = this.parentElement;
    var previousElement = document.querySelector('.clicked');


    if (this.classList.length <= 2) {
    this.classList.add("clicked");

  } else {
    this.classList.remove("clicked");

  }
}


var choices = [];
var clicks = [];

function changeChoice(choice, index) {
    if (clicks[index] == undefined) {
      clicks[index] = true;
    
      choices.push(choice);
    }
    else {
      clicks[index] = undefined;
      let indexOfChoicesElement = choices.indexOf(choice);
      choices.splice(indexOfChoicesElement, 1);
    }
}


function sendRequest() {
var sessionID = sessionStorage.getItem('sessionID');
if (sessionID) {
  const data = { sessionID: sessionID, manufacturers: choices};
  const url = 'http://localhost:8080/selection/manufacturers';
  fetch(url, {
      method: 'POST',
      headers: {
          'Content-Type': 'application/json'
      },
      body: JSON.stringify(data)
  })
  .then(response => {
      if (response.ok) {
          window.location.href = "http://localhost:8080/selection/choice";
      } else {
          throw new Error('HTTP Error: ' + response.status);
      }
  })
  .catch(error => console.error(error));
} else {
  console.log("Key 'sessionID' not found in sessionStorage");
}
}

