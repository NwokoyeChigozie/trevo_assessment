function UserDetailsInterval() {
    if (ready == 1) {
        numberOfPositions = userDetails.portfolio_positions.length;
        document.querySelector(".fullname").innerHTML = userDetails.full_name;
        customerEmail = userDetails.email;
        customerName = userDetails.full_name;

        for (let i = 0; i < numberOfPositions; i++) {
            var equity_value = userDetails.portfolio_positions[i].equity_value;
            totalPositionsValue += equity_value;
        }
        document.querySelector(".total_portfolio_value").innerHTML = totalPositionsValue;
        clearInterval(id);
    }
}
var id = setInterval(UserDetailsInterval, 10);

function LoanDetailsInterval() {
    if (loanready == 1) {
        var loan_value = 0;
        var loan_balance = 0;
        if (loanDetails === null) {
            loan_value = 0;
        } else {

            loan_balance = loanDetails.balance;
            loan_value = loanDetails.total_amount;
        }


        if (loan_balance <= 0) {
            document.getElementById("payback_section").style.display = "none";
        }

        document.querySelector("#due_amount").innerHTML = loanDetails.actualDue;
        document.querySelector("#due_next_info").innerHTML = '$<span id="amount_due_next" class="st_amount">' + loanDetails.amountNextDue + '</span> will be Due on ' + loanDetails.whenNextDue;
        document.querySelector(".total_loan_value").innerHTML = loan_value;
        document.querySelector("#loan_balance").innerHTML = loan_balance;
        clearInterval(loanInterval);
    }
}
var loanInterval = setInterval(LoanDetailsInterval, 10);