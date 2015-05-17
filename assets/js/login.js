var centerLogin = function () {
	$('#login').css({
        'position' : 'fixed',
        'left' : '50%',
        'top' : '50%',
        'margin-left' : -$('#login').width() / 2,
        'margin-top' : -$('#login').height() / 2
    });
}

$(document).ready(function () {
	centerLogin();
	$(window).resize(centerLogin)
})

var app = angular.module('OtpApp', ['ngMaterial']);

app.controller('loginCtrl', ['$scope', function ($scope) {

	
}]);