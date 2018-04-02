let count = 0;

const option = document.getElementById("option");
const options = document.getElementById("options");
const password1 = document.getElementById("password");
const password2 = document.getElementById("password2");
const name = document.getElementById("name");
const email = document.getElementById("email");

function addOption(){
       if (option.value == "") {
           return;
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
        return;
    }
    option.value = "";
    option.value = options.innerText;
    document.getElementById("poll-form").submit();
}

function submitNewUser(){
    var emailReg = /\S+@\S+\.\S+/;
    if ( name.value.length < 3 ){
        alert("NAME MUST BE AT LEAST 3 CHARACTERES LONG !!!");
        return;
    } else if (!emailReg.test(email.value)) {
        alert("PLEASE ENTER VALID EMAIL ADDRESS !!!");
        return;
    } else if (password1.value.length < 3) {
        alert("PASSWORD MUST BE AT LEAST 3 CHARACTERES LONG !!!")
        return;
    } else if (password1.value !== password2.value) {
        alert("PASSWORD DID NOT MATCH...")
        return;
    }
    document.getElementById("new-user-form").submit();
}
    
function appendOption() {  
    var newOption = `<div class='alert alert-info my-3'>
    <button type='button' class='close' 
    data-dismiss='alert'>&times;</button>` + option.value + `</div>`;
    options.insertAdjacentHTML("beforeend", newOption);
}