var count = 0;

var option = document.getElementById("option");
var options = document.getElementById("options");
var password1 = document.getElementById("password");
var password2 = document.getElementById("password2")

function addOption(){
       if (option.value == "") {
           return
       }
       if( count >= 2 ){
            option.required = false;
       }
       appendOption();
       option.value = "";
       count ++
}

function submitNewPoll(){
    if (count < 2) {
        alert("PLEASE ADD AT LEAST TWO OPTIONS!")
        return
    }
    option.value = "";
    option.value = options.innerText;
    document.getElementById("poll-form").submit();
}

function submitNewUser(){
    if (password1.value.length < 6) {
        alert("PASSWORD MUST BE AT LEAST 6 CHARACTERES LONG...")
        return
    }else if (password1.value !== password2.value) {
        alert("PASSWORD DID NOT MATCH...")
        return
    }
    document.getElementById("new-user-form").submit();
}
    
function appendOption() {  
    var newOption = `<div class='alert alert-info my-3'>
    <button type='button' class='close' 
    data-dismiss='alert'>&times;</button>` + option.value + `</div>`;
    options.insertAdjacentHTML("beforeend", newOption);
}