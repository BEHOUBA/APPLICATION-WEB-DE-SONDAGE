var count = 0;


function addOption(){
       var option = document.getElementById("option");
       if (option.value == "") {
           return
       }
       if( count >= 2 ){
            option.required = false;
       }
       appendOption();
       var allOptions = document.getElementById("all-options");
       allOptions.value += " " + option.value;
       option.value = "";
       count ++
    }


    
function appendOption() {  
    var optionsField = document.getElementById("options");
    var newOption = `<div class='alert alert-info my-3'>
    <button type='button' class='close' 
    data-dismiss='alert'>&times;</button>` + option.value + `</div>`;
    optionsField.insertAdjacentHTML("beforeend", newOption);
}