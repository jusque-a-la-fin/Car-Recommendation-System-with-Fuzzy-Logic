
var preloader = document.getElementById("preloader");
preloader.style.display = "none";
var mark;

 const brandSelect = document.getElementById('mark');
 const modelSelect = document.getElementById('model');

 const initialModelOptions = modelSelect.innerHTML;


 brandSelect.addEventListener('change', function () {
   const selectedBrandId = brandSelect.options[brandSelect.selectedIndex].id;
   modelSelect.innerHTML = initialModelOptions;

   if (!selectedBrandId) {
	 modelSelect.disabled = true;
   } else {
	 modelSelect.disabled = false;
	 Array.from(modelSelect.options).forEach(function (option) {
	   const brandId = option.getAttribute('id'); 
	   if (brandId === selectedBrandId) {
		 option.style.display = 'block'; 
	   } else {
		 option.style.display = 'none'; 
	   }
	 });
   }
 });



const lowPriceInput = document.getElementById('low_price_limit');
const highPriceInput = document.getElementById('high_price_limit');

lowPriceInput.addEventListener('blur', () => {
const value = lowPriceInput.value.trim();

if (isNaN(value) || value < 0 || value[0] === '0') {
lowPriceInput.value = ''; 
lowPriceInput.classList.add('is-invalid');
document.getElementsByClassName('error low_price')[0].innerText = 'Введите число больше 0';
} else {
lowPriceInput.classList.remove('is-invalid');
document.getElementsByClassName('error low_price')[0].innerText = '';
}
});

highPriceInput.addEventListener('blur', () => {
const value = highPriceInput.value.trim();
if (isNaN(value) || value < 0 || value[0] === '0') {
highPriceInput.value = ''; 
highPriceInput.classList.add('is-invalid');
document.getElementsByClassName('error high_price')[0].innerText = 'Введите число больше 0';
} else {
highPriceInput.classList.remove('is-invalid');
document.getElementsByClassName('error high_price')[0].innerText = '';
}
});


function hideAllAndShowLoading() {
	var form = document.getElementById("usual_search")
	form.style.display = "none";
	var fuzzy_algorithm = document.getElementById("fuzzy_algorithm")
	fuzzy_algorithm.style.display = "none"
	preloader.style.display = "block";

    const formData = new FormData(form);
    let formDataObject = {};
formData.forEach(function(value, key){
    formDataObject[key] = value;
});

    var sessionID = sessionStorage.getItem('sessionID');
    if (sessionID) {
        const data = { sessionID: sessionID, form: formDataObject };
        const url = 'http://localhost:8080/main';
        fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        })
        .then(response => {
            if (response.ok) {
                window.location.href = "http://localhost:8080/search?guest="+sessionID;
            } else {
                throw new Error('HTTP Error: ' + response.status);
            }
        })
        .catch(error => console.error(error));
    } else {
        console.log("Key 'sessionID' not found in sessionStorage");
    }
}
