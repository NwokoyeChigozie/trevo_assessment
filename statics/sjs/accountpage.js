function UserDetailsInterval() {
    if (ready == 1) {
        var tableData = "";
        console.log(userDetails);
        numberOfPositions = userDetails.portfolio_positions.length;
        document.querySelector(".fullname").innerHTML = userDetails.full_name;
        document.querySelector(".email").innerHTML = userDetails.email;
        // document.getElementById("edit_email").value = userDetails.email;
        document.querySelector("#number_of_positions").innerHTML = numberOfPositions;


        for (let i = 0; i < numberOfPositions; i++) {
            var equity_value = userDetails.portfolio_positions[i].equity_value;
            var price_per_share = userDetails.portfolio_positions[i].price_per_share;
            var symbol = userDetails.portfolio_positions[i].symbol;
            var total_quantity = userDetails.portfolio_positions[i].total_quantity;

            var row = "<tr><td>" + symbol + "</td><td>" + total_quantity + "</td><td>$" + price_per_share + "</td><td>$" + equity_value + "</td></tr>";
            tableData += row
            totalPositionsValue += equity_value;
            console.log(userDetails.portfolio_positions[i]);
        }
        document.querySelector("#total_portfolio_value").innerHTML = totalPositionsValue;
        document.querySelector("#positons_table").innerHTML = tableData;
        clearInterval(id);
    }
}
var id = setInterval(UserDetailsInterval, 10);

function LoanDetailsInterval() {
    if (loanready == 1) {
        var loan_value = 0;
        if (loanDetails === null) {
            loan_value = 0;
        } else {
            loan_value = loanDetails.balance;
        }

        document.querySelector("#total_loan_value").innerHTML = loan_value;
        clearInterval(loanInterval);
    }
}
var loanInterval = setInterval(LoanDetailsInterval, 10);