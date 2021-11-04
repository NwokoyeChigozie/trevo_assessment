function UserDetailsInterval() {
    if (ready == 1) {
        var tableData = "";
        console.log(userDetails);
        numberOfPositions = userDetails.portfolio_positions.length;
        // document.querySelector(".email").innerHTML = userDetails.email;
        document.getElementById("edit_full_name").value = echoFieldvalue(userDetails.full_name);
        document.getElementById("edit_email").value = echoFieldvalue(userDetails.email);
        document.getElementById("edit_phone").value = echoFieldvalue(userDetails.phone);

        clearInterval(id);
    }
}
var id = setInterval(UserDetailsInterval, 10);

function echoFieldvalue(value) {
    if (value !== undefined) {
        return value;
    } else {
        return "";
    }
}