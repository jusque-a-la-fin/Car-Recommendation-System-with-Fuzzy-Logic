var loading = document.getElementById('loader');

loading.style.display = "none";

function goToInternet() {
    hideAllAndShowLoading();
    var sessionID = sessionStorage.getItem('sessionID');
    if (sessionID) {
      const url = 'http://localhost:8080/selection/internet';
      fetch(url, {
          method: 'POST',
          headers: {
              'Content-Type': 'application/json'
          },
          body: JSON.stringify(sessionID)
      })
      .then(response => {
          if (response.ok) {
            window.location.href = "http://localhost:8080/selection/internet?guest="+sessionID;
          } else {
              throw new Error('HTTP Error: ' + response.status);
          }
      })
    }
}

function goToInternalDB() {
    hideAllAndShowLoading();
    var sessionID = sessionStorage.getItem('sessionID');
    if (sessionID) {
      const url = 'http://localhost:8080/selection/internal_db';
      fetch(url, {
          method: 'POST',
          headers: {
              'Content-Type': 'application/json'
          },
          body: JSON.stringify(sessionID)
      })
      .then(response => {
          if (response.ok) {
            window.location.href = "http://localhost:8080/selection/internal_db?guest="+sessionID;
          } else {
              throw new Error('HTTP Error: ' + response.status);
          }
      })
    }
}


function hideAllAndShowLoading() {
    var choice_box = document.getElementById('choice_box');
    choice_box.style.display = "none";
    var loading = document.getElementById("loader");
    loading.style.display = "block";
}

