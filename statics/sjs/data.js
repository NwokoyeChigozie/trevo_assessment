var userDetails;
var loanDetails;
var ready = 0;
var loanready = 0;
var totalPositionsValue = 0;
var totalLoansValue = 0;
var numberOfPositions = 0;
var token = "";
var customerEmail = "";
var customerName = "";

function navigateToLogin() {
    window.location = '/login';
}
$(document).ready(function() {
    var claims = localStorage.getItem('claims')
    if (claims === null) {
        navigateToLogin();
    } else {
        token = claims;
        var settings = {
            "url": "/v1/verify-token",
            "method": "POST",
            "timeout": 0,
            "headers": {
                "Authorization": "Bearer " + claims
            },
        };

        $.ajax(settings).done(function(response) {
            if (response.status == 200) {
                var settings = {
                    "url": "/v1/user",
                    "method": "GET",
                    "timeout": 0,
                    "headers": {
                        "Authorization": "Bearer " + claims
                    },
                };

                var settings1 = {
                    "url": "/v1/loan",
                    "method": "GET",
                    "timeout": 0,
                    "headers": {
                        "Authorization": "Bearer " + claims
                    },
                };

                $.ajax(settings).done(function(response) {
                    // console.log(response);
                    userDetails = response.data
                    ready = 1;
                });

                $.ajax(settings1).done(function(response) {
                    console.log(response);
                    loanDetails = response.data
                    loanready += 1;
                });

            } else {
                localStorage.clear();
                navigateToLogin();
            }
        });

    }



});