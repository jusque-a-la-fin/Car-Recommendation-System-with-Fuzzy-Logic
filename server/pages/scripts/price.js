function isNumber(value) {
  return !isNaN(parseFloat(value)) && isFinite(value);
}

function validateInput(inputElement, errorElement) {
  let value = inputElement.value;
  if (value.startsWith('0')) {
    value = '';
  }

  if (isNumber(value) && value >= 0) {
      inputElement.classList.remove('error');
      errorElement.textContent = '';
      inputElement.placeholder = 'Самая низкая цена';
  } else {
      inputElement.classList.add('error');
      if (value < 0) {
          errorElement.textContent = 'Введите число больше нуля';
      } else {
          errorElement.textContent = 'Введите число больше нуля';
      }
      inputElement.value = ''; 
      inputElement.placeholder = ''; 
  }
}


const lstPriceInput = document.getElementById('lst_price');
const hstPriceInput = document.getElementById('hst_price');
const lstPriceError = document.getElementById('lst_price_error');
const hstPriceError = document.getElementById('hst_price_error');

lstPriceInput.addEventListener('input', function() {
  validateInput(lstPriceInput, lstPriceError);
});

hstPriceInput.addEventListener('input', function() {
  validateInput(hstPriceInput, hstPriceError);
});

 


function sendRequest() {
    var sessionID = sessionStorage.getItem('sessionID');
    if (sessionID) {
      const data = { sessionID: sessionID, minPrice: document.getElementById('lst_price').value, maxPrice: document.getElementById('hst_price').value};
      const url = 'http://localhost:8080/selection/price';
      fetch(url, {
          method: 'POST',
          headers: {
              'Content-Type': 'application/json'
          },
          body: JSON.stringify(data)
      })
      .then(response => {
          if (response.ok) {
              window.location.href = "http://localhost:8080/selection/manufacturers";
          } else {
              throw new Error('HTTP Error: ' + response.status);
          }
      })
      .catch(error => console.error(error));
  } else {
      console.log("Key 'sessionID' not found in sessionStorage");
  }

}

